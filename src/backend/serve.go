package flicq

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"

	"appengine"
	"appengine/datastore"
)

const dataKind = "Shot"

type Shot struct {
	ID    string `datastore:""`
	Foo   string
	Bar   string
	Duck  string
	Goose string
}

func init() {
	r := mux.NewRouter().PathPrefix("/api/").Subrouter()

	r.Handle("/add", appHandler(addshot)).Methods("POST")
	r.Handle("/get", appHandler(getshot)).Methods("GET")
	//This create REST API is just to test if the data inserts properly.
	r.Handle("/create", appHandler(create)).Methods("POST")
	http.Handle("/api/", r)
}

func getshot(w io.Writer, r *http.Request) error {
	c := appengine.NewContext(r)

	shots := []Shot{}
	keys, err := datastore.NewQuery(dataKind).GetAll(c, &shots)
	if err != nil {
		return appErrorf(http.StatusBadRequest, "get all content failed", err)
	}

	for i, k := range keys {
		shots[i].ID = k.Encode()
	}
	return json.NewEncoder(w).Encode(&shots)
}

func create(w io.Writer, r *http.Request) error {
	c := appengine.NewContext(r)
	shot := Shot{
		Foo:   "1.0",
		Bar:   "1.0",
		Duck:  "1.0",
		Goose: "1.0",
	}

	key := datastore.NewIncompleteKey(c, dataKind, nil)
	key, err := datastore.Put(c, key, &shot)
	if err != nil {
		return fmt.Errorf("create shot data: %v", err)
	}

	shot.ID = key.Encode()
	return json.NewEncoder(w).Encode(&shot)
}

func addshot(w io.Writer, r *http.Request) error {
	c := appengine.NewContext(r)

	shot := Shot{}
	err := json.NewDecoder(r.Body).Decode(&shot)
	if err != nil {
		return appErrorf(http.StatusBadRequest, "decode list: %v", err)
	}

	key := datastore.NewIncompleteKey(c, dataKind, nil)
	key, err = datastore.Put(c, key, &shot)
	if err != nil {
		return fmt.Errorf("add shot data: %v", err)
	}

	shot.ID = key.Encode()
	//We don't need to send the data back in latter days. This is just for
	//informational purposes.
	return json.NewEncoder(w).Encode(&shot)
}
