package lrys

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	tool "autotool"
)

type HandlerFunc func(Handler)

type Handler struct {
	Response http.ResponseWriter
	Request  *http.Request
}

type Engine struct {
	router map[string]HandlerFunc
	port   int
}

type Rest struct {
	Client *http.Client
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	D := Handler{
		Response: w,
		Request:  r,
	}
	key := r.Method + "-" + r.URL.Path
	if handler, ok := engine.router[key]; ok {
		handler(D)
	} else {
		fmt.Fprintf(w, "404 Not Founda: %s\n", r.URL.Path)
	}
}

func (Hand *Handler) GetFrom(str string) (string, bool) {
	if Hand.Request.Form.Has(str) {
		return Hand.Request.Form.Get(str), true
	} else {
		return "", false
	}
}

func (Hand *Handler) GetBody() ([]byte, error) {
	body, err := io.ReadAll(io.Reader(Hand.Request.Body))
	return body, err
}

func (engine *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	key := method + "-" + pattern
	engine.router[key] = handler
}

func (engine *Engine) GET(pattern string, handler HandlerFunc) {
	engine.addRoute("GET", pattern, handler)
}

func (engine *Engine) RunServer(port int) {
	engine.port = port

	fmt.Printf("%s%d\n", "http server running on port ", engine.port)
	err := http.ListenAndServe(":"+strconv.Itoa(engine.port), engine)
	if err != nil {
		panic(tool.Error(tool.LogSprint("RunServerError:"+err.Error(), tool.ERROR, 1)))
	}

}

func (engine *Engine) POST(pattern string, handler HandlerFunc) {
	engine.addRoute("POST", pattern, handler)
}

func New() *Engine {
	return &Engine{router: make(map[string]HandlerFunc)}
}

func (engine *Engine) Static(static string, root string) {
	if strings.Contains(root, ":") || strings.Contains(root, "*") {
		panic("URL parameters can not be used when serving a static folder")
	}
	http.Handle(static, http.FileServer(http.Dir(root)))
}
