package skill

import "atlas-cos/job"

const (
	BeginnerThreeSnails uint32 = 1000
	BeginnerRecovery    uint32 = 1001
	BeginnerNimbleFeet  uint32 = 1002

	WarriorImprovedHPRecovery uint32 = 1000000
	WarriorImprovedHPIncrease uint32 = 1000001
	WarriorEndure             uint32 = 1000002
	WarriorIronBody           uint32 = 1001003
	WarriorPowerStrike        uint32 = 1001004
	WarriorSlashBlast         uint32 = 1001005

	MagicianImprovedMPRecovery uint32 = 2000000
	MagicianImprovedMPIncrease uint32 = 2000001
	MagicianMagicGuard         uint32 = 2001002
	MagicianMagicArmor         uint32 = 2001003
	MagicianEnergyBolt         uint32 = 2001004
	MagicianMagicClaw          uint32 = 2001005

	BowmanBlessingOfAmazon uint32 = 3000000
	BowmanCriticalShot     uint32 = 3000001
	BowmanTheEyeOfAmazon   uint32 = 3000002
	BowmanFocus            uint32 = 3001003
	BowmanArrowBlow        uint32 = 3001004
	BowmanDoubleShot       uint32 = 3001005

	ThiefNimbleBody uint32 = 4000000
	ThiefKeenEyes   uint32 = 4000001
	ThiefDisorder   uint32 = 4001002
	ThiefDarkSight  uint32 = 4001003
	ThiefDoubleStab uint32 = 4001334
	ThiefLuckySeven uint32 = 4001344

	PirateBulletTime     uint32 = 5000000
	PirateFlashFist      uint32 = 5001001
	PirateSomersaultKick uint32 = 5001002
	PirateDoubleShot     uint32 = 5001003
	PirateDash           uint32 = 5001005

	BrawlerImproveMaxHP uint32 = 5100000

	NoblesseThreeSnails uint32 = 10001000
	NoblesseRecovery    uint32 = 10001001
	NoblesseNimbleFeet  uint32 = 10001002

	DawnWarriorMaxHPEnhancement uint32 = 11000000

	BlazeWizardIncreasingMaxMP uint32 = 12000000

	ThunderBreakerImproveMaxHP uint32 = 15100000

	EvanMagicMastery uint32 = 22170001
	EvanFlameWheel   uint32 = 22171003
	EvanHerosWill    uint32 = 22171004
	EvanDarkFog      uint32 = 22181002
	EvanSoulStone    uint32 = 22181003
)

// Is determines if the reference skill matches with any of the choices provided. Returns true if so.
func Is(reference uint32, choices ...uint32) bool {
	for _, c := range choices {
		if reference == c {
			return true
		}
	}
	return false
}

// GetSkillBook retrieves the index for the skill book, given the characters job.
func GetSkillBook(jobId uint32) uint32 {
	if jobId >= 2210 && jobId <= 2218 {
		return jobId - 2209
	}
	return 0
}

// IsFourthJob determines if the skill provided is a fourth job skill. Returns true if so.
func IsFourthJob(jobId uint16, skillId uint32) bool {
	if jobId == job.Evan4 {
		return false
	}

	if Is(skillId, EvanMagicMastery, EvanFlameWheel, EvanHerosWill, EvanDarkFog, EvanSoulStone) {
		return true
	}

	return jobId%10 == 2
}
