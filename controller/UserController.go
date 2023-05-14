package controller

import (
	"log"
	"net/http"
	"slip/common"
	"slip/model"
	"slip/util"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Register(ctx *gin.Context) {
	DB := common.GetDB()
	//获取参数
	name := ctx.PostForm("name")
	telephone := ctx.PostForm("telephone")
	password := ctx.PostForm("password")

	//数据验证
	if len(telephone) != 11 {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "手机号必须为11位"}) //这段代码使用了 Gin 框架中的上下文对象（Context）提供的 JSON 方法，将一个 JSON 格式的数据作为 HTTP 响应返回给客户端。
		return
	}
	if len(password) < 6 {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "密码不能少于6位"})
		return
	}
	//如果名称没有传，给一个10位随机字符串
	if len(name) == 0 {
		name = util.RandomString(10)
	}

	log.Println(name, telephone, password)
	//判断手机号是否存在
	if isTelephoneExist(DB, telephone) {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "用户已经存在！"})
		return
	}

	//创建用户
	hasedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 500, "msg": "密码错误！"})
		return
	}
	newUser := model.User{
		Name:      name,
		Telephone: telephone,
		Password:  string(hasedPassword),
	}
	DB.Create(&newUser)

	//返回结果

	ctx.JSON(200, gin.H{
		"code": 200,
		"msg":  "注册成功",
	})
}

func Login(ctx *gin.Context) {
	//获取参数
	DB := common.GetDB()
	//获取参数
	telephone := ctx.PostForm("telephone")
	password := ctx.PostForm("password")
	//数据验证
	if len(telephone) != 11 {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "手机号必须为11位"}) //这段代码使用了 Gin 框架中的上下文对象（Context）提供的 JSON 方法，将一个 JSON 格式的数据作为 HTTP 响应返回给客户端。
		return
	}
	if len(password) < 6 {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "密码不能少于6位"})
		return
	}
	//判断手机号是否存在
	var user model.User
	//使用了 GORM 库提供的 API，对数据库进行查询操作
	DB.Where("telephone = ?", telephone).First(&user)
	if user.ID == 0 {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "用户不存在"})
		return
	}
	//判断密码是否正确
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "密码错误！"})
		return
	}
	//发放token
	token, err := common.ReleaseToken(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "系统异常！"})
		log.Printf("token generate error:%v", err)
		return
	}
	//返回结果
	ctx.JSON(200, gin.H{
		"code": 200,
		"data": gin.H{"token": token},
		"msg":  "登陆成功",
	})
}

func Info(ctx *gin.Context) {
	user, _ := ctx.Get("user")
	ctx.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"user": user}})
}

func isTelephoneExist(db *gorm.DB, telephone string) bool {
	var user model.User
	//使用了 GORM 库提供的 API，对数据库进行查询操作
	db.Where("telephone = ?", telephone).First(&user)
	//Where("telephone = ?", telephone) 表示将查询条件设置为 telephone 字段等于指定的 telephone 参数的值
	//.First(&user) 则表示根据上述查询条件，从数据库中查找第一条符合条件的记录，并将其数据保存到 user 变量中。
	return user.ID != 0
	//如果查询结果为空，则 user.ID 的值为 0，否则为数据库中对应记录的 ID 值
	//检查数据库中是否已存在具有指定手机号码的用户记录。如果存在，则返回 true，否则返回 false。
}
