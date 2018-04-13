package main

import (
	"github.com/arizalsaputro/go-ddd/infrastructure"
	"log"
	"flag"
	"github.com/arizalsaputro/go-ddd/interface/webcontrollers"
	"github.com/arizalsaputro/go-ddd/interface/webcontrollers/middleware"
	"github.com/arizalsaputro/go-ddd/usecase"
	"github.com/arizalsaputro/go-ddd/interface/repositories/mongo"
	"github.com/gin-gonic/gin"
	"net/http"
)

func init()  {
	err:= infrastructure.ConnectMongo()
	if err != nil{
		log.Fatal(err)
	}
}

var (
	webServiceHandler	= new(webcontrollers.WebServiceHandler)
)

func main()  {
	flag.Parse()

	userInteractor := new(usecase.UserInteractor)
	userInteractor.UserRepository = mongo.MongoRepo{}
	webServiceHandler.UserInteractor = userInteractor

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	type Empty struct {

	}

	router.GET("/", func(c *gin.Context) {
		c.JSON(200,Empty{})
	})

	router.GET("/ping", func(context *gin.Context) {
		context.JSON(http.StatusOK,gin.H{"message":"pong"})
	})

	v1 := router.Group("/v1")
	{
		v1.POST("/user/register",middleware.ValidateRegister(usecase.UserRegister{}),func(context *gin.Context) {
			webServiceHandler.RegisterUser(context)
		})
		v1.POST("/user/login",middleware.ValidateLogin(usecase.UserLogin{}),func(context *gin.Context) {
			webServiceHandler.LoginUser(context)
		})
		v1.GET("/user/me",middleware.JwtGinAuthUser().MiddlewareFunc(), func(c *gin.Context) {
			webServiceHandler.GetMe(c)
		})
	}


	router.Run(":3000")
}