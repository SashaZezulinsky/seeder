package usecase

import (
	"context"
	"seeder/pkg/errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"seeder/internal/domain"
	"seeder/internal/node/mock"
)

var node = &domain.Node{
	IP:      "127.0.0.1",
	Name:    "name",
	Version: "version",
	Client:  "client",
}

func TestNodeUsecase_AddNode(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockNodeRepo := mock.NewMockNodeRepository(ctrl)
	authUC := NewNodeUseCase(mockNodeRepo)

	mockNodeRepo.EXPECT().AddNode(gomock.Any(), node).Return(nil)
	mockNodeRepo.EXPECT().FindNode(gomock.Any(), node).Return(errors.ErrNotFound)

	err := authUC.AddNode(context.Background(), node)

	require.NoError(t, err)
}

func TestNodeUsecase_AddExistingNode(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockNodeRepo := mock.NewMockNodeRepository(ctrl)
	authUC := NewNodeUseCase(mockNodeRepo)

	mockNodeRepo.EXPECT().FindNode(gomock.Any(), node).Return(nil)

	err := authUC.AddNode(context.Background(), node)

	require.NoError(t, err)
}

func TestNodeUsecase_GetNodesList(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockNodeRepo := mock.NewMockNodeRepository(ctrl)
	authUC := NewNodeUseCase(mockNodeRepo)

	nodes := []*domain.Node{node}
	mockNodeRepo.EXPECT().GetNodesList(gomock.Any()).Return(nodes, nil)

	list, err := authUC.GetNodesList(context.Background())

	require.NoError(t, err)
	require.Equal(t, nodes, list)
}
