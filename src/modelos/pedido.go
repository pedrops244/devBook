package modelos

import (
	"database/sql"
	"errors"
	"time"
)

type Pedido struct {
	ID          uint          `json:"id,omitempty"`
	Status      string        `json:"status,omitempty"`
	CriadoEm    time.Time     `json:"criadoEm,omitempty"`
	RecebidoEm  sql.NullTime  `json:"recebidoEm,omitempty"`
	ConferidoEm sql.NullTime  `json:"conferidoEm,omitempty"`
	Itens       []ItensPedido `json:"itens,omitempty"` // Relacionamento com itens
}

func (pedido *Pedido) Validar() error {
	if pedido.Status == "" {
		return errors.New("o status do pedido é obrigatório")
	}
	if len(pedido.Itens) == 0 {
		return errors.New("o pedido deve conter pelo menos um item")
	}
	return nil
}
