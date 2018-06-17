package players

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/hoop33/roster/models"
	"github.com/hoop33/roster/pb"
	"github.com/stretchr/testify/assert"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestGRPCListPlayersShouldReturnPlayersWhenDatabaseReturnsPlayers(t *testing.T) {
	db, mock, err := createDB()
	assert.Nil(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "number"}).
		AddRow(1, "Blake Bortles", "5").
		AddRow(2, "Jalen Ramsey", "20")

	mock.ExpectQuery(`^SELECT \* FROM players ORDER BY number ASC$`).
		WillReturnRows(rows)

	es := NewEndpoints(NewService(db))

	tr := NewGRPCTransport(es, log.NewNopLogger())
	req := &pb.ListPlayersRequest{}
	resp, err := tr.ListPlayers(context.Background(), req)
	assert.Nil(t, err)
	players := resp.GetPlayers()
	assert.Equal(t, 2, len(players))
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestGRPCListPlayersShouldReturnErrorWhenDatabaseReturnsNoPlayers(t *testing.T) {
	db, mock, err := createDB()
	assert.Nil(t, err)
	defer db.Close()

	mock.ExpectQuery(`^SELECT \* FROM players ORDER BY number ASC$`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "number"}))

	es := NewEndpoints(NewService(db))

	tr := NewGRPCTransport(es, log.NewNopLogger())
	req := &pb.ListPlayersRequest{}
	resp, err := tr.ListPlayers(context.Background(), req)
	assert.Error(t, errNotFound, err)
	players := resp.GetPlayers()
	assert.Equal(t, 0, len(players))
	assert.Equal(t, "not found", resp.GetErr())
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestGRPCListPlayersShouldReturnErrorWhenDatabaseReturnsError(t *testing.T) {
	db, mock, err := createDB()
	assert.Nil(t, err)
	defer db.Close()

	mock.ExpectQuery(`^SELECT \* FROM players ORDER BY number ASC$`).
		WillReturnError(errors.New("database error"))

	es := NewEndpoints(NewService(db))

	tr := NewGRPCTransport(es, log.NewNopLogger())
	req := &pb.ListPlayersRequest{}
	resp, err := tr.ListPlayers(context.Background(), req)
	assert.Error(t, errNotFound, err)
	players := resp.GetPlayers()
	assert.Equal(t, 0, len(players))
	assert.Equal(t, "database error", resp.GetErr())
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestGRPCGetPlayerShouldReturnPlayerWhenDatabaseReturnsPlayer(t *testing.T) {
	db, mock, err := createDB()
	assert.Nil(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "number"}).
		AddRow(1, "Blake Bortles", "5")

	mock.ExpectQuery(`^SELECT \* FROM players WHERE id = \$1$`).
		WithArgs(1).
		WillReturnRows(rows)

	es := NewEndpoints(NewService(db))

	tr := NewGRPCTransport(es, log.NewNopLogger())
	req := &pb.GetPlayerRequest{
		Id: 1,
	}
	resp, err := tr.GetPlayer(context.Background(), req)
	assert.Nil(t, err)
	assert.Equal(t, "Blake Bortles", resp.GetPlayer().GetName())
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestGRPCGetPlayerShouldReturnErrorWhenDatabaseReturnsNoRows(t *testing.T) {
	db, mock, err := createDB()
	assert.Nil(t, err)
	defer db.Close()

	mock.ExpectQuery(`^SELECT \* FROM players WHERE id = \$1$`).
		WithArgs(1).
		WillReturnError(sql.ErrNoRows)

	es := NewEndpoints(NewService(db))

	tr := NewGRPCTransport(es, log.NewNopLogger())
	req := &pb.GetPlayerRequest{
		Id: 1,
	}
	resp, err := tr.GetPlayer(context.Background(), req)
	assert.Nil(t, err)
	assert.Equal(t, "not found", resp.GetErr())
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestGRPCGetPlayerShouldReturnErrorWhenDatabaseReturnsError(t *testing.T) {
	db, mock, err := createDB()
	assert.Nil(t, err)
	defer db.Close()

	mock.ExpectQuery(`^SELECT \* FROM players WHERE id = \$1$`).
		WithArgs(1).
		WillReturnError(errors.New("database error"))

	es := NewEndpoints(NewService(db))

	tr := NewGRPCTransport(es, log.NewNopLogger())
	req := &pb.GetPlayerRequest{
		Id: 1,
	}
	resp, err := tr.GetPlayer(context.Background(), req)
	assert.Nil(t, err)
	assert.Equal(t, "database error", resp.GetErr())
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestGRPCSavePlayerShouldCreatePlayerWhenNoID(t *testing.T) {
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

	es := NewEndpoints(NewService(db))

	player := modelsPlayerToProtoPlayer(*p)
	tr := NewGRPCTransport(es, log.NewNopLogger())
	req := &pb.SavePlayerRequest{
		Player: &player,
	}
	resp, err := tr.SavePlayer(context.Background(), req)
	assert.Nil(t, err)
	assert.Equal(t, "Blake Bortles", resp.GetPlayer().GetName())
	assert.True(t, resp.GetCreated())
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestGRPCSavePlayerShouldUpdatePlayerWhenHasID(t *testing.T) {
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

	es := NewEndpoints(NewService(db))

	player := modelsPlayerToProtoPlayer(*p)
	tr := NewGRPCTransport(es, log.NewNopLogger())
	req := &pb.SavePlayerRequest{
		Player: &player,
	}
	resp, err := tr.SavePlayer(context.Background(), req)
	assert.Nil(t, err)
	assert.Equal(t, "Blake Bortles", resp.GetPlayer().GetName())
	assert.False(t, resp.GetCreated())
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestGRPCSavePlayerShouldReturnErrorWhenCreateFails(t *testing.T) {
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

	es := NewEndpoints(NewService(db))

	player := modelsPlayerToProtoPlayer(*p)
	tr := NewGRPCTransport(es, log.NewNopLogger())
	req := &pb.SavePlayerRequest{
		Player: &player,
	}
	resp, err := tr.SavePlayer(context.Background(), req)
	assert.Nil(t, err)
	assert.Equal(t, "database error", resp.GetErr())
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestGRPCSavePlayerShouldReturnErrorWhenNotFound(t *testing.T) {
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
		WillReturnError(sql.ErrNoRows)

	es := NewEndpoints(NewService(db))

	player := modelsPlayerToProtoPlayer(*p)
	tr := NewGRPCTransport(es, log.NewNopLogger())
	req := &pb.SavePlayerRequest{
		Player: &player,
	}
	resp, err := tr.SavePlayer(context.Background(), req)
	assert.Nil(t, err)
	assert.Equal(t, "not found", resp.GetErr())
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestGRPCSavePlayerShouldReturnErrorWhenDatabaseReturnsError(t *testing.T) {
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

	es := NewEndpoints(NewService(db))

	player := modelsPlayerToProtoPlayer(*p)
	tr := NewGRPCTransport(es, log.NewNopLogger())
	req := &pb.SavePlayerRequest{
		Player: &player,
	}
	resp, err := tr.SavePlayer(context.Background(), req)
	assert.Nil(t, err)
	assert.Equal(t, "database error", resp.GetErr())
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestGRPCDeletePlayerShouldReturnNoErrorWhenDeleteSucceeds(t *testing.T) {
	db, mock, err := createDB()
	assert.Nil(t, err)
	defer db.Close()

	mock.ExpectExec(`^DELETE FROM players
		WHERE id=\$1$`).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	es := NewEndpoints(NewService(db))

	tr := NewGRPCTransport(es, log.NewNopLogger())
	req := &pb.DeletePlayerRequest{
		Id: 1,
	}
	resp, err := tr.DeletePlayer(context.Background(), req)
	assert.Nil(t, err)
	assert.Equal(t, "", resp.GetErr())
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestGRPCDeletePlayerShouldReturnErrorWhenNotFound(t *testing.T) {
	db, mock, err := createDB()
	assert.Nil(t, err)
	defer db.Close()

	mock.ExpectExec(`^DELETE FROM players
		WHERE id=\$1$`).
		WithArgs(1).
		WillReturnError(sql.ErrNoRows)

	es := NewEndpoints(NewService(db))

	tr := NewGRPCTransport(es, log.NewNopLogger())
	req := &pb.DeletePlayerRequest{
		Id: 1,
	}
	resp, err := tr.DeletePlayer(context.Background(), req)
	assert.Nil(t, err)
	assert.Equal(t, "not found", resp.GetErr())
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestGRPCDeletePlayerShouldReturnErrorWhenDatabaseReturnsError(t *testing.T) {
	db, mock, err := createDB()
	assert.Nil(t, err)
	defer db.Close()

	mock.ExpectExec(`^DELETE FROM players
		WHERE id=\$1$`).
		WithArgs(1).
		WillReturnError(errors.New("database error"))

	es := NewEndpoints(NewService(db))

	tr := NewGRPCTransport(es, log.NewNopLogger())
	req := &pb.DeletePlayerRequest{
		Id: 1,
	}
	resp, err := tr.DeletePlayer(context.Background(), req)
	assert.Nil(t, err)
	assert.Equal(t, "database error", resp.GetErr())
	assert.Nil(t, mock.ExpectationsWereMet())
}
