package data

import (
	"encoding/json"
	"fmt"
	"os"
	"rogue/application/dto"
	"rogue/domain/model"
)

func SaveGame(gs *model.GameSession) error {
	if gsJson, err := json.MarshalIndent(GameSessionToDTO(gs), "", "\t"); err == nil {
		os.WriteFile(dto.FILENAME, gsJson, 0777)
		return nil
	} else {
		return fmt.Errorf("save game: %w", err)
	}
}
