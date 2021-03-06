package server

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/phuc0302/go-server/expected_format"
	"github.com/phuc0302/go-server/util"
)

func Test_ServeHTTP_InvalidResource(t *testing.T) {
	defer os.Remove(Debug)
	Initialize(true)

	// Setup test server
	BindGet("/sample", func(c *RequestContext) {
		c.OutputJSON(util.Status200(), map[string]string{"apple": "apple"})
	})

	ts := httptest.NewServer(ServeHTTP())
	defer ts.Close()

	response, _ := http.Get(fmt.Sprintf("%s/%s", ts.URL, "resources/README"))
	if response.StatusCode != 404 {
		t.Errorf(expectedFormat.NumberButFoundNumber, 404, response.StatusCode)
	}
}

func Test_ServeHTTP_ValidResource(t *testing.T) {
	defer os.Remove(Debug)
	Initialize(true)

	// Setup test server
	BindGet("/sample", func(c *RequestContext) {
		c.OutputJSON(util.Status200(), map[string]string{"apple": "apple"})
	})

	ts := httptest.NewServer(ServeHTTP())
	defer ts.Close()

	response, _ := http.Get(fmt.Sprintf("%s/%s", ts.URL, "resources/LICENSE"))
	if response.StatusCode != 404 {
		t.Errorf(expectedFormat.NumberButFoundNumber, 404, response.StatusCode)
	}
}

func Test_ServeHTTP_InvalidHTTPMethod(t *testing.T) {
	defer os.Remove(Debug)
	Initialize(true)

	// Update allow methods
	Cfg.AllowMethods = []string{Get, Post, Patch, Delete}

	// Setup test server
	BindGet("/sample", func(c *RequestContext) {
		c.OutputJSON(util.Status200(), map[string]string{"apple": "apple"})
	})

	ts := httptest.NewServer(ServeHTTP())
	defer ts.Close()

	request, _ := http.NewRequest("LINK", fmt.Sprintf("%s/%s", ts.URL, "token"), nil)
	response, _ := http.DefaultClient.Do(request)
	if response.StatusCode != 405 {
		t.Errorf(expectedFormat.NumberButFoundNumber, 405, response.StatusCode)
	}
}

func Test_ServeHTTP_InvalidURL(t *testing.T) {
	defer os.Remove(Debug)
	Initialize(true)

	// Setup test server
	BindGet("/sample", func(c *RequestContext) {
		c.OutputJSON(util.Status200(), map[string]string{"apple": "apple"})
	})

	ts := httptest.NewServer(ServeHTTP())
	defer ts.Close()

	request, _ := http.NewRequest("POST", fmt.Sprintf("%s/%s", ts.URL, "sample"), strings.NewReader(""))
	request.Header.Set("content-type", "application/x-www-form-urlencoded")

	response, _ := http.DefaultClient.Do(request)
	if response.StatusCode != 503 {
		t.Errorf(expectedFormat.NumberButFoundNumber, 503, response.StatusCode)
	}
}

func Test_ServeHTTP_ValidURL(t *testing.T) {
	defer os.Remove(Debug)
	Initialize(true)

	// Setup test server
	BindGet("/sample", func(c *RequestContext) {
		c.OutputJSON(util.Status200(), map[string]string{"apple": "apple"})
	})

	ts := httptest.NewServer(ServeHTTP())
	defer ts.Close()

	response, _ := http.Get(fmt.Sprintf("%s/%s", ts.URL, "sample"))
	if response.StatusCode != 200 {
		t.Errorf(expectedFormat.NumberButFoundNumber, 200, response.StatusCode)
	} else {
		bytes, _ := ioutil.ReadAll(response.Body)

		if string(bytes) != "{\"apple\":\"apple\"}" {
			t.Errorf(expectedFormat.StringButFoundString, "{\"apple\":\"apple\"}", string(bytes))
		}
	}
}
