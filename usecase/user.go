package usecase

import (
	"github.com/arizalsaputro/go-ddd/domain"
	"time"
	"gopkg.in/mgo.v2/bson"
)

type UserRegister struct {
	Email				string			`json:"email" form:"email" binding:"required,email" `
	Password			string			`json:"password" form:"password" binding:"required,min=6" `
	Name				string			`json:"name" form:"name" binding:"required"`
	Code				string			`json:"-" binding:"omitempty"`
	Device 				domain.UserDevice	`json:"-" binding:"omitempty"`
}

type UserLogin struct {
	Email 				string			`json:"email" form:"email" binding:"required,email"`
	Password			string			`json:"password" form:"password" binding:"required,min=6" `
}

type UserInteractor struct {
	UserRepository domain.UserRepository
}

func (interactor *UserInteractor)GetMe(id string)(domain.User,error)  {
	user,err := interactor.UserRepository.FindUserById(id)
	return user,err
}

func (interactor *UserInteractor)RegisterUser (user UserRegister) (domain.User,error){
	newUser,err := interactor.UserRepository.StoreUser(mapToDomainUser(user))
	if err != nil{
		return newUser,err
	}
	theUser,err := interactor.UserRepository.FindUserById(newUser.ID.Hex())
	return theUser,err
}

func (interactor *UserInteractor)LoginUser (user UserLogin)(domain.User,error)  {
	usr,err := interactor.UserRepository.GetByEmail(user.Email)
	return usr,err
}

func (interactor *UserInteractor)UpdateDevice(userId bson.ObjectId,deviceList []domain.UserDevice)() {
	interactor.UserRepository.UpdateDevice(userId,deviceList)
}

func mapToDomainUser(user UserRegister) (dUser domain.User) {
	dUser.ID = bson.NewObjectId()
	dUser.DeviceList = []domain.UserDevice{user.Device}
	dUser.Name = user.Name
	dUser.Email = domain.Data{Data:user.Email,Code:user.Code}
	dUser.Password = user.Password
	dUser.UpdatedAt = time.Now().UTC()
	dUser.CreatedAt = time.Now().UTC()
	return
}