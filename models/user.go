package models

import (
	"errors"
	"fmt"
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
	Email    string `bson:"email" json:"email"`
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

// Insert method
func (u *User) Insert() error {
	u.ID = bson.NewObjectId()
	// Password should not be empty.
	// TODO :: Add email, usename validation rules
	if u.Password == "" {
		return errors.New("user password should not be empty")
	}

	where := bson.M{
		"$or": []bson.M{
			{"username": u.Username},
			{"email": u.Email},
		},
	}
	// make sure there is no duplication when adding users
	count, err := db.Find(collection, where).Count()
	if err != nil {
		log.Printf("Error! duplication user %v", err)
		return errors.New("internal server error, please try again")
	}

	if count > 0 {
		return fmt.Errorf("username %s or email %s already in use", u.Username, u.Email)
	}
	//bcrypt hashed password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error! generating password hash %v", err)
		return errors.New("internal server error, please try again")
	}
	u.Password = string(hashedPassword)
	if err = db.Insert(collection, u); err != nil {
		log.Printf("New User Error:: %v", err)
		return errors.New("something went wrong while creating your account")
	}

	u.Password = ""
	return nil
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
func GetUserByCredentials(email, password string) (string, error) {
	u := &User{Email: email}
	log.Println(email)
	err := db.Find(collection, bson.M{"email": email}).One(&u)
	if err != nil {
		log.Printf("Get User error 1 :: %v ", err)
		return "", errors.New("User doesn't exist")
	}

	log.Println(u.Password)
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		log.Printf("Get User error 2 :: %v ", err)
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
