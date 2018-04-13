package domain

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

type UserRepository interface {
	StoreUser(user User)(User,error)
	GetByEmail(email string)(User,error)
	UpdateDevice(userId bson.ObjectId,deviceList []UserDevice)
	FindUserById(id string)(User,error)
}



type Data struct {
	Data		string	`json:"data" bson:"data,omitempty"`
	Code 		string	`json:"-"`
	Verified	bool	`json:"verified"`
}


type UserDevice struct {
	Agent 				string		`json:"agent"`
	Ip					string		`json:"ip"`
	DeviceId 			string		`json:"device_id" bson:"device_id"`
	Blocked 			bool		`json:"blocked"`
	ClientToken 		string		`json:"-" bson:"client_token"`
	Verified			bool		`json:"verified"`
	VerificationCode	int			`json:"-" bson:"verification_code"`
	IsLogin 			bool		`json:"is_login" bson:"is_login"`
	LastActive 			time.Time	`json:"last_active" bson:"last_active"`
}

type User struct {
	ID 					bson.ObjectId 	`json:"id" bson:"_id,omitempty"`
	ProfilePictureUrl	string			`json:"profile_picture_url" bson:"profile_picture_url"`
	Password			string			`json:"-"`
	TmpPassword 		string			`json:"-" bson:"tmp_password,omitempty"`
	Name				string			`json:"name"`
	Email				Data			`json:"email"`
	PhoneNumber			Data			`json:"phone_number" bson:"phone_number"`
	IsSuspend 			bool			`json:"is_suspend" bson:"is_suspend"`
	VerifiedUser		bool			`json:"verified_user" bson:"verified_user"`
	DeviceList			[]UserDevice	`json:"-" bson:"device_list"`
	CreatedAt 			time.Time		`json:"created_at" bson:"created_at"`
	UpdatedAt 			time.Time		`json:"updated_at" bson:"updated_at"`
	AccessToken 		string			`json:"access_token" bson:"-"`
}