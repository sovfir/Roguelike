package model

import "math/rand/v2"

type Room struct {
	Width         int
	Height        int
	TopLeftCorner Coordinate
}

func NewRoom(topLeftCorner Coordinate, width int, height int) Room {
	return Room{
		Width:         width,
		Height:        height,
		TopLeftCorner: topLeftCorner,
	}
}

func (r Room) IsRoomCoordinate(checkedCoordinate Coordinate) bool {
	minX := r.TopLeftCorner.col
	maxX := minX + r.Width - 1
	minY := r.TopLeftCorner.row
	maxY := minY + r.Height - 1
	x := checkedCoordinate.col
	y := checkedCoordinate.row

	return x >= minX && x <= maxX && y >= minY && y <= maxY
}

func (r Room) GetRandomCoordinate() Coordinate {
	randX := rand.IntN(r.Width-2) + 1
	randY := rand.IntN(r.Height-2) + 1
	return NewCoordinate(randY+r.TopLeftCorner.row, randX+r.TopLeftCorner.col)
}

func (r Room) AllCoordinates() []Coordinate {
	result := make([]Coordinate, 0, r.Height*r.Width)
	for y := range r.Height {
		for x := range r.Width {
			result = append(result, NewCoordinate(r.TopLeftCorner.row+y, r.TopLeftCorner.col+x))
		}
	}
	return result
}
