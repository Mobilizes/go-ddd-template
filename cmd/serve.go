package cmd

import (
	"fmt"
	"mob/ddd-template/internal/app"
	"mob/ddd-template/internal/domain/repository"
	"mob/ddd-template/internal/domain/service"
	"mob/ddd-template/internal/infra/persistence"
	"mob/ddd-template/internal/infra/security"
	"mob/ddd-template/internal/presentation"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
	"gorm.io/gorm"
)

func injectInfra(injector do.Injector, db *gorm.DB) {
	do.Provide(injector, func(i do.Injector) (repository.UserRepository, error) {
		return persistence.NewUserPersistence(db), nil
	})

	do.Provide(injector, func(i do.Injector) (service.PasswordHasher, error) {
		return security.NewBcryptPasswordHasher(), nil
	})

	do.Provide(injector, func(i do.Injector) (service.TokenGenerator, error) {
		return security.NewJWTTokenGenerator(os.Getenv("JWT_SECRET"), 15*time.Minute), nil
	})
}

func injectApp(injector do.Injector) {
	do.Provide(injector, func(i do.Injector) (app.UserUseCase, error) {
		return app.NewUserUseCase(i), nil
	})

	do.Provide(injector, func(i do.Injector) (app.AuthUseCase, error) {
		return app.NewAuthUseCase(i), nil
	})
}

func injectPresentation(injector do.Injector, server *gin.Engine) {
	presentation.RegisterUserRoutes(server, presentation.NewUserPresentation(injector))
}

func Serve() {
	injector := do.New()

	db := SetUpDatabaseConnectionOrFail()

	do.Provide(injector, func(i do.Injector) (*gorm.DB, error) {
		return db, nil
	})

	server := gin.Default()

	port := os.Getenv("GOLANG_PORT")
	if port == "" {
		port = "8888"
	}

	var serve string
	if os.Getenv("APP_ENV") == "localhost" {
		serve = "0.0.0.0:" + port
	} else {
		serve = ":" + port
	}

	injectInfra(injector, db)
	injectApp(injector)
	injectPresentation(injector, server)

	if err := server.Run(serve); err != nil {
		fmt.Errorf("error running server: %v", err)
	}
}
