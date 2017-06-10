package middleware

import (
	log "github.com/Sirupsen/logrus"
	"github.com/mstovicek/go-chapter-route-recommendation/api/response"
	"net/http"
)

type recoveryHandler struct {
	logger  *log.Logger
	handler http.Handler
}

func NewRecovery(l *log.Logger, h http.Handler) http.Handler {
	return &recoveryHandler{
		logger:  l,
		handler: h,
	}
}

func (h *recoveryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		// recover allows you to continue execution in case of panic
		if err := recover(); err != nil {
			response.WriteError(w, http.StatusInternalServerError, "Internal Server Error")

			h.logger.WithFields(log.Fields{
				"method": r.Method,
				"path":   r.URL.Path,
				"error":  err,
			}).Error("Internal Server Error")
		}
	}()

	h.handler.ServeHTTP(w, r)
}
