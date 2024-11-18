package main

import (
	"flag"
	"fmt"
	"log"
)


func seedAccount(store Storage, fname, lname, password string) *Account {
	acc, err := NewAccount(fname,lname,password)
	if err != nil {
		fmt.Printf("1 %+v\n",err)
		log.Fatal(err)
	}
	if err := store.CreateAccount(acc); err != nil{
		fmt.Printf("2 %+v\n",err)
		log.Fatal(err)
	}
	fmt.Println("new Account => ", acc.Number)
	return acc
}

func seedAccounts(s Storage) {
	seedAccount(s ,"antony","GG","hunter8888"); 
}

func main() {
	seed := flag.Bool("seed", false, "seed the db")
	flag.Parse()

	store, err := NewPostgressStore()
	if err != nil {
		log.Fatal(err)
	}

	if err = store.Init(); err != nil{
		log.Fatal(err)
	}


	if *seed {
		seedAccounts(store)
		fmt.Println("seeding the database")
		
	}

	server := NewAPIServer(":3000",store)
	server.Run()
}