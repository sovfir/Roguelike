package data

import (
	"rogue/application/dto"
	"rogue/domain/model"
	"rogue/infrastructure/constants"
)

func GameSessionToDTO(gs *model.GameSession) dto.SaveGameDTO {
	result := dto.SaveGameDTO{
		Level:       levelToDTO(gs),
		LevelNumber: gs.LevelNumber,
		Hero: dto.HeroDTO{
			Coordinate:   coordinateToInt(gs.Hero.Coordinate),
			Health:       gs.Hero.Health,
			MaxHealth:    gs.Hero.MaxHealth,
			Agility:      gs.Hero.Agility,
			Strength:     gs.Hero.Strength,
			HasFirstKey:  gs.Hero.HasFirstKey,
			HasSecondKey: gs.Hero.HasSecondKey,
			Gold:         gs.Hero.Gold,
			IsSleeping:   gs.Hero.IsSleeping,
		},
		Backpack: make(map[constants.EntityType][constants.BACKPACK_SIZE]model.Item, 4),
		Effects:  make(map[int]dto.EffectSaveDTO, len(gs.Effects)),
		// Добавляем статистику
		SessionStats: gameStatsToDTO(gs.SessionStats),
	}
	if gs.Hero.Weapon != nil {
		result.Hero.Weapon = *gs.Hero.Weapon
	}

	result.Backpack[constants.FOOD] = backpackToDTO(gs.GetItemsByType(constants.FOOD))
	result.Backpack[constants.ELIXIR] = backpackToDTO(gs.GetItemsByType(constants.ELIXIR))
	result.Backpack[constants.SCROLL] = backpackToDTO(gs.GetItemsByType(constants.SCROLL))
	result.Backpack[constants.WEAPON] = backpackToDTO(gs.GetItemsByType(constants.WEAPON))

	for key, val := range gs.Effects {
		effectDTO := dto.EffectSaveDTO{
			TimeLeft: val.TimeLeft,
		}
		if val.Item != nil {
			effectDTO.Item = *val.Item
		}
		result.Effects[key] = effectDTO
	}
	return result
}

func coordinateToInt(coordinate model.Coordinate) int {
	return coordinate.Row()*constants.FIELD_WIDTH + coordinate.Col()
}

func intToCoordinate(value int) model.Coordinate {
	row := value / constants.FIELD_WIDTH
	col := value % constants.FIELD_WIDTH
	return model.NewCoordinate(row, col)
}

func levelToDTO(gs *model.GameSession) dto.LevelSaveDTO {
	result := dto.LevelSaveDTO{}

	rooms := gs.Rooms()
	roomsDTO := make([]dto.RoomSaveDTO, 0, len(rooms))
	for _, room := range rooms {
		roomDTO := dto.RoomSaveDTO{
			Width:         room.Width,
			Height:        room.Height,
			TopLeftCorner: coordinateToInt(room.TopLeftCorner),
		}
		roomsDTO = append(roomsDTO, roomDTO)
	}
	result.Rooms = roomsDTO

	cells := gs.Cells()
	cellsDTO := make(map[int]dto.CellDTO, len(cells))
	for key, val := range cells {
		cellsDTO[coordinateToInt(key)] = dto.CellDTO{
			Ground:  val.Ground(),
			Visible: val.IsVisible(),
			Visited: val.IsVisited(),
		}
	}
	result.Cells = cellsDTO

	items := gs.Items()
	itemsDTO := make(map[int]model.Item, len(items))
	for key, val := range items {
		itemsDTO[coordinateToInt(key)] = *val
	}
	result.Items = itemsDTO

	enemies := gs.Enemies()
	enemiesDTO := make(map[int]dto.MonsterDTO, len(enemies))
	for key, val := range enemies {
		enemiesDTO[coordinateToInt(key)] = monsterToDTO(val)
	}
	result.Enemies = enemiesDTO

	return result
}

