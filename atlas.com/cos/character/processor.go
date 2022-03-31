package character

import (
	"atlas-cos/configuration"
	"atlas-cos/database"
	"atlas-cos/equipment"
	"atlas-cos/equipment/statistics"
	"atlas-cos/inventory"
	"atlas-cos/item"
	"atlas-cos/job"
	_map "atlas-cos/map"
	"atlas-cos/portal"
	"atlas-cos/skill"
	"atlas-cos/skill/information"
	"errors"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"math"
	"math/rand"
)

// characterFunc - Function which does something about the character, and returns whether or not further processing should continue.
type characterFunc func(Model) error

// Returns a function which accepts a character model,and updates the persisted state of the character given a set of
// modifying functions.
func characterDatabaseUpdate(_ logrus.FieldLogger, db *gorm.DB) func(modifiers ...EntityUpdateFunction) characterFunc {
	return func(modifiers ...EntityUpdateFunction) characterFunc {
		return func(c Model) error {
			if len(modifiers) > 0 {
				err := update(db, c.Id(), modifiers...)
				if err != nil {
					return err
				}
			}
			return nil
		}
	}
}

// For the characterId, perform the updaterFunc, and if successful, call the successFunc, otherwise log an error.
func characterUpdate(l logrus.FieldLogger, db *gorm.DB) func(characterId uint32, functions ...characterFunc) {
	return func(characterId uint32, functions ...characterFunc) {
		c, err := GetById(l, db)(characterId)
		if err != nil {
			l.WithError(err).Errorf("Unable to locate character %d for update.", characterId)
			return
		}

		err = nil
		for _, f := range functions {
			err = f(c)
			if err != nil {
				l.WithError(err).Errorln("Unable to complete character update.")
				break
			}
		}
	}
}

// AdjustHealth - Adjusts the Health statistic for a character, and emits a CharacterStatUpdateEvent when successful.
func AdjustHealth(l logrus.FieldLogger, db *gorm.DB, span opentracing.Span) func(characterId uint32, amount int16) {
	return func(characterId uint32, amount int16) {
		characterUpdate(l, db)(characterId, persistHealthUpdate(l, db)(amount), healthUpdateSuccess(l, span))
	}
}

// Produces a function which persists a character health update, given the amount, respecting the MaxHP bound.
func persistHealthUpdate(l logrus.FieldLogger, db *gorm.DB) func(amount int16) characterFunc {
	return func(amount int16) characterFunc {
		return func(c Model) error {
			adjustedAmount := enforceBounds(amount, c.HP(), c.MaxHP(), 0)
			l.Debugf("Adjusting health of character %d by %d to %d.", c.Id(), amount, adjustedAmount)
			return characterDatabaseUpdate(l, db)(SetHealth(adjustedAmount))(c)
		}
	}
}

// When a Health update is successful, emit a CharacterStatUpdateEvent.
func healthUpdateSuccess(l logrus.FieldLogger, span opentracing.Span) characterFunc {
	return statisticUpdateSuccess(l, span)("HP")
}

// AdjustMana - Adjusts the Mana statistic for a character, and emits a CharacterStatUpdateEvent when successful.
func AdjustMana(l logrus.FieldLogger, db *gorm.DB, span opentracing.Span) func(characterId uint32, amount int16) {
	return func(characterId uint32, amount int16) {
		characterUpdate(l, db)(characterId, persistManaUpdate(l, db)(amount), manaUpdateSuccess(l, span))
	}
}

// Produces a function which persists a character mana update, given the amount, respecting the MaxMP bound.
func persistManaUpdate(l logrus.FieldLogger, db *gorm.DB) func(amount int16) characterFunc {
	return func(amount int16) characterFunc {
		return func(c Model) error {
			adjustedAmount := enforceBounds(amount, c.MP(), c.MaxMP(), 0)
			l.Debugf("Adjusting mana of character %d by %d to %d.", c.Id(), amount, adjustedAmount)
			return characterDatabaseUpdate(l, db)(SetMana(adjustedAmount))(c)
		}
	}
}

func enforceBounds(change int16, current uint16, upperBound uint16, lowerBound uint16) uint16 {
	var adjusted = int16(current) + change
	return uint16(math.Min(math.Max(float64(adjusted), float64(lowerBound)), float64(upperBound)))
}

// When a Mana update is successful, emit a CharacterStatUpdateEvent.
func manaUpdateSuccess(l logrus.FieldLogger, span opentracing.Span) characterFunc {
	return statisticUpdateSuccess(l, span)("MP")
}

// Produces a function which emits a CharacterStatUpdateEvent for the given characterId and statistic
func statisticUpdateSuccess(l logrus.FieldLogger, span opentracing.Span) func(statistic string) characterFunc {
	return func(statistic string) characterFunc {
		return func(c Model) error {
			emitStatUpdateEvent(l, span)(c.Id(), []string{statistic})
			return nil
		}
	}
}

// Produces a function which emits a CharacterStatUpdateEvent for the given characterId and statistic
func statisticsUpdateSuccess(l logrus.FieldLogger, span opentracing.Span) func(statistics ...string) characterFunc {
	return func(statistics ...string) characterFunc {
		return func(c Model) error {
			emitStatUpdateEvent(l, span)(c.Id(), statistics)
			return nil
		}
	}
}

// AdjustMeso - Adjusts the Meso count for a character, and emits a MesoGainedEvent when successful.
func AdjustMeso(l logrus.FieldLogger, db *gorm.DB, span opentracing.Span) func(characterId uint32, amount int32, show bool) {
	return func(characterId uint32, amount int32, show bool) {
		characterUpdate(l, db)(characterId, persistMesoUpdate(l, db)(amount), statisticUpdateSuccess(l, span)("MESO"), mesoUpdateSuccess(l, span)(amount, show))
	}
}

// Produces a function which persists a character meso update, given the amount.
func persistMesoUpdate(l logrus.FieldLogger, db *gorm.DB) func(amount int32) characterFunc {
	return func(amount int32) characterFunc {
		return func(c Model) error {
			final := uint32(math.Max(0, float64(amount)+float64(c.Meso())))
			l.Debugf("Adjusting meso of character %d by %d to %d.", c.Id(), amount, final)
			return characterDatabaseUpdate(l, db)(SetMeso(final))(c)
		}
	}
}

// Produces a function which emits a MesoGainedEvent for the given characterId and amount.
func mesoUpdateSuccess(l logrus.FieldLogger, span opentracing.Span) func(amount int32, show bool) characterFunc {
	return func(amount int32, show bool) characterFunc {
		return func(c Model) error {
			if show {
				emitMesoGainedEvent(l, span)(c.Id(), amount)
			}
			return nil
		}
	}
}

