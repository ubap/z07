package packets

import (
	"errors"
	"z07/internal/game/domain"
	"z07/internal/protocol"
	"z07/internal/protocol/crypto"
)

type LoginRequest struct {
	Protocol      uint8
	ClientOS      uint16
	ClientVersion uint16

	XTEAKey       [4]uint32
	Gamemaster    bool
	AccountNumber uint32
	CharacterName string
	Password      string
}

func (lr *LoginRequest) Encode(pw *protocol.PacketWriter) {
	pw.WriteUint8(lr.Protocol)
	pw.WriteUint16(lr.ClientOS)
	pw.WriteUint16(lr.ClientVersion)

	// RSA Encrypted part starts here
	toEncrypt := protocol.NewPacketWriter()

	toEncrypt.WriteUint8(0x00) // Write the check check byte
	toEncrypt.WriteUint32(lr.XTEAKey[0])
	toEncrypt.WriteUint32(lr.XTEAKey[1])
	toEncrypt.WriteUint32(lr.XTEAKey[2])
	toEncrypt.WriteUint32(lr.XTEAKey[3])
	toEncrypt.WriteBool(lr.Gamemaster)
	toEncrypt.WriteUint32(lr.AccountNumber)
	toEncrypt.WriteString(lr.CharacterName)
	toEncrypt.WriteString(lr.Password)

	// Encrypt the data block with the target server's public key.
	unencodedBytes, err := toEncrypt.GetBytes()
	pw.SetError(err)

	encryptedBlock, err := crypto.EncryptRSA(crypto.RSA.GameServerPublicKey, unencodedBytes)
	pw.SetError(err)

	pw.WriteBytes(encryptedBlock)
}

func ParseLoginRequest(packetReader *protocol.PacketReader) (*LoginRequest, error) {
	packet := &LoginRequest{}

	packet.Protocol = packetReader.ReadUint8()
	packet.ClientOS = packetReader.ReadUint16()
	packet.ClientVersion = packetReader.ReadUint16()

	encryptedBlock := packetReader.ReadAll()
	if packetReader.Err() != nil {
		return nil, packetReader.Err()
	}

	decryptedBlock := crypto.DecryptRSA(encryptedBlock)
	decryptedBlockReader := protocol.NewPacketReader(decryptedBlock)
	checkByte := decryptedBlockReader.ReadUint8()
	if checkByte != 0x00 {
		return nil, errors.New("invalid checkByte")
	}

	packet.XTEAKey[0] = decryptedBlockReader.ReadUint32()
	packet.XTEAKey[1] = decryptedBlockReader.ReadUint32()
	packet.XTEAKey[2] = decryptedBlockReader.ReadUint32()
	packet.XTEAKey[3] = decryptedBlockReader.ReadUint32()
	packet.Gamemaster = decryptedBlockReader.ReadBool()
	packet.AccountNumber = decryptedBlockReader.ReadUint32()
	packet.CharacterName = decryptedBlockReader.ReadString()
	packet.Password = decryptedBlockReader.ReadString()

	return packet, packetReader.Err()
}

func NewGameLoginRequest(accountNumber uint32, characterName string, password string) *LoginRequest {
	return &LoginRequest{
		Protocol:      1,
		ClientOS:      1,
		ClientVersion: 772,
		XTEAKey:       [4]uint32{0x11111111, 0x22222222, 0x33333333, 0x44444444},
		Gamemaster:    false,
		AccountNumber: accountNumber,
		CharacterName: characterName,
		Password:      password,
	}
}

func (lr *LoginRequest) GetXTEAKey() [4]uint32 {
	return lr.XTEAKey
}

type LookRequest struct {
	Pos      domain.Position
	ItemId   uint16
	StackPos uint8
}

func ParseLookRequest(pr *protocol.PacketReader) (*LookRequest, error) {
	lr := &LookRequest{}

	lr.Pos = readPosition(pr)
	lr.ItemId = pr.ReadUint16()
	lr.StackPos = pr.ReadUint8()

	return lr, nil
}

type UseItemWithCrosshairRequest struct {
	FromPos      domain.Position
	FromItemId   uint16
	FromStackPos uint8

	ToPos      domain.Position
	ToItemId   uint16
	ToStackPos uint8
}

func ParseUseItemWithCrosshairRequest(pr *protocol.PacketReader) (*UseItemWithCrosshairRequest, error) {
	ur := &UseItemWithCrosshairRequest{}

	ur.FromPos = readPosition(pr)
	ur.FromItemId = pr.ReadUint16()
	ur.FromStackPos = pr.ReadUint8()

	ur.ToPos = readPosition(pr)
	ur.ToItemId = pr.ReadUint16()
	ur.ToStackPos = pr.ReadUint8()

	return ur, nil
}

func (ur *UseItemWithCrosshairRequest) Encode(pw *protocol.PacketWriter) {
	pw.WriteUint8(byte(C2SUseItemWithCrosshair))

	writePosition(pw, ur.FromPos)
	pw.WriteUint16(ur.FromItemId)
	pw.WriteUint8(ur.FromStackPos)

	writePosition(pw, ur.ToPos)
	pw.WriteUint16(ur.ToItemId)
	pw.WriteUint8(ur.ToStackPos)
}
