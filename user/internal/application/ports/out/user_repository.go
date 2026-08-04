package out

import (
	"database/sql"
	"github.com/aakashloyar/elevate/user/internal/domain"
)

type UserRepository interface {
	Save(user domain.User) error
	FindByID(userID string) (domain.User, error)
	Delete(userID string) (sql.Result, error)
	ExistsByUsername(username string) (bool, error)
	ExistsByEmail(email string) (bool, error)
}
