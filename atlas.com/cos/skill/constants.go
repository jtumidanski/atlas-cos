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

	FighterSwordMastery     uint32 = 1100000
	FighterAxeMastery       uint32 = 1100001
	FighterFinalAttackSword uint32 = 1100002
	FighterFinalAttackAxe   uint32 = 1100003
	FighterSwordBooster     uint32 = 1101004
	FighterAxeBooster       uint32 = 1101005
	FighterRage             uint32 = 1101006
	FighterPowerGuard       uint32 = 1101007

	PageSwordMastery           uint32 = 1200000
	PageBluntWeaponMastery     uint32 = 1200001
	PageFinalAttackSword       uint32 = 1200002
	PageFinalAttackBluntWeapon uint32 = 1200003
	PageSwordBooster           uint32 = 1201004
	PageBluntWeaponBooster     uint32 = 1201005
	PageThreaten               uint32 = 1201006
	PagePowerGuard             uint32 = 1201007

	SpearmanSpearMastery       uint32 = 1300000
	SpearmanPolearmMastery     uint32 = 1300001
	SpearmanFinalAttackSpear   uint32 = 1300002
	SpearmanFinalAttackPolearm uint32 = 1300003
	SpearmanSpearBooster       uint32 = 1301004
	SpearmanPolearmBooster     uint32 = 1301005
	SpearmanIronWill           uint32 = 1301006
	SpearmanHyperBody          uint32 = 1301007

	MagicianImprovedMPRecovery uint32 = 2000000
	MagicianImprovedMPIncrease uint32 = 2000001
	MagicianMagicGuard         uint32 = 2001002
	MagicianMagicArmor         uint32 = 2001003
	MagicianEnergyBolt         uint32 = 2001004
	MagicianMagicClaw          uint32 = 2001005

	FirePoisonWizardMPEater      uint32 = 2100000
	FirePoisonWizardMeditation   uint32 = 2101001
	FirePoisonWizardTeleport     uint32 = 2101002
	FirePoisonWizardSlow         uint32 = 2101003
	FirePoisonWizardFireArrow    uint32 = 2101004
	FirePoisonWizardPoisonBreath uint32 = 2101005

	IceLightningWizardMPEater     uint32 = 2200000
	IceLightningWizardMeditation  uint32 = 2201001
	IceLightningWizardTeleport    uint32 = 2201002
	IceLightningWizardSlow        uint32 = 2201003
	IceLightningWizardColdBeam    uint32 = 2201004
	IceLightningWizardThunderBolt uint32 = 2201005

	ClericMPEater    uint32 = 2300000
	ClericTeleport   uint32 = 2301001
	ClericHeal       uint32 = 2301002
	ClericInvincible uint32 = 2301003
	ClericBless      uint32 = 2301004
	ClericHolyArrow  uint32 = 2301005

	BowmanBlessingOfAmazon uint32 = 3000000
	BowmanCriticalShot     uint32 = 3000001
	BowmanTheEyeOfAmazon   uint32 = 3000002
	BowmanFocus            uint32 = 3001003
	BowmanArrowBlow        uint32 = 3001004
	BowmanDoubleShot       uint32 = 3001005

	HunterBowMastery     uint32 = 3100000
	HunterFinalAttack    uint32 = 3100001
	HunterBowBooster     uint32 = 3101002
	HunterPowerKnockback uint32 = 3101003
	HunterSoulArrow      uint32 = 3101004
	HunterArrowBomb      uint32 = 3101005

	CrossbowmanCrossbowMastery uint32 = 3200000
	CrossbowmanFinalAttack     uint32 = 3200001
	CrossbowmanCrossbowBooster uint32 = 3201002
	CrossbowmanPowerKnockback  uint32 = 3201003
	CrossbowmanSoulArrow       uint32 = 3201004
	CrossbowmanIronArrow       uint32 = 3201005

	ThiefNimbleBody uint32 = 4000000
	ThiefKeenEyes   uint32 = 4000001
	ThiefDisorder   uint32 = 4001002
	ThiefDarkSight  uint32 = 4001003
	ThiefDoubleStab uint32 = 4001334
	ThiefLuckySeven uint32 = 4001344

	AssassinClawMastery   uint32 = 4100000
	AssassinCriticalThrow uint32 = 4100001
	AssassinEndure        uint32 = 4100002
	AssassinClawBooster   uint32 = 4101003
	AssassinHaste         uint32 = 4101004
	AssassinDrain         uint32 = 4101005

	BanditDaggerMastery uint32 = 4200000
	BanditEndure        uint32 = 4200001
	BanditDaggerBooster uint32 = 4201002
	BanditHaste         uint32 = 4201003
	BanditSteal         uint32 = 4201004
	BanditSavageBlow    uint32 = 4201005

	PirateBulletTime     uint32 = 5000000
	PirateFlashFist      uint32 = 5001001
	PirateSomersaultKick uint32 = 5001002
	PirateDoubleShot     uint32 = 5001003
	PirateDash           uint32 = 5001005

	BrawlerImproveMaxHP    uint32 = 5100000
	BrawlerKnucklerMastery uint32 = 5100001
	BrawlerBackSpinBlow    uint32 = 5101002
	BrawlerDoubleUppercut  uint32 = 5101003
	BrawlerCorkscrewBlow   uint32 = 5101004
	BrawlerMPRecovery      uint32 = 5101005
	BrawlerKnucklerBooster uint32 = 5101006
	BrawlerOakBarrel       uint32 = 5101007

	GunslingerGunMastery    uint32 = 5200000
	GunslingerInvisibleShot uint32 = 5201001
	GunslingerGrenade       uint32 = 5201002
	GunslingerGunBooster    uint32 = 5201003
	GunslingerBlankShot     uint32 = 5201004
	GunslingerWings         uint32 = 5201005
	GunslingerRecoilShot    uint32 = 5201006

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
