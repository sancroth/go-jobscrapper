package models

import (
	"time"

	"../database"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

// Post struct
type Post struct {
	ID        bson.ObjectId `bson:"_id,omitempty"`
	Path      string
	Name      string
	Host      string
	Title     string
	Salary    string
	Position  string
	Company   string
	Apply     string
	Processed bool
	Created   time.Time
	Update    time.Time
}

// PostsHandler struct
type PostsHandler struct {
	posts *mgo.Collection
}

// NewPostsHandler returns a handler for managing Post structs
func NewPostsHandler() *PostsHandler {
	return &PostsHandler{
		posts: database.CreateConn().Posts,
	}
}

func (ph *PostsHandler) FindPosts(limit int) ([]*Post, error) {
	var ps []*Post
	return ps, ph.posts.Find(bson.M{
		"processed": false,
	}).Limit(limit).Sort("-created").All(&ps)
}

func (ph *PostsHandler) Processed(pa []*Post) error {
	batch := ph.posts.Bulk()
	for _, p := range pa {
		batch.UpdateAll(bson.M{"_id": p.ID}, bson.M{
			"$set": bson.M{"processed": true, "updated": time.Now()},
		})
	}
	_,err := batch.Run()
	if err !=nil{
		return err
	}
	return nil
}

func (ph *PostsHandler) GetPostsCount(name, path string) (int, error) {
	return ph.posts.Find(bson.M{
		"name": name,
		"path": path,
	}).Count()
}

func (ph *PostsHandler) SavePost(p *Post) error {
	p.Created = time.Now()
	p.Processed = false
	return ph.posts.Insert(p)
}
