package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAccount(*Account) error
	UpdateAccount(*Account) error
	GetAccounts()([]*Account,error)
	GetAccountByNumber(int) (*Account, error)
	GetAccountByID(int) (*Account, error)
	DeleteAccount(int) error
}

type PostgressStore struct {
	db *sql.DB
}

func NewPostgressStore() (*PostgressStore, error) {
	conStr := "user=postgres dbname=postgres password=gobank sslmode=disable"
	
	db,err := sql.Open("postgres",conStr)
	if err!=nil {
		return nil,err
	}
	err = db.Ping(); 
	if err != nil{
		return nil,err
	}
	
	return &PostgressStore{
		db: db,
	},nil
}

func (s *PostgressStore) Init()error{
 return s.CreateAccountTable()
}

func (s *PostgressStore) CreateAccountTable() error{
	query:= `CREATE TABLE IF NOT EXISTS account (
		id SERIAL PRIMARY KEY,
		first_name VARCHAR(100),
		last_name VARCHAR(100),
		number SERIAL,
		encrypted_password VARCHAR(100),
		balance SERIAL,
		created_at timestamp
	)`

	_,err := s.db.Exec(query)

	return err
}

// --------------------Create an account--------------------
func (s *PostgressStore) CreateAccount(acc *Account) error{
	query := `
	INSERT INTO account(first_name, last_name, number, encrypted_password, balance, created_at )
	VALUES($1, $2, $3, $4, $5, $6 )
	`
	stmt,err:= s.db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_ ,err = stmt.Exec(acc.FirstName, acc.LastName, acc.Number, acc.EncryptedPassword, acc.Balance, acc.CreatedTime)
	if err != nil {
		return err
	}

	return nil
}
func (s *PostgressStore) UpdateAccount(*Account) error{
	return nil
}


// --------------------Delete an account--------------------
func (s *PostgressStore) DeleteAccount(id int) error{

	_,err := s.GetAccountByID(id)
	if err != nil {
		return err
	}

	query := `DELETE FROM account WHERE id = $1`
	stmt, err := s.db.Prepare(query)
	if err != nil{
		return err
	}
	defer stmt.Close()

	_,err = stmt.Exec(id)

	return err
}


// --------------------Get all accounts--------------------
func (s *PostgressStore) GetAccounts()([]*Account,error){
	rows,err := s.db.Query(`SELECT * FROM account`)
	if err != nil {
		return nil, err
	}
	accounts:=[]*Account{}
	for rows.Next(){
		acc := Account{}
	err := rows.Scan(
		&acc.ID,
		&acc.FirstName,
		&acc.LastName,
		&acc.Number,
		&acc.EncryptedPassword,
		&acc.Balance,
		&acc.CreatedTime)

		if err != nil{
			return nil, err
		}
		accounts = append(accounts, &acc)
	}
	return accounts,nil
}

// --------------------Get account by Number--------------------
func(s *PostgressStore) GetAccountByNumber(number int)(*Account,error){
	var acc Account

	query := `SELECT * FROM account WHERE number = $1`

	row := s.db.QueryRow(query,number)
	if err := row.Scan(
		&acc.ID,
		&acc.FirstName,
		&acc.LastName,
		&acc.Number,
		&acc.EncryptedPassword,
		&acc.Balance,
		&acc.CreatedTime); err != nil{
			if err ==sql.ErrNoRows{
				return &acc, fmt.Errorf("accountById %d: no such account",number)
			}
		return &acc, fmt.Errorf("accountById %d : %v",number,err)
	}

	return &acc,nil

}


// --------------------Get account by ID--------------------
func (s *PostgressStore) GetAccountByID(id int) (*Account,error){
	var acc Account

	query := `SELECT * FROM account WHERE id = $1`

	row := s.db.QueryRow(query,id)
	if err := row.Scan(
		&acc.ID,
		&acc.FirstName,
		&acc.LastName,
		&acc.Number,
		&acc.EncryptedPassword,
		&acc.Balance,
		&acc.CreatedTime); err != nil{
			if err ==sql.ErrNoRows{
				return &acc, fmt.Errorf("accountById %d: no such account",id)
			}
		return &acc, fmt.Errorf("accountById %d : %v",id,err)
	}

	return &acc,nil
}
