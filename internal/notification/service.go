package notification

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/gorilla/websocket"
	"ws-test/internal/log"
)

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (r *Service) HandleConn(ctx context.Context, c *websocket.Conn) error {
	// subscribed channels
	var subscribeLeaderboard bool
	var subscribeOutcomes bool

	done := make(chan struct{})

	go func() {
		defer func() {
			close(done)
		}()
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				var wsErr *websocket.CloseError
				if errors.As(err, &wsErr) {
					if wsErr.Code == websocket.CloseNormalClosure {
						break
					}
				}

				log.Errorw("read ws", "err", err)
				break
			}

			log.Infow("ws message received", "message", string(message))

			m := &Message{}
			if err := json.Unmarshal(message, &m); err != nil {
				log.Errorw("parse ws message", "err", err)
			}

			if m.Type == NotificationTypeSubscribe {
				switch m.Value {
				case NotificationValueOutcomes:
					subscribeOutcomes = true
				case NotificationValueLeaderboard:
					subscribeLeaderboard = true
				}
			}
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

Loop:
	for {
		select {
		case <-done:
			break Loop
		case t := <-ticker.C:
			if subscribeLeaderboard {
				m := &Message{
					Type:  NotificationValueLeaderboard,
					Value: t.String(),
				}
				buf, _ := json.Marshal(m)
				if err := c.WriteMessage(websocket.TextMessage, buf); err != nil {
					log.Errorw("write ws", "err", err)
					break Loop
				}
				log.Infow("ws message written", "message", string(buf))
			}
			if subscribeOutcomes {
				m := &Message{
					Type:  NotificationValueLeaderboard,
					Value: t.String(),
				}
				buf, _ := json.Marshal(m)
				if err := c.WriteMessage(websocket.TextMessage, buf); err != nil {
					log.Errorw("write ws", "err", err)
					break Loop
				}
				log.Infow("ws message written", "message", string(buf))
			}
		case <-ctx.Done():
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Errorw("write last ws", "err", err)
				break Loop
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			break Loop
		}
	}

	return nil
}
