package server

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/phuc0302/go-server/expected_format"
)

func Test_GroupRoute(t *testing.T) {
	router := new(Router)
	router.GroupRoute("/user/profile", func() {
		router.BindRoute(Get, "", func(request *RequestContext) {})
		router.BindRoute(Get, "/{profileID}", func(request *RequestContext) {})
		router.BindRoute(Post, "/{profileID}", func(request *RequestContext) {})
	})

	if router.routes == nil {
		t.Error(expectedFormat.NotNil)
	} else {
		if len(router.routes) != 2 {
			t.Errorf(expectedFormat.NumberButFoundNumber, 2, len(router.routes))
		} else {
			route0 := router.routes[0]
			if route0.regex.String() != "^/user/profile/?$" {
				t.Errorf(expectedFormat.StringButFoundString, "^/user/profile/?$", route0.regex.String())
			}
			if route0.handlers[Get] == nil {
				t.Error(expectedFormat.NotNil)
			}

			route1 := router.routes[1]
			if route1.regex.String() != "^/user/profile/(?P<profileID>[^/#?]+)/?$" {
				t.Errorf(expectedFormat.StringButFoundString, "^/user/profile/(?P<profileID>[^/#?]+)/?$", route1.regex.String())
			}
			if route1.handlers[Get] == nil {
				t.Error(expectedFormat.NotNil)
			}
			if route1.handlers[Post] == nil {
				t.Error(expectedFormat.NotNil)
			}
		}
	}
}

func Test_BindRoute(t *testing.T) {
	router := new(Router)

	// [Test 1] First bind
	router.BindRoute(Get, "/", func(c *RequestContext) {})
	if len(router.routes) != 1 {
		t.Errorf(expectedFormat.NumberButFoundNumber, 1, len(router.routes))
	}

	// [Test 2] Second bind
	router.BindRoute(Get, "/sample", func(c *RequestContext) {})
	if len(router.routes) != 2 {
		t.Errorf(expectedFormat.NumberButFoundNumber, 2, len(router.routes))
	}
}

func Test_MatchRoute_InvalidPath(t *testing.T) {
	// Setup router
	router := new(Router)
	router.BindRoute(Get, "/", func(request *RequestContext) {})
	router.GroupRoute("/user/profile(.htm[l]?)?", func() {
		router.BindRoute(Get, "", func(request *RequestContext) {})
		router.BindRoute(Post, "", func(request *RequestContext) {})
		router.BindRoute(Get, "/{profileID}", func(request *RequestContext) {})
	})
	router.GroupRoute("/private", func() {
		router.BindRoute(Get, "", func(request *RequestContext) {})
		router.BindRoute(Get, "/{profileID}", func(request *RequestContext) {})
	})

	// Setup test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := httprouter.CleanPath(r.URL.Path)
		method := strings.ToLower(r.Method)

		route, pathParams := router.MatchRoute(method, path)
		if route != nil {
			t.Error(expectedFormat.Nil)
		}
		if pathParams != nil {
			t.Error(expectedFormat.Nil)
		}
	}))
	defer ts.Close()
	http.Get(fmt.Sprintf("%s/user", ts.URL))
}

func Test_MatchRoute_InvalidHTTPMethod(t *testing.T) {
	// Setup router
	router := new(Router)
	router.BindRoute(Get, "/", func(request *RequestContext) {})
	router.GroupRoute("/user/profile(.htm[l]?)?", func() {
		router.BindRoute(Get, "", func(request *RequestContext) {})
		router.BindRoute(Get, "/{profileID}", func(request *RequestContext) {})
	})
	router.GroupRoute("/private", func() {
		router.BindRoute(Get, "", func(request *RequestContext) {})
		router.BindRoute(Get, "/{profileID}", func(request *RequestContext) {})
	})

	// Setup test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := httprouter.CleanPath(r.URL.Path)
		method := strings.ToLower(r.Method)

		route, pathParams := router.MatchRoute(method, path)
		if route != nil {
			t.Error(expectedFormat.Nil)
		}
		if pathParams != nil {
			t.Error(expectedFormat.Nil)
		}
	}))
	defer ts.Close()
	http.Post(fmt.Sprintf("%s/user/profile", ts.URL), "application/x-www-form-urlencoded", nil)
}

func Test_MatchRoute_ValidHTTPMethodAndPath(t *testing.T) {
	// Setup router
	router := new(Router)
	router.BindRoute(Get, "/", func(request *RequestContext) {})
	router.GroupRoute("/user/profile(.htm[l]?)?", func() {
		router.BindRoute(Get, "", func(request *RequestContext) {})
		router.BindRoute(Post, "", func(request *RequestContext) {})
		router.BindRoute(Get, "/{profileID}", func(request *RequestContext) {})
	})
	router.GroupRoute("/private", func() {
		router.BindRoute(Get, "", func(request *RequestContext) {})
		router.BindRoute(Get, "/{profileID}", func(request *RequestContext) {})
	})

	// Setup test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := httprouter.CleanPath(r.URL.Path)
		method := strings.ToLower(r.Method)

		route, _ := router.MatchRoute(method, path)
		if route == nil {
			t.Error(path)
		}
	}))
	defer ts.Close()

	http.Get(ts.URL)
	http.Get(fmt.Sprintf("%s/user/profile", ts.URL))
	http.Get(fmt.Sprintf("%s/user/profile?userID=1", ts.URL))
	http.Get(fmt.Sprintf("%s/user/profile/", ts.URL))
	http.Get(fmt.Sprintf("%s/user/profile/?userID=1", ts.URL))
	http.Get(fmt.Sprintf("%s/user/profile/1", ts.URL))
	http.Get(fmt.Sprintf("%s/user/profile/1?userID=1", ts.URL))
	http.Get(fmt.Sprintf("%s/user/profile/1/", ts.URL))
	http.Get(fmt.Sprintf("%s/user/profile/1/?userID=1", ts.URL))

	http.Get(fmt.Sprintf("%s/user/profile.htm", ts.URL))
	http.Get(fmt.Sprintf("%s/user/profile.htm?userID=1", ts.URL))
	http.Get(fmt.Sprintf("%s/user/profile.htm/", ts.URL))
	http.Get(fmt.Sprintf("%s/user/profile.html/?userID=1", ts.URL))
}
