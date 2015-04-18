package flicq

import (
	"appengine/datastore"
	"log"
)

import "github.com/GoogleCloudPlatform/go-endpoints/endpoints"

const dataKind = "Shot"

type Shot struct {
	ID             *datastore.Key `json:"id" datastore:"-"`
	Counter        string         `json:"counter" datastore:",noindex" endpoints:"List"`
	Quaternion0    string         `json:"q0"`
	Quaternion1    string         `json:"q1"`
	Quaternion2    string         `json:"q2"`
	Quaternion3    string         `json:"q3"`
	AccelerometerX string         `json:"aX"`
	AccelerometerY string         `json:"aY"`
	AccelerometerZ string         `json:"aZ"`
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
	register("Create", "FlicqEndpointService.Shots.Create", "POST", "shots", "Create a shot info with random data")
	endpoints.HandleHTTP()
}

func (service *FlicqEndpointService) List(c endpoints.Context, r *FlicqRequest) (*Shots, error) {

	if r.Limit <= 0 {
		r.Limit = 10
	}

	q := datastore.NewQuery(dataKind).Limit(r.Limit)
	shots := make([]*Shot, 0, r.Limit)
	keys, err := q.GetAll(c, &shots)
	if err != nil {
		return nil, err
	}
	for i, id := range keys {
		shots[i].ID = id
	}

	return &Shots{shots}, nil
}

func (service *FlicqEndpointService) Add(c endpoints.Context, shot *Shot) error {

	key := datastore.NewIncompleteKey(c, dataKind, nil)
	_, err := datastore.Put(c, key, shot)
	return err
}

func (service *FlicqEndpointService) Create(c endpoints.Context) error {
	shot := Shot{
		Quaternion0:    "0",
		Quaternion1:    "1",
		Quaternion2:    "2",
		Quaternion3:    "3",
		AccelerometerX: "4",
		AccelerometerY: "5",
		AccelerometerZ: "6",
	}

	key := datastore.NewIncompleteKey(c, dataKind, nil)
	_, err := datastore.Put(c, key, &shot)
	return err
}
