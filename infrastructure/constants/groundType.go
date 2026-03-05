package constants

type GroundType int

const (
	WALL GroundType = iota
	FLOOR
	CORRIDOR
	PASSAGE
	EXIT
)
