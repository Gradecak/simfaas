package simfaas

import (
	"net/http"
	"regexp"
	"strings"
)

// RegexpHandler is simple http.Handler to enable the use of
// wildcards in routes.
type RegexpHandler struct {
	routes []*route
}

func (h *RegexpHandler) Handler(pattern *regexp.Regexp, handler http.Handler) {
	h.routes = append(h.routes, &route{pattern, handler})
}

func (h *RegexpHandler) HandleFunc(pattern *regexp.Regexp, handler func(http.ResponseWriter, *http.Request)) {
	h.routes = append(h.routes, &route{pattern, http.HandlerFunc(handler)})
}

func (h *RegexpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// strip ? character to ignore query paramaters when matching routes
	path := strings.Split(r.URL.Path, "?")[0]
	// Reverse match, such that newer routes have precedence
	for i := len(h.routes) - 1; i >= 0; i-- {
		route := h.routes[i]
		if route.pattern.MatchString(path) {
			route.handler.ServeHTTP(w, r)
			return
		}
	}
	// no pattern matched; send 404 response
	http.NotFound(w, r)
}

type route struct {
	pattern *regexp.Regexp
	handler http.Handler
}
