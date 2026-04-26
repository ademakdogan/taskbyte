package app

import (
	"fmt"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/adem/taskbyte/internal/config"
	"github.com/adem/taskbyte/internal/db"
	"github.com/adem/taskbyte/internal/model"
	"github.com/adem/taskbyte/internal/service"
	"github.com/adem/taskbyte/internal/ui"
)

// Mode represents the current application mode.
type Mode int

const (
	ModeViewer Mode = iota
	ModeInsert
	ModeEdit
	ModeSearch
	ModeGoto
	ModeStats
	ModeSettings
)

// Model is the top-level Bubble Tea model.
type Model struct {
	mode        Mode
	svc         *service.TaskService
	cfg         config.Config
	styles      ui.Styles
	currentDate string // YYYY-MM-DD
	tasks       []model.Task
	cursor      int
	width       int
	height      int
	err         error

	// Insert/Edit mode
	inputValue   string
	inputCursor  int
	editTaskID   int
	inputHistory []string
	historyIdx   int

	// Search mode
	searchQuery   string
	searchResults []model.Task
	searchCursor  int

	// Delete confirmation
	deleteConfirm bool

	// Goto mode
	gotoStats  []db.DateStats
	gotoCursor int
	gotoInput  string

	// Stats mode
	statsData  ui.StatsData
	statsRange string

	// Settings mode
	settingsCursor    int
	settingsDropdown  bool
	settingsItems     []ui.SettingItem
	settingsOptCursor int
}

// New creates a new application model.
func New(svc *service.TaskService, cfg config.Config) Model {
	return Model{
		mode:        ModeViewer,
		svc:         svc,
		cfg:         cfg,
		styles:      ui.NewStyles(cfg),
		currentDate: service.TodayString(),
		historyIdx:  -1,
	}
}

// Init loads the initial data.
func (m Model) Init() tea.Cmd {
	return m.loadTasks()
}

type tasksLoadedMsg struct {
	tasks []model.Task
	err   error
}

func (m Model) loadTasks() tea.Cmd {
	return func() tea.Msg {
		tasks, err := m.svc.GetTasksForDate(m.currentDate)
		return tasksLoadedMsg{tasks: tasks, err: err}
	}
}

// Update handles all messages.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tasksLoadedMsg:
		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}
		m.tasks = msg.tasks
		if m.cursor >= len(m.tasks) {
			m.cursor = max(0, len(m.tasks)-1)
		}
		return m, nil

	case searchResultsMsg:
		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}
		m.searchResults = msg.results
		return m, nil

	case gotoStatsMsg:
		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}
		m.gotoStats = msg.stats
		return m, nil

	case statsDataMsg:
		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}
		m.statsData = ui.AggregateStats(msg.stats)
		return m, nil

	case tea.KeyMsg:
		// Global quit
		if msg.String() == "ctrl+c" || msg.String() == "q" && m.mode == ModeViewer {
			return m, tea.Quit
		}

		switch m.mode {
		case ModeViewer:
			return m.updateViewer(msg)
		case ModeInsert:
			return m.updateInsert(msg)
		case ModeEdit:
			return m.updateEdit(msg)
		case ModeSearch:
			return m.updateSearch(msg)
		case ModeGoto:
			return m.updateGoto(msg)
		case ModeStats:
			return m.updateStats(msg)
		case ModeSettings:
			return m.updateSettings(msg)
		}
	}

	return m, nil
}

