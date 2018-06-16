package players

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/hoop33/roster/models"
)

var errBadRoute = errors.New("bad route")
var errBadRequest = errors.New("bad request")

// NewHTTPTransport returns a handler for HTTP transport
func NewHTTPTransport(ep *Endpoints, logger log.Logger) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(log.With(logger, "tag", "http")),
		kithttp.ServerErrorEncoder(encodeHTTPError),
	}

	listPlayersHandler := kithttp.NewServer(
		ep.listPlayersEndpoint,
		decodeHTTPListPlayersRequest,
		encodeHTTPListPlayersResponse,
		opts...,
	)

	getPlayerHandler := kithttp.NewServer(
		ep.getPlayerEndpoint,
		decodeHTTPGetPlayerRequest,
		encodeHTTPGetPlayerResponse,
		opts...,
	)

	createPlayerHandler := kithttp.NewServer(
		ep.savePlayerEndpoint,
		decodeHTTPCreatePlayerRequest,
		encodeHTTPSavePlayerResponse,
		opts...,
	)

	updatePlayerHandler := kithttp.NewServer(
		ep.savePlayerEndpoint,
		decodeHTTPUpdatePlayerRequest,
		encodeHTTPSavePlayerResponse,
		opts...,
	)

	deletePlayerHandler := kithttp.NewServer(
		ep.deletePlayerEndpoint,
		decodeHTTPDeletePlayerRequest,
		encodeHTTPDeletePlayerResponse,
		opts...,
	)

	r := mux.NewRouter()
	r.Handle("/v1/players", listPlayersHandler).Methods("GET")
	r.Handle("/v1/players/{id}", getPlayerHandler).Methods("GET")
	r.Handle("/v1/players", createPlayerHandler).Methods("POST")
	r.Handle("/v1/players/{id}", updatePlayerHandler).Methods("PUT")
	r.Handle("/v1/players/{id}", deletePlayerHandler).Methods("DELETE")

	return accessControl(r)
}

func accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")

		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)
	})
}

func decodeHTTPListPlayersRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return listPlayersRequest{
		Position: r.URL.Query().Get("position"),
	}, nil
}

func encodeHTTPListPlayersResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	lpr := response.(listPlayersResponse)
	if lpr.Err == "" {
		return encodeHTTPResponse(ctx, http.StatusOK, w, response)
	}
	encodeHTTPError(ctx, getHTTPError(lpr.Err), w)
	return nil
}

func decodeHTTPGetPlayerRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, errBadRoute
	}

	ID, err := strconv.Atoi(id)
	if err != nil {
		return nil, errBadRequest
	}

	return getPlayerRequest{
		ID: ID,
	}, nil
}

func encodeHTTPGetPlayerResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	gpr := response.(getPlayerResponse)
	if gpr.Err == "" {
		return encodeHTTPResponse(ctx, http.StatusOK, w, response)
	}
	encodeHTTPError(ctx, getHTTPError(gpr.Err), w)
	return nil
}

func decodeHTTPCreatePlayerRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var player models.Player
	if err := json.NewDecoder(r.Body).Decode(&player); err != nil {
		return nil, errBadRequest
	}

	if player.ID > 0 {
		return nil, errBadRequest
	}

	return savePlayerRequest{
		Player: &player,
	}, nil
}

func decodeHTTPUpdatePlayerRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, errBadRoute
	}

	ID, err := strconv.Atoi(id)
	if err != nil {
		return nil, errBadRequest
	}

	var player models.Player
	if err := json.NewDecoder(r.Body).Decode(&player); err != nil {
		return nil, errBadRequest
	}

	if ID != player.ID {
		return nil, errBadRequest
	}

	return savePlayerRequest{
		Player: &player,
	}, nil
}

func encodeHTTPSavePlayerResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	spr := response.(savePlayerResponse)
	if spr.Err == "" {
		sc := http.StatusOK
		if spr.Created {
			sc = http.StatusCreated
		}
		return encodeHTTPResponse(ctx, sc, w, response)
	}
	encodeHTTPError(ctx, getHTTPError(spr.Err), w)
	return nil
}

func decodeHTTPDeletePlayerRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, errBadRoute
	}

	ID, err := strconv.Atoi(id)
	if err != nil {
		return nil, errBadRequest
	}

	return deletePlayerRequest{
		ID: ID,
	}, nil
}

func encodeHTTPDeletePlayerResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	dpr := response.(deletePlayerResponse)
	if dpr.Err == "" {
		return encodeHTTPResponse(ctx, http.StatusNoContent, w, nil)
	}
	encodeHTTPError(ctx, getHTTPError(dpr.Err), w)
	return nil
}

func encodeHTTPResponse(_ context.Context, statusCode int, w http.ResponseWriter, response interface{}) error {
	w.WriteHeader(statusCode)
	if response == nil {
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func encodeHTTPError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	switch err {
	case errBadRequest:
		w.WriteHeader(http.StatusBadRequest)
	case errNotFound:
		w.WriteHeader(http.StatusNotFound)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	e := json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
	if e != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func getHTTPError(str string) error {
	if str == "not found" {
		return errNotFound
	}
	return errors.New(str)
}
