package equipment

import (
	"atlas-cos/rest/requests"
	"errors"
	"gorm.io/gorm"
	"log"
	"strconv"
)

var characterCreationItems = []uint32{
	1302000, 1312004, 1322005, 1442079, // weapons
	1040002, 1040006, 1040010, 1041002, 1041006, 1041010, 1041011, 1042167, // bottom
	1060002, 1060006, 1061002, 1061008, 1062115, // top
	1072001, 1072005, 1072037, 1072038, 1072383, // shoes
	30000, 30010, 30020, 30030, 31000, 31040, 31050, // hair
	20000, 20001, 20002, 21000, 21001, 21002, 21201, 20401, 20402, 21700, 20100, //face
}

type processor struct {
	l  *log.Logger
	db *gorm.DB
}

var Processor = func(l *log.Logger, db *gorm.DB) *processor {
	return &processor{l, db}
}

func (p processor) CreateForCharacter(characterId uint32, itemId uint32, characterCreation bool) (*Model, error) {
	if characterCreation {
		if invalidCharacterCreationItem(itemId) {
			p.l.Printf("[ERROR] received a request to create an item %d for character %d which is not valid for character creation. This is usually a hack.")
			return nil, errors.New("not valid item for character creation")
		}
	}

	nextOpen, err := GetNextFreeEquipmentSlot(p.db, characterId)
	if err != nil {
		nextOpen = 0
	}

	ro, err := requests.EquipmentRegistry().Create(itemId)
	if err != nil {
		p.l.Printf("[ERROR] generating equipment item %d for character %d, they were not awarded this item. Check request in ESO service.")
		return nil, err
	}
	eid, err := strconv.Atoi(ro.Data.Id)
	if err != nil {
		p.l.Printf("[ERROR] generating equipment item %d for character %d, they were not awarded this item. Invalid ID from ESO service.")
		return nil, err
	}
	eq, err := Create(p.db, characterId, uint32(eid), nextOpen)
	if err != nil {
		p.l.Printf("[ERROR] persisting equipment %d association for character %d in slot %d.", eid, characterId, nextOpen)
		return nil, err
	}
	return eq, nil
}

func (p processor) GetEquipmentForCharacter(characterId uint32) ([]*Model, error) {
	return GetEquipmentForCharacter(p.db, characterId)
}

func (p processor) GetEquippedItemForCharacterBySlot(characterId uint32, slot int16) (*Model, error) {
	return GetEquipmentForCharacterBySlot(p.db, characterId, slot)
}

func invalidCharacterCreationItem(itemId uint32) bool {
	for _, v := range characterCreationItems {
		if itemId == v {
			return false
		}
	}
	return true
}
