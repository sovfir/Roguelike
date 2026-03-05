package dto

type HeroInfoDTO struct {
	Row        int
	Col        int
	Health     int
	MaxHealth  int
	Agility    int
	Strength   int
	Gold       int
	IsSleeping bool
}

func newHeroInfoDTO(row, col, health, maxHealth, agility, strength, gold int, isSleeping bool) HeroInfoDTO {
	return HeroInfoDTO{
		Row:        row,
		Col:        col,
		Health:     health,
		MaxHealth:  maxHealth,
		Agility:    agility,
		Strength:   strength,
		Gold:       gold,
		IsSleeping: isSleeping,
	}
}
