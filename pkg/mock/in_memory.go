/*
Copyright 2024 API Testing Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package mock

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/linuxsuren/api-testing/pkg/render"
	"github.com/linuxsuren/api-testing/pkg/util"
	"github.com/linuxsuren/api-testing/pkg/version"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
)

type inMemoryServer struct {
	data     map[string][]map[string]interface{}
	mux      *mux.Router
	port     int
	listener net.Listener
}

func NewInMemoryServer(port int) DynamicServer {
	return &inMemoryServer{
		port: port,
	}
}

func (s *inMemoryServer) Start(reader Reader) (err error) {
	var server *Server
	if server, err = reader.Parse(); err != nil {
		return
	}

	// init the data
	s.data = make(map[string][]map[string]interface{})
	s.mux = mux.NewRouter()

	log.Println("start to run all the APIs from objects")
	for _, obj := range server.Objects {
		s.startObject(obj)
		s.initObjectData(obj)
	}

	log.Println("start to run all the APIs from items")
	for _, item := range server.Items {
		s.startItem(item)
	}

	s.listener, err = net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	go func() {
		err = http.Serve(s.listener, s.mux)
	}()
	return
}

func (s *inMemoryServer) startObject(obj Object) {
	// create a simple CRUD server
	s.mux.HandleFunc("/"+obj.Name, func(w http.ResponseWriter, req *http.Request) {
		method := req.Method
		w.Header().Set(util.ContentType, util.JSON)

		switch method {
		case http.MethodGet:
			// list all items
			allItems := s.data[obj.Name]
			filteredItems := make([]map[string]interface{}, 0)

			for i, item := range allItems {
				exclude := false

				for k, v := range req.URL.Query() {
					if v == nil || len(v) == 0 {
						continue
					}

					if val, ok := item[k]; ok && val != v[0] {
						exclude = true
						break
					}
				}

				if !exclude {
					filteredItems = append(filteredItems, allItems[i])
				}
			}

			if len(filteredItems) != len(allItems) {
				allItems = filteredItems
			}

			if data, err := json.Marshal(allItems); err == nil {
				w.Write(data)
			}
		case http.MethodPost:
			// create an item
			if data, err := io.ReadAll(req.Body); err == nil {
				objData := map[string]interface{}{}

				jsonErr := json.Unmarshal(data, &objData)
				if jsonErr != nil {
					log.Println(jsonErr)
					return
				}

				s.data[obj.Name] = append(s.data[obj.Name], objData)

				_, _ = w.Write(data)
			} else {
				log.Println("failed to read from body", err)
			}
		case http.MethodDelete:
			// delete an item
			if data, err := io.ReadAll(req.Body); err == nil {
				objData := map[string]interface{}{}

				jsonErr := json.Unmarshal(data, &objData)
				if jsonErr != nil {
					log.Println(jsonErr)
					return
				}

				for i, item := range s.data[obj.Name] {
					if objData["name"] == item["name"] {
						if len(s.data[obj.Name]) == i+1 {
							s.data[obj.Name] = s.data[obj.Name][:i]
						} else {
							s.data[obj.Name] = append(s.data[obj.Name][:i], s.data[obj.Name][i+1])
						}
						break
					}
				}

				_, _ = w.Write(data)
			} else {
				log.Println("failed to read from body", err)
			}
		case http.MethodPut:
			if data, err := io.ReadAll(req.Body); err == nil {
				objData := map[string]interface{}{}

				jsonErr := json.Unmarshal(data, &objData)
				if jsonErr != nil {
					log.Println(jsonErr)
					return
				}

				for i, item := range s.data[obj.Name] {
					if objData["name"] == item["name"] {
						s.data[obj.Name][i] = objData
						break
					}
				}

				_, _ = w.Write(data)
			} else {
				log.Println("failed to read from body", err)
			}
		default:
			w.WriteHeader(http.StatusBadRequest)
		}
	})

	// get a single object
	s.mux.HandleFunc(fmt.Sprintf("/%s/{name:[a-z]+}", obj.Name), func(w http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set(util.ContentType, util.JSON)
		objects := s.data[obj.Name]
		if objects != nil {
			name := mux.Vars(req)["name"]

			for _, obj := range objects {
				if obj["name"] == name {

					data, err := json.Marshal(obj)
					if err == nil {
						w.Write(data)
					} else {
						w.Write([]byte(err.Error()))
						w.WriteHeader(http.StatusBadRequest)
					}
					return
				}
			}
		}
	})
	return
}

func (s *inMemoryServer) startItem(item Item) {
	s.mux.HandleFunc(item.Request.Path, func(w http.ResponseWriter, req *http.Request) {
		item.Response.Headers = append(item.Response.Headers, Header{
			Key:   headerMockServer,
			Value: fmt.Sprintf("api-testing: %s", version.GetVersion()),
		})
		for _, header := range item.Response.Headers {
			w.Header().Set(header.Key, header.Value)
		}
		body, err := render.Render("start-item", item.Response.Body, req)
		if err == nil {
			w.Write([]byte(body))
		} else {
			w.Write([]byte(err.Error()))
		}
		w.WriteHeader(util.ZeroThenDefault(item.Response.StatusCode, http.StatusOK))
	}).Methods(util.EmptyThenDefault(item.Request.Method, http.MethodGet))
}

func (s *inMemoryServer) initObjectData(obj Object) {
	if obj.Sample == "" {
		return
	}

	defaultCount := 1
	if obj.InitCount == nil {
		obj.InitCount = &defaultCount
	}

	for i := 0; i < *obj.InitCount; i++ {
		objData, jsonErr := jsonStrToInterface(obj.Sample)
		if jsonErr == nil {
			s.data[obj.Name] = append(s.data[obj.Name], objData)
		} else {
			log.Println(jsonErr)
		}
	}
}

func jsonStrToInterface(jsonStr string) (objData map[string]interface{}, err error) {
	if jsonStr, err = render.Render("init object", jsonStr, nil); err == nil {
		objData = map[string]interface{}{}
		err = json.Unmarshal([]byte(jsonStr), &objData)
	}
	return
}

func (s *inMemoryServer) GetPort() string {
	addr := s.listener.Addr().String()
	items := strings.Split(addr, ":")
	return items[len(items)-1]
}

func (s *inMemoryServer) Stop() (err error) {
	if s.listener != nil {
		err = s.listener.Close()
	}
	return
}
