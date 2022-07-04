package service

import (
	"base/pkg/router"
	// "base/service/index"
	"base/service/users/controller"
)

// LoadRoutes to Load Routes to Router
func LoadRoutes() {

	// Set Endpoint for User Functions
	router.Router.Get(router.RouterBasePath+"/users", controller.GetUser)
	router.Router.Post(router.RouterBasePath+"/users", controller.AddUser)
	router.Router.Get(router.RouterBasePath+"/users/{id}", controller.GetUserByID)
	router.Router.Put(router.RouterBasePath+"/users/{id}", controller.PutUserByID)
	router.Router.Patch(router.RouterBasePath+"/users/{id}", controller.PutUserByID)
	router.Router.Delete(router.RouterBasePath+"/users/{id}", controller.DelUserByID)

}
