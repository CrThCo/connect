package couch

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/astaxie/beego"
)

var httpapi *HTTPAPI
var userDB string

func init() {
	httpapi = &HTTPAPI{
		URL: beego.AppConfig.String("DBServer"),
	}
	userDB = "connect_users"
}

// CreateDB function will use to create the default database.
func CreateDB() bool {
	httpapi.Endpoint = userDB
	resp, err := httpapi.Put()

	if err != nil {
		log.Printf("Create Database: %v", err)
		return false
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Create DB: Response parsing error: %v", err)
		return false
	}

	log.Printf("Data: %s", string(data))
	return true
}

// CreateUser function
func CreateUser(id string, user interface{}) bool {
	httpapi.Endpoint = fmt.Sprintf("%s/%s", userDB, id)

	// res, err := httpapi.Get()

	// if err != nil {
	// 	log.Printf("User Exits: %v", err)
	// 	return false
	// }

	httpapi.Data = user
	resp, err := httpapi.Put()

	if err != nil {
		log.Printf("New User: %v", err)
		return false
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("New User: Response parsing error: %v", err)
		return false
	}

	log.Printf("New User: %s", string(data))
	return true
}