// View renders the UI.
func (m Model) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	var s strings.Builder

	// Header
	dateDisplay := service.StorageToDisplay(m.currentDate, m.cfg.DateFormat)
	if service.IsToday(m.currentDate) {
		dateDisplay = "Today (" + dateDisplay + ")"
	}

	header := m.styles.Title.Render("TaskByte") + "  " + m.styles.Subtle.Render(dateDisplay)
	s.WriteString(header + "\n")
	s.WriteString(m.styles.Subtle.Render(strings.Repeat("─", min(m.width, 60))) + "\n\n")

	switch m.mode {
	case ModeViewer, ModeInsert, ModeEdit:
		s.WriteString(m.viewTaskList())
		if m.mode == ModeInsert || m.mode == ModeEdit {
			s.WriteString("\n")
			s.WriteString(m.viewInput())
		}
	case ModeSearch:
		s.WriteString(m.viewSearch())
	case ModeGoto:
		s.WriteString(ui.RenderGoto(m.gotoStats, m.gotoCursor, m.gotoInput, m.cfg.DateFormat, m.width, m.styles))
	case ModeStats:
		s.WriteString(ui.RenderStats(m.statsData, m.statsRange, m.styles))
	case ModeSettings:
		s.WriteString(ui.RenderSettings(m.settingsItems, m.settingsCursor, m.settingsDropdown, m.settingsOptCursor, m.styles))
	}

	// Error display
	if m.err != nil {
		s.WriteString("\n" + m.styles.ErrorStyle.Render("Error: "+m.err.Error()))
	}

	// Help bar
	s.WriteString("\n\n")
	s.WriteString(m.viewHelp())

	return s.String()
}

func (m Model) viewTaskList() string {
	if len(m.tasks) == 0 {
		return m.styles.Subtle.Render("  No tasks for this date. Press 'i' to add one.\n")
	}


	var s strings.Builder
	for i, task := range m.tasks {
		// Auto-hide completed/cancelled tasks if enabled
		if m.cfg.AutoHideCompleted && (task.Status == model.StatusDone || task.Status == model.StatusCancelled) {
			continue
		}

		cursor := "  "
		if i == m.cursor && m.mode == ModeViewer {
			cursor = m.styles.FocusedItem.Render("> ")
		}

		line := m.renderTask(task, i == m.cursor && m.mode == ModeViewer)

		if m.deleteConfirm && i == m.cursor {
			line += m.styles.ErrorStyle.Render(" [Press 'r' again to confirm delete]")
		}

		s.WriteString(cursor + line + "\n")
	}
	return s.String()
}

func (m Model) renderTask(task model.Task, focused bool) string {
	symbol := task.Status.Symbol()
	label := task.Status.Label()

	var style lipgloss.Style
	switch task.Status {
	case model.StatusTodo:
		style = m.styles.TodoStyle
	case model.StatusInProgress:
		style = m.styles.InProgressStyle
	case model.StatusDone:
		style = m.styles.DoneStyle
	case model.StatusCancelled:
		style = m.styles.CancelledStyle
	}

	text := fmt.Sprintf("%s %s", symbol, task.Text)
	if label != "" {
		text += " " + label
	}

	// Add status timestamp
	if task.StatusChangedAt != nil && task.Status != model.StatusTodo {
		ts := task.StatusChangedAt.Local().Format("02.01.2006 15:04")
		text += style.Render(" - " + ts)
	}

	if focused {
		return m.styles.FocusedItem.Render(symbol) + " " + style.Render(task.Text) + func() string {
			if label != "" {
				return " " + style.Render(label)
			}
			return ""
		}() + func() string {
			if task.StatusChangedAt != nil && task.Status != model.StatusTodo {
				ts := task.StatusChangedAt.Local().Format("02.01.2006 15:04")
				return " " + style.Render("- "+ts)
			}
			return ""
		}()
	}

	return style.Render(text)
}

