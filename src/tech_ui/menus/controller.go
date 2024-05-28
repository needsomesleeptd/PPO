package menus

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/dixonwille/wmenu/v5"
)

func (m *Menu) RunControllerMenu(client *http.Client) error {
	m.cMenu = wmenu.NewMenu("Enter your option:")
	m.AddOptionsController(client)

	for {
		err := m.cMenu.Run()
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
