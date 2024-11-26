package repositorios

import (
	"api/src/modelos"
	"database/sql"
	"errors"
)

// ItensPedidosRepositorio representa o repositório de itens do pedido
type ItensPedidosRepositorio struct {
	db *sql.DB
}

// NovoRepositorioDeItensPedidos cria um repositório de itens do pedido
func NovoRepositorioDeItensPedidos(db *sql.DB) *ItensPedidosRepositorio {
	return &ItensPedidosRepositorio{db}
}

// Adicionar adiciona um item ao pedido
func (repositorio *ItensPedidosRepositorio) Adicionar(item modelos.ItensPedido) (int64, error) {
	query := `
		INSERT INTO ItensPedidos 
		(PedidoID, QuantidadeSolicitada, QuantidadeRecebida, QuantidadeConferida, Codigo)
		OUTPUT INSERTED.ID
		VALUES (@PedidoID, @QuantidadeSolicitada, 0, 0, @Codigo)
	`
	stmt, err := repositorio.db.Prepare(query)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	var id int64
	err = stmt.QueryRow(
		sql.Named("PedidoID", item.PedidoID),
		sql.Named("QuantidadeSolicitada", item.QuantidadeSolicitada),
		sql.Named("Codigo", item.Codigo),
	).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// BuscarPorPedido retorna todos os itens associados a um pedido
func (repositorio *ItensPedidosRepositorio) BuscarPorPedido(pedidoID uint) ([]modelos.ItensPedido, error) {
	query := `
		SELECT ID, PedidoID, QuantidadeSolicitada, QuantidadeRecebida, QuantidadeConferida, Codigo
		FROM ItensPedidos
		WHERE PedidoID = @PedidoID
	`
	rows, err := repositorio.db.Query(query, sql.Named("PedidoID", pedidoID))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var itens []modelos.ItensPedido
	for rows.Next() {
		var item modelos.ItensPedido
		if err := rows.Scan(
			&item.ID,
			&item.PedidoID,
			&item.QuantidadeSolicitada,
			&item.QuantidadeRecebida,
			&item.QuantidadeConferida,
			&item.Codigo,
		); err != nil {
			return nil, err
		}
		itens = append(itens, item)
	}
	return itens, nil
}

// DeletarPorPedido remove todos os itens de um pedido
func (repositorio *ItensPedidosRepositorio) DeletarPorPedido(pedidoID uint) error {
	query := `
		DELETE FROM ItensPedidos
		WHERE PedidoID = @PedidoID
	`
	stmt, err := repositorio.db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(sql.Named("PedidoID", pedidoID))
	return err
}

// BuscarPorID busca um item específico pelo ID
func (repositorio *ItensPedidosRepositorio) BuscarPorID(itemID uint) (modelos.ItensPedido, error) {
	query := `
		SELECT ID, PedidoID, QuantidadeSolicitada, QuantidadeRecebida, QuantidadeConferida, Codigo
		FROM ItensPedidos
		WHERE ID = @ItemID
	`
	var item modelos.ItensPedido
	err := repositorio.db.QueryRow(query, sql.Named("ItemID", itemID)).
		Scan(
			&item.ID,
			&item.PedidoID,
			&item.QuantidadeSolicitada,
			&item.QuantidadeRecebida,
			&item.QuantidadeConferida,
			&item.Codigo,
		)
	if err == sql.ErrNoRows {
		return item, errors.New("item não encontrado")
	} else if err != nil {
		return item, err
	}
	return item, nil
}
