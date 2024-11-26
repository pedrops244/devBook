package modelos

import (
	"errors"
)

type ItensPedido struct {
	ID                   uint   `json:"id,omitempty"`
	PedidoID             uint   `json:"pedidoID,omitempty"`
	QuantidadeSolicitada int    `json:"quantidadeSolicitada,omitempty"`
	QuantidadeRecebida   int    `json:"quantidadeRecebida,omitempty"`
	QuantidadeConferida  int    `json:"quantidadeConferida,omitempty"`
	Codigo               string `json:"codigo,omitempty"`
}

func (item *ItensPedido) Preparar() error {
	if erro := item.validar(); erro != nil {
		return erro
	}
	return nil
}

func (item *ItensPedido) validar() error {
	if item.QuantidadeSolicitada <= 0 {
		return errors.New("a quantidade solicitada deve ser maior que zero")
	}
	if item.Codigo == "" {
		return errors.New("o código do item é obrigatório")
	}
	return nil
}
