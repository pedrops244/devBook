package main

import (
	"api/src/config"
	"api/src/router"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	config.Carregar()

	// Gera o roteador configurado
	r := router.Gerar()

	// Cria um handler com o middleware de CORS
	handler := aplicarCORS(r)

	// Inicia o servidor
	fmt.Printf("Escutando na porta: %d\n", config.Porta)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Porta), handler))
}

// Middleware para adicionar cabeçalhos CORS
func aplicarCORS(router *mux.Router) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent) // Resposta para requisições preflight
			return
		}

		router.ServeHTTP(w, r)
	})
}
