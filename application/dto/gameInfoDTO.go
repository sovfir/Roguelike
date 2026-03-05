package dto

type GameInfoDTO struct {
	cells   []CellInfoDTO
	hero    HeroInfoDTO
	level   int
	message string
	status  int
	weapon  string
}

func newGameInfoDTO(cells []CellInfoDTO, hero HeroInfoDTO, level int, message string,
	status int, weapon string) GameInfoDTO {

	return GameInfoDTO{
		cells:   cells,
		hero:    hero,
		level:   level,
		message: message,
		status:  status,
		weapon:  weapon,
	}
}

func (t GameInfoDTO) FieldInfo() []CellInfoDTO {
	return t.cells
}

func (t GameInfoDTO) HeroInfo() HeroInfoDTO {
	return t.hero
}

func (t GameInfoDTO) Level() int {
	return t.level
}

func (t GameInfoDTO) Message() string {
	return t.message
}

// 0 - игра идет, 1 - выигрыш, -1 - проигрыш
func (t GameInfoDTO) GameStatus() int {
	return t.status
}

func (t GameInfoDTO) Weapon() string {
	return t.weapon
}
