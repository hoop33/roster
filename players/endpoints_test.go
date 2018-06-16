package players

import (
	"context"
	"errors"
	"testing"

	"github.com/hoop33/roster/models"
	"github.com/stretchr/testify/assert"
)

var jr = models.Player{Name: "Jalen Ramsey"}

type mockSuccessService struct{}

func (m *mockSuccessService) ListPlayers(context.Context, string) ([]models.Player, error) {
	return []models.Player{jr}, nil
}

func (m *mockSuccessService) GetPlayer(context.Context, int) (*models.Player, error) {
	return &jr, nil
}

func (m *mockSuccessService) SavePlayer(context.Context, *models.Player) (*models.Player, bool, error) {
	return &jr, false, nil
}

func (m *mockSuccessService) DeletePlayer(context.Context, int) error {
	return nil
}

var successSvc = &mockSuccessService{}

type mockFailService struct{}

func (m *mockFailService) ListPlayers(context.Context, string) ([]models.Player, error) {
	return nil, errors.New("fail")
}

func (m *mockFailService) GetPlayer(context.Context, int) (*models.Player, error) {
	return nil, errors.New("fail")
}

func (m *mockFailService) SavePlayer(context.Context, *models.Player) (*models.Player, bool, error) {
	return nil, false, errors.New("fail")
}

func (m *mockFailService) DeletePlayer(context.Context, int) error {
	return errors.New("fail")
}

var failSvc = &mockFailService{}

func TestMakeListPlayersEndpointShouldReturnFuncThatReturnsListPlayersResponse(t *testing.T) {
	ep := NewEndpoints(successSvc)
	resp, err := ep.listPlayersEndpoint(context.Background(), listPlayersRequest{})
	assert.Nil(t, err)
	lpr, ok := resp.(listPlayersResponse)
	assert.True(t, ok)
	assert.Equal(t, 1, len(lpr.Players))
	assert.Equal(t, "Jalen Ramsey", lpr.Players[0].Name)
}

func TestMakeGetPlayerEndpointShouldReturnFuncThatReturnsGetPlayerResponse(t *testing.T) {
	ep := NewEndpoints(successSvc)
	resp, err := ep.getPlayerEndpoint(context.Background(), getPlayerRequest{})
	assert.Nil(t, err)
	gpr, ok := resp.(getPlayerResponse)
	assert.True(t, ok)
	assert.Equal(t, "Jalen Ramsey", gpr.Player.Name)
}

func TestMakeSavePlayerEndpointShouldReturnFuncThatReturnsSavePlayerResponse(t *testing.T) {
	ep := NewEndpoints(successSvc)
	resp, err := ep.savePlayerEndpoint(context.Background(), savePlayerRequest{Player: &jr})
	assert.Nil(t, err)
	spr, ok := resp.(savePlayerResponse)
	assert.True(t, ok)
	assert.Equal(t, "Jalen Ramsey", spr.Player.Name)
}

func TestMakeDeletePlayerEndpointShouldReturnFuncThatReturnsDeletePlayerResponse(t *testing.T) {
	ep := NewEndpoints(successSvc)
	resp, err := ep.deletePlayerEndpoint(context.Background(), deletePlayerRequest{ID: 20})
	assert.Nil(t, err)
	dpr, ok := resp.(deletePlayerResponse)
	assert.True(t, ok)
	assert.Equal(t, "", dpr.Err)
}

func TestMakeListPlayersEndpointShouldReturnFuncThatReturnsListPlayersResponseWithErrorWhenError(t *testing.T) {
	ep := NewEndpoints(failSvc)
	resp, err := ep.listPlayersEndpoint(context.Background(), listPlayersRequest{})
	assert.Nil(t, err)
	lpr, ok := resp.(listPlayersResponse)
	assert.True(t, ok)
	assert.Equal(t, "fail", lpr.Err)
}

func TestMakeGetPlayerEndpointShouldReturnFuncThatReturnsGetPlayerResponseWithErrorWhenError(t *testing.T) {
	ep := NewEndpoints(failSvc)
	resp, err := ep.getPlayerEndpoint(context.Background(), getPlayerRequest{})
	assert.Nil(t, err)
	gpr, ok := resp.(getPlayerResponse)
	assert.True(t, ok)
	assert.Equal(t, "fail", gpr.Err)
}

func TestMakeSavePlayerEndpointShouldReturnFuncThatReturnsSavePlayerResponseWithErrorWhenError(t *testing.T) {
	ep := NewEndpoints(failSvc)
	resp, err := ep.savePlayerEndpoint(context.Background(), savePlayerRequest{Player: &jr})
	assert.Nil(t, err)
	spr, ok := resp.(savePlayerResponse)
	assert.True(t, ok)
	assert.Equal(t, "fail", spr.Err)
}

func TestMakeDeletePlayerEndpointShouldReturnFuncThatReturnsDeletePlayerResponseWithErrorWhenError(t *testing.T) {
	ep := NewEndpoints(failSvc)
	resp, err := ep.deletePlayerEndpoint(context.Background(), deletePlayerRequest{ID: 20})
	assert.Nil(t, err)
	dpr, ok := resp.(deletePlayerResponse)
	assert.True(t, ok)
	assert.Equal(t, "fail", dpr.Err)
}
