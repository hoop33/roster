package players

import (
	"context"
	"database/sql"
	"errors"

	"github.com/hoop33/roster/models"
	"github.com/jmoiron/sqlx"
)

// Service defines the functions for a players service
type Service interface {
	ListPlayers(context.Context, string) ([]models.Player, error)
	GetPlayer(context.Context, int) (*models.Player, error)
	SavePlayer(context.Context, *models.Player) (*models.Player, bool, error)
	DeletePlayer(context.Context, int) error
}

type service struct {
	db *sqlx.DB
}

var errNotFound = errors.New("not found")

// NewService returns a new service for interacting with players
func NewService(db *sqlx.DB) Service {
	return &service{
		db: db,
	}
}

func (p *service) ListPlayers(_ context.Context, position string) ([]models.Player, error) {
	players, err := models.ListPlayers(p.db, position)
	if err == sql.ErrNoRows {
		return nil, errNotFound
	}
	return players, err
}

func (p *service) GetPlayer(_ context.Context, id int) (*models.Player, error) {
	player, err := models.GetPlayer(p.db, id)
	if err == sql.ErrNoRows {
		return nil, errNotFound
	}
	return player, err
}

func (p *service) SavePlayer(_ context.Context, player *models.Player) (*models.Player, bool, error) {
	player, created, err := player.Save(p.db)
	if err == sql.ErrNoRows {
		return nil, false, errNotFound
	}
	return player, created, err
}

func (p *service) DeletePlayer(_ context.Context, id int) error {
	player := &models.Player{
		ID: id,
	}

	err := player.Delete(p.db)
	if err == sql.ErrNoRows {
		return errNotFound
	}
	return err
}