func (m Model) viewInput() string {
	var modeLabel string
	switch m.mode {
	case ModeInsert:
		modeLabel = "—insert—"
	case ModeEdit:
		modeLabel = "—edit—"
	}

	dateDisplay := service.StorageToDisplay(m.currentDate, m.cfg.DateFormat)
	if service.IsToday(m.currentDate) {
		dateDisplay = "Today"
	}

	top := m.styles.Subtle.Render("┌─ " + dateDisplay + " " + strings.Repeat("─", max(0, 40-len(dateDisplay))) + "┐")

	// Ghost text for slash commands
	ghostText := m.getGhostText()
	displayInput := m.inputValue + "█"
	if ghostText != "" {
		displayInput = m.inputValue + m.styles.Subtle.Render(ghostText) + "█"
	}

	input := "│ › " + displayInput + strings.Repeat(" ", max(0, 37-len(m.inputValue)-len(ghostText))) + " │"
	bottom := m.styles.Subtle.Render("└" + strings.Repeat("─", max(0, 34-len(modeLabel))) + " " + modeLabel + " ─┘")

	return top + "\n" + input + "\n" + bottom
}

func (m Model) getGhostText() string {
	if m.inputValue == "" {
		return ""
	}

	commands := map[string]string{
		"/":            "date | sort | export | hide | show | migrate | migrate-all | stats | settings",
		"/d":           "ate DD.MM.YYYY",
		"/da":          "te DD.MM.YYYY",
		"/dat":         "e DD.MM.YYYY",
		"/date":        " DD.MM.YYYY",
		"/s":           "ort | stats | settings",
		"/so":          "rt [date/progress/date-reverse/progress-reverse]",
		"/sor":         "t [date/progress/date-reverse/progress-reverse]",
		"/sort":        " [date/progress/date-reverse/progress-reverse]",
		"/st":          "ats [day/week/month/all]",
		"/sta":         "ts [day/week/month/all]",
		"/stat":        "s [day/week/month/all]",
		"/stats":       " [day/week/month/all]",
		"/se":          "ttings",
		"/set":         "tings",
		"/sett":        "ings",
		"/e":           "xport [Path]",
		"/ex":          "port [Path]",
		"/exp":         "ort [Path]",
		"/expo":        "rt [Path]",
		"/expor":       "t [Path]",
		"/export":      " [Path]",
		"/h":           "ide",
		"/hi":          "de",
		"/hid":         "e",
		"/sh":          "ow",
		"/sho":         "w",
		"/m":           "igrate [DD.MM.YYYY]",
		"/mi":          "grate [DD.MM.YYYY]",
		"/mig":         "rate [DD.MM.YYYY]",
		"/migr":        "ate [DD.MM.YYYY]",
		"/migra":       "te [DD.MM.YYYY]",
		"/migrat":      "e [DD.MM.YYYY]",
		"/migrate":     " [DD.MM.YYYY]",
		"/migrate-":    "all",
		"/migrate-a":   "ll",
		"/migrate-al":  "l",
	}

	if ghost, ok := commands[m.inputValue]; ok {
		return ghost
	}
	return ""
}

func (m Model) viewSearch() string {
	var s strings.Builder

	top := m.styles.Subtle.Render("┌─ Search " + strings.Repeat("─", 34) + "┐")
	input := "│ › " + m.searchQuery + "█" + strings.Repeat(" ", max(0, 37-len(m.searchQuery))) + " │"
	bottom := m.styles.Subtle.Render("└" + strings.Repeat("─", 26) + " —search— ─┘")
	s.WriteString(top + "\n" + input + "\n" + bottom + "\n\n")

	if len(m.searchResults) == 0 && m.searchQuery != "" {
		s.WriteString(m.styles.Subtle.Render("  No results found.\n"))
	} else {
		for i, task := range m.searchResults {
			cursor := "  "
			if i == m.searchCursor {
				cursor = m.styles.FocusedItem.Render("> ")
			}
			dateStr := service.StorageToDisplay(task.Date, m.cfg.DateFormat)
			line := m.renderTask(task, i == m.searchCursor)
			s.WriteString(cursor + m.styles.Subtle.Render("["+dateStr+"] ") + line + "\n")
		}
	}

	return s.String()
}

