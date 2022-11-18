package domain

import (
	"fmt"
	"sync"

	"github.com/DaniilShd/BlackWallGroup/dto"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type MapRequest struct {
	MapRequestId map[string][]dto.ClientRequest
	mu           sync.Mutex
	DB           *sqlx.DB
}

func NewMapRequest(dbClient *sqlx.DB) ClientMapCache {
	return &MapRequest{
		MapRequestId: make(map[string][]dto.ClientRequest),
		DB:           dbClient,
	}
}

func (m *MapRequest) AddToMapRequest(r *dto.ClientRequest) {
	m.mu.Lock()
	defer m.mu.Unlock()
	r.UUID = uuid.New()
	pos := *r
	m.MapRequestId[r.ClientId] = append(m.MapRequestId[r.ClientId], pos)
	_, err := m.DB.Exec(`INSERT INTO cache (uuid, client_id, type_transaction, amount) values ($1, $2, $3, $4)`, r.UUID, r.ClientId, r.Type, r.Amount)
	if err != nil {
		fmt.Println(err)
	}
}

func (m *MapRequest) DeleteToMapRequest(r *dto.ClientRequest) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for i, val := range m.MapRequestId[r.ClientId] {
		if r.UUID == val.UUID {
			m.MapRequestId[r.ClientId][i] = m.MapRequestId[r.ClientId][len(m.MapRequestId[r.ClientId])-1]
			m.MapRequestId[r.ClientId] = m.MapRequestId[r.ClientId][:len(m.MapRequestId[r.ClientId])-1]
		}
	}
	fmt.Println(r.UUID)
	_, err := m.DB.Exec(`DELETE FROM cache WHERE uuid = $1`, r.UUID)
	if err != nil {
		fmt.Println(err)
	}
}
