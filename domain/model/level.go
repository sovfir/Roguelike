package model

import (
	"math"
	"math/rand/v2"

	"rogue/infrastructure/constants"
)

type edge struct {
	room1 int
	room2 int
}

type Level struct {
	rooms          []Room
	StartRoomIndex int
	endRoomIndex   int
	cells          map[Coordinate]Cell
	items          map[Coordinate]*Item
	enemies        map[Coordinate]Enemy
}

func NewLevel(hero Hero, levelNumber int) *Level {
	newLevel := Level{
		rooms:   make([]Room, constants.ROOM_COUNT_X*constants.ROOM_COUNT_Y),
		cells:   make(map[Coordinate]Cell),
		items:   make(map[Coordinate]*Item),
		enemies: make(map[Coordinate]Enemy),
	}
	newLevel.generate(hero, levelNumber)
	return &newLevel
}

func CreateLevel(rooms []Room, cells map[Coordinate]Cell, items map[Coordinate]*Item,
	enemies map[Coordinate]Enemy) *Level {

	return &Level{
		rooms:          rooms,
		StartRoomIndex: 0,
		endRoomIndex:   0,
		cells:          cells,
		items:          items,
		enemies:        enemies,
	}
}

func (l *Level) generate(hero Hero, levelNumber int) {
	l.createRooms()
	l.createCorridors()
	l.createExit()
	balancer := calcBalancer(hero, levelNumber)
	SpawnItems(hero, levelNumber, l, balancer)
	SpawnEnemies(levelNumber, l, balancer)
}

func (l *Level) createRooms() {
	l.generateRooms()
	for i, room := range l.rooms {
		for row := 0; row < room.Height; row++ {
			for col := 0; col < room.Width; col++ {
				ground := constants.FLOOR
				if row == 0 || row == room.Height-1 || col == 0 || col == room.Width-1 {
					ground = constants.WALL
				}
				coordinate := NewCoordinate(row+room.TopLeftCorner.row, col+room.TopLeftCorner.col)
				cell := *NewCell(ground)
				if i == l.StartRoomIndex {
					cell.SetVisible()
				}
				l.cells[coordinate] = cell
			}
		}
	}
}

func (l *Level) generateRooms() {
	gridWidth := constants.FIELD_WIDTH / constants.ROOM_COUNT_X
	gridHeight := constants.FIELD_HEIGHT / constants.ROOM_COUNT_Y
	randomRangeY := gridHeight - constants.GRID_LINE_SPACES - constants.MIN_ROOM_SIZE
	randomRangeX := gridWidth - constants.GRID_LINE_SPACES - constants.MIN_ROOM_SIZE
	for row := range constants.ROOM_COUNT_Y {
		for col := range constants.ROOM_COUNT_X {
			startY := rand.IntN(randomRangeY + 1)
			maxHeiht := int(math.Min(constants.MAX_ROOM_SIZE-constants.MIN_ROOM_SIZE, float64(randomRangeY+1-startY)))
			height := rand.IntN(maxHeiht) + constants.MIN_ROOM_SIZE

			startX := rand.IntN(randomRangeX + 1)
			maxWidth := int(math.Min(constants.MAX_ROOM_SIZE-constants.MIN_ROOM_SIZE, float64(randomRangeX+1-startX)))
			width := rand.IntN(maxWidth) + constants.MIN_ROOM_SIZE
			topLeftCorner := NewCoordinate(startY+row*gridHeight, startX+col*gridWidth)
			l.rooms[row*constants.ROOM_COUNT_X+col] = NewRoom(topLeftCorner, width, height)
		}
	}
	startRoom := rand.IntN(len(l.rooms))
	l.StartRoomIndex = startRoom
}

func (l *Level) createCorridors() {
	edges := l.generateEdges()
	corridors := l.generateCorridors(edges)
	for _, corridor := range corridors {
		for _, coordinate := range corridor {
			if cell, ok := l.cells[coordinate]; ok && cell.ground == constants.WALL {
				cell.ground = constants.PASSAGE
				l.cells[coordinate] = cell
			} else {
				l.cells[coordinate] = *NewCell(constants.CORRIDOR)
			}
		}
	}
}

func (l *Level) generateEdges() (edges []edge) {
	// Генерация горизонтальных ребер между комнатами
	for i := 0; i < constants.ROOM_COUNT_Y; i++ {
		for j := 0; j < constants.ROOM_COUNT_X-1; j++ {
			currentRoom := i*constants.ROOM_COUNT_X + j
			edges = append(edges, edge{currentRoom, currentRoom + 1})
		}
	}
	// Генерация вертикальных ребер между комнатами
	for i := 0; i < constants.ROOM_COUNT_Y-1; i++ {
		for j := 0; j < constants.ROOM_COUNT_X; j++ {
			currentRoom := i*constants.ROOM_COUNT_X + j
			edges = append(edges, edge{currentRoom, currentRoom + constants.ROOM_COUNT_X})
		}
	}
	return
}

