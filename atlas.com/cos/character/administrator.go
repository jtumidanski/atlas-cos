package character

import "gorm.io/gorm"

func SetHealth(db *gorm.DB, characterId uint32, amount uint16) error {
	c := entity{ID: characterId}
	err := db.Where(&c).First(&c).Error
	if err != nil {
		return err
	}

	c.HP = amount
	err = db.Save(&c).Error
	return err
}

func SetMana(db *gorm.DB, characterId uint32, amount uint16) error {
	c := entity{ID: characterId}
	err := db.Where(&c).First(&c).Error
	if err != nil {
		return err
	}

	c.MP = amount
	err = db.Save(&c).Error
	return err
}

func AdjustMeso(db *gorm.DB, characterId uint32, amount uint32) error {
	c := entity{ID: characterId}
	err := db.Where(&c).First(&c).Error
	if err != nil {
		return err
	}

	c.Meso = c.Meso + amount
	err = db.Save(&c).Error
	return err
}

func SpendOnStrength(db *gorm.DB, characterId uint32, strength uint16, ap uint16) error {
	c := entity{ID: characterId}
	err := db.Where(&c).First(&c).Error
	if err != nil {
		return err
	}

	c.Strength = strength
	c.AP = ap
	err = db.Save(&c).Error
	return err
}

func SpendOnDexterity(db *gorm.DB, characterId uint32, dexterity uint16, ap uint16) error {
	c := entity{ID: characterId}
	err := db.Where(&c).First(&c).Error
	if err != nil {
		return err
	}

	c.Dexterity = dexterity
	c.AP = ap
	err = db.Save(&c).Error
	return err
}

func SpendOnIntelligence(db *gorm.DB, characterId uint32, intelligence uint16, ap uint16) error {
	c := entity{ID: characterId}
	err := db.Where(&c).First(&c).Error
	if err != nil {
		return err
	}

	c.Intelligence = intelligence
	c.AP = ap
	err = db.Save(&c).Error
	return err
}

func SpendOnLuck(db *gorm.DB, characterId uint32, luck uint16, ap uint16) error {
	c := entity{ID: characterId}
	err := db.Where(&c).First(&c).Error
	if err != nil {
		return err
	}

	c.Luck = luck
	c.AP = ap
	err = db.Save(&c).Error
	return err
}

func SetMaxHP(db *gorm.DB, characterId uint32, hp uint16) error {
	c := entity{ID: characterId}
	err := db.Where(&c).First(&c).Error
	if err != nil {
		return err
	}

	c.HP = hp
	err = db.Save(&c).Error
	return err
}

func SetMaxMP(db *gorm.DB, characterId uint32, mp uint16) error {
	c := entity{ID: characterId}
	err := db.Where(&c).First(&c).Error
	if err != nil {
		return err
	}

	c.MP = mp
	err = db.Save(&c).Error
	return err
}

func SetMapId(db *gorm.DB, characterId uint32, mapId uint32) error {
	c := entity{ID: characterId}
	err := db.Where(&c).First(&c).Error
	if err != nil {
		return err
	}

	c.MapId = mapId
	err = db.Save(&c).Error
	return err
}

func SetExperience(db *gorm.DB, characterId uint32, experience uint32) error {
	c := entity{ID: characterId}
	err := db.Where(&c).First(&c).Error
	if err != nil {
		return err
	}

	c.Experience = experience
	err = db.Save(&c).Error
	return err
}

func IncreaseExperience(db *gorm.DB, characterId uint32, gain uint32) error {
	c := entity{ID: characterId}
	err := db.Where(&c).First(&c).Error
	if err != nil {
		return err
	}

	c.Experience = c.Experience + gain
	err = db.Save(&c).Error
	return err
}

func UpdateSpawnPoint(db *gorm.DB, characterId uint32, spawnPoint uint32) error {
	c := entity{ID: characterId}
	err := db.Where(&c).First(&c).Error
	if err != nil {
		return err
	}

	c.SpawnPoint = spawnPoint
	err = db.Save(&c).Error
	return err
}
