package handler

import (
	"dinosaur-emoji-service/internal/model"
	"dinosaur-emoji-service/pkg/database"
	"dinosaur-emoji-service/pkg/jwtMain"
	"dinosaur-emoji-service/pkg/redisMain"
	"dinosaur-emoji-service/pkg/utils"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserDTO struct {
	Username string `json:"username"`
	Phone    int64  `json:"phone"`
	Email    string `json:"email"`
	Nickname string `json:"nickname"`
	Id       int64  `json:"id"`
	Uuid     string `json:"uuid"`
}

// @Summary 获取用户信息
// @Description 获取用户信息
// @Tags user
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization"
// @Success 200 {object} map[string]interface{} "{"code":200,"msg":"获取用户信息成功","data":UserDTO}"
// @Failure 400 {object} map[string]interface{} "{"code":400,"msg":"xxxxx"}"
// @Router /user/info [get]
func UserInfo(c *gin.Context) {
	uuid := c.MustGet("uuid").(string)
	var userTemp UserDTO
	err := database.DB.
		Model(&model.User{}).
		Select("id,username, phone,email,nickname").
		Where("uuid = ?", uuid).
		Scan(&userTemp).Error
	userTemp.Uuid = uuid
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "用户不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "获取用户信息成功", "data": userTemp})
}

type UserRequest struct {
	Username  string `json:"username" binding:"required"`
	Birthdate string `json:"birthdate"`
	Phone     string `json:"phone" binding:"required"`
	Email     string `json:"email"`
	Password  string `json:"password" binding:"required"`
	Nickname  string `json:"nickname"`
	Code      string `json:"code" binding:"required"` // 验证码
}

// @Summary 用户注册
// @Description 用户注册
// @Tags user
// @Accept json
// @Produce json
// @Param user body UserRequest true "用户信息"
// @Success 200 {object} map[string]interface{} "{"code":200,"msg":"注册成功"}"
// @Failure 400 {object} map[string]interface{} "{"code":400,"msg":"xxxx"}"
// @Router /register [post]
func UserRegister(c *gin.Context) {
	var params UserRequest

	jsonErr := c.ShouldBindJSON(&params)
	if jsonErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误", "error": jsonErr.Error()})
		return
	}
	//初始化参数
	//Nickname
	if params.Nickname == "" {
		params.Nickname = utils.GenerateRandomNickname()
	}
	//Birthdate
	birthdate, birthdateErr := time.Parse("2006-01-02", params.Birthdate)
	if birthdateErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "error": "出生日期格式错误，请使用YYYY-MM-DD"})
		return
	}
	//Password
	hashedPassword, passwordErr := bcrypt.GenerateFromPassword(
		[]byte(params.Password),
		bcrypt.DefaultCost, // 默认成本系数10，范围4-31
	)
	if passwordErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 400, "error": "密码加密失败"})
		return
	}
	hashedPasswordString := string(hashedPassword)

	code, getCodeErr := redisMain.Rdb.Get(redisMain.Ctx, params.Phone).Result()
	if getCodeErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "请先发送验证码", "error": getCodeErr.Error()})
		return
	}
	codeTrue := code != params.Code
	if codeTrue {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "验证码错误"})
		return
	}

	//创建model实例
	newUser := model.User{
		Username:  params.Username,
		Birthdate: birthdate,
		Phone:     params.Phone,
		Email:     params.Email,
		Password:  hashedPasswordString,
		Nickname:  params.Nickname,
		Uuid:      uuid.New().String(),
	}

	selectErr01 := database.DB.Where("phone = ?", params.Phone).First(&newUser).Error
	if selectErr01 == nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "手机号已注册"})
		return
	}

	selectErr02 := database.DB.Where("username = ?", params.Username).First(&newUser).Error
	if selectErr02 == nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "用户名已存在"})
		return
	}

	selectErr03 := database.DB.Create(&newUser).Error
	if selectErr03 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "注册失败，请稍后重试", "error": selectErr03.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "注册成功",
	})
}

type LoginRequest struct {
	Phone    string `json:"phone" binding:""`
	Code     string `json:"code" binding:""`
	Username string `json:"username" binding:""`     // 用户名
	Password string `json:"password" binding:""`     // 密码
	Type     string `json:"type" binding:"required"` // 登录类型，phone 或 username
}

