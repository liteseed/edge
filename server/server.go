package server

import (
	"github.com/gin-gonic/gin"
	"github.com/liteseed/bungo/server/routes"
)

var APIv1 *gin.RouterGroup

func Register(router *gin.Engine) {
	routes.GetStatus(&router.RouterGroup)
}
