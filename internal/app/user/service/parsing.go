package service

import (
	"encoding/json"
	"sakucita/internal/domain"
	"sakucita/internal/infra/postgres/repository"

	"github.com/gofiber/fiber/v3"
)

func userRepoToUserDomain(u repository.User) domain.User {
	return domain.User{
		ID:            u.ID,
		Email:         u.Email,
		EmailVerified: u.EmailVerified,
		Phone:         &u.Phone.String,
		Name:          u.Name,
		Nickname:      u.Nickname,
		ImageUrl:      &u.ImageUrl.String,
		SingleSession: u.SingleSession,
		Meta:          u.Meta,
		CreatedAt:     u.CreatedAt.Time,
		UpdatedAt:     u.UpdatedAt.Time,
		DeletedAt:     u.DeletedAt.Time,
	}
}

func userRepoWithRolesToUserDomain(u repository.GetUserByIDWithRolesRow) (domain.UserWithRoles, error) {
	var roles []domain.Role
	if len(u.Roles) > 0 {
		if err := json.Unmarshal(u.Roles, &roles); err != nil {
			return domain.UserWithRoles{}, domain.NewAppError(fiber.StatusInternalServerError, domain.ErrMsgInternalServerError, domain.ErrInternalServerError)
		}
	}
	return domain.UserWithRoles{
		User: domain.User{
			ID:            u.ID,
			Email:         u.Email,
			EmailVerified: u.EmailVerified,
			Phone:         &u.Phone.String,
			Name:          u.Name,
			Nickname:      u.Nickname,
			ImageUrl:      &u.ImageUrl.String,
			SingleSession: u.SingleSession,
			Meta:          u.Meta,
			CreatedAt:     u.CreatedAt.Time,
			UpdatedAt:     u.UpdatedAt.Time,
			DeletedAt:     u.DeletedAt.Time,
		},
		Roles: roles,
	}, nil
}
