package handler

import (
	"dinosaur-emoji-service/internal/model"
	"dinosaur-emoji-service/pkg/database"
	"fmt"

	"encoding/json"

	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
)

type EmojiAddPackRequest struct {
	Name    string   `json:"name" binding:"required"`
	IconURL string   `json:"iconUrl"`
	Tags    []string `json:"tags"`
}

// EmojiPackAdd 创建表情包合集
// @Description 创建表情包合集
// @Tags emojiPack
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization"
// @Param emojiPack body EmojiAddPackRequest true "EmojiPack信息"
// @Success 200 {object} map[string]interface{} "{"code":200,"msg":"添加表情包成功","data":EmojiPack}"
// @Failure 400 {object} map[string]interface{} "{"code":400,"msg":"xxxx"}"
// @Router /emojiPack/add [post]
func EmojiPackAdd(c *gin.Context) {
	var params EmojiAddPackRequest
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(400, gin.H{"code": 400, "msg": "创建失败", "error": err.Error()})
		return
	}

	//处理参数
	temp, _ := json.Marshal(params.Tags)
	tagsTemp := datatypes.JSON(json.RawMessage(temp))

	EmojiPackInstance := model.EmojiPack{
		Name:             params.Name,
		IconURL:          params.IconURL,
		Tags:             tagsTemp,
		View_count:       0,
		Collection_count: 0,
		AuthorUUID:       c.GetString("uuid"),
	}
	fmt.Println("EmojiPackInstance==", EmojiPackInstance)
	if err := database.DB.Create(&EmojiPackInstance).Error; err != nil {
		c.JSON(500, gin.H{"code": 500, "msg": "创建表情包合集失败", "error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"code": 200, "msg": "创建表情包合集成功"})
}

type EmojiEditPackRequest struct {
	Id      uint     `json:"id" binding:"required"`
	Name    string   `json:"name" binding:"required"`
	IconURL string   `json:"iconUrl"`
	Tags    []string `json:"tags"`
}

// EmojiPackEdit 编辑表情包合集
// @Description 编辑表情包合集
// @Tags emojiPack
// @Accept json
// @Produce json
// @Param Authorization header string true "
// @Param emojiPack body EmojiEditPackRequest true "EmojiPack信息"
// @Success 200 {object} map[string]interface{} "{"code":200,"msg":"编辑表情包合集成功"}"
// @Failure 400 {object} map[string]interface{} "{"code":400,"msg":"xxxx"}"
// @Router /emojiPack/edit [post]
func EmojiPackEdit(c *gin.Context) {
	var params EmojiEditPackRequest
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(400, gin.H{"code": 400, "msg": "编辑失败", "error": err.Error()})
		return
	}

	// 处理参数
	temp, _ := json.Marshal(params.Tags)
	tagsTemp := datatypes.JSON(json.RawMessage(temp))
	// 更新表情包合集信息
	emojiPackInstance := model.EmojiPack{
		Name:    params.Name,
		IconURL: params.IconURL,
		Tags:    tagsTemp,
	}
	fmt.Println("emojiPackInstance==", emojiPackInstance)

	if database.DB.Model(&model.EmojiPack{}).Where("id = ?", params.Id).Updates(emojiPackInstance).RowsAffected == 0 {
		c.JSON(404, gin.H{"code": 404, "msg": "表情包合集未找到"})
		return
	}

	c.JSON(200, gin.H{"code": 200, "msg": "编辑表情包合集成功"})
}

type EmojiPackDeleteRequest struct {
	ID uint `json:"id" binding:"required"`
}

// EmojiPackDelete 删除表情包合集
// @Description 删除表情包合集
// @Tags emojiPack
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization"
// @Param emojiPack body EmojiPackDeleteRequest true "EmojiPack ID"
// @Success 200 {object} map[string]interface{} "{"code":200,"msg":"删除表情包合集成功"}"
// @Failure 400 {object} map[string]interface{} "{"code":400,"msg":"xxxx"}"
// @Router /emojiPack/delete [post]
func EmojiPackDelete(c *gin.Context) {
	var params EmojiPackDeleteRequest
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(400, gin.H{"code": 400, "msg": "参数错误", "error": err.Error()})
		return
	}
	fmt.Println("删除表情包合集，ID：", params.ID)
	if err := database.DB.Model(&model.EmojiPack{}).Where("id = ?", params.ID).Delete(&model.EmojiPack{}).Error; err != nil {
		c.JSON(404, gin.H{"code": 404, "msg": "表情包合集未找到", "error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"code": 200, "msg": "删除表情包合集成功"})
}

type EmojiPackDetailRequest struct {
	ID uint `json:"id" binding:"required"`
}
type EmojiPackDetailDTO struct {
	ID               uint     `json:"id"`
	Name             string   `json:"name"`
	IconURL          string   `json:"iconUrl"`
	Tags             []string `json:"tags"`
	View_count       int      `json:"view_count"`
	Collection_count int      `json:"collection_count"`
	AuthorUUID       string   `json:"authorUUID"`
}

// EmojiPackDetail 获取表情包合集详情
// @Description 获取表情包合集详情
// @Tags emojiPack
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization"
// @Param emojiPack body EmojiPackDetaildTO true "EmojiPack ID"
// @Success 200 {object} map[string]interface{} "{"code":200,"msg":"获取表情包合集详情成功","data":EmojiPackDetaildTO}"
// @Failure 400 {object} map[string]interface{} "{"code":400,"msg":"xxxx"}"
// @Router /emojiPack/detail [post]
func EmojiPackDetail(c *gin.Context) {
	var params EmojiPackDeleteRequest

	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(400, gin.H{"code": 400, "msg": "参数错误", "error": err.Error()})
		return
	}
	var emojiPack model.EmojiPack
	if err := database.DB.Where("id = ?", params.ID).First(&emojiPack).Error; err != nil {
		c.JSON(404, gin.H{"code": 404, "msg": "表情包合集未找到", "error": err.Error()})
		return
	}
	var emojiPackDetail model.EmojiPack
	var EmojiPackDetailDTO EmojiPackDetailDTO
	if err := database.DB.Model(&model.EmojiPack{}).Where("id = ?", params.ID).First(&emojiPackDetail).Scan(&EmojiPackDetailDTO).Error; err != nil {
		c.JSON(500, gin.H{"code": 500, "msg": "获取表情包合集详情失败", "error": err.Error()})
		return
	}

	//data解析
	EmojiPackDetailDTO.Tags = make([]string, 0)
	if err := json.Unmarshal(emojiPackDetail.Tags, &EmojiPackDetailDTO.Tags); err != nil {
		c.JSON(500, gin.H{"code": 500, "msg": "解析表情包合集标签失败", "error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"code": 200, "msg": "获取表情包合集详情成功", "data": EmojiPackDetailDTO})
}

type EmojiListReq struct {
	Page     int    `form:"page" binding:"required,min=1"`   // 页码
	PageSize int    `form:"pageSize" binding:"min=1,max=50"` // 每页数量
	Keyword  string `form:"keyword"`                         // 搜索关键词
}
type EmojiPackListDTO struct {
	ID               uint     `json:"id"`
	Name             string   `json:"name"`
	IconURL          string   `json:"iconUrl"`
	View_count       int      `json:"view_count"`
	Collection_count int      `json:"collection_count"`
	AuthorUUID       string   `json:"authorUUID"`
	Tags             []string `json:"tags"`
}
type EmojiPackListPageResp[T any] struct {
	List     []T   `json:"list"`
	Total    int64 `json:"total"`
	Page     int   `json:"page"`
	PageSize int   `json:"pageSize"`
}

// EmojiPackList 获取表情包合集列表(公有)
// @Description 获取表情包合集列表(公有)
// @Tags emojiPack
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "{"code":200,"msg":"获取表情包合集列表成功","data":[]}"
// @Failure 500 {object} map[string]interface{} "{"code":500,"msg":"获取表情包合集列表失败"}"
// @Router /emojiPack/list [get]
func EmojiPackList(c *gin.Context) {
	var req EmojiListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		req.Page = 1
		req.PageSize = 12
	}

	var emojiPacks []model.EmojiPack
	var total int64

	db := database.DB.Model(&model.EmojiPack{}).Where("is_default = ?", false)
	if req.Keyword != "" {
		// 根据名称模糊搜索
		db = db.Where("name LIKE ?", "%"+req.Keyword+"%")
	}
	// 统计总条数
	if err := db.Count(&total).Error; err != nil {
		c.JSON(500, gin.H{"code": 500, "msg": "统计总数失败", "error": err.Error()})
		return
	}

	// 分页查询数据
	offset := (req.Page - 1) * req.PageSize
	if err := db.Limit(req.PageSize).Offset(offset).Find(&emojiPacks).Error; err != nil {
		c.JSON(500, gin.H{"code": 500, "msg": "获取表情包合集列表失败", "error": err.Error()})
		return
	}

	// model转DTO
	listDto := make([]EmojiPackListDTO, 0, len(emojiPacks))
	for _, pack := range emojiPacks {
		var tagTemp []string
		// json字符串转切片
		tagsData := pack.Tags
		if string(tagsData) != "null" && len(tagsData) > 0 {
			if err := json.Unmarshal(pack.Tags, &tagTemp); err != nil {
				c.JSON(500, gin.H{"code": 500, "msg": "解析表情包合集标签失败", "error": err.Error()})
				return
			}
		}
		listDto = append(listDto, EmojiPackListDTO{
			ID:               pack.ID,
			Name:             pack.Name,
			IconURL:          pack.IconURL,
			View_count:       pack.View_count,
			Collection_count: pack.Collection_count,
			AuthorUUID:       pack.AuthorUUID,
			Tags:             tagTemp,
		})
	}
	// 组装分页返回数据
	pageData := EmojiPackListPageResp[EmojiPackListDTO]{
		List:     listDto,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "获取表情包合集列表成功",
		"data": pageData,
	})
}

type EmojiPackListByUserPageResp[T any] struct {
	List     []T   `json:"list"`
	Total    int64 `json:"total"`
	Page     int   `json:"page"`
	PageSize int   `json:"pageSize"`
}

// EmojiPackListByUser 获取表情包合集列表(用户)
// @Description 获取表情包合集列表(用户)
// @Tags emojiPack
// @Accept json
// @Produce json
// @Param Authorization header  true "Authorization"
// @Success 200 {object} map[string]interface{} "{"code":200,"msg":"获取表情包合集列表成功","data":[]}"
// @Failure 500 {object} map[string]interface{} "{"code":500,"msg":"获取表情包合集列表失败"}"
// @Router /emojiPack/listByUser [get]
func EmojiPackListByUser(c *gin.Context) {
	userUUID := c.GetString("uuid")

	var req EmojiListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		req.Page = 1
		req.PageSize = 12
	}

	var emojiPacks []model.EmojiPack
	var total int64

	db := database.DB.Model(&model.EmojiPack{}).
		Where("is_default = ?", false).
		Where("author_uuid = ?", userUUID)

	if req.Keyword != "" {
		// 根据名称模糊搜索
		db = db.Where("name LIKE ?", "%"+req.Keyword+"%")
	}

	// 统计总条数
	if err := db.Count(&total).Error; err != nil {
		c.JSON(500, gin.H{"code": 500, "msg": "统计总数失败", "error": err.Error()})
		return
	}

	// 分页查询数据
	offset := (req.Page - 1) * req.PageSize
	if err := db.Limit(req.PageSize).Offset(offset).Find(&emojiPacks).Error; err != nil {
		c.JSON(500, gin.H{"code": 500, "msg": "获取表情包合集列表失败", "error": err.Error()})
		return
	}

	// model转DTO
	listDto := make([]EmojiPackListDTO, 0, len(emojiPacks))
	for _, pack := range emojiPacks {
		tagTemp := make([]string, 0)
		if err := json.Unmarshal(pack.Tags, &tagTemp); err != nil {
			c.JSON(500, gin.H{"code": 500, "msg": "解析表情包合集标签失败", "error": err.Error()})
			return
		}

		// json字符串转切片
		tagsData := pack.Tags
		if string(tagsData) != "null" && len(tagsData) > 0 {
			if err := json.Unmarshal(pack.Tags, &tagTemp); err != nil {
				c.JSON(500, gin.H{"code": 500, "msg": "解析表情包合集标签失败", "error": err.Error()})
				return
			}
		}

		listDto = append(listDto, EmojiPackListDTO{
			ID:               pack.ID,
			Name:             pack.Name,
			IconURL:          pack.IconURL,
			View_count:       pack.View_count,
			Collection_count: pack.Collection_count,
			AuthorUUID:       pack.AuthorUUID,
			Tags:             tagTemp,
		})
	}
	// 组装分页返回数据
	pageData := EmojiPackListPageResp[EmojiPackListDTO]{
		List:     listDto,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "获取表情包合集列表成功",
		"data": pageData,
	})
}

type EmojiPackAddEmojiRequest struct {
	EmojiPackID uint `json:"emojiPackId" binding:"required"`
	EmojiID     uint `json:"emojiId" binding:"required"`
}

// EmojiPackAddEmoji 添加表情到表情包合集
// @Description 添加表情到表情包合集
// @Tags emojiPack
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization"
// @Param emojiPack body EmojiPackAddEmojiRequest true "EmojiPack ID和Emoji ID"
// @Success 200 {object} map[string]interface{} "{"code":200,"msg":"添加表情到表情包合集成功"}"
// @Failure 400 {object} map[string]interface{} "{"code":400,"msg":"xxxx"}"
// @Router /emojiPack/addEmoji [post]
func EmojiPackAddEmoji(c *gin.Context) {
	var params EmojiPackAddEmojiRequest
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(400, gin.H{"code": 400, "msg": "参数错误", "error": err.Error()})
		return
	}
	emojiPackEmoji := model.EmojiPack_Emoji{
		EmojiPackID: params.EmojiPackID,
		EmojiID:     params.EmojiID,
	}
	if err := database.DB.Create(&emojiPackEmoji).Error; err != nil {
		c.JSON(500, gin.H{"code": 500, "msg": "添加表情到表情包合集失败", "error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"code": 200, "msg": "添加表情到表情包合集成功"})
}

// EmojiPackAddEmoji 从表情包合集移除表情
// @Description 从表情包合集移除表情
// @Tags emojiPack
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization"
// @Param emojiPack body EmojiPackAddEmojiRequest true "EmojiPack ID和Emoji ID"
// @Success 200 {object} map[string]interface{} "{"code":200,"msg":"从表情包合集移除表情成功"}"
// @Failure 400 {object} map[string]interface{} "{"code":400,"msg":"xxxx"}"
// @Router /emojiPack/removeEmoji [post]
func EmojiPackRemoveEmoji(c *gin.Context) {
	var params EmojiPackAddEmojiRequest
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(400, gin.H{"code": 400, "msg": "参数错误", "error": err.Error()})
		return
	}
	if err := database.DB.Where("emoji_pack_id = ? AND emoji_id = ?", params.EmojiPackID, params.EmojiID).Delete(&model.EmojiPack_Emoji{}).Error; err != nil {
		c.JSON(500, gin.H{"code": 500, "msg": "从表情包合集移除表情失败", "error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"code": 200, "msg": "从表情包合集移除表情成功"})
}

type EmojiPackGetEmojisRequest struct {
	EmojiPackID uint `json:"emojiPackId" binding:"required"`
}
type EmojiPackGetEmojisDTO struct {
	Name string   `json:"name"`
	URL  string   `json:"url"`
	Tags []string `json:"tags"`
}

// emojiPackGetemojis 获取表情包合集内的表情列表
// @Description 获取表情包合集内的表情列表
// @Tags emojiPack
// @Accept json
// @Produce json
// @Param emojiPack body EmojiPackDeleteRequest true "EmojiPack ID"
// @Success 200 {object} map[string]interface{} "{"code":200,"msg":"获取表情包合集内的表情列表成功","data":[]}"
// @Failure 400 {object} map[string]interface{} "{"code":400,"msg":"xxxx"}"
// @Router /emojiPack/getemojis [post]
func EmojiPackGetEmojis(c *gin.Context) {
	var params EmojiPackGetEmojisRequest
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(400, gin.H{"code": 400, "msg": "参数错误", "error": err.Error()})
		return
	}
	var emojis []model.Emoji
	if err := database.DB.Table("emoji_pack_emojis").Select("emojis.name, emojis.url, emojis.tags").
		Joins("join emojis on emoji_pack_emojis.emoji_id = emojis.id").
		Where("emoji_pack_emojis.emoji_pack_id = ?", params.EmojiPackID).
		Scan(&emojis).Error; err != nil {
		c.JSON(500, gin.H{"code": 500, "msg": "获取表情包合集内的表情列表失败", "error": err.Error()})
		return
	}

	// 处理数据
	resp := make([]EmojiPackGetEmojisDTO, 0, len(emojis))
	for _, pack := range emojis {
		tagsTemp := make([]string, 0)
		if err := json.Unmarshal(pack.Tags, &tagsTemp); err != nil {
			c.JSON(500, gin.H{"code": 500, "msg": "解析表情标签失败", "error": err.Error()})
			return
		}
		resp = append(resp, EmojiPackGetEmojisDTO{
			Name: pack.Name,
			URL:  pack.URL,
			Tags: tagsTemp,
		})
	}

	c.JSON(200, gin.H{"code": 200, "msg": "获取表情包合集内的表情列表成功", "data": resp})
}
