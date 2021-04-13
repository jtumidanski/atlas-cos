package character

import (
	"atlas-cos/configuration"
	"atlas-cos/job"
	"atlas-cos/kafka/producers"
	_map "atlas-cos/map"
	"atlas-cos/portal"
	"context"
	"gorm.io/gorm"
	"log"
	"math/rand"
)

type processor struct {
	l  *log.Logger
	db *gorm.DB
}

var Processor = func(l *log.Logger, db *gorm.DB) *processor {
	return &processor{l, db}
}

func (p *processor) AdjustHealth(characterId uint32, amount uint16) {
	c, err := GetById(p.db, characterId)
	if err != nil {
		return
	}
	adjustedAmount := p.enforceBounds(amount, c.HP(), c.MaxHP(), 0)
	err = SetHealth(p.db, characterId, adjustedAmount)
	if err != nil {
		log.Printf("[ERROR] unable to persist health adjustment for character %d.", characterId)
		return
	}
	producers.CharacterStatUpdate(p.l, context.Background()).Emit(characterId, []string{"HP"})
}

func (p *processor) AdjustMana(characterId uint32, amount uint16) {
	c, err := GetById(p.db, characterId)
	if err != nil {
		return
	}
	adjustedAmount := p.enforceBounds(amount, c.MP(), c.MaxMP(), 0)
	err = SetMana(p.db, characterId, adjustedAmount)
	if err != nil {
		log.Printf("[ERROR] unable to persist mana adjustment for character %d.", characterId)
		return
	}
	producers.CharacterStatUpdate(p.l, context.Background()).Emit(characterId, []string{"MP"})
}

func (p *processor) AdjustMeso(characterId uint32, amount uint32, show bool) {
	_, err := GetById(p.db, characterId)
	if err != nil {
		return
	}
	err = AdjustMeso(p.db, characterId, amount)
	if err != nil {
		log.Printf("[ERROR] unable to persist meso adjustment for character %d.", characterId)
		return
	}
	if show {
		producers.MesoGained(p.l, context.Background()).Emit(characterId, amount)
	}
}

func (p *processor) AssignStrength(characterId uint32, strength uint16) {
	c, err := GetById(p.db, characterId)
	if err != nil {
		return
	}

	adjustedStrength := c.Strength() + strength
	if p.outOfRange(adjustedStrength, strength) {
		return
	}
	err = SpendOnStrength(p.db, characterId, adjustedStrength, c.ap-1)
	if err != nil {
		log.Printf("[ERROR] unable to persist strength adjustment for character %d.", characterId)
		return
	}
	producers.CharacterStatUpdate(p.l, context.Background()).Emit(characterId, []string{"STRENGTH, AVAILABLE_AP"})
}

func (p *processor) AssignDexterity(characterId uint32, dexterity uint16) {
	c, err := GetById(p.db, characterId)
	if err != nil {
		return
	}

	adjustedDexterity := c.Dexterity() + dexterity
	if p.outOfRange(adjustedDexterity, dexterity) {
		return
	}
	err = SpendOnDexterity(p.db, characterId, adjustedDexterity, c.ap-1)
	if err != nil {
		log.Printf("[ERROR] unable to persist dexterity adjustment for character %d.", characterId)
		return
	}
	producers.CharacterStatUpdate(p.l, context.Background()).Emit(characterId, []string{"DEXTERITY, AVAILABLE_AP"})
}

func (p *processor) AssignIntelligence(characterId uint32, intelligence uint16) {
	c, err := GetById(p.db, characterId)
	if err != nil {
		return
	}

	adjustedIntelligence := c.Intelligence() + intelligence
	if p.outOfRange(adjustedIntelligence, intelligence) {
		return
	}
	err = SpendOnIntelligence(p.db, characterId, adjustedIntelligence, c.ap-1)
	if err != nil {
		log.Printf("[ERROR] unable to persist intelligence adjustment for character %d.", characterId)
		return
	}
	producers.CharacterStatUpdate(p.l, context.Background()).Emit(characterId, []string{"INTELLIGENCE, AVAILABLE_AP"})
}

func (p *processor) AssignLuck(characterId uint32, luck uint16) {
	c, err := GetById(p.db, characterId)
	if err != nil {
		return
	}

	adjustedLuck := c.Luck() + luck
	if p.outOfRange(adjustedLuck, luck) {
		return
	}
	err = SpendOnLuck(p.db, characterId, adjustedLuck, c.ap-1)
	if err != nil {
		log.Printf("[ERROR] unable to persist luck adjustment for character %d.", characterId)
		return
	}
	producers.CharacterStatUpdate(p.l, context.Background()).Emit(characterId, []string{"LUCK, AVAILABLE_AP"})
}

