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
	m.uMenu.Option("Getting loaded files description", ClientEntity{client}, false, m.GettingReportsMetaData)
	m.uMenu.Option("Getting report by ID", ClientEntity{client}, false, m.GettingReportByID)
	m.uMenu.Option("Getting document by ID", ClientEntity{client}, false, m.GettingDocumentByID)
	m.uMenu.Option("Exit", ClientEntity{client}, false, func(_ wmenu.Opt) error {
		return errExit
	})
}

func (m *Menu) AddOptionsController(client *http.Client) {
	m.cMenu.Option("Select file for check", ClientEntity{client}, false, m.checkingDocumentOpt)
	m.cMenu.Option("Getting loaded files description", ClientEntity{client}, false, m.GettingReportsMetaData)
	m.cMenu.Option("Getting report by ID", ClientEntity{client}, false, m.GettingReportByID)
	m.cMenu.Option("Getting document by ID", ClientEntity{client}, false, m.GettingDocumentByID)

	m.cMenu.Option("Add markup", ClientEntity{client}, false, m.AddingAnotattion)
	m.cMenu.Option("Add markupType", ClientEntity{client}, false, m.AddingAnotattionType)
	m.cMenu.Option("Deleting markup", ClientEntity{client}, false, m.DeletingAnotattion)
	m.cMenu.Option("Getting all your markup types", ClientEntity{client}, false, m.GettingAnotattionType)
	m.cMenu.Option("Getting all your markups", ClientEntity{client}, false, m.GettingAnotattionsByUserID)

	m.cMenu.Option("Get all markupTypes", ClientEntity{client}, false, m.GettingAllAnottationTypes)
	m.cMenu.Option("Get all markups", ClientEntity{client}, false, m.GettingAllAnottations)

	m.cMenu.Option("Exit", ClientEntity{client}, false, func(_ wmenu.Opt) error {
		return errExit
	})
}

func (m *Menu) AddOptionsAdmin(client *http.Client) {
	m.aMenu.Option("Select file for check", ClientEntity{client}, false, m.checkingDocumentOpt)
	m.aMenu.Option("Getting loaded files description", ClientEntity{client}, false, m.GettingReportsMetaData)
	m.aMenu.Option("Getting report by ID", ClientEntity{client}, false, m.GettingReportByID)
	m.aMenu.Option("Getting document by ID", ClientEntity{client}, false, m.GettingDocumentByID)

	m.aMenu.Option("Add markup", ClientEntity{client}, false, m.AddingAnotattion)
	m.aMenu.Option("Add markupType", ClientEntity{client}, false, m.AddingAnotattionType)
	m.aMenu.Option("Deleting markup", ClientEntity{client}, false, m.DeletingAnotattion)
	m.aMenu.Option("Getting all your markup types", ClientEntity{client}, false, m.GettingAnotattionType)
	m.aMenu.Option("Getting all your markups", ClientEntity{client}, false, m.GettingAnotattionsByUserID)

	m.aMenu.Option("Change user role", ClientEntity{client}, false, m.ChangeUserRole)
	m.aMenu.Option("Getting all users Data", ClientEntity{client}, false, m.GettingAllUsers)
	m.aMenu.Option("Delete the whole markupType", ClientEntity{client}, false, m.DeletingAnotattionType)

	m.aMenu.Option("Get all markupTypes", ClientEntity{client}, false, m.GettingAllAnottationTypes)
	m.aMenu.Option("Get all markups", ClientEntity{client}, false, m.GettingAllAnottations)

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
