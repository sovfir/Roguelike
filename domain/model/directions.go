package model

import "rogue/infrastructure/constants"

func Directions() []constants.Directions {
	return []constants.Directions{
		constants.Up(),
		constants.Down(),
		constants.Left(),
		constants.Right(),
	}
}
