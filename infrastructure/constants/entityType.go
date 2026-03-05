package constants

type EntityType int

const (
	// ключи и двери
	NONE EntityType = iota
	FIRST_KEY
	FIRST_DOOR
	SECOND_KEY
	SECOND_DOOR
	// предметы
	FOOD
	ELIXIR
	SCROLL
	WEAPON
	TREASURE
	// монстры
	VAMPIRE
	ZOMBIE
	OGRE
	GHOST
	SNAKE
	MIMIC
)
