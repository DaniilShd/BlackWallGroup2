package domain

import (
	"github.com/DaniilShd/BlackWallGroup/dto"
)

type ClientMapCache interface {
	AddToMapRequest(r *dto.ClientRequest)
	DeleteFromMapRequest(r *dto.ClientRequest)
}

type ClientRepository interface {
	SaveTransaction(t dto.ClientRequest) (*dto.ClientResponse, error)
	GetHistory(id string) (*dto.ClientHistoryTransaction, error)
}
