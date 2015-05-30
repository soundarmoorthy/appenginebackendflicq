package flicq

import (
	"appengine/datastore"
	"log"
	"time"
)

import "github.com/GoogleCloudPlatform/go-endpoints/endpoints"

const dataKind = "Shot"

type Shot struct {
	KEY   *datastore.Key `json:"key" datastore:"-"`
	ID    string
	Items []float32
	Time  time.Time
}

type Shots struct {
	Items []*Shot `json:"items"`
}

type FlicqRequest struct {
	Limit int `json:"limit" endpoints:"d=10"`
}

type FlicqEndpointService struct {
}

func init() {
	service := &FlicqEndpointService{}
	api, err := endpoints.RegisterService(service, "flicq", "v1", "Flicq Backend Data managemenet API", true)
	if err != nil {
		log.Fatalf("Register Service : %v", err)
	}

	register := func(orig, name, httpMethod, path, desc string) {
		method := api.MethodByName(orig)
		if method == nil {
			log.Fatalf("Missing method : %s", orig)
		}
		info := method.Info()
		info.Name, info.HTTPMethod, info.Path, info.Desc = name, httpMethod, path, desc
	}

	register("Add", "FlicqEndpointService.Shots.Add", "PUT", "shots", "Add a shot")
	register("List", "FlicqEndpointService.Shots.List", "GET", "shots", "List all the shots")
	endpoints.HandleHTTP()
}

func (service *FlicqEndpointService) List(c endpoints.Context, r *FlicqRequest) (*Shots, error) {
	if r.Limit <= 0 {
		r.Limit = 10
	}
	q := datastore.NewQuery(dataKind)
	shots := make([]*Shot, 0, r.Limit)
	keys, err := q.GetAll(c, &shots)
	if err != nil {
		return nil, err
	}
	for i, key := range keys {
		shots[i].KEY = key
	}
	return &Shots{shots}, nil
}

func (service *FlicqEndpointService) Add(context endpoints.Context, shot *Shot) error {
	key := datastore.NewIncompleteKey(context, dataKind, nil)
	_, result := datastore.Put(context, key, shot)
	return result
}
