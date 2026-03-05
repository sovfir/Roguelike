package dto

import (
	"rogue/domain/model"
	"rogue/infrastructure/constants"
)

const FILENAME = "rogue.json"

type LevelSaveDTO struct {
	Rooms   []RoomSaveDTO
	Cells   map[int]CellDTO
	Items   map[int]model.Item
	Enemies map[int]MonsterDTO
}

type MonsterDTO struct {
	MonsterType constants.EntityType
	Name        string
	Health      int
	Agility     int
	Strength    int
	Hostility   int
	InBattle    bool
	FirstHit    bool
	Visibility  bool
	Resting     bool
	ItemMimic   constants.EntityType
}

type HeroDTO struct {
	Coordinate   int
	Health       int
	MaxHealth    int
	Agility      int
	Strength     int
	Weapon       model.Item
	HasFirstKey  bool
	HasSecondKey bool
	Gold         int
	IsSleeping   bool
}

type EffectSaveDTO struct {
	Item     model.Item
	TimeLeft int
}

type RoomSaveDTO struct {
	Width         int
	Height        int
	TopLeftCorner int
}

type CellDTO struct {
	Ground  constants.GroundType
	Visible bool
	Visited bool
}

type GameStatsDTO struct {
    TotalTreasureCollected int
    DeepestLvlReached      int
    EnemysKilled           int
    FoodConsumed           int
    ElixirsConsumed        int
    ScrollsConsumed        int
    TotalHitsTaken         int
    TilesWalked            int
}

type SaveGameDTO struct {
	Level       LevelSaveDTO
	LevelNumber int
	Hero        HeroDTO
	Backpack    map[constants.EntityType][constants.BACKPACK_SIZE]model.Item
	Effects     map[int]EffectSaveDTO
	SessionStats GameStatsDTO // добавляем статистику
}
