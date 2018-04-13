package util

import (
	"golang.org/x/crypto/bcrypt"
	"github.com/gin-gonic/gin"
	"net"
	"fmt"
	"net/http"
	"log"
	"github.com/arizalsaputro/go-ddd/config"
	"time"
	"github.com/dgrijalva/jwt-go"
)

func HashPassword(password string) (string,error)  {
	bytes,err := bcrypt.GenerateFromPassword([]byte(password),8)
	return string(bytes),err
}


func CheckPassword(password string,hash string)bool  {
	err := bcrypt.CompareHashAndPassword([]byte(hash),[]byte(password))
	return err == nil
}

func ServeError(c *gin.Context,status int,message string)  {
	c.JSON(status,gin.H{"error":true,"message":message})
}

func ServeOk(c *gin.Context,status int,data interface{})  {
	c.JSON(status,gin.H{"error":false,"data":data})
}

func GetClientIPByRequest(req *http.Request) (ip string) {

	// Try via request
	ip, port, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		log.Printf("debug: Getting req.RemoteAddr %v", err)
		return ""
	} else {
		log.Printf("debug: With req.RemoteAddr found IP:%v; Port: %v", ip, port)
	}

	userIP := net.ParseIP(ip)
	if userIP == nil {
		message := fmt.Sprintf("debug: Parsing IP from Request.RemoteAddr got nothing.")
		log.Printf(message)
		return ""

	}
	log.Printf("debug: Found IP: %v", userIP)
	return userIP.String()
}

func GenerateAccessToken(userId string,userType string)(string,error){
	fmt.Println(userId)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,jwt.MapClaims{
		"id":userId,
		"type":userType,
		"iat":time.Now().Unix(),
	})
	tokenString, err := token.SignedString([]byte(config.Conf.SecretJwt))
	return tokenString,err
}