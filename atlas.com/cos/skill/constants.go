package skill

import "atlas-cos/job"

const (
	BeginnerThreeSnails uint32 = 1000
	BeginnerRecovery    uint32 = 1001
	BeginnerNimbleFeet  uint32 = 1002

	NoblesseThreeSnails uint32 = 10001000
	NoblesseRecovery    uint32 = 10001001
	NoblesseNimbleFeet  uint32 = 10001002

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
