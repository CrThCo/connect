package db

import (
	"errors"
	"fmt"
	"log"
	"regexp"

	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"
)

var (
	conn          *MongoConnection
	collection    string
	emailRegex    *regexp.Regexp
	usernameRegex *regexp.Regexp
)

func init() {
	conn = GetMongo()
	collection = "user"
	emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	usernameRegex = regexp.MustCompile(`^[a-z]+.*`)
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
	Password    string        `bson:"password,omitempty" json:"password"`
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

// Validate method
func (u *User) Validate() error {
	if len(u.Username) < 1 {
		return errors.New("username is missing")
	}
	if !emailRegex.MatchString(u.Email) {
		return errors.New("email address is invalid")
	}
	if !usernameRegex.MatchString(u.Username) {
		return errors.New("username should be starts with a-z")
	}
	if len(u.Username) < 5 || len(u.Username) > 32 {
		return errors.New("username should be between 5 to 32 charecters")
	}
	if len(u.Password) < 8 {
		return errors.New("password should be be 8 charecters long")
	}
	return nil
}

// Insert method
func (u *User) Insert() error {
	if err := u.Validate(); err != nil {
		return err
	}
	u.ID = bson.NewObjectId()
	where := bson.M{
		"$or": []bson.M{
			{"username": u.Username},
			{"email": u.Email},
		},
	}
	// make sure there is no duplication when adding users
	count, err := conn.Find(collection, where).Count()
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
	if err = conn.Insert(collection, u); err != nil {
		log.Printf("New User Error:: %v", err)
		return errors.New("something went wrong while creating your account")
	}

	u.Password = ""
	return nil
}

// GetUser function
func GetUser(uid string, u *User) error {
	err := conn.FindByID(collection, uid).One(&u)
	if err != nil {
		log.Printf("Get User error :: %v ", err)
		return errors.New("User not exists")
	}
	return nil
}

// GetUserByEmail for fetching with username and password
func GetUserByEmail(email string) (*User, error) {
	u := &User{Email: email}
	err := conn.Find(collection, bson.M{"email": email}).One(&u)
	if err != nil {
		log.Printf("Get User error 1 :: %v ", err)
		return nil, errors.New("User doesn't exist")
	}
	return u, nil
}

// GetUserByCredentials for fetching with username and password
func GetUserByCredentials(email, password string) (string, error) {
	u := &User{Email: email}
	err := conn.Find(collection, bson.M{"email": email}).One(&u)
	if err != nil {
		log.Printf("Get User error 1 :: %v ", err)
		return "", errors.New("User doesn't exist")
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		log.Printf("Get User error 2 :: %v ", err)
		return "", errors.New("Bad Password")
	}
	return u.ID.Hex(), nil
}

// GetUserList function
func GetUserList() ([]User, error) {
	users := []User{}
	err := conn.Find(collection, nil).All(&users)
	if err != nil {
		log.Printf("User fetch error :: %v ", err)
		return nil, err
	}
	return users, nil
}

// UpdateUser function
func UpdateUser(uid string, uu *User) error {
	if err := conn.UpdateByID(collection, uid, uu); err != nil {
		return fmt.Errorf("User Not Exist %v", err)
	}
	return nil
}
