package models

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestPlayerStringShouldIncludeAllFields(t *testing.T) {
	player := &Player{
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
	assert.Equal(t, "[1] Blake Bortles (QB) -- #5, 6-5, 236lb, 26yo, 5exp -- Central Florida", player.String())
}

func TestListPlayersShouldReturnAllPlayers(t *testing.T) {
	db, mock, err := createDB()
	assert.Nil(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "number"}).
		AddRow(1, "Blake Bortles", "5").
		AddRow(2, "Jalen Ramsey", "20")

	mock.ExpectQuery(`^SELECT \* FROM players ORDER BY number ASC$`).
		WillReturnRows(rows)

	players, err := ListPlayers(db, "")
	assert.Nil(t, err)
	assert.Equal(t, 2, len(players))
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestListPlayersShouldReturnNoRowsWhenNoPlayers(t *testing.T) {
	db, mock, err := createDB()
	assert.Nil(t, err)
	defer db.Close()

	mock.ExpectQuery(`^SELECT \* FROM players ORDER BY number ASC$`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "number"}))

	players, err := ListPlayers(db, "")
	assert.Error(t, sql.ErrNoRows, err)
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

	players, err := ListPlayers(db, "QB")
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

	players, err := ListPlayers(db, "")
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

	player, err := GetPlayer(db, 1)
	assert.Nil(t, err)
	assert.Equal(t, "Blake Bortles", player.Name)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestGetPlayerShouldReturnErrorWhenNotExists(t *testing.T) {
	db, mock, err := createDB()
	assert.Nil(t, err)
	defer db.Close()

	mock.ExpectQuery(`^SELECT \* FROM players WHERE id = \$1$`).
		WithArgs(1).
		WillReturnError(sql.ErrNoRows)

	player, err := GetPlayer(db, 1)
	assert.NotNil(t, err)
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

	player, err := GetPlayer(db, 1)
	assert.NotNil(t, err)
	assert.Nil(t, player)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestPlayerSaveShouldCreatePlayerWhenNoID(t *testing.T) {
	db, mock, err := createDB()
	assert.Nil(t, err)
	defer db.Close()

	p := &Player{
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

	_, created, err := p.Save(db)
	assert.Nil(t, err)
	assert.True(t, created)
	assert.Equal(t, 1, p.ID)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestPlayerSaveShouldReturnErrorWhenNoIDAndDatabaseError(t *testing.T) {
	db, mock, err := createDB()
	assert.Nil(t, err)
	defer db.Close()

	p := &Player{
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

	_, created, err := p.Save(db)
	assert.NotNil(t, err)
	assert.False(t, created)
	assert.Equal(t, 0, p.ID)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestPlayerSaveShouldUpdatePlayerWhenHasID(t *testing.T) {
	db, mock, err := createDB()
	assert.Nil(t, err)
	defer db.Close()

	p := &Player{
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

	_, created, err := p.Save(db)
	assert.Nil(t, err)
	assert.False(t, created)
	assert.Equal(t, 1, p.ID)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestPlayerSaveShouldReturnErrWhenUpdatePlayerNotFound(t *testing.T) {
	db, mock, err := createDB()
	assert.Nil(t, err)
	defer db.Close()

	p := &Player{
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

	_, created, err := p.Save(db)
	assert.Error(t, sql.ErrNoRows, err)
	assert.False(t, created)
	assert.Equal(t, 1, p.ID)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestPlayerSaveShouldReturnErrorWhenHasIDAndDatabaseError(t *testing.T) {
	db, mock, err := createDB()
	assert.Nil(t, err)
	defer db.Close()

	p := &Player{
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

	_, created, err := p.Save(db)
	assert.NotNil(t, err)
	assert.False(t, created)
	assert.Equal(t, 1, p.ID)
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

	player := &Player{
		ID: 1,
	}
	err = player.Delete(db)
	assert.Nil(t, err)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestDeleteShouldReturnNoRowsWhenNotFound(t *testing.T) {
	db, mock, err := createDB()
	assert.Nil(t, err)
	defer db.Close()

	mock.ExpectExec(`^DELETE FROM players
		WHERE id=\$1$`).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 0))

	player := &Player{
		ID: 1,
	}
	err = player.Delete(db)
	assert.Error(t, sql.ErrNoRows, err)
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

	player := &Player{
		ID: 1,
	}
	err = player.Delete(db)
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
