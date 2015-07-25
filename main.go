package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/otiai10/gosseract"
	r "gopkg.in/dancannon/gorethink.v1"
)

func main() {
	// verify tesseract is available...
	goss, err := gosseract.NewClient()
	if err != nil {
		log.Fatalln(err.Error())
	}

	session, err := r.Connect(r.ConnectOpts{
		Address:  "localhost:28015",
		Database: "test",
	})
	if err != nil {
		log.Fatalln(err.Error())
	}

	mux := mux.NewRouter()
	mux.Handle("/post/new", createImageHandler(goss, NewPostingStore(session))).
		Methods("POST")

	mw := negroni.Classic()
	mw.UseHandler(mux)
	mw.Run(":1234")
}

func unmarshal(r io.Reader, v interface{}) (interface{}, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, v)
	if err != nil {
		return nil, err
	}
	return v, nil
}
