package main

import (
	nn_adapter "annotater/internal/bl/NN/NNAdapter"
	nn_model_handler "annotater/internal/bl/NN/NNAdapter/NNmodelhandler"
	annot_service "annotater/internal/bl/annotationService"
	annot_repo_adapter "annotater/internal/bl/annotationService/annotattionRepo/anotattionRepoAdapter"
	annot_type_service "annotater/internal/bl/anotattionTypeService"
	annot_type_repo_adapter "annotater/internal/bl/anotattionTypeService/anottationTypeRepo/anotattionTypeRepoAdapter"
	auth_service "annotater/internal/bl/auth"
	document_service "annotater/internal/bl/documentService"
	doc_data_repo_adapter "annotater/internal/bl/documentService/documentDataRepo/documentDataRepo"
	document_repo_adapter "annotater/internal/bl/documentService/documentMetaDataRepo/documentMetaDataRepoAdapter"
	rep_data_repo_adapter "annotater/internal/bl/documentService/reportDataRepo/reportDataRepoAdapter"
	rep_creator_service "annotater/internal/bl/reportCreatorService"
	report_creator "annotater/internal/bl/reportCreatorService/reportCreator"
	service "annotater/internal/bl/userService"
	user_repo_adapter "annotater/internal/bl/userService/userRepo/userRepoAdapter"
	annot_handler "annotater/internal/http-server/handlers/annot"
	annot_type_handler "annotater/internal/http-server/handlers/annotType"
	auth_handler "annotater/internal/http-server/handlers/auth"
	document_handler "annotater/internal/http-server/handlers/document"
	user_handler "annotater/internal/http-server/handlers/user"
	"annotater/internal/middleware/access_middleware"
	"annotater/internal/middleware/auth_middleware"
	models_da "annotater/internal/models/modelsDA"
	auth_utils "annotater/internal/pkg/authUtils"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	CONN_POSTGRES_STR    = "host=localhost user=andrew password=1 database=lab01db port=5432" //TODO:: export through parameters
	POSTGRES_CFG         = postgres.Config{DSN: CONN_POSTGRES_STR}
	MODEL_ROUTE          = "http://0.0.0.0:5000/pred"
	REPORTS_CREATOR_PATH = "../../../reportsCreator"
	REPORTS_PATH         = "../../../reports"
	DOCUMENTS_PATH       = "../../../documents"
	DOCUMENTS_EXT        = ".pdf"
	REPORTS_EXT          = ".pdf"
)

// andrew1 2
// admin admin
// control control

func migrate(db *gorm.DB) error {
	err := db.AutoMigrate(&models_da.Document{})
	if err != nil {
		return err
	}
	err = db.AutoMigrate(&models_da.User{})
	if err != nil {
		return err
	}
	err = db.AutoMigrate(&models_da.MarkupType{})
	if err != nil {
		return err
	}
	err = db.AutoMigrate(&models_da.Markup{})
	if err != nil {
		return err
	}
	return nil
}

func setuplog() *slog.Logger {
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	return log
}

