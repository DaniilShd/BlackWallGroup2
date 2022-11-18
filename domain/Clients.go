package domain

import (
	"github.com/DaniilShd/BlackWallGroup/dto"
)

type ClientMapCache interface {
	AddToMapRequest(r *dto.ClientRequest)
	DeleteToMapRequest(r *dto.ClientRequest)
}

type ClientRepository interface {
	SaveTransaction(t dto.ClientRequest) (*dto.ClientResponse, error)
	GetHistory(id string) (*dto.ClientHistoryTransaction, error)
}
