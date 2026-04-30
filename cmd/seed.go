package cmd

import (
	"fmt"
	"mob/ddd-template/internal/domain/entity"
	"mob/ddd-template/internal/infra/persistence"
	"mob/ddd-template/internal/infra/security"

	"github.com/google/uuid"
)

func Seed() {
	db := SetUpDatabaseConnectionOrFail()
	fmt.Println("Database connection established for seeding.")

	hasher := security.NewBcryptHasher()
	hashedPassword, err := hasher.Hash("password123")
	if err != nil {
		fmt.Printf("Error hashing password: %v\n", err)
		return
	}

	userRepo := persistence.NewUserPersistence(db)

	adminUser := entity.NewUser(uuid.NewString(), "Admin", "admin@example.com", hashedPassword)
	err = userRepo.Create(adminUser)
	if err != nil {
		fmt.Printf("Error seeding admin user: %v\n", err)
		return
	}

	testUser := entity.NewUser(uuid.NewString(), "Test User", "test@example.com", hashedPassword)
	err = userRepo.Create(testUser)
	if err != nil {
		fmt.Printf("Error seeding test user: %v\n", err)
		return
	}

	fmt.Println("Seeding completed successfully.")
}
