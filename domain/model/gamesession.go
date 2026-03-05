package model

import (
	"encoding/json"
	"os"
	"path/filepath"
	"rogue/infrastructure/constants"
	"sort"
	"fmt"
)
// структура для хранения статистики сессии
type GameStats struct {
	TotalTreasureCollected int
	DeepestLvlReached      int
	EnemysKilled           int
	FoodConsumed           int
	ElixirsConsumed        int
	ScrollsConsumed        int
	TotalHitsTaken         int
	TilesWalked            int
}

type GameSession struct {
	level        *Level
	LevelNumber  int
	Hero         Hero
	backpack     map[constants.EntityType][constants.BACKPACK_SIZE]*Item
	Message      string
	Effects      map[int]Effect
	SessionStats GameStats
	Allstats     []GameStats
}

func NewGameSession() *GameSession {
	hero := NewHero(
		NewCoordinate(0, 0),
		constants.HERO_MAX_HEALTH,
		constants.HERO_MAX_HEALTH,
		constants.HERO_AGILITY,
		constants.HERO_STRENGTH,
	)
	level := NewLevel(hero, 1)
	heroCoordinate := level.rooms[level.StartRoomIndex].GetRandomCoordinate()
	hero.Coordinate = heroCoordinate
	backpack := make(map[constants.EntityType][9]*Item, 4)
	effects := make(map[int]Effect)
	
	// Создаем сессию
	gameSession := &GameSession{
		level: level, 
		LevelNumber: 1, 
		Hero: hero, 
		backpack: backpack, 
		Effects: effects,
	}
	
	// Загружаем статистику
	if err := gameSession.LoadStats(); err != nil {
		// Если ошибка загрузки, просто пишем в сообщение или логируем
		gameSession.Message = "Failed to load stats: " + err.Error()
	}
	
	return gameSession
}

func CreateGameSession(
	level *Level,
	levelNumber int,
	hero Hero,
	backpack map[constants.EntityType][constants.BACKPACK_SIZE]*Item,
	message string,
	effects map[int]Effect,
) *GameSession {

	gameSession := &GameSession{
		level: level, 
		LevelNumber: levelNumber, 
		Hero: hero,
		backpack: backpack, 
		Message: message, 
		Effects: effects,
	}
	
	// Загружаем статистику
	if err := gameSession.LoadStats(); err != nil {
		gameSession.Message += " | Failed to load stats: " + err.Error()
	}
	
	return gameSession
}

func (gs GameSession) Cells() map[Coordinate]Cell {
	return gs.level.cells
}

func (gs GameSession) Rooms() []Room {
	return gs.level.rooms
}

func (gs GameSession) Enemies() map[Coordinate]Enemy {
	return gs.level.enemies
}

func (gs *GameSession) MoveEnemy(oldCoordinate, newCoordinate Coordinate, enemy Enemy) {
	delete(gs.level.enemies, oldCoordinate)
	gs.level.enemies[newCoordinate] = enemy
}

func (gs GameSession) Items() map[Coordinate]*Item {
	return gs.level.items
}

func (gs GameSession) AddItem(coordinate Coordinate, item *Item) {
	gs.level.items[coordinate] = item
}

func (gs GameSession) RemoveItem(coordinate Coordinate) {
	delete(gs.level.items, coordinate)
}

func (gs *GameSession) NextLevel() {
	gs.LevelNumber++
	//сохраняем статистику по достигнутой глубине
	gs.SessionStats.DeepestLvlReached = gs.LevelNumber
	gs.level = NewLevel(gs.Hero, gs.LevelNumber)
	heroCoordinate := gs.level.rooms[gs.level.StartRoomIndex].GetRandomCoordinate()
	gs.Hero.Coordinate = heroCoordinate
	gs.Hero.HasFirstKey = false
	gs.Hero.HasSecondKey = false
}

func (gs *GameSession) AddItemToBackpack(item *Item) bool {
	items := gs.backpack[item.ItemType]
	for i := range constants.BACKPACK_SIZE {
		if items[i] == nil {
			items[i] = item
			gs.backpack[item.ItemType] = items
			return true
		}
	}
	return false
}

