package main

import (
	"rogue/application/controller"
	"rogue/application/dto"
	"rogue/front"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// session := domain.NewGameSession()

	// канал ввода с клавиатуры
	inputCh := make(chan string)
	defer close(inputCh)
	updateCh := make(chan dto.DomainToViewDTO)
	itemCh := make(chan dto.BackpackDTO)

	// Запускаем горутину для обработки игровой логики
	go controller.UseCasesController(inputCh, updateCh, itemCh)
	// Получаем первое DTO
	firstDto := <-updateCh

	// Создаем начальную модель и применяем обновление
	initialModel := front.InitialModel(front.Background, inputCh, updateCh, itemCh)
	updatedModel := initialModel.ProcessUpdate(firstDto)

	// Запускаем UI с уже обновленной моделью
	p := tea.NewProgram(updatedModel, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		panic(err)
	}
}