func backpackToDTO(items [constants.BACKPACK_SIZE]*model.Item) (result [constants.BACKPACK_SIZE]model.Item) {
	for i, item := range items {
		if item != nil {
			result[i] = *item
		}
	}
	return
}

func backpackFromDTO(items [constants.BACKPACK_SIZE]model.Item) (result [constants.BACKPACK_SIZE]*model.Item) {
	for i, item := range items {
		if item.ItemType != 0 {
			result[i] = &item
		}
	}
	return
}

func DTOtoGameSession(gsDTO dto.SaveGameDTO) *model.GameSession {
	rooms := make([]model.Room, 0, len(gsDTO.Level.Rooms))
	for _, room := range gsDTO.Level.Rooms {
		rooms = append(rooms, model.NewRoom(intToCoordinate(room.TopLeftCorner), room.Width, room.Height))
	}
	cells := make(map[model.Coordinate]model.Cell, len(gsDTO.Level.Cells))
	for key, cell := range gsDTO.Level.Cells {
		cells[intToCoordinate(key)] = model.CreateCell(cell.Ground, cell.Visible, cell.Visited)
	}
	items := make(map[model.Coordinate]*model.Item, len(gsDTO.Level.Items))
	for key, item := range gsDTO.Level.Items {
		items[intToCoordinate(key)] = &item
	}
	enemies := make(map[model.Coordinate]model.Enemy, len(gsDTO.Level.Enemies))
	for key, enemy := range gsDTO.Level.Enemies {
		enemies[intToCoordinate(key)] = monsterFromDTO(enemy)
	}
	level := model.CreateLevel(rooms, cells, items, enemies)

	hero := model.NewHero(intToCoordinate(gsDTO.Hero.Coordinate), gsDTO.Hero.Health, gsDTO.Hero.MaxHealth,
		gsDTO.Hero.Agility, gsDTO.Hero.Strength)
	hero.Weapon = &gsDTO.Hero.Weapon
	hero.HasFirstKey = gsDTO.Hero.HasFirstKey
	hero.HasSecondKey = gsDTO.Hero.HasSecondKey
	hero.IsSleeping = gsDTO.Hero.IsSleeping

	backpack := make(map[constants.EntityType][constants.BACKPACK_SIZE]*model.Item, 4)
	backpack[constants.FOOD] = backpackFromDTO(gsDTO.Backpack[constants.FOOD])
	backpack[constants.ELIXIR] = backpackFromDTO(gsDTO.Backpack[constants.ELIXIR])
	backpack[constants.SCROLL] = backpackFromDTO(gsDTO.Backpack[constants.SCROLL])
	backpack[constants.WEAPON] = backpackFromDTO(gsDTO.Backpack[constants.WEAPON])

	effects := make(map[int]model.Effect, len(gsDTO.Effects))
	for key, value := range gsDTO.Effects {
		effects[key] = model.NewEffect(&value.Item, value.TimeLeft)
	}

	// Создаем GameSession со статистикой
	gameSession := model.CreateGameSession(level, gsDTO.LevelNumber, hero, backpack, "Game loaded", effects)
	
	// Восстанавливаем статистику
	gameSession.SessionStats = gameStatsFromDTO(gsDTO.SessionStats)
	return gameSession
}

