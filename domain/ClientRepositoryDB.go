package domain

import (
	"errors"
	"fmt"

	"reflect"

	"github.com/DaniilShd/BlackWallGroup/dto"
	"github.com/jmoiron/sqlx"
)

type ClientDomainDb struct {
	client *sqlx.DB
}

func NewClientRepositoryDb(dbClient *sqlx.DB) ClientRepository {
	return &ClientDomainDb{dbClient}
}

func (c *ClientDomainDb) SaveTransaction(t dto.ClientRequest) (*dto.ClientResponse, error) {
	// start transaction
	tx, err := c.client.Begin()
	if err != nil {
		fmt.Println("Error while starting a new transaction for bank account transaction: " + err.Error())
		return nil, errors.New("Unexpected database error")
	}

	// inserting bank account transaction
	_, err = tx.Exec(`INSERT INTO transactions (client_id, type_transaction, amount) values ($1, $2, $3)`, t.ClientId, t.Type, t.Amount)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("Unexpected database error")
	}

	//Select client info
	sql_statement := "SELECT * from clients where id = $1;"
	row := c.client.QueryRow(sql_statement, t.ClientId)

	checkSum := dto.ClientResponse{}
	err = row.Scan(&checkSum.ClientId, &checkSum.Sum)
	if err != nil {
		fmt.Println(err)
	}

	// updating account balance and Checking a non-negative account balance
	if t.IsWithdrawal() {
		if !checkSum.CheckSum(t.Amount) {
			tx.Rollback()
			return nil, errors.New("Insufficient funds on the account")
		}
		_, err = tx.Exec(`UPDATE clients SET bank_account = bank_account - $1 where id = $2`, t.Amount, t.ClientId)
	} else {
		_, err = tx.Exec(`UPDATE clients SET bank_account = bank_account + $1 where id = $2`, t.Amount, t.ClientId)
	}

	// in case of error Rollback, and changes from both the tables will be reverted
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// commit the transaction when all is good
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// Getting info about client
	sql_st := "SELECT * from clients where id = $1;"
	rowRes := c.client.QueryRow(sql_st, t.ClientId)

	Response := dto.ClientResponse{}
	err = rowRes.Scan(&Response.ClientId, &Response.Sum)
	if err != nil {
		fmt.Println(err)
	}

	return &Response, nil

}

func (c *ClientDomainDb) GetHistory(id string) (*dto.ClientHistoryTransaction, error) {

	//check id client
	sqlStatClient := "SELECT * FROM clients WHERE id = $1;"
	rowClient := c.client.QueryRow(sqlStatClient, id)

	client := dto.ClientResponse{}
	err := rowClient.Scan(&client.ClientId, &client.Sum)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(client.ClientId)
	fmt.Println(reflect.TypeOf(client.ClientId))

	//Checking if the client exists
	if client.ClientId == 0 {
		return nil, errors.New("Client does not exist")
	}

	sqlStatTransaction := "SELECT id_transaction, type_transaction, amount FROM transactions WHERE client_id = $1;"
	rows, err := c.client.Query(sqlStatTransaction, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	transactions := dto.ClientHistoryTransaction{}
	for rows.Next() {
		var transaction dto.ClientTransaction
		err = rows.Scan(&transaction.TransactionId, &transaction.Type, &transaction.Amount)
		if err != nil {
			fmt.Println(err)
		}
		transactions.History = append(transactions.History, transaction)
	}

	return &transactions, nil
}
