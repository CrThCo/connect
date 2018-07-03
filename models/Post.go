package models

import (
	"crypto/sha256"
	"fmt"
	"log"
	"time"

	"gopkg.in/mgo.v2/bson"
)

var postCollecton string

func init() {
	postCollecton = "posts"
}

// Post struct to map post
type Post struct {
	ID        bson.ObjectId `bson:"_id" json:"id"`
	Content   string        `bson:"content" json:"content"`
	Hash      string        `bson:"hash" json:"hash"`
	Verified  bool          `bson:"verified" json:"verified"`
	Poster    string        `bson:"poster" json:"poster"`
	Votes     int           `bson:"votes" json:"votes"`
	CreatedAt time.Time     `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time     `bson:"updated_at" json:"updated_at"`
}

// Insert post
func (p *Post) Insert() error {
	p.ID = bson.NewObjectId()
	p.CreatedAt = time.Now().UTC()
	p.UpdatedAt = time.Now().UTC()
	h := sha256.New()
	_, err := h.Write([]byte(p.Content))
	if err != nil {
		return err
	}
	p.Hash = fmt.Sprintf("%x", h.Sum(nil))
	if err := GetMongo().Insert(postCollecton, p); err != nil {
		log.Println(err)
		return err
	}
	return nil
}
