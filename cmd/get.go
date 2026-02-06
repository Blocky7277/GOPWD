package cmd

import (
	"os"
	// "golang.org/x/term"
	"github.com/sahilm/fuzzy"
	"internal/cryptoutil"
	"internal/util"
	"fmt"
	// "syscall"
	// "strings"
	tea "github.com/charmbracelet/bubbletea"
	"encoding/json"
	"github.com/atotto/clipboard"
	"time"
	"strings"
)

func Get() {
	mtrPath, pwdPath := util.VerifyInit()
	masterPassword := util.AuthMasterPassword(mtrPath)
	data, err := os.ReadFile(pwdPath)
    if err != nil {
        panic(err)
    }
	var jsonPwd map[string]string
	err = json.Unmarshal(data, &jsonPwd)
    if err != nil {
        panic(err)
    }
	var passwords []string
	for password := range jsonPwd {
		passwords = append(passwords, password)
	}
	p := tea.NewProgram(initialFuzzyModel(passwords))
    finalModel, err := p.Run()
	if err != nil {
        fmt.Println("Error running program:", err)
    }
	if m, ok := finalModel.(fuzzyModel); ok {
		pwd, err := cryptoutil.DecryptString(jsonPwd[m.selected], masterPassword)
		if err != nil {
			panic(err)
		}
		err = clipboard.WriteAll(pwd)
		if err != nil {
			fmt.Println("Failed to copy to clipboard password will display then clear")
			fmt.Print(pwd)
			time.Sleep(5 * time.Second)
			fmt.Print("\r", strings.Repeat(" ", len(pwd)+10), "\r") // overwrite with spaces

			return
		}

		fmt.Println("Password copied to clipboard!")
	}else {
        fmt.Println("No selection made.")
    }
}

type fuzzyModel struct {
    passwords []string
    query     string
    filtered  []string
	cursor int
	selected string
}

func initialFuzzyModel(pwds []string) fuzzyModel {
    return fuzzyModel{
        passwords: pwds,
        filtered:  pwds,
		cursor: 0,
		selected: "",
    }
}

func (m fuzzyModel) Init() tea.Cmd { return nil }

func (m fuzzyModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "ctrl+c", "q":
            return m, tea.Quit
        case "backspace":
            if len(m.query) > 0 {
                m.query = m.query[:len(m.query)-1]
            }
		case "enter":	
			if len(m.filtered) > 0 {
				m.selected = m.filtered[m.cursor]
				return m, tea.Quit
			}else {
            return m, nil
			}
        case "up":
            if m.cursor > 0 {
                m.cursor--
            }
        case "down":
            if m.cursor < len(m.filtered)-1 {
                m.cursor++
            }
        default:
            if len(msg.String()) == 1 { // letters/numbers
                m.query += msg.String()
            }
        }

        // Fuzzy search every time query changes
        if m.query == "" {
            m.filtered = m.passwords
        } else {
            matches := fuzzy.Find(m.query, m.passwords)
            m.filtered = nil
            for _, match := range matches {
                m.filtered = append(m.filtered, m.passwords[match.Index])
            }
        }

		if m.cursor >= len(m.filtered) {
			m.cursor = len(m.filtered) - 1
		}
		if m.cursor < 0 {
			m.cursor = 0
		}
    }

    return m, nil
}

func (m fuzzyModel) View() string {
    s := "Search: " + m.query + "\n"
    for i, pwd := range m.filtered {
		line := pwd
        if i == m.cursor {
            line = "   > " + pwd
        } else {
            line = "    " + pwd
        }
        s += line + "\n"
    }
    s += "\n(q to quit)"
    return s
}
