package model

import (
	"container/list"
	"math/rand/v2"
	"slices"
	"strconv"
	"rogue/infrastructure/constants"
)

type Enemy interface {
	RegularMove(room Room, coordinate Coordinate, gs *GameSession) Coordinate
	PathFinder(start, end Coordinate, gs *GameSession) Coordinate
	Damage(heroAgility int, isOgre bool) (int, string)
	Attack(gs *GameSession)
	LooksLike() constants.EntityType
	Type() constants.EntityType
	HitCheck(heroAgility int) bool
	GetAgility() int
	GetName() string
	InPursuing() bool
	InAgrDistance(coordinate, heroCoordinate Coordinate) bool
}

type Monster struct {
	MonsterType constants.EntityType
	Name        string
	Health      int
	Agility     int
	Strength    int
	Hostility   int
	InBattle    bool
}

func (m Monster) LooksLike() constants.EntityType {
	return m.MonsterType
}

func (m Monster) Type() constants.EntityType {
	return m.MonsterType
}

func (m Monster) HitCheck(heroAgility int) bool {
	return hitCheck(heroAgility, m.Agility)
}

func (m Monster) GetAgility() int {
	return m.Agility
}

func (m Monster) GetName() string {
	return m.Name
}

func (m Monster) InPursuing() bool {
	return m.InBattle
}

func (m Monster) Damage(heroAgility int, isOgre bool) (int, string) {
	if !isOgre && !hitCheck(m.Agility, heroAgility) {
		return 0, " " + m.Name + " missed"
	}

	damageCoefficient := 1 - constants.DAMAGE_DISPERSION/2 + rand.Float64()*constants.DAMAGE_DISPERSION
	damage := float64(m.Strength) * damageCoefficient
	critChance := constants.BASE_CRIT_CHANCE * float64(m.Agility) / float64(heroAgility)
	if rand.Float64() < critChance {
		damage *= constants.CRIT_DAMAGE_COEFFICIENT
	}

	return int(damage), " " + m.Name + " dealt " + strconv.Itoa(int(damage)) + " damage"
}

func (m Monster) Attack(gs *GameSession) {
	damage, message := m.Damage(gs.Hero.Agility, false)
	gs.Hero.Health -= damage
	// статы урона
	gs.SessionStats.TotalHitsTaken++
	gs.Message += message
}

func (m Monster) InAgrDistance(coordinate, heroCoordinate Coordinate) bool {
	distance := abs(coordinate.col-heroCoordinate.col) + abs(coordinate.row-heroCoordinate.row)
	return distance <= m.Hostility
}

func (m Monster) PathFinder(start, end Coordinate, gs *GameSession) Coordinate {
	frontier := list.New()
	frontier.PushBack(start)
	cameFrom := make(map[Coordinate]Coordinate)

	for frontier.Len() > 0 {
		currentPtr := frontier.Front()
		current := currentPtr.Value.(Coordinate)
		frontier.Remove(currentPtr)
		if current == end {
			break
		}

		for _, d := range Directions() {
			neighbour := NewCoordinate(current.row+d.DeltaY, current.col+d.DeltaX)
			if neighbour == end {
				cameFrom[end] = current
				frontier.Init()
				break
			}
			if _, exists := cameFrom[neighbour]; !exists && !gs.IsOccupied(neighbour) {
				cameFrom[neighbour] = current
				frontier.PushBack(neighbour)
			}
		}
	}

	next := end
	for {
		if coordinate, exists := cameFrom[next]; !exists {
			return start
		} else if coordinate == start {
			return next
		} else {
			next = coordinate
		}
	}
}

func (m *Monster) TakeDamage(damage int) int {
	m.Health -= damage
	if m.Health <= 0 {
		return (m.Agility + m.Strength + m.Hostility) / 10
	}
	return 0
}

func (m *Monster) StartPursuing() {
	m.InBattle = true
}

type Zombie struct {
	Monster
}

func (z Zombie) RegularMove(room Room, coordinate Coordinate, gs *GameSession) Coordinate {
	return selectCoordinate(room, coordinate, Directions(), gs)
}

type Vampire struct {
	Monster
	FirstHit bool
}

