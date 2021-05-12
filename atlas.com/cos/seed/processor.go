package seed

import (
	"atlas-cos/character"
	"atlas-cos/configuration"
	"atlas-cos/equipment"
	"atlas-cos/inventory"
	"atlas-cos/item"
	"atlas-cos/job"
	"atlas-cos/skill"
	"errors"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type processor struct {
	l  log.FieldLogger
	db *gorm.DB
}

var Processor = func(l log.FieldLogger, db *gorm.DB) *processor {
	return &processor{l, db}
}

type BuilderCreator func(*character.Builder) (*character.Model, error)

func (p *processor) CreateFromSeed(accountId uint32, worldId byte, name string, jobIndex uint32, face uint32, hair uint32, hairColor uint32, skinColor byte, gender byte, top uint32, bottom uint32, shoes uint32, weapon uint32) (*character.Model, error) {
	if jobId, ok := job.GetJobFromIndex(jobIndex); ok {
		if bc, ok := p.getCreator(jobId); ok {
			config := character.NewBuilderConfiguration(configuration.Get().UseStarting4Ap, configuration.Get().UseAutoAssignStartersAp)
			c, err := bc(character.NewBuilder(config, accountId, worldId, name, skinColor, gender, hair+hairColor, face))
			if err != nil {
				return nil, err
			}
			p.addEquippedItems(c, top, bottom, shoes, weapon)
			p.addOtherItems(c)
			p.addSkills(c)
			return c, nil
		}
		return nil, errors.New("creator not available for job")
	}
	return nil, errors.New("invalid job creator index")
}

func (p *processor) getCreator(jobId uint16) (BuilderCreator, bool) {
	if jobId == job.Beginner {
		return p.createBeginner, true
	} else if jobId == job.Noblesse {
		return p.createNoblesse, true
	} else if jobId == job.Legend {
		return p.createLegend, true
	}
	return nil, false
}

func (p *processor) createBeginner(b *character.Builder) (*character.Model, error) {
	b.SetJobId(job.Beginner)
	b.SetMapId(10000)
	return character.Create(p.l, p.db)(b)
}

func (p *processor) createNoblesse(b *character.Builder) (*character.Model, error) {
	b.SetJobId(job.Noblesse)
	b.SetMapId(130030000)
	return character.Create(p.l, p.db)(b)
}

func (p *processor) createLegend(b *character.Builder) (*character.Model, error) {
	b.SetJobId(job.Legend)
	b.SetMapId(914000000)
	return character.Create(p.l, p.db)(b)
}

func (p *processor) addEquippedItems(c *character.Model, top uint32, bottom uint32, shoes uint32, weapon uint32) {
	equipment.Processor(p.l, p.db).CreateAndEquip(c.Id(), top, bottom, shoes, weapon)
}

func (p *processor) addOtherItems(c *character.Model) {
	if job.IsA(c.JobId(), job.Beginner) {
		_, err := item.Processor(p.l, p.db).CreateItemForCharacter(c.Id(), inventory.TypeValueETC, 4161001, 1)
		if err != nil {
			p.l.WithError(err).Errorf("Unable to give character %d item %d.", c.Id(), 4161001)
		}
	} else if job.IsA(c.JobId(), job.Noblesse) {
		_, err := item.Processor(p.l, p.db).CreateItemForCharacter(c.Id(), inventory.TypeValueETC, 4161047, 1)
		if err != nil {
			p.l.WithError(err).Errorf("Unable to give character %d item %d.", c.Id(), 4161047)
		}
	} else if job.IsA(c.JobId(), job.Legend) {
		_, err := item.Processor(p.l, p.db).CreateItemForCharacter(c.Id(), inventory.TypeValueETC, 4161048, 1)
		if err != nil {
			p.l.WithError(err).Errorf("Unable to give character %d item %d.", c.Id(), 4161048)
		}
	}
}

func (p *processor) addSkills(c *character.Model) {
	if job.IsA(c.JobId(), job.Beginner) {
		p.awardBeginnerSkills(c)
	} else if job.IsA(c.JobId(), job.Noblesse) {
		p.awardNoblesseBeginnerSkills(c)
	}
}

func (p *processor) awardBeginnerSkills(c *character.Model) {
	err := skill.Processor(p.l, p.db).AwardSkills(c.Id(), skill.BeginnerRecovery, skill.BeginnerNimbleFeet, skill.BeginnerThreeSnails)
	if err != nil {
		p.l.WithError(err).Errorf("Unable to award character %d beginner skills.", c.Id())
	}
}

func (p *processor) awardNoblesseBeginnerSkills(c *character.Model) {
	err := skill.Processor(p.l, p.db).AwardSkills(c.Id(), skill.NoblesseRecovery, skill.NoblesseNimbleFeet, skill.NoblesseThreeSnails)
	if err != nil {
		p.l.WithError(err).Errorf("Unable to award character %d beginner skills.", c.Id())
	}
}
