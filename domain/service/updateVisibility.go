package service

import (
	"fmt"

	"rogue/domain/model"
	"rogue/infrastructure/constants"
)

func UpdateVisibility(oldCoordinate, newCoordinate model.Coordinate, gs *model.GameSession) error {
	oldCell, _ := CellByCoordinate(oldCoordinate, gs)
	oldGround := oldCell.Ground()
	newCell, _ := CellByCoordinate(newCoordinate, gs)
	newGround := newCell.Ground()
	switch {
	case oldGround == constants.FLOOR && newGround == constants.FLOOR:
	case oldGround == constants.FLOOR && newGround == constants.PASSAGE:
		showNeighbours(newCoordinate, gs)
	case oldGround == constants.PASSAGE && newGround == constants.FLOOR:
		hideNeighbours(oldCoordinate, gs)
		showNeighbours(newCoordinate, gs)
	case oldGround == constants.PASSAGE && newGround == constants.CORRIDOR:
		room, err := RoomByCoordinate(oldCoordinate, gs)
		if err != nil {
			return fmt.Errorf("from passage to corridor: %w", err)
		}
		hideRoom(room, gs)
		showNeighbours(newCoordinate, gs)
		setPartialVisibility(newCoordinate, oldCoordinate, gs)
	case oldGround == constants.CORRIDOR && newGround == constants.PASSAGE:
		room, err := RoomByCoordinate(newCoordinate, gs)
		if err != nil {
			return fmt.Errorf("from corridor to passage: %w", err)
		}
		hideNeighbours(oldCoordinate, gs)
		showRoom(room, gs)
		showNeighbours(newCoordinate, gs)
	case oldGround == constants.CORRIDOR && newGround == constants.CORRIDOR:
		hideNeighbours(oldCoordinate, gs)
		showNeighbours(newCoordinate, gs)
		if passageCoordinate, ok := searchPassage(newCoordinate, gs); ok {
			setPartialVisibility(newCoordinate, passageCoordinate, gs)
		}
		if passageCoordinate, ok := searchPassage(oldCoordinate, gs); ok {
			room, err := RoomByCoordinate(passageCoordinate, gs)
			if err != nil {
				return fmt.Errorf("from corridor to corridor: %w", err)
			}
			hideRoom(room, gs)
		}
	default:
		return fmt.Errorf("")
	}
	return nil
}

func showRoom(room model.Room, gs *model.GameSession) {
	cells := gs.Cells()
	for _, coordinate := range room.AllCoordinates() {
		cell := cells[coordinate]
		cell.SetVisible()
		cells[coordinate] = cell
	}
}

func hideRoom(room model.Room, gs *model.GameSession) {
	cells := gs.Cells()
	for _, coordinate := range room.AllCoordinates() {
		cell := cells[coordinate]
		cell.SetUnvisible()
		cells[coordinate] = cell
	}
}

func showNeighbours(coordinate model.Coordinate, gs *model.GameSession) {
	cells := gs.Cells()
	for _, d := range model.Directions() {
		neighbour := model.NewCoordinate(coordinate.Row()+d.DeltaY, coordinate.Col()+d.DeltaX)
		cell, ok := cells[neighbour]
		if ok {
			cell.SetVisible()
			cells[neighbour] = cell
		}
	}
	cell, ok := cells[coordinate]
	if ok {
		cell.SetVisible()
		cells[coordinate] = cell
	}
}

func hideNeighbours(coordinate model.Coordinate, gs *model.GameSession) {
	cells := gs.Cells()
	for _, d := range model.Directions() {
		neighbour := model.NewCoordinate(coordinate.Row()+d.DeltaY, coordinate.Col()+d.DeltaX)
		cell, ok := cells[neighbour]
		if ok {
			cell.SetUnvisible()
			cells[neighbour] = cell
		}
	}
}

func searchPassage(coordinate model.Coordinate, gs *model.GameSession) (model.Coordinate, bool) {
	cells := gs.Cells()
	for _, d := range model.Directions() {
		neighbour := model.NewCoordinate(coordinate.Row()+d.DeltaY, coordinate.Col()+d.DeltaX)
		cell, ok := cells[neighbour]
		if ok && cell.Ground() == constants.PASSAGE {
			return neighbour, true
		}
	}
	return coordinate, false
}

func setPartialVisibility(heroCoordinate, passageCoordinate model.Coordinate, gs *model.GameSession) {
	cells := gs.Cells()
	isVertical := heroCoordinate.Row() == passageCoordinate.Row()
	room, _ := RoomByCoordinate(passageCoordinate, gs)
	for _, coordinate := range room.AllCoordinates() {
		dx := coordinate.Row() - passageCoordinate.Row()
		if dx < 0 {
			dx *= -1
		}
		dy := coordinate.Col() - passageCoordinate.Col()
		if dy < 0 {
			dy *= -1
		}
		if (isVertical && dx <= dy/2) || (!isVertical && dy <= dx/2) {
			cell, ok := cells[coordinate]
			if ok {
				cell.SetVisible()
				cells[coordinate] = cell
			}
		}
	}
}
