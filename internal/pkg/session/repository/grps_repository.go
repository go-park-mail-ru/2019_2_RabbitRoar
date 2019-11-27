package repository

import (
	"context"
	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/models"
	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/session"
	_grpc "github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/session/delivery/grpc"
	"google.golang.org/grpc"
)

//go:generate protoc -I ../../session/delivery/grpc/ --go_out=plugins=grpc:../../session/delivery/grpc/ ../../session/delivery/grpc/session.proto

type grpcSessionRepository struct {
	client _grpc.SessionServiceClient
}

func NewGrpcSessionRepository(conn *grpc.ClientConn) session.Repository {
	return &grpcSessionRepository{
		client: _grpc.NewSessionServiceClient(conn),
	}
}

func (g grpcSessionRepository) Create(u models.User) (*string, error) {
	sessionID, err := g.client.Create(
		context.Background(),
		&_grpc.User{
			ID:       int32(u.ID),
			Username: u.Username,
			Email:    u.Email,
			Rating:   int32(u.Rating),
			Avatar:   u.AvatarUrl,
		},
	)

	if err != nil {
		return nil, err
	}
	return &sessionID.ID, nil
}

func (g grpcSessionRepository) Destroy(sessionID string) error {
	_, err := g.client.Delete(
		context.Background(),
		&_grpc.SessionID{
			ID: sessionID,
		},
	)
	return err
}

func (g grpcSessionRepository) GetByID(sessionID string) (*models.Session, error) {
	sess, err := g.client.GetByID(
		context.Background(),
		&_grpc.SessionID{
			ID: sessionID,
		},
	)

	if err != nil {
		return nil, err
	}

	return &models.Session{
		ID:   sess.SessionID,
		User: models.User{
			ID:        int(sess.User.ID),
			Username:  sess.User.Username,
			Email:     sess.User.Email,
			Rating:    int(sess.User.Rating),
			AvatarUrl: sess.User.Avatar,
		},
	}, nil
}
