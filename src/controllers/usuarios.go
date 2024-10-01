package controllers

import (
	"api/src/banco"
	"api/src/modelos"
	"api/src/repositorios"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func CriarUsuario(w http.ResponseWriter, r *http.Request) {
	corpoRequest, erro := io.ReadAll(r.Body)
	if erro != nil {
		log.Fatal(erro)
	}

	var usuario modelos.Usuario
	if erro = json.Unmarshal(corpoRequest, &usuario); erro != nil {
		log.Fatal(erro)
	}
	db, erro := banco.Conectar()
	if erro != nil {
		log.Fatal(erro)
	}

	repositorio := repositorios.NovoRepositorioDeUsuario(db)
	usuarioID, erro := repositorio.Criar(usuario)
	if erro != nil {
		log.Fatal(erro)
	}

	w.Write([]byte(fmt.Sprintf("Id inserido: %d", usuarioID)))
}
func BuscarUsuarios(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Buscando todos os usuarios"))

}
func BuscarUsuario(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Criando usuário"))

}
func AtualizarUsuario(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Atualizando usuário"))

}
func DeletarUsuario(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Deletando usuário"))

}
