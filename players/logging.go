package players

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/hoop33/roster/models"
)

type loggingService struct {
	logger log.Logger
	next   Service
}

// NewLoggingService returns a new logging service
func NewLoggingService(logger log.Logger, next Service) Service {
	return &loggingService{
		logger: logger,
		next:   next,
	}
}

func (l *loggingService) ListPlayers(ctx context.Context, position string) (players []models.Player, err error) {
	defer func(begin time.Time) {
		l.logger.Log("msg", "listing players", "pos", position, "num", len(players), "err", err, "took", time.Since(begin))
	}(time.Now())
	return l.next.ListPlayers(ctx, position)
}

func (l *loggingService) GetPlayer(ctx context.Context, id int) (player *models.Player, err error) {
	defer func(begin time.Time) {
		l.logger.Log("msg", "getting a player", "id", id, "err", err, "took", time.Since(begin))
	}(time.Now())
	return l.next.GetPlayer(ctx, id)
}

func (l *loggingService) SavePlayer(ctx context.Context, player *models.Player) (p *models.Player, created bool, err error) {
	defer func(begin time.Time) {
		l.logger.Log("msg", "saving a player", "id", player.ID, "name", player.Name, "created", created, "err", err, "took", time.Since(begin))
	}(time.Now())
	return l.next.SavePlayer(ctx, player)
}

func (l *loggingService) DeletePlayer(ctx context.Context, id int) (err error) {
	defer func(begin time.Time) {
		l.logger.Log("msg", "deleting a player", "id", id, "err", err, "took", time.Since(begin))
	}(time.Now())
	return l.next.DeletePlayer(ctx, id)
}
