package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/arizalsaputro/go-ddd/interface/webcontrollers/util"
	"github.com/arizalsaputro/go-ddd/usecase"
	"net/http"
	"gopkg.in/go-playground/validator.v8"
	"strings"
	"github.com/dchest/uniuri"
	"github.com/arizalsaputro/go-ddd/config"
	ginJwt "github.com/appleboy/gin-jwt"
)



func JwtGinAuthUser() *ginJwt.GinJWTMiddleware{
	authMiddleware := &ginJwt.GinJWTMiddleware{
		Realm:"bikiniBottom",
		Key:[]byte(config.Conf.SecretJwt),
		Unauthorized: func(c *gin.Context, code int, message string) {
			au := c.GetHeader("Authorization")
			if au == ""{
				util.ServeError(c,http.StatusUnauthorized,"Authorization Header required")
				return
			}
			util.ServeError(c,http.StatusUnauthorized,"Invalid Access Token")
		},
		Authorizator: func(userId string, c *gin.Context) bool {
			claims := ginJwt.ExtractClaims(c)
			if claims["type"] != "user" {
				return false
			}
			c.Set("USER",claims["id"])
			return true
		},
		TokenLookup: "header:Authorization",
		TokenHeadName: "Bearer",
	}
	return  authMiddleware
}


func ValidateLogin(class interface{})gin.HandlerFunc  {
	return func(c *gin.Context) {
		login,ok := class.(usecase.UserLogin)
		if ok{
			err:=c.ShouldBindJSON(&login)
			if err != nil{
				HandleErrorBindJson(c,err)
				c.Abort()
				return
			}
			login.Email = strings.Trim(strings.ToLower(login.Email)," ")
			login.Password = strings.Trim(login.Password," ")
			c.Set("JSON",login)
			c.Next()
		}else{
			c.JSON(http.StatusInternalServerError,gin.H{"error":true,"message":"highly trained monkey has been dispatch"})
		}
	}
}

func ValidateRegister(class interface{}) gin.HandlerFunc  {
	
	return func(c *gin.Context) {
		register,ok := class.(usecase.UserRegister)
		if ok {
			err:=c.ShouldBindJSON(&register)
			if err != nil{
				HandleErrorBindJson(c,err)
				c.Abort()
				return
			}
			register.Email = strings.Trim(strings.ToLower(register.Email)," ")
			register.Name = strings.Trim(register.Name," ")
			register.Password = strings.Trim(register.Password," ")
			register.Code = uniuri.NewLen(40)
			hash,err := util.HashPassword(register.Password)
			if err != nil{
				c.JSON(http.StatusInternalServerError,gin.H{"error":true,"message":err.Error()})
				c.Abort()
				return
			}
			register.Password = hash
			c.Set("JSON",register)
			c.Next()
		}else{
			c.JSON(http.StatusInternalServerError,gin.H{"error":true,"message":"highly trained monkey has been dispatch"})
		}
	}
}





func HandleErrorBindJson(c *gin.Context,errs error)  {
	lerr,ok := errs.(validator.ValidationErrors)
	if !ok {
		util.ServeError(c,http.StatusBadRequest,errs.Error())
		return
	}

	for _,err := range lerr{
		if err.Field == "Name" {
			if err.Tag == "required" {
				util.ServeError(c,http.StatusBadRequest,"name required")
				return
			}
		}
		if err.Field == "Key" {
			if err.Tag == "required" {
				util.ServeError(c,http.StatusBadRequest,"unique key required")
				return
			}
		}
		if err.Field == "Email" {
			if err.Tag == "required" {
				util.ServeError(c,http.StatusBadRequest,"email required")
			}
			if err.Tag == "email" {
				util.ServeError(c,http.StatusBadRequest,"invalid email address")
			}
			return
		}
		if err.Field == "Password" {
			if err.Tag == "required" {
				util.ServeError(c,http.StatusBadRequest,"password required")
			}
			if err.Tag == "min" {
				util.ServeError(c,http.StatusBadRequest,"password to short")
			}
			return
		}
	}
	c.JSON(http.StatusInternalServerError,gin.H{"error":true,"message":errs.(validator.ValidationErrors)})

}

