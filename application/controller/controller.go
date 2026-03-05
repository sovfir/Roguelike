package controller

import (
	"fmt"
	"log"
	"rogue/application/dto"
	"rogue/application/usecases"
	"rogue/data"
	"rogue/domain/model"
	"rogue/infrastructure/constants"
)

func UseCasesController(
	inputCh <-chan string,
	updateGameInfoCh chan<- dto.DomainToViewDTO,
	backpackInfoCh chan<- dto.BackpackDTO) {

	gs := model.NewGameSession()

	// Загружаем статистику при старте
	if err := gs.LoadStats(); err != nil {
		gs.Message = "Failed to load stats: " + err.Error()
	}

	for {
		updateGameInfoCh <- dto.GetGameInfo(gs)
		signal := <-inputCh
		if signal == "Esc" {
			close(updateGameInfoCh)
			break
		}
		gs.Message = ""
		var err error
		switch signal {
		case "W", "w":
			err = usecases.MoveHero(constants.Up(), gs)
		case "S", "s":
			err = usecases.MoveHero(constants.Down(), gs)
		case "A", "a":
			err = usecases.MoveHero(constants.Left(), gs)
		case "D", "d":
			err = usecases.MoveHero(constants.Right(), gs)
		case "H", "h":
			usecases.BackpackInfo(constants.WEAPON, gs, inputCh, backpackInfoCh)
		case "J", "j":
			usecases.BackpackInfo(constants.FOOD, gs, inputCh, backpackInfoCh)
		case "K", "k":
			usecases.BackpackInfo(constants.ELIXIR, gs, inputCh, backpackInfoCh)
		case "E", "e":
			usecases.BackpackInfo(constants.SCROLL, gs, inputCh, backpackInfoCh)
		case "P", "p":
			if errSave := data.SaveGame(gs); errSave != nil {
				gs.Message = fmt.Sprintf("save game error: %s", errSave)
			}
			gs.Message += " Game saved"
		case "L", "l":
			if loadedGamesession, errSave := data.LoadGame(); errSave != nil {
				gs.Message = fmt.Sprintf("load game error: %s", errSave)
			} else {
				gs = loadedGamesession
			}
		case "O", "o":
			// Выводим текущую статистику в message
			gs.Message = gs.FormatCurrentStats()
		case "u", "U":
			// Выводим leaderboard в message
			gs.Message = gs.GetLeaderboard()
		}

		if err != nil {
			close(updateGameInfoCh)
			log.Fatalf("move hero error: %s", err)
		}
	}
}
