package db

import (
	"log"
	"os"
	"sync"
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// MongoConnection struct
type MongoConnection struct {
	DailInfo *mgo.DialInfo
	Instance *mgo.Database
}

var (
	instance *MongoConnection
	once     sync.Once
)

// GetMongo function return the MongoConnection reference
func GetMongo() *MongoConnection {
	once.Do(func() {
		info, err := mgo.ParseURL(os.Getenv("MONGODB_URI"))
		if err != nil {
			log.Fatalln("Invalid MONGODB_URI")

		}
		instance = &MongoConnection{
			DailInfo: info,
		}
		instance.Connect()
	})
	return instance
}

// Connect method
func (db *MongoConnection) Connect() {
	for {
		session, err := mgo.DialWithInfo(db.DailInfo)
		if err != nil {
			log.Printf("Mongo Connection Error:: %v", err)
			select {
			case <-time.After(30 * time.Second):
				continue
			}
		}
		log.Println("Connection established with mongodb!")
		db.Instance = session.DB("")
		break
	}

}

// Insert method
func (db *MongoConnection) Insert(collection string, document interface{}) error {
	return db.Instance.C(collection).Insert(document)
}

// Find Method
func (db *MongoConnection) Find(collection string, filter interface{}) *mgo.Query {
	return db.Instance.C(collection).Find(filter)
}

// Pipe method
func (db *MongoConnection) Pipe(collection string, pipe interface{}) *mgo.Pipe {
	return db.Instance.C(collection).Pipe(pipe)
}

// FindByID Method
func (db *MongoConnection) FindByID(collection string, id string) *mgo.Query {
	return db.Instance.C(collection).FindId(bson.ObjectIdHex(id))
}

// UpdateByID method
func (db *MongoConnection) UpdateByID(collection string, id string, doc interface{}) error {
	return db.Instance.C(collection).UpdateId(bson.ObjectIdHex(id), doc)
}
