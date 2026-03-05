package service

import "rogue/domain/model"

func DecreaseEffectTime(gs *model.GameSession) {
	newEffects := make(map[int]model.Effect)
	i := 0
	for _, effect := range gs.Effects {
		effect.TimeLeft--
		if effect.TimeLeft > 0 {
			newEffects[i] = effect
			i++
		} else {
			gs.Hero.DecreaseStats(effect.Item)
		}
	}
	gs.Effects = newEffects
}