func monsterToDTO(enemy model.Enemy) dto.MonsterDTO {
	switch enemy.Type() {
	case constants.ZOMBIE:
		e := enemy.(model.Zombie)
		return dto.MonsterDTO{
			MonsterType: e.MonsterType,
			Name:        e.Name,
			Health:      e.Health,
			Agility:     e.Agility,
			Strength:    e.Strength,
			Hostility:   e.Hostility,
			InBattle:    e.InBattle,
		}
	case constants.VAMPIRE:
		e := enemy.(model.Vampire)
		return dto.MonsterDTO{
			MonsterType: e.MonsterType,
			Name:        e.Name,
			Health:      e.Health,
			Agility:     e.Agility,
			Strength:    e.Strength,
			Hostility:   e.Hostility,
			InBattle:    e.InBattle,
			FirstHit:    e.FirstHit,
		}
	case constants.GHOST:
		e := enemy.(model.Ghost)
		return dto.MonsterDTO{
			MonsterType: e.MonsterType,
			Name:        e.Name,
			Health:      e.Health,
			Agility:     e.Agility,
			Strength:    e.Strength,
			Hostility:   e.Hostility,
			InBattle:    e.InBattle,
			Visibility:  e.Visibility,
		}
	case constants.OGRE:
		e := enemy.(model.Ogre)
		return dto.MonsterDTO{
			MonsterType: e.MonsterType,
			Name:        e.Name,
			Health:      e.Health,
			Agility:     e.Agility,
			Strength:    e.Strength,
			Hostility:   e.Hostility,
			InBattle:    e.InBattle,
			Resting:     e.Resting,
		}
	case constants.SNAKE:
		e := enemy.(model.Snake)
		return dto.MonsterDTO{
			MonsterType: e.MonsterType,
			Name:        e.Name,
			Health:      e.Health,
			Agility:     e.Agility,
			Strength:    e.Strength,
			Hostility:   e.Hostility,
			InBattle:    e.InBattle,
		}
	case constants.MIMIC:
		e := enemy.(model.Mimic)
		return dto.MonsterDTO{
			MonsterType: e.MonsterType,
			Name:        e.Name,
			Health:      e.Health,
			Agility:     e.Agility,
			Strength:    e.Strength,
			Hostility:   e.Hostility,
			InBattle:    e.InBattle,
			ItemMimic:   e.ItemMimic,
		}
	}
	return dto.MonsterDTO{}
}

func monsterFromDTO(enemy dto.MonsterDTO) model.Enemy {
	monster := model.Monster{
		MonsterType: enemy.MonsterType,
		Name:        enemy.Name,
		Health:      enemy.Health,
		Agility:     enemy.Agility,
		Strength:    enemy.Strength,
		Hostility:   enemy.Hostility,
		InBattle:    enemy.InBattle,
	}
	switch enemy.MonsterType {
	case constants.VAMPIRE:
		return model.Vampire{Monster: monster, FirstHit: enemy.FirstHit}
	case constants.GHOST:
		return model.Ghost{Monster: monster, Visibility: enemy.Visibility}
	case constants.OGRE:
		return model.Ogre{Monster: monster, Resting: enemy.Resting}
	case constants.SNAKE:
		return model.Snake{Monster: monster}
	case constants.MIMIC:
		return model.Mimic{Monster: monster, ItemMimic: enemy.ItemMimic}
	}
	return model.Zombie{Monster: monster}
}

// gameStatsToDTO преобразует модель статистики в DTO
func gameStatsToDTO(stats model.GameStats) dto.GameStatsDTO {
	return dto.GameStatsDTO{
		TotalTreasureCollected: stats.TotalTreasureCollected,
		DeepestLvlReached:      stats.DeepestLvlReached,
		EnemysKilled:           stats.EnemysKilled,
		FoodConsumed:           stats.FoodConsumed,
		ElixirsConsumed:        stats.ElixirsConsumed,
		ScrollsConsumed:        stats.ScrollsConsumed,
		TotalHitsTaken:         stats.TotalHitsTaken,
		TilesWalked:            stats.TilesWalked,
	}
}


// gameStatsFromDTO преобразует DTO статистики в модель
func gameStatsFromDTO(statsDTO dto.GameStatsDTO) model.GameStats {
	return model.GameStats{
		TotalTreasureCollected: statsDTO.TotalTreasureCollected,
		DeepestLvlReached:      statsDTO.DeepestLvlReached,
		EnemysKilled:           statsDTO.EnemysKilled,
		FoodConsumed:           statsDTO.FoodConsumed,
		ElixirsConsumed:        statsDTO.ElixirsConsumed,
		ScrollsConsumed:        statsDTO.ScrollsConsumed,
		TotalHitsTaken:         statsDTO.TotalHitsTaken,
		TilesWalked:            statsDTO.TilesWalked,
	}
}


