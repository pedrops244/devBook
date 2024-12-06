package repositorios

import (
	"api/src/modelos"
	"database/sql"
	"errors"
)

// pedidos representa o repositório de pedidos
type pedidos struct {
	db *sql.DB
}

// NovoRepositorioDePedidos cria um repositório de pedidos
func NovoRepositorioDePedidos(db *sql.DB) *pedidos {
	return &pedidos{db}
}

// Criar adiciona um novo pedido com seus itens
func (repositorio pedidos) Criar(pedido modelos.Pedido) (uint, error) {
	tx, err := repositorio.db.Begin()
	if err != nil {
		return 0, err
	}

	// Inserir o pedido
	queryPedido := `
		INSERT INTO Pedidos (Status, UsuarioId, CriadoEm)
		OUTPUT INSERTED.ID
		VALUES (@Status, @UsuarioId, @CriadoEm)
	`
	stmtPedido, err := tx.Prepare(queryPedido)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	defer stmtPedido.Close()

	var pedidoID uint
	err = stmtPedido.QueryRow(
		sql.Named("Status", pedido.Status),
		sql.Named("UsuarioId", pedido.UsuarioId),
		sql.Named("CriadoEm", pedido.CriadoEm),
	).Scan(&pedidoID)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// Inserir os itens do pedido
	queryItem := `
		INSERT INTO ItensPedidos 
		(PedidoID, QuantidadeSolicitada, QuantidadeRecebida, QuantidadeConferida, Codigo)
		VALUES (@PedidoID, @QuantidadeSolicitada, @QuantidadeRecebida, @QuantidadeConferida, @Codigo)
	`
	stmtItem, err := tx.Prepare(queryItem)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	defer stmtItem.Close()

	for _, item := range pedido.Itens {
		_, err := stmtItem.Exec(
			sql.Named("PedidoID", pedidoID),
			sql.Named("QuantidadeSolicitada", item.QuantidadeSolicitada),
			sql.Named("QuantidadeRecebida", 0),
			sql.Named("QuantidadeConferida", 0),
			sql.Named("Codigo", item.Codigo),
		)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	// Commit da transação
	if err := tx.Commit(); err != nil {
		return 0, err
	}

	// Retornar o ID do pedido criado
	return pedidoID, nil
}

// BuscarPorID busca o pedido pelo ID
func (repositorio pedidos) BuscarPorID(id uint) (modelos.Pedido, error) {
	queryPedido := `
		SELECT ID, Status, UsuarioId, CriadoEm
		FROM Pedidos
		WHERE ID = @ID
	`
	var pedido modelos.Pedido
	err := repositorio.db.QueryRow(queryPedido, sql.Named("ID", id)).
		Scan(&pedido.ID, &pedido.Status, &pedido.UsuarioId, &pedido.CriadoEm)
	if err == sql.ErrNoRows {
		return pedido, errors.New("pedido não encontrado")
	} else if err != nil {
		return pedido, err
	}

	// Buscar os itens associados ao pedido
	queryItens := `
		SELECT ID, PedidoID, QuantidadeSolicitada, QuantidadeRecebida, QuantidadeConferida, Codigo
		FROM ItensPedidos
		WHERE PedidoID = @PedidoID
	`
	rows, err := repositorio.db.Query(queryItens, sql.Named("PedidoID", pedido.ID))
	if err != nil {
		return pedido, err
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
			return pedido, err
		}
		itens = append(itens, item)
	}
	pedido.Itens = itens
	return pedido, nil
}

// Listar lista todos os pedidos
func (repositorio pedidos) Listar() ([]modelos.Pedido, error) {
	queryPedidos := `
		SELECT ID, Status, UsuarioId, CriadoEm, RecebidoEm, ConferidoEm
		FROM Pedidos
	`
	rows, err := repositorio.db.Query(queryPedidos)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pedidos []modelos.Pedido
	for rows.Next() {
		var pedido modelos.Pedido
		if err := rows.Scan(&pedido.ID, &pedido.Status, &pedido.UsuarioId, &pedido.CriadoEm, &pedido.RecebidoEm, &pedido.ConferidoEm); err != nil {
			return nil, err
		}

		// Buscar os itens do pedido
		itens, err := repositorio.BuscarItensDoPedido(pedido.ID)
		if err != nil {
			return nil, err
		}
		pedido.Itens = itens
		pedidos = append(pedidos, pedido)
	}
	return pedidos, nil
}

// BuscarItensDoPedido busca todos os itens dentro de um pedido especifico pelo ID
func (repositorio pedidos) BuscarItensDoPedido(pedidoID uint) ([]modelos.ItensPedido, error) {
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

// DeletarPedido remove um pedido pelo ID
func (repositorio pedidos) DeletarPedido(pedidoID uint) error {
	tx, err := repositorio.db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	query := `
		DELETE FROM Pedidos
		WHERE ID = @PedidoID
	`

	_, err = tx.Exec(query, sql.Named("PedidoID", pedidoID))
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

// AtualizarRecebimento atualiza a QuantidadeRecebida no banco
func (repositorio pedidos) AtualizarRecebimento(pedidoID uint, pedido modelos.Pedido) error {
	tx, err := repositorio.db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	query := `
		UPDATE Pedidos
		SET Status = @Status, RecebidoEm = @RecebidoEm
		WHERE ID = @PedidoID
	`

	_, err = tx.Exec(
		query,
		sql.Named("Status", pedido.Status),
		sql.Named("RecebidoEm", pedido.RecebidoEm),
		sql.Named("PedidoID", pedidoID),
	)
	if err != nil {
		return err
	}

	queryItem := `
	UPDATE ItensPedidos
	SET QuantidadeRecebida = @QuantidadeRecebida
	WHERE PedidoID = @PedidoID AND Codigo = @Codigo
`
	for _, itemRecebido := range pedido.Itens {
		_, err = tx.Exec(
			queryItem,
			sql.Named("QuantidadeRecebida", itemRecebido.QuantidadeRecebida),
			sql.Named("PedidoID", pedidoID),
			sql.Named("Codigo", itemRecebido.Codigo),
		)
		if err != nil {
			return err
		}
	}
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

// AtualizarConferencia atualiza a QuantidadeConferida no banco
func (repositorio pedidos) AtualizarConferencia(pedidoID uint, pedido modelos.Pedido) error {
	tx, err := repositorio.db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	query := `
		UPDATE Pedidos
		SET Status = @Status, ConferidoEm = @ConferidoEm
		WHERE ID = @PedidoID
	`
	_, err = tx.Exec(
		query,
		sql.Named("Status", pedido.Status),
		sql.Named("ConferidoEm", pedido.ConferidoEm),
		sql.Named("PedidoID", pedidoID),
	)
	if err != nil {
		return err
	}

	queryItem := `
	UPDATE ItensPedidos
	SET QuantidadeConferida = @QuantidadeConferida
	WHERE PedidoID = @PedidoID AND Codigo = @Codigo
`
	for _, itemRecebido := range pedido.Itens {
		_, err = tx.Exec(
			queryItem,
			sql.Named("QuantidadeConferida", itemRecebido.QuantidadeConferida),
			sql.Named("PedidoID", pedidoID),
			sql.Named("Codigo", itemRecebido.Codigo),
		)
		if err != nil {
			return err
		}
	}
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
