package lib

import (
	"encoding/json"
	"flag"
	"github.com/SmartEnergyPlatform/jwt-http-router"
	"github.com/SmartEnergyPlatform/process-deployment/lib/util"
	"github.com/satori/go.uuid"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestQuery(t *testing.T) {
	called := false
	size := 28447 //max known id count
	handler := http.NewServeMux()
	handler.HandleFunc("/test", func(writer http.ResponseWriter, request *http.Request) {
		called = true
		q := request.URL.Query()
		if len(strings.Split(q.Get("ids"), ",")) != size {
			t.Error("unexpected size:", q)
		}
	})
	s := httptest.NewServer(handler)
	defer s.Close()

	ids := []string{}
	for i := 0; i < size; i++ {
		ids = append(ids, uuid.NewV4().String())
	}

	_, err := http.Get(s.URL + "/test?ids=" + strings.Join(ids, ","))
	if err != nil {
		t.Error(err)
	}
	if !called {
		t.Error("not called")
	}
}

func TestDependenciesEndpoint(t *testing.T) {
	testjwt := jwt_http_router.Jwt{}
	err := jwt_http_router.GetJWTPayload(string(ownerJwt), &testjwt)
	if err != nil {
		t.Error(err)
		return
	}
	if testjwt.UserId != owner {
		t.Error("jwt assertion error:", testjwt.UserId, owner)
		return
	}

	eventmanagerMock := http.NewServeMux()
	eventmanagerMock.HandleFunc("/filter/f1", func(writer http.ResponseWriter, request *http.Request) {
		json.NewEncoder(writer).Encode(map[string]string{"state": "running"})
	})
	eventmanagerMock.HandleFunc("/filter/f2", func(writer http.ResponseWriter, request *http.Request) {
		json.NewEncoder(writer).Encode(map[string]string{"state": "running"})
	})
	eventmanagerMock.HandleFunc("/filter/f3", func(writer http.ResponseWriter, request *http.Request) {
		json.NewEncoder(writer).Encode(map[string]string{"state": "running"})
	})
	eventmanagerMockServer := httptest.NewServer(eventmanagerMock)
	defer eventmanagerMockServer.Close()

	closer, mongoport, _, err := testHelper_getMongoDependency()
	defer closer()
	if err != nil {
		t.Error(err)
		return
	}

	configLocation := flag.String("config", "../config.json", "configuration file")
	flag.Parse()

	err = util.LoadConfig(*configLocation)
	if err != nil {
		t.Error(err)
		return
	}

	util.Config.MongoUrl = "mongodb://localhost:" + mongoport
	util.Config.EventManagerUrl = eventmanagerMockServer.URL

	s := httptest.NewServer(getRoutes())
	defer s.Close()

	err = testHelper_putProcess("1", "f1")
	if err != nil {
		t.Error(err)
		return
	}

	err = testHelper_putProcess("2", "f2")
	if err != nil {
		t.Error(err)
		return
	}

	err = testHelper_putProcess("3", "f3")
	if err != nil {
		t.Error(err)
		return
	}

	resp, err := ownerJwt.Get(s.URL + "/deployment/2/dependencies")
	if err != nil {
		t.Error(err)
		return
	}

	metadata := Metadata{}
	err = json.NewDecoder(resp.Body).Decode(&metadata)
	if err != nil {
		t.Error(err)
		return
	}

	if metadata.Process != "2" {
		t.Error("unexpected result", metadata)
		return
	}

	resp, err = ownerJwt.Get(s.URL + "/dependencies?deployments=1,3")
	if err != nil {
		t.Error(err)
		return
	}

	result := []Metadata{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		t.Error(err)
		return
	}

	if len(result) != 2 {
		t.Error("unexpected result length", result)
		return
	}
	for _, metadata := range result {
		if metadata.Process != "1" && metadata.Process != "3" {
			t.Error("unexpected result", metadata)
			return
		}
	}

}