// ChangeMap - Changes the map for a character in the database, updates the temporal position, and emits a MapChangedEvent when successful.
func ChangeMap(l logrus.FieldLogger, db *gorm.DB, span opentracing.Span) func(characterId uint32, worldId byte, channelId byte, mapId uint32, portalId uint32) {
	return func(characterId uint32, worldId byte, channelId byte, mapId uint32, portalId uint32) {
		characterUpdate(l, db)(characterId, performChangeMap(l, db, span)(mapId, portalId), changeMapSuccess(l, span)(worldId, channelId, mapId, portalId))
	}
}

// Produces a function which persists a character map update, then updates the temporal position.
func performChangeMap(l logrus.FieldLogger, db *gorm.DB, span opentracing.Span) func(mapId uint32, portalId uint32) characterFunc {
	return func(mapId uint32, portalId uint32) characterFunc {
		return func(c Model) error {
			err := characterDatabaseUpdate(l, db)(SetMapId(mapId))(c)
			if err != nil {
				return err
			}
			por, err := portal.GetById(l, span)(mapId, portalId)
			if err != nil {
				return err
			}
			GetTemporalRegistry().UpdatePosition(c.Id(), por.X(), por.Y())
			return nil
		}
	}
}

// Produces a function which emits a MapChangedEvent for the given characterId.
func changeMapSuccess(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, mapId uint32, portalId uint32) characterFunc {
	return func(worldId byte, channelId byte, mapId uint32, portalId uint32) characterFunc {
		return func(c Model) error {
			emitMapChangedEvent(l, span)(worldId, channelId, mapId, portalId, c.Id())
			return nil
		}
	}
}

// GainExperience - Updates the character based on the experience gained, may trigger level updates depending on amount gained.
func GainExperience(l logrus.FieldLogger, db *gorm.DB, span opentracing.Span) func(characterId uint32, amount uint32) {
	return func(characterId uint32, amount uint32) {
		characterUpdate(l, db)(characterId, performGainExperience(l, db, span)(amount))
	}
}

func performGainExperience(l logrus.FieldLogger, db *gorm.DB, span opentracing.Span) func(amount uint32) characterFunc {
	return func(amount uint32) characterFunc {
		return func(c Model) error {
			gainExperience(l, db, span)(c.Id(), c.Level(), c.MaxClassLevel(), c.Experience(), amount)
			return nil
		}
	}
}

func gainExperience(l logrus.FieldLogger, db *gorm.DB, span opentracing.Span) func(characterId uint32, level byte, masterLevel byte, experience uint32, gain uint32) {
	return func(characterId uint32, level byte, masterLevel byte, experience uint32, gain uint32) {
		if level < masterLevel {
			l.Debugf("Character %d received experience, and is not max level.", characterId)
			maxNext := GetExperienceNeededForLevel(level)
			toNext := maxNext - experience
			l.Debugf("Character %d needs a total of %d experience for their level. They have %d, and gained %d.", characterId, maxNext, experience, gain)
			if toNext <= gain {
				l.Debugf("Character %d leveled. Set experience to 0 during the level, and perform level.", characterId)
				setExperience(l, db, span)(characterId, 0)
				emitLevelEvent(l, span)(characterId)
				if gain-toNext > 0 {
					l.Debugf("Character %d has %d experience left to process.", characterId, gain-toNext)
					gainExperience(l, db, span)(characterId, level+1, masterLevel, 0, gain-toNext)
				}
			} else {
				l.Debugf("Character %d received less experience than what is needed to level.", characterId)
				increaseExperience(l, db, span)(characterId, gain)
			}
		} else {
			l.Debugf("Character %d received experience while at max level, retain 0 experience.", characterId)
			setExperience(l, db, span)(characterId, 0)
		}
	}
}

func persistSetExperience(l logrus.FieldLogger, db *gorm.DB) func(experience uint32) characterFunc {
	return func(experience uint32) characterFunc {
		return func(c Model) error {
			l.Debugf("Setting character %d experience to %d from %d.", c.Id(), experience, c.Experience())
			return characterDatabaseUpdate(l, db)(SetExperience(experience))(c)
		}
	}
}

func experienceChangeSuccess(l logrus.FieldLogger, span opentracing.Span) characterFunc {
	return statisticUpdateSuccess(l, span)("EXPERIENCE")
}

func setExperience(l logrus.FieldLogger, db *gorm.DB, span opentracing.Span) func(characterId uint32, experience uint32) {
	return func(characterId uint32, experience uint32) {
		characterUpdate(l, db)(characterId, persistSetExperience(l, db)(experience), experienceChangeSuccess(l, span))
	}
}

func persistIncreaseExperience(l logrus.FieldLogger, db *gorm.DB) func(experience uint32) characterFunc {
	return func(experience uint32) characterFunc {
		return func(c Model) error {
			l.Debugf("Increasing character %d experience by %d to %d.", c.Id(), experience, experience+c.Experience())
			return characterDatabaseUpdate(l, db)(SetExperience(experience + c.Experience()))(c)
		}
	}
}

func increaseExperience(l logrus.FieldLogger, db *gorm.DB, span opentracing.Span) func(characterId uint32, gain uint32) {
	return func(characterId uint32, gain uint32) {
		characterUpdate(l, db)(characterId, persistIncreaseExperience(l, db)(gain), experienceChangeSuccess(l, span))
	}
}

func MoveCharacter(l logrus.FieldLogger, db *gorm.DB, span opentracing.Span) func(characterId uint32, x int16, y int16, stance byte) {
	return func(characterId uint32, x int16, y int16, stance byte) {
		characterUpdate(l, db)(characterId, updateTemporalPosition(x, y, stance), updateSpawnPoint(l, db, span)(x, y))
	}
}

func updateTemporalPosition(x int16, y int16, stance byte) characterFunc {
	return func(c Model) error {
		GetTemporalRegistry().Update(c.Id(), x, y, stance)
		return nil
	}
}

func updateSpawnPoint(l logrus.FieldLogger, db *gorm.DB, span opentracing.Span) func(x int16, y int16) characterFunc {
	return func(x int16, y int16) characterFunc {
		return func(c Model) error {
			sp, err := _map.FindClosestSpawnPoint(l, span)(c.MapId(), x, y)
			if err != nil {
				return err
			}
			return characterDatabaseUpdate(l, db)(UpdateSpawnPoint(sp.Id()))(c)
		}
	}
}

func UpdateStance(l logrus.FieldLogger, db *gorm.DB) func(characterId uint32, stance byte) {
	return func(characterId uint32, stance byte) {
		characterUpdate(l, db)(characterId, updateStance(stance))
	}
}

func updateStance(stance byte) characterFunc {
	return func(c Model) error {
		GetTemporalRegistry().UpdateStance(c.Id(), stance)
		return nil
	}
}