func (p *processor) AssignHp(characterId uint32) {
	c, err := GetById(p.db, characterId)
	if err != nil {
		return
	}
	adjustedHP := p.calculateHPChange(c, false)
	err = SetMaxHP(p.db, characterId, adjustedHP)
	if err != nil {
		log.Printf("[ERROR] unable to persist HP adjustment for character %d.", characterId)
		return
	}
	producers.CharacterStatUpdate(p.l, context.Background()).Emit(characterId, []string{"MAX_HP"})
}

func (p *processor) AssignMp(characterId uint32) {
	c, err := GetById(p.db, characterId)
	if err != nil {
		return
	}
	adjustedMP := p.calculateMPChange(c, false)
	err = SetMaxMP(p.db, characterId, adjustedMP)
	if err != nil {
		log.Printf("[ERROR] unable to persist MP adjustment for character %d.", characterId)
		return
	}
	producers.CharacterStatUpdate(p.l, context.Background()).Emit(characterId, []string{"MAX_MP"})
}

func (p *processor) AssignSP(characterId uint32, skillId uint32) {

}

func (p *processor) ChangeMap(characterId uint32, worldId byte, channelId byte, mapId uint32, portalId uint32) {
	err := SetMapId(p.db, characterId, mapId)
	if err != nil {
		log.Printf("[ERROR] unable to persist map adjustment for character %d.", characterId)
		return
	}
	por, err := portal.Processor(p.l).GetMapPortalById(mapId, portalId)
	if por != nil {
		GetTemporalRegistry().UpdatePosition(characterId, por.X(), por.Y())
	}
	producers.MapChanged(p.l, context.Background()).Emit(worldId, channelId, mapId, portalId, characterId)
}

func (p *processor) GainExperience(characterId uint32, amount uint32) {
	c, err := GetById(p.db, characterId)
	if err != nil {
		return
	}
	p.gainExperience(characterId, c.Level(), c.MaxClassLevel(), c.Experience(), amount)
}

func (p *processor) GainLevel(characterId uint32) {

}

func (p *processor) enforceBounds(change uint16, current uint16, upperBound uint16, lowerBound uint16) uint16 {
	var adjusted = current + change
	if adjusted < lowerBound {
		return lowerBound
	}
	if adjusted > upperBound {
		return upperBound
	}
	return adjusted
}

func (p *processor) outOfRange(new uint16, change uint16) bool {
	return new < 4 && change != 0 || new > configuration.GetUINT16(configuration.Configuration.MaxAp, 9999)
}

func (p *processor) calculateHPChange(character *Model, usedAPReset bool) uint16 {
	jobId := character.JobId()
	var maxHP uint16 = 0

	if job.IsA(jobId, job.Warrior) || job.IsA(jobId, job.DawnWarrior1) {
		//TODO apply IMPROVED HP INCREASE OR IMPROVED MAX HP
		maxHP = p.adjustHPMPGain(usedAPReset, maxHP, 20, 22, 18, 18, 20)
	} else if job.IsA(jobId, job.Aran1) {
		maxHP = p.adjustHPMPGain(usedAPReset, maxHP, 20, 30, 26, 26, 28)
	} else if job.IsA(jobId, job.Magician) || job.IsA(jobId, job.BlazeWizard1) {
		maxHP = p.adjustHPMPGain(usedAPReset, maxHP, 6, 9, 5, 5, 6)
	} else if job.IsA(jobId, job.Thief) || job.IsA(jobId, job.NightWalker1) {
		maxHP = p.adjustHPMPGain(usedAPReset, maxHP, 16, 18, 14, 14, 16)
	} else if job.IsA(jobId, job.Bowman) || job.IsA(jobId, job.WindArcher1) {
		maxHP = p.adjustHPMPGain(usedAPReset, maxHP, 16, 18, 14, 14, 16)
	} else if job.IsA(jobId, job.Pirate) || job.IsA(jobId, job.ThunderBreaker1) {
		//TODO apply IMPROVE HP INCREASE OR IMPROVE MAX HP
		maxHP = p.adjustHPMPGain(usedAPReset, maxHP, 18, 20, 16, 16, 18)
	} else {
		maxHP = p.adjustHPMPGain(usedAPReset, maxHP, 8, 12, 8, 8, 10)
	}
	return maxHP
}

