package expressgo

import (
	"log"
	"net/http"
)

type ExpressRoute struct {
	path string
}

func (app *Express) NewRoute(path string) *ExpressRoute {
	return &ExpressRoute{
		path: path,
	}
}

func (route *ExpressRoute) Route(path string) *ExpressRoute {
	return &ExpressRoute{
		path: route.path + path,
	}
}

func (route *ExpressRoute) Use(mw MiddlewareFunction) *ExpressRoute {
	_GlobalContext.middlewares = append(_GlobalContext.middlewares, _Middleware{
		execute: mw,
		path:    route.path,
	})
	return route
}

func (route *ExpressRoute) Get(handler func(w http.ResponseWriter, r *http.Request, context *RequestContext)) {
	route.handleRequest(http.MethodGet, handler)
}

func (route *ExpressRoute) Post(handler func(w http.ResponseWriter, r *http.Request, context *RequestContext)) {
	route.handleRequest(http.MethodPost, handler)
}

func (route *ExpressRoute) Put(handler func(w http.ResponseWriter, r *http.Request, context *RequestContext)) {
	route.handleRequest(http.MethodPut, handler)
}

func (route *ExpressRoute) Delete(handler func(w http.ResponseWriter, r *http.Request, context *RequestContext)) {
	route.handleRequest(http.MethodDelete, handler)
}

func (route *ExpressRoute) handleRequest(method string, handler func(w http.ResponseWriter, r *http.Request, context *RequestContext)) {
	log.Default().Printf("Registered %s request for route %s", method, route.path)
	_GlobalContext.mux.HandleFunc(route.path, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		requestContext := &RequestContext{
			context: make(map[string]any),
		}
		for _, mw := range _GlobalContext.middlewares {
			ok := mw.execute(r, requestContext)

			if !ok {
				_GlobalContext.errorHandler(w, r, requestContext)
				return
			}
		}
		handler(w, r, requestContext)
	})
}
