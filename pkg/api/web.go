package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jfbramlett/go-aop/pkg/aop"
	"github.com/jfbramlett/go-aop/pkg/web"
	"html/template"
	"net/http"
	"path"
)

var (
	homeTpl *template.Template
)

type Resource struct {
	Name string
	URL  string
}

type WebResponse struct {
	Msg       string
	Error     error
	Type      string
	Resources []Resource
}

type Runner interface {
	Run()
}

func NewWebRunner(webRoot string) Runner {
	return &webRunner{webRoot: webRoot}
}

type webRunner struct {
	webRoot string
}

func (wr *webRunner) Run() {
	base := template.Must(template.ParseFiles(path.Join(wr.webRoot, "web", "index.html")))
	homeTpl = template.Must(template.Must(base.Clone()).ParseFiles(path.Join(wr.webRoot, "web", "templates", "home.html")))

	// set up HTTP sequence
	router := mux.NewRouter()
	router.HandleFunc("/", wr.indexHandler)
	router.HandleFunc("/endpoint", wr.getMessage)

	loggingMiddleware := &web.LoggingMiddleware{}
	spanMiddleware := &web.SpanMiddleware{}
	router.Use(loggingMiddleware.Middleware)
	router.Use(spanMiddleware.Middleware)

	fmt.Println(http.ListenAndServe(":8085", router))
}

func (wr *webRunner) indexHandler(w http.ResponseWriter, r *http.Request) {
	webResp := WebResponse{Msg: "Successful call", Type: "Sample",
		Resources: []Resource{{Name: "Resource1", URL: "resource1.html"}, {Name: "Resource2", URL: "resource2.html"}}}
	_ = homeTpl.Execute(w, webResp)
}

func (wr *webRunner) getMessage(w http.ResponseWriter, r *http.Request) {
	var err error

	ctx := aop.Before(context.Background())
	defer func() { aop.After(ctx, err) }()

	if r.Method == http.MethodGet {
		writeResponse(map[string]string{"msg": "hello world"}, http.StatusOK, w)
		return
	}

	writeResponse(map[string]string{"error": "GET is only supported method"}, http.StatusBadRequest, w)
}

func writeResponse(data interface{}, status int, writer http.ResponseWriter) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)

	enc := json.NewEncoder(writer)
	enc.SetIndent("", "    ")
	if err := enc.Encode(&data); err != nil {
		return
	}
}
