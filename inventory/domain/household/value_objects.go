package household

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

const (
	MinHouseholdNameLength = 3
	MaxHouseholdNameLength = 50

	MinRoomNameLength = 3
	MaxRoomNameLength = 50

	MinItemNameLength = 3
	MaxItemNameLength = 100
)

type HouseholdName string

func NewHouseholdName(name string) (HouseholdName, error) {
	if len(name) < MinHouseholdNameLength || len(name) > MaxHouseholdNameLength {
		return "", fmt.Errorf("household name must be between %d and %d characters: %s", MinHouseholdNameLength, MaxHouseholdNameLength, name)
	}

	return HouseholdName(name), nil
}

func (n HouseholdName) String() string {
	return string(n)
}

type Rooms map[RoomID]Room

func NewRooms() Rooms {
	return make(Rooms)
}

func (r Rooms) Add(room Room) error {
	if _, ok := r.Get(room.ID); ok {
		return fmt.Errorf("room with ID %s already exists", room.ID)
	}

	if _, ok := r.FindByName(room.Name); ok {
		return fmt.Errorf("room with name %s already exists", room.Name)
	}

	r[room.ID] = room

	return nil
}

func (r Rooms) Get(id RoomID) (Room, bool) {
	room, ok := r[id]

	return room, ok
}

func (r Rooms) FindByName(name RoomName) (Room, bool) {
	for _, room := range r {
		if room.Name == name {
			return room, true
		}
	}

	return Room{}, false
}

type Room struct {
	ID    RoomID
	Name  RoomName
	Items Items
}

func NewRoom(id, name string) (Room, error) {
	roomID, err := NewRoomID(id)
	if err != nil {
		return Room{}, err
	}

	roomName, err := NewRoomName(name)
	if err != nil {
		return Room{}, err
	}

	return Room{
		ID:    roomID,
		Name:  roomName,
		Items: make(Items),
	}, nil
}

type RoomID string

func NewRoomID(id string) (RoomID, error) {
	uuid, err := uuid.Parse(id)
	if err != nil {
		return "", fmt.Errorf("invalid room ID. must be valid UUID: %s", id)
	}

	return RoomID(uuid.String()), nil
}

func (id RoomID) String() string {
	return string(id)
}

type RoomName string

func NewRoomName(name string) (RoomName, error) {
	clean := strings.TrimSpace(name)

	if len(clean) < MinRoomNameLength || len(clean) > MaxRoomNameLength {
		return "", fmt.Errorf("room name must be between %d and %d characters", MinRoomNameLength, MaxRoomNameLength)
	}

	return RoomName(clean), nil
}

func (n RoomName) String() string {
	return string(n)
}

type Items map[ItemID]Item

func (i Items) Add(item Item) error {
	if _, ok := i.Get(item.ID); ok {
		return fmt.Errorf("item with ID %s already exists", item.ID)
	}

	i[item.ID] = item

	return nil
}

func (i Items) Update(item Item) error {
	if _, ok := i.Get(item.ID); !ok {
		return fmt.Errorf("item with ID %s does not exist", item.ID)
	}

	i[item.ID] = item

	return nil
}

func (i Items) Remove(id ItemID) {
	delete(i, id)
}

func (i Items) Get(id ItemID) (Item, bool) {
	item, ok := i[id]
	if !ok {
		return Item{}, false
	}

	return item, true
}

type Item struct {
	ID       ItemID
	Name     ItemName
	Barcode  ItemBarcode
	Quantity ItemQuantity
}

func NewItem(id, name, barcode string, quantity uint) (Item, error) {
	itemID, err := NewItemID(id)
	if err != nil {
		return Item{}, err
	}

	itemName, err := NewItemName(name)
	if err != nil {
		return Item{}, err
	}

	itemBarcode, err := NewItemBarcode(barcode)
	if err != nil {
		return Item{}, err
	}

	itemQuantity, err := NewItemQuantity(quantity)
	if err != nil {
		return Item{}, err
	}

	return Item{
		ID:       itemID,
		Name:     itemName,
		Barcode:  itemBarcode,
		Quantity: itemQuantity,
	}, nil
}

func (i Item) Update(name, barcode string, quantity uint) (Item, error) {
	itemName, err := NewItemName(name)
	if err != nil {
		return Item{}, err
	}

	itemBarcode, err := NewItemBarcode(barcode)
	if err != nil {
		return Item{}, err
	}

	itemQuantity, err := NewItemQuantity(quantity)
	if err != nil {
		return Item{}, err
	}

	return Item{
		ID:       i.ID,
		Name:     itemName,
		Barcode:  itemBarcode,
		Quantity: itemQuantity,
	}, nil
}

type ItemID string

func NewItemID(id string) (ItemID, error) {
	uuid, err := uuid.Parse(id)
	if err != nil {
		return "", fmt.Errorf("item ID is invalid. must be valid UUID: %s", id)
	}

	return ItemID(uuid.String()), nil
}

func (id ItemID) String() string {
	return string(id)
}

type ItemName string

func NewItemName(name string) (ItemName, error) {
	clean := strings.TrimSpace(name)

	if len(clean) < MinItemNameLength || len(clean) > MaxItemNameLength {
		return "", fmt.Errorf("item name is invalid. must be between %d and %d characters: %s", MinItemNameLength, MaxItemNameLength, name)
	}

	return ItemName(clean), nil
}

func (n ItemName) String() string {
	return string(n)
}

type ItemBarcode string

func NewItemBarcode(barcode string) (ItemBarcode, error) {
	// TODO: validate barcode
	return ItemBarcode(barcode), nil
}

func (b ItemBarcode) String() string {
	return string(b)
}

type ItemQuantity uint

func NewItemQuantity(quantity uint) (ItemQuantity, error) {
	return ItemQuantity(quantity), nil
}

func (q ItemQuantity) Uint() uint {
	return uint(q)
}
