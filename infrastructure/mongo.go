package infrastructure

import (
	"gopkg.in/mgo.v2"
	"time"
	"github.com/arizalsaputro/go-ddd/config"
	"log"
)

type MgoSession struct {
	Session *mgo.Session
}

var mongoSession *mgo.Session


func ConnectMongo() error {
	var err error
	/*
	 * Create the connection
	 */
	log.Println("Connecting to mongo database...")

	// We need this object to establish a session to our MongoDB.
	mongoDBDialInfo := &mgo.DialInfo{
		Addrs:    []string{config.Conf.Db.Host},
		Timeout:  60 * time.Second,
		Database: config.Conf.Db.Database,
		Username: config.Conf.Db.User,
		Password: config.Conf.Db.Password,
	}

	// Create a session which maintains a pool of socket connections
	// to our MongoDB.
	mongoSession, err = mgo.DialWithInfo(mongoDBDialInfo)
	if err != nil {
		log.Fatalf("CreateSession: %s\n", err)
	}

	mongoSession.SetMode(mgo.Monotonic, true)
	go ensureIndex(mongoSession)
	return nil
}

func ensureIndex(ms *mgo.Session)  {
	userIndex := mgo.Index{
		Key: 	[]string{"email.data","phone_number.data"},
		Unique:		true,
		DropDups:   true,
		Sparse:		true,
		Background:	true,
	}
	ms.DB(config.Conf.Db.Database).C("users").EnsureIndex(userIndex)

}

func MongoSession() *MgoSession {
	return &MgoSession{mongoSession.Copy()}
}

func (ms *MgoSession) UserCol() *mgo.Collection {
	return ms.Session.DB(config.Conf.Db.Database).C("users")
}
