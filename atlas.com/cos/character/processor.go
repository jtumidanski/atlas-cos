package character

import (
	"atlas-cos/configuration"
	"atlas-cos/equipment"
	"atlas-cos/equipment/statistics"
	"atlas-cos/item"
	"atlas-cos/job"
	"atlas-cos/kafka/producers"
	_map "atlas-cos/map"
	"atlas-cos/portal"
	"atlas-cos/skill"
	"atlas-cos/skill/information"
	"context"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"math"
	"math/rand"
)

type processor struct {
	l  log.FieldLogger
	db *gorm.DB
}

var Processor = func(l log.FieldLogger, db *gorm.DB) *processor {
	return &processor{l, db}
}

// characterFunc - Function which does something about the character, and returns whether or not further processing should continue.
type characterFunc func(*Model) error

// Returns a function which accepts a character model,and updates the persisted state of the character given a set of
// modifying functions.
func (p *processor) characterDatabaseUpdate(modifiers ...EntityUpdateFunction) characterFunc {
	return func(c *Model) error {
		if len(modifiers) > 0 {
			err := Update(p.db, c.Id(), modifiers...)
			if err != nil {
				return err
			}
		}
		return nil
	}
}

// For the characterId, perform the updaterFunc, and if successful, call the successFunc, otherwise log an error.
func (p *processor) characterUpdate(characterId uint32, functions ...characterFunc) {
	c, err := GetById(p.db, characterId)
	if err != nil {
		return
	}

	err = nil
	for _, f := range functions {
		err = f(c)
		if err != nil {
			break
		}
	}
}

// AdjustHealth - Adjusts the Health statistic for a character, and emits a CharacterStatUpdateEvent when successful.
func (p *processor) AdjustHealth(characterId uint32, amount uint16) {
	p.characterUpdate(characterId, p.persistHealthUpdate(amount), p.healthUpdateSuccess())
}

// Produces a function which persists a character health update, given the amount, respecting the MaxHP bound.
func (p *processor) persistHealthUpdate(amount uint16) characterFunc {
	return func(c *Model) error {
		adjustedAmount := p.enforceBounds(amount, c.HP(), c.MaxHP(), 0)
		return p.characterDatabaseUpdate(SetHealth(adjustedAmount))(c)
	}
}

// When a Health update is successful, emit a CharacterStatUpdateEvent.
func (p *processor) healthUpdateSuccess() characterFunc {
	return p.statisticUpdateSuccess("HP")
}

// AdjustMana - Adjusts the Mana statistic for a character, and emits a CharacterStatUpdateEvent when successful.
func (p *processor) AdjustMana(characterId uint32, amount uint16) {
	p.characterUpdate(characterId, p.persistManaUpdate(amount), p.manaUpdateSuccess())
}

// Produces a function which persists a character mana update, given the amount, respecting the MaxMP bound.
func (p *processor) persistManaUpdate(amount uint16) characterFunc {
	return func(c *Model) error {
		adjustedAmount := p.enforceBounds(amount, c.MP(), c.MaxMP(), 0)
		return p.characterDatabaseUpdate(SetMana(adjustedAmount))(c)
	}
}

func (p *processor) enforceBounds(change uint16, current uint16, upperBound uint16, lowerBound uint16) uint16 {
	var adjusted = current + change
	return uint16(math.Min(math.Max(float64(adjusted), float64(lowerBound)), float64(upperBound)))
}

// When a Mana update is successful, emit a CharacterStatUpdateEvent.
func (p *processor) manaUpdateSuccess() characterFunc {
	return p.statisticUpdateSuccess("MP")
}

// Produces a function which emits a CharacterStatUpdateEvent for the given characterId and statistic
func (p *processor) statisticUpdateSuccess(statistic string) characterFunc {
	return func(c *Model) error {
		producers.CharacterStatUpdate(p.l, context.Background()).Emit(c.Id(), []string{statistic})
		return nil
	}
}

// Produces a function which emits a CharacterStatUpdateEvent for the given characterId and statistic
func (p *processor) statisticsUpdateSuccess(statistics ...string) characterFunc {
	return func(c *Model) error {
		producers.CharacterStatUpdate(p.l, context.Background()).Emit(c.Id(), statistics)
		return nil
	}
}

