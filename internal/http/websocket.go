package http

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"ws-test/internal/log"
)

type NotificationService interface {
	HandleConn(ctx context.Context, c *websocket.Conn) error
}

var upgrader = websocket.Upgrader{}

func (r *Server) handleWebsocket(ctx *gin.Context) {
	c, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Errorw("upgrade:", "err", err)
		return
	}
	defer func() {
		_ = c.Close()
	}()

	if err := r.notificationService.HandleConn(ctx.Request.Context(), c); err != nil {
		log.Errorw("handle ws conn", "err", err)
	}
}
