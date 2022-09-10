package character

import (
	"atlas-cos/database"
	"atlas-cos/model"
	"gorm.io/gorm"
)

func getById(characterId uint32) database.EntityProvider[entity] {
	return func(db *gorm.DB) model.Provider[entity] {
		return database.Query[entity](db, &entity{ID: characterId})
	}
}

func getForAccountInWorld(accountId uint32, worldId byte) database.EntityProvider[[]entity] {
	return func(db *gorm.DB) model.Provider[[]entity] {
		return database.SliceQuery[entity](db, &entity{AccountId: accountId, World: worldId})
	}
}

func getForMapInWorld(worldId byte, mapId uint32) database.EntityProvider[[]entity] {
	return func(db *gorm.DB) model.Provider[[]entity] {
		return database.SliceQuery[entity](db, &entity{World: worldId, MapId: mapId})
	}
}

func getForName(name string) database.EntityProvider[[]entity] {
	return func(db *gorm.DB) model.Provider[[]entity] {
		return database.SliceQuery[entity](db, &entity{Name: name})
	}
}

func makeCharacter(e entity) (Model, error) {
	r := NewModelBuilder().
		SetId(e.ID).
		SetAccountId(e.AccountId).
		SetWorldId(e.World).
		SetName(e.Name).
		SetLevel(e.Level).
		SetExperience(e.Experience).
		SetGachaponExperience(e.GachaponExperience).
		SetStrength(e.Strength).
		SetDexterity(e.Dexterity).
		SetLuck(e.Luck).
		SetIntelligence(e.Intelligence).
		SetHp(e.HP).
		SetMp(e.MP).
		SetMaxHp(e.MaxHP).
		SetMaxMp(e.MaxMP).
		SetMeso(e.Meso).
		SetHpMpUsed(e.HPMPUsed).
		SetJobId(e.JobId).
		SetSkinColor(e.SkinColor).
		SetGender(e.Gender).
		SetFame(e.Fame).
		SetHair(e.Hair).
		SetFace(e.Face).
		SetAp(e.AP).
		SetSp(e.SP).
		SetMapId(e.MapId).
		SetSpawnPoint(e.SpawnPoint).
		SetGm(e.GM).
		Build()
	return r, nil
}
