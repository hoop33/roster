package players

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/hoop33/roster/models"
)

// Endpoints contains all the endpoints for the players service
type Endpoints struct {
	listPlayersEndpoint  endpoint.Endpoint
	getPlayerEndpoint    endpoint.Endpoint
	savePlayerEndpoint   endpoint.Endpoint
	deletePlayerEndpoint endpoint.Endpoint
}

type listPlayersRequest struct {
	Position string `json:"position,omitempty"`
}

type listPlayersResponse struct {
	Players []models.Player `json:"players,omitempty"`
	Err     string          `json:"error,omitempty"`
}

type getPlayerRequest struct {
	ID int `json:"id,omitempty"`
}

type getPlayerResponse struct {
	Player *models.Player `json:"player,omitempty"`
	Err    string         `json:"error,omitempty"`
}

type savePlayerRequest struct {
	Player *models.Player `json:"player,omitempty"`
}

type savePlayerResponse struct {
	Player  *models.Player `json:"player,omitempty"`
	Created bool           `json:"created,omitempty"`
	Err     string         `json:"error,omitempty"`
}

type deletePlayerRequest struct {
	ID int `json:"id,omitempty"`
}

type deletePlayerResponse struct {
	Err string `json:"error,omitempty"`
}

// NewEndpoints creates the endpoints
func NewEndpoints(s Service) *Endpoints {
	return &Endpoints{
		listPlayersEndpoint:  makeListPlayersEndpoint(s),
		getPlayerEndpoint:    makeGetPlayerEndpoint(s),
		savePlayerEndpoint:   makeSavePlayerEndpoint(s),
		deletePlayerEndpoint: makeDeletePlayerEndpoint(s),
	}
}

func makeListPlayersEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(listPlayersRequest)
		players, err := s.ListPlayers(ctx, req.Position)
		if err != nil {
			return listPlayersResponse{
				Err: err.Error(),
			}, nil
		}
		return listPlayersResponse{
			Players: players,
		}, nil
	}
}

func makeGetPlayerEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(getPlayerRequest)
		player, err := s.GetPlayer(ctx, req.ID)
		if err != nil {
			return getPlayerResponse{
				Err: err.Error(),
			}, nil
		}
		return getPlayerResponse{
			Player: player,
		}, nil
	}
}

func makeSavePlayerEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(savePlayerRequest)
		player, created, err := s.SavePlayer(ctx, req.Player)
		if err != nil {
			return savePlayerResponse{
				Err: err.Error(),
			}, nil
		}
		return savePlayerResponse{
			Player:  player,
			Created: created,
		}, nil
	}
}

func makeDeletePlayerEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(deletePlayerRequest)
		err := s.DeletePlayer(ctx, req.ID)
		if err != nil {
			return deletePlayerResponse{
				Err: err.Error(),
			}, nil
		}
		return deletePlayerResponse{}, nil
	}
}
