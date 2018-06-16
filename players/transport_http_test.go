package players

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/hoop33/roster/models"
	"github.com/stretchr/testify/assert"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestHTTPListPlayersShouldReturnPlayersWhenDatabaseReturnsPlayers(t *testing.T) {
	db, mock, err := createDB()
	assert.Nil(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "number"}).
		AddRow(1, "Blake Bortles", "5").
		AddRow(2, "Jalen Ramsey", "20")

	mock.ExpectQuery(`^SELECT \* FROM players ORDER BY number ASC$`).
		WillReturnRows(rows)

	es := NewEndpoints(NewService(db))

	req := httptest.NewRequest("GET", "/v1/players", nil)
	resp := httptest.NewRecorder()
	NewHTTPTransport(es, log.NewNopLogger()).ServeHTTP(resp, req)
	body, err := ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)
	assert.True(t, strings.Contains(string(body), "Bortles"))
	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestHTTPListPlayersShouldReturnErrorWhenDatabaseReturnsNoRows(t *testing.T) {
	db, mock, err := createDB()
	assert.Nil(t, err)
	defer db.Close()

	mock.ExpectQuery(`^SELECT \* FROM players ORDER BY number ASC$`).
		WillReturnError(sql.ErrNoRows)

	es := NewEndpoints(NewService(db))

	req := httptest.NewRequest("GET", "/v1/players", nil)
	resp := httptest.NewRecorder()
	NewHTTPTransport(es, log.NewNopLogger()).ServeHTTP(resp, req)
	body, err := ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)
	assert.True(t, strings.Contains(string(body), `"error":"not found"`))
	assert.Equal(t, http.StatusNotFound, resp.Code)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestHTTPListPlayersShouldReturnErrorWhenDatabaseReturnsError(t *testing.T) {
	db, mock, err := createDB()
	assert.Nil(t, err)
	defer db.Close()

	mock.ExpectQuery(`^SELECT \* FROM players ORDER BY number ASC$`).
		WillReturnError(errors.New("database error"))

	es := NewEndpoints(NewService(db))

	req := httptest.NewRequest("GET", "/v1/players", nil)
	resp := httptest.NewRecorder()
	NewHTTPTransport(es, log.NewNopLogger()).ServeHTTP(resp, req)
	body, err := ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)
	assert.True(t, strings.Contains(string(body), `"error":"database error"`))
	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestHTTPGetPlayerShouldReturnPlayerWhenExists(t *testing.T) {
	db, mock, err := createDB()
	assert.Nil(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "number"}).
		AddRow(1, "Blake Bortles", "5")

	mock.ExpectQuery(`^SELECT \* FROM players WHERE id = \$1$`).
		WithArgs(1).
		WillReturnRows(rows)

	es := NewEndpoints(NewService(db))

	req := httptest.NewRequest("GET", "/v1/players/1", nil)
	resp := httptest.NewRecorder()
	NewHTTPTransport(es, log.NewNopLogger()).ServeHTTP(resp, req)
	body, err := ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)
	assert.True(t, strings.Contains(string(body), "Bortles"))
	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestHTTPGetPlayerShouldReturnErrorWhenNotExists(t *testing.T) {
	db, mock, err := createDB()
	assert.Nil(t, err)
	defer db.Close()

	mock.ExpectQuery(`^SELECT \* FROM players WHERE id = \$1$`).
		WithArgs(1).
		WillReturnError(sql.ErrNoRows)

	es := NewEndpoints(NewService(db))

	req := httptest.NewRequest("GET", "/v1/players/1", nil)
	resp := httptest.NewRecorder()
	NewHTTPTransport(es, log.NewNopLogger()).ServeHTTP(resp, req)
	body, err := ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)
	assert.True(t, strings.Contains(string(body), `"error":"not found"`))
	assert.Equal(t, http.StatusNotFound, resp.Code)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestHTTPGetPlayerShouldReturnErrorWhenDatabaseReturnsError(t *testing.T) {
	db, mock, err := createDB()
	assert.Nil(t, err)
	defer db.Close()

	mock.ExpectQuery(`^SELECT \* FROM players WHERE id = \$1$`).
		WithArgs(1).
		WillReturnError(errors.New("database error"))

	es := NewEndpoints(NewService(db))

	req := httptest.NewRequest("GET", "/v1/players/1", nil)
	resp := httptest.NewRecorder()
	NewHTTPTransport(es, log.NewNopLogger()).ServeHTTP(resp, req)
	body, err := ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)
	assert.True(t, strings.Contains(string(body), `"error":"database error"`))
	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestHTTPGetPlayerShouldReturnErrorWhenIDIsNonNumeric(t *testing.T) {
	db, mock, err := createDB()
	assert.Nil(t, err)
	defer db.Close()

	es := NewEndpoints(NewService(db))

	req := httptest.NewRequest("GET", "/v1/players/a", nil)
	resp := httptest.NewRecorder()
	NewHTTPTransport(es, log.NewNopLogger()).ServeHTTP(resp, req)
	body, err := ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)
	assert.True(t, strings.Contains(string(body), `"error":"bad request"`))
	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestHTTPSavePlayerShouldSavePlayerWhenPost(t *testing.T) {
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

	b, err := json.Marshal(p)
	assert.Nil(t, err)

	req := httptest.NewRequest("POST", "/v1/players", bytes.NewReader(b))
	resp := httptest.NewRecorder()
	NewHTTPTransport(es, log.NewNopLogger()).ServeHTTP(resp, req)
	body, err := ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)
	assert.True(t, strings.Contains(string(body), "Bortles"))
	assert.Equal(t, http.StatusCreated, resp.Code)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestHTTPSavePlayerShouldReturnErrorWhenPostAndDatabaseError(t *testing.T) {
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

	b, err := json.Marshal(p)
	assert.Nil(t, err)

	req := httptest.NewRequest("POST", "/v1/players", bytes.NewReader(b))
	resp := httptest.NewRecorder()
	NewHTTPTransport(es, log.NewNopLogger()).ServeHTTP(resp, req)
	body, err := ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)
	assert.True(t, strings.Contains(string(body), `"error":"database error"`))
	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestHTTPSavePlayerShouldReturnErrorWhenPostContainsID(t *testing.T) {
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

	es := NewEndpoints(NewService(db))

	b, err := json.Marshal(p)
	assert.Nil(t, err)

	req := httptest.NewRequest("POST", "/v1/players", bytes.NewReader(b))
	resp := httptest.NewRecorder()
	NewHTTPTransport(es, log.NewNopLogger()).ServeHTTP(resp, req)
	body, err := ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)
	assert.True(t, strings.Contains(string(body), `"error":"bad request"`))
	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestHTTPSavePlayerShouldReturnErrorWhenPostBodyIsMalformatted(t *testing.T) {
	db, mock, err := createDB()
	assert.Nil(t, err)
	defer db.Close()

	es := NewEndpoints(NewService(db))

	req := httptest.NewRequest("POST", "/v1/players", strings.NewReader("bad format"))
	resp := httptest.NewRecorder()
	NewHTTPTransport(es, log.NewNopLogger()).ServeHTTP(resp, req)
	body, err := ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)
	assert.True(t, strings.Contains(string(body), `"error":"bad request"`))
	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestHTTPSavePlayerShouldSavePlayerWhenPut(t *testing.T) {
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

	b, err := json.Marshal(p)
	assert.Nil(t, err)

	req := httptest.NewRequest("PUT", "/v1/players/1", bytes.NewReader(b))
	resp := httptest.NewRecorder()
	NewHTTPTransport(es, log.NewNopLogger()).ServeHTTP(resp, req)
	body, err := ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)
	assert.True(t, strings.Contains(string(body), "Bortles"))
	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestHTTPSavePlayerShouldReturnErrorWhenPutAndDatabaseError(t *testing.T) {
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

	b, err := json.Marshal(p)
	assert.Nil(t, err)

	req := httptest.NewRequest("PUT", "/v1/players/1", bytes.NewReader(b))
	resp := httptest.NewRecorder()
	NewHTTPTransport(es, log.NewNopLogger()).ServeHTTP(resp, req)
	body, err := ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)
	assert.True(t, strings.Contains(string(body), `"error":"database error"`))
	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestHTTPSavePlayerShouldReturnErrorWhenPutDoesNotContainID(t *testing.T) {
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

	es := NewEndpoints(NewService(db))

	b, err := json.Marshal(p)
	assert.Nil(t, err)

	req := httptest.NewRequest("PUT", "/v1/players/1", bytes.NewReader(b))
	resp := httptest.NewRecorder()
	NewHTTPTransport(es, log.NewNopLogger()).ServeHTTP(resp, req)
	body, err := ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)
	assert.True(t, strings.Contains(string(body), `"error":"bad request"`))
	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestHTTPSavePlayerShouldReturnErrorWhenPutIDsDoNotMatch(t *testing.T) {
	db, mock, err := createDB()
	assert.Nil(t, err)
	defer db.Close()

	p := &models.Player{
		ID:         2,
		Name:       "Blake Bortles",
		Number:     "5",
		Position:   "QB",
		Height:     "6-5",
		Weight:     "236",
		Age:        "26",
		Experience: 5,
		College:    "Central Florida",
	}

	es := NewEndpoints(NewService(db))

	b, err := json.Marshal(p)
	assert.Nil(t, err)

	req := httptest.NewRequest("PUT", "/v1/players/1", bytes.NewReader(b))
	resp := httptest.NewRecorder()
	NewHTTPTransport(es, log.NewNopLogger()).ServeHTTP(resp, req)
	body, err := ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)
	assert.True(t, strings.Contains(string(body), `"error":"bad request"`))
	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestHTTPSavePlayerShouldReturnErrorWhenPutBodyIsMalformatted(t *testing.T) {
	db, mock, err := createDB()
	assert.Nil(t, err)
	defer db.Close()

	es := NewEndpoints(NewService(db))

	req := httptest.NewRequest("PUT", "/v1/players/1", strings.NewReader("bad format"))
	resp := httptest.NewRecorder()
	NewHTTPTransport(es, log.NewNopLogger()).ServeHTTP(resp, req)
	body, err := ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)
	assert.True(t, strings.Contains(string(body), `"error":"bad request"`))
	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestHTTPSavePlayerShouldReturnErrorWhenPutIDIsNonNumeric(t *testing.T) {
	db, mock, err := createDB()
	assert.Nil(t, err)
	defer db.Close()

	es := NewEndpoints(NewService(db))

	req := httptest.NewRequest("PUT", "/v1/players/a", nil)
	resp := httptest.NewRecorder()
	NewHTTPTransport(es, log.NewNopLogger()).ServeHTTP(resp, req)
	body, err := ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)
	assert.True(t, strings.Contains(string(body), `"error":"bad request"`))
	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestHTTPSavePlayerShouldReturnErrorWhenPutIDIsMissing(t *testing.T) {
	db, mock, err := createDB()
	assert.Nil(t, err)
	defer db.Close()

	es := NewEndpoints(NewService(db))

	req := httptest.NewRequest("PUT", "/v1/players", nil)
	resp := httptest.NewRecorder()
	NewHTTPTransport(es, log.NewNopLogger()).ServeHTTP(resp, req)
	body, err := ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)
	assert.Equal(t, "", string(body))
	assert.Equal(t, http.StatusMethodNotAllowed, resp.Code)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestHTTPDeletePlayerShouldDeletePlayerWhenExists(t *testing.T) {
	db, mock, err := createDB()
	assert.Nil(t, err)
	defer db.Close()

	mock.ExpectExec(`^DELETE FROM players
		WHERE id=\$1$`).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	es := NewEndpoints(NewService(db))

	req := httptest.NewRequest("DELETE", "/v1/players/1", nil)
	resp := httptest.NewRecorder()
	NewHTTPTransport(es, log.NewNopLogger()).ServeHTTP(resp, req)
	body, err := ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)
	assert.Equal(t, "", string(body))
	assert.Equal(t, http.StatusNoContent, resp.Code)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestHTTPDeletePlayerShouldReturnErrorWhenNotExists(t *testing.T) {
	db, mock, err := createDB()
	assert.Nil(t, err)
	defer db.Close()

	mock.ExpectExec(`^DELETE FROM players
		WHERE id=\$1$`).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 0))

	es := NewEndpoints(NewService(db))

	req := httptest.NewRequest("DELETE", "/v1/players/1", nil)
	resp := httptest.NewRecorder()
	NewHTTPTransport(es, log.NewNopLogger()).ServeHTTP(resp, req)
	body, err := ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)
	assert.True(t, strings.Contains(string(body), `"error":"not found"`))
	assert.Equal(t, http.StatusNotFound, resp.Code)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestHTTPDeletePlayerShouldReturnErrorWhenDatabaseReturnsError(t *testing.T) {
	db, mock, err := createDB()
	assert.Nil(t, err)
	defer db.Close()

	mock.ExpectExec(`^DELETE FROM players
		WHERE id=\$1$`).
		WithArgs(1).
		WillReturnError(errors.New("database error"))

	es := NewEndpoints(NewService(db))

	req := httptest.NewRequest("DELETE", "/v1/players/1", nil)
	resp := httptest.NewRecorder()
	NewHTTPTransport(es, log.NewNopLogger()).ServeHTTP(resp, req)
	body, err := ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)
	assert.True(t, strings.Contains(string(body), `"error":"database error"`))
	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestHTTPDeletePlayerShouldReturnErrorWhenIDIsNonNumeric(t *testing.T) {
	db, mock, err := createDB()
	assert.Nil(t, err)
	defer db.Close()

	es := NewEndpoints(NewService(db))

	req := httptest.NewRequest("DELETE", "/v1/players/a", nil)
	resp := httptest.NewRecorder()
	NewHTTPTransport(es, log.NewNopLogger()).ServeHTTP(resp, req)
	body, err := ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)
	assert.True(t, strings.Contains(string(body), `"error":"bad request"`))
	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestHTTPDeletePlayerShouldReturnErrorWhenIDIsMissing(t *testing.T) {
	db, mock, err := createDB()
	assert.Nil(t, err)
	defer db.Close()

	es := NewEndpoints(NewService(db))

	req := httptest.NewRequest("DELETE", "/v1/players", nil)
	resp := httptest.NewRecorder()
	NewHTTPTransport(es, log.NewNopLogger()).ServeHTTP(resp, req)
	body, err := ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)
	assert.Equal(t, "", string(body))
	assert.Equal(t, http.StatusMethodNotAllowed, resp.Code)
	assert.Nil(t, mock.ExpectationsWereMet())
}
