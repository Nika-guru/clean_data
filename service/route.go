package service

import (
	"base/pkg/auth"
	"base/pkg/router"
	"base/service/index"
	"base/service/users/controller"
)

// LoadRoutes to Load Routes to Router
func LoadRoutes() {

	// Set Endpoint for Authorization Functions
	router.Router.With(auth.Basic).Get(router.RouterBasePath+"/auth", index.GetAuth)

	// Set Endpoint for admin
	router.Router.Get(router.RouterBasePath+"/", index.GetIndex)
	router.Router.Get(router.RouterBasePath+"/health", index.GetHealth)

	// Set Endpoint for User Functions
	router.Router.Get(router.RouterBasePath+"/users", controller.GetUser)
	router.Router.Post(router.RouterBasePath+"/users", controller.AddUser)
	router.Router.Get(router.RouterBasePath+"/users/{id}", controller.GetUserByID)
	router.Router.Put(router.RouterBasePath+"/users/{id}", controller.PutUserByID)
	router.Router.Patch(router.RouterBasePath+"/users/{id}", controller.PutUserByID)
	router.Router.Delete(router.RouterBasePath+"/users/{id}", controller.DelUserByID)

}
