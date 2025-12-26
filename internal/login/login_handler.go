package login

import (
	"goTibia/internal/login/packets"
	"goTibia/internal/protocol"
	"goTibia/internal/proxy"
	"log"
	"strconv"
	"time"
)

type LoginHandler struct {
	TargetAddr string
	ProxyMOTD  string
}

func (h *LoginHandler) Handle(protoClientConn protocol.Connection) {
	log.Printf("[Login] New Connection: %s", protoClientConn.RemoteAddr())

	_, protoServerConn, err := proxy.InitSession(
		"Login",
		protoClientConn,
		h.TargetAddr,
		packets.ParseCredentialsPacket,
	)
	defer protoServerConn.Close()
	if err != nil {
		log.Printf("[Login]: Failed to initialize session for %s: %v", protoClientConn.RemoteAddr(), err)
		return
	}

	rawMsg, err := protoServerConn.ReadMessage()
	if err != nil {
		log.Printf("[Login]: Failed to read server response for %s: %v", protoClientConn.RemoteAddr(), err)
		return
	}

	packetReader := protocol.NewPacketReader(rawMsg)
	loginResultMessage, err := packets.ParseLoginResultMessage(packetReader)
	if err != nil {
		log.Printf("[Login]: Failed to receive login result message for %s: %v", protoClientConn.RemoteAddr(), err)
		return
	}

	if !loginResultMessage.ClientDisconnected {
		injectMotd(loginResultMessage, h.ProxyMOTD)
		injectProxyGameworldIP(loginResultMessage)
	}

	err = protoClientConn.SendPacket(loginResultMessage)
	if err != nil {
		log.Printf("[Login]: Failed to send login result message for %s: %v", protoClientConn.RemoteAddr(), err)
		return
	}

	log.Printf("[Login]: Connection for %s finished.", protoClientConn.RemoteAddr())
}

func injectMotd(message *packets.LoginResultMessage, motd string) {
	message.Motd = &packets.Motd{
		MotdId:  strconv.Itoa(int(time.Now().Unix())),
		Message: motd,
	}
}

func injectProxyGameworldIP(message *packets.LoginResultMessage) {
	for _, c := range message.CharacterList.Characters {

		// TODO: Extract configuration - do not hardcode
		ip, err := protocol.StringToIP("192.168.1.142")
		if err != nil {
			panic("StringToIP failed")
		}

		c.WorldIp = ip
		c.WorldPort = 7172
		c.WorldName = "Proxy"
	}
}