func GetById(_ logrus.FieldLogger, db *gorm.DB) func(characterId uint32) (Model, error) {
	return func(characterId uint32) (Model, error) {
		return database.ModelProvider[Model, entity](db)(getById(characterId), makeCharacter)()
	}
}

func outOfRange(new uint16, change uint16) bool {
	return new < 4 && change != 0 || new > configuration.Get().MaxAp
}

func persistAttributeUpdate(l logrus.FieldLogger, db *gorm.DB) func(getter func(Model) uint16, modifierGetter func(uint16, uint16) []EntityUpdateFunction) characterFunc {
	return func(getter func(Model) uint16, modifierGetter func(uint16, uint16) []EntityUpdateFunction) characterFunc {
		return func(c Model) error {
			value := getter(c) + 1
			if outOfRange(value, 1) {
				return nil
			}
			return characterDatabaseUpdate(l, db)(modifierGetter(value, c.AP()-1)...)(c)
		}
	}
}

func AssignStrength(l logrus.FieldLogger, db *gorm.DB, span opentracing.Span) func(characterId uint32) {
	return func(characterId uint32) {
		characterUpdate(l, db)(characterId, persistStrengthUpdate(l, db), strengthUpdateSuccess(l, span))
	}
}

func persistStrengthUpdate(l logrus.FieldLogger, db *gorm.DB) characterFunc {
	return persistAttributeUpdate(l, db)(Model.Strength, SpendOnStrength)
}

func strengthUpdateSuccess(l logrus.FieldLogger, span opentracing.Span) characterFunc {
	return statisticsUpdateSuccess(l, span)("STRENGTH", "AVAILABLE_AP")
}

func AssignDexterity(l logrus.FieldLogger, db *gorm.DB, span opentracing.Span) func(characterId uint32) {
	return func(characterId uint32) {
		characterUpdate(l, db)(characterId, persistDexterityUpdate(l, db), dexterityUpdateSuccess(l, span))
	}
}

func persistDexterityUpdate(l logrus.FieldLogger, db *gorm.DB) characterFunc {
	return persistAttributeUpdate(l, db)(Model.Dexterity, SpendOnDexterity)
}

func dexterityUpdateSuccess(l logrus.FieldLogger, span opentracing.Span) characterFunc {
	return statisticsUpdateSuccess(l, span)("DEXTERITY", "AVAILABLE_AP")
}

func AssignIntelligence(l logrus.FieldLogger, db *gorm.DB, span opentracing.Span) func(characterId uint32) {
	return func(characterId uint32) {
		characterUpdate(l, db)(characterId, persistIntelligenceUpdate(l, db), intelligenceUpdateSuccess(l, span))
	}
}

func persistIntelligenceUpdate(l logrus.FieldLogger, db *gorm.DB) characterFunc {
	return persistAttributeUpdate(l, db)(Model.Intelligence, SpendOnIntelligence)
}

func intelligenceUpdateSuccess(l logrus.FieldLogger, span opentracing.Span) characterFunc {
	return statisticsUpdateSuccess(l, span)("INTELLIGENCE", "AVAILABLE_AP")
}

func AssignLuck(l logrus.FieldLogger, db *gorm.DB, span opentracing.Span) func(characterId uint32) {
	return func(characterId uint32) {
		characterUpdate(l, db)(characterId, persistLuckUpdate(l, db), luckUpdateSuccess(l, span))
	}
}

func persistLuckUpdate(l logrus.FieldLogger, db *gorm.DB) characterFunc {
	return persistAttributeUpdate(l, db)(Model.Luck, SpendOnLuck)
}

func luckUpdateSuccess(l logrus.FieldLogger, span opentracing.Span) characterFunc {
	return statisticsUpdateSuccess(l, span)("LUCK", "AVAILABLE_AP")
}

func AssignHp(l logrus.FieldLogger, db *gorm.DB, span opentracing.Span) func(characterId uint32) {
	return func(characterId uint32) {
		characterUpdate(l, db)(characterId, persistHpUpdate(l, db, span), hpUpdateSuccess(l, span))
	}
}

func persistHpUpdate(l logrus.FieldLogger, db *gorm.DB, span opentracing.Span) characterFunc {
	return func(c Model) error {
		adjustedHP := calculateHPChange(l, db, span)(c, false)
		return characterDatabaseUpdate(l, db)(SetMaxHP(adjustedHP))(c)
	}
}

func calculateHPChange(l logrus.FieldLogger, db *gorm.DB, span opentracing.Span) func(c Model, usedAPReset bool) uint16 {
	return func(c Model, usedAPReset bool) uint16 {
		var maxHP uint16 = 0
		if job.IsA(c.JobId(), job.Warrior, job.DawnWarrior1) {
			if !usedAPReset {
				maxHP += getHPIncreaseEnhancementForAPAssignment(l, db, span)(c)
			}
			maxHP = adjustHPMPGain(usedAPReset, maxHP, 20, 22, 18, 18, 20)
		} else if job.IsA(c.JobId(), job.Aran1) {
			maxHP = adjustHPMPGain(usedAPReset, maxHP, 20, 30, 26, 26, 28)
		} else if job.IsA(c.JobId(), job.Magician, job.BlazeWizard1) {
			maxHP = adjustHPMPGain(usedAPReset, maxHP, 6, 9, 5, 5, 6)
		} else if job.IsA(c.JobId(), job.Thief, job.NightWalker1) {
			maxHP = adjustHPMPGain(usedAPReset, maxHP, 16, 18, 14, 14, 16)
		} else if job.IsA(c.JobId(), job.Bowman, job.WindArcher1) {
			maxHP = adjustHPMPGain(usedAPReset, maxHP, 16, 18, 14, 14, 16)
		} else if job.IsA(c.JobId(), job.Pirate, job.ThunderBreaker1) {
			if !usedAPReset {
				maxHP += getHPIncreaseEnhancementForAPAssignment(l, db, span)(c)
			}
			maxHP = adjustHPMPGain(usedAPReset, maxHP, 18, 20, 16, 16, 18)
		} else {
			maxHP = adjustHPMPGain(usedAPReset, maxHP, 8, 12, 8, 8, 10)
		}
		return maxHP
	}
}

func getHPIncreaseEnhancementSkill(c Model) uint32 {
	var skillId uint32
	if job.IsA(c.JobId(), job.DawnWarrior1) {
		skillId = skill.DawnWarriorMaxHPEnhancement
	} else if job.IsA(c.JobId(), job.Warrior) {
		skillId = skill.WarriorImprovedHPIncrease
	} else if job.IsA(c.JobId(), job.ThunderBreaker1) {
		skillId = skill.ThunderBreakerImproveMaxHP
	} else if job.IsA(c.JobId(), job.Pirate) {
		skillId = skill.BrawlerImproveMaxHP
	}
	return skillId
}

