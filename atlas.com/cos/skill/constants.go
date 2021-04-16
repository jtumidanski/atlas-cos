package skill

const (
	EvanMagicMastery uint32 = 22170001
	EvanFlameWheel   uint32 = 22171003
	EvanHerosWill    uint32 = 22171004
	EvanDarkFog      uint32 = 22181002
	EvanSoulStone    uint32 = 22181003
)

func Is(reference uint32, choices ...uint32) bool {
	for _, c := range choices {
		if reference == c {
			return true
		}
	}
	return false
}
