package drivers

import (
	"context"

	"github.com/JetBrainer/sso/internal/domain/models"
)

type RolesRepository interface {
	Create(ctx context.Context, role *models.Role) error
	RoleByName(ctx context.Context, name string) (*models.Role, error)
	Roles(ctx context.Context) ([]models.Role, error)
	Update(ctx context.Context, role *models.Role) error
	DeleteByName(ctx context.Context, name string) error
}

