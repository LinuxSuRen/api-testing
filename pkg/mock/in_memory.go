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
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/swaggest/openapi-go/openapi3"
	"github.com/swaggest/rest/gorillamux"

	"github.com/linuxsuren/api-testing/pkg/version"

	"github.com/linuxsuren/api-testing/pkg/logging"
	"github.com/linuxsuren/api-testing/pkg/render"
	"github.com/linuxsuren/api-testing/pkg/util"

	"github.com/gorilla/mux"
)

var (
	memLogger = logging.DefaultLogger(logging.LogLevelInfo).WithName("memory")
)

type inMemoryServer struct {
	data       map[string][]map[string]interface{}
	mux        *mux.Router
	listener   net.Listener
	port       int
	wg         sync.WaitGroup
	ctx        context.Context
	cancelFunc context.CancelFunc
	reader     Reader
}

func NewInMemoryServer(port int) DynamicServer {
	ctx, cancel := context.WithCancel(context.TODO())
	return &inMemoryServer{
		port:       port,
		wg:         sync.WaitGroup{},
		ctx:        ctx,
		cancelFunc: cancel,
	}
}

func (s *inMemoryServer) SetupHandler(reader Reader, prefix string) (handler http.Handler, err error) {
	s.reader = reader
	// init the data
	s.data = make(map[string][]map[string]interface{})
	s.mux = mux.NewRouter().PathPrefix(prefix).Subrouter()
	handler = s.mux
	err = s.Load()
	return
}

func (s *inMemoryServer) Load() (err error) {
	var server *Server
	if server, err = s.reader.Parse(); err != nil {
		return
	}

	memLogger.Info("start to run all the APIs from objects", "count", len(server.Objects))
	for _, obj := range server.Objects {
		memLogger.Info("start mock server from object", "name", obj.Name)
		s.startObject(obj)
		s.initObjectData(obj)
	}

	memLogger.Info("start to run all the APIs from items", "count", len(server.Items))
	for _, item := range server.Items {
		s.startItem(item)
	}

	memLogger.Info("start webhook servers", "count", len(server.Webhooks))
	for _, item := range server.Webhooks {
		if err = s.startWebhook(&item); err != nil {
			return
		}
	}

	s.handleOpenAPI()
	return
}

func (s *inMemoryServer) Start(reader Reader, prefix string) (err error) {
	var handler http.Handler
	if handler, err = s.SetupHandler(reader, prefix); err != nil {
		return
	}

	if s.listener, err = net.Listen("tcp", fmt.Sprintf(":%d", s.port)); err != nil {
		return
	}
	go func() {
		err = http.Serve(s.listener, handler)
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

			data, err := json.Marshal(allItems)
			writeResponse(w, data, err)
		case http.MethodPost:
			// create an item
			if data, err := io.ReadAll(req.Body); err == nil {
				objData := map[string]interface{}{}

				jsonErr := json.Unmarshal(data, &objData)
				if jsonErr != nil {
					memLogger.Info(jsonErr.Error())
					return
				}

				s.data[obj.Name] = append(s.data[obj.Name], objData)

				_, _ = w.Write(data)
			} else {
				memLogger.Info("failed to read from body", "error", err)
			}
		case http.MethodDelete:
			// delete an item
			if data, err := io.ReadAll(req.Body); err == nil {
				objData := map[string]interface{}{}

				jsonErr := json.Unmarshal(data, &objData)
				if jsonErr != nil {
					memLogger.Info(jsonErr.Error())
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
				memLogger.Info("failed to read from body", "error", err)
			}
		case http.MethodPut:
			if data, err := io.ReadAll(req.Body); err == nil {
				objData := map[string]interface{}{}

				jsonErr := json.Unmarshal(data, &objData)
				if jsonErr != nil {
					memLogger.Info(jsonErr.Error())
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
				memLogger.Info("failed to read from body", "error", err)
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
					writeResponse(w, data, err)
					return
				}
			}
		}
	})
	return
}

func (s *inMemoryServer) startItem(item Item) {
	method := util.EmptyThenDefault(item.Request.Method, http.MethodGet)
	memLogger.Info("register mock service", "method", method, "path", item.Request.Path, "encoder", item.Response.Encoder)

	var headerSlices []string
	for k, v := range item.Request.Header {
		headerSlices = append(headerSlices, k, v)
	}

	adHandler := &advanceHandler{item: &item}
	s.mux.HandleFunc(item.Request.Path, adHandler.handle).Methods(strings.Split(method, ",")...).Headers(headerSlices...)
}

type advanceHandler struct {
	item *Item
}

