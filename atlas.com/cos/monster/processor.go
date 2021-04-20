package monster

import (
	"atlas-cos/character"
	"atlas-cos/configuration"
	"atlas-cos/kafka/producers"
	"atlas-cos/rest/attributes"
	"atlas-cos/rest/requests"
	"context"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"math"
)

type processor struct {
	l  log.FieldLogger
	db *gorm.DB
}

var Processor = func(l log.FieldLogger, db *gorm.DB) *processor {
	return &processor{l, db}
}

func (p processor) GetMonster(monsterId uint32) (*Model, bool) {
	resp, err := requests.Monster().GetById(monsterId)
	if err != nil {
		p.l.WithError(err).Errorf("Retrieving monster %d information.", monsterId)
		return nil, false
	}
	return makeMonster(resp), true
}

func makeMonster(resp *attributes.MonsterDataContainer) *Model {
	return &Model{
		experience: resp.Data.Attributes.Experience,
		hp:         resp.Data.Attributes.HP,
	}
}

func (p processor) DistributeExperience(worldId byte, channelId byte, mapId uint32, m *Model, entries []*DamageEntry) {
	d := p.produceDistribution(mapId, m, entries)
	for k, v := range d.Solo() {
		experience := float64(v) * d.ExperiencePerDamage()
		c, err := character.Processor(p.l, p.db).GetById(k)
		if err != nil {
			p.l.WithError(err).Errorf("Unable to locate character %d whose for distributing experience from monster death.", k)
		} else {
			whiteExperienceGain := p.isWhiteExperienceGain(c.Id(), d.PersonalRatio(), d.StandardDeviationRatio())
			p.distributeCharacterExperience(c.Id(), c.Level(), experience, 0.0, c.Level(), true, whiteExperienceGain, false)
		}
	}
}

type distribution struct {
	solo                   map[uint32]uint64
	party                  map[uint32]map[uint32]uint64
	personalRatio          map[uint32]float64
	experiencePerDamage    float64
	standardDeviationRatio float64
}

func (d distribution) Solo() map[uint32]uint64 {
	return d.solo
}

func (d distribution) ExperiencePerDamage() float64 {
	return d.experiencePerDamage
}

func (d distribution) PersonalRatio() map[uint32]float64 {
	return d.personalRatio
}

func (d distribution) StandardDeviationRatio() float64 {
	return d.standardDeviationRatio
}

func (p processor) produceDistribution(mapId uint32, monster *Model, entries []*DamageEntry) distribution {
	totalEntries := 0
	//TODO incorporate party distribution.
	partyDistribution := make(map[uint32]map[uint32]uint64)
	soloDistribution := make(map[uint32]uint64)

	for _, entry := range entries {
		if character.Processor(p.l, p.db).InMap(entry.CharacterId(), mapId) {
			soloDistribution[entry.CharacterId()] = entry.Damage()
		}
		totalEntries += 1
	}

	//TODO account for healing
	totalDamage := monster.HP()
	epd := float64(monster.Experience()) / float64(totalDamage)

	personalRatio := make(map[uint32]float64)
	entryExperienceRatio := make([]float64, 0)

	for k, v := range soloDistribution {
		ratio := float64(v) / float64(totalDamage)
		personalRatio[k] = ratio
		entryExperienceRatio = append(entryExperienceRatio, ratio)
	}

	for _, party := range partyDistribution {
		ratio := 0.0
		for k, v := range party {
			cr := float64(v) / float64(totalDamage)
			personalRatio[k] = cr
			ratio += cr
		}
		entryExperienceRatio = append(entryExperienceRatio, ratio)
	}

	stdr := p.calculateExperienceStandardDeviationThreshold(entryExperienceRatio, totalEntries)
	return distribution{
		solo:                   soloDistribution,
		party:                  partyDistribution,
		personalRatio:          personalRatio,
		experiencePerDamage:    epd,
		standardDeviationRatio: stdr,
	}
}

func (p processor) calculateExperienceStandardDeviationThreshold(entryExperienceRatio []float64, totalEntries int) float64 {
	averageExperienceReward := 0.0
	for _, v := range entryExperienceRatio {
		averageExperienceReward += v
	}
	averageExperienceReward /= float64(totalEntries)

	varExperienceReward := 0.0
	for _, v := range entryExperienceRatio {
		varExperienceReward += math.Pow(v-averageExperienceReward, 2)
	}
	varExperienceReward /= float64(len(entryExperienceRatio))

	return averageExperienceReward + math.Sqrt(varExperienceReward)
}

func (p processor) isWhiteExperienceGain(characterId uint32, personalRatio map[uint32]float64, standardDeviationRatio float64) bool {
	if val, ok := personalRatio[characterId]; ok {
		return val >= standardDeviationRatio
	} else {
		return false
	}
}

func (p processor) distributeCharacterExperience(characterId uint32, level byte, experience float64, partyBonusMod float64, totalPartyLevel byte, hightestPartyDamage bool, whiteExperienceGain bool, hasPartySharers bool) {
	expSplitCommonMod := configuration.Get().ExpSplitCommonMod
	characterExperience := (float64(expSplitCommonMod) * float64(level)) / float64(totalPartyLevel)
	if hightestPartyDamage {
		characterExperience += float64(configuration.Get().ExpSplitMvpMod)
	}
	characterExperience *= experience
	bonusExperience := partyBonusMod * characterExperience

	p.giveExperienceToCharacter(characterId, characterExperience, bonusExperience, whiteExperienceGain, hasPartySharers)
}

func (p processor) giveExperienceToCharacter(characterId uint32, experience float64, bonus float64, whiteExperienceGain bool, hasPartySharers bool) {
	correctedPersonal := p.experienceValueToInteger(experience)
	correctedParty := p.experienceValueToInteger(bonus)
	producers.CharacterExperienceGain(p.l, context.Background()).Emit(characterId, correctedPersonal, correctedParty, true, false, whiteExperienceGain)
}

func (p processor) experienceValueToInteger(experience float64) uint32 {
	if experience > math.MaxInt32 {
		experience = math.MaxInt32
	} else if experience < math.MinInt32 {
		experience = math.MinInt32
	}
	return uint32(math.RoundToEven(experience))
}
