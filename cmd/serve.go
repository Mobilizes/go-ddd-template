package cmd

import (
	"fmt"
	"mob/ddd-template/internal/app/port"
	"mob/ddd-template/internal/app/usecase"
	"mob/ddd-template/internal/domain/repository"
	"mob/ddd-template/internal/infra/persistence"
	"mob/ddd-template/internal/infra/security"
	"mob/ddd-template/internal/infra/storage"
	"mob/ddd-template/internal/presentation/handler"
	"mob/ddd-template/internal/presentation/route"
	"os"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/samber/do/v2"
	"gorm.io/gorm"
)

func injectRepository(injector do.Injector, db *gorm.DB) {
	do.Provide(injector, func(i do.Injector) (repository.UserRepository, error) {
		return persistence.NewUserPersistence(db), nil
	})

	do.Provide(injector, func(i do.Injector) (repository.RefreshTokenRepository, error) {
		return persistence.NewRefreshTokenPersistence(db), nil
	})

	do.Provide(injector, func(i do.Injector) (repository.FileRepository, error) {
		return persistence.NewFilePersistence(db), nil
	})
}

func injectPort(injector do.Injector, db *gorm.DB) {
	do.Provide(injector, func(i do.Injector) (port.UnitOfWork, error) {
		return persistence.NewGormUnitOfWork(db), nil
	})

	do.Provide(injector, func(i do.Injector) (port.Hasher, error) {
		return security.NewHasher(), nil
	})

	do.Provide(injector, func(i do.Injector) (port.TokenGenerator, error) {
		return security.NewJWTTokenGenerator(os.Getenv("JWT_SECRET"), 15*time.Minute), nil
	})

	do.Provide(injector, func(i do.Injector) (port.FileStorage, error) {
		return storage.NewLocalStorage(os.Getenv("STORAGE_PATH")), nil
	})
}

func injectInfra(injector do.Injector, db *gorm.DB) {
	injectRepository(injector, db)
	injectPort(injector, db)
}

func injectApp(injector do.Injector) {
	do.Provide(injector, func(i do.Injector) (usecase.UserUseCase, error) {
		return usecase.NewUserUseCase(i), nil
	})

	do.Provide(injector, func(i do.Injector) (usecase.AuthUseCase, error) {
		return usecase.NewAuthUseCase(i), nil
	})
}

func injectPresentation(injector do.Injector, app *fiber.App) {
	route.RegisterUser(handler.NewUserHandler(injector), app)
	route.RegisterAuth(handler.NewAuthHandler(injector), app)
	route.RegisterHealthCheck(app)
}

func Serve() {
	injector := do.New()

	db := SetUpDatabaseConnectionOrFail()

	do.Provide(injector, func(i do.Injector) (*gorm.DB, error) {
		return db, nil
	})

	app := fiber.New()

	app.Use(cors.New())

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
	injectPresentation(injector, app)

	if err := app.Listen(serve, fiber.ListenConfig{EnablePrintRoutes: os.Getenv("PRINT_ROUTES") == "true"}); err != nil {
		fmt.Printf("error running server: %v", err)
	}
}
