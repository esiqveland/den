package main

import (
	r "gopkg.in/dancannon/gorethink.v1"
	"log"
	"time"
	"errors"
)

const postingsTable = "POSTINGS"

var ErrPostNotFound = errors.New("Post Not Found")

type PostingFetcher interface {
	Get(id string) (Posting, error)
}
type PostingCreator interface {
	Create(post Posting) (id string, err error)
}
type PostingStore interface {
	PostingFetcher
	PostingCreator
}

func NewPostingStore(session *r.Session) PostingStore {
	return &RethinkPosting{session: session}
}
type RethinkPosting struct {
	session *r.Session
}

// Get fetches a posting. Make sure to check for err.
// returns ErrPostNotFound for empty result.
func (self *RethinkPosting) Get(id string) (Posting, error) {
	post := Posting{}
	res, err := r.Table(postingsTable).Get(id).Run(self.session)
	if err != nil {
		return post, err
	}
	if res.IsNil() {
		return post, ErrPostNotFound
	}
	err = res.One(&post)
	if err != nil {
		return post, err
	}
	return post, nil
}

func (self *RethinkPosting) Create(post Posting) (string, error) {
	wr, err := r.Table(postingsTable).Insert(&post).RunWrite(self.session)
	if err != nil {
		return "", err
	}
	log.Printf("%v", wr.GeneratedKeys)
	return wr.GeneratedKeys[0], nil
}

type Posting struct {
	Id          string    `gorethink:"id,omitempty"`
	Title       string    `gorethink:"title"`
	Description string    `gorethink:"description,omitempty"`
	Images      Images    `gorethink:"images,omitempty":json:images"` // B64 encoded images
	Created     time.Time `gorethink:"created,omitempty"`
	Updated     time.Time `gorethink:"updated,omitempty"`
}

type Content struct {
	Created string `gorethink:"created,omitempty"`
	Origin  string `gorethink:"origin,omitempty"` // Creating system
	Data    string `gorethink:"id,omitempty"`     //
}
type Image struct {
	Content   Contents `gorethink:"contents"`
	Data      string   `gorethink:"data"`
	Footprint string   `gorethink:"footprint,omitempty"`
}

type Contents []Content
type Images []Image