func (l *Level) generateCorridors(edges []edge) (corridors [][]Coordinate) {
	roomCount := constants.ROOM_COUNT_X * constants.ROOM_COUNT_Y
	linkedRooms := map[int]struct{}{l.StartRoomIndex: {}}
	currentZone := 1
	roomsByZones := map[int]int{l.StartRoomIndex: currentZone}

	for len(linkedRooms) < roomCount {
		switch {
		case len(roomsByZones) >= constants.SECOND_ZONE_CAP && currentZone == 2:
			l.addKey(currentZone, roomsByZones, l.StartRoomIndex)
			currentZone = 3
		case len(roomsByZones) >= constants.FIRST_ZONE_CAP && currentZone == 1:
			l.addKey(currentZone, roomsByZones, l.StartRoomIndex)
			currentZone = 2
		}
		possibleLinks := []int{}
		for i, edge := range edges {
			_, ok1 := linkedRooms[edge.room1]
			_, ok2 := linkedRooms[edge.room2]
			if (ok1 && !ok2) || (!ok1 && ok2) {
				possibleLinks = append(possibleLinks, i)
			}
		}

		randomEdgeIndex := rand.IntN(len(possibleLinks))
		newLinkIndex := possibleLinks[randomEdgeIndex]
		newLink := edges[newLinkIndex]
		possibleLinks = append(possibleLinks[:randomEdgeIndex], possibleLinks[randomEdgeIndex+1:]...)
		addedRoom := newLink.room2
		if _, ok := linkedRooms[addedRoom]; ok {
			addedRoom = newLink.room1
		}
		roomsByZones[addedRoom] = currentZone

		newCorridor := l.generateCorridor(newLink.room1, newLink.room2)
		corridors = append(corridors, newCorridor)
		l.addDoor(newCorridor, newLink.room1, newLink.room2, roomsByZones, currentZone)

		for _, i := range possibleLinks {
			if edges[i].room1 == addedRoom || edges[i].room2 == addedRoom {
				if rnd := rand.Float64(); rnd > 0.5 {
					newCorridor := l.generateCorridor(edges[i].room1, edges[i].room2)
					corridors = append(corridors, newCorridor)
					l.addDoor(newCorridor, edges[i].room1, edges[i].room2, roomsByZones, currentZone)
				}
			}
		}

		linkedRooms[addedRoom] = struct{}{}
		if len(linkedRooms) == roomCount {
			l.endRoomIndex = addedRoom
		}
	}
	return
}

func (l *Level) generateCorridor(room1, room2 int) []Coordinate {
	if room2 == room1+1 {
		return l.generateHorizontalCorridor(room1, room2)
	} else {
		return l.generateVerticalCorridor(room1, room2)
	}
}

func (l *Level) generateHorizontalCorridor(room1Idx, room2Idx int) []Coordinate {
	room1 := l.rooms[room1Idx]
	room2 := l.rooms[room2Idx]
	// У первой комнаты берем случайную координату и сдвигаем на правую стену
	corridorBegin := room1.GetRandomCoordinate()
	corridorBegin.col = room1.TopLeftCorner.col + room1.Width - 1
	// У второй комнаты берем случайную координату и сдвигаем на левую стену
	corridorEnd := room2.GetRandomCoordinate()
	corridorEnd.col = room2.TopLeftCorner.col

	var newCorridor []Coordinate
	// Если Y координаты равны, то строится прямой коридор, иначе - с изгибом
	if corridorBegin.row == corridorEnd.row {
		newCorridor = connectCoordinates(corridorBegin, corridorEnd)
	} else {
		bend := corridorBegin.col + 1 + rand.IntN(corridorEnd.col-corridorBegin.col-1)
		startBend := NewCoordinate(corridorBegin.row, bend)
		endBend := NewCoordinate(corridorEnd.row, bend)
		newCorridor = connectCoordinates(corridorBegin, startBend)
		newCorridor = append(newCorridor, connectCoordinates(startBend, endBend)...)
		newCorridor = append(newCorridor, connectCoordinates(endBend, corridorEnd)...)
	}

	return newCorridor
}

