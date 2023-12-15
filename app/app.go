package app

import (
	"github.com/colinjlacy/mocking-http-requests/endpoints"
	"github.com/colinjlacy/mocking-http-requests/models"
	"github.com/emicklei/go-restful/v3"
	"log"
	"net/http"
)

type app interface {
	Run() error
}

type appImpl struct{}

var (
	App  app = appImpl{}
	user     = endpoints.UserResource
)

func (app appImpl) Run() error {
	ws := app.registerEndpoints()

	restful.DefaultContainer.Add(ws)

	log.Fatal(http.ListenAndServe(":8080", nil))
	return nil
}

func (app appImpl) registerEndpoints() *restful.WebService {
	ws := new(restful.WebService)
	ws.Path("/users").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	ws.Route(ws.GET("/").To(user.FindAllUsers).
		Doc("get all users").
		Writes([]models.User{}).
		Returns(200, "OK", []models.User{}))

	ws.Route(ws.GET("/{user-id}").To(endpoints.UserResource.FindUser).
		Doc("get a user").
		Param(ws.PathParameter("user-id", "identifier of the user").DataType("string").DefaultValue("1")).
		Writes(models.User{}). // on the response
		Returns(200, "OK", models.User{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.PUT("/{user-id}").To(endpoints.UserResource.UpsertUser).
		Doc("update a user").
		Param(ws.PathParameter("user-id", "identifier of the user").DataType("string")).
		Reads(models.User{})) // from the request

	ws.Route(ws.POST("").To(endpoints.UserResource.CreateUser).
		Doc("create a user").
		Reads(models.User{})) // from the request

	ws.Route(ws.DELETE("/{user-id}").To(endpoints.UserResource.RemoveUser).
		Doc("delete a user").
		Param(ws.PathParameter("user-id", "identifier of the user").DataType("string")))

	return ws
}
