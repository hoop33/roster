package players

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/hoop33/roster/models"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestListPlayersShouldReturnAllPlayers(t *testing.T) {
	db, mock, err := createDB()
	assert.Nil(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "number"}).
		AddRow(1, "Blake Bortles", "5").
		AddRow(2, "Jalen Ramsey", "20")

	mock.ExpectQuery(`^SELECT \* FROM players ORDER BY number ASC$`).
		WillReturnRows(rows)

	players, err := NewService(db).ListPlayers(context.Background(), "")
	assert.Nil(t, err)
	assert.Equal(t, 2, len(players))
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestListPlayersShouldReturnNotFoundWhenNoRows(t *testing.T) {
	db, mock, err := createDB()
	assert.Nil(t, err)
	defer db.Close()

	mock.ExpectQuery(`^SELECT \* FROM players ORDER BY number ASC$`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "number"}))

	players, err := NewService(db).ListPlayers(context.Background(), "")
	assert.Error(t, errNotFound, err)
	assert.Equal(t, 0, len(players))
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestListPlayersShouldUsePositionWhenSpecified(t *testing.T) {
	db, mock, err := createDB()
	assert.Nil(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "number", "position"}).
		AddRow(1, "Blake Bortles", "5", "QB").
		AddRow(2, "Cody Kessler", "6", "QB")

	mock.ExpectQuery(`^SELECT \* FROM players WHERE position = \$1 ORDER BY number ASC$`).
		WithArgs("QB").
		WillReturnRows(rows)

	players, err := NewService(db).ListPlayers(context.Background(), "QB")
	assert.Nil(t, err)
	assert.Equal(t, 2, len(players))
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestListPlayersShouldReturnErrorWhenDatabaseError(t *testing.T) {
	db, mock, err := createDB()
	assert.Nil(t, err)
	defer db.Close()

	mock.ExpectQuery(`^SELECT \* FROM players ORDER BY number ASC`).
		WillReturnError(errors.New("database error"))

	players, err := NewService(db).ListPlayers(context.Background(), "")
	assert.NotNil(t, err)
	assert.Equal(t, 0, len(players))
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestGetPlayerShouldReturnPlayerWhenExists(t *testing.T) {
	db, mock, err := createDB()
	assert.Nil(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "number"}).
		AddRow(1, "Blake Bortles", "5")

	mock.ExpectQuery(`^SELECT \* FROM players WHERE id = \$1$`).
		WithArgs(1).
		WillReturnRows(rows)

	player, err := NewService(db).GetPlayer(context.Background(), 1)
	assert.Nil(t, err)
	assert.Equal(t, "Blake Bortles", player.Name)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestGetPlayerShouldReturnErrorWhenNotFound(t *testing.T) {
	db, mock, err := createDB()
	assert.Nil(t, err)
	defer db.Close()

	mock.ExpectQuery(`^SELECT \* FROM players WHERE id = \$1$`).
		WithArgs(1).
		WillReturnError(sql.ErrNoRows)

	player, err := NewService(db).GetPlayer(context.Background(), 1)
	assert.Error(t, errNotFound, err)
	assert.Nil(t, player)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestGetPlayerShouldReturnErrorWhenDatabaseError(t *testing.T) {
	db, mock, err := createDB()
	assert.Nil(t, err)
	defer db.Close()

	mock.ExpectQuery(`^SELECT \* FROM players WHERE id = \$1$`).
		WithArgs(1).
		WillReturnError(errors.New("database error"))

	player, err := NewService(db).GetPlayer(context.Background(), 1)
	assert.NotNil(t, err)
	assert.Nil(t, player)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestSavePlayerShouldCreatePlayerWhenNoID(t *testing.T) {
	db, mock, err := createDB()
	assert.Nil(t, err)
	defer db.Close()

	p := &models.Player{
		Name:       "Blake Bortles",
		Number:     "5",
		Position:   "QB",
		Height:     "6-5",
		Weight:     "236",
		Age:        "26",
		Experience: 5,
		College:    "Central Florida",
	}
	mock.ExpectQuery(`^INSERT INTO players
		\(name, number, position, height, weight, age, experience, college\)
		VALUES \(\$1, \$2, \$3, \$4, \$5, \$6, \$7, \$8\)
		RETURNING id$`).
		WithArgs(p.Name, p.Number, p.Position, p.Height, p.Weight, p.Age, p.Experience, p.College).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	player, created, err := NewService(db).SavePlayer(context.Background(), p)
	assert.Nil(t, err)
	assert.True(t, created)
	assert.Equal(t, 1, player.ID)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestSavePlayerShouldReturnErrorWhenNoIDAndDatabaseError(t *testing.T) {
	db, mock, err := createDB()
	assert.Nil(t, err)
	defer db.Close()

	p := &models.Player{
		Name:       "Blake Bortles",
		Number:     "5",
		Position:   "QB",
		Height:     "6-5",
		Weight:     "236",
		Age:        "26",
		Experience: 5,
		College:    "Central Florida",
	}
	mock.ExpectQuery(`^INSERT INTO players
		\(name, number, position, height, weight, age, experience, college\)
		VALUES \(\$1, \$2, \$3, \$4, \$5, \$6, \$7, \$8\)
		RETURNING id$`).
		WithArgs(p.Name, p.Number, p.Position, p.Height, p.Weight, p.Age, p.Experience, p.College).
		WillReturnError(errors.New("database error"))

	player, created, err := NewService(db).SavePlayer(context.Background(), p)
	assert.NotNil(t, err)
	assert.False(t, created)
	assert.Equal(t, 0, player.ID)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestSavePlayerShouldUpdatePlayerWhenHasID(t *testing.T) {
	db, mock, err := createDB()
	assert.Nil(t, err)
	defer db.Close()

	p := &models.Player{
		ID:         1,
		Name:       "Blake Bortles",
		Number:     "5",
		Position:   "QB",
		Height:     "6-5",
		Weight:     "236",
		Age:        "26",
		Experience: 5,
		College:    "Central Florida",
	}
	mock.ExpectExec(`^UPDATE players
		SET name=\$1, number=\$2, position=\$3, height=\$4, weight=\$5, age=\$6, experience=\$7, college=\$8
		WHERE id=\$9$`).
		WithArgs(p.Name, p.Number, p.Position, p.Height, p.Weight, p.Age, p.Experience, p.College, p.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	player, created, err := NewService(db).SavePlayer(context.Background(), p)
	assert.Nil(t, err)
	assert.False(t, created)
	assert.Equal(t, 1, player.ID)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestSavePlayerShouldReturnErrWhenUpdatePlayerNotFound(t *testing.T) {
	db, mock, err := createDB()
	assert.Nil(t, err)
	defer db.Close()

	p := &models.Player{
		ID:         1,
		Name:       "Blake Bortles",
		Number:     "5",
		Position:   "QB",
		Height:     "6-5",
		Weight:     "236",
		Age:        "26",
		Experience: 5,
		College:    "Central Florida",
	}
	mock.ExpectExec(`^UPDATE players
		SET name=\$1, number=\$2, position=\$3, height=\$4, weight=\$5, age=\$6, experience=\$7, college=\$8
		WHERE id=\$9$`).
		WithArgs(p.Name, p.Number, p.Position, p.Height, p.Weight, p.Age, p.Experience, p.College, p.ID).
		WillReturnResult(sqlmock.NewResult(0, 0))

	player, created, err := NewService(db).SavePlayer(context.Background(), p)
	assert.Error(t, errNotFound, err)
	assert.False(t, created)
	assert.Nil(t, player)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestSavePlayerShouldReturnErrorWhenHasIDAndDatabaseError(t *testing.T) {
	db, mock, err := createDB()
	assert.Nil(t, err)
	defer db.Close()

	p := &models.Player{
		ID:         1,
		Name:       "Blake Bortles",
		Number:     "5",
		Position:   "QB",
		Height:     "6-5",
		Weight:     "236",
		Age:        "26",
		Experience: 5,
		College:    "Central Florida",
	}
	mock.ExpectExec(`^UPDATE players
		SET name=\$1, number=\$2, position=\$3, height=\$4, weight=\$5, age=\$6, experience=\$7, college=\$8
		WHERE id=\$9$`).
		WithArgs(p.Name, p.Number, p.Position, p.Height, p.Weight, p.Age, p.Experience, p.College, p.ID).
		WillReturnError(errors.New("database error"))

	player, created, err := NewService(db).SavePlayer(context.Background(), p)
	assert.NotNil(t, err)
	assert.False(t, created)
	assert.Equal(t, 1, player.ID)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestDeleteShouldDeletePlayer(t *testing.T) {
	db, mock, err := createDB()
	assert.Nil(t, err)
	defer db.Close()

	mock.ExpectExec(`^DELETE FROM players
		WHERE id=\$1$`).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = NewService(db).DeletePlayer(context.Background(), 1)
	assert.Nil(t, err)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestDeleteShouldReturnErrorWhenNotFound(t *testing.T) {
	db, mock, err := createDB()
	assert.Nil(t, err)
	defer db.Close()

	mock.ExpectExec(`^DELETE FROM players
		WHERE id=\$1$`).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err = NewService(db).DeletePlayer(context.Background(), 1)
	assert.Error(t, errNotFound, err)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestDeleteShouldReturnErrorWhenError(t *testing.T) {
	db, mock, err := createDB()
	assert.Nil(t, err)
	defer db.Close()

	mock.ExpectExec(`^DELETE FROM players
		WHERE id=\$1$`).
		WithArgs(1).
		WillReturnError(errors.New("database error"))

	err = NewService(db).DeletePlayer(context.Background(), 1)
	assert.NotNil(t, err)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func createDB() (*sqlx.DB, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}

	xdb := sqlx.NewDb(db, "sqlmock")
	return xdb, mock, nil
}