func (m Model) viewHelp() string {
	var help string
	switch m.mode {
	case ModeViewer:
		help = "[i]nsert  [e]dit  [s]earch  [g]oto  [r]emove  [Enter] cycle status  [/] commands  [q]uit"
	case ModeInsert:
		help = "[Enter] add task  [Esc] back  [/] commands  [↑/↓] history"
	case ModeEdit:
		help = "[Enter] save  [Esc] cancel"
	case ModeSearch:
		help = "[Enter] go to task  [Esc] back  [↑/↓] navigate  [p/d/x/space] status"
	case ModeGoto:
		help = "[↑/↓]: Navigate  |  [Enter]: Open List  |  [i]: Insert  |  [Esc]: Back"
	case ModeStats:
		help = "[Esc]: Back"
	case ModeSettings:
		help = "[↑/↓]: Navigate  |  [Enter]: Change  |  [Esc]: Return"
	}
	return m.styles.HelpStyle.Render(help)
}

// --- Viewer Mode ---

func (m Model) updateViewer(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case ui.Key(msg, "up", "k"):
		if m.cursor > 0 {
			m.cursor--
		}
		m.deleteConfirm = false

	case ui.Key(msg, "down", "j"):
		if m.cursor < len(m.tasks)-1 {
			m.cursor++
		}
		m.deleteConfirm = false

	case ui.Key(msg, "i"):
		m.mode = ModeInsert
		m.inputValue = ""
		m.deleteConfirm = false

	case ui.Key(msg, "e"):
		if len(m.tasks) > 0 {
			m.mode = ModeEdit
			m.editTaskID = m.tasks[m.cursor].ID
			m.inputValue = m.tasks[m.cursor].Text
		}

	case ui.Key(msg, "s"):
		m.mode = ModeSearch
		m.searchQuery = ""
		m.searchResults = nil
		m.searchCursor = 0

	case ui.Key(msg, "g"):
		m.mode = ModeGoto
		m.gotoInput = ""
		m.gotoCursor = 0
		return m, m.loadGotoStats()

	case ui.Key(msg, "/"):
		// Enter insert mode with slash pre-filled
		m.mode = ModeInsert
		m.inputValue = "/"
		m.deleteConfirm = false

	case ui.Key(msg, "r"):
		if len(m.tasks) > 0 {
			if m.deleteConfirm {
				m.svc.DeleteTask(m.tasks[m.cursor].ID)
				m.deleteConfirm = false
				return m, m.loadTasks()
			}
			m.deleteConfirm = true
		}

	case ui.Key(msg, "enter"):
		if len(m.tasks) > 0 {
			task := m.tasks[m.cursor]
			m.svc.CycleStatus(task.ID)
			return m, m.loadTasks()
		}

	case ui.Key(msg, "p"):
		if len(m.tasks) > 0 {
			m.svc.SetStatus(m.tasks[m.cursor].ID, model.StatusInProgress)
			return m, m.loadTasks()
		}

	case ui.Key(msg, "d"):
		if len(m.tasks) > 0 {
			m.svc.SetStatus(m.tasks[m.cursor].ID, model.StatusDone)
			return m, m.loadTasks()
		}

	case ui.Key(msg, "x"):
		if len(m.tasks) > 0 {
			m.svc.SetStatus(m.tasks[m.cursor].ID, model.StatusCancelled)
			return m, m.loadTasks()
		}

	case ui.Key(msg, " "):
		if len(m.tasks) > 0 {
			m.svc.SetStatus(m.tasks[m.cursor].ID, model.StatusTodo)
			return m, m.loadTasks()
		}
	}

	return m, nil
}

// --- Insert Mode ---