// AdjustMeso - Adjusts the Meso count for a character, and emits a MesoGainedEvent when successful.
func (p *processor) AdjustMeso(characterId uint32, amount uint32, show bool) {
	p.characterUpdate(characterId, p.persistMesoUpdate(amount), p.mesoUpdateSuccess(amount, show))
}

// Produces a function which persists a character meso update, given the amount.
func (p *processor) persistMesoUpdate(amount uint32) characterFunc {
	return func(c *Model) error {
		return p.characterDatabaseUpdate(IncreaseMeso(amount))(c)
	}
}

// Produces a function which emits a MesoGainedEvent for the given characterId and amount.
func (p *processor) mesoUpdateSuccess(amount uint32, show bool) characterFunc {
	return func(c *Model) error {
		if show {
			producers.MesoGained(p.l, context.Background()).Emit(c.Id(), amount)
		}
		return nil
	}
}

// ChangeMap - Changes the map for a character in the database, updates the temporal position, and emits a MapChangedEvent when successful.
func (p *processor) ChangeMap(characterId uint32, worldId byte, channelId byte, mapId uint32, portalId uint32) {
	p.characterUpdate(characterId, p.performChangeMap(mapId, portalId), p.changeMapSuccess(worldId, channelId, mapId, portalId))
}

// Produces a function which persists a character map update, then updates the temporal position.
func (p *processor) performChangeMap(mapId uint32, portalId uint32) characterFunc {
	return func(c *Model) error {
		err := p.characterDatabaseUpdate(SetMapId(mapId))(c)
		if err != nil {
			return err
		}
		por, err := portal.Processor(p.l).GetMapPortalById(mapId, portalId)
		if err != nil {
			return err
		}
		GetTemporalRegistry().UpdatePosition(c.Id(), por.X(), por.Y())
		return nil
	}
}

// Produces a function which emits a MapChangedEvent for the given characterId.
func (p *processor) changeMapSuccess(worldId byte, channelId byte, mapId uint32, portalId uint32) characterFunc {
	return func(c *Model) error {
		producers.MapChanged(p.l, context.Background()).Emit(worldId, channelId, mapId, portalId, c.Id())
		return nil
	}
}

// GainExperience - Updates the character based on the experience gained, may trigger level updates depending on amount gained.
func (p *processor) GainExperience(characterId uint32, amount uint32) {
	p.characterUpdate(characterId, p.performGainExperience(amount))
}

