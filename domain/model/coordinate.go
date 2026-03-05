package model

type Coordinate struct {
	row int
	col int
}

func NewCoordinate(row int, col int) Coordinate {
	return Coordinate{row: row, col: col}
}

func (c Coordinate) Row() int {
	return c.row
}

func (c Coordinate) Col() int {
	return c.col
}