func (m Model) updateInsert(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case ui.Key(msg, "esc"):
		m.mode = ModeViewer
		m.inputValue = ""

	case ui.Key(msg, "enter"):
		if m.inputValue != "" {
			// Check for slash commands
			if strings.HasPrefix(m.inputValue, "/") {
				return m.handleSlashCommand(m.inputValue)
			}

			m.svc.AddTask(m.inputValue, m.currentDate)
			if m.cfg.InsertPromptHistory {
				m.inputHistory = append(m.inputHistory, m.inputValue)
			}
			m.inputValue = ""
			m.historyIdx = -1
			return m, m.loadTasks()
		}

	case ui.Key(msg, "up"):
		if m.cfg.InsertPromptHistory && len(m.inputHistory) > 0 {
			if m.historyIdx < len(m.inputHistory)-1 {
				m.historyIdx++
				m.inputValue = m.inputHistory[len(m.inputHistory)-1-m.historyIdx]
			}
		}

	case ui.Key(msg, "down"):
		if m.cfg.InsertPromptHistory {
			if m.historyIdx > 0 {
				m.historyIdx--
				m.inputValue = m.inputHistory[len(m.inputHistory)-1-m.historyIdx]
			} else {
				m.historyIdx = -1
				m.inputValue = ""
			}
		}

	case ui.Key(msg, "backspace"):
		if len(m.inputValue) > 0 {
			m.inputValue = m.inputValue[:len(m.inputValue)-1]
		}

	default:
		if msg.Type == tea.KeyRunes {
			m.inputValue += string(msg.Runes)
			m.historyIdx = -1
		}
	}

	return m, nil
}

// --- Edit Mode ---

func (m Model) updateEdit(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case ui.Key(msg, "esc"):
		m.mode = ModeViewer
		m.inputValue = ""

	case ui.Key(msg, "enter"):
		if m.inputValue != "" {
			m.svc.EditTask(m.editTaskID, m.inputValue)
			m.mode = ModeViewer
			m.inputValue = ""
			return m, m.loadTasks()
		}

	case ui.Key(msg, "backspace"):
		if len(m.inputValue) > 0 {
			m.inputValue = m.inputValue[:len(m.inputValue)-1]
		}

	default:
		if msg.Type == tea.KeyRunes {
			m.inputValue += string(msg.Runes)
		}
	}

	return m, nil
}

// --- Search Mode ---

func (m Model) updateSearch(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case ui.Key(msg, "esc"):
		m.mode = ModeViewer
		m.searchQuery = ""
		m.searchResults = nil

	case ui.Key(msg, "enter"):
		if len(m.searchResults) > 0 {
			task := m.searchResults[m.searchCursor]
			m.currentDate = task.Date
			m.mode = ModeViewer
			m.searchQuery = ""
			m.searchResults = nil
			return m, m.loadTasks()
		}

	case ui.Key(msg, "up"):
		if m.searchCursor > 0 {
			m.searchCursor--
		}

	case ui.Key(msg, "down"):
		if m.searchCursor < len(m.searchResults)-1 {
			m.searchCursor++
		}

	case ui.Key(msg, "p"):
		if len(m.searchResults) > 0 {
			m.svc.SetStatus(m.searchResults[m.searchCursor].ID, model.StatusInProgress)
			return m, m.doSearch()
		}

	case ui.Key(msg, "backspace"):
		if len(m.searchQuery) > 0 {
			m.searchQuery = m.searchQuery[:len(m.searchQuery)-1]
			m.searchCursor = 0
			return m, m.doSearch()
		}

	default:
		if msg.Type == tea.KeyRunes {
			m.searchQuery += string(msg.Runes)
			m.searchCursor = 0
			return m, m.doSearch()
		}
	}

	return m, nil
}

type searchResultsMsg struct {
	results []model.Task
	err     error
}

func (m Model) doSearch() tea.Cmd {
	return func() tea.Msg {
		results, err := m.svc.SearchTasks(m.searchQuery)
		return searchResultsMsg{results: results, err: err}
	}
}

