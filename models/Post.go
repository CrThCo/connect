package models

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"image/jpeg"
	"log"
	"os"
	"time"

	"github.com/astaxie/beego"
	"gopkg.in/mgo.v2/bson"
)

const (
	postCollection = "posts"
	voteCollection = "votes"
)

type VoteOptions struct {
	Name string  `json:"name"`
	Value string `json:"value"`
}

type VoteStruct struct {
	Options []VoteOptions `json:"options"`
	Image string `json:"image"`

}

// Vote struct to store vote per 
type Vote struct {
	ID       bson.ObjectId `bson:"_id" json:"id"`
	PostID   bson.ObjectId `bson:"post_id" json:"vote_id"`
	VoterID  bson.ObjectId `bson:"voter_id" json:"voter_id"`
	Vote     []string    `bson:"votes" json:"votes"`
	CastedAt time.Time `bson:"created_at" json:"created_at"`
}

// Post struct to map post
type Post struct {
	ID        bson.ObjectId `bson:"_id" json:"id"`
	Content   string        `bson:"content" json:"content"`
	Hash      string        `bson:"hash" json:"hash"`
	Image     string        `bson:"image" json:"image"`
	Verified  bool          `bson:"verified" json:"verified"`
	Poster    string        `bson:"poster" json:"poster"`
	VoteCount int           `bson:"vote_count" json:"vote_count"`
	CreatedAt time.Time     `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time     `bson:"updated_at" json:"updated_at"`
}

// SaveImage method
func (p *Post) SaveImage() error {

	d, err := base64.StdEncoding.DecodeString(p.Image)
	if err != nil {
		return err
	}

	// decoded and error
	dimg, err := jpeg.Decode(bytes.NewReader(d))
	// d, err := base64.StdEncoding.DecodeString(p.Image)
	if err != nil {
		return err
	}
	filename := fmt.Sprintf("image_%d.jpg", time.Now().Nanosecond())
	// path
	path := fmt.Sprintf("%s/%s", beego.AppConfig.String("FileStoragePath"), filename)

	f, err := os.Create(path)
	defer f.Close()

	opts := jpeg.Options{
		Quality: 100,
	}
	if err := jpeg.Encode(f, dimg, &opts); err != nil {
		return err
	}
	// if err := ioutil.WriteFile(path, []byte(""), 0777); err != nil {
	// 	return err
	// }
	// file and error
	// f, err := os.Create(path)
	// if err != nil {
	// 	return err
	// }
	// defer f.Close()
	// if _, err := f.Write(d); err != nil {
	// 	return err
	// }

	// if err := f.Sync(); err != nil {
	// 	return err
	// }

	p.Image = filename
	return nil
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
	if err := GetMongo().Insert(postCollection, p); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// GetByUser method
func (p *Post) GetByUser() ([]*Post, error) {
	var posts []*Post
	err := GetMongo().Find(postCollection, nil).Sort("-$natural").All(&posts)
	if err != nil {
		log.Printf("Post -> GetByUser %v", err)
		return nil, errors.New("unable to reterive posts for user")
	}
	return posts, nil
}

// GetByID function
func GetByID(uid string) (*Post, error) {
	var post *Post
	err := GetMongo().FindByID(postCollection, uid).One(&post)
	if err != nil {
		log.Printf("Post -> GetByID %v", err)
		return nil, errors.New("unable to retrieve posts for user")
	}
	return post, nil
}

// Insert post
func (v *VoteStruct) AddVote(postid, voterid bson.ObjectId) error {
	if (!postid.Valid() || !voterid.Valid()) {
		return errors.New("Post id or voter id not valid")
	}
	if len(v.Options) == 0 {
		return errors.New("No vote options")
	}
	var vote *Vote
	vote.Vote = make([]string, 0, len(v.Options))
	vote.ID = bson.NewObjectId()
	vote.CastedAt = time.Now().UTC()
	vote.PostID = postid
	vote.VoterID = voterid
	for _, o := range v.Options {
		log.Println(o)
		vote.Vote = append(vote.Vote, o.Name)
	}
	
	if err := GetMongo().Insert(voteCollection, v); err != nil {
		log.Println(err)
		return err
	}
	return nil
}