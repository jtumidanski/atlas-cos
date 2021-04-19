package character

import (
	"gorm.io/gorm"
	"math/rand"
	"strconv"
	"strings"
)

type EntityUpdateFunction func(e *entity)

func Create(db *gorm.DB, accountId uint32, worldId byte, name string, level byte, strength uint16, dexterity uint16, intelligence uint16, luck uint16, maxHP uint16, maxMP uint16, jobId uint16, gender byte, hair uint32, face uint32, skinColor byte, mapId uint32) (*Model, error) {
	e := &entity{
		AccountId:    accountId,
		World:        worldId,
		Name:         name,
		Level:        level,
		Strength:     strength,
		Dexterity:    dexterity,
		Intelligence: intelligence,
		Luck:         luck,
		HP:           maxHP,
		MP:           maxMP,
		MaxHP:        maxHP,
		MaxMP:        maxMP,
		JobId:        jobId,
		SkinColor:    skinColor,
		Gender:       gender,
		Hair:         hair,
		Face:         face,
		MapId:        mapId,
		SP:           "0, 0, 0, 0, 0, 0, 0, 0, 0, 0",
	}

	err := db.Create(e).Error
	if err != nil {
		return nil, err
	}
	return makeCharacter(e), nil
}

func Update(db *gorm.DB, characterId uint32, modifiers ...EntityUpdateFunction) error {
	c := &entity{ID: characterId}
	err := db.Where(c).First(c).Error
	if err != nil {
		return err
	}

	for _, modifier := range modifiers {
		modifier(c)
	}

	err = db.Save(c).Error
	return err
}

func SetHealth(amount uint16) EntityUpdateFunction {
	return func(e *entity) {
		e.HP = amount
	}
}

func SetMana(amount uint16) EntityUpdateFunction {
	return func(e *entity) {
		e.MP = amount
	}
}

func IncreaseMeso(amount uint32) EntityUpdateFunction {
	return func(e *entity) {
		e.Meso += amount
	}
}

func SetAP(amount uint16) EntityUpdateFunction {
	return func(e *entity) {
		e.AP = amount
	}
}

func IncreaseAP(amount uint16) EntityUpdateFunction {
	return func(e *entity) {
		e.AP += amount
	}
}

func SetStrength(amount uint16) EntityUpdateFunction {
	return func(e *entity) {
		e.Strength = amount
	}
}

func SetDexterity(amount uint16) EntityUpdateFunction {
	return func(e *entity) {
		e.Dexterity = amount
	}
}

func SetIntelligence(amount uint16) EntityUpdateFunction {
	return func(e *entity) {
		e.Intelligence = amount
	}
}

func SetLuck(amount uint16) EntityUpdateFunction {
	return func(e *entity) {
		e.Luck = amount
	}
}

func IncreaseHP(amount uint16) EntityUpdateFunction {
	return func(e *entity) {
		e.HP += amount
		e.MaxHP += amount
	}
}

func IncreaseHPRange(lowerBound uint16, upperBound uint16) EntityUpdateFunction {
	return func(e *entity) {
		amount := uint16(rand.Int31n(int32(upperBound-lowerBound))) + lowerBound
		IncreaseHP(amount)(e)
	}
}

func IncreaseMP(amount uint16) EntityUpdateFunction {
	return func(e *entity) {
		e.MP += amount
		e.MaxMP += amount
	}
}

func IncreaseMPRange(lowerBound uint16, upperBound uint16) EntityUpdateFunction {
	return func(e *entity) {
		amount := uint16(rand.Int31n(int32(upperBound-lowerBound))) + lowerBound
		IncreaseMP(amount)(e)
	}
}

func SpendOnStrength(strength uint16, ap uint16) []EntityUpdateFunction {
	return []EntityUpdateFunction{SetStrength(strength), SetAP(ap)}
}

func SpendOnDexterity(dexterity uint16, ap uint16) []EntityUpdateFunction {
	return []EntityUpdateFunction{SetDexterity(dexterity), SetAP(ap)}
}

func SpendOnIntelligence(intelligence uint16, ap uint16) []EntityUpdateFunction {
	return []EntityUpdateFunction{SetIntelligence(intelligence), SetAP(ap)}
}

func SpendOnLuck(luck uint16, ap uint16) []EntityUpdateFunction {
	return []EntityUpdateFunction{SetLuck(luck), SetAP(ap)}
}

func SetMaxHP(hp uint16) EntityUpdateFunction {
	return func(e *entity) {
		e.MaxHP = hp
	}
}

func SetMaxMP(mp uint16) EntityUpdateFunction {
	return func(e *entity) {
		e.MaxMP = mp
	}
}

func SetMapId(mapId uint32) EntityUpdateFunction {
	return func(e *entity) {
		e.MapId = mapId
	}
}

func SetExperience(experience uint32) EntityUpdateFunction {
	return func(e *entity) {
		e.Experience = experience
	}
}

func IncreaseExperience(gain uint32) EntityUpdateFunction {
	return func(e *entity) {
		e.Experience += gain
	}
}

func UpdateSpawnPoint(spawnPoint uint32) EntityUpdateFunction {
	return func(e *entity) {
		e.SpawnPoint = spawnPoint
	}
}

func SetSP(amount uint32, bookId uint32) EntityUpdateFunction {
	return func(e *entity) {
		sps := strings.Split(e.SP, ",")
		sps[bookId] = strconv.Itoa(int(amount))
		e.SP = strings.Join(sps, ",")
	}
}
