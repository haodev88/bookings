package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type postData struct {
	key string
	value string
}

// http.StatusOK

var theTests = []struct {
	name string
	url string
	method string
	params []postData
	expectedStatusCode int
}{
	// For GET method
	{"home", "/", "GET", []postData{}, http.StatusOK},
	{"about", "/about", "GET", []postData{}, http.StatusOK},
	{"generals-quarters", "/generals-quarters", "GET", []postData{}, http.StatusOK},
	{"majors-suite", "/majors-suite", "GET", []postData{}, http.StatusOK},
	{"make-res", "/make-reservation", "GET", []postData{}, http.StatusOK},
	{"search-availability", "/search-availability", "GET", []postData{}, http.StatusOK},
	{"contact", "/contact", "GET", []postData{}, http.StatusOK},

	// For Post Method
	{"post-search-avail","/search-availability", "POST", []postData{
		{key:"start", value: "2024-03-20"},
		{key:"end", value: "2024-03-21"},
	}, http.StatusOK},

	{"post-search-avail-json","/search-availability-json", "POST", []postData{
		{key:"start", value: "2024-03-20"},
		{key:"end", value: "2024-03-21"},
	}, http.StatusOK},

	{"post make reservation","/make-reservation", "POST", []postData{
		{key:"first_name", value: "Hao"},
		{key:"last_name", value: "Tran"},
		{key: "email", value: "me@here.com"},
		{key: "phone", value: "555-555-5555"},
	}, http.StatusOK},

}

func TestHandlers(t *testing.T)  {
	routes:= getRoutes()
	ts:= httptest.NewTLSServer(routes)
	defer ts.Close()
	for _,e:= range theTests {
		if e.method == "GET" {
			resp, err:= ts.Client().Get(ts.URL + e.url)
			if err!=nil {
				t.Log(err)
				t.Fatal(err)
			}
			if resp.StatusCode != e.expectedStatusCode {
				t.Errorf("for %s expected %d but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
			}
		} else {
			values:=url.Values{}
			for _,x:= range e.params {
				values.Add(x.key, x.value)
			}
			resp,err:=ts.Client().PostForm(ts.URL + e.url, values)
			if err!=nil {
				t.Log(err)
				t.Fatal(err)
			}
			if resp.StatusCode != e.expectedStatusCode {
				t.Errorf("for %s expected %d but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
			}
		}
	}
}