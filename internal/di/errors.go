package di

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type ErrorHandler struct{}

func NewErrorHandler() *ErrorHandler {
	return &ErrorHandler{}
}

func (h *ErrorHandler) HandleConfigError(err error) error {
	if err == nil {
		return nil
	}
	if strings.Contains(err.Error(), "PUBLIC_BFF_URL") {
		return h.formatEnvironmentError(err)
	}
	return h.formatGenericError("Configuration", err)
}

func (h *ErrorHandler) HandleAPIError(err error) error {
	if err == nil {
		return nil
	}
	return h.formatGenericError("API Client", err)
}

func (h *ErrorHandler) formatEnvironmentError(err error) error {
	errorStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("1")).
		Border(lipgloss.RoundedBorder()).
		Padding(1, 2)

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("3")).
		MarginBottom(1)

	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("8")).
		MarginTop(1)

	title := titleStyle.Render("Configuration Error")
	message := errorStyle.Render(err.Error())
	help := helpStyle.Render("Set PUBLIC_BFF_URL (and optional PROJECT_ID/AUTH_HEADER) in environment or .env file.")

	return fmt.Errorf("%s\n\n%s\n\n%s", title, message, help)
}

func (h *ErrorHandler) formatGenericError(service string, err error) error {
	errorStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("1")).
		Border(lipgloss.RoundedBorder()).
		Padding(1, 2)

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("3")).
		MarginBottom(1)

	title := titleStyle.Render(fmt.Sprintf("%s Error", service))
	message := errorStyle.Render(err.Error())

	return fmt.Errorf("%s\n\n%s", title, message)
}
