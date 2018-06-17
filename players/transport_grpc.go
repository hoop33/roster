package players

import (
	"context"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport/grpc"
	"github.com/hoop33/roster/models"
	"github.com/hoop33/roster/pb"
)

type grpcTransport struct {
	listPlayers  grpc.Handler
	getPlayer    grpc.Handler
	savePlayer   grpc.Handler
	deletePlayer grpc.Handler
}

// NewGRPCTransport returns a handler for GRPC transport
func NewGRPCTransport(ep *Endpoints, logger log.Logger) pb.PlayersServer {
	opts := []grpc.ServerOption{
		grpc.ServerErrorLogger(log.With(logger, "tag", "grpc")),
	}

	return &grpcTransport{
		listPlayers: grpc.NewServer(
			ep.listPlayersEndpoint,
			decodeGRPCListPlayersRequest,
			encodeGRPCListPlayersResponse,
			opts...,
		),
		getPlayer: grpc.NewServer(
			ep.getPlayerEndpoint,
			decodeGRPCGetPlayerRequest,
			encodeGRPCGetPlayerResponse,
			opts...,
		),
		savePlayer: grpc.NewServer(
			ep.savePlayerEndpoint,
			decodeGRPCSavePlayerRequest,
			encodeGRPCSavePlayerResponse,
			opts...,
		),
		deletePlayer: grpc.NewServer(
			ep.deletePlayerEndpoint,
			decodeGRPCDeletePlayerRequest,
			encodeGRPCDeletePlayerResponse,
			opts...,
		),
	}
}

func (s *grpcTransport) ListPlayers(ctx context.Context, r *pb.ListPlayersRequest) (*pb.ListPlayersResponse, error) {
	_, resp, err := s.listPlayers.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.ListPlayersResponse), nil
}

func (s *grpcTransport) GetPlayer(ctx context.Context, r *pb.GetPlayerRequest) (*pb.GetPlayerResponse, error) {
	_, resp, err := s.getPlayer.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.GetPlayerResponse), nil
}

func (s *grpcTransport) SavePlayer(ctx context.Context, r *pb.SavePlayerRequest) (*pb.SavePlayerResponse, error) {
	_, resp, err := s.savePlayer.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.SavePlayerResponse), nil
}

func (s *grpcTransport) DeletePlayer(ctx context.Context, r *pb.DeletePlayerRequest) (*pb.DeletePlayerResponse, error) {
	_, resp, err := s.deletePlayer.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.DeletePlayerResponse), nil
}

func decodeGRPCListPlayersRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.ListPlayersRequest)
	return listPlayersRequest{
		Position: req.Position,
	}, nil
}

func encodeGRPCListPlayersResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(listPlayersResponse)

	players := make([]*pb.Player, len(resp.Players))
	for i, p := range resp.Players {
		player := modelsPlayerToProtoPlayer(p)
		players[i] = &player
	}

	return &pb.ListPlayersResponse{
		Players: players,
		Err:     resp.Err,
	}, nil
}

func decodeGRPCGetPlayerRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.GetPlayerRequest)
	return getPlayerRequest{
		ID: int(req.Id),
	}, nil
}

func encodeGRPCGetPlayerResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(getPlayerResponse)

	if resp.Player == nil {
		return &pb.GetPlayerResponse{
			Err: resp.Err,
		}, nil
	}

	player := modelsPlayerToProtoPlayer(*resp.Player)
	return &pb.GetPlayerResponse{
		Player: &player,
		Err:    resp.Err,
	}, nil
}

func decodeGRPCSavePlayerRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.SavePlayerRequest)
	player := protoPlayerToModelsPlayer(*req.Player)
	return savePlayerRequest{
		Player: &player,
	}, nil
}

func encodeGRPCSavePlayerResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(savePlayerResponse)

	if resp.Player == nil {
		return &pb.SavePlayerResponse{
			Created: resp.Created,
			Err:     resp.Err,
		}, nil
	}

	player := modelsPlayerToProtoPlayer(*resp.Player)
	return &pb.SavePlayerResponse{
		Player:  &player,
		Created: resp.Created,
		Err:     resp.Err,
	}, nil
}

func decodeGRPCDeletePlayerRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.DeletePlayerRequest)
	return deletePlayerRequest{
		ID: int(req.Id),
	}, nil
}

func encodeGRPCDeletePlayerResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(deletePlayerResponse)

	return &pb.DeletePlayerResponse{
		Err: resp.Err,
	}, nil
}

func modelsPlayerToProtoPlayer(p models.Player) pb.Player {
	return pb.Player{
		Id:         int32(p.ID),
		Name:       p.Name,
		Number:     p.Number,
		Position:   p.Position,
		Height:     p.Height,
		Weight:     p.Weight,
		Age:        p.Age,
		Experience: int32(p.Experience),
		College:    p.College,
	}
}

func protoPlayerToModelsPlayer(p pb.Player) models.Player {
	return models.Player{
		ID:         int(p.Id),
		Name:       p.Name,
		Number:     p.Number,
		Position:   p.Position,
		Height:     p.Height,
		Weight:     p.Weight,
		Age:        p.Age,
		Experience: int(p.Experience),
		College:    p.College,
	}
}