func (p *processor) performGainExperience(amount uint32) characterFunc {
	return func(c *Model) error {
		p.gainExperience(c.Id(), c.Level(), c.MaxClassLevel(), c.Experience(), amount)
		return nil
	}
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

func (p *processor) persistSetExperience(experience uint32) characterFunc {
	return func(c *Model) error {
		return p.characterDatabaseUpdate(SetExperience(experience))(c)
	}
}

func (p *processor) experienceChangeSuccess() characterFunc {
	return p.statisticUpdateSuccess("EXPERIENCE")
}

func (p *processor) setExperience(characterId uint32, experience uint32) {
	p.characterUpdate(characterId, p.persistSetExperience(experience), p.experienceChangeSuccess())
}

func (p *processor) persistIncreaseExperience(experience uint32) characterFunc {
	return func(c *Model) error {
		return p.characterDatabaseUpdate(IncreaseExperience(experience))(c)
	}
}

func (p *processor) increaseExperience(characterId uint32, gain uint32) {
	p.characterUpdate(characterId, p.persistIncreaseExperience(gain), p.experienceChangeSuccess())
}

func (p *processor) MoveCharacter(characterId uint32, x int16, y int16, stance byte) {
	p.characterUpdate(characterId, p.updateTemporalPosition(x, y, stance), p.updateSpawnPoint(x, y))
}

func (p *processor) updateTemporalPosition(x int16, y int16, stance byte) characterFunc {
	return func(c *Model) error {
		GetTemporalRegistry().Update(c.Id(), x, y, stance)
		return nil
	}
}

func (p *processor) updateSpawnPoint(x int16, y int16) characterFunc {
	return func(c *Model) error {
		sp, err := _map.Processor(p.l).FindClosestSpawnPoint(c.MapId(), x, y)
		if err != nil {
			return err
		}
		return p.characterDatabaseUpdate(UpdateSpawnPoint(sp.Id()))(c)
	}
}

func (p *processor) UpdateStance(characterId uint32, stance byte) {
	p.characterUpdate(characterId, p.updateStance(stance))
}

func (p *processor) updateStance(stance byte) characterFunc {
	return func(c *Model) error {
		GetTemporalRegistry().UpdateStance(c.Id(), stance)
		return nil
	}
}

func (p *processor) GetById(characterId uint32) (*Model, error) {
	return GetById(p.db, characterId)
}

func (p *processor) InMap(characterId uint32, mapId uint32) bool {
	c, err := p.GetById(characterId)
	if err != nil {
		p.l.Errorf("Unable to validate character %d is in map %d. Assuming not.", characterId, mapId)
		return false
	}
	return c.MapId() == mapId
}

func (p *processor) outOfRange(new uint16, change uint16) bool {
	return new < 4 && change != 0 || new > configuration.Get().MaxAp
}

func (p *processor) persistAttributeUpdate(getter func(*Model) uint16, modifierGetter func(uint16, uint16) []EntityUpdateFunction) characterFunc {
	return func(c *Model) error {
		value := getter(c) + 1
		if p.outOfRange(value, 1) {
			return nil
		}
		return p.characterDatabaseUpdate(modifierGetter(value, c.AP()-1)...)(c)
	}
}

func (p *processor) AssignStrength(characterId uint32) {
	p.characterUpdate(characterId, p.persistStrengthUpdate(), p.strengthUpdateSuccess())
}

func (p *processor) persistStrengthUpdate() characterFunc {
	return p.persistAttributeUpdate((*Model).Strength, SpendOnStrength)
}

func (p *processor) strengthUpdateSuccess() characterFunc {
	return p.statisticsUpdateSuccess("STRENGTH", "AVAILABLE_AP")
}

func (p *processor) AssignDexterity(characterId uint32) {
	p.characterUpdate(characterId, p.persistDexterityUpdate(), p.dexterityUpdateSuccess())
}

func (p *processor) persistDexterityUpdate() characterFunc {
	return p.persistAttributeUpdate((*Model).Dexterity, SpendOnDexterity)
}

func (p *processor) dexterityUpdateSuccess() characterFunc {
	return p.statisticsUpdateSuccess("DEXTERITY", "AVAILABLE_AP")
}

func (p *processor) AssignIntelligence(characterId uint32) {
	p.characterUpdate(characterId, p.persistIntelligenceUpdate(), p.intelligenceUpdateSuccess())
}

func (p *processor) persistIntelligenceUpdate() characterFunc {
	return p.persistAttributeUpdate((*Model).Intelligence, SpendOnIntelligence)
}

func (p *processor) intelligenceUpdateSuccess() characterFunc {
	return p.statisticsUpdateSuccess("INTELLIGENCE", "AVAILABLE_AP")
}

func (p *processor) AssignLuck(characterId uint32) {
	p.characterUpdate(characterId, p.persistLuckUpdate(), p.luckUpdateSuccess())
}

func (p *processor) persistLuckUpdate() characterFunc {
	return p.persistAttributeUpdate((*Model).Luck, SpendOnLuck)
}

func (p *processor) luckUpdateSuccess() characterFunc {
	return p.statisticsUpdateSuccess("LUCK", "AVAILABLE_AP")
}

func (p *processor) AssignHp(characterId uint32) {
	p.characterUpdate(characterId, p.persistHpUpdate(), p.hpUpdateSuccess())
}

func (p *processor) persistHpUpdate() characterFunc {
	return func(c *Model) error {
		adjustedHP := p.calculateHPChange(c, false)
		return p.characterDatabaseUpdate(SetMaxHP(adjustedHP))(c)
	}
}

func (p *processor) calculateHPChange(c *Model, usedAPReset bool) uint16 {
	var maxHP uint16 = 0
	if job.IsA(c.JobId(), job.Warrior, job.DawnWarrior1) {
		//TODO apply IMPROVED HP INCREASE OR IMPROVED MAX HP
		maxHP = p.adjustHPMPGain(usedAPReset, maxHP, 20, 22, 18, 18, 20)
	} else if job.IsA(c.JobId(), job.Aran1) {
		maxHP = p.adjustHPMPGain(usedAPReset, maxHP, 20, 30, 26, 26, 28)
	} else if job.IsA(c.JobId(), job.Magician, job.BlazeWizard1) {
		maxHP = p.adjustHPMPGain(usedAPReset, maxHP, 6, 9, 5, 5, 6)
	} else if job.IsA(c.JobId(), job.Thief, job.NightWalker1) {
		maxHP = p.adjustHPMPGain(usedAPReset, maxHP, 16, 18, 14, 14, 16)
	} else if job.IsA(c.JobId(), job.Bowman, job.WindArcher1) {
		maxHP = p.adjustHPMPGain(usedAPReset, maxHP, 16, 18, 14, 14, 16)
	} else if job.IsA(c.JobId(), job.Pirate, job.ThunderBreaker1) {
		//TODO apply IMPROVE HP INCREASE OR IMPROVE MAX HP
		maxHP = p.adjustHPMPGain(usedAPReset, maxHP, 18, 20, 16, 16, 18)
	} else {
		maxHP = p.adjustHPMPGain(usedAPReset, maxHP, 8, 12, 8, 8, 10)
	}
	return maxHP
}

func (p *processor) adjustHPMPGain(usedAPReset bool, maxHP uint16, apResetAmount uint16, upperBound uint16, lowerBound uint16, floor uint16, staticAmount uint16) uint16 {
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

func (p *processor) hpUpdateSuccess() characterFunc {
	return p.statisticsUpdateSuccess("MAX_HP")
}

func (p *processor) AssignMp(characterId uint32) {
	p.characterUpdate(characterId, p.persistMpUpdate(), p.mpUpdateSuccess())
}

func (p *processor) persistMpUpdate() characterFunc {
	return func(c *Model) error {
		adjustedMP := p.calculateMPChange(c, false)
		return p.characterDatabaseUpdate(SetMaxMP(adjustedMP))(c)
	}
}

func (p *processor) calculateMPChange(c *Model, usedAPReset bool) uint16 {
	jobId := c.JobId()
	var maxMP uint16 = 0

	if job.IsA(jobId, job.Warrior, job.DawnWarrior1, job.Aran1) {
		maxMP = p.adjustHPMPGain(usedAPReset, maxMP, 2, 4, 2, c.Intelligence()/10, 3)
	} else if job.IsA(jobId, job.Magician, job.BlazeWizard1) {
		//TODO apply IMPROVED MP INCREASE OR IMPROVED MAX MP
		maxMP = p.adjustHPMPGain(usedAPReset, maxMP, 18, 16, 12, c.Intelligence()/20, 18)
	} else if job.IsA(jobId, job.Thief, job.NightWalker1) {
		maxMP = p.adjustHPMPGain(usedAPReset, maxMP, 10, 8, 6, c.Intelligence()/10, 10)
	} else if job.IsA(jobId, job.Bowman, job.WindArcher1) {
		maxMP = p.adjustHPMPGain(usedAPReset, maxMP, 10, 8, 6, c.Intelligence()/10, 10)
	} else if job.IsA(jobId, job.Pirate, job.ThunderBreaker1) {
		maxMP = p.adjustHPMPGain(usedAPReset, maxMP, 14, 9, 7, c.Intelligence()/10, 14)
	} else {
		maxMP = p.adjustHPMPGain(usedAPReset, maxMP, 6, 6, 4, c.Intelligence()/10, 6)
	}
	return maxMP
}

func (p *processor) mpUpdateSuccess() characterFunc {
	return p.statisticsUpdateSuccess("MAX_MP")
}

func (p *processor) TotalIntelligence(c *Model) uint16 {
	return p.totalStat(*c, (*Model).Intelligence, (*statistics.Model).Intelligence)
}

func (p *processor) TotalDexterity(c *Model) uint16 {
	return p.totalStat(*c, (*Model).Dexterity, (*statistics.Model).Dexterity)
}

func (p *processor) TotalStrength(c *Model) uint16 {
	return p.totalStat(*c, (*Model).Strength, (*statistics.Model).Strength)
}

func (p *processor) TotalLuck(c *Model) uint16 {
	return p.totalStat(*c, (*Model).Luck, (*statistics.Model).Luck)
}

func (p *processor) totalStat(c Model, baseGetter func(*Model) uint16, equipGetter func(*statistics.Model) uint16) uint16 {
	value := baseGetter(&c)

	//TODO apply MapleWarrior

	equips, err := equipment.Processor(p.l, p.db).GetEquipmentForCharacter(c.Id())
	if err != nil {
		p.l.Errorf("Unable to retrieve equipment for character.")
	}
	for _, equip := range equips {
		es, err := statistics.Processor(p.l, p.db).GetEquipmentStatistics(equip.EquipmentId())
		if err != nil {
			p.l.Errorf("Unable to retrieve statistics for equipment %d.", equip.EquipmentId())
		} else {
			value += equipGetter(es)
		}
	}
	return value
}

func (p *processor) GainLevel(characterId uint32) {
	p.characterUpdate(characterId, p.persistLevelUpdate(), p.levelUpdateSuccess())
}

func (p *processor) persistLevelUpdate() characterFunc {
	return func(c *Model) error {
		var modifiers = make([]EntityUpdateFunction, 0)
		modifiers = append(modifiers, p.onLevelAdjustAP(c)...)
		modifiers = append(modifiers, p.onLevelAdjustHealthAndMana(c)...)
		return p.characterDatabaseUpdate(modifiers...)(c)
	}
}

func (p *processor) onLevelAdjustAP(c *Model) []EntityUpdateFunction {
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
		modifiers = append(modifiers, SetAP(5))
		if c.Cygnus() && c.Level() > 10 {
			if c.Level() <= 17 {
				modifiers = append(modifiers, IncreaseAP(2))
			} else if c.Level() < 77 {
				modifiers = append(modifiers, IncreaseAP(1))
			}
		}
	}
	return modifiers
}

func (p *processor) onLevelAdjustHealthAndMana(c *Model) []EntityUpdateFunction {
	var modifiers = make([]EntityUpdateFunction, 0)
	if c.IsBeginner() {
		modifiers = append(modifiers, IncreaseHPRange(12, 16))
		modifiers = append(modifiers, IncreaseMPRange(10, 12))
	} else if job.IsA(c.JobId(), job.Warrior, job.DawnWarrior1) {
		//TODO process DawnWarrior.MAX_HP_INCREASE : Warrior.IMPROVED_MAX_HP
		modifiers = append(modifiers, IncreaseHPRange(24, 28))
		modifiers = append(modifiers, IncreaseMPRange(4, 6))
	} else if job.IsA(c.JobId(), job.Magician, job.BlazeWizard1) {
		//TODO process BlazeWizard.INCREASING_MAX_MP : Magician.IMPROVED_MAX_MP_INCREASE
		modifiers = append(modifiers, IncreaseHPRange(10, 14))
		modifiers = append(modifiers, IncreaseMPRange(22, 24))
	} else if job.IsA(c.JobId(), job.Bowman, job.WindArcher1, job.Thief, job.NightWalker1) {
		modifiers = append(modifiers, IncreaseHPRange(20, 24))
		modifiers = append(modifiers, IncreaseMPRange(14, 16))
	} else if job.IsA(c.JobId(), job.GM) {
		modifiers = append(modifiers, IncreaseHP(30000))
		modifiers = append(modifiers, IncreaseMP(30000))
	} else if job.IsA(c.JobId(), job.Pirate, job.ThunderBreaker1) {
		//TODO process ThunderBreaker.IMPROVE_MAX_HP : Brawler.IMPROVE_MAX_HP
		modifiers = append(modifiers, IncreaseHPRange(22, 28))
		modifiers = append(modifiers, IncreaseMPRange(18, 23))
	} else if job.IsA(c.JobId(), job.Aran1) {
		mpSeed := rand.Int31n(8-4) + 4
		modifiers = append(modifiers, IncreaseHPRange(44, 48))
		modifiers = append(modifiers, IncreaseMP(uint16(mpSeed)+uint16(math.Floor(float64(mpSeed)*0.1))))
	}

	if configuration.Get().UseRandomizeHpMpGain {
		if job.GetJobStyle(c.JobId(), c.Strength(), c.Dexterity()) == job.Magician {
			modifiers = append(modifiers, IncreaseMP(p.TotalIntelligence(c)/20))
		} else {
			modifiers = append(modifiers, IncreaseMP(p.TotalIntelligence(c)/10))
		}
	}
	return modifiers
}

func (p *processor) levelUpdateSuccess() characterFunc {
	return p.statisticsUpdateSuccess("EXPERIENCE", "LEVEL", "AVAILABLE_AP", "HP", "MP", "MAX_HP", "MAX_MP", "STRENGTH", "DEXTERITY", "LUCK", "INTELLIGENCE")
}

func (p *processor) AssignSP(characterId uint32, skillId uint32) {
	p.characterUpdate(characterId, p.assignSP(skillId))
}

func (p *processor) assignSP(skillId uint32) characterFunc {
	return func(c *Model) error {
		if s, ok := skill.Processor(p.l, p.db).GetSkill(c.Id(), skillId); ok {
			skillBookId := skill.GetSkillBook(skillId / 10000)
			remainingSP := c.SP(int(skillBookId))

			beginnerSkill := false
			if skillId%10000000 > 999 && skillId%10000000 < 1003 {
				total := uint32(0)
				for i := uint32(0); i < 3; i++ {
					if bs, ok := skill.Processor(p.l, p.db).GetSkill(c.Id(), uint32(c.JobType())*10000000+1000+i); ok {
						total += bs.Level()
					}
				}
				remainingSP = uint32(math.Min(float64(c.Level()-1), 6)) - total
				beginnerSkill = true
			}

			skillMaxLevel := uint32(20)
			if si, ok := information.Processor(p.l, p.db).GetSkillInformation(skillId); ok {
				skillMaxLevel = uint32(len(si.Effects()))
			}
			var maxLevel = uint32(0)
			if skill.IsFourthJob(c.JobId(), skillId) {
				maxLevel = s.MasterLevel()
			} else {
				maxLevel = skillMaxLevel
			}

			if remainingSP > 0 && uint32(c.Level()+1) <= maxLevel {
				if !beginnerSkill {
					err := p.adjustSP(c, -1, skillBookId)
					if err != nil {
						return err
					}
				} else {
					producers.EnableActions(p.l, context.Background()).Emit(c.Id())
				}

				//TODO special handling for aran full swing and over swing.
				err := skill.Processor(p.l, p.db).UpdateSkill(c.Id(), skillId, s.Level()+1, s.MasterLevel(), s.Expiration())
				if err != nil {
					return err
				}
				producers.CharacterSkillUpdate(p.l, context.Background()).Emit(c.Id(), skillId, s.Level()+1, s.MasterLevel(), s.Expiration())
			}
		}
		return nil
	}
}

func (p *processor) adjustSP(c *Model, amount int32, bookId uint32) error {
	nv := uint32(math.Max(0, float64(int32(c.SP(int(bookId)))+amount)))
	err := p.characterDatabaseUpdate(SetSP(nv, bookId))(c)
	if err != nil {
		return err
	}
	return p.statisticUpdateSuccess("AVAILABLE_SP")(c)
}

func (p *processor) UpdateLoginPosition(characterId uint32) {
	p.characterUpdate(characterId, p.updateTemporalPositionLogin())
}

func (p *processor) updateTemporalPositionLogin() characterFunc {
	return func(c *Model) error {
		port, err := portal.Processor(p.l).GetMapPortalById(c.MapId(), c.SpawnPoint())
		if err != nil {
			p.l.Warnf("Unable to find spawn point %d in map %d for character %d.", c.SpawnPoint(), c.MapId(), c.Id())
			port, err = portal.Processor(p.l).GetMapPortalById(c.MapId(), 0)
			if err != nil {
				p.l.Errorf("Unable to get a portal in map %d to update character %d position to.", c.MapId(), c.Id())
				return err
			}
		}
		GetTemporalRegistry().UpdatePosition(c.Id(), port.X(), port.Y())
		return nil
	}
}

func (p *processor) GetForAccountInWorld(accountId uint32, worldId byte) ([]*Model, error) {
	return GetForAccountInWorld(p.db, accountId, worldId)
}

func (p *processor) GetForMapInWorld(worldId byte, mapId uint32) ([]*Model, error) {
	return GetForMapInWorld(p.db, worldId, mapId)
}

func (p *processor) GetForName(name string) ([]*Model, error) {
	return GetForName(p.db, name)
}

func (p *processor) GetMaximumBaseDamage(characterId uint32) uint32 {
	c, err := p.GetById(characterId)
	if err != nil {
		p.l.Errorf("Unable to retrieve character %d for damage information request.")
		return 0
	}

	wa := p.WeaponAttack(c)

	equip, err := equipment.Processor(p.l, p.db).GetEquippedItemForCharacterBySlot(c.Id(), -11)
	if err != nil {
		p.l.Errorf("Retrieving equipment for character %d.", c.Id())
		return p.getMaximumBaseDamageNoWeapon(c)
	}
	es, err := statistics.Processor(p.l, p.db).GetEquipmentStatistics(equip.EquipmentId())
	if err != nil {
		p.l.Errorf("Retrieving equipment %d statistics for character %d.", equip.EquipmentId(), c.Id())
		return p.getMaximumBaseDamageNoWeapon(c)
	}
	return p.getMaximumBaseDamage(c, wa, item.GetWeaponType(es.ItemId()))
}

func (p *processor) WeaponAttack(c *Model) uint16 {
	wa := uint16(0)

	equips, err := equipment.Processor(p.l, p.db).GetEquipmentForCharacter(c.Id())
	if err != nil {
		p.l.Errorf("Retrieving equipment for character %d.", c.Id())
		return 0
	}
	for _, equip := range equips {
		es, err := statistics.Processor(p.l, p.db).GetEquipmentStatistics(equip.EquipmentId())
		if err != nil {
			p.l.Errorf("Retrieving equipment %d statistics for character %d.", equip.EquipmentId(), c.Id())
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

func (p *processor) getMaximumBaseDamage(c *Model, weaponAttack uint16, weaponType int) uint32 {
	workingWeaponType := weaponType

	if job.IsA(c.JobId(), job.Thief) && workingWeaponType == item.WeaponTypeDaggerOther {
		workingWeaponType = item.WeaponTypeDaggerThieves
	}

	var mainStat uint16
	var secondaryStat uint16

	if workingWeaponType == item.WeaponTypeBow || workingWeaponType == item.WeaponTypeCrossbow || workingWeaponType == item.WeaponTypeGun {
		mainStat = p.TotalDexterity(c)
		secondaryStat = p.TotalStrength(c)
	} else if workingWeaponType == item.WeaponTypeClaw || workingWeaponType == item.WeaponTypeDaggerThieves {
		mainStat = p.TotalLuck(c)
		secondaryStat = p.TotalDexterity(c) + p.TotalStrength(c)
	} else {
		mainStat = p.TotalStrength(c)
		secondaryStat = p.TotalDexterity(c)
	}

	return uint32(math.Ceil(((item.GetWeaponDamageMultiplier(workingWeaponType) * float64(mainStat*secondaryStat)) / 100.0) * float64(weaponAttack)))
}

func (p *processor) getMaximumBaseDamageNoWeapon(c *Model) uint32 {
	if job.IsA(c.JobId(), job.Pirate, job.ThunderBreaker1) {
		wm := 3.0
		if c.JobId()%100 != 0 {
			wm = 4.2
		}
		attack := uint32(math.Min(math.Floor(float64(2*c.Level()+31)/3.0), 31))
		strength := p.TotalStrength(c)
		dexterity := p.TotalDexterity(c)
		return uint32(math.Ceil((float64(strength) * wm * float64(dexterity)) * float64(attack) / 100.0))
	}
	return 1
}

func (p *processor) Create(b *Builder) (*Model, error) {
	c := b.Build()
	c, err := Create(p.db, c.AccountId(), c.WorldId(), c.Name(), c.Level(), c.Strength(), c.Dexterity(), c.Intelligence(), c.Luck(), c.MaxHP(), c.MaxMP(), c.JobId(), c.Gender(), c.Hair(), c.Face(), c.SkinColor(), c.MapId())
	if err != nil {
		return nil, err
	}
	producers.CharacterCreated(p.l, context.Background()).Emit(c.Id(), c.WorldId(), c.Name())
	return c, nil
}
