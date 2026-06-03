package handler

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"dinosaur-emoji-service/pkg/jwtMain"
	"dinosaur-emoji-service/pkg/redisMain"
	"dinosaur-emoji-service/pkg/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SendCodeRequest struct {
	Phone string `json:"phone" binding:"required"`
}

// @Summary 发送注册手机验证码
// @Description 发送注册手机验证码
// @Tags global
// @Accept json
// @Produce json
// @Param user body SendCodeRequest true "手机号"
// @Success 200 {object} map[string]interface{} "{"code":200,"msg":"验证码发送成功"}"
// @Failure 400 {object} map[string]interface{} "{"code":400,"msg":"参数错误"}"
// @Router /sendCode [post]
func SendCode(c *gin.Context) {
	var params SendCodeRequest

	jsonErr := c.ShouldBindJSON(&params)
	if jsonErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "error": jsonErr.Error(), "msg": "参数错误"})
		return
	}

	// var userTemp model.User
	// selectErr := database.DB.Where("phone = ?", params.Phone).First(&userTemp).Error
	// if selectErr == nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "手机号已注册"})
	// 	return
	// }

	_, getErr := redisMain.Rdb.Get(redisMain.Ctx, params.Phone).Result()
	if getErr == nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "验证码未过期，请稍后再试"})
		return
	}

	setErr := redisMain.Rdb.Set(redisMain.Ctx, params.Phone, "123", 1*time.Minute).Err()
	if setErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "设置验证码失败", "error": setErr.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "验证码发送成功"})
}

// @Summary 刷新重置token
// @Description 刷新重置token
// @Tags global
// @Accept json
// @Produce json
// @Param refresh_token header string true "refresh_token" // 需要在请求头中传递 refresh_token
// @Success 200 {object} map[string]interface{} "{"code":200,"msg":"刷新重置token成功","token":"xxxx"}"
// @Failure 400 {object} map[string]interface{} "{"code":401,"msg":"xxxxx"}"
// @Router /refresh [get]
func Refresh(c *gin.Context) {
	tokenStr := c.GetHeader("refresh_token")
	if tokenStr == "" || !strings.HasPrefix(tokenStr, "Bearer ") {
		c.AbortWithStatusJSON(401, gin.H{"code": 400, "msg": "no token"})
		return
	}
	tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")
	claims, claimsErr := jwtMain.ParseToken(tokenStr)
	if claimsErr != nil {
		c.AbortWithStatusJSON(401, gin.H{"code": 401, "msg": "非法或过期 token", "error": claimsErr.Error()})
		return
	}

	token01, _ := jwtMain.GenerateToken(claims.Uuid, time.Now().Add(5*time.Minute))
	token02, _ := jwtMain.GenerateToken(claims.Uuid, time.Now().Add(7*24*time.Minute))
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "刷新重置token成功", "data": gin.H{"accessToken": token01, "refreshToken": token02}})
}

// @Summary 文件上传
// @Description 文件上传
// @Tags global
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "文件"
// @Success 200 {object} map[string]interface{} "{"code":200,"msg":"文件上传成功","url":"xxxx"}"
// @Failure 400 {object} map[string]interface{} "{"code":400,"msg":"xxxxx"}"
// @Router /upload [post]
func Upload(c *gin.Context) {
	// 获取图片文件
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "error": "未上传文件"})
		return
	}
	defer file.Close()

	// 计算文件的 MD5 哈希值
	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "error": "计算文件哈希值失败"})
		return
	}
	fileHash := hex.EncodeToString(hash.Sum(nil))

	// 重置文件指针
	file.Seek(0, 0)

	// 检查文件是否已存在
	if existingFileName, ok := utils.CheckHash(fileHash); ok {
		// 文件内容已存在，返回文件的访问 URL
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "文件已存在，无需重复上传",
			"url":     fmt.Sprintf("/uploads/%s", existingFileName),
		})
		return
	}

	// 创建文件夹
	dirPath := "./uploads"
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		os.MkdirAll(dirPath, os.ModePerm) // 创建目录
	}

	// 生成新的文件名（加入 UUID）
	fileExtension := filepath.Ext(header.Filename)
	newFileName := uuid.New().String() + fileExtension
	newFilePath := filepath.Join(dirPath, newFileName)

	// 保存文件
	if err := c.SaveUploadedFile(header, newFilePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "error": "保存文件失败"})
		return
	}

	// 保存哈希值到 Redis
	utils.SaveHash(fileHash, newFileName)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "文件上传成功",
		"url":     fmt.Sprintf("/uploads/%s", newFileName),
	})
}
