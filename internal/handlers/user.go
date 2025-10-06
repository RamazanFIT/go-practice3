package handlers

import (
	"encoding/json"
	"net/http"
	"practice_2/internal/i18n"
	"practice_2/internal/logger"
	"strconv"
)

func GetUser(w http.ResponseWriter, r *http.Request) {
	lang := i18n.ParseAcceptLanguage(r.Header.Get("Accept-Language"))
	msg := i18n.Get(lang)

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		logger.Warn("GetUser: missing id parameter")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": msg.InvalidID})
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		logger.Warn("GetUser: invalid id parameter: %s", idStr)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": msg.InvalidID})
		return
	}

	logger.Debug("GetUser: fetching user with id=%d", id)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]int{"user_id": id})
	logger.Info("GetUser: successfully returned user with id=%d", id)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	lang := i18n.ParseAcceptLanguage(r.Header.Get("Accept-Language"))
	msg := i18n.Get(lang)

	var req struct {
		Name string `json:"name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("CreateUser: failed to decode request body: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": msg.InvalidName})
		return
	}

	if req.Name == "" {
		logger.Warn("CreateUser: empty name provided")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": msg.InvalidName})
		return
	}

	logger.Debug("CreateUser: creating user with name=%s", req.Name)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"created": req.Name})
	logger.Info("CreateUser: successfully created user: %s", req.Name)
}