func getHPIncreaseEnhancementForAPAssignment(l logrus.FieldLogger, db *gorm.DB, span opentracing.Span) func(c Model) uint16 {
	return func(c Model) uint16 {
		return getHPIncreaseEnhancementAmount(l, db, span)(c, func(effect *information.Effect) uint16 {
			return uint16(effect.Y())
		})
	}
}

func getHPIncreaseEnhancementForLevelUp(l logrus.FieldLogger, db *gorm.DB, span opentracing.Span) func(c Model) uint16 {
	return func(c Model) uint16 {
		return getHPIncreaseEnhancementAmount(l, db, span)(c, func(effect *information.Effect) uint16 {
			return uint16(effect.X())
		})
	}
}

func getHPIncreaseEnhancementAmount(l logrus.FieldLogger, db *gorm.DB, span opentracing.Span) func(c Model, getter func(*information.Effect) uint16) uint16 {
	return func(c Model, getter func(*information.Effect) uint16) uint16 {
		skillId := getHPIncreaseEnhancementSkill(c)
		if effect, ok := skill.IfHasSkillGetEffect(l, db, span)(c.Id(), skillId); ok {
			val := getter(effect)
			l.Debugf("Character %d HP increase impacted by %d due to skill %d.", c.Id(), val, skillId)
			return val
		}
		return 0
	}
}

