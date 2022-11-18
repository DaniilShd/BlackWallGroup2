package app

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/DaniilShd/BlackWallGroup/domain"
	"github.com/DaniilShd/BlackWallGroup/dto"
	"github.com/gorilla/mux"
)

type ClientHandler struct {
	domain domain.ClientRepository
}

func (c ClientMapRequestAndHandler) makeTransaction(w http.ResponseWriter, r *http.Request) {

	//get id from request
	vars := mux.Vars(r)
	clientId := vars["id_client"]

	var request dto.ClientRequest
	//decode body to json
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		writeResponse(w, http.StatusBadRequest, err.Error())
	} else {
		//request validation (function of validation is in the folder dto)
		if ok, err := dto.Valid(request); !ok {
			writeResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		request.ClientId = clientId

		// cache to map
		c.MapRequest.domain.AddToMapRequest(&request)
		fmt.Println(request)

		//save Transaction
		clientResp, err := c.Handler.domain.SaveTransaction(request)

		//delete request from caheMap
		c.MapRequest.domain.DeleteFromMapRequest(&request)

		if err != nil {
			writeResponse(w, http.StatusBadRequest, err.Error())
		} else {
			writeResponse(w, http.StatusCreated, clientResp)
		}
	}
}

func writeResponse(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		panic(err)
	}
}
