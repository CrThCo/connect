package models

import (
	"log"
	"sync"
	"time"

	"github.com/astaxie/beego"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// MongoConnection struct
type MongoConnection struct {
	ConnectionString string
	DatabaseName     string
	Instance         *mgo.Database
}

var (
	instance *MongoConnection
	once     sync.Once
)

// GetMongo function return the MongoConnection reference
func GetMongo() *MongoConnection {
	once.Do(func() {
		instance = &MongoConnection{
			ConnectionString: beego.AppConfig.String("DBServer"),
			DatabaseName:     "connect-test",
		}
		instance.Connect()
	})
	return instance
}

// Connect method
func (db *MongoConnection) Connect() {
	// for {
	// 	client, err := mongo.Connect(context.Background(), db.ConnectionString, nil)
	// 	if err != nil {
	// 		log.Printf("Mongo Connection Error:: %v :: retry in 30 seconds", err)
	// 		select {
	// 		case <-time.After(30 * time.Second):
	// 			continue
	// 		}
	// 	}
	// 	log.Println("Connection established with mongodb!")
	// 	db.Instance = client.Database(db.DatabaseName)
	// 	break
	// }

	for {
		session, err := mgo.Dial(db.ConnectionString)

		if err != nil {
			log.Printf("Mongo Connection Error:: %v", err)
			select {
			case <-time.After(30 * time.Second):
				continue
			}
		}
		log.Println("Connection established with mongodb!")
		db.Instance = session.DB(db.DatabaseName)
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

// FindByID Method
func (db *MongoConnection) FindByID(collection string, id string) *mgo.Query {
	return db.Instance.C(collection).FindId(bson.ObjectIdHex(id))
}

// UpdateByID method
func (db *MongoConnection) UpdateByID(collection string, id string, doc interface{}) error {
	return db.Instance.C(collection).UpdateId(bson.ObjectIdHex(id), doc)
}
