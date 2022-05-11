package main

import (
	"context"
	"fmt"
	"os"

	pgx "github.com/jackc/pgx/v4"
)

func connectDB() (*pgx.Conn, error) {
	connstr := os.Getenv("DATATESTE_URL")
	if len(connstr) == 0 {
		err := fmt.Errorf("sem url de conexão")
		return nil, err
	}
	conn, err := pgx.Connect(context.Background(), connstr)
	if err != nil {
		fmt.Errorf("impossivel estaelecer conexão[%v]: %v", connstr, err)
	}
	return conn, nil
}

var conn *pgx.Conn

func init() {
	var err error
	conn, err = connectDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, "acesso negado, contate suporte %v", err)
		os.Exit(1)

	}
}

type user struct {
	cod       string
	nome      string
	sobrenome string
	telefone  string
	cidade    string
}

func searchForName(searchName string) (*user, error) {

	var sql string = "SELECT cod, nome, sobrenome, telefone, cidade FROM contato WHERE nome = $1"
	var user user
	err := conn.QueryRow(context.Background(), sql, searchName).Scan(&user.cod, &user.nome, &user.sobrenome, &user.telefone, &user.cidade)
	if err != nil {
		fmt.Errorf("erro na consulta %v", err)
		return nil, err

	}
	return &user, nil
}

func main() {

	defer conn.Close(context.Background())

	user, err := searchForName("fulano")
	if err != nil {
		fmt.Fprintf(os.Stderr, "impossivel acessar o banco de dados %v", err)
		os.Exit(2)

	}
	fmt.Printf("dados gerados pelo banco de dados: \n")
	fmt.Printf("cod\t%v\n", user.cod)
	fmt.Printf("nome\t%v\n", user.nome)
	fmt.Printf("sobrenome\t%v\n", user.sobrenome)
	fmt.Printf("telenone\t%v\n", user.telefone)
	fmt.Printf("cidade\t%v\n", user.cidade)
}
