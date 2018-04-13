package mongo

import (
	"log"
	"github.com/arizalsaputro/go-ddd/domain"
	"github.com/arizalsaputro/go-ddd/infrastructure"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type MongoRepo struct{}

func (MongoRepo) StoreUser(user domain.User)(domain.User,error) {
	ms := infrastructure.MongoSession()
	defer ms.Session.Close()

	err := ms.UserCol().Insert(&user)
	log.Println("Inserted user into db")
	return user,err
}

func (MongoRepo) GetByEmail(email string)(domain.User,error)  {
	ms := infrastructure.MongoSession()
	defer ms.Session.Close()
	var user domain.User
	err := ms.UserCol().Find(bson.M{"email.data":email}).One(&user)
	return user,err
}

func (MongoRepo)UpdateDevice(userId bson.ObjectId,deviceList []domain.UserDevice)  {
	ms := infrastructure.MongoSession()
	defer ms.Session.Close()
	log.Println("Updating device")
	err := ms.UserCol().UpdateId(userId,bson.M{"$set":bson.M{"device_list":deviceList,"updated_at":time.Now().UTC()}})
	if err != nil{
		log.Println("error updating device")
	}
}

func (MongoRepo) FindUserById(id string) (domain.User,error) {
	ms := infrastructure.MongoSession()
	defer ms.Session.Close()
	var user domain.User
	err := ms.UserCol().FindId(bson.ObjectIdHex(id)).One(&user)
	return user,err
}

func (MongoRepo) FindUserByEmailAddress(email string)  {

}

