package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/thehaung/juggernaut/domain"
	"github.com/thehaung/juggernaut/internal/logger"
	"github.com/thehaung/juggernaut/internal/utils/httputil"
	"net/http"
)

type juggernautHandler struct {
	logger            logger.Interface
	juggernautUseCase domain.JuggernautUseCase
}

func NewJuggernautHandler(router *chi.Mux, juggernautUseCase domain.JuggernautUseCase) {
	handler := &juggernautHandler{
		logger:            logger.GetLogger(),
		juggernautUseCase: juggernautUseCase,
	}

	router.Route("/juggernaut", func(r chi.Router) {
		r.Get("/current-ip", handler.getCurrentIP)
	})
}

func (h *juggernautHandler) getCurrentIP(w http.ResponseWriter, r *http.Request) {
	ip, err := h.juggernautUseCase.GetCurrentIP(r.Context())
	if err != nil {
		httputil.WriteError(w, http.StatusBadRequest, err)
		return
	}

	httputil.WriteResponse(w, http.StatusOK, ip)
}
