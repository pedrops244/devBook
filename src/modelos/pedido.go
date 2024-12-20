package modelos

import (
	"errors"
	"time"
)

type Pedido struct {
	ID           uint          `json:"id,omitempty"`
	Status       string        `json:"status,omitempty"`
	UsuarioId    int           `json:"usuarioId,omitempty"`
	RepositorId  int           `json:"repositorId,omitempty"`
	ConferenteId int           `json:"conferenteId,omitempty"`
	CriadoEm     time.Time     `json:"criadoEm,omitempty"`
	RecebidoEm   *time.Time    `json:"recebidoEm,omitempty"`
	ConferidoEm  *time.Time    `json:"conferidoEm,omitempty"`
	Itens        []ItensPedido `json:"itens,omitempty"` // Relacionamento com itens
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
