package main

import (
	"irctalk/common"
	"log"
)

// relay server -> irctalk server 
// chat message와 같은 실시간 푸시 메세지만 전달한다.

func InitHandler(z *common.ZmqMessenger) {
	z.HandleFunc("CHAT", func(msg *common.ZmqMsg) {
		user, err := manager.user.GetConnectedUser(msg.UserId)
		if err != nil {
			return
		}
		irclog := msg.Body().(*common.ZmqChat).Log

		log.Printf("Msg Recv: %+v\n", irclog)
		packet := MakePacket(&SendPushLog{Log:irclog})
		user.Send(packet, nil)
	})

	z.HandleFunc("SERVER_STATUS", func(msg *common.ZmqMsg) {
		user, err := manager.user.GetConnectedUser(msg.UserId)
		if err != nil {
			log.Println("[ZMQMSG]SERVER_STATUS ERROR:", err)
			return
		}
		active := msg.Body().(*common.ZmqServerStatus).Active
		user.ChangeServerActive(msg.ServerId, active)
	})

	z.HandleFunc("ADD_CHANNEL", func(msg *common.ZmqMsg) {
		user, err := manager.user.GetConnectedUser(msg.UserId)
		if err != nil {
			log.Println("[ZMQMSG]ADD_CHANNEL ERROR:", err)
			return
		}

		channel := msg.Body().(*common.ZmqAddChannel).Channel
		packet := MakePacket(&ResAddChannel{Channel:channel})
		user.Send(packet, nil)
	})
}
