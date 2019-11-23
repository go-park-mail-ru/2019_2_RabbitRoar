package session

import (
	"context"
	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/models"
	session "github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/session"
	_grpc "github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/session/delivery/grpc"
)

//go:generate protoc -I ../../pkg/session/delivery/grpc/ --go_out=plugins=grpc:../../pkg/session/delivery/grpc/ ../../pkg/session/delivery/grpc/session.proto

type manager struct {
	sessionUseCase session.UseCase
}

func NewManager(sessionUseCase session.UseCase) *manager {
	return &manager{
		sessionUseCase: sessionUseCase,
	}
}

func (m *manager) Create(ctx context.Context, in *_grpc.Session) (*_grpc.SessionID, error) {
	var user = models.User{
		ID:        int(in.User.ID),
		Username:  in.User.Username,
		Email:     in.User.Email,
		Rating:    int(in.User.Rating),
		AvatarUrl: in.User.Avatar,
	}
	sessionID, err := m.sessionUseCase.Create(user)
	if err != nil {
		return nil, err
	}
	return &_grpc.SessionID{ID: *sessionID}, nil
}

func (m *manager) GetByID(ctx context.Context, in *_grpc.SessionID) (*_grpc.Session, error) {
	sess, err := m.sessionUseCase.GetByID(in.ID)
	if err != nil {
		return nil, err
	}
	var u = &_grpc.User{
		ID:                   int32(sess.User.ID),
		Username:             sess.User.Username,
		Email:                sess.User.Email,
		Rating:               int32(sess.User.Rating),
		Avatar:               sess.User.AvatarUrl,
	}
	return &_grpc.Session{
		SessionID:            sess.ID,
		User:                 u,
	}, nil
}

func (m *manager) Delete(ctx context.Context, in *_grpc.SessionID) (*_grpc.Nothing, error) {
	return &_grpc.Nothing{}, m.sessionUseCase.Destroy(in.ID)
}
