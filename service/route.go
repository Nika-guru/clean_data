package service

import (
	"crawler/pkg/router"
	"crawler/service/controller"
	"crawler/service/index"
)

// LoadRoutes to Load Routes to Router
func LoadRoutes() {

	// Set Endpoint for admin
	router.Router.Get(router.RouterBasePath+"/", index.GetIndex)
	router.Router.Get(router.RouterBasePath+"/health", index.GetHealth)

	router.Router.Get(router.RouterBasePath+"/info", controller.Info)

}
