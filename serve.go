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
	ID       string `datastore:"-"`
	time_abs time.Time
	time_rel int32
	Ax       float32
	Ay       float32
	Az       float32
	Mx       float32
	My       float32
	Mz       float32
	Gx       float32
	Gy       float32
	Gz       float32
	quat0    float32
	quat1    float32
	quat2    float32
	quat3    float32
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

func create(w io.Writer, r *http.Request) error {
	c := appengine.NewContext(r)
	shot := Shot{
		time_abs: time.Now(),
		time_rel: 1000,
		Ax:       10.0,
		Ay:       12.0,
		Az:       13.0,
		Mx:       14.0,
		My:       15.0,
		Mz:       16.0,
		Gx:       17.0,
		Gy:       18.0,
		Gz:       19.0,
		quat0:    1.0,
		quat1:    2.0,
		quat2:    3.0,
		quat3:    4.0,
	}

	key := datastore.NewIncompleteKey(c, dataKind, nil)
	key, err := datastore.Put(c, key, &shot)
	if err != nil {
		return fmt.Errorf("create list : %v", err)
	}

	shot.ID = key.Encode()
	return json.NewEncoder(w).Encode(shot)
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
