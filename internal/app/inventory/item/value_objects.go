package item

import (
	"errors"

	"github.com/google/uuid"
)

var (
	ErrItemIDInvalid   = errors.New("item id is invalid")
	ErrItemNameInvalid = errors.New("item name is invalid")
)

const (
	MinItemNameLength = 3
	MaxItemNameLength = 100
)

type Items map[ItemID]Item

func (i Items) Add(item Item) {
	i[item.ID] = item
}

func (i Items) Update(item Item) {
	i[item.ID] = item
}

func (i Items) Remove(id ItemID) {
	delete(i, id)
}

func (i Items) Exists(id ItemID) bool {
	_, ok := i[id]

	return ok
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

func NewItemFromEvent(e ItemAddedEvent) (Item, error) {
	return NewItem(e.ID, e.Name, e.Barcode, e.Quantity)
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

func (i Item) UpdateFromEvent(e ItemUpdatedEvent) (Item, error) {
	return i.Update(e.Name, e.Barcode, e.Quantity)
}

func (i Item) ToItemAddedEvent() ItemAddedEvent {
	return ItemAddedEvent{
		ID:       string(i.ID),
		Name:     string(i.Name),
		Barcode:  string(i.Barcode),
		Quantity: uint(i.Quantity),
	}
}

func (i Item) ToItemUpdatedEvent() ItemUpdatedEvent {
	return ItemUpdatedEvent{
		ID:       string(i.ID),
		Name:     string(i.Name),
		Barcode:  string(i.Barcode),
		Quantity: uint(i.Quantity),
	}
}

type ItemID string

func NewItemID(id string) (ItemID, error) {
	uuid, err := uuid.Parse(id)
	if err != nil {
		return "", ErrItemIDInvalid
	}

	return ItemID(uuid.String()), nil
}

type ItemName string

func NewItemName(name string) (ItemName, error) {
	if len(name) < MinItemNameLength || len(name) > MaxItemNameLength {
		return "", ErrItemNameInvalid
	}

	return ItemName(name), nil
}

type ItemBarcode string

func NewItemBarcode(barcode string) (ItemBarcode, error) {
	// TODO: validate barcode
	return ItemBarcode(barcode), nil
}

type ItemQuantity uint

func NewItemQuantity(quantity uint) (ItemQuantity, error) {
	return ItemQuantity(quantity), nil
}
