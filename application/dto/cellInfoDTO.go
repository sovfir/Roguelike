package dto

import "rogue/infrastructure/constants"

type CellInfoDTO struct {
	Row    int
	Col    int
	Ground constants.GroundType
	Entity constants.EntityType
}

func newCellInfoDTO(row, col int, ground constants.GroundType, entity constants.EntityType) CellInfoDTO {
	return CellInfoDTO{
		Row:    row,
		Col:    col,
		Ground: ground,
		Entity: entity,
	}
}
