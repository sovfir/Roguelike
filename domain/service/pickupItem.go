package service

import (
	"rogue/domain/model"
	"rogue/infrastructure/constants"
)

func PickUpItem(coordinate model.Coordinate, item *model.Item, gs *model.GameSession) bool {
	var pickUpMessage string
	switch item.ItemType {
	case constants.FIRST_KEY:
		gs.Hero.HasFirstKey = true
		pickUpMessage = item.Name + " recieved"
	case constants.FIRST_DOOR:
		if !gs.Hero.HasFirstKey {
			gs.Message = gs.Message + " " + constants.FIRST_KEY_NAME + " required"
			return false
		}
		pickUpMessage = constants.FIRST_DOOR_NAME + " was open"
	case constants.SECOND_KEY:
		gs.Hero.HasSecondKey = true
		pickUpMessage = item.Name + " recieved"
	case constants.SECOND_DOOR:
		if !gs.Hero.HasSecondKey {
			gs.Message = gs.Message + " " + constants.SECOND_KEY_NAME + " required"
			return false
		}
		pickUpMessage = constants.SECOND_DOOR_NAME + " was open"
	case constants.TREASURE:
		gs.Hero.Gold += item.Value
		// Сохраняем статистику по золоту
		gs.SessionStats.TotalTreasureCollected += item.Value
		pickUpMessage = item.Name + " was picked up"
	default:
		ok := gs.AddItemToBackpack(item)
		if !ok {
			gs.Message = gs.Message + " " + " backpack is full"
			return false
		}
		pickUpMessage = item.Name + " was picked up"
	}

	gs.RemoveItem(coordinate)
	gs.Message = gs.Message + " " + pickUpMessage

	return true
}
