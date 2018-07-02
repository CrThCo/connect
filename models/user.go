package models

import (
	"errors"
	"log"

	"golang.org/x/crypto/bcrypt"
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

// Auth struct
type Auth struct {
	Username string `bson:"username" json:"username"`
	Password string `bson:"password" json:"password"`
}

// User struct
type User struct {
	ID          bson.ObjectId `bson:"_id" json:"id"`
	Username    string        `bson:"username" json:"username"`
	Email       string        `bson:"email" json:"email"`
	Password    string        `bson:"password" json:"password"`
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
	//invalid user
	if u.Password == "" {
		log.Printf("New User Error:: %v", "Invalid user")
		return ""
	}
	//bcrypt hashed password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	u.Password = string(hashedPassword)
	// make sure there is no duplication when adding users
	i, err := db.Find(collection, bson.M{"$or": []bson.M{
		{"username": u.Username},
		{"email": u.Email},
	}}).Count()
	if i > 0 {
		log.Printf("New User Error:: %v", "username or email is already in use")
		return ""
	}
	err = db.Insert(collection, u)
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

// GetUserByCredentials for fetching with username and password
func GetUserByCredentials(username, password string) (string, error) {
	u := &User{Username: username}
	err := db.Find(collection, u).One(&u)
	if err != nil {
		log.Printf("Get User error :: %v ", err)
		return "", errors.New("User doesn't exist")
	}
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		log.Printf("Get User error :: %v ", err)
		return "", errors.New("Bad Password")
	}
	return string(u.ID), nil
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
