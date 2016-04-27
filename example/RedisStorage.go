package example

import (
	"errors"
	"fmt"
	"github.com/RangelReale/osin"
	"gopkg.in/redis.v3"
  	"encoding/json"
)

var client = redis.NewClient(&redis.Options{
  Addr: "192.168.64.2:6379",
  Password: "", // no password set
  DB: 0, // use default DB
})

type TestStorage struct {
	clients   map[string]osin.Client
	authorize map[string]*osin.AuthorizeData
	access    map[string]*osin.AccessData
	refresh   map[string]string
}

func NewTestStorage() *TestStorage {
	r := &TestStorage{
		clients:   make(map[string]osin.Client),
		authorize: make(map[string]*osin.AuthorizeData),
		access:    make(map[string]*osin.AccessData),
		refresh:   make(map[string]string),
	}

	r.clients["1234"] = &osin.DefaultClient{
		Id:          "1234",
		Secret:      "aabbccdd",
		RedirectUri: "http://localhost:14000/appauth",
	}

	return r
}

func (s *TestStorage) Clone() osin.Storage {
	return s
}

func (s *TestStorage) Close() {
}

func (s *TestStorage) GetClient(id string) (osin.Client, error) {
	fmt.Printf("GetClient: %s\n", id)
	if c, ok := s.clients[id]; ok {
		return c, nil
	}
	return nil, errors.New("Client not found")
}

func (s *TestStorage) SetClient(id string, client osin.Client) error {
	fmt.Printf("SetClient: %s\n", id)
	s.clients[id] = client
	return nil
}

func (s *TestStorage) SaveAuthorize(data *osin.AuthorizeData) error {
	fmt.Printf("SaveAuthorize: %s\n", data.Code)
	s.authorize[data.Code] = data
	b, err := json.Marshal(data); if err != nil {
		fmt.Println(err)
	}
	set(data.Code, string(b))
	return nil
}

func (s *TestStorage) LoadAuthorize(code string) (*osin.AuthorizeData, error) {
	fmt.Printf("LoadAuthorize: %s\n", code)
	
	val := get(code)
	if val != "" {
		var authData *osin.AuthorizeData
		err := json.Unmarshal([]byte(val), authData); if err != nil {
			fmt.Println(err)
		}
		return authData, nil
	}
	return nil, errors.New("Authorize not found")
}

func (s *TestStorage) RemoveAuthorize(code string) error {
	fmt.Printf("RemoveAuthorize: %s\n", code)
	delete(s.authorize, code)
	return nil
}

func (s *TestStorage) SaveAccess(data *osin.AccessData) error {
	fmt.Printf("SaveAccess: %s\n", data.AccessToken)
	s.access[data.AccessToken] = data
	if data.RefreshToken != "" {
		s.refresh[data.RefreshToken] = data.AccessToken
	}
	return nil
}

func (s *TestStorage) LoadAccess(code string) (*osin.AccessData, error) {
	fmt.Printf("LoadAccess: %s\n", code)
	if d, ok := s.access[code]; ok {
		return d, nil
	}
	return nil, errors.New("Access not found")
}

func (s *TestStorage) RemoveAccess(code string) error {
	fmt.Printf("RemoveAccess: %s\n", code)
	delete(s.access, code)
	return nil
}

func (s *TestStorage) LoadRefresh(code string) (*osin.AccessData, error) {
	fmt.Printf("LoadRefresh: %s\n", code)
	if d, ok := s.refresh[code]; ok {
		return s.LoadAccess(d)
	}
	return nil, errors.New("Refresh not found")
}

func (s *TestStorage) RemoveRefresh(code string) error {
	fmt.Printf("RemoveRefresh: %s\n", code)
	delete(s.refresh, code)
	return nil
}

// redis
func get(key string) string {
  val, err := client.Get(key).Result()
  if err != nil {
    fmt.Println(err)
  } else if err == redis.Nil {
    fmt.Println(key + " does not exists")
  }
  return val
}

func set(key string, val string) {
  b, err := json.Marshal(val)
  if err != nil {
    fmt.Println(val + " is not json format")
  } else {
    err := client.Set(key, b, 0).Err()
    if err != nil {
        panic(err)
    }
  }
}