package httpapi

import (
	"net/http"

	"example.com/paleo-prep-station/internal/contracts"
)

func (api *API) returnSpecimen(w http.ResponseWriter, r *http.Request) {
	var request contracts.ReturnRequest
	if !decodeJSON(w, r, &request) {
		return
	}
	value, err := api.service.ReturnSpecimen(r.Context(), r.PathValue("id"), request)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, value)
}
