package service

import "rogue/domain/model"

func Unequip(gs *model.GameSession) {
	item := gs.Hero.Weapon
	if ok := gs.AddItemToBackpack(item); ok {
		gs.Hero.Weapon = nil
		gs.Message = gs.Message + item.Name + " moved to backpack"
	} else {
		gs.Message = gs.Message + " backpack is full"
	}

}
