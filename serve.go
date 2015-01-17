package shotcontent

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"appengine"
	"appengine/datastore"
)

const dataKind = "Shot"

type Shot struct {
	ID        string `datastore:"-"`
	Timestamp time.Time
	Pitch     float32
	Roll      float32
	Yaw       float32
}

func init() {
	r := mux.NewRouter().PathPrefix("/api/").Subrouter()

	r.Handle("/add", appHandler(addshot)).Methods("POST")
	r.Handle("/get", appHandler(getshot)).Methods("GET")
	http.Handle("/api/", r)
}

func getshot(w io.Writer, r *http.Request) error {
	c := appengine.NewContext(r)

	shotdata := []Shot{}
	keys, err := datastore.NewQuery(dataKind).GetAll(c, &shotdata)
	if err != nil {
		return appErrorf(http.StatusBadRequest, "get all content failed", err)
	}

	for i, k := range keys {
		shotdata[i].ID = k.Encode()
	}
	return json.NewEncoder(w).Encode(shotdata)
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
		return fmt.Errorf("create list : %v", err)
	}

	shot.ID = key.Encode()
	return json.NewEncoder(w).Encode(shot)
}
