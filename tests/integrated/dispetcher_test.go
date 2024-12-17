package integrated

import (
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/l1qwie/JWTAuth/app/logs"
	"github.com/l1qwie/JWTAuth/tests"
)

func callserver(method, url, key string, body io.Reader, t *testing.T, withheader bool) *http.Response {
	var resp *http.Response
	if req, err := http.NewRequest(method, fmt.Sprint("http://localhost:8080", url), body); err != nil {
		t.Fatal(err)
	} else {
		if withheader {
			req.Header.Add("Refresh-Token", key)
		}
		client := &http.Client{}
		if resp, err = client.Do(req); err != nil {
			t.Fatal(err)
		}
	}
	return resp
}

func invalidreq(resp *http.Response, t *testing.T) {
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatal("status code isn't 400")
	}
}

func okreq(resp *http.Response, t *testing.T) {
	if resp.StatusCode != http.StatusOK {
		t.Fatal("status code isn't 200")
	}
}

func TestLoginNoGUID(t *testing.T) {
	resp := callserver("GET", "/login/", "", nil, t, false)
	invalidreq(resp, t)
}

func TestLoginOK(t *testing.T) {
	resp := callserver("GET", "/login/?id=123e4567-e89b-12d3-a456-426614174000&justacall=true", "", nil, t, false)
	okreq(resp, t)
}

func TestRefreshNoToken(t *testing.T) {
	resp := callserver("PATCH", "/refresh", "", nil, t, true)
	invalidreq(resp, t)
}

func TestRefreshOK(t *testing.T) {
	resp := callserver("PATCH", "/refresh/?justacall=true", "a refresh token", nil, t, true)
	okreq(resp, t)
}

func init() {
	logs.SetDebug()
	tests.PutEnvVal()
	tests.StartAPI()
}
