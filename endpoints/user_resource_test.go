package endpoints

import (
	"bytes"
	"github.com/emicklei/go-restful/v3"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUserResourceImpl_FindAllUsers(t *testing.T) {
	rec := httptest.NewRecorder()
	httpReq := httptest.NewRequest(http.MethodGet, "http://test/users", &bytes.Buffer{})
	httpReq.Header.Set("Accept", restful.MIME_JSON)
	req := &restful.Request{
		Request: httpReq,
	}
	res := restful.NewResponse(rec)
	UserResource.FindAllUsers(req, res)
	if res.StatusCode() != http.StatusOK {
		t.Fatalf("expected status code %d, got %d", http.StatusOK, res.StatusCode())
	}
}
