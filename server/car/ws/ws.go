package ws

import (
	"context"
	"coolcar/car/mq"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

func Handler(u *websocket.Upgrader, sub mq.Subscriber,logger *zap.Logger) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		//http 升级
		// u := &websocket.Upgrader{
		// 	//检测是否同源，同host下的才能使用websocket
		// 	CheckOrigin: func(r *http.Request) bool { return true },
		// }
		c, err := u.Upgrade(w, r, nil)
		if err != nil {
			logger.Warn("cannot upgrade", zap.Error(err))
			return
		}
		defer c.Close()

		message,cleanUp,err:= sub.Subscribe(context.Background())
		defer cleanUp()
		if err != nil {
			logger.Error("cannot subscribe",zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		done := make(chan struct{})
		go func() {
			for {
				_, _, err := c.ReadMessage()
				if err != nil {
					//非正常close
					if !websocket.IsCloseError(err,
						websocket.CloseGoingAway,
						websocket.CloseNormalClosure) {
						logger.Warn("unexpected read error", zap.Error(err))
					}
					done <- struct{}{}
					break
				}
			}
		}() 
		i := 0
		for {
			select {
			case msg:= <- message:
				err:=c.WriteJSON(msg)
				if err != nil {
					logger.Warn("cannot write JSON",zap.Error(err))
				}
			case <-done:
				return
			}

			i++
			err := c.WriteJSON(map[string]string{
				"hello":  "websocket",
				"msg_id": strconv.Itoa(i),
			})
			if err != nil {
				fmt.Printf("cannot write json: %v\n", err)
			}

		}
	}
}
