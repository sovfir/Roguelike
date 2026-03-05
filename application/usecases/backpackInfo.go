package usecases

import (
	"rogue/application/dto"
	"rogue/domain/model"
	"rogue/domain/service"
	"rogue/infrastructure/constants"
	"strconv"
)

func BackpackInfo(
	itemType constants.EntityType,
	gs *model.GameSession,
	inputCh <-chan string,
	backpackInfoCh chan<- dto.BackpackDTO) {

	items := gs.GetItemsByType(itemType)
	backpackInfoDTO := dto.NewBackpackInfoDTO(items, itemType)
	backpackInfoCh <- backpackInfoDTO
	itemIndex := <-inputCh
	if itemType == constants.WEAPON && itemIndex == "0" {
		service.Unequip(gs)
		return
	}
	if itemIndex == "0" {
		return
	}
	if i, err := strconv.Atoi(itemIndex); err == nil {
		service.UseItem(items[i-1], gs)
		gs.Message += " " + items[i-1].Name + " used"
		gs.RemoveItemFromBackpack(itemType, i-1)
	}
}
