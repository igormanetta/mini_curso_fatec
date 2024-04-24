package services

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"exemplo.com/exemplo/database"
	_ "github.com/lib/pq" // Driver PostgreSQL
)

type Services struct {
	dao *database.DAO
}

func New(d *database.DAO) *Services {
	return &Services{
		dao: d,
	}
}

// Start inicia os serviços da aplicação.
func Start() {
	fmt.Println("Iniciando os serviços da aplicação...")
}

// HandleDownload DownloadFile
func (s *Services) HandleDownload(w http.ResponseWriter, r *http.Request) {
	// Obter o valor do parâmetro "url" da consulta
	url := r.URL.Query().Get("url")

	if url == "" {
		http.Error(w, "Parâmetro 'url' não encontrado na consulta", http.StatusBadRequest)
		return
	}

	// Diretório onde você deseja salvar o arquivo baixado
	DownLoadDir := `./basedir/downloads/`

	// Baixar o arquivo a partir da URL fornecida
	filename, err := s.downloadFile(url, DownLoadDir)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao baixar o arquivo: %v /URL: %s /Dir: %s", err, url, DownLoadDir), http.StatusInternalServerError)
		return
	}

	// Alterando o path para que seja armazenado no banco de dados de forma que possa ser utilizado pela nossa api
	postpath := fmt.Sprintf(`http://localhost:3000/listar/downloads/%s`, filename)

	// Adicionar a URL do arquivo à tabela 'fileserver'
	err = s.dao.ImputBD(postpath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao adicionar a URL do arquivo no banco de dados: %v", err), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Arquivo baixado e salvo com sucesso em: %s", DownLoadDir)
}

// DownloadFile irá baixar o arquivo na minha pasta temp e em seguida chamar a rotina para gravar na tabela do banco o link para consulta posterior
func (s *Services) downloadFile(url string, destinationDir string) (string, error) {
	// Obter o nome do arquivo a partir da URL
	filename := filepath.Base(url)

	// Caminho completo para salvar o arquivo
	filepath := filepath.Join(destinationDir, filename)

	// Criar o arquivo no caminho especificado e retorna um ponteiro para que o arquivo
	out, err := os.Create(filepath)
	if err != nil {
		return filepath, fmt.Errorf("erro ao criar o arquivo local: %v", err)
	}
	defer out.Close()

	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // Ignorar a verificação do certificado
		},
	}

	// Fazer a requisição HTTP GET para baixar o arquivo, aqui temos de fato o nosso download
	response, err := httpClient.Get(url)
	if err != nil {
		return filepath, fmt.Errorf("erro ao fazer a requisição GET: %v", err)
	}
	defer response.Body.Close()

	// Verificar se o código de resposta é OK (200)
	if response.StatusCode != http.StatusOK {
		return filepath, fmt.Errorf("erro: resposta do servidor não está OK: %s", response.Status)
	}

	// Copia o conteúdo da resposta para o arquivo local
	_, err = io.Copy(out, response.Body)
	if err != nil {
		return filepath, fmt.Errorf("erro ao salvar o conteúdo do arquivo: %v", err)
	}

	return filename, nil
}