func adjustHPMPGain(usedAPReset bool, maxHP uint16, apResetAmount uint16, upperBound uint16, lowerBound uint16, floor uint16, staticAmount uint16) uint16 {
	if configuration.Get().UseRandomizeHpMpGain {
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

func hpUpdateSuccess(l logrus.FieldLogger, span opentracing.Span) characterFunc {
	return statisticsUpdateSuccess(l, span)("MAX_HP")
}

func AssignMp(l logrus.FieldLogger, db *gorm.DB, span opentracing.Span) func(characterId uint32) {
	return func(characterId uint32) {
		characterUpdate(l, db)(characterId, persistMpUpdate(l, db, span), mpUpdateSuccess(l, span))
	}
}

func persistMpUpdate(l logrus.FieldLogger, db *gorm.DB, span opentracing.Span) characterFunc {
	return func(c Model) error {
		adjustedMP := calculateMPChange(l, db, span)(c, false)
		return characterDatabaseUpdate(l, db)(SetMaxMP(adjustedMP))(c)
	}
}

func getMPIncreaseEnhancementSkill(c Model) uint32 {
	var skillId uint32
	if job.IsA(c.JobId(), job.BlazeWizard1) {
		skillId = skill.BlazeWizardIncreasingMaxMP
	} else {
		skillId = skill.MagicianImprovedMPIncrease
	}
	return skillId
}

func getMPIncreaseEnhancementForAPAssignment(l logrus.FieldLogger, db *gorm.DB, span opentracing.Span) func(c Model) uint16 {
	return func(c Model) uint16 {
		return getMPIncreaseEnhancementAmount(l, db, span)(c, func(effect *information.Effect) uint16 {
			return uint16(effect.Y())
		})
	}
}

func getMPIncreaseEnhancementForLevelUp(l logrus.FieldLogger, db *gorm.DB, span opentracing.Span) func(c Model) uint16 {
	return func(c Model) uint16 {
		return getMPIncreaseEnhancementAmount(l, db, span)(c, func(effect *information.Effect) uint16 {
			return uint16(effect.X())
		})
	}
}

func getMPIncreaseEnhancementAmount(l logrus.FieldLogger, db *gorm.DB, span opentracing.Span) func(c Model, getter func(*information.Effect) uint16) uint16 {
	return func(c Model, getter func(*information.Effect) uint16) uint16 {
		skillId := getMPIncreaseEnhancementSkill(c)
		if effect, ok := skill.IfHasSkillGetEffect(l, db, span)(c.Id(), skillId); ok {
			val := getter(effect)
			l.Debugf("Character %d MP increase impacted by %d due to skill %d.", c.Id(), val, skillId)
			return val
		}
		return 0
	}
}

func calculateMPChange(l logrus.FieldLogger, db *gorm.DB, span opentracing.Span) func(c Model, usedAPReset bool) uint16 {
	return func(c Model, usedAPReset bool) uint16 {
		jobId := c.JobId()
		var maxMP uint16 = 0

		if job.IsA(jobId, job.Warrior, job.DawnWarrior1, job.Aran1) {
			maxMP = adjustHPMPGain(usedAPReset, maxMP, 2, 4, 2, c.Intelligence()/10, 3)
		} else if job.IsA(jobId, job.Magician, job.BlazeWizard1) {
			if !usedAPReset {
				maxMP += getMPIncreaseEnhancementForAPAssignment(l, db, span)(c)
			}

			maxMP = adjustHPMPGain(usedAPReset, maxMP, 18, 16, 12, c.Intelligence()/20, 18)
		} else if job.IsA(jobId, job.Thief, job.NightWalker1) {
			maxMP = adjustHPMPGain(usedAPReset, maxMP, 10, 8, 6, c.Intelligence()/10, 10)
		} else if job.IsA(jobId, job.Bowman, job.WindArcher1) {
			maxMP = adjustHPMPGain(usedAPReset, maxMP, 10, 8, 6, c.Intelligence()/10, 10)
		} else if job.IsA(jobId, job.Pirate, job.ThunderBreaker1) {
			maxMP = adjustHPMPGain(usedAPReset, maxMP, 14, 9, 7, c.Intelligence()/10, 14)
		} else {
			maxMP = adjustHPMPGain(usedAPReset, maxMP, 6, 6, 4, c.Intelligence()/10, 6)
		}
		return maxMP
	}
}

func mpUpdateSuccess(l logrus.FieldLogger, span opentracing.Span) characterFunc {
	return statisticsUpdateSuccess(l, span)("MAX_MP")
}

func TotalIntelligence(l logrus.FieldLogger, db *gorm.DB, span opentracing.Span) func(c Model) uint16 {
	return func(c Model) uint16 {
		return totalStat(l, db, span)(c, (*Model).Intelligence, statistics.Model.Intelligence)
	}
}

func TotalDexterity(l logrus.FieldLogger, db *gorm.DB, span opentracing.Span) func(c Model) uint16 {
	return func(c Model) uint16 {
		return totalStat(l, db, span)(c, (*Model).Dexterity, statistics.Model.Dexterity)
	}
}

func TotalStrength(l logrus.FieldLogger, db *gorm.DB, span opentracing.Span) func(c Model) uint16 {
	return func(c Model) uint16 {
		return totalStat(l, db, span)(c, (*Model).Strength, statistics.Model.Strength)
	}
}

func TotalLuck(l logrus.FieldLogger, db *gorm.DB, span opentracing.Span) func(c Model) uint16 {
	return func(c Model) uint16 {
		return totalStat(l, db, span)(c, (*Model).Luck, statistics.Model.Luck)
	}
}

func totalStat(l logrus.FieldLogger, db *gorm.DB, span opentracing.Span) func(c Model, baseGetter func(*Model) uint16, equipGetter func(statistics.Model) uint16) uint16 {
	return func(c Model, baseGetter func(*Model) uint16, equipGetter func(statistics.Model) uint16) uint16 {
		value := baseGetter(&c)

		//TODO apply MapleWarrior

		equips, err := equipment.GetEquipmentForCharacter(l, db)(c.Id())
		if err != nil {
			l.WithError(err).Errorf("Unable to retrieve equipment for character %d.", c.Id())
		}
		for _, equip := range equips {
			es, err := statistics.GetEquipmentStatistics(l, span)(equip.EquipmentId())
			if err != nil {
				l.WithError(err).Errorf("Unable to retrieve statistics for equipment %d.", equip.EquipmentId())
			} else {
				value += equipGetter(es)
			}
		}
		return value
	}
}

func GainLevel(l logrus.FieldLogger, db *gorm.DB, span opentracing.Span) func(characterId uint32) {
	return func(characterId uint32) {
		characterUpdate(l, db)(characterId, persistLevelUpdate(l, db, span), levelUpdateSuccess(l, span))
	}
}

func persistLevelUpdate(l logrus.FieldLogger, db *gorm.DB, span opentracing.Span) characterFunc {
	return func(c Model) error {
		var modifiers = make([]EntityUpdateFunction, 0)
		modifiers = append(modifiers, SetLevel(c.Level()+1))
		modifiers = append(modifiers, onLevelAdjustAP(c)...)
		modifiers = append(modifiers, onLevelAdjustSP(c)...)
		modifiers = append(modifiers, onLevelAdjustHealthAndMana(l, db, span)(c)...)
		return characterDatabaseUpdate(l, db)(modifiers...)(c)
	}
}

func onLevelAdjustAP(c Model) []EntityUpdateFunction {
	var modifiers = make([]EntityUpdateFunction, 0)
	autoAssignStarterAp := configuration.Get().UseAutoAssignStartersAp
	if autoAssignStarterAp && c.IsBeginner() && c.Level() <= 10 {
		if c.Level() <= 5 {
			modifiers = append(modifiers, SetStrength(5))
		} else {
			modifiers = append(modifiers, SetStrength(4))
			modifiers = append(modifiers, SetDexterity(1))
		}
	} else {
		ap := c.AP() + 5
		if c.Cygnus() && c.Level() > 10 {
			if c.Level() <= 17 {
				ap += 2
			} else if c.Level() < 77 {
				ap += 1
			}
		}
		modifiers = append(modifiers, SetAP(ap))
	}
	return modifiers
}

func onLevelAdjustSP(c Model) []EntityUpdateFunction {
	var modifiers = make([]EntityUpdateFunction, 0)
	if c.IsBeginner() {
		return modifiers
	}

	//TODO account for Evan like SP.
	modifiers = append(modifiers, SetSP(c.SP(0)+3, 0))
	return modifiers
}

func randRange(lowerBound uint16, upperBound uint16) uint16 {
	return uint16(rand.Int31n(int32(upperBound-lowerBound))) + lowerBound
}

func onLevelAdjustHealthAndMana(l logrus.FieldLogger, db *gorm.DB, span opentracing.Span) func(c Model) []EntityUpdateFunction {
	return func(c Model) []EntityUpdateFunction {
		var modifiers = make([]EntityUpdateFunction, 0)
		hp := c.MaxHP()
		mp := c.MaxMP()
		l.Debugf("Adjusting HP/MP for character %d on level. Current HP/MP is %d/%d.", c.Id(), hp, mp)
		if c.IsBeginner() {
			hp += randRange(12, 16)
			mp += randRange(10, 12)
		} else if job.IsA(c.JobId(), job.Warrior, job.DawnWarrior1) {
			hp += getHPIncreaseEnhancementForLevelUp(l, db, span)(c)
			hp += randRange(24, 28)
			mp += randRange(4, 6)
		} else if job.IsA(c.JobId(), job.Magician, job.BlazeWizard1) {
			mp += getMPIncreaseEnhancementForLevelUp(l, db, span)(c)
			hp += randRange(10, 14)
			mp += randRange(22, 24)
		} else if job.IsA(c.JobId(), job.Bowman, job.WindArcher1, job.Thief, job.NightWalker1) {
			hp += randRange(20, 24)
			mp += randRange(14, 16)
		} else if job.IsA(c.JobId(), job.GM) {
			hp += 30000
			mp += 30000
		} else if job.IsA(c.JobId(), job.Pirate, job.ThunderBreaker1) {
			hp += getHPIncreaseEnhancementForLevelUp(l, db, span)(c)
			hp += randRange(22, 28)
			mp += randRange(18, 23)
		} else if job.IsA(c.JobId(), job.Aran1) {
			mpSeed := rand.Int31n(8-4) + 4
			hp += randRange(44, 48)
			mp += uint16(mpSeed) + uint16(math.Floor(float64(mpSeed)*0.1))
		}

		if configuration.Get().UseRandomizeHpMpGain {
			if job.GetJobStyle(c.JobId(), c.Strength(), c.Dexterity()) == job.Magician {
				mp += TotalIntelligence(l, db, span)(c) / 20
			} else {
				mp += TotalIntelligence(l, db, span)(c) / 10
			}
		}
		l.Debugf("HP/MP for character %d on level will become %d/%d.", c.Id(), hp, mp)
		modifiers = append(modifiers, SetHealth(hp), SetMana(mp), SetMaxHP(hp), SetMaxMP(mp))
		return modifiers
	}
}

func levelUpdateSuccess(l logrus.FieldLogger, span opentracing.Span) characterFunc {
	return statisticsUpdateSuccess(l, span)("EXPERIENCE", "LEVEL", "AVAILABLE_AP", "AVAILABLE_SP", "HP", "MP", "MAX_HP", "MAX_MP", "STRENGTH", "DEXTERITY", "LUCK", "INTELLIGENCE")
}

func AssignSP(l logrus.FieldLogger, db *gorm.DB, span opentracing.Span) func(characterId uint32, skillId uint32) {
	return func(characterId uint32, skillId uint32) {
		characterUpdate(l, db)(characterId, assignSP(l, db, span)(skillId))
	}
}

func assignSP(l logrus.FieldLogger, db *gorm.DB, span opentracing.Span) func(skillId uint32) characterFunc {
	return func(skillId uint32) characterFunc {
		return func(c Model) error {
			if s, err := skill.GetSkill(l, db)(c.Id(), skillId); err == nil {
				skillBookId := skill.GetSkillBook(skillId / 10000)
				remainingSP := c.SP(int(skillBookId))

				beginnerSkill := false
				if skillId%10000000 > 999 && skillId%10000000 < 1003 {
					total := uint32(0)
					for i := uint32(0); i < 3; i++ {
						if bs, err := skill.GetSkill(l, db)(c.Id(), uint32(c.JobType())*10000000+1000+i); err == nil {
							total += bs.Level()
						}
					}
					remainingSP = uint32(math.Min(float64(c.Level()-1), 6)) - total
					beginnerSkill = true
					l.Debugf("Skill %d was identified as a beginner skill.", skillId)
				}

				skillMaxLevel := uint32(20)
				if si, err := information.GetById(l, span)(skillId); err == nil {
					skillMaxLevel = uint32(len(si.Effects()))
				}
				var maxLevel = uint32(0)
				if skill.IsFourthJob(c.JobId(), skillId) {
					maxLevel = s.MasterLevel()
				} else {
					maxLevel = skillMaxLevel
				}

				if remainingSP <= 0 {
					l.Debugf("Skill %d update for character %d skipped. Needs more SP.")
					return nil
				}

				if s.Level()+1 > maxLevel {
					l.Debugf("Skill %d update for character %d skipped. Increasing level would push above max level.")
					return nil
				}

				if !beginnerSkill {
					err := adjustSP(l, db, span)(c, -1, skillBookId)
					if err != nil {
						return err
					}
				} else {
					emitEnableActionsCommand(l, span)(c.Id())
				}

				//TODO special handling for aran full swing and over swing.
				err := skill.UpdateSkill(l, db)(c.Id(), skillId, s.Level()+1, s.MasterLevel(), s.Expiration())
				if err != nil {
					return err
				}
				emitSkillUpdateEvent(l, span)(c.Id(), skillId, s.Level()+1, s.MasterLevel(), s.Expiration())
			} else {
				l.Warnf("Received a skill %d assignment for character %d who does not have the skill.", skillId, c.Id())
			}
			return nil
		}
	}
}

func adjustSP(l logrus.FieldLogger, db *gorm.DB, span opentracing.Span) func(c Model, amount int32, bookId uint32) error {
	return func(c Model, amount int32, bookId uint32) error {
		nv := uint32(math.Max(0, float64(int32(c.SP(int(bookId)))+amount)))
		err := characterDatabaseUpdate(l, db)(SetSP(nv, bookId))(c)
		if err != nil {
			return err
		}
		return statisticUpdateSuccess(l, span)("AVAILABLE_SP")(c)
	}
}

func UpdateLoginPosition(l logrus.FieldLogger, db *gorm.DB, span opentracing.Span) func(characterId uint32) {
	return func(characterId uint32) {
		characterUpdate(l, db)(characterId, updateTemporalPositionLogin(l, span))
	}
}

func updateTemporalPositionLogin(l logrus.FieldLogger, span opentracing.Span) characterFunc {
	return func(c Model) error {
		port, err := portal.GetById(l, span)(c.MapId(), c.SpawnPoint())
		if err != nil {
			l.Warnf("Unable to find spawn point %d in map %d for character %d.", c.SpawnPoint(), c.MapId(), c.Id())
			port, err = portal.GetById(l, span)(c.MapId(), 0)
			if err != nil {
				l.Errorf("Unable to get a portal in map %d to update character %d position to.", c.MapId(), c.Id())
				return err
			}
		}
		GetTemporalRegistry().UpdatePosition(c.Id(), port.X(), port.Y())
		return nil
	}
}

func GetForAccountInWorld(_ logrus.FieldLogger, db *gorm.DB) func(accountId uint32, worldId byte) ([]Model, error) {
	return func(accountId uint32, worldId byte) ([]Model, error) {
		return database.ModelSliceProvider[Model, entity](db)(getForAccountInWorld(accountId, worldId), makeCharacter)()
	}
}

func GetForMapInWorld(_ logrus.FieldLogger, db *gorm.DB) func(worldId byte, mapId uint32) ([]Model, error) {
	return func(worldId byte, mapId uint32) ([]Model, error) {
		return database.ModelSliceProvider[Model, entity](db)(getForMapInWorld(worldId, mapId), makeCharacter)()
	}
}

func GetForName(_ logrus.FieldLogger, db *gorm.DB) func(name string) ([]Model, error) {
	return func(name string) ([]Model, error) {
		return database.ModelSliceProvider[Model, entity](db)(getForName(name), makeCharacter)()
	}
}

func GetMaximumBaseDamage(l logrus.FieldLogger, db *gorm.DB, span opentracing.Span) func(characterId uint32) uint32 {
	return func(characterId uint32) uint32 {
		c, err := GetById(l, db)(characterId)
		if err != nil {
			l.WithError(err).Errorf("Unable to retrieve character %d for damage information request.", c.Id())
			return 0
		}

		wa := WeaponAttack(l, db, span)(c)

		equip, err := equipment.GetEquippedItemForCharacterBySlot(l, db)(c.Id(), -11)
		if err != nil {
			l.WithError(err).Errorf("Retrieving equipment for character %d.", c.Id())
			return getMaximumBaseDamageNoWeapon(l, db, span)(c)
		}
		es, err := statistics.GetEquipmentStatistics(l, span)(equip.EquipmentId())
		if err != nil {
			l.WithError(err).Errorf("Retrieving equipment %d statistics for character %d.", equip.EquipmentId(), c.Id())
			return getMaximumBaseDamageNoWeapon(l, db, span)(c)
		}
		return getMaximumBaseDamage(l, db, span)(c, wa, item.GetWeaponType(es.ItemId()))
	}
}

func WeaponAttack(l logrus.FieldLogger, db *gorm.DB, span opentracing.Span) func(c Model) uint16 {
	return func(c Model) uint16 {
		wa := uint16(0)

		equips, err := equipment.GetEquipmentForCharacter(l, db)(c.Id())
		if err != nil {
			l.WithError(err).Errorf("Retrieving equipment for character %d.", c.Id())
			return 0
		}
		for _, equip := range equips {
			es, err := statistics.GetEquipmentStatistics(l, span)(equip.EquipmentId())
			if err != nil {
				l.WithError(err).Errorf("Retrieving equipment %d statistics for character %d.", equip.EquipmentId(), c.Id())
			} else {
				wa += es.WeaponAttack()
			}
		}

		//TODO
		// apply Aran Combo
		// apply ThunderBreaker Marauder energy charge
		// apply Marksman Boost or Bowmaster Expert
		// apply weapon attack buffs
		// apply blessing
		// apply active projectile
		return wa
	}
}

func getMaximumBaseDamage(l logrus.FieldLogger, db *gorm.DB, span opentracing.Span) func(c Model, weaponAttack uint16, weaponType int) uint32 {
	return func(c Model, weaponAttack uint16, weaponType int) uint32 {
		workingWeaponType := weaponType

		if job.IsA(c.JobId(), job.Thief) && workingWeaponType == item.WeaponTypeDaggerOther {
			workingWeaponType = item.WeaponTypeDaggerThieves
		}

		var mainStat uint16
		var secondaryStat uint16

		if workingWeaponType == item.WeaponTypeBow || workingWeaponType == item.WeaponTypeCrossbow || workingWeaponType == item.WeaponTypeGun {
			mainStat = TotalDexterity(l, db, span)(c)
			secondaryStat = TotalStrength(l, db, span)(c)
		} else if workingWeaponType == item.WeaponTypeClaw || workingWeaponType == item.WeaponTypeDaggerThieves {
			mainStat = TotalLuck(l, db, span)(c)
			secondaryStat = TotalDexterity(l, db, span)(c) + TotalStrength(l, db, span)(c)
		} else {
			mainStat = TotalStrength(l, db, span)(c)
			secondaryStat = TotalDexterity(l, db, span)(c)
		}

		return uint32(math.Ceil(((item.GetWeaponDamageMultiplier(workingWeaponType) * float64(mainStat*secondaryStat)) / 100.0) * float64(weaponAttack)))
	}
}

func getMaximumBaseDamageNoWeapon(l logrus.FieldLogger, db *gorm.DB, span opentracing.Span) func(c Model) uint32 {
	return func(c Model) uint32 {
		if job.IsA(c.JobId(), job.Pirate, job.ThunderBreaker1) {
			wm := 3.0
			if c.JobId()%100 != 0 {
				wm = 4.2
			}
			attack := uint32(math.Min(math.Floor(float64(2*c.Level()+31)/3.0), 31))
			strength := TotalStrength(l, db, span)(c)
			dexterity := TotalDexterity(l, db, span)(c)
			return uint32(math.Ceil((float64(strength) * wm * float64(dexterity)) * float64(attack) / 100.0))
		}
		return 1
	}
}

func Create(l logrus.FieldLogger, db *gorm.DB, span opentracing.Span) func(b *Builder) (Model, error) {
	return func(b *Builder) (Model, error) {
		c := b.Build()
		c, err := create(db, c.AccountId(), c.WorldId(), c.Name(), c.Level(), c.Strength(), c.Dexterity(), c.Intelligence(), c.Luck(), c.MaxHP(), c.MaxMP(), c.JobId(), c.Gender(), c.Hair(), c.Face(), c.SkinColor(), c.MapId())
		if err != nil {
			return Model{}, err
		}

		err = inventory.CreateInitialInventories(l, db)(c.Id())
		if err != nil {
			return Model{}, err
		}

		emitCreatedEvent(l, span)(c.Id(), c.WorldId(), c.Name())
		return c, nil
	}
}

func AdjustJob(l logrus.FieldLogger, db *gorm.DB, span opentracing.Span) func(characterId uint32, jobId uint16) error {
	return func(characterId uint32, jobId uint16) error {
		characterUpdate(l, db)(characterId, adjustJob(l, db)(jobId), awardSkillsForJobUpdate(l, db, span)(jobId), jobUpdateSuccess(l, span))
		return nil
	}
}

func awardSkillsForJobUpdate(l logrus.FieldLogger, db *gorm.DB, span opentracing.Span) func(jobId uint16) characterFunc {
	return func(jobId uint16) characterFunc {
		return func(c Model) error {
			skills := make([]uint32, 0)
			switch jobId {
			case job.Warrior:
				skills = []uint32{skill.WarriorImprovedHPRecovery, skill.WarriorImprovedHPIncrease, skill.WarriorEndure, skill.WarriorIronBody, skill.WarriorPowerStrike, skill.WarriorSlashBlast}
			case job.Fighter:
				skills = []uint32{skill.FighterSwordMastery, skill.FighterAxeMastery, skill.FighterFinalAttackSword, skill.FighterFinalAttackAxe, skill.FighterSwordBooster, skill.FighterAxeBooster, skill.FighterRage, skill.FighterPowerGuard}
			case job.Page:
				skills = []uint32{skill.PageSwordMastery, skill.PageBluntWeaponMastery, skill.PageFinalAttackSword, skill.PageFinalAttackBluntWeapon, skill.PageSwordBooster, skill.PageBluntWeaponBooster, skill.PageThreaten, skill.PagePowerGuard}
			case job.Spearman:
				skills = []uint32{skill.SpearmanSpearMastery, skill.SpearmanPolearmMastery, skill.SpearmanFinalAttackSpear, skill.SpearmanFinalAttackPolearm, skill.SpearmanSpearBooster, skill.SpearmanPolearmBooster, skill.SpearmanIronWill, skill.SpearmanHyperBody}
			case job.Magician:
				skills = []uint32{skill.MagicianImprovedMPRecovery, skill.MagicianImprovedMPIncrease, skill.MagicianMagicGuard, skill.MagicianMagicArmor, skill.MagicianEnergyBolt, skill.MagicianMagicClaw}
			case job.FirePoisonWizard:
				skills = []uint32{skill.FirePoisonWizardMPEater, skill.FirePoisonWizardMeditation, skill.FirePoisonWizardTeleport, skill.FirePoisonWizardSlow, skill.FirePoisonWizardFireArrow, skill.FirePoisonWizardPoisonBreath}
			case job.IceLighteningWizard:
				skills = []uint32{skill.IceLightningWizardMPEater, skill.IceLightningWizardMeditation, skill.IceLightningWizardTeleport, skill.IceLightningWizardSlow, skill.IceLightningWizardColdBeam, skill.IceLightningWizardThunderBolt}
			case job.Cleric:
				skills = []uint32{skill.ClericMPEater, skill.ClericTeleport, skill.ClericHeal, skill.ClericInvincible, skill.ClericBless, skill.ClericHolyArrow}
			case job.Bowman:
				skills = []uint32{skill.BowmanBlessingOfAmazon, skill.BowmanCriticalShot, skill.BowmanTheEyeOfAmazon, skill.BowmanFocus, skill.BowmanArrowBlow, skill.BowmanDoubleShot}
			case job.Hunter:
				skills = []uint32{skill.HunterBowMastery, skill.HunterFinalAttack, skill.HunterBowBooster, skill.HunterPowerKnockback, skill.HunterSoulArrow, skill.HunterArrowBomb}
			case job.CrossBowman:
				skills = []uint32{skill.CrossbowmanCrossbowMastery, skill.CrossbowmanFinalAttack, skill.CrossbowmanCrossbowBooster, skill.CrossbowmanPowerKnockback, skill.CrossbowmanSoulArrow, skill.CrossbowmanIronArrow}
			case job.Thief:
				skills = []uint32{skill.ThiefNimbleBody, skill.ThiefKeenEyes, skill.ThiefDisorder, skill.ThiefDarkSight, skill.ThiefDoubleStab, skill.ThiefLuckySeven}
			case job.Assassin:
				skills = []uint32{skill.AssassinClawMastery, skill.AssassinCriticalThrow, skill.AssassinEndure, skill.AssassinClawBooster, skill.AssassinHaste, skill.AssassinDrain}
			case job.Bandit:
				skills = []uint32{skill.BanditDaggerMastery, skill.BanditEndure, skill.BanditDaggerBooster, skill.BanditHaste, skill.BanditSteal, skill.BanditSavageBlow}
			case job.Pirate:
				skills = []uint32{skill.PirateBulletTime, skill.PirateFlashFist, skill.PirateSomersaultKick, skill.PirateDoubleShot, skill.PirateDash}
			case job.Brawler:
				skills = []uint32{skill.BrawlerImproveMaxHP, skill.BrawlerKnucklerMastery, skill.BrawlerBackSpinBlow, skill.BrawlerDoubleUppercut, skill.BrawlerCorkscrewBlow, skill.BrawlerMPRecovery, skill.BrawlerKnucklerBooster, skill.BrawlerOakBarrel}
			case job.Gunslinger:
				skills = []uint32{skill.GunslingerGunMastery, skill.GunslingerInvisibleShot, skill.GunslingerGrenade, skill.GunslingerGunBooster, skill.GunslingerBlankShot, skill.GunslingerWings, skill.GunslingerRecoilShot}
			}

			err := skill.AwardSkills(l, db, span)(c.Id(), skills...)
			if err != nil {
				l.WithError(err).Errorf("Unable to award skills to character %d for job advancement to %d.", c.Id(), jobId)
			}
			return err
		}
	}
}

func jobUpdateSuccess(l logrus.FieldLogger, span opentracing.Span) characterFunc {
	return statisticsUpdateSuccess(l, span)("HP", "MP", "MAX_HP", "MAX_MP", "AVAILABLE_AP", "AVAILABLE_SP", "JOB")
}

func adjustJob(l logrus.FieldLogger, db *gorm.DB) func(jobId uint16) characterFunc {
	return func(jobId uint16) characterFunc {
		return func(c Model) error {
			hp := uint16(0)
			mp := uint16(0)

			j := c.JobId() % 1000
			if j == 100 {
				// 1st job warrior
				hp += randRange(200, 250)
			} else if j == 200 {
				// 1st job magician
				mp += randRange(100, 150)
			} else if j%100 == 0 {
				hp += randRange(100, 150)
				mp += randRange(25, 50)
			} else if j > 0 && j < 200 {
				// 2nd-4th warrior
				hp += randRange(300, 350)
			} else if j < 300 {
				mp += randRange(450, 500)
			} else {
				hp += randRange(300, 350)
				mp += randRange(150, 200)
			}

			modifiers := []EntityUpdateFunction{
				SetHealth(c.HP() + hp),
				SetMana(c.MP() + mp),
				SetMaxHP(c.MaxHP() + hp),
				SetMaxMP(c.MaxMP() + mp),
				SetJob(jobId),
			}
			return characterDatabaseUpdate(l, db)(modifiers...)(c)
		}
	}
}

func ResetAP(l logrus.FieldLogger, db *gorm.DB, span opentracing.Span) func(characterId uint32) error {
	return func(characterId uint32) error {
		characterUpdate(l, db)(characterId, resetAP(l, db), apResetSuccess(l, span))
		return nil
	}
}

func apResetSuccess(l logrus.FieldLogger, span opentracing.Span) characterFunc {
	return statisticsUpdateSuccess(l, span)("AVAILABLE_AP", "AVAILABLE_SP", "STRENGTH", "DEXTERITY", "LUCK", "INTELLIGENCE")
}

func resetAP(l logrus.FieldLogger, db *gorm.DB) characterFunc {
	return func(c Model) error {
		tap := c.AP() + c.Strength() + c.Dexterity() + c.Intelligence() + c.Luck()
		tstr := uint16(4)
		tdex := uint16(4)
		tint := uint16(4)
		tluk := uint16(4)
		tsp := uint32(1)

		switch c.JobId() {
		case job.Warrior, job.DawnWarrior1, job.Aran1:
			tstr = 35
			tsp += uint32((c.Level() - 10) * 3)
			break
		case job.Magician, job.BlazeWizard1:
			tint = 20
			tsp += uint32((c.Level() - 10) * 3)
			break
		case job.Bowman, job.WindArcher1, job.Thief, job.NightWalker1:
			tdex = 25
			tsp += uint32((c.Level() - 10) * 3)
		case job.Pirate, job.ThunderBreaker1:
			tdex = 20
			tsp += uint32((c.Level() - 10) * 3)
		}

		tap -= tstr
		tap -= tdex
		tap -= tint
		tap -= tluk

		if tap < 0 {
			l.Errorf("Cannot reset statistics, character does not have base AP needed.")
			return errors.New("not enough ap")
		}

		modifiers := []EntityUpdateFunction{SetStrength(tstr), SetDexterity(tdex), SetIntelligence(tint), SetLuck(tluk), SetAP(tap), SetSP(tsp, 0)}
		return characterDatabaseUpdate(l, db)(modifiers...)(c)
	}
}

func MoveItem(l logrus.FieldLogger, db *gorm.DB, span opentracing.Span) func(characterId uint32, inventoryType int8, source int16, destination int16) error {
	return func(characterId uint32, inventoryType int8, source int16, destination int16) error {
		_, err := GetById(l, db)(characterId)
		if err != nil {
			l.WithError(err).Errorf("Cannot retrieve character %d performing the drop.", characterId)
			return err
		}
		return inventory.Move(l, db, span)(characterId, inventoryType, source, destination)
	}
}

func DropItem(l logrus.FieldLogger, db *gorm.DB, span opentracing.Span) func(worldId byte, channelId byte, characterId uint32, inventoryType int8, slot int16, quantity int16) error {
	return func(worldId byte, channelId byte, characterId uint32, inventoryType int8, slot int16, quantity int16) error {
		c, err := GetById(l, db)(characterId)
		if err != nil {
			l.WithError(err).Errorf("Cannot retrieve character %d performing the drop.", characterId)
			return err
		}
		ctd := GetTemporalRegistry().GetById(characterId)
		if ctd == nil {
			return errors.New("unable to locate character temporal data")
		}
		return inventory.DropFromInventorySlot(l, db, span)(worldId, channelId, c.MapId(), c.Id(), ctd.X(), ctd.Y(), inventoryType, slot, quantity)
	}
}
