package players

import (
	"context"
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/hoop33/roster/models"
	"github.com/stretchr/testify/assert"
)

type mockNextService struct {
	called bool
}

func (m *mockNextService) ListPlayers(_ context.Context, _ string) ([]models.Player, error) {
	m.called = true
	return nil, nil
}

func (m *mockNextService) GetPlayer(_ context.Context, _ int) (*models.Player, error) {
	m.called = true
	return nil, nil
}

func (m *mockNextService) SavePlayer(_ context.Context, _ *models.Player) (*models.Player, bool, error) {
	m.called = true
	return nil, false, nil
}

func (m *mockNextService) DeletePlayer(_ context.Context, _ int) error {
	m.called = true
	return nil
}

func TestListPlayersShouldCallNext(t *testing.T) {
	m := &mockNextService{}
	s := NewLoggingService(log.NewNopLogger(), m)
	assert.False(t, m.called)
	_, err := s.ListPlayers(context.Background(), "")
	assert.Nil(t, err)
	assert.True(t, m.called)
}

func TestGetPlayerShouldCallNext(t *testing.T) {
	m := &mockNextService{}
	s := NewLoggingService(log.NewNopLogger(), m)
	assert.False(t, m.called)
	_, err := s.GetPlayer(context.Background(), 0)
	assert.Nil(t, err)
	assert.True(t, m.called)
}

func TestSavePlayerShouldCallNext(t *testing.T) {
	m := &mockNextService{}
	s := NewLoggingService(log.NewNopLogger(), m)
	assert.False(t, m.called)
	_, _, err := s.SavePlayer(context.Background(), &models.Player{})
	assert.Nil(t, err)
	assert.True(t, m.called)
}

func TestDeletePlayerShouldCallNext(t *testing.T) {
	m := &mockNextService{}
	s := NewLoggingService(log.NewNopLogger(), m)
	assert.False(t, m.called)
	err := s.DeletePlayer(context.Background(), 0)
	assert.Nil(t, err)
	assert.True(t, m.called)
}
