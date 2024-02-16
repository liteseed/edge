package server

import (
	"github.com/gin-gonic/gin"
	"github.com/liteseed/bungo/server/api"
)

var APIv1 *gin.RouterGroup

func Register(router *gin.Engine) {
	api.GetStatus(&router.RouterGroup)
	api.GetTransactionPrice(&router.RouterGroup)
}

func Run(port string) {
	router := gin.Default()

	Register(router)

	router.Run(port)
}
