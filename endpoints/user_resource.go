package endpoints

import (
	"fmt"
	"github.com/colinjlacy/mocking-http-requests/models"
	"github.com/emicklei/go-restful/v3"
	"log"
	"net/http"
	"time"
)

type userResource interface {
	FindAllUsers(request *restful.Request, response *restful.Response)
	FindUser(request *restful.Request, response *restful.Response)
	UpsertUser(request *restful.Request, response *restful.Response)
	CreateUser(request *restful.Request, response *restful.Response)
	RemoveUser(request *restful.Request, response *restful.Response)
}

type userResourceImpl struct {
}

var (
	UserResource userResource           = &userResourceImpl{} // normally one would use DAO (data access object)
	users        map[string]models.User = map[string]models.User{
		"1": {"1", "Mario", 35},
		"2": {"2", "Luigi", 32},
		"3": {"3", "Toad", 481},
		"4": {"4", "Peach", 27},
	}
)

// FindAllUsers GET http://localhost:8080/users
func (u *userResourceImpl) FindAllUsers(request *restful.Request, response *restful.Response) {
	log.Println("findAllUsers")
	var list []models.User
	for _, each := range users {
		list = append(list, each)
	}
	response.WriteAsJson(list)
}

// FindUser GET http://localhost:8080/users/1
func (u *userResourceImpl) FindUser(request *restful.Request, response *restful.Response) {
	log.Println("findUser")
	id := request.PathParameter("user-id")
	usr := users[id]
	if len(usr.ID) == 0 {
		response.WriteErrorString(http.StatusNotFound, "models.User could not be found.")
	} else {
		response.WriteAsJson(usr)
	}
}

// UpsertUser PUT http://localhost:8080/users/1
// <models.User><Id>1</Id><Name>Melissa Raspberry</Name></models.User>
func (u *userResourceImpl) UpsertUser(request *restful.Request, response *restful.Response) {
	log.Println("upsertUser")
	usr := models.User{ID: request.PathParameter("user-id")}
	err := request.ReadEntity(&usr)
	if err == nil {
		users[usr.ID] = usr
		response.WriteAsJson(usr)
	} else {
		response.WriteError(http.StatusBadRequest, err)
	}
}

// CreateUser POST http://localhost:8080/users
// <models.User><Id>1</Id><Name>Melissa</Name></models.User>
func (u *userResourceImpl) CreateUser(request *restful.Request, response *restful.Response) {
	log.Println("createUser")
	usr := models.User{}
	err := request.ReadEntity(&usr)
	if err == nil {
		usr.ID = fmt.Sprintf("%d", time.Now().Unix())
		users[usr.ID] = usr
		response.WriteHeaderAndJson(http.StatusCreated, usr, restful.MIME_JSON)
	} else {
		response.WriteError(http.StatusBadRequest, err)
		return
	}
}

// RemoveUser DELETE http://localhost:8080/users/1
func (u *userResourceImpl) RemoveUser(request *restful.Request, response *restful.Response) {
	log.Println("removeUser")
	id := request.PathParameter("user-id")
	delete(users, id)
	response.WriteHeader(http.StatusNoContent)
}
