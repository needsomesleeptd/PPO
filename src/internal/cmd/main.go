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
	document_repo_adapter "annotater/internal/bl/documentService/documentRepo/documentRepoAdapter"
	service "annotater/internal/bl/userService"
	user_repo_adapter "annotater/internal/bl/userService/userRepo/userRepoAdapter"
	annot_handler "annotater/internal/http-server/handlers/annot"
	annot_type_handler "annotater/internal/http-server/handlers/annotType"
	auth_handler "annotater/internal/http-server/handlers/auth"
	document_handler "annotater/internal/http-server/handlers/document"
	"annotater/internal/middleware/access_middleware"
	"annotater/internal/middleware/auth_middleware"
	models_da "annotater/internal/models/modelsDA"
	auth_utils "annotater/internal/pkg/authUtils"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	CONN_POSTGRES_STR = "host=localhost user=andrew password=1 database=lab01db port=5432" //TODO:: export through parameters
	POSTGRES_CFG      = postgres.Config{DSN: CONN_POSTGRES_STR}
	MODEL_ROUTE       = "http://0.0.0.0:5000/pred"
)

func main() {
	db, err := gorm.Open(postgres.New(POSTGRES_CFG), &gorm.Config{})
	db.AutoMigrate(&models_da.Document{})
	db.AutoMigrate(&models_da.User{})
	db.AutoMigrate(&models_da.MarkupType{})
	db.AutoMigrate(&models_da.Markup{})

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

	documentRepo := document_repo_adapter.NewDocumentRepositoryAdapter(db)
	documentService := document_service.NewDocumentService(documentRepo, model)

	//userService 0_0
	userService := service.NewUserService(userRepo)

	//auth service
	router := chi.NewRouter()

	authMiddleware := (func(h http.Handler) http.Handler {
		return auth_middleware.JwtAuthMiddleware(h, auth_service.SECRET, tokenHandler)
	})

	accesMiddleware := access_middleware.NewAccessMiddleware(userService)
	router.Group(func(r chi.Router) { // group for which auth middleware is required
		r.Use(authMiddleware)

		// Document
		r.Route("/document", func(r chi.Router) {
			r.Use(accesMiddleware.ControllersAndHigherMiddleware) // apply the desired middleware here
			r.Post("/load", document_handler.LoadDocument(documentService))
			r.Get("/check", document_handler.CheckDocument(documentService))
		})

		// AnnotType
		r.Route("/annotType", func(r chi.Router) {
			r.Use(accesMiddleware.ControllersAndHigherMiddleware) // apply the desired middleware here
			r.Post("/add", annot_type_handler.AddAnnotType(annotTypeService))
			r.Get("/get", annot_type_handler.GetAnnotType(annotTypeService))
			r.Delete("/delete", annot_type_handler.DeleteAnnotType(annotTypeService))

		})
		//Annot
		r.Route("/annot", func(r chi.Router) {
			r.Use(accesMiddleware.ControllersAndHigherMiddleware) // apply the desired middleware here
			r.Post("/add", annot_handler.AddAnnot(annotService))
			r.Get("/get", annot_handler.GetAnnot(annotService))
			r.Delete("/delete", annot_handler.DeleteAnnot(annotService))
		})

	})

	//auth
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