// @Summary 用户登录
// @Description 用户登录
// @Tags user
// @Accept json
// @Produce json
// @Param user body LoginRequest true "登录信息"
// @Success 200 {object} map[string]interface{} "{"code":200,"msg":"登录成功"}"
// @Failure 400 {object} map[string]interface{} "{"code":400,"msg":"xxxxx",token:"xxxx"}"
// @Router /login [post]
func Login(c *gin.Context) {
	var params LoginRequest

	jsonErr := c.ShouldBindJSON(&params)
	if jsonErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "error": jsonErr.Error(), "msg": "参数错误"})
		return
	}

	var userTemp model.User
	if params.Type == "phone" {
		selectErr01 := database.DB.Where("phone = ?", params.Phone).First(&userTemp).Error
		if selectErr01 != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "账号错误"})
			return
		}

	}
	if params.Type == "username" {
		selectErr02 := database.DB.Where("username = ?", params.Username).First(&userTemp).Error
		if selectErr02 != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "用户不存在"})
			return
		}

		passWordErr := bcrypt.CompareHashAndPassword([]byte(userTemp.Password), []byte(params.Password))
		switch {
		case passWordErr == bcrypt.ErrMismatchedHashAndPassword:
			c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "密码错误"})
			return
		case passWordErr != nil:
			c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "密码校验失败"})
			return
		default:
		}
	}
	fmt.Println("userTemp", userTemp)
	token01, _ := jwtMain.GenerateToken(userTemp.Uuid, time.Now().Add(5*time.Minute))
	token02, _ := jwtMain.GenerateToken(userTemp.Uuid, time.Now().Add(7*24*time.Minute))
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "登录成功", "data": gin.H{"accessToken": token01, "refreshToken": token02}})
}

type LoginAndRegisterReq struct {
	Phone string `json:"phone" binding:"required"`
	Code  string `json:"code" binding:"required"`
}

// @Summary 用户登录以及自动注册
// @Description 用户登录以及自动注册
// @Tags global
// @Accept json
// @Produce json
// @Param user body LoginAndRegisterRequest true "登录信息"
// @Success 200 {object} map[string]interface{} "{"code":200,"msg":"登录成功"}"
// @Failure 400 {object} map[string]interface{} "{"code":400,"msg":"xxxxx",token:"xxxx"}"
// @Router /login [post]
func LoginAndRegister(c *gin.Context) {
	var params LoginAndRegisterReq

	jsonErr := c.ShouldBindJSON(&params)
	if jsonErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "error": jsonErr.Error(), "msg": "参数错误"})
		return
	}

	code, getCodeErr := redisMain.Rdb.Get(redisMain.Ctx, params.Phone).Result()
	if getCodeErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "请先发送验证码", "error": getCodeErr.Error()})
		return
	}
	codeTrue := code != params.Code
	if codeTrue {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "验证码错误"})
		return
	}

	var newUser model.User
	selectErr01 := database.DB.Where("phone = ?", params.Phone).First(&newUser).Error

	if selectErr01 != nil {
		newUser = model.User{
			Phone:    params.Phone,
			Uuid:     uuid.New().String(),
			Nickname: utils.GenerateRandomNickname(),
		}
		createErr := database.DB.Omit("birthdate").Create(&newUser).Error
		if createErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "自动注册失败，请稍后重试", "error": createErr.Error()})
			return
		}
		fmt.Println("newUser01", newUser)
		token01, _ := jwtMain.GenerateToken(newUser.Uuid, time.Now().Add(5*time.Minute))
		token02, _ := jwtMain.GenerateToken(newUser.Uuid, time.Now().Add(7*24*time.Minute))
		c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "登录成功", "data": gin.H{"accessToken": token01, "refreshToken": token02}})
		return
	}
	fmt.Println("newUser02", newUser)
	token01, _ := jwtMain.GenerateToken(newUser.Uuid, time.Now().Add(5*time.Minute))
	token02, _ := jwtMain.GenerateToken(newUser.Uuid, time.Now().Add(7*24*time.Minute))
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "登录成功", "data": gin.H{"accessToken": token01, "refreshToken": token02}})
}
