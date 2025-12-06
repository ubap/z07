package domain

type InventorySlot uint8

const (
	SlotNone     InventorySlot = 0
	SlotHead     InventorySlot = 1
	SlotNeck     InventorySlot = 2
	SlotBackpack InventorySlot = 3
	SlotArmor    InventorySlot = 4
	SlotRight    InventorySlot = 5
	SlotLeft     InventorySlot = 6
	SlotLegs     InventorySlot = 7
	SlotFeet     InventorySlot = 8
	SlotRing     InventorySlot = 9
	SlotAmmo     InventorySlot = 10
)

func (s InventorySlot) String() string {
	switch s {
	case SlotNone:
		return "None"
	case SlotHead:
		return "Head"
	case SlotNeck:
		return "Neck"
	case SlotBackpack:
		return "Backpack"
	case SlotArmor:
		return "Armor"
	case SlotRight:
		return "RightHand"
	case SlotLeft:
		return "LeftHand"
	case SlotLegs:
		return "Legs"
	case SlotFeet:
		return "Feet"
	case SlotRing:
		return "Ring"
	case SlotAmmo:
		return "Ammo"
	default:
		return "UnknownSlot"
	}
}
