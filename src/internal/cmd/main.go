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
	"annotater/internal/config"
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
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	easy "github.com/t-tomalak/logrus-easy-formatter"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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

func setuplog(conf *config.Config) *logrus.Logger {

	log := logrus.New()
	useFile := conf.Logger.UseFile
	if useFile {
		log.Printf("using file %s\n", conf.OutputFilePath)
		f, err := os.OpenFile(conf.OutputFilePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			log.Printf("Failed to create logfile %s:%s, defaulting to stderr", conf.OutputFilePath, err.Error())
			useFile = false
		}
		defer f.Close()
		log.SetOutput(f)
	} else {
		log.SetOutput(os.Stderr)
	}

	easyFormatter := &easy.Formatter{
		TimestampFormat: conf.TimestampFormat,
		LogFormat:       conf.LogFormat,
	}

	log.SetFormatter(easyFormatter)
	if conf.OutputFormat == "text" {
		log.SetFormatter(&logrus.TextFormatter{})
	}
	if conf.OutputFormat == "json" {
		log.SetFormatter(&logrus.JSONFormatter{})
	}
	log.SetReportCaller(true)
	return log
}

func main() {

	config := config.MustLoad()
	postgresConStr := config.Database.GetGormConnectStr()
	db, err := gorm.Open(postgres.New(postgres.Config{DSN: postgresConStr}), &gorm.Config{TranslateError: true})

	log := setuplog(config)

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
	annotTypeService := annot_type_service.NewAnotattionTypeService(log, annotTypeRepo)

	//document service
	//setting up NN
	modelhandler := nn_model_handler.NewHttpModelHandler(log, config.Model.Route)
	model := nn_adapter.NewDetectionModel(modelhandler)

	reportCreator := report_creator.NewPDFReportCreator(config.ReportCreatorPath)
	reportCreatorService := rep_creator_service.NewDocumentService(model, annotTypeRepo, reportCreator)

	documentStorage := doc_data_repo_adapter.NewDocumentRepositoryAdapter(config.DocumentPath, config.DocumentExt)

	reportStorage := rep_data_repo_adapter.NewDocumentRepositoryAdapter(config.ReportPath, config.ReportExt)

	documentRepo := document_repo_adapter.NewDocumentRepositoryAdapter(db)
	documentService := document_service.NewDocumentService(log, documentRepo, documentStorage, reportStorage, reportCreatorService)

	//userService 0_0
	userService := service.NewUserService(userRepo)

	//handlers
	userHandler := user_handler.NewDocumentHandler(log, userService)
	documentHandler := document_handler.NewDocumentHandler(log, documentService)
	annotHandler := annot_handler.NewAnnotHandler(log, annotService)
	annotTypeHandler := annot_type_handler.NewAnnotTypehandler(log, annotTypeService)

	authHandler := auth_handler.NewAuthHandler(log, authService)
	//auth service
	router := chi.NewRouter()
	//router.Use(middleware.Logger)

	authMiddleware := (func(h http.Handler) http.Handler {
		return auth_middleware.JwtAuthMiddleware(h, auth_service.SECRET, tokenHandler)
	})

	accesMiddleware := access_middleware.NewAccessMiddleware(userService)

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

			r.Post("/add", annotTypeHandler.AddAnnotType())
			r.Get("/get", annotTypeHandler.GetAnnotType())

			r.Get("/creatorID", annotTypeHandler.GetAnnotTypesByCreatorID())

			r.Get("/gets", annotTypeHandler.GetAnnotTypesByIDs())

			adminOnlyAnnotTypes.Delete("/delete", annotTypeHandler.DeleteAnnotType())
			r.Get("/getsAll", annotTypeHandler.GetAllAnnotTypes())

		})
		//Annot
		r.Route("/annot", func(r chi.Router) {
			r.Use(accesMiddleware.ControllersAndHigherMiddleware)
			//adminOnlyAnnots := r.Group(nil)
			//adminOnlyAnnots.Use(accesMiddleware.AdminOnlyMiddleware)

			r.Post("/add", annotHandler.AddAnnot())
			r.Get("/get", annotHandler.GetAnnot())
			r.Get("/creatorID", annotHandler.GetAnnotsByUserID())

			r.Delete("/delete", annotHandler.DeleteAnnot())
			r.Get("/getsAll", annotHandler.GetAllAnnots())
		})
		//user
		r.Route("/user", func(r chi.Router) {
			r.Use(accesMiddleware.AdminOnlyMiddleware)
			r.Post("/role", userHandler.ChangeUserPerms())
			r.Get("/getUsers", userHandler.GetAllUsers())
		})

	})

	//auth, no middleware is required
	router.Post("/user/SignUp", authHandler.SignUp())
	router.Post("/user/SignIn", authHandler.SignIn())

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	srv := &http.Server{
		Addr:         config.Addr,
		Handler:      router,
		ReadTimeout:  config.ReadTimeout,
		WriteTimeout: config.WriteTimeout,
		IdleTimeout:  config.IdleTimeout,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			fmt.Println("error with server")
		}
	}()

	<-done
}
