package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/edmarqueslima/memorygame/internal/application"
	"github.com/edmarqueslima/memorygame/internal/infrastructure"
	apihttp "github.com/edmarqueslima/memorygame/internal/interfaces/http"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// 1. Instanciar Infraestrutura (Repositório em Memória)
	repo := infrastructure.NewMemoryRepo()

	// 2. Instanciar Aplicação (Serviço do Jogo)
	service := application.NewGameService(repo)

	// 3. Instanciar Interface (Handlers e Router HTTP)
	handler := apihttp.NewGameHandler(service)
	router := apihttp.NewRouter(handler)

	log.Printf("Servidor rodando em http://localhost:%s", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), router); err != nil {
		log.Fatal(err)
	}
}
