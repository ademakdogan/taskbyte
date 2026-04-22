package ui

import (
"fmt"
"strings"

"github.com/charmbracelet/lipgloss"

"github.com/adem/taskbyte/internal/db"
"github.com/adem/taskbyte/internal/service"
)

// RenderGoto renders the goto/calendar view.
func RenderGoto(stats []db.DateStats, cursor int, input string, dateFormat string, width int, styles Styles) string {
	var s strings.Builder

	title := "=================== CALENDAR & GOTO ==================="
	s.WriteString(styles.Title.Render(title) + "\n\n")

	// Date input
	s.WriteString("  Type to jump (" + dateFormat + ") : " + input + "█\n\n")

	// Table header
	headerFmt := "    %-14s | %4s | %4s | %6s | %4s | %4s"
	header := fmt.Sprintf(headerFmt, "Date", "Tot", "Todo", "InProg", "Done", "Canc")
	s.WriteString(styles.Subtle.Render(header) + "\n")
	s.WriteString(styles.Subtle.Render("    "+strings.Repeat("─", min(width-8, 56))) + "\n")

	if len(stats) == 0 {
		s.WriteString(styles.Subtle.Render("    No data available.\n"))
	} else {
		for i, st := range stats {
			dateDisplay := service.StorageToDisplay(st.Date, dateFormat)

			todoStr := coloredNumber(st.Todo, styles.TodoStyle)
			inProgStr := coloredNumber(st.InProgress, styles.InProgressStyle)
			doneStr := coloredNumber(st.Done, styles.DoneStyle)
			cancStr := coloredNumber(st.Cancelled, styles.CancelledStyle)

			prefix := "    "
			suffix := ""
			if i == cursor {
				prefix = styles.FocusedItem.Render("  > ")
				suffix = styles.FocusedItem.Render(" <")
			}

			line := fmt.Sprintf("%-14s | %4d | %s | %s | %s | %s",
dateDisplay, st.Total, todoStr, inProgStr, doneStr, cancStr)

			if i == cursor {
				line = styles.Highlight.Render(fmt.Sprintf("%-14s", dateDisplay)) +
					fmt.Sprintf(" | %4d | %s | %s | %s | %s",
st.Total, todoStr, inProgStr, doneStr, cancStr)
			}

			s.WriteString(prefix + line + suffix + "\n")
		}
	}

	s.WriteString("\n" + styles.Subtle.Render("    "+strings.Repeat("=", min(width-8, 56))) + "\n")
	help := "[↑/↓]: Navigate  |  [Enter]: Open List  |  [i]: Insert  |  [Esc]: Back"
	s.WriteString(styles.HelpStyle.Render(help))

	return s.String()
}

func coloredNumber(n int, style lipgloss.Style) string {
	return style.Render(fmt.Sprintf("%4d", n))
}
