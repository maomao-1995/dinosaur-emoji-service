package handler

import (
	"dinosaur-emoji-service/internal/model"
	"dinosaur-emoji-service/pkg/database"
	"fmt"

	"encoding/json"

	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type EmojiAddPackRequest struct {
	Name        string   `json:"name" binding:"required"`
	IconURL     string   `json:"iconUrl"`
	Tags        []string `json:"tags"`
	Description string   `json:"description"`
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
		c.JSON(400, gin.H{"code": 400, "msg": "参数错误", "error": err.Error()})
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
		Description:      params.Description,
	}
	fmt.Println("EmojiPackInstance==", EmojiPackInstance)
	if err := database.DB.Create(&EmojiPackInstance).Error; err != nil {
		c.JSON(500, gin.H{"code": 500, "msg": "创建表情包合集失败", "error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"code": 200, "msg": "创建表情包合集成功"})
}

type EmojiEditPackRequest struct {
	Id          uint     `json:"id" binding:"required"`
	Name        string   `json:"name" binding:"required"`
	IconURL     string   `json:"iconUrl"`
	Tags        []string `json:"tags"`
	Description string   `json:"description"`
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
		Name:        params.Name,
		IconURL:     params.IconURL,
		Tags:        tagsTemp,
		Description: params.Description,
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
	ID uint `form:"id" binding:"required"`
}
type AuthorInfo struct {
	// UUID        string `json:"uuid"`
	Username    string `json:"username"`
	Avatar      string `json:"avatar"`
	Description string `json:"description"`
}

type EmojiPackDetailDTO struct {
	ID              uint       `json:"id"`
	Name            string     `json:"name"`
	IconURL         string     `json:"iconUrl"`
	Tags            []string   `json:"tags"`
	ViewCount       int        `json:"viewCount"`
	Description     string     `json:"description"`
	AuthorUUID      string     `json:"authorUUID"`
	AuthorInfo      AuthorInfo `json:"authorInfo"`
	CollectionCount int        `json:"collectionCount"`
	IsCollection    bool       `json:"isCollection"`
	CanEdit         bool       `json:"canEdit"`
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
	var params EmojiPackDetailRequest

	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(400, gin.H{"code": 400, "msg": "参数错误", "error": err.Error()})
		return
	}

	var emojiPackDetail model.EmojiPack
	var emojiPackDetailDTO EmojiPackDetailDTO
	if err := database.DB.Model(&model.EmojiPack{}).Where("id = ?", params.ID).First(&emojiPackDetail).Error; err != nil {
		c.JSON(500, gin.H{"code": 500, "msg": "获取表情包合集详情失败", "error": err.Error()})
		return
	}
	//通过作者UUID获取作者信息
	var authorInfo model.User
	if err := database.DB.Model(&model.User{}).Where("uuid = ?", emojiPackDetail.AuthorUUID).First(&authorInfo).Error; err != nil {
		c.JSON(500, gin.H{"code": 500, "msg": "获取作者信息失败", "error": err.Error()})
		return
	}

	//data解析
	emojiPackDetailDTO = EmojiPackDetailDTO{
		ID:          emojiPackDetail.ID,
		Name:        emojiPackDetail.Name,
		IconURL:     emojiPackDetail.IconURL,
		ViewCount:   emojiPackDetail.View_count,
		Description: emojiPackDetail.Description,
		AuthorUUID:  emojiPackDetail.AuthorUUID,
	}
	//tag是null时，解析会报错，所以先初始化一个空切片
	if emojiPackDetail.Tags == nil {
		emojiPackDetailDTO.Tags = make([]string, 0)
	} else {
		emojiPackDetailDTO.Tags = make([]string, 0)
		if err := json.Unmarshal(emojiPackDetail.Tags, &emojiPackDetailDTO.Tags); err != nil {
			c.JSON(500, gin.H{"code": 500, "msg": "解析表情包合集标签失败", "error": err.Error()})
			return
		}
	}

	emojiPackDetailDTO.AuthorInfo = AuthorInfo{
		// UUID:        authorInfo.Uuid,
		Username:    authorInfo.Username,
		Avatar:      authorInfo.Avatar,
		Description: authorInfo.Description,
	}
	//查询收藏数
	Db := database.DB.Model(&model.EmojiPackCollection{}).Where("emoji_pack_id = ?", params.ID)
	var collectionCount int64
	if err := Db.Count(&collectionCount).Error; err != nil {
		c.JSON(500, gin.H{"code": 500, "msg": "获取收藏数失败", "error": err.Error()})
		return
	}
	//查询当前用户是否收藏
	userUUID := c.GetString("uuid")
	var isCollection bool
	if userUUID != "" {
		var collection model.EmojiPackCollection
		if err := database.DB.Model(&model.EmojiPackCollection{}).Where("emoji_pack_id = ? AND user_uuid = ?", params.ID, userUUID).First(&collection).Error; err == nil {
			isCollection = true
		}
	}

	emojiPackDetailDTO.CollectionCount = int(collectionCount)
	emojiPackDetailDTO.IsCollection = isCollection

	//查询当前用户是否是作者
	canEdit := false
	if userUUID == emojiPackDetail.AuthorUUID {
		canEdit = true
	}
	emojiPackDetailDTO.CanEdit = canEdit

	c.JSON(200, gin.H{"code": 200, "msg": "获取表情包合集详情成功", "data": emojiPackDetailDTO})
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

	var emojiPacks []model.EmojiPack

	db := database.DB.Model(&model.EmojiPack{}).
		// Where("is_default = ?", false).
		Where("author_uuid = ?", userUUID)
	if err := db.Find(&emojiPacks).Error; err != nil {
		c.JSON(500, gin.H{"code": 500, "msg": "获取表情包合集列表失败", "error": err.Error()})
		return
	}
	// model转DTO
	listDto := make([]EmojiPackListDTO, 0, len(emojiPacks))
	for _, pack := range emojiPacks {
		tagTemp := make([]string, 0)
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

	db02 := database.DB.Model(&model.EmojiPackCollection{}).Where("user_uuid = ?", userUUID)
	var collections []model.EmojiPackCollection
	if err := db02.Find(&collections).Error; err != nil {
		c.JSON(500, gin.H{"code": 500, "msg": "获取表情包合集列表失败", "error": err.Error()})
		return
	}
	collectionPackIds := make([]uint, 0, len(collections))
	for _, collection := range collections {
		collectionPackIds = append(collectionPackIds, collection.EmojiPackID)
	}
	listDto02 := make([]EmojiPackListDTO, 0, len(emojiPacks))
	if len(collectionPackIds) > 0 {
		var collectionPacks []model.EmojiPack
		if err := database.DB.Model(&model.EmojiPack{}).Where("id IN ?", collectionPackIds).Find(&collectionPacks).Error; err != nil {
			c.JSON(500, gin.H{"code": 500, "msg": "获取表情包合集列表失败", "error": err.Error()})
			return
		}
		for _, pack := range collectionPacks {
			tagTemp := make([]string, 0)
			// json字符串转切片
			tagsData := pack.Tags
			if string(tagsData) != "null" && len(tagsData) > 0 {
				if err := json.Unmarshal(pack.Tags, &tagTemp); err != nil {
					c.JSON(500, gin.H{"code": 500, "msg": "解析表情包合集标签失败", "error": err.Error()})
					return
				}
			}
			listDto02 = append(listDto02, EmojiPackListDTO{
				ID:               pack.ID,
				Name:             pack.Name,
				IconURL:          pack.IconURL,
				View_count:       pack.View_count,
				Collection_count: pack.Collection_count,
				AuthorUUID:       pack.AuthorUUID,
				Tags:             tagTemp,
			})
		}

	}

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "获取表情包合集列表成功",
		"data": gin.H{"createList": listDto, "collectionList": listDto02},
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
	// 如有数据库中已存在该表情包和表情的关联，则不再添加
	var existingRelation model.EmojiPackEmoji
	if err := database.DB.Where("emoji_pack_id = ? AND emoji_id = ?", params.EmojiPackID, params.EmojiID).First(&existingRelation).Error; err == nil {
		c.JSON(400, gin.H{"code": 200, "msg": "表情已存在于该表情包合集中"})
		return
	}

	emojiPackEmoji := model.EmojiPackEmoji{
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
	if err := database.DB.Where("emoji_pack_id = ? AND emoji_id = ?", params.EmojiPackID, params.EmojiID).Delete(&model.EmojiPackEmoji{}).Error; err != nil {
		c.JSON(500, gin.H{"code": 500, "msg": "从表情包合集移除表情失败", "error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"code": 200, "msg": "从表情包合集移除表情成功"})

}

type EmojiPackGetEmojisRequest struct {
	EmojiPackID uint `form:"emojiPackId" binding:"required"`
	Page        int  `form:"page" binding:"omitempty,min=1"`
	PageSize    int  `form:"pageSize" binding:"omitempty,min=1,max=50"`
}
type EmojiPackGetEmojisDTO struct {
	ID   uint     `json:"id"`
	Name string   `json:"name"`
	URL  string   `json:"url"`
	Tags []string `json:"tags"`
}
type EmojiPackGetEmojisResponse[T any] struct {
	List     []T   `json:"list"`
	Total    int64 `json:"total"`
	Page     int   `json:"page"`
	PageSize int   `json:"pageSize"`
}

// emojiPackGetemojis 获取表情包合集内的表情列表
// @Description 获取表情包合集内的表情列表
// @Tags emojiPack
// @Accept json
// @Produce json
// @Param emojiPack query EmojiPackGetEmojisRequest true "EmojiPack ID and pagination"
// @Success 200 {object} map[string]interface{} "{"code":200,"msg":"获取表情包合集内的表情列表成功","data":[]}"
// @Failure 400 {object} map[string]interface{} "{"code":400,"msg":"xxxx"}"
// @Router /emojiPack/getemojis [get]
func EmojiPackGetEmojis(c *gin.Context) {
	var params EmojiPackGetEmojisRequest
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(400, gin.H{"code": 400, "msg": "参数错误", "error": err.Error()})
		return
	}

	if params.Page == 0 {
		params.Page = 1
	}
	if params.PageSize == 0 {
		params.PageSize = 12
	}

	var total int64
	countDB := database.DB.Table("emoji_pack_emojis").
		Joins("join emojis on emoji_pack_emojis.emoji_id = emojis.id").
		Where("emoji_pack_emojis.emoji_pack_id = ?", params.EmojiPackID)
	if err := countDB.Count(&total).Error; err != nil {
		c.JSON(500, gin.H{"code": 500, "msg": "统计表情包内表情总数失败", "error": err.Error()})
		return
	}

	offset := (params.Page - 1) * params.PageSize
	var emojis []model.Emoji
	queryDB := database.DB.Table("emoji_pack_emojis").Select("emojis.id, emojis.name, emojis.url, emojis.tags").
		Joins("join emojis on emoji_pack_emojis.emoji_id = emojis.id").
		Where("emoji_pack_emojis.emoji_pack_id = ?", params.EmojiPackID).
		Limit(params.PageSize).Offset(offset)
	if err := queryDB.Scan(&emojis).Error; err != nil {
		c.JSON(500, gin.H{"code": 500, "msg": "获取表情包合集内的表情列表失败", "error": err.Error()})
		return
	}

	resp := make([]EmojiPackGetEmojisDTO, 0, len(emojis))
	for _, pack := range emojis {
		tagsTemp := make([]string, 0)
		if err := json.Unmarshal(pack.Tags, &tagsTemp); err != nil {
			c.JSON(500, gin.H{"code": 500, "msg": "解析表情标签失败", "error": err.Error()})
			return
		}
		resp = append(resp, EmojiPackGetEmojisDTO{
			ID:   pack.ID,
			Name: pack.Name,
			URL:  pack.URL,
			Tags: tagsTemp,
		})
	}

	pageData := EmojiPackGetEmojisResponse[EmojiPackGetEmojisDTO]{
		List:     resp,
		Total:    total,
		Page:     params.Page,
		PageSize: params.PageSize,
	}

	c.JSON(200, gin.H{"code": 200, "msg": "获取表情包合集内的表情列表成功", "data": pageData})
}

type EmojiPackCollectionRequest struct {
	EmojiPackID  uint  `json:"emojiPackId" binding:"required"`
	IsCollection *bool `json:"isCollection" binding:"required"`
}

//收藏/取消 表情包集合
// @Description 收藏表情包集合
// @Tags emojiPack
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization"

// @Param emojiPack body EmojiPackCollectionRequest true "EmojiPack ID和收藏状态"
// @Success 200 {object} map[string]interface{} "{"code":200,"msg":"收藏表情包合集成功"}"
// @Failure 400 {object} map[string]interface{} "{"code":400,"msg":"xxxx"}"
// @Router /emojiPack/collection [post]

func EmojiPackCollection(c *gin.Context) {
	//IsCollection false 取消收藏 true 收藏

	var params EmojiPackCollectionRequest
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(400, gin.H{"code": 400, "msg": "参数错误", "error": err.Error()})
		return
	}
	// 获取用户UUID
	userUUID := c.GetString("uuid")
	if params.IsCollection == nil {
		c.JSON(400, gin.H{"code": 400, "msg": "参数错误", "error": "isCollection is required"})
		return
	}

	if *params.IsCollection {
		// 添加收藏
		collection := model.EmojiPackCollection{
			EmojiPackID: params.EmojiPackID,
			UserUUID:    userUUID,
		}
		// 检查是否已经收藏过
		var existingCollection model.EmojiPackCollection
		if err := database.DB.Where("emoji_pack_id = ? AND user_uuid = ?", params.EmojiPackID, userUUID).First(&existingCollection).Error; err == nil {
			c.JSON(200, gin.H{"code": 200, "msg": "已经收藏过该表情包合集"})
			return
		}
		if err := database.DB.Create(&collection).Error; err != nil {
			c.JSON(500, gin.H{"code": 500, "msg": "收藏表情包合集失败", "error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"code": 200, "msg": "收藏表情包合集成功"})
	} else {
		// 取消收藏
		if err := database.DB.Where("emoji_pack_id = ? AND user_uuid = ?", params.EmojiPackID, userUUID).Delete(&model.EmojiPackCollection{}).Error; err != nil {
			c.JSON(500, gin.H{"code": 500, "msg": "取消收藏表情包合集失败", "error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"code": 200, "msg": "取消收藏表情包合集成功"})
	}
}

type EmojiPackAddViewCountRequest struct {
	EmojiPackID uint `json:"emojiPackId" binding:"required"`
}

// 增加emojiPack viewCount
func EmojiPackAddViewCount(c *gin.Context) {
	var params EmojiPackAddViewCountRequest
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(400, gin.H{"code": 400, "msg": "参数错误", "error": err.Error()})
		return
	}
	// 使用原子操作增加 view_count
	if err := database.DB.Model(&model.EmojiPack{}).Where("id = ?", params.EmojiPackID).UpdateColumn("view_count", gorm.Expr("view_count + ?", 1)).Error; err != nil {
		c.JSON(500, gin.H{"code": 500, "msg": "增加表情包合集浏览量失败", "error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"code": 200, "msg": "增加表情包合集浏览量成功"})
}
