package login

import (
	"bytes"
	"encoding/binary"
	"goTibia/protocol"
)

type LoginResultMessage struct {
	ClientDisconnected       bool
	ClientDisconnectedReason string
}

func (lp *LoginResultMessage) Marshal() ([]byte, error) {
	// 1. Create the buffer to build the message.
	buf := new(bytes.Buffer)

	// 2. Write a 2-byte (uint16) placeholder for the length. We'll overwrite it later.
	// We write two zero bytes for now.
	err := binary.Write(buf, binary.LittleEndian, uint16(0))
	if err != nil {
		return nil, err
	}

	// 3. Write the actual message payload.
	if lp.ClientDisconnected {
		buf.WriteByte(S2COpcodeDisconnectClient)
		protocol.WriteString(buf, lp.ClientDisconnectedReason)
	}

	// 4. Get the final byte slice from the buffer.
	finalBytes := buf.Bytes()

	// 5. Calculate the length of the PAYLOAD (total length minus the 2 placeholder bytes).
	payloadLength := len(finalBytes) - 2

	binary.LittleEndian.PutUint16(finalBytes, uint16(payloadLength))

	// 7. Return the complete message with the correct length prefix.
	return finalBytes, nil
}
