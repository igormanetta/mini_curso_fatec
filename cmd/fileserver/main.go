package main

import (
	"fmt"
	"log"
	"net/http"

	"exemplo.com/exemplo/config"
	"exemplo.com/exemplo/database"
	"exemplo.com/exemplo/services"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
)

func main() {
	// Criar uma instância de PostgreSQLConfig
	appConfig := config.LoadConfig()

	// Criar uma nova instância do Database
	dao, err := database.New(&appConfig)
	if err != nil {
		log.Fatalf("erro ao inicializar o Database: %v", err)
	}

	defer dao.Close()

	s := services.New(dao)

	// Mensagens exibidas no terminal
	fmt.Println("Iniciando o servidor...")

	// Configuração do roteador
	r := chi.NewRouter()

	// Configuração do CORS, aqui definimos o que nossa api irá aceitar metodos
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://*", "https://*"},
		AllowedMethods:   []string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"X-Requested-With", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Rota padrão da API, util para o caso de testes de funcionamento
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Bem-vindo ao servidor de arquivos!")
	})

	//Definição do diretorio base local, que em momento algum sera exposto ao nosso navegador
	basedir := appConfig.ImagesFolderPath

	// Representa um sistema de arquivos baseado em um diretório específico no caminho
	fs := http.FileServer(http.Dir(basedir))

	// Forma de manipular rotas para servir conteúdo estático de um diretório específico no sistema de arquivos
	fs = http.StripPrefix("/listar/", fs)

	//Delegando ao handlepara servir um conteudo neste caso todos retornados pela rota listar
	r.Handle("/listar/*", fs)

	// Configurar rota para listar itens da tabela 'fileserver' no banco de dados
	r.Handle("/consultar", dao.ConsultaFileServer())

	// Configurar rota para lidar com o download
	r.HandleFunc("/downloads/", s.HandleDownload)

	// Iniciando os serviços
	services.Start()

	// Printando informações do serviço no console/terminal
	port := ":3000"
	fmt.Printf("Servidor iniciado e aguardando na porta %s\n", port)
	log.Fatal(http.ListenAndServe(port, r))
}
