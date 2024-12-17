package e2e

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/l1qwie/JWTAuth/app/database"
	"github.com/l1qwie/JWTAuth/app/logs"
	"github.com/l1qwie/JWTAuth/app/types"
	"github.com/l1qwie/JWTAuth/tests"
)

var refreshtoken *string

const testguid = "123e4567-e89b-12d3-a456-426614174000"
const testip = "127.0.0.1"

type test struct {
	method     string
	url        string
	key        *string
	withheader bool
}

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

func testAction(arrt []test, t *testing.T) {
	for _, test := range arrt {
		resp := callserver(test.method, test.url, *test.key, nil, t, test.withheader)
		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			t.Log(string(body))
			t.Fatal("status response is bad")
		}
		if body, err := io.ReadAll(resp.Body); err != nil {
			t.Fatal(err)
		} else {
			tokens := new(types.Tokens)
			if err = json.Unmarshal(body, tokens); err != nil {
				t.Log(string(body))
				t.Fatal(err)
			}
			t.Logf("Access-Token: %s", tokens.Access)
			t.Logf("Refresh-Token: %s", tokens.Refresh)
			*refreshtoken = tokens.Refresh
		}
	}
}

func createEnv(t *testing.T) {
	var err error
	if database.Conn, err = database.Connect(); err != nil {
		t.Fatal(err)
	}
	if err = database.Conn.CreateMokData(testguid, testip); err != nil {
		t.Fatal(err)
	}
}

func TestE2E(t *testing.T) {
	defer database.Conn.DeleteUsers()
	var s string
	createEnv(t)
	refreshtoken = &s
	arrt := []test{{"GET", fmt.Sprintf("/login?id=%s", testguid), refreshtoken, false},
		{"PATCH", "/refresh", refreshtoken, true}}
	testAction(arrt, t)
}

func init() {
	logs.SetDebug()
	tests.PutEnvVal()
	tests.StartAPI()
}
