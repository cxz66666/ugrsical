package zjuicalsrv

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func PingEp(ctx *gin.Context) {
	ctx.String(http.StatusOK, "pong")
}
