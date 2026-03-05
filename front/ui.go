package front

import (
	"fmt"
	"rogue/application/dto"
	"rogue/infrastructure/constants"
	"strings"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	ScreenWidth  = 80
	ScreenHeight = 24
	GameWidth    = 78
	GameHeight   = 21
	PanelHeight  = 1
)

type Model struct {
	grid        [GameHeight][GameWidth]string
	msgPanel    string
	statusPanel string
	menuPanel   string
	inputCh     chan<- string
	updateCh    <-chan dto.DomainToViewDTO
	inventoryCh chan dto.BackpackDTO
	gameStarted bool
	showInventory bool
	inventory dto.BackpackDTO
}

// инициализируем ее
func InitialModel(background [GameHeight][GameWidth]string, inputCh chan<- string, updateCh <-chan dto.DomainToViewDTO, inventory chan dto.BackpackDTO) Model {
	m := Model{}

	m.grid = background
	m.msgPanel = "arrowcle/gjacinta for School-21"
	m.statusPanel = "Press any key to start"
	m.menuPanel = "[Q]Выход"
	m.inputCh = inputCh
	m.updateCh = updateCh
	m.gameStarted = false
	m.showInventory = false
	m.inventoryCh = inventory
	return m
}

// Обязательный метод для tea.Model
func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) ProcessUpdate(dto dto.DomainToViewDTO) Model {
	if dto.GameStatus() == -1 {
		m.msgPanel = "Game Over"
		m.grid = looseBackground
	
	}
	if dto.GameStatus() == 1 {
		m.msgPanel = "VICTORY!"
		m.grid = winBackground
	}
	if dto.GameStatus() == 0 {
		// Обновляем сообщение
	m.msgPanel = dto.Message()

	// Обновляем статус с информацией о герое
	hero := dto.HeroInfo()
	m.statusPanel = fmt.Sprintf("Lvl:%d|GLD:%d|❤:%d/%d|Str:%d|Agl:%d|⚔:%s",
		dto.Level(), hero.Gold, hero.Health, hero.MaxHealth, hero.Strength, hero.Agility,dto.Weapon())

	// Преобразуем клетки в grid с учетом GroundType
	m.updateGrid(dto.FieldInfo(), hero)


	}
	return m
}