func (m Model) handleSlashCommand(input string) (tea.Model, tea.Cmd) {
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return m, nil
	}

	cmd := strings.ToLower(parts[0])
	switch cmd {
	case "/date":
		if len(parts) > 1 {
			storageDate, err := service.FormatDate(parts[1], m.cfg.DateFormat)
			if err == nil {
				m.currentDate = storageDate
				m.inputValue = ""
				return m, m.loadTasks()
			}
			m.err = err
		}

	case "/hide":
		m.cfg.AutoHideCompleted = true
		m.inputValue = ""
		return m, m.loadTasks()

	case "/show":
		m.cfg.AutoHideCompleted = false
		m.inputValue = ""
		return m, m.loadTasks()

	case "/migrate":
		if len(parts) > 1 {
			storageDate, err := service.FormatDate(parts[1], m.cfg.DateFormat)
			if err == nil {
				m.svc.MigrateTasks(storageDate)
				m.inputValue = ""
				return m, m.loadTasks()
			}
		} else {
			m.svc.MigrateTasks(service.YesterdayString())
			m.inputValue = ""
			return m, m.loadTasks()
		}

	case "/migrate-all":
		m.svc.MigrateAllTasks()
		m.inputValue = ""
		return m, m.loadTasks()

	case "/stats":
		rangeLabel := "all"
		if len(parts) > 1 {
			rangeLabel = parts[1]
		}
		m.statsRange = rangeLabel
		m.mode = ModeStats
		m.inputValue = ""
		return m, m.loadStatsData()

	case "/settings":
		m.mode = ModeSettings
		m.settingsItems = ui.BuildSettingsItems(m.cfg)
		m.settingsCursor = 0
		m.settingsDropdown = false
		m.inputValue = ""

	case "/sort":
		if len(parts) > 1 {
			sortType := parts[1]
			switch sortType {
			case "date":
				sort.Slice(m.tasks, func(i, j int) bool {
					return m.tasks[i].CreatedAt.Before(m.tasks[j].CreatedAt)
				})
			case "date-reverse":
				sort.Slice(m.tasks, func(i, j int) bool {
					return m.tasks[i].CreatedAt.After(m.tasks[j].CreatedAt)
				})
			case "progress":
				statusOrder := map[model.Status]int{
					model.StatusInProgress: 0,
					model.StatusTodo:       1,
					model.StatusDone:       2,
					model.StatusCancelled:  3,
				}
				sort.Slice(m.tasks, func(i, j int) bool {
					return statusOrder[m.tasks[i].Status] < statusOrder[m.tasks[j].Status]
				})
			case "progress-reverse":
				statusOrder := map[model.Status]int{
					model.StatusInProgress: 0,
					model.StatusTodo:       1,
					model.StatusDone:       2,
					model.StatusCancelled:  3,
				}
				sort.Slice(m.tasks, func(i, j int) bool {
					return statusOrder[m.tasks[i].Status] > statusOrder[m.tasks[j].Status]
				})
			}
		}
		m.inputValue = ""

	case "/export":
		exportPath := ""
		if len(parts) > 1 {
			exportPath = strings.Join(parts[1:], " ")
		}
		path, err := ui.ExportToMarkdown(m.tasks, m.currentDate, m.cfg.DateFormat, exportPath)
		if err != nil {
			m.err = err
		} else {
			m.err = fmt.Errorf("exported to: %s", path) // show as info
		}
		m.inputValue = ""
	}

	m.inputValue = ""
	return m, nil
}

// --- Goto Mode ---

type gotoStatsMsg struct {
	stats []db.DateStats
	err   error
}

func (m Model) loadGotoStats() tea.Cmd {
	return func() tea.Msg {
		stats, err := m.svc.GetDateStats()
		return gotoStatsMsg{stats: stats, err: err}
	}
}

