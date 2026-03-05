package model

import (
	"rogue/infrastructure/constants"
)

type Cell struct {
	ground  constants.GroundType
	visible bool
	visited bool
}

func NewCell(ground constants.GroundType) *Cell {
	return &Cell{ground: ground, visible: false, visited: false}
}

func CreateCell(ground constants.GroundType, visible bool, visited bool) Cell {
	return Cell{ground: ground, visible: visible, visited: visited}
}

func (c Cell) Ground() constants.GroundType {
	return c.ground
}

func (c Cell) IsVisible() bool {
	return c.visible
}

func (c *Cell) SetVisible() {
	c.visible = true
	if c.ground == constants.WALL || c.ground == constants.PASSAGE || c.ground == constants.CORRIDOR {
		c.visited = true
	}
}

func (c *Cell) SetUnvisible() {
	c.visible = false
}

func (c Cell) IsVisited() bool {
	return c.visited
}