func (m *Model) updateGrid(cells []dto.CellInfoDTO, hero dto.HeroInfoDTO) {
	// Сначала очищаем grid
	for i := range m.grid {
		for j := range m.grid[i] {
			m.grid[i][j] = " "
		}
	}

	// Маппинг GroundType на символы для отображения
	groundSymbols := map[constants.GroundType]string{
		constants.WALL:     "█",
		constants.FLOOR:    "░",
		constants.CORRIDOR: "░",
		constants.PASSAGE:  "▒",
		constants.EXIT:     "E",
	}

	// Стили для дверей и ключей
	redStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000"))    // Красный
	blueStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#0000FF"))   // Синий
	goldStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#ffea00ff")) // Золото
	greenStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#00ff55ff")) // Зелень
	whiteStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#ffffffff")) // Зелень
	purpleStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#d000ffff")) // Зелень
	// Заполняем клетками из DTO
	for _, cell := range cells {
		if cell.Row < GameHeight && cell.Col < GameWidth {
			if cell.Entity == constants.NONE {
				if symbol, exists := groundSymbols[cell.Ground]; exists {
					m.grid[cell.Row][cell.Col] = symbol
				} else {
					m.grid[cell.Row][cell.Col] = "?"
				}
			} else {
				switch cell.Entity {
				case constants.FIRST_DOOR:
					m.grid[cell.Row][cell.Col] = redStyle.Render("▒") // Красная дверь
				case constants.FIRST_KEY:
					m.grid[cell.Row][cell.Col] = redStyle.Render("✜") // Красный ключ
				case constants.SECOND_DOOR:
					m.grid[cell.Row][cell.Col] = blueStyle.Render("▒") // Синяя дверь
				case constants.SECOND_KEY:
					m.grid[cell.Row][cell.Col] = blueStyle.Render("✜") // Синий ключ
				case constants.FOOD:
					m.grid[cell.Row][cell.Col] = purpleStyle.Render("") 
				case constants.SCROLL:
					m.grid[cell.Row][cell.Col] = purpleStyle.Render("🜏") 
				case constants.WEAPON:
					m.grid[cell.Row][cell.Col] = goldStyle.Render("⸸")
				case constants.GHOST:
					m.grid[cell.Row][cell.Col] = whiteStyle.Render("G")
				case constants.OGRE:
					m.grid[cell.Row][cell.Col] = goldStyle.Render("O")
				case constants.SNAKE:
					m.grid[cell.Row][cell.Col] = whiteStyle.Render("S")
				case constants.VAMPIRE:
					m.grid[cell.Row][cell.Col] = redStyle.Render("V")
				case constants.ZOMBIE:
					m.grid[cell.Row][cell.Col] = greenStyle.Render("Z")
				case constants.MIMIC:
					m.grid[cell.Row][cell.Col] = whiteStyle.Render("M")
				case constants.ELIXIR:
					m.grid[cell.Row][cell.Col] = purpleStyle.Render("🜬")
				case constants.TREASURE:
					m.grid[cell.Row][cell.Col] = purpleStyle.Render("◛")
				}

			}
		}
	}

	// Добавляем героя поверх карты
	if hero.Row < GameHeight && hero.Col < GameWidth {
		m.grid[hero.Row][hero.Col] = goldStyle.Render("◙")
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Игра не запущена
	if !m.gameStarted {
		m.grid = Background
		m.msgPanel = "arrowcle/gjacinta for School-21"
		m.statusPanel = "Press any key to start"
		m.menuPanel = "[Q]Quit [L]Load [P]Save [WASD]Move [K]Potions [H]Weapons [E]Scrolls [J]Food [O]Stats [U]Leaderborad"
		   // Если игра не запущена и нажата любая клавиша - начинаем игру
        if _, isKeyMsg := msg.(tea.KeyMsg); isKeyMsg {
            m.gameStarted = true
        }
	}
	
	// Игра запущена
	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := msg.String()
		if key == "q" || key == "Q" {
			return m, tea.Quit
		}

		m.inputCh <- key
		m.showInventory = false
		select {
		case updateData := <-m.updateCh:
			return m.ProcessUpdate(updateData), nil
		case itemsInfo := <-m.inventoryCh:
			//вывод менб с предметами
			m.inventory = itemsInfo
			m.showInventory = true
		}

	default:
	}

	return m, nil
}

// Метод отображения
func (m Model) View() string {
	var gridLines []string
	for row := 0; row < GameHeight; row++ {
		var line string
		for col := 0; col < GameWidth; col++ {
			line += m.grid[row][col]
		}
		gridLines = append(gridLines, line)
	}
	gameArea := strings.Join(gridLines, "\n")
	//стиль модального окна
	inventoryStyle := lipgloss.NewStyle().
            Width(60).
            Height(15).
            Padding(1, 2).
            Border(lipgloss.RoundedBorder()).
            BorderForeground(lipgloss.Color("62")).
            Background(lipgloss.Color("0")).
            Foreground(lipgloss.Color("15"))

	// Разные стили для каждой панели
	msgStyle := lipgloss.NewStyle().
		Width(ScreenWidth).
		Background(lipgloss.Color("#3E3E3E")).
		Foreground(lipgloss.Color("#E8B4A5"))

	statusStyle := lipgloss.NewStyle().
		Width(ScreenWidth).
		Background(lipgloss.Color("#F5E9DE")).
		Foreground(lipgloss.Color("#5A4C3C"))

	menuStyle := lipgloss.NewStyle().
		Width(ScreenWidth).
		Background(lipgloss.Color("#6A5D52")).
		Foreground(lipgloss.Color("#FFFFFF"))

	if m.showInventory {
   items := m.inventory.Items()
var lines []string
lines = append(lines, "INVENTORY:")
lines = append(lines, "──────────")

startFromZero := len(items) == 10

for i, item := range items {
    if startFromZero {
        lines = append(lines, fmt.Sprintf("%d. %s", i, item))
    } else {
        lines = append(lines, fmt.Sprintf("%d. %s", i+1, item))
    }
}

modalContent := strings.Join(lines, "\n")
    
    return lipgloss.Place(
        80, 30,
        lipgloss.Center, lipgloss.Center,
        inventoryStyle.Render(modalContent),
    )
}

	return lipgloss.JoinVertical(lipgloss.Left,
		gameArea,
		msgStyle.Render(m.msgPanel),
		statusStyle.Render(m.statusPanel),
		menuStyle.Render(m.menuPanel),
	)
}
