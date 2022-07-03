package service

import (
	"base/pkg/router"
	"base/service/index"
	"base/service/users"
)

// LoadRoutes to Load Routes to Router
func LoadRoutes() {
	// Set Endpoint for Root Functions
	router.Router.Get(router.RouterBasePath+"/", index.GetIndex)
	router.Router.Get(router.RouterBasePath+"/health", index.GetHealth)

	// Set Endpoint for Authorization Functions
	// router.Router.With(auth.Basic).Get(router.RouterBasePath+"/auth", index.GetAuth)

	// Set Endpoint for User Functions
	router.Router.Get(router.RouterBasePath+"/users", users.GetUser)
	router.Router.Post(router.RouterBasePath+"/users", users.AddUser)
	router.Router.Get(router.RouterBasePath+"/users/{id}", users.GetUserByID)
	router.Router.Put(router.RouterBasePath+"/users/{id}", users.PutUserByID)
	router.Router.Patch(router.RouterBasePath+"/users/{id}", users.PutUserByID)
	router.Router.Delete(router.RouterBasePath+"/users/{id}", users.DelUserByID)

	// Set Endpoint for Upload Function
	// router.Router.With(auth.JWT).Post(router.RouterBasePath+"/uploads", uploads.UploadFile)
}