func (gs *GameSession) RemoveItemFromBackpack(itemType constants.EntityType, index int) {
	items := gs.backpack[itemType]
	items[index] = nil
	gs.backpack[itemType] = items
}

func (gs *GameSession) GetItemsByType(itemType constants.EntityType) [constants.BACKPACK_SIZE]*Item {
	return gs.backpack[itemType]
}

func (gs GameSession) IsOccupied(coordinate Coordinate) bool {
	if gs.level.cells[coordinate].ground == constants.WALL || gs.level.cells[coordinate].ground == constants.EXIT {
		return true
	}
	if _, exists := gs.level.items[coordinate]; exists {
		return true
	}
	if _, exists := gs.level.enemies[coordinate]; exists {
		return true
	}
	if coordinate == gs.Hero.Coordinate {
		return true
	}
	return false
}


// SaveStats сохраняет текущую статистику в JSON файл
func (gs *GameSession) SaveStats() error {
	// Загружаем существующую статистику
	existingStats, err := gs.loadStatsFromFile()
	if err != nil {
		return err
	}

	// Добавляем новую статистику
	existingStats = append(existingStats, gs.SessionStats)

	// Сортируем по золотишку (или другому критерию)
	sort.Slice(existingStats, func(i, j int) bool {
		return existingStats[i].TotalTreasureCollected > existingStats[j].TotalTreasureCollected
	})

	// Оставляем только 5 последних записей
	if len(existingStats) > 5 {
		existingStats = existingStats[:5]
	}

	// Сохраняем обратно в файл
	return gs.saveStatsToFile(existingStats)
}

// LoadStats загружает статистику из JSON файла в Allstats
func (gs *GameSession) LoadStats() error {
	stats, err := gs.loadStatsFromFile()
	if err != nil {
		return err
	}
	
	gs.Allstats = stats
	return nil
}

// loadStatsFromFile загружает статистику из файла
func (gs *GameSession) loadStatsFromFile() ([]GameStats, error) {
	filePath := gs.getStatsFilePath()
	
	// Если файла нет, возвращаем пустой слайс
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return []GameStats{}, nil
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var stats []GameStats
	if err := json.Unmarshal(data, &stats); err != nil {
		return nil, err
	}

	return stats, nil
}

// saveStatsToFile сохраняет статистику в файл
func (gs *GameSession) saveStatsToFile(stats []GameStats) error {
	filePath := gs.getStatsFilePath()
	
	// Создаем директорию если нужно
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(stats, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, data, 0644)
}

// getStatsFilePath возвращает путь к файлу со статистикой
func (gs *GameSession) getStatsFilePath() string {
	return filepath.Join(".", "game_stats.json")
}

// GetLeaderboard возвращает статистику для отображения в виде строки
func (gs *GameSession) GetLeaderboard() string {
	if len(gs.Allstats) == 0 {
		return "No statistics available"
	}

	var result string
	result += "Leaderboard:\n"
	
	for i, stats := range gs.Allstats {
		result += fmt.Sprintf("%d. Lvl:%d Kills:%d Gold:%d Food:%d Elixirs:%d Scrolls:%d Hits:%d Steps:%d\n", 
			i+1,
			stats.DeepestLvlReached,
			stats.EnemysKilled,
			stats.TotalTreasureCollected,
			stats.FoodConsumed,
			stats.ElixirsConsumed,
			stats.ScrollsConsumed,
			stats.TotalHitsTaken,
			stats.TilesWalked)
	}
	
	return result
}

// FormatCurrentStats форматирует текущую статистику для отображения
func (gs *GameSession) FormatCurrentStats() string {
	stats := gs.SessionStats
	
	return fmt.Sprintf(
		"Stats:LVL:%d KILL:%d GLD:%d FOOD:%d DRINKS:%d SCRLS:%d HITS:%d STPS:%d",
		stats.DeepestLvlReached,
		stats.EnemysKilled,
		stats.TotalTreasureCollected,
		stats.FoodConsumed,
		stats.ElixirsConsumed,
		stats.ScrollsConsumed,
		stats.TotalHitsTaken,
		stats.TilesWalked,
	)
}