func main() {
	db, err := gorm.Open(postgres.New(POSTGRES_CFG), &gorm.Config{TranslateError: true})
	log := setuplog()

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	err = migrate(db)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	//auth service
	userRepo := user_repo_adapter.NewUserRepositoryAdapter(db)
	hasher := auth_utils.NewPasswordHashCrypto()
	tokenHandler := auth_utils.NewJWTTokenHandler()
	authService := auth_service.NewAuthService(userRepo, hasher, tokenHandler, auth_service.SECRET)

	//annot service
	annotRepo := annot_repo_adapter.NewAnotattionRepositoryAdapter(db)
	annotService := annot_service.NewAnnotattionService(annotRepo)

	//annotType service
	annotTypeRepo := annot_type_repo_adapter.NewAnotattionTypeRepositoryAdapter(db)
	annotTypeService := annot_type_service.NewAnotattionTypeService(annotTypeRepo)

	//document service
	//setting up NN
	modelhandler := nn_model_handler.NewHttpModelHandler(MODEL_ROUTE)
	model := nn_adapter.NewDetectionModel(modelhandler)

	reportCreator := report_creator.NewPDFReportCreator(REPORTS_CREATOR_PATH)
	reportCreatorService := rep_creator_service.NewDocumentService(model, annotTypeRepo, reportCreator)

	documentStorage := doc_data_repo_adapter.NewDocumentRepositoryAdapter(DOCUMENTS_PATH, DOCUMENTS_EXT)

	reportStorage := rep_data_repo_adapter.NewDocumentRepositoryAdapter(REPORTS_PATH, REPORTS_EXT)

	documentRepo := document_repo_adapter.NewDocumentRepositoryAdapter(db)
	documentService := document_service.NewDocumentService(documentRepo, documentStorage, reportStorage, reportCreatorService)

	//userService 0_0
	userService := service.NewUserService(userRepo)

	userHandler := user_handler.NewDocumentHandler(log, userService)
	//auth service
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	authMiddleware := (func(h http.Handler) http.Handler {
		return auth_middleware.JwtAuthMiddleware(h, auth_service.SECRET, tokenHandler)
	})

	accesMiddleware := access_middleware.NewAccessMiddleware(userService)

	documentHandler := document_handler.NewDocumentHandler(log, documentService)

	router.Group(func(r chi.Router) { // group for which auth middleware is required
		r.Use(authMiddleware)

		// Document
		r.Route("/document", func(r chi.Router) {
			r.Post("/report", documentHandler.CreateReport())
			r.Get("/getDocument", documentHandler.GetDocumentByID())
			r.Get("/getReport", documentHandler.GetReportByID())
			r.Get("/getDocumentsMeta", documentHandler.GetDocumentsMetaData())
		})

		// AnnotType
		r.Route("/annotType", func(r chi.Router) {
			r.Use(accesMiddleware.ControllersAndHigherMiddleware) // apply the desired middleware here

			adminOnlyAnnotTypes := r.Group(nil)
			adminOnlyAnnotTypes.Use(accesMiddleware.AdminOnlyMiddleware)

			r.Post("/add", annot_type_handler.AddAnnotType(annotTypeService))
			r.Get("/get", annot_type_handler.GetAnnotType(annotTypeService))

			r.Get("/creatorID", annot_type_handler.GetAnnotTypesByCreatorID(annotTypeService))

			r.Get("/gets", annot_type_handler.GetAnnotTypesByIDs(annotTypeService))

			adminOnlyAnnotTypes.Delete("/delete", annot_type_handler.DeleteAnnotType(annotTypeService))
			r.Get("/getsAll", annot_type_handler.GetAllAnnotTypes(annotTypeService))

		})
		//Annot
		r.Route("/annot", func(r chi.Router) {
			r.Use(accesMiddleware.ControllersAndHigherMiddleware)
			//adminOnlyAnnots := r.Group(nil)
			//adminOnlyAnnots.Use(accesMiddleware.AdminOnlyMiddleware)

			r.Post("/add", annot_handler.AddAnnot(annotService))
			r.Get("/get", annot_handler.GetAnnot(annotService))
			r.Get("/creatorID", annot_handler.GetAnnotsByUserID(annotService))

			r.Delete("/delete", annot_handler.DeleteAnnot(annotService))
			r.Get("/getsAll", annot_handler.GetAllAnnots(annotService))
		})
		//user
		r.Route("/user", func(r chi.Router) {
			r.Use(accesMiddleware.AdminOnlyMiddleware)
			r.Post("/role", userHandler.ChangeUserPerms())
			r.Get("/getUsers", userHandler.GetAllUsers())
		})

	})

	//auth, no middleware is required
	router.Post("/user/SignUp", auth_handler.SignUp(authService))
	router.Post("/user/SignIn", auth_handler.SignIn(authService))

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	srv := &http.Server{
		Addr:         "0.0.0.0:8080",
		Handler:      router,
		ReadTimeout:  40 * time.Second,
		WriteTimeout: 40 * time.Second,
		IdleTimeout:  40 * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			fmt.Println("error with server")
		}
	}()

	<-done
}
