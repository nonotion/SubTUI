package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type Config struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Domain   string `yaml:"domain"`
}

var AppConfig Config

func main() {
	if err := initConfig(); err != nil {
		fmt.Fprintf(os.Stderr, "Config Error: %v\n", err)
		os.Exit(1)
	}

	if err := subsonicPing(); err != nil {
		fmt.Fprintf(os.Stderr, "Auth Error: %v\n", err)
		os.Exit(1)
	}

	if err := initPlayer(); err != nil {
		panic(err)
	}
	defer shutdownPlayer()

	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error while running program:", err)
		os.Exit(1)
	}
}