func (p *processor) adjustHPMPGain(usedAPReset bool, maxHP uint16, apResetAmount uint16, upperBound uint16, lowerBound uint16, floor uint16, staticAmount uint16) uint16 {
	if configuration.GetBool(configuration.Configuration.UseRandomizeHPMPGain, true) {
		if usedAPReset {
			maxHP = maxHP + apResetAmount
		} else {
			maxHP = maxHP + uint16(rand.Int31n(int32(upperBound)-int32(lowerBound))+int32(floor))
		}
	} else {
		maxHP = maxHP + staticAmount
	}
	return maxHP
}

func (p *processor) calculateMPChange(character *Model, usedAPReset bool) uint16 {
	jobId := character.JobId()
	var maxMP uint16 = 0

	if job.IsA(jobId, job.Warrior) || job.IsA(jobId, job.DawnWarrior1) || job.IsA(jobId, job.Aran1) {
		maxMP = p.adjustHPMPGain(usedAPReset, maxMP, 2, 4, 2, character.Intelligence()/10, 3)
	} else if job.IsA(jobId, job.Magician) || job.IsA(jobId, job.BlazeWizard1) {
		//TODO apply IMPROVED MP INCREASE OR IMPROVED MAX MP
		maxMP = p.adjustHPMPGain(usedAPReset, maxMP, 18, 16, 12, character.Intelligence()/20, 18)
	} else if job.IsA(jobId, job.Thief) || job.IsA(jobId, job.NightWalker1) {
		maxMP = p.adjustHPMPGain(usedAPReset, maxMP, 10, 8, 6, character.Intelligence()/10, 10)
	} else if job.IsA(jobId, job.Bowman) || job.IsA(jobId, job.WindArcher1) {
		maxMP = p.adjustHPMPGain(usedAPReset, maxMP, 10, 8, 6, character.Intelligence()/10, 10)
	} else if job.IsA(jobId, job.Pirate) || job.IsA(jobId, job.ThunderBreaker1) {
		maxMP = p.adjustHPMPGain(usedAPReset, maxMP, 14, 9, 7, character.Intelligence()/10, 14)
	} else {
		maxMP = p.adjustHPMPGain(usedAPReset, maxMP, 6, 6, 4, character.Intelligence()/10, 6)
	}
	return maxMP
}

func (p *processor) gainExperience(characterId uint32, level byte, masterLevel byte, experience uint32, gain uint32) {
	if level < masterLevel {
		toNext := GetExperienceNeededForLevel(level) - experience
		if toNext <= gain {
			p.setExperience(characterId, 0)
			producers.CharacterLevel(p.l, context.Background()).Emit(characterId)
			p.gainExperience(characterId, level+1, masterLevel, 0, gain-toNext)
		} else {
			p.increaseExperience(characterId, gain)
		}
	} else {
		p.setExperience(characterId, 0)
	}
}

func (p *processor) setExperience(characterId uint32, experience uint32) {
	err := SetExperience(p.db, characterId, experience)
	if err != nil {
		log.Printf("[ERROR] unable to persist experience adjustment for character %d.", characterId)
		return
	}
	producers.CharacterStatUpdate(p.l, context.Background()).Emit(characterId, []string{"EXPERIENCE"})
}

func (p *processor) increaseExperience(characterId uint32, gain uint32) {
	err := IncreaseExperience(p.db, characterId, gain)
	if err != nil {
		log.Printf("[ERROR] unable to persist experience increase for character %d.", characterId)
		return
	}
	producers.CharacterStatUpdate(p.l, context.Background()).Emit(characterId, []string{"EXPERIENCE"})
}

func (p *processor) MoveCharacter(characterId uint32, x int16, y int16, stance byte) {
	GetTemporalRegistry().Update(characterId, x, y, stance)
	c, err := GetById(p.db, characterId)
	if err != nil {
		return
	}
	spawnPoint, err := _map.Processor(p.l).FindClosestSpawnPoint(c.MapId(), x, y)
	if err != nil {
		return
	}
	err = UpdateSpawnPoint(p.db, characterId, spawnPoint.Id())
	if err != nil {
		log.Printf("[ERROR] unable to persist spawn point for character %d.", characterId)
	}
}

func (p *processor) UpdateStance(characterId uint32, stance byte) {
	GetTemporalRegistry().UpdateStance(characterId, stance)
}

func (p *processor) GetById(characterId uint32) (*Model, error) {
	return GetById(p.db, characterId)
}
