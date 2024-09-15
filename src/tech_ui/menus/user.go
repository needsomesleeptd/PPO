package menus

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/dixonwille/wmenu/v5"
)

func (m *Menu) RunUserMenu(client *http.Client) error {
	m.uMenu = wmenu.NewMenu("Enter your option:")
	m.AddOptionsUser(client)

	for {
		err := m.uMenu.Run()
		fmt.Println()
		if err != nil {
			if errors.Is(err, errExit) {
				break
			}

			fmt.Printf("ERROR: %v\n\n", err)
		}
	}

	fmt.Printf("Exited menu.\n")
	return nil
}