func (h *advanceHandler) handle(w http.ResponseWriter, req *http.Request) {
	memLogger.Info("receiving mock request", "name", h.item.Name, "method", req.Method, "path", req.URL.Path,
		"encoder", h.item.Response.Encoder)

	var err error
	if h.item.Response.Encoder == "base64" {
		h.item.Response.BodyData, err = base64.StdEncoding.DecodeString(h.item.Response.Body)
	} else if h.item.Response.Encoder == "url" {
		var resp *http.Response
		if resp, err = http.Get(h.item.Response.Body); err == nil {
			h.item.Response.BodyData, err = io.ReadAll(resp.Body)
		}
	} else {
		h.item.Response.BodyData, err = render.RenderAsBytes("start-item", h.item.Response.Body, h.item)
	}

	if err == nil {
		h.item.Param = mux.Vars(req)
		if h.item.Param == nil {
			h.item.Param = make(map[string]string)
		}
		h.item.Param["Host"] = req.Host
		if h.item.Response.Header == nil {
			h.item.Response.Header = make(map[string]string)
		}
		h.item.Response.Header[headerMockServer] = fmt.Sprintf("api-testing: %s", version.GetVersion())
		h.item.Response.Header[util.ContentLength] = fmt.Sprintf("%d", len(h.item.Response.BodyData))
		for k, v := range h.item.Response.Header {
			hv, hErr := render.Render("mock-server-header", v, &h.item)
			if hErr != nil {
				hv = v
				memLogger.Error(hErr, "failed render mock-server-header", "value", v)
			}

			w.Header().Set(k, hv)
		}
		w.WriteHeader(util.ZeroThenDefault(h.item.Response.StatusCode, http.StatusOK))
	}

	writeResponse(w, h.item.Response.BodyData, err)
}

func writeResponse(w http.ResponseWriter, data []byte, err error) {
	if err == nil {
		w.Write(data)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}
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
			memLogger.Info(jsonErr.Error())
		}
	}
}

func (s *inMemoryServer) startWebhook(webhook *Webhook) (err error) {
	if webhook.Timer == "" || webhook.Name == "" {
		return
	}

	var duration time.Duration
	duration, err = time.ParseDuration(webhook.Timer)
	if err != nil {
		memLogger.Error(err, "Error parsing webhook timer")
		return
	}

	s.wg.Add(1)
	go func(wh *Webhook) {
		defer s.wg.Done()

		memLogger.Info("start webhook server", "name", wh.Name)
		timer := time.NewTimer(duration)
		for {
			timer.Reset(duration)
			select {
			case <-s.ctx.Done():
				memLogger.Info("stop webhook server", "name", wh.Name)
				return
			case <-timer.C:
				client := http.DefaultClient

				payload, err := render.RenderAsReader("mock webhook server payload", wh.Request.Body, wh)
				if err != nil {
					memLogger.Error(err, "Error when render payload")
					continue
				}

				method := util.EmptyThenDefault(wh.Request.Method, http.MethodPost)
				api, err := render.Render("webhook request api", wh.Request.Path, s)
				if err != nil {
					memLogger.Error(err, "Error when render api", "raw", wh.Request.Path)
					continue
				}

				req, err := http.NewRequestWithContext(s.ctx, method, api, payload)
				if err != nil {
					memLogger.Error(err, "Error when create request")
					continue
				}

				resp, err := client.Do(req)
				if err != nil {
					memLogger.Error(err, "Error when sending webhook")
				} else {
					memLogger.Info("received from webhook", "code", resp.StatusCode)
				}
			}
		}
	}(webhook)
	return
}

func (s *inMemoryServer) handleOpenAPI() {
	s.mux.HandleFunc("/api.json", func(w http.ResponseWriter, req *http.Request) {
		// Setup OpenAPI schema
		reflector := openapi3.NewReflector()
		reflector.SpecSchema().SetTitle("Mock Server API")
		reflector.SpecSchema().SetVersion(version.GetVersion())
		reflector.SpecSchema().SetDescription("Powered by https://github.com/linuxsuren/api-testing")

		// Walk the router with OpenAPI collector
		c := gorillamux.NewOpenAPICollector(reflector)

		_ = s.mux.Walk(c.Walker)

		// Get the resulting schema
		if jsonData, err := reflector.Spec.MarshalJSON(); err == nil {
			w.Header().Set(util.ContentType, util.JSON)
			_, _ = w.Write(jsonData)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
		}
	})
}

func jsonStrToInterface(jsonStr string) (objData map[string]interface{}, err error) {
	if jsonStr, err = render.Render("init object", jsonStr, nil); err == nil {
		objData = map[string]interface{}{}
		err = json.Unmarshal([]byte(jsonStr), &objData)
	}
	return
}

func (s *inMemoryServer) GetPort() string {
	return util.GetPort(s.listener)
}

func (s *inMemoryServer) Stop() (err error) {
	if s.listener != nil {
		err = s.listener.Close()
	}
	if s.cancelFunc != nil {
		s.cancelFunc()
	}
	s.wg.Wait()
	return
}
