package webcontrollers

import (
	"github.com/arizalsaputro/go-ddd/usecase"
	"github.com/arizalsaputro/go-ddd/domain"
	"github.com/gin-gonic/gin"
	"github.com/arizalsaputro/go-ddd/interface/webcontrollers/util"
	"net/http"
	"time"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

)

type UserInteractor interface {
	RegisterUser(user usecase.UserRegister) (domain.User,error)
	LoginUser(user usecase.UserLogin)(domain.User,error)
	UpdateDevice(userId bson.ObjectId,deviceList []domain.UserDevice)
	GetMe(id string)(domain.User,error)
}

type WebServiceHandler struct {
	UserInteractor UserInteractor
}

func (handler WebServiceHandler)GetMe(c *gin.Context)  {
	userId := c.MustGet("USER").(string)
	user,err := handler.UserInteractor.GetMe(userId)
	if err != nil{
		util.ServeError(c,http.StatusBadRequest,err.Error())
		return
	}
	util.ServeOk(c,http.StatusOK,user)
}

func (handler WebServiceHandler)RegisterUser (c *gin.Context) {
	form := c.MustGet("JSON").(usecase.UserRegister)
	agent := c.GetHeader("api-client")
	if agent == ""{
		agent = "etc"
	}
	deviceId := c.GetHeader("device-id")
	if deviceId == ""{
		deviceId = "123"
	}
	clientToken := c.GetHeader("client-token")
	ipAdd := util.GetClientIPByRequest(c.Request)
	form.Device = domain.UserDevice{Agent:agent,DeviceId:deviceId,Verified:true,IsLogin:true,LastActive:time.Now().UTC(),ClientToken:clientToken,Ip:ipAdd}
	user,err := handler.UserInteractor.RegisterUser(form)
	if err != nil{
		if mgo.IsDup(err){
			util.ServeError(c,http.StatusBadRequest,"email address already used")
			return
		}
		util.ServeError(c,http.StatusInternalServerError,err.Error())
		return
	}
	token,err := util.GenerateAccessToken(string(user.ID.Hex()),"user")
	if err != nil{
		util.ServeError(c,http.StatusInternalServerError,err.Error())
		return
	}
	user.AccessToken = token
	util.ServeOk(c,http.StatusCreated,user)
}

func (handler WebServiceHandler)LoginUser (c *gin.Context) {
	form := c.MustGet("JSON").(usecase.UserLogin)
	agent := c.GetHeader("api-client")
	if agent == ""{
		agent = "etc"
	}
	deviceId := c.GetHeader("device-id")
	if deviceId == ""{
		deviceId = "123"
	}
	clientToken := c.GetHeader("client-token")
	ipAdd := util.GetClientIPByRequest(c.Request)
	device := domain.UserDevice{Agent:agent,DeviceId:deviceId,ClientToken:clientToken,Ip:ipAdd,IsLogin:true,Verified:false,LastActive:time.Now().UTC()}
	user,err := handler.UserInteractor.LoginUser(form)
	if err != nil{
		util.ServeError(c,http.StatusBadRequest,err.Error())
		return
	}
	if util.CheckPassword(form.Password,user.Password){
		token,err := util.GenerateAccessToken(string(user.ID.Hex()),"user")
		if err != nil{
			util.ServeError(c,http.StatusInternalServerError,err.Error())
			return
		}
		user.AccessToken = token
		util.ServeOk(c,http.StatusOK,user)
		var exist = false
		for i:=0;i<len(user.DeviceList);i++{
			if user.DeviceList[i].DeviceId == device.DeviceId{
				user.DeviceList[i].LastActive = device.LastActive
				user.DeviceList[i].IsLogin = true
				user.DeviceList[i].Ip =device.Ip
				exist = true
				break
			}
		}
		if !exist {
			user.DeviceList = append(user.DeviceList, device)
		}
		handler.UserInteractor.UpdateDevice(user.ID,user.DeviceList)
	}else{
		util.ServeError(c,http.StatusBadRequest,"invalid password")
	}
}