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
	login    string
	password string
	jwt      string
	role     models.Role
	ID       uint64
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
	m.uMenu.Option("Load file", ClientEntity{client}, false, m.LoadingDocument)
	m.uMenu.Option("Exit", ClientEntity{client}, false, func(_ wmenu.Opt) error {
		return errExit
	})
}

func (m *Menu) AddOptionsController(client *http.Client) {
	m.cMenu.Option("Select file for checking", ClientEntity{client}, false, m.checkingDocumentOpt)
	m.cMenu.Option("Load file", ClientEntity{client}, false, m.LoadingDocument)
	m.cMenu.Option("Add markup", ClientEntity{client}, false, m.AddingAnotattion)
	m.cMenu.Option("Add markupType", ClientEntity{client}, false, m.AddingAnotattionType)
	m.cMenu.Option("Deleting markup", ClientEntity{client}, false, m.SignUpMenu)
	m.cMenu.Option("Getting all your markup types", ClientEntity{client}, false, m.GettingAnotattionType)
	m.cMenu.Option("Exit", ClientEntity{client}, false, func(_ wmenu.Opt) error {
		return errExit
	})
}

func (m *Menu) AddOptionsAdmin(client *http.Client) {
	m.aMenu.Option("Select file for checking", ClientEntity{client}, false, m.checkingDocumentOpt)
	m.aMenu.Option("Load file", ClientEntity{client}, false, m.LoadingDocument)
	m.aMenu.Option("Add markup", ClientEntity{client}, false, m.AddingAnotattion)
	m.aMenu.Option("Add markupType", ClientEntity{client}, false, m.AddingAnotattionType)
	m.aMenu.Option("Deleting markup", ClientEntity{client}, false, m.DeletingAnotattion)
	m.aMenu.Option("Getting all your markup types", ClientEntity{client}, false, m.GettingAnotattionType)
	m.aMenu.Option("Change user role", ClientEntity{client}, false, m.ChangeUserRole)
	m.aMenu.Option("Delete the whole anotattion type", ClientEntity{client}, false, nil)
	m.aMenu.Option("Exit", ClientEntity{client}, false, func(_ wmenu.Opt) error {
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
			fmt.Printf("ERROR: %v\n\n", err)
		}

	}

	fmt.Printf("Exited menu.\n")

	return nil
}