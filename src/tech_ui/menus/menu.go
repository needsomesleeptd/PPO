package menus

import (
	"annotater/internal/models"
	"errors"
	"fmt"
	"net/http"

	"github.com/dixonwille/wmenu/v5"
)

var (
	errExit = errors.New("exiting")
)

type ClientEntity struct {
	Client *http.Client
}

type Menu struct {
	mainMenu *wmenu.Menu
	uMenu    *wmenu.Menu
	cMenu    *wmenu.Menu
	aMenu    *wmenu.Menu
	Role     models.Role
	login    string
	password string
	jwt      string
	role     models.Role
}

func NewMenu() *Menu {
	return &Menu{}
}

func (m *Menu) AddOptionsMain(client *http.Client) {
	m.mainMenu.Option("SignUp", ClientEntity{client}, false, m.SignUpMenu)
	m.mainMenu.Option("SignIn", ClientEntity{client}, false, m.SignInMenu)
	m.mainMenu.Option("Exit", ClientEntity{client}, false, func(_ wmenu.Opt) error {
		return errExit
	})
}

func (m *Menu) AddOptionsUser(client *http.Client) {
	m.uMenu.Option("Select file for check", ClientEntity{client}, false, m.checkingDocumentOpt)
	m.uMenu.Option("Load file", ClientEntity{client}, false, m.SignUpMenu)
	m.uMenu.Option("Exit", ClientEntity{client}, false, func(_ wmenu.Opt) error {
		return errExit
	})
}

func (m *Menu) AddOptionsController(client *http.Client) {
	m.uMenu.Option("Select file for checking", ClientEntity{client}, false, m.SignInMenu)
	m.uMenu.Option("Load file", ClientEntity{client}, false, m.SignUpMenu)
	m.uMenu.Option("Add markup", ClientEntity{client}, false, m.SignUpMenu)
	m.uMenu.Option("Getting all your markups", ClientEntity{client}, false, m.SignUpMenu)
	m.uMenu.Option("Deleting markup", ClientEntity{client}, false, m.SignUpMenu)
	m.uMenu.Option("Getting all markups", ClientEntity{client}, false, m.SignUpMenu)
	m.uMenu.Option("Exit", ClientEntity{client}, false, func(_ wmenu.Opt) error {
		return errExit
	})
}

func (m *Menu) RunMenu(client *http.Client) error {
	m.mainMenu = wmenu.NewMenu("Please SignIn or SignUp")
	m.AddOptionsMain(client)

	for {
		err := m.mainMenu.Run()
		fmt.Println()
		if err != nil {
			if errors.Is(err, errExit) {
				break
			}
		}
	}

	fmt.Printf("Exited menu.\n")

	return nil
}
