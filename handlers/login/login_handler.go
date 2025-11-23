package login

import (
	loginpkt "goTibia/packets/login"
	"goTibia/protocol"
	"goTibia/proxy"
	"log"
	"strconv"
	"time"
)

type LoginHandler struct {
	TargetAddr string
	ProxyMOTD  string
	// You could add "DB *sql.DB" here later!
}

func (h *LoginHandler) Handle(client *protocol.Connection) {
	log.Printf("[Login] New Connection: %s", client.RemoteAddr())

	packetReader, err := client.ReadMessage()
	if err != nil {
		log.Printf("error reading message from %s: %v", client.RemoteAddr(), err)
		return
	}

	loginPacket, err := loginpkt.ParseCredentialsPacket(packetReader)
	if err != nil {
		log.Printf("Login: Failed to parse login packet: %v", err)
		return
	}

	protoServerConn, err := proxy.ConnectToBackend(h.TargetAddr)
	if err != nil {
		log.Printf("Login: Failed to connect to %s: %v", client.RemoteAddr(), err)
		return
	}
	defer protoServerConn.Close()

	if err := protoServerConn.SendPacket(loginPacket); err != nil {
		log.Printf("Login: Failed to forward credentials to backend: %v", err)
		return
	}

	log.Println("Login: Credentials forwarded to backend.")

	protoServerConn.EnableXTEA(loginPacket.XTEAKey)
	client.EnableXTEA(loginPacket.XTEAKey)

	message, err := protoServerConn.ReadMessage()
	if err != nil {
		log.Printf("Login: Failed to read server response for %s: %v", client.RemoteAddr(), err)
		return
	}

	loginResultMessage, err := loginpkt.ParseLoginResultMessage(message)
	if err != nil {
		log.Printf("Login: Failed to receive login result message for %s: %v", client.RemoteAddr(), err)
		return
	}

	injectMotd(loginResultMessage, h.ProxyMOTD)
	injectProxyGameworldIP(loginResultMessage)

	err = client.SendPacket(loginResultMessage)
	if err != nil {
		log.Printf("Login: Failed to send login result message for %s: %v", client.RemoteAddr(), err)
		return
	}

	log.Printf("Login: Connection for %s finished.", client.RemoteAddr())
}

func injectMotd(message *loginpkt.LoginResultMessage, motd string) {
	message.Motd = &loginpkt.Motd{
		MotdId:  strconv.Itoa(int(time.Now().Unix())),
		Message: motd,
	}
}

func injectProxyGameworldIP(message *loginpkt.LoginResultMessage) {
	for _, c := range message.CharacterList.Characters {

		// TODO: Extract configuration - do not hardcode
		ip, err := protocol.StringToIP("192.168.1.140")
		if err != nil {
			panic("StringToIP failed")
		}

		c.WorldIp = ip
		c.WorldPort = 7172
		c.WorldName = "Proxy"
	}
}
