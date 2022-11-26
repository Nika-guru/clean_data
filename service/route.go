package service

import (
	"review-service/pkg/router"
	"review-service/service/index"
	"review-service/service/review"
)

// LoadRoutes to Load Routes to Router
func LoadRoutes() {

	// Set Endpoint for admin
	router.Router.Get(router.RouterBasePath+"/", index.GetIndex)
	router.Router.Get(router.RouterBasePath+"/health", index.GetHealth)

	router.Router.Mount(router.RouterBasePath+"/review", review.ExchangeInfoServiceSubRoute)
}
