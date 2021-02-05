package database

import (
	"github.com/globalsign/mgo"
	"log"
	"os"
)

// Mongo struct
type Mongo struct {
	Posts *mgo.Collection
}

type MongoConfig struct {
	mongoHost string
	mongoPort string
	mongoDb   string
	username  string
	password  string
	msgCol    string
}

// CreateConn returns a mongo db connection object
func CreateConn() *Mongo {
	var mongoConfig = MongoConfig{
		mongoHost: getEnv("MONGO_HOST", "localhost"),
		mongoPort: getEnv("MONGO_PORT", "27017"),
		mongoDb:   getEnv("MONGO_DB", "jobscrapper"),
		username:  getEnv("MONGO_USER", "admin"),
		password:  getEnv("MONGO_PASS", "pass"),
		msgCol:    getEnv("MONGO_MSG_COLLECTION", "msgs"),
	}

	info := &mgo.DialInfo{
		Addrs:    []string{mongoConfig.mongoHost},
		Database: mongoConfig.mongoDb,
		Username: mongoConfig.username,
		Password: mongoConfig.password,
	}

	s, err := mgo.DialWithInfo(info)
	if err != nil {
		log.Printf("ERROR connecting mongo, %s ", err.Error())
		os.Exit(999)
	}

	db := s.DB("jobscrapper")

	return &Mongo{
		Posts: db.C("posts"),
	}

}

func getEnv(key, ifnotset string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return ifnotset
}