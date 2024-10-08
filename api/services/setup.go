package services

import (
	"context"

	"github.com/shellhub-io/shellhub/pkg/api/authorizer"
	"github.com/shellhub-io/shellhub/pkg/api/requests"
	"github.com/shellhub-io/shellhub/pkg/clock"
	"github.com/shellhub-io/shellhub/pkg/models"
)

type SetupService interface {
	Setup(ctx context.Context, req requests.Setup) error
}

func (s *service) Setup(ctx context.Context, req requests.Setup) error {
	data := models.UserData{
		Name:          req.Name,
		Email:         req.Email,
		Username:      req.Username,
		RecoveryEmail: "",
	}

	if ok, err := s.validator.Struct(data); !ok || err != nil {
		return NewErrUserInvalid(nil, err)
	}

	password, err := models.HashUserPassword(req.Password)
	if err != nil {
		return NewErrUserPasswordInvalid(err)
	}

	if ok, err := s.validator.Struct(password); !ok || err != nil {
		return NewErrUserPasswordInvalid(err)
	}

	user := &models.User{
		UserData: data,
		Password: password,
		// NOTE: user's created from the setup screen doesn't need to be confirmed.
		Status:    models.UserStatusConfirmed,
		CreatedAt: clock.Now(),
	}

	if err := s.store.UserCreate(ctx, user); err != nil {
		return NewErrUserDuplicated([]string{req.Username}, err)
	}

	namespace := &models.Namespace{
		Name:       req.Namespace,
		Owner:      user.ID,
		MaxDevices: 0,
		Members: []models.Member{
			{
				ID:   user.ID,
				Role: authorizer.RoleOwner,
			},
		},
		CreatedAt: clock.Now(),
		Settings: &models.NamespaceSettings{
			SessionRecord:          false,
			ConnectionAnnouncement: "",
		},
	}

	if _, err = s.store.NamespaceCreate(ctx, namespace); err != nil {
		return NewErrNamespaceDuplicated(err)
	}

	return nil
}
