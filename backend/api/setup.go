package api

import (
	"github.com/gin-gonic/gin"
)

type API struct {
	router *gin.Engine
}

func New() *API {
	a := &API{
		router: gin.Default(),
	}
	a.registerRoutes()
	return a
}

func (a *API) Run(addr string) error {
	return a.router.Run(addr)
}
