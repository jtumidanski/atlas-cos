package character

import "gorm.io/gorm"

func GetById(db *gorm.DB, characterId uint32) (*Model, error) {
	var result entity
	err := db.First(&result, characterId).Error
	if err != nil {
		return nil, err
	}
	return makeCharacter(&result), nil
}

func listGet(db *gorm.DB, query interface{}, args ...interface{}) ([]*Model, error) {
	var results []entity
	err := db.First(&results).Where(query, args).Error
	if err != nil {
		return nil, err
	}

	var character = make([]*Model, 0)
	for _, e := range results {
		character = append(character, makeCharacter(&e))
	}
	return character, nil
}

func GetForAccountInWorld(db *gorm.DB, accountId uint32, worldId byte) ([]*Model, error) {
	return listGet(db, "accountId = ? AND worldId = ?", accountId, worldId)
}

func GetForMapInWorld(db *gorm.DB, worldId byte, mapId uint32) ([]*Model, error) {
	return listGet(db, "worldId = ? AND mapId = ?", worldId, mapId)
}

func GetForName(db *gorm.DB, name string) ([]*Model, error) {
	return listGet(db, "name = ?", name)
}

func makeCharacter(e *entity) *Model {
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
	return &r
}
