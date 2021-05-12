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

type BuilderCreator func(*character.Builder) (*character.Model, error)

func CreateFromSeed(l logrus.FieldLogger, db *gorm.DB) func(accountId uint32, worldId byte, name string, jobIndex uint32, face uint32, hair uint32, hairColor uint32, skinColor byte, gender byte, top uint32, bottom uint32, shoes uint32, weapon uint32) (*character.Model, error) {
	return func(accountId uint32, worldId byte, name string, jobIndex uint32, face uint32, hair uint32, hairColor uint32, skinColor byte, gender byte, top uint32, bottom uint32, shoes uint32, weapon uint32) (*character.Model, error) {
		if jobId, ok := job.GetJobFromIndex(jobIndex); ok {
			if bc, ok := getCreator(l, db)(jobId); ok {
				config := character.NewBuilderConfiguration(configuration.Get().UseStarting4Ap, configuration.Get().UseAutoAssignStartersAp)
				c, err := bc(character.NewBuilder(config, accountId, worldId, name, skinColor, gender, hair+hairColor, face))
				if err != nil {
					return nil, err
				}
				addEquippedItems(l, db)(c, top, bottom, shoes, weapon)
				addOtherItems(l, db)(c)
				addSkills(l, db)(c)
				return c, nil
			}
			return nil, errors.New("creator not available for job")
		}
		return nil, errors.New("invalid job creator index")
	}
}

func getCreator(l logrus.FieldLogger, db *gorm.DB) func(jobId uint16) (BuilderCreator, bool) {
	return func(jobId uint16) (BuilderCreator, bool) {
		if jobId == job.Beginner {
			return createBeginner(l, db), true
		} else if jobId == job.Noblesse {
			return createNoblesse(l, db), true
		} else if jobId == job.Legend {
			return createLegend(l, db), true
		}
		return nil, false
	}
}

func createBeginner(l logrus.FieldLogger, db *gorm.DB) func(b *character.Builder) (*character.Model, error) {
	return func(b *character.Builder) (*character.Model, error) {
		b.SetJobId(job.Beginner)
		b.SetMapId(10000)
		return character.Create(l, db)(b)
	}
}

func createNoblesse(l logrus.FieldLogger, db *gorm.DB) func(b *character.Builder) (*character.Model, error) {
	return func(b *character.Builder) (*character.Model, error) {
		b.SetJobId(job.Noblesse)
		b.SetMapId(130030000)
		return character.Create(l, db)(b)
	}
}

func createLegend(l logrus.FieldLogger, db *gorm.DB) func(b *character.Builder) (*character.Model, error) {
	return func(b *character.Builder) (*character.Model, error) {
		b.SetJobId(job.Legend)
		b.SetMapId(914000000)
		return character.Create(l, db)(b)
	}
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
		if job.IsA(c.JobId(), job.Beginner) {
			awardBeginnerSkills(l, db)(c)
		} else if job.IsA(c.JobId(), job.Noblesse) {
			awardNoblesseBeginnerSkills(l, db)(c)
		}
	}
}

func awardBeginnerSkills(l logrus.FieldLogger, db *gorm.DB) func(c *character.Model) {
	return func(c *character.Model) {
		err := skill.AwardSkills(l, db)(c.Id(), skill.BeginnerRecovery, skill.BeginnerNimbleFeet, skill.BeginnerThreeSnails)
		if err != nil {
			l.WithError(err).Errorf("Unable to award character %d beginner skills.", c.Id())
		}
	}
}

func awardNoblesseBeginnerSkills(l logrus.FieldLogger, db *gorm.DB) func(c *character.Model) {
	return func(c *character.Model) {
		err := skill.AwardSkills(l, db)(c.Id(), skill.NoblesseRecovery, skill.NoblesseNimbleFeet, skill.NoblesseThreeSnails)
		if err != nil {
			l.WithError(err).Errorf("Unable to award character %d beginner skills.", c.Id())
		}
	}
}
