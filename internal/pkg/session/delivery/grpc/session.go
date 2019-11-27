package session

import (
	"context"
	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/models"
	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/session"
)

//go:generate protoc -I . --go_out=plugins=grpc:. ./session.proto

type manager struct {
	sessionRepo session.Repository
}

func NewManager(sessionRepo session.Repository) SessionServiceServer {
	return &manager{
		sessionRepo: sessionRepo,
	}
}

func (m *manager) Create(ctx context.Context, in *User) (*SessionID, error) {
	var user = models.User{
		ID:        int(in.ID),
		Username:  in.Username,
		Email:     in.Email,
		Rating:    int(in.Rating),
		AvatarUrl: in.Avatar,
	}
	sessionID, err := m.sessionRepo.Create(user)
	if err != nil {
		return nil, err
	}
	return &SessionID{ID: *sessionID}, nil
}

func (m *manager) GetByID(ctx context.Context, in *SessionID) (*Session, error) {
	sess, err := m.sessionRepo.GetByID(in.ID)
	if err != nil {
		return nil, err
	}
	var u = &User{
		ID:                   int32(sess.User.ID),
		Username:             sess.User.Username,
		Email:                sess.User.Email,
		Rating:               int32(sess.User.Rating),
		Avatar:               sess.User.AvatarUrl,
	}
	return &Session{
		SessionID:            sess.ID,
		User:                 u,
	}, nil
}

func (m *manager) Delete(ctx context.Context, in *SessionID) (*Nothing, error) {
	return &Nothing{}, m.sessionRepo.Destroy(in.ID)
}

