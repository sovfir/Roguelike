package data

import (
	"encoding/json"
	"fmt"
	"os"
	"rogue/application/dto"
	"rogue/domain/model"
)

func LoadGame() (*model.GameSession, error) {
	var gsDTO dto.SaveGameDTO
	gsBytes, err := os.ReadFile(dto.FILENAME)
	if err != nil {
		return &model.GameSession{}, fmt.Errorf("reading file: %w", err)
	}
	err = json.Unmarshal(gsBytes, &gsDTO)
	if err != nil {
		return &model.GameSession{}, fmt.Errorf("unmarshall: %w", err)
	}

	return DTOtoGameSession(gsDTO), nil
}
