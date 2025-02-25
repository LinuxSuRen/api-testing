/*
Copyright 2024-2025 API Testing Authors.

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
	"bytes"
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
	prefix     string
	wg         sync.WaitGroup
	ctx        context.Context
	cancelFunc context.CancelFunc
	reader     Reader
	metrics    RequestMetrics
}

func NewInMemoryServer(ctx context.Context, port int) DynamicServer {
	ctx, cancel := context.WithCancel(ctx)
	return &inMemoryServer{
		port:       port,
		wg:         sync.WaitGroup{},
		ctx:        ctx,
		cancelFunc: cancel,
		metrics:    NewNoopMetrics(),
	}
}

func (s *inMemoryServer) SetupHandler(reader Reader, prefix string) (handler http.Handler, err error) {
	s.reader = reader
	// init the data
	s.data = make(map[string][]map[string]interface{})
	s.mux = mux.NewRouter().PathPrefix(prefix).Subrouter()
	s.prefix = prefix
	handler = s.mux
	s.metrics.AddMetricsHandler(s.mux)
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
			continue
		}
	}

	s.handleOpenAPI()

	for _, proxy := range server.Proxies {
		memLogger.Info("start to proxy", "target", proxy.Target)
		s.mux.HandleFunc(proxy.Path, func(w http.ResponseWriter, req *http.Request) {
			apiRaw := fmt.Sprintf("%s/%s", proxy.Target, strings.TrimPrefix(req.URL.Path, s.prefix))
			var api string
			api, err = render.Render("proxy api", apiRaw, s)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				memLogger.Error(err, "failed to render proxy api", "api", apiRaw)
				return
			}
			memLogger.Info("redirect to", "target", api)

			targetReq, err := http.NewRequestWithContext(req.Context(), req.Method, api, req.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				memLogger.Error(err, "failed to create proxy request")
				return
			}

			resp, err := http.DefaultClient.Do(targetReq)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				memLogger.Error(err, "failed to do proxy request")
				return
			}

			data, err := io.ReadAll(resp.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				memLogger.Error(err, "failed to read response body")
				return
			}

			for k, v := range resp.Header {
				w.Header().Add(k, v[0])
			}
			w.Write(data)
		})
	}
	return
}

func (s *inMemoryServer) Start(reader Reader, prefix string) (err error) {
	var handler http.Handler
	if handler, err = s.SetupHandler(reader, prefix); err == nil {
		if s.listener, err = net.Listen("tcp", fmt.Sprintf(":%d", s.port)); err == nil {
			go func() {
				err = http.Serve(s.listener, handler)
			}()
		}
	}
	return
}

func (s *inMemoryServer) EnableMetrics() {
	s.metrics = NewInMemoryMetrics()
}

func (s *inMemoryServer) startObject(obj Object) {
	// create a simple CRUD server
	s.mux.HandleFunc("/"+obj.Name, func(w http.ResponseWriter, req *http.Request) {
		fmt.Println("mock server received request", req.URL.Path)
		s.metrics.RecordRequest(req.URL.Path)
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
					if len(v) == 0 {
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
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	// handle a single object
	s.mux.HandleFunc(fmt.Sprintf("/%s/{name}", obj.Name), func(w http.ResponseWriter, req *http.Request) {
		s.metrics.RecordRequest(req.URL.Path)
		w.Header().Set(util.ContentType, util.JSON)
		objects := s.data[obj.Name]
		if objects != nil {
			name := mux.Vars(req)["name"]
			var data []byte
			for _, obj := range objects {
				if obj["name"] == name {

					data, _ = json.Marshal(obj)
					break
				}
			}

			if len(data) == 0 {
				w.WriteHeader(http.StatusNotFound)
				return
			}

			method := req.Method
			switch method {
			case http.MethodGet:
				writeResponse(w, data, nil)
			case http.MethodPut:
				objData := map[string]interface{}{}
				if data, err := io.ReadAll(req.Body); err == nil {

					jsonErr := json.Unmarshal(data, &objData)
					if jsonErr != nil {
						memLogger.Info(jsonErr.Error())
						return
					}
					for i, item := range s.data[obj.Name] {
						if item["name"] == name {
							s.data[obj.Name][i] = objData
							break
						}
					}
					_, _ = w.Write(data)
				}
			case http.MethodDelete:
				for i, item := range s.data[obj.Name] {
					if item["name"] == name {
						if len(s.data[obj.Name]) == i+1 {
							s.data[obj.Name] = s.data[obj.Name][:i]
						} else {
							s.data[obj.Name] = append(s.data[obj.Name][:i], s.data[obj.Name][i+1])
						}

						writeResponse(w, []byte(`{"msg": "deleted"}`), nil)
					}
				}
			default:
				w.WriteHeader(http.StatusMethodNotAllowed)
			}

		}
	})
}

func (s *inMemoryServer) startItem(item Item) {
	method := util.EmptyThenDefault(item.Request.Method, http.MethodGet)
	memLogger.Info("register mock service", "method", method, "path", item.Request.Path, "encoder", item.Response.Encoder)

	var headerSlices []string
	for k, v := range item.Request.Header {
		headerSlices = append(headerSlices, k, v)
	}

	adHandler := &advanceHandler{
		item:    &item,
		metrics: s.metrics,
		mu:      sync.Mutex{},
	}
	s.mux.HandleFunc(item.Request.Path, adHandler.handle).Methods(strings.Split(method, ",")...).Headers(headerSlices...)
}

type advanceHandler struct {
	item    *Item
	metrics RequestMetrics
	mu      sync.Mutex
}

func (h *advanceHandler) handle(w http.ResponseWriter, req *http.Request) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.metrics.RecordRequest(req.URL.Path)
	memLogger.Info("receiving mock request", "name", h.item.Name, "method", req.Method, "path", req.URL.Path,
		"encoder", h.item.Response.Encoder)

	h.item.Param = mux.Vars(req)
	if h.item.Param == nil {
		h.item.Param = make(map[string]string)
	}
	h.item.Param["Host"] = req.Host
	if h.item.Response.Header == nil {
		h.item.Response.Header = make(map[string]string)
	}
	h.item.Response.Header[headerMockServer] = fmt.Sprintf("api-testing: %s", version.GetVersion())
	for k, v := range h.item.Response.Header {
		hv, hErr := render.Render("mock-server-header", v, &h.item)
		if hErr != nil {
			hv = v
			memLogger.Error(hErr, "failed render mock-server-header", "value", v)
		}

		w.Header().Set(k, hv)
	}

	var err error
	if h.item.Response.Encoder == "base64" {
		h.item.Response.BodyData, err = base64.StdEncoding.DecodeString(h.item.Response.Body)
	} else if h.item.Response.Encoder == "url" {
		var resp *http.Response
		if resp, err = http.Get(h.item.Response.Body); err == nil {
			h.item.Response.BodyData, err = io.ReadAll(resp.Body)
		}
	} else {
		if h.item.Response.BodyData, err = render.RenderAsBytes("start-item", h.item.Response.Body, h.item); err != nil {
			fmt.Printf("failed to render body: %v", err)
		}
	}

	if err == nil {
		h.item.Response.Header[util.ContentLength] = fmt.Sprintf("%d", len(h.item.Response.BodyData))
		w.Header().Set(util.ContentLength, h.item.Response.Header[util.ContentLength])
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
				if err = runWebhook(s.ctx, s, wh); err != nil {
					memLogger.Error(err, "Error when run webhook")
				}
			}
		}
	}(webhook)
	return
}

func runWebhook(ctx context.Context, objCtx interface{}, wh *Webhook) (err error) {
	client := http.DefaultClient

	var payload io.Reader
	payload, err = render.RenderAsReader("mock webhook server payload", wh.Request.Body, wh)
	if err != nil {
		err = fmt.Errorf("error when render payload: %w", err)
		return
	}

	method := util.EmptyThenDefault(wh.Request.Method, http.MethodPost)
	var api string
	api, err = render.Render("webhook request api", wh.Request.Path, objCtx)
	if err != nil {
		err = fmt.Errorf("error when render api: %w, template: %s", err, wh.Request.Path)
		return
	}

	var bearerToken string
	bearerToken, err = getBearerToken(ctx, wh.Request)
	if err != nil {
		memLogger.Error(err, "Error when render bearer token")
		return
	}

	var req *http.Request
	req, err = http.NewRequestWithContext(ctx, method, api, payload)
	if err != nil {
		memLogger.Error(err, "Error when create request")
		return
	}

	if bearerToken != "" {
		memLogger.V(7).Info("set bearer token", "token", bearerToken)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", bearerToken))
	}

	for k, v := range wh.Request.Header {
		req.Header.Set(k, v)
	}

	memLogger.Info("send webhook request", "api", api)
	resp, err := client.Do(req)
	if err != nil {
		err = fmt.Errorf("error when sending webhook")
	} else {
		if resp.StatusCode != http.StatusOK {
			memLogger.Info("unexpected status", "code", resp.StatusCode)
		}

		data, _ := io.ReadAll(resp.Body)
		memLogger.V(7).Info("received from webhook", "code", resp.StatusCode, "response", string(data))
	}
	return
}

type bearerToken struct {
	Token string `json:"token"`
}

func getBearerToken(ctx context.Context, request RequestWithAuth) (token string, err error) {
	if request.BearerAPI == "" {
		return
	}

	if request.BearerAPI, err = render.Render("bearer token request", request.BearerAPI, &request); err != nil {
		return
	}

	var data []byte
	if data, err = json.Marshal(&request); err == nil {
		client := http.DefaultClient
		var req *http.Request
		if req, err = http.NewRequestWithContext(ctx, http.MethodPost, request.BearerAPI, bytes.NewBuffer(data)); err == nil {
			req.Header.Set(util.ContentType, util.JSON)

			var resp *http.Response
			if resp, err = client.Do(req); err == nil && resp.StatusCode == http.StatusOK {
				if data, err = io.ReadAll(resp.Body); err == nil {
					var tokenObj bearerToken
					if err = json.Unmarshal(data, &tokenObj); err == nil {
						token = tokenObj.Token
					}
				}
			}
		}
	}

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
	} else {
		memLogger.Info("listener is nil")
	}
	if s.cancelFunc != nil {
		s.cancelFunc()
	}
	s.wg.Wait()
	return
}
