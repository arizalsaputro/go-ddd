package config

import "os"

type DB struct {
	Host 		string
	Database 	string
	User 		string
	Password	string
}

type Config struct {
	Db DB
	SecretJwt string
}

var Conf Config

func init()  {
	db := DB{}
	db.Host = os.Getenv("MONGO_HOST")
	if db.Host == ""{
		db.Host = "localhost:27017"
	}
	db.Database = os.Getenv("MONGO_DATABASE")
	if db.Database == ""{
		db.Database = "lowcost"
	}
	db.User = os.Getenv("MONGO_USER")
	if db.User == ""{
		db.User = "lowcostuser"
	}
	db.Password = os.Getenv("MONGO_PASSWORD")
	if db.Password == ""{
		db.Password = "bukaajaditrello"
	}
	Conf.Db = db

	Conf.SecretJwt = os.Getenv("JWT_SECRET")
	if Conf.SecretJwt == ""{
		Conf.SecretJwt = "C]:zJF2vpk_yP@dJ"
	}
}