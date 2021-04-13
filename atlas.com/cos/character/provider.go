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
