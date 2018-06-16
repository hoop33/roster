package models

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

// Player is a player on the roster
type Player struct {
	ID         int    `db:"id" json:"id"`
	Name       string `db:"name" json:"name"`
	Number     string `db:"number" json:"number"`
	Position   string `db:"position" json:"position"`
	Height     string `db:"height" json:"height"`
	Weight     string `db:"weight" json:"weight"`
	Age        string `db:"age" json:"age"`
	Experience int    `db:"experience" json:"experience"`
	College    string `db:"college" json:"college"`
}

// String returns a String version of a player
func (p *Player) String() string {
	return fmt.Sprintf("[%d] %s (%s) -- #%s, %s, %slb, %syo, %dexp -- %s",
		p.ID,
		p.Name,
		p.Position,
		p.Number,
		p.Height,
		p.Weight,
		p.Age,
		p.Experience,
		p.College)
}

// ListPlayers lists all the players, optionally restricted to a position
func ListPlayers(db *sqlx.DB, position string) ([]Player, error) {
	var players []Player
	var err error
	if position == "" {
		err = db.Select(&players, "SELECT * FROM players ORDER BY number ASC")
	} else {
		err = db.Select(&players, "SELECT * FROM players WHERE position = $1 ORDER BY number ASC", position)
	}
	if err != nil {
		return nil, err
	}
	if len(players) == 0 {
		return nil, sql.ErrNoRows
	}
	return players, nil
}

// GetPlayer gets a player by ID
func GetPlayer(db *sqlx.DB, id int) (*Player, error) {
	player := Player{}
	err := db.Get(&player, "SELECT * FROM players WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	return &player, nil
}

// Save saves a player (insert or update)
func (p *Player) Save(db *sqlx.DB) (*Player, bool, error) {
	created := false
	var err error
	if p.ID <= 0 {
		err = p.create(db)
		created = err == nil
	} else {
		err = p.update(db)
	}
	return p, created, err
}

// Delete deletes a player
func (p *Player) Delete(db *sqlx.DB) error {
	result, err := db.Exec(`DELETE FROM players
		WHERE id=$1`,
		p.ID)
	if err != nil {
		return err
	}

	count, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if count != 1 {
		return sql.ErrNoRows
	}

	return nil
}

func (p *Player) create(db *sqlx.DB) error {
	var id int
	err := db.QueryRow(`INSERT INTO players
		(name, number, position, height, weight, age, experience, college)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id`,
		p.Name, p.Number, p.Position, p.Height, p.Weight, p.Age, p.Experience, p.College).
		Scan(&id)
	if err != nil {
		return err
	}

	p.ID = id
	return nil
}

func (p *Player) update(db *sqlx.DB) error {
	result, err := db.Exec(`UPDATE players
		SET name=$1, number=$2, position=$3, height=$4, weight=$5, age=$6, experience=$7, college=$8
		WHERE id=$9`,
		p.Name, p.Number, p.Position, p.Height, p.Weight, p.Age, p.Experience, p.College, p.ID)
	if err != nil {
		return err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return sql.ErrNoRows
	}
	return nil
}
