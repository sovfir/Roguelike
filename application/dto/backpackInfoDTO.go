package dto

import (
	"rogue/domain/model"
	"rogue/infrastructure/constants"
)

type BackpackDTO interface {
	Items() []string
}

type BackpackInfoDTO struct {
	// itemsType constants.EntityType
	items []string
}

func (b BackpackInfoDTO) Items() []string {
	return b.items
}

func NewBackpackInfoDTO(backpackItems [constants.BACKPACK_SIZE]*model.Item,
	itemType constants.EntityType) BackpackInfoDTO {

	items := make([]string, 0, 10)
	if itemType == constants.WEAPON {
		items = append(items, "unequip current weapon")
	}
	for _, item := range backpackItems {
		if item == nil {
			items = append(items, "empty slot")
		} else {
			items = append(items, item.Name)
		}
	}
	return BackpackInfoDTO{items: items}
}
