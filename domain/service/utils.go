package service

import (
	"fmt"
	"rogue/domain/model"
)

func CellByCoordinate(coordinate model.Coordinate, gs *model.GameSession) (model.Cell, bool) {
	cell, ok := gs.Cells()[coordinate]
	return cell, ok
}

func RoomByCoordinate(coordinate model.Coordinate, gs *model.GameSession) (model.Room, error) {
	rooms := gs.Rooms()
	for _, room := range rooms {
		if room.IsRoomCoordinate(coordinate) {
			return room, nil
		}
	}
	return rooms[0], fmt.Errorf("room not found by coordinates %v", coordinate)
}