func (l *Level) generateVerticalCorridor(room1Idx, room2Idx int) []Coordinate {
	room1 := l.rooms[room1Idx]
	room2 := l.rooms[room2Idx]
	// У первой комнаты берем случайную координату и сдвигаем на нижнюю стену
	corridorBegin := room1.GetRandomCoordinate()
	corridorBegin.row = room1.TopLeftCorner.row + room1.Height - 1
	// У второй комнаты берем случайную координату и сдвигаем на верхнюю стену
	corridorEnd := room2.GetRandomCoordinate()
	corridorEnd.row = room2.TopLeftCorner.row

	var newCorridor []Coordinate
	// Если X координаты равны, то строится прямой коридор, иначе - с изгибом
	if corridorBegin.col == corridorEnd.col {
		newCorridor = connectCoordinates(corridorBegin, corridorEnd)
	} else {
		bend := corridorBegin.row + 1 + rand.IntN(corridorEnd.row-corridorBegin.row-1)
		startBend := NewCoordinate(bend, corridorBegin.col)
		endBend := NewCoordinate(bend, corridorEnd.col)
		newCorridor = connectCoordinates(corridorBegin, startBend)
		newCorridor = append(newCorridor, connectCoordinates(startBend, endBend)...)
		newCorridor = append(newCorridor, connectCoordinates(endBend, corridorEnd)...)
	}

	return newCorridor
}

func connectCoordinates(c1, c2 Coordinate) (result []Coordinate) {
	if c1.row == c2.row {
		start := math.Min(float64(c1.col), float64(c2.col))
		end := math.Max(float64(c1.col), float64(c2.col))
		for i := int(start); i <= int(end); i++ {
			result = append(result, NewCoordinate(c1.row, i))
		}
	} else {
		start := math.Min(float64(c1.row), float64(c2.row))
		end := math.Max(float64(c1.row), float64(c2.row))
		for i := int(start); i <= int(end); i++ {
			result = append(result, NewCoordinate(i, c1.col))
		}
	}
	return
}

func (l *Level) createExit() {
	end := l.rooms[l.endRoomIndex].GetRandomCoordinate()
	switch {
	case l.cells[NewCoordinate(end.row+1, end.col)].ground == constants.PASSAGE:
		end = NewCoordinate(end.row-1, end.col)
	case l.cells[NewCoordinate(end.row-1, end.col)].ground == constants.PASSAGE:
		end = NewCoordinate(end.row+1, end.col)
	case l.cells[NewCoordinate(end.row, end.col+1)].ground == constants.PASSAGE:
		end = NewCoordinate(end.row, end.col-1)
	case l.cells[NewCoordinate(end.row, end.col-1)].ground == constants.PASSAGE:
		end = NewCoordinate(end.row, end.col+1)
	}
	l.cells[end] = *NewCell(constants.EXIT)
}

func (l *Level) addDoor(corridor []Coordinate, room1Idx, room2Idx int,
	roomsByZone map[int]int, currentZone int) {
	room1Zone := roomsByZone[room1Idx]
	room2Zone := roomsByZone[room2Idx]
	if room1Zone != currentZone {
		l.items[corridor[0]] = doorItem(currentZone)
	}
	if room2Zone != currentZone {
		l.items[corridor[len(corridor)-1]] = doorItem(currentZone)
	}
}

func doorItem(currentZone int) *Item {
	if currentZone == 2 {
		return NewItem(constants.FIRST_DOOR, constants.FIRST_DOOR_NAME, 0, 0, 0, 0, 0)
	}
	return NewItem(constants.SECOND_DOOR, constants.SECOND_DOOR_NAME, 0, 0, 0, 0, 0)
}

func (l *Level) addKey(currentZone int, roomsByZone map[int]int, startRoomIndex int) {
	roomsList := make([]int, 0, len(roomsByZone))
	for key, val := range roomsByZone {
		if val == currentZone && key != startRoomIndex {
			roomsList = append(roomsList, key)
		}
	}
	randomIndex := rand.IntN(len(roomsList))
	roomIndex := roomsList[randomIndex]
	coordinateKey := l.rooms[roomIndex].GetRandomCoordinate()
	if currentZone == 1 {
		l.items[coordinateKey] = NewItem(constants.FIRST_KEY, constants.FIRST_KEY_NAME, 0, 0, 0, 0, 0)
	} else {
		l.items[coordinateKey] = NewItem(constants.SECOND_KEY, constants.SECOND_KEY_NAME, 0, 0, 0, 0, 0)
	}
}

func calcBalancer(hero Hero, levelNumber int) float64 {
	normal := 1.0 + float64((levelNumber-1)*10/100)
	healthNorm := constants.HERO_MAX_HEALTH * normal
	strNorm := constants.HERO_STRENGTH * normal
	aglNorm := constants.HERO_AGILITY * normal
	return (float64(hero.Health)/healthNorm + float64(hero.Strength)/strNorm + float64(hero.Agility)/aglNorm) / 3
}
