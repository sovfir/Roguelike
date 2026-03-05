package usecases

import (
	"fmt"
	"rogue/domain/model"
	"rogue/domain/service"
	"rogue/infrastructure/constants"
)

func MoveHero(direction constants.Directions, gs *model.GameSession) error {
	if gs.Hero.IsSleeping {
		gs.Hero.IsSleeping = false
	} else {
		oldCoordinate := gs.Hero.Coordinate
		//считаем шаги героя для статистики
		gs.SessionStats.TilesWalked++
		newCoordinate := model.NewCoordinate(oldCoordinate.Row()+direction.DeltaY, oldCoordinate.Col()+direction.DeltaX)
		newCell, ok := service.CellByCoordinate(newCoordinate, gs)
		if !ok || newCell.Ground() == constants.WALL {
			return nil
		}

		if newCell.Ground() == constants.EXIT {
			gs.NextLevel()
			return nil
		}

		if item, exists := gs.Items()[newCoordinate]; exists {
			if !service.PickUpItem(newCoordinate, item, gs) {
				return nil
			}
		}

		if _, exists := gs.Enemies()[newCoordinate]; exists {
			gs.Hero.Attack(gs, newCoordinate)
		} else {
			gs.Hero.Coordinate = newCoordinate

			if err := service.UpdateVisibility(oldCoordinate, newCoordinate, gs); err != nil {
				return fmt.Errorf("update visibility: %w", err)
			}
		}
	}

	EnemiesMove(gs)

	service.DecreaseEffectTime(gs)

	return nil
}
