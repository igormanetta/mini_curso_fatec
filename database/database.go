package database

import (
	"fmt"
	"net/http"

	"exemplo.com/exemplo/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Driver PostgreSQL
)

// FileServerRow representa uma linha da tabela fileserver e será utilizado para exibir no navegador as informações ao usuário
type FileServerRow struct {
	ID  int    `db:"id_lcto"`
	URL string `db:"url"`
}

// Dao é uma estrutura para gerenciar a conexão com o banco de dados PostgreSQL.
type DAO struct {
	db *sqlx.DB
}

func New(cfg *config.AppConfig) (*DAO, error) {
	db, err := NewDatabase(cfg)
	if err != nil {
		return nil, fmt.Errorf("erro ao conectar com o banco: %w", err)
	}

	return &DAO{db: db}, nil
}

// NewDatabase cria uma nova instância do Database
func NewDatabase(cfg *config.AppConfig) (*sqlx.DB, error) {

	fmt.Println("Conectando ao banco postgres...")
	connectionString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database)

	db, err := sqlx.Open("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("erro ao conectar ao banco de dados: %v", err)
	}
	fmt.Println("Conectado OK")

	return db, nil
}

// Close fecha a conexão com o banco de dados.
func (d *DAO) Close() error {
	if d.db != nil {
		return d.db.Close()
	}
	return nil
}

// GetConection fornece a conexão com o banco de dados.
func (d *DAO) GetConection() *sqlx.DB {
	return d.db
}

// ConsultaFileServer consulta a tabela 'fileserver' e imprime os campos 'id_lcto' e 'url' no console.
func (d *DAO) ConsultaFileServer() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Consultando...")
		// Consulta SQL para selecionar todos os registros da tabela 'fileserver'
		query := `SELECT id_lcto, url FROM fileserver`

		// Executa a consulta no banco de dados Postgres associando para rows, um array da nossa struct FileServerRow
		var rows []FileServerRow
		err := d.db.Select(&rows, query)
		if err != nil {
			http.Error(w, fmt.Sprintf("erro ao consultar tabela 'fileserver': %v", err), http.StatusInternalServerError)
			return
		}

		// Gera uma lista usando HTML com os resultados vindos do banco de dados
		var htmlList string
		htmlList += "<ul>"
		for _, row := range rows {
			htmlList += fmt.Sprintf("<li><a href='%s'>%s</a></li>", row.URL, row.URL)
		}
		htmlList += "</ul>"

		// Escreve a lista em HTML na resposta HTTP, e retorna o StatusOK muito utilizado em praticamente toda troca de informações com APIs
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "<h1>Urls inseridas na tabela 'FileServer'</h1>")
		fmt.Fprintf(w, htmlList)

		// Resposta de sucesso = 200
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Consulta realizada com sucesso. Clique no link para acessar.\n"))
	}
}

// ImputBD Função para adicionar o caminho do arquivo baixado à tabela 'fileserver'
func (d *DAO) InputBD(filepath string) error {

	// Consulta SQL para inserir o caminho do arquivo na tabela 'fileserver'
	query := `INSERT INTO fileserver (url) VALUES ($1)`

	// Executar a consulta SQL
	_, err := d.db.Exec(query, filepath)
	if err != nil {
		return fmt.Errorf("erro ao inserir o caminho do arquivo na tabela 'fileserver': %v", err)
	}

	return nil
}
