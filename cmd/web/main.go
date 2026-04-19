package main

import (
	"log"
	"net/http"

	"github.com/edmarqueslima/memorygame/internal/application"
	"github.com/edmarqueslima/memorygame/internal/infrastructure"
	apihttp "github.com/edmarqueslima/memorygame/internal/interfaces/http"
)

func main() {
	// 1. Instanciar Infraestrutura (Repositório em Memória)
	repo := infrastructure.NewMemoryRepo()

	// 2. Instanciar Aplicação (Serviço do Jogo)
	service := application.NewGameService(repo)

	// 3. Instanciar Interface (Handlers e Router HTTP)
	handler := apihttp.NewGameHandler(service)
	router := apihttp.NewRouter(handler)

	log.Println("Servidor rodando em http://localhost:8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}
