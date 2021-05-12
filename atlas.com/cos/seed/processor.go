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
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type BuilderModifier func(*character.Builder) *character.Builder

func CreateFromSeed(l logrus.FieldLogger, db *gorm.DB) func(accountId uint32, worldId byte, name string, jobIndex uint32, face uint32, hair uint32, hairColor uint32, skinColor byte, gender byte, top uint32, bottom uint32, shoes uint32, weapon uint32) (*character.Model, error) {
	return func(accountId uint32, worldId byte, name string, jobIndex uint32, face uint32, hair uint32, hairColor uint32, skinColor byte, gender byte, top uint32, bottom uint32, shoes uint32, weapon uint32) (*character.Model, error) {
		jobId, ok := job.GetJobFromIndex(jobIndex)
		if !ok {
			return nil, errors.New("invalid job creator index")
		}

		bc, ok := getBuilderModifier(jobId)
		if !ok {
			return nil, errors.New("creator not available for job")
		}

		config := character.NewBuilderConfiguration(configuration.Get().UseStarting4Ap, configuration.Get().UseAutoAssignStartersAp)
		builder := character.NewBuilder(config, accountId, worldId, name, skinColor, gender, hair+hairColor, face)
		c, err := character.Create(l, db)(bc(builder))
		if err != nil {
			return nil, err
		}
		addEquippedItems(l, db)(c, top, bottom, shoes, weapon)
		addOtherItems(l, db)(c)
		addSkills(l, db)(c)
		return c, nil
	}
}

func getBuilderModifier(jobId uint16) (BuilderModifier, bool) {
	if jobId == job.Beginner {
		return modifyForBeginner, true
	} else if jobId == job.Noblesse {
		return modifyForNoblesse, true
	} else if jobId == job.Legend {
		return modifyForLegend, true
	}
	return nil, false
}

func modifyForBeginner(b *character.Builder) *character.Builder {
	return b.SetJobId(job.Beginner).SetMapId(10000)
}

func modifyForNoblesse(b *character.Builder) *character.Builder {
	return b.SetJobId(job.Noblesse).SetMapId(130030000)
}

func modifyForLegend(b *character.Builder) *character.Builder {
	return b.SetJobId(job.Legend).SetMapId(914000000)
}

func addEquippedItems(l logrus.FieldLogger, db *gorm.DB) func(c *character.Model, top uint32, bottom uint32, shoes uint32, weapon uint32) {
	return func(c *character.Model, top uint32, bottom uint32, shoes uint32, weapon uint32) {
		equipment.CreateAndEquip(l, db)(c.Id(), top, bottom, shoes, weapon)
	}
}

func addOtherItems(l logrus.FieldLogger, db *gorm.DB) func(c *character.Model) {
	return func(c *character.Model) {
		if job.IsA(c.JobId(), job.Beginner) {
			_, err := item.CreateItemForCharacter(l, db)(c.Id(), inventory.TypeValueETC, 4161001, 1)
			if err != nil {
				l.WithError(err).Errorf("Unable to give character %d item %d.", c.Id(), 4161001)
			}
		} else if job.IsA(c.JobId(), job.Noblesse) {
			_, err := item.CreateItemForCharacter(l, db)(c.Id(), inventory.TypeValueETC, 4161047, 1)
			if err != nil {
				l.WithError(err).Errorf("Unable to give character %d item %d.", c.Id(), 4161047)
			}
		} else if job.IsA(c.JobId(), job.Legend) {
			_, err := item.CreateItemForCharacter(l, db)(c.Id(), inventory.TypeValueETC, 4161048, 1)
			if err != nil {
				l.WithError(err).Errorf("Unable to give character %d item %d.", c.Id(), 4161048)
			}
		}
	}
}

func addSkills(l logrus.FieldLogger, db *gorm.DB) func(c *character.Model) {
	return func(c *character.Model) {
		var skills []uint32
		if job.IsA(c.JobId(), job.Beginner) {
			skills = beginnerSkills()
		} else if job.IsA(c.JobId(), job.Noblesse) {
			noblesseBeginnerSkills()
		}

		err := skill.AwardSkills(l, db)(c.Id(), skills...)
		if err != nil {
			l.WithError(err).Errorf("Unable to award character %d skills.", c.Id())
		}
	}
}

func beginnerSkills() []uint32 {
	return []uint32{skill.BeginnerRecovery, skill.BeginnerNimbleFeet, skill.BeginnerThreeSnails}
}

func noblesseBeginnerSkills() []uint32 {
	return []uint32{skill.NoblesseRecovery, skill.NoblesseNimbleFeet, skill.NoblesseThreeSnails}
}
