package dto

type DomainToViewDTO interface {
	FieldInfo() []CellInfoDTO
	HeroInfo() HeroInfoDTO
	Level() int
	Message() string
	GameStatus() int
	Weapon() string
}