func (s Vampire) RegularMove(room Room, coordinate Coordinate, gs *GameSession) Coordinate {
	directions := []constants.Directions{
		constants.LeftUp(),
		constants.Up(),
		constants.RightUp(),
		constants.Left(),
		constants.Right(),
		constants.LeftDown(),
		constants.Down(),
		constants.RightDown(),
	}

	return selectCoordinate(room, coordinate, directions, gs)
}

func (v Vampire) Attack(gs *GameSession) {
	damage, message := v.Damage(gs.Hero.Agility, false)
	gs.Hero.Health -= damage
	gs.SessionStats.TotalHitsTaken++
	gs.Message += message
	if damage != 0 {
		maxHealthDecrease := float64(damage) * constants.VAMPIRE_ABILITY_RATE
		gs.Hero.MaxHealth -= int(maxHealthDecrease)
	}
}

type Ghost struct {
	Monster
	Visibility bool
}

func (g Ghost) LooksLike() constants.EntityType {
	if g.Visibility {
		return g.MonsterType
	}
	return constants.NONE
}

func (g *Ghost) StartPursuing() {
	g.InBattle = true
	g.Visibility = true
}

func (g Ghost) RegularMove(room Room, coordinate Coordinate, gs *GameSession) Coordinate {
	for range 5 {
		newCoordinate := room.GetRandomCoordinate()
		if !gs.IsOccupied(newCoordinate) {
			return newCoordinate
		}
	}
	return coordinate
}

type Ogre struct {
	Monster
	Resting bool
}

func (o Ogre) RegularMove(room Room, coordinate Coordinate, gs *GameSession) Coordinate {
	directions := []constants.Directions{
		constants.LeftUp(),
		constants.Up(),
		constants.RightUp(),
		constants.Left(),
		constants.Right(),
		constants.LeftDown(),
		constants.Down(),
		constants.RightDown(),
	}

	firstStep := selectCoordinate(room, coordinate, directions, gs)
	return selectCoordinate(room, firstStep, directions, gs)
}

func (o Ogre) Attack(gs *GameSession) {
	damage, message := o.Damage(gs.Hero.Agility, true)
	gs.Hero.Health -= damage
	// статы урона
	gs.SessionStats.TotalHitsTaken++
	gs.Message += message
}

type Snake struct {
	Monster
}

func (s Snake) RegularMove(room Room, coordinate Coordinate, gs *GameSession) Coordinate {
	directions := []constants.Directions{
		constants.LeftUp(),
		constants.RightUp(),
		constants.LeftDown(),
		constants.RightDown(),
	}

	return selectCoordinate(room, coordinate, directions, gs)
}

func (s Snake) Attack(gs *GameSession) {
	damage, message := s.Damage(gs.Hero.Agility, false)
	gs.Hero.Health -= damage
	gs.Message += message
	if rand.Float64() < constants.SNAKE_SLEEP_CHANCE {
		gs.Hero.IsSleeping = true
		gs.Message += " Hero is sleeping"
	}
}

type Mimic struct {
	Monster
	ItemMimic constants.EntityType
}

func (m Mimic) LooksLike() constants.EntityType {
	if m.InBattle {
		return m.MonsterType
	}
	return m.ItemMimic
}

func (m Mimic) RegularMove(room Room, coordinate Coordinate, gs *GameSession) Coordinate {
	return coordinate
}

func hitCheck(attAgility, defAgility int) bool {
	chance := float64(attAgility) / float64(defAgility)
	if chance > 1 {
		chance = 1
	}

	return rand.Float64() <= constants.BASE_HIT_CHANCE+chance*constants.RAND_HIT_CHANCE
}

func selectCoordinate(room Room, coordinate Coordinate, directions []constants.Directions, gs *GameSession) Coordinate {
	inRoom := room.AllCoordinates()
	newCoordinates := make([]Coordinate, 0, 8)

	for _, d := range directions {
		newCoordinate := NewCoordinate(coordinate.row+d.DeltaY, coordinate.col+d.DeltaX)
		if !gs.IsOccupied(newCoordinate) && slices.Contains(inRoom, newCoordinate) {
			newCoordinates = append(newCoordinates, newCoordinate)
		}
	}

	if len(newCoordinates) == 0 {
		return coordinate
	}
	return newCoordinates[rand.IntN(len(newCoordinates))]
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