func (m Model) updateGoto(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case ui.Key(msg, "esc"):
		m.mode = ModeViewer
		m.gotoInput = ""

	case ui.Key(msg, "up"):
		if m.gotoCursor > 0 {
			m.gotoCursor--
		}

	case ui.Key(msg, "down"):
		if m.gotoCursor < len(m.gotoStats)-1 {
			m.gotoCursor++
		}

	case ui.Key(msg, "enter"):
		if len(m.gotoStats) > 0 {
			m.currentDate = m.gotoStats[m.gotoCursor].Date
			m.mode = ModeViewer
			m.gotoInput = ""
			return m, m.loadTasks()
		}
		// If we have date input, try to parse and jump
		if m.gotoInput != "" {
			storageDate, err := service.FormatDate(m.gotoInput, m.cfg.DateFormat)
			if err == nil {
				m.currentDate = storageDate
				m.mode = ModeViewer
				m.gotoInput = ""
				return m, m.loadTasks()
			}
			m.err = err
		}

	case ui.Key(msg, "i"):
		m.mode = ModeInsert
		m.inputValue = ""

	case ui.Key(msg, "backspace"):
		if len(m.gotoInput) > 0 {
			m.gotoInput = m.gotoInput[:len(m.gotoInput)-1]
		}

	default:
		if msg.Type == tea.KeyRunes {
			m.gotoInput += string(msg.Runes)
			// Try to find matching date and focus
			for i, st := range m.gotoStats {
				display := service.StorageToDisplay(st.Date, m.cfg.DateFormat)
				if strings.HasPrefix(display, m.gotoInput) {
					m.gotoCursor = i
					break
				}
			}
		}
	}

	return m, nil
}

// --- Stats Mode ---

type statsDataMsg struct {
	stats []db.DateStats
	err   error
}

func (m Model) loadStatsData() tea.Cmd {
	statsRange := m.statsRange
	return func() tea.Msg {
		stats, err := m.svc.GetDateStats()
		if err != nil {
			return statsDataMsg{stats: nil, err: err}
		}
		filtered := service.FilterStatsByRange(stats, statsRange)
		return statsDataMsg{stats: filtered, err: nil}
	}
}

func (m Model) updateStats(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if ui.Key(msg, "esc") {
		m.mode = ModeViewer
	}
	return m, nil
}

// --- Settings Mode ---

func (m Model) updateSettings(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.settingsDropdown {
		return m.updateSettingsDropdown(msg)
	}

	switch {
	case ui.Key(msg, "esc"):
		// Save and exit
		newCfg := ui.ApplySettingsItems(m.settingsItems)
		config.Save(newCfg)
		m.cfg = newCfg
		m.styles = ui.NewStyles(newCfg)
		m.mode = ModeViewer

	case ui.Key(msg, "up"):
		if m.settingsCursor > 0 {
			m.settingsCursor--
		}

	case ui.Key(msg, "down"):
		if m.settingsCursor < len(m.settingsItems)-1 {
			m.settingsCursor++
		}

	case ui.Key(msg, "enter"):
		item := m.settingsItems[m.settingsCursor]
		if item.Type == "bool" {
			if item.Value == "True" {
				m.settingsItems[m.settingsCursor].Value = "False"
			} else {
				m.settingsItems[m.settingsCursor].Value = "True"
			}
		} else if item.Type == "select" {
			m.settingsDropdown = true
			m.settingsOptCursor = 0
			// Find current value in options
			for i, opt := range item.Options {
				if opt == item.Value {
					m.settingsOptCursor = i
					break
				}
			}
		}
	}

	return m, nil
}

func (m Model) updateSettingsDropdown(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	item := m.settingsItems[m.settingsCursor]

	switch {
	case ui.Key(msg, "esc"):
		m.settingsDropdown = false

	case ui.Key(msg, "up"):
		if m.settingsOptCursor > 0 {
			m.settingsOptCursor--
		}

	case ui.Key(msg, "down"):
		if m.settingsOptCursor < len(item.Options)-1 {
			m.settingsOptCursor++
		}

	case ui.Key(msg, "enter"):
		m.settingsItems[m.settingsCursor].Value = item.Options[m.settingsOptCursor]
		m.settingsDropdown = false
	}

	return m, nil
}
