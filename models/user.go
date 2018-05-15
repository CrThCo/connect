package models

import (
	"errors"
	"log"

	"gopkg.in/mgo.v2/bson"
)

var (
	db         *MongoConnection
	collection string
)

func init() {
	db = GetMongo()
	collection = "user"
}

// User struct
type User struct {
	ID          bson.ObjectId `bson:"_id" json:"id"`
	Username    string        `bson:"username" json:"username"`
	Email       string        `bson:"email" json:"email"`
	FirstName   string        `bson:"first_name" json:"first_name"`
	LastName    string        `bson:"last_name" json:"last_name"`
	Image       string        `bson:"image" json:"image"`
	SocialMedia []SocialMedia `bson:"social_media" json:"social_media"`
}

// SocialMedia struct
type SocialMedia struct {
	Name string `bson:"name" json:"name"`
	Key  string `bson:"key" json:"key"`
	URL  string `bson:"url" json:"url"`
}

// AddUser function
func AddUser(u User) string {
	u.ID = bson.NewObjectId()
	err := db.Insert(collection, u)
	if err != nil {
		log.Printf("New User Error:: %v", err)
		return ""
	}
	return u.ID.Hex()
}

// GetUser function
func GetUser(uid string) (*User, error) {
	u := &User{}
	err := db.FindByID(collection, uid).One(&u)
	if err != nil {
		log.Printf("Get User error :: %v ", err)
		return nil, errors.New("User not exists")
	}
	return u, nil
}

// GetAllUsers function
func GetAllUsers() []User {
	users := []User{}
	err := db.Find(collection, nil).All(&users)
	if err != nil {
		log.Printf("User FindAll Error :: %v ", err)
		return nil
	}
	return users
}

// UpdateUser function
func UpdateUser(uid string, uu *User) (*User, error) {

	if err := db.UpdateByID(collection, uid, uu); err != nil {
		log.Printf("Update Error :: %v", err)
		return nil, errors.New("User Not Exist")
	}
	err := db.FindByID(collection, uid).One(&uu)
	if err != nil {
		return nil, err
	}
	return uu, nil
}
