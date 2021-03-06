package server

import "regexp"

// Route describes a route component implementation.
type Route struct {
	regex    *regexp.Regexp
	handlers map[string]HandleContextFunc
}

// DefaultRoute creates new route component.
//
// @param
// - patternURL {string} (the URL matching pattern)
//
// @return
// - route {Route} (a Route's new instance)
func DefaultRoute(patternURL string) *Route {
	route := &Route{
		regex:    regexp.MustCompile(patternURL),
		handlers: make(map[string]HandleContextFunc),
	}
	return route
}

// BindHandler binds HTTP request method with handler.
//
// @param
// - method {string} (HTTP request method)
// - handler {HandleContextFunc} (the callback func)
func (r *Route) BindHandler(method string, handler HandleContextFunc) {
	/* Condition validation: only accept function */
	if handler == nil {
		panic("Request handler must not be nil.")
	}

	/* Condition validation: only accept if there is none associated handler */
	if r.handlers[method] != nil {
		panic("This HTTP request method had been associated with another handler.")
	}
	r.handlers[method] = handler
}

// InvokeHandler invokes handler.
//
// @param
// - c {RequestContext} (the request context)
func (r *Route) InvokeHandler(c *RequestContext) {
	handler := r.handlers[c.Method]
	handler(c)
}

// Match matchs request path against route's regex pattern.
//
// @param
// - method {string} (HTTP request method)
// - pathURL {string} (request's path that will be matched)
//
// @return
// - flag {bool} (indicate flag if it is a matched or not)
// - pathParams {map[string]string} (a path params)
func (r *Route) Match(method string, pathURL string) (bool, map[string]string) {
	if matches := r.regex.FindStringSubmatch(pathURL); len(matches) > 0 && matches[0] == pathURL {
		if handler := r.handlers[method]; handler != nil {

			// Find path params if there is any
			var params map[string]string
			if names := r.regex.SubexpNames(); len(names) > 1 {

				params = make(map[string]string)
				for i, name := range names {
					if len(name) > 0 {
						params[name] = matches[i]
					}
				}
			}

			// Return result
			return true, params
		}
	}
	return false, nil
}
