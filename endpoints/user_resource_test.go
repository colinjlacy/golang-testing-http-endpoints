package endpoints

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/colinjlacy/mocking-http-requests/models"
	"github.com/emicklei/go-restful/v3"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var testError = "test error"

type errReader struct{}

func (er errReader) Read([]byte) (int, error) {
	return 0, fmt.Errorf(testError)
}

func TestUserResourceImpl_FindAllUsers(t *testing.T) {
	rec := httptest.NewRecorder()
	httpReq := httptest.NewRequest(http.MethodGet, "/users", nil)
	req := restful.NewRequest(httpReq)
	res := restful.NewResponse(rec)
	UserResource.FindAllUsers(req, res)
	if res.StatusCode() != http.StatusOK {
		t.Fatalf("expected status code %d, got %d", http.StatusOK, res.StatusCode())
	}
	var users []models.User
	err := json.Unmarshal(rec.Body.Bytes(), &users)
	if err != nil {
		t.Fatalf("error parsing response body: %s", err.Error())
	}
	if len(users) != 4 {
		t.Fatalf("expected 4 users, found %d", len(users))
	}
}

func TestUserResourceImpl_FindUser(t *testing.T) {
	ws := new(restful.WebService)
	ws.Path("/users").Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON)
	ws.Route(ws.GET("/{user-id}").To(UserResource.FindUser).
		Doc("get a user").
		Param(ws.PathParameter("user-id", "identifier of the user").DataType("string").DefaultValue("1")).
		Writes(models.User{}). // on the response
		Returns(200, "OK", models.User{}).
		Returns(404, "Not Found", nil))
	wc := restful.NewContainer()
	wc.Add(ws)

	rec := httptest.NewRecorder()
	res := restful.NewResponse(rec)
	httpReq := httptest.NewRequest(http.MethodGet, "/users/1", nil)
	wc.ServeHTTP(res, httpReq)

	if res.StatusCode() != http.StatusOK {
		t.Fatalf("expected status code %d, got %d", http.StatusOK, res.StatusCode())
	}
	var user models.User
	err := json.Unmarshal(rec.Body.Bytes(), &user)
	if err != nil {
		t.Fatalf("error parsing response body: %s", err.Error())
	}
	if user.Name != "Mario" {
		t.Fatalf("expected Mario, got %s", user.Name)
	}

}

func TestUserResourceImpl_FindUser_UserNotFound(t *testing.T) {
	ws := new(restful.WebService)
	ws.Path("/users").Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON)
	ws.Route(ws.GET("/{user-id}").To(UserResource.FindUser).
		Doc("get a user").
		Param(ws.PathParameter("user-id", "identifier of the user").DataType("string").DefaultValue("1")).
		Writes(models.User{}). // on the response
		Returns(200, "OK", models.User{}).
		Returns(404, "Not Found", nil))
	wc := restful.NewContainer()
	wc.Add(ws)

	rec := httptest.NewRecorder()
	res := restful.NewResponse(rec)
	httpReq := httptest.NewRequest(http.MethodGet, "/users/notFound", nil)
	wc.ServeHTTP(res, httpReq)

	if res.StatusCode() != http.StatusNotFound {
		t.Fatalf("expected status code %d, got %d", http.StatusNotFound, res.StatusCode())
	}
}

func TestUserResourceImpl_CreateUser(t *testing.T) {
	rec := httptest.NewRecorder()
	httpReq, _ := http.NewRequest(http.MethodPost, "/users", strings.NewReader(`{"Name": "Bowser", "Age": 13}`))
	httpReq.Header.Set("Content-Type", restful.MIME_JSON)
	req := restful.NewRequest(httpReq)
	res := restful.NewResponse(rec)
	UserResource.CreateUser(req, res)

	if res.StatusCode() != http.StatusCreated {
		t.Fatalf("expected status code %d, got %d", http.StatusCreated, res.StatusCode())
	}
	if len(users) != 5 {
		t.Fatalf("expected 5 users, found %d", len(users))
	}
}

func TestUserResourceImpl_CreateUser_BadRequest(t *testing.T) {
	rec := httptest.NewRecorder()
	// passing a struct that throws an error on Read(), error should propagate back to the response
	httpReq, _ := http.NewRequest(http.MethodPost, "/users", errReader{})
	httpReq.Header.Set("Content-Type", restful.MIME_JSON)
	req := restful.NewRequest(httpReq)
	res := restful.NewResponse(rec)
	UserResource.CreateUser(req, res)

	if res.StatusCode() != http.StatusBadRequest {
		t.Fatalf("expected status code %d, got %d", http.StatusBadRequest, res.StatusCode())
	}
	if rec.Body.String() != testError {
		t.Fatalf("expected error as %s, got %s", testError, rec.Body.String())
	}
}

