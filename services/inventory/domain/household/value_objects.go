package household

import (
	"strings"

	"github.com/bnkamalesh/errors"
	"github.com/google/uuid"
)

const (
	MinHouseholdNameLength = 3
	MaxHouseholdNameLength = 50

	MinHouseholdLocationLength = 3
	MaxHouseholdLocationLength = 50

	MaxHouseholdDescriptionLength = 200

	MinRoomNameLength = 3
	MaxRoomNameLength = 50

	MinItemNameLength = 3
	MaxItemNameLength = 100
)

type HouseholdName string

func NewHouseholdName(name string) (HouseholdName, error) {
	name = strings.TrimSpace(name)

	if len(name) < MinHouseholdNameLength || len(name) > MaxHouseholdNameLength {
		return "", errors.InputBodyf("household name must be between %d and %d characters: %s", MinHouseholdNameLength, MaxHouseholdNameLength, name)
	}

	return HouseholdName(name), nil
}

func (n HouseholdName) String() string {
	return string(n)
}

type HouseholdLocation string

func NewHouseholdLocation(location string) (HouseholdLocation, error) {
	location = strings.TrimSpace(location)

	if len(location) < MinHouseholdLocationLength || len(location) > MaxHouseholdLocationLength {
		return "", errors.InputBodyf("household location must be between %d and %d characters: %s", MinHouseholdLocationLength, MaxHouseholdLocationLength, location)
	}

	return HouseholdLocation(location), nil
}

func (l HouseholdLocation) String() string {
	return string(l)
}

type HouseholdDescription string

func NewHouseholdDescription(description string) (HouseholdDescription, error) {
	description = strings.TrimSpace(description)

	if len(description) > MaxHouseholdDescriptionLength {
		return "", errors.InputBodyf("household description must be less than %d characters", MaxHouseholdDescriptionLength)
	}

	return HouseholdDescription(description), nil
}

func (d HouseholdDescription) String() string {
	return string(d)
}

type Rooms map[RoomID]Room

func NewRooms() Rooms {
	return make(Rooms)
}

func (r Rooms) Add(room Room) error {
	if _, ok := r.Get(room.ID); ok {
		return errors.Duplicatef("room with ID %s already exists", room.ID)
	}

	if _, ok := r.FindByName(room.Name); ok {
		return errors.Duplicatef("room with name %s already exists", room.Name)
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

func (r Rooms) Update(room Room) error {
	if _, ok := r.Get(room.ID); !ok {
		return errors.NotFoundf("room with ID %s does not exist", room.ID)
	}

	if _, ok := r.Without(room.ID).FindByName(room.Name); ok {
		return errors.Duplicatef("room with name %s already exists", room.Name)
	}

	r[room.ID] = room

	return nil
}

func (r Rooms) Without(id RoomID) Rooms {
	rooms := make(Rooms)

	for roomID, room := range r {
		if roomID != id {
			rooms[roomID] = room
		}
	}

	return rooms
}

func (r Rooms) Remove(id RoomID) {
	delete(r, id)
}

func (r Rooms) Count() int {
	return len(r)
}

type Room struct {
	ID    RoomID
	Name  RoomName
	Items Items
	Order uint
}

func NewRoom(id, name string, order uint) (Room, error) {
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
		Order: order,
	}, nil
}

func (r Room) Update(name string) (Room, error) {
	roomName, err := NewRoomName(name)
	if err != nil {
		return Room{}, err
	}

	return Room{
		ID:    r.ID,
		Name:  roomName,
		Items: r.Items,
		Order: r.Order,
	}, nil
}

type RoomID string

func NewRoomID(id string) (RoomID, error) {
	uuid, err := uuid.Parse(id)
	if err != nil {
		return "", errors.InputBodyf("invalid room ID. must be valid UUID: %s", id)
	}

	return RoomID(uuid.String()), nil
}

func (id RoomID) String() string {
	return string(id)
}

type RoomName string

func NewRoomName(name string) (RoomName, error) {
	name = strings.TrimSpace(name)

	if len(name) < MinRoomNameLength || len(name) > MaxRoomNameLength {
		return "", errors.InputBodyf("room name must be between %d and %d characters", MinRoomNameLength, MaxRoomNameLength)
	}

	return RoomName(name), nil
}

func (n RoomName) String() string {
	return string(n)
}

type Items map[ItemID]Item

func (i Items) Add(item Item) error {
	if _, ok := i.Get(item.ID); ok {
		return errors.Duplicatef("item with ID %s already exists", item.ID)
	}

	i[item.ID] = item

	return nil
}

func (i Items) Update(item Item) error {
	if _, ok := i.Get(item.ID); !ok {
		return errors.NotFoundf("item with ID %s does not exist", item.ID)
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

func (i Items) Count() int {
	return len(i)
}

type Item struct {
	ID       ItemID
	Name     ItemName
	Barcode  ItemBarcode
	Quantity uint
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

	return Item{
		ID:       itemID,
		Name:     itemName,
		Barcode:  itemBarcode,
		Quantity: quantity,
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

	return Item{
		ID:       i.ID,
		Name:     itemName,
		Barcode:  itemBarcode,
		Quantity: quantity,
	}, nil
}

type ItemID string

func NewItemID(id string) (ItemID, error) {
	uuid, err := uuid.Parse(id)
	if err != nil {
		return "", errors.InputBodyf("item ID is invalid. must be valid UUID: %s", id)
	}

	return ItemID(uuid.String()), nil
}

func (id ItemID) String() string {
	return string(id)
}

type ItemName string

func NewItemName(name string) (ItemName, error) {
	name = strings.TrimSpace(name)

	if len(name) < MinItemNameLength || len(name) > MaxItemNameLength {
		return "", errors.InputBodyf("item name is invalid. must be between %d and %d characters: %s", MinItemNameLength, MaxItemNameLength, name)
	}

	return ItemName(name), nil
}

func (n ItemName) String() string {
	return string(n)
}

type ItemBarcode string

func NewItemBarcode(barcode string) (ItemBarcode, error) {
	barcode = strings.TrimSpace(barcode)
	// TODO: validate barcode
	return ItemBarcode(barcode), nil
}

func (b ItemBarcode) String() string {
	return string(b)
}
