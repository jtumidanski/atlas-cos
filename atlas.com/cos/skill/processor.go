package skill

import (
	"atlas-cos/job"
	"gorm.io/gorm"
	"log"
)

type processor struct {
	l  *log.Logger
	db *gorm.DB
}

var Processor = func(l *log.Logger, db *gorm.DB) *processor {
	return &processor{l, db}
}

func (p processor) GetSkill(characterId uint32, skillId uint32) (*Model, bool) {
	s, err := GetById(p.db, characterId, skillId)
	if err != nil {
		p.l.Printf("[ERROR] unable to retrieve skill %d for character %d.", skillId, characterId)
		return nil, false
	}
	return s, true
}

func GetSkillBook(jobId uint32) uint32 {
	if jobId >= 2210 && jobId <= 2218 {
		return jobId - 2209
	}
	return 0
}

func IsFourthJob(jobId uint16, skillId uint32) bool {
	if jobId == job.Evan4 {
		return false
	}

	if Is(skillId, EvanMagicMastery, EvanFlameWheel, EvanHerosWill, EvanDarkFog, EvanSoulStone) {
		return true
	}

	return jobId%10 == 2
}

func (p processor) UpdateSkill(characterId uint32, skillId uint32, level uint32, masterLevel uint32, expiration uint64) error {
	if skill, ok := p.GetSkill(characterId, skillId); ok {
		return Update(p.db, skill.Id(), SetLevel(level), SetMasterLevel(masterLevel), SetExpiration(expiration))
	}
	_, err := Create(p.db, characterId, skillId, level, masterLevel, expiration)
	return err
}

func (p processor) GetSkills(characterId uint32) ([]*Model, error) {
	return GetForCharacter(p.db, characterId)
}