func TestUserResourceImpl_UpsertUser_Update(t *testing.T) {
	ws := new(restful.WebService)
	ws.Path("/users").Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON)
	ws.Route(ws.PUT("/{user-id}").To(UserResource.UpsertUser).
		Doc("update a user").
		Param(ws.PathParameter("user-id", "identifier of the user").DataType("string")).
		Reads(models.User{})) // from the request
	wc := restful.NewContainer()
	wc.Add(ws)

	mario := &models.User{ID: "1", Name: "Mario", Age: 40}
	b, err := json.Marshal(mario)
	if err != nil {
		t.Fatalf("error marshaling test model.User to json: %s", err.Error())
	}
	httpReq, _ := http.NewRequest(http.MethodPut, "/users/1", bytes.NewBuffer(b))
	httpReq.Header.Set("Content-Type", restful.MIME_JSON)
	rec := httptest.NewRecorder()
	res := restful.NewResponse(rec)
	wc.ServeHTTP(res, httpReq)

	if res.StatusCode() != http.StatusOK {
		t.Fatalf("expected status code %d, got %d", http.StatusOK, res.StatusCode())
	}

	if u, _ := users["1"]; u.Age != 40 {
		t.Fatalf("expected Mario's age to be updated to 40, seeing %d", u.Age)
	}
}

func TestUserResourceImpl_UpsertUser_Insert(t *testing.T) {
	ws := new(restful.WebService)
	ws.Path("/users").Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON)
	ws.Route(ws.PUT("/{user-id}").To(UserResource.UpsertUser).
		Doc("update a user").
		Param(ws.PathParameter("user-id", "identifier of the user").DataType("string")).
		Reads(models.User{})) // from the request
	wc := restful.NewContainer()
	wc.Add(ws)

	mario := &models.User{ID: "6", Name: "Rosalina", Age: 200}
	b, err := json.Marshal(mario)
	if err != nil {
		t.Fatalf("error marshaling test model.User to json: %s", err.Error())
	}
	httpReq, _ := http.NewRequest(http.MethodPut, "/users/6", bytes.NewBuffer(b))
	httpReq.Header.Set("Content-Type", restful.MIME_JSON)
	rec := httptest.NewRecorder()
	res := restful.NewResponse(rec)
	wc.ServeHTTP(res, httpReq)

	if res.StatusCode() != http.StatusOK {
		t.Fatalf("expected status code %d, got %d", http.StatusOK, res.StatusCode())
	}

	if u, ok := users["6"]; !ok {
		t.Fatalf("expected a user with the ID of 6, found none")
	} else if u.Name != "Rosalina" {
		t.Fatalf("expected user 6 to have a name Rosalina, got %q", u.Name)
	}
}

func TestUserResourceImpl_UpsertUser_Error(t *testing.T) {
	ws := new(restful.WebService)
	ws.Path("/users").Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON)
	ws.Route(ws.PUT("/{user-id}").To(UserResource.UpsertUser).
		Doc("update a user").
		Param(ws.PathParameter("user-id", "identifier of the user").DataType("string")).
		Reads(models.User{})) // from the request
	wc := restful.NewContainer()
	wc.Add(ws)

	httpReq, _ := http.NewRequest(http.MethodPut, "/users/1", errReader{})
	httpReq.Header.Set("Content-Type", restful.MIME_JSON)
	rec := httptest.NewRecorder()
	res := restful.NewResponse(rec)
	wc.ServeHTTP(res, httpReq)

	if res.StatusCode() != http.StatusBadRequest {
		t.Fatalf("expected status code %d, got %d", http.StatusBadRequest, res.StatusCode())
	}

	if rec.Body.String() != testError {
		t.Fatalf("expected error message %q, got %q", testError, rec.Body.String())
	}
}

func TestUserResourceImpl_RemoveUser(t *testing.T) {
	ws := new(restful.WebService)
	ws.Path("/users").Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON)
	ws.Route(ws.DELETE("/{user-id}").To(UserResource.RemoveUser).
		Doc("delete a user").
		Param(ws.PathParameter("user-id", "identifier of the user").DataType("string")))
	wc := restful.NewContainer()
	wc.Add(ws)

	httpReq, _ := http.NewRequest(http.MethodDelete, "/users/1", nil)
	rec := httptest.NewRecorder()
	res := restful.NewResponse(rec)
	wc.ServeHTTP(res, httpReq)

	if res.StatusCode() != http.StatusNoContent {
		t.Fatalf("expected status code %d, got %d", http.StatusNoContent, res.StatusCode())
	}

	if u, ok := users["1"]; ok {
		t.Fatalf("expected user 1 to be deleted, but %s is still there", u.Name)
	}
}
