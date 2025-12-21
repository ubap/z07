package packets

// To sort:
//go:generate go run ../../tools/sortcon/main.go

type S2COpcode uint8
type C2SOpcode uint8

const (
	S2CLoginSuccessful     S2COpcode = 0x0A
	S2CLoginAsAdmin        S2COpcode = 0x0B
	S2CServerClosed        S2COpcode = 0x14
	S2CSLoginQueue         S2COpcode = 0x16
	S2CPing                S2COpcode = 0x1E
	S2CMapDescription      S2COpcode = 0x64
	S2CMapSliceNorth       S2COpcode = 0x65
	S2CMapSliceEast        S2COpcode = 0x66
	S2CMapSliceSouth       S2COpcode = 0x67
	S2CMapSliceWest        S2COpcode = 0x68
	S2CAddTileThing        S2COpcode = 0x6A
	S2CUpdateTileItem      S2COpcode = 0x6B
	S2CRemoveTileThing     S2COpcode = 0x6C
	S2CMoveCreature        S2COpcode = 0x6D
	S2COpenContainer       S2COpcode = 0x6E
	S2CCloseContainer      S2COpcode = 0x6F
	S2CAddContainerItem    S2COpcode = 0x70
	S2CUpdateContainerItem S2COpcode = 0x71
	S2CRemoveContainerItem S2COpcode = 0x72
	S2CAddInventoryItem    S2COpcode = 0x78
	S2CRemoveInventoryItem S2COpcode = 0x79
	S2CWorldLight          S2COpcode = 0x82
	S2CMagicEffect         S2COpcode = 0x83
	S2CCreatureHealth      S2COpcode = 0x8C
	S2CCreatureLight       S2COpcode = 0x8D
	S2CPlayerStats         S2COpcode = 0xA0
	S2CPlayerSkills        S2COpcode = 0xA1
	S2CPlayerIcons         S2COpcode = 0xA2
	S2CSay                 S2COpcode = 0xAA
)

const (
	C2SLookRequest C2SOpcode = 0x8C
)
