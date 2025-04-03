package userHandler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/kourai55k/booking-service/internal/domain"
)

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	const op = "http.userHandler.DeleteUser"

	log := h.logger

	log.Debug("request received", "method", r.Method, "path", r.URL.Path)

	// Extract the 'id' path parameter using Go 1.22's PathValue
	idStr := r.PathValue("id")

	// Convert 'id' to uint
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil || idStr == "" {
		http.Error(w, "bad request", http.StatusBadRequest)
		log.Error("bad request", "err", fmt.Errorf("%s: bad request", op).Error())
		return
	}

	err = h.userService.DeleteUser(uint(id))
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			http.Error(w, "user not found", http.StatusNotFound)
			log.Error("user not found", "err", fmt.Errorf("%s: %w", op, err).Error())
			return
		}
		http.Error(w, "failed to delete user", http.StatusInternalServerError)
		log.Error("failed to delete user", "err", err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}
