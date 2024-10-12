package expressgo

import (
	"fmt"
	"log"
	"net/http"
)

type Express struct {
}

type ExpressConfig struct {
	Host string
	Port uint16
}

type MiddlewareFunction func(r *http.Request, context *RequestContext) bool

type ErrorMiddlewareFunction func(w http.ResponseWriter, r *http.Request, context *RequestContext)

type _Middleware struct {
	execute MiddlewareFunction
	path    string
}

type RequestContext struct {
	context map[string]any
}

type _ExpressGlobalContext struct {
	mux          *http.ServeMux
	config       *ExpressConfig
	middlewares  []_Middleware
	errorHandler ErrorMiddlewareFunction
}

var _GlobalContext *_ExpressGlobalContext

func NewExpress(config *ExpressConfig) *Express {
	mux := http.NewServeMux()
	_GlobalContext = &_ExpressGlobalContext{
		mux:          mux,
		config:       config,
		middlewares:  make([]_Middleware, 0),
		errorHandler: nil,
	}
	return &Express{}
}

func (app *Express) StartServer() {
	log.Default().Printf("Starting server on %v", _GlobalContext.config.Port)
	http.ListenAndServe(fmt.Sprintf(":%v", _GlobalContext.config.Port), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if e := recover(); e != nil {
				em := "Unhandled Exception"

				switch e := e.(type) {
				case error:
					em = e.Error()
				case string:
					em = e
				}
				w.Write([]byte(em))
				return
			}
		}()
		_GlobalContext.mux.ServeHTTP(w, r)
	}))
}

func (app *Express) Use(mw MiddlewareFunction) {
	_GlobalContext.middlewares = append(_GlobalContext.middlewares, _Middleware{
		execute: mw,
		path:    "*",
	})
}

func (app *Express) UseErrorHandler(handler ErrorMiddlewareFunction) {
	_GlobalContext.errorHandler = handler
}

func (requestContext *RequestContext) Get(key string) any {
	return requestContext.context[key]
}

func (requestContext *RequestContext) Set(key string, value any) {
	requestContext.context[key] = value
}

func (requestContext *RequestContext) GetRaw() *map[string]any {
	return &requestContext.context
}
