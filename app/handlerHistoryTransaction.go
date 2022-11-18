package app

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (c ClientHandler) getHistoryTransaction(w http.ResponseWriter, r *http.Request) {

	//get id from request
	vars := mux.Vars(r)
	clientId := vars["id_client"]

	// get history transaction of client
	clientHistoryResp, err := c.domain.GetHistory(clientId)
	if err != nil {
		writeResponse(w, http.StatusBadRequest, err.Error())
	} else {
		writeResponse(w, http.StatusCreated, clientHistoryResp)
	}
}
