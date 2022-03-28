package handlers

import (
	"net/http"
	"strings"

	"github.com/gopherlearning/go-advanced-devops/internal/repositories"
)

type Handler struct {
	s repositories.Repository
}

// NewHandler создаёт новый экземпляр обработчика запросов, привязанный к хранилищу
func NewHandler(s repositories.Repository) *Handler {
	return &Handler{s: s}
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed!", http.StatusMethodNotAllowed)
		return
	}
	// if r.Header.Get("Content-Type") != "text/plain" {
	// 	http.Error(w, "Only text/plain content are allowed!", http.StatusBadRequest)
	// 	return
	// }
	target := strings.Split(r.RemoteAddr, ":")[0]
	if err := h.s.Update(target, r.RequestURI); err != nil {
		switch err {
		case repositories.ErrBadMetric:
			w.WriteHeader(http.StatusNotFound)
		case repositories.ErrWrongMetricType:
			w.WriteHeader(http.StatusNotImplemented)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)

}
