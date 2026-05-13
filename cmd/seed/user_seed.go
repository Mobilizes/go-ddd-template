package seed

import (
	"fmt"
	"mob/ddd-template/internal/domain/entity"
	"mob/ddd-template/internal/infra/security"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type userSeed struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	PasswordHash string `json:"password_hash"`
}

func SeedUsers(db *gorm.DB, path string) error {
	var rows []userSeed
	if err := readSeedJSON(path, &rows); err != nil {
		return err
	}

	hasher := security.NewHasher()
	users := make([]*entity.User, 0, len(rows))
	for _, row := range rows {
		if row.ID == "" {
			row.ID = uuid.NewString()
		}

		password := row.PasswordHash
		if password == "" {
			if row.Password == "" {
				return fmt.Errorf("user %q must have password or password_hash", row.Email)
			}

			hashedPassword, err := hasher.RandomHash(row.Password)
			if err != nil {
				return err
			}
			password = hashedPassword
		}

		users = append(users, entity.NewUser(row.ID, row.Name, row.Email, password))
	}

	if len(users) == 0 {
		return nil
	}

	return db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "email"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"name",
			"password",
			"updated_at",
			"deleted_at",
		}),
	}).Create(&users).Error
}


