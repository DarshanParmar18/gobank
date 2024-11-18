package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
)

type APIServer struct {
	listenAddr string
	store Storage
}

func NewAPIServer(listenAddr string, store Storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store: store,
	}
}

func (s *APIServer) Run()  {
	router := mux.NewRouter()

	router.HandleFunc("/login",makeHttpHandleFunc(s.handleLogin))
	router.HandleFunc("/account",makeHttpHandleFunc(s.handleAccount))
	router.HandleFunc("/account/{id}",withJWTAuth(makeHttpHandleFunc(s.handleGetAccountByID),s.store))
	router.HandleFunc("/transfer",makeHttpHandleFunc(s.handleTransfer))

	log.Printf("API server running on %s", s.listenAddr)

	http.ListenAndServe(s.listenAddr, router)
}

func (s *APIServer) handleAccount(w http.ResponseWriter, r *http.Request) error {

	switch r.Method {
	case "GET":
		return s.handleGetAccount(w ,r)
	case "POST":
		return s.handleCreateAccount(w ,r)
	default:
		return fmt.Errorf("invalid method")
	}
}

// --------------------/  Create an account  /--------------------
func (s *APIServer) handleCreateAccount(w http.ResponseWriter, req *http.Request) error {
	createAccountReq := CreateAccountRequest{}
	if err := json.NewDecoder(req.Body).Decode(&createAccountReq); err != nil{
		return err
	}
	
	account, err := NewAccount(createAccountReq.FirstName,createAccountReq.LastName,createAccountReq.Password)
	if err != nil{
		return err
	}

	if err := s.store.CreateAccount(account); err != nil{
		return err
	}

	// tokenString,err := createJWT(account)
	// if err !=nil {
	// 	return err
	// }

	// fmt.Println("JWT TOKEN: ",tokenString)

	return WriteJson(w, http.StatusOK,account)
}

// 84130

func (s *APIServer) handleLogin(w http.ResponseWriter, req *http.Request) error{
	if req.Method == "POST" {
	return fmt.Errorf("invalid method")
	}
	
	loginRequest := LoginRequest{}
	if err := json.NewDecoder(req.Body).Decode(&loginRequest); err != nil {
		return err
	}

	acc, err := s.store.GetAccountByNumber(int(loginRequest.Number))
	if err != nil {
		return err
	}

	if !acc.ValidatePassword(loginRequest.Password){
		return fmt.Errorf("not authenticated")
	}

	token, err := createJWT(acc)
	if err != nil {
		return err
	}

	resp := LoginResponse{
		Token: token,
		Number: acc.Number,
	}


 return WriteJson(w,http.StatusOK,resp)
}

// --------------------/  Transfer Amt  /--------------------
func (s *APIServer) handleTransfer(w http.ResponseWriter, req *http.Request) error {
	transferRequest := TransferRequest{}
	if err := json.NewDecoder(req.Body).Decode(&transferRequest); err != nil{
		return err
	}

	defer req.Body.Close()
	return WriteJson(w,http.StatusOK,transferRequest)
}


// --------------------/  Get all accounts  /--------------------
func (s *APIServer) handleGetAccount(w http.ResponseWriter, req *http.Request) error {
	accounts,err := s.store.GetAccounts()
	if err != nil  {
		return err
	}
	return WriteJson(w,http.StatusOK,accounts)
}

// --------------------/   Get account by ID   /--------------------
func (s *APIServer) handleGetAccountByID(w http.ResponseWriter, req *http.Request) error {
	if req.Method == "GET" {
	id := mux.Vars(req)["id"]

	intId,err := strconv.Atoi(id)
	if err != nil {
		return err
	}
	
	
	account, err := s.store.GetAccountByID(intId)
	if err != nil {
		return err
	}

	return WriteJson(w,http.StatusOK,account)
}else if req.Method == "DELETE" {
	return s.handleDeleteAccount(w,req)
}

	return fmt.Errorf("invalid method")

}

func getID(r *http.Request)(int,error){
	id := mux.Vars(r)["id"]
	intId, err :=strconv.Atoi(id)
	if err != nil {
		return intId, err
	}
	return intId, nil
}

// --------------------/   Delete an account   /--------------------
func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, req *http.Request) error {
	id := mux.Vars(req)["id"]
	intId,err :=strconv.Atoi(id)
	if err != nil {
		return fmt.Errorf("invalid id %v",err)
	}

	if err = s.store.DeleteAccount(intId); err !=nil{
		return WriteJson(w,http.StatusBadRequest,err)
	}

	return WriteJson(w,http.StatusOK,id)
}



func WriteJson(w http.ResponseWriter, status int, v any) error{
	w.Header().Add("Content-Type","application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}


// --------------------/   CREATE JWT   /--------------------
func createJWT(account *Account) (string,error)  {
	claims := jwt.MapClaims{
		"expiresAt": 15000,
		"accountNumber": account.Number,
	}
	key:=os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,claims)

	return token.SignedString([]byte(key))
}


// --------------------/   JWT AUTHENTICATION MIDDLEWARE   /--------------------
func withJWTAuth(handleFunc http.HandlerFunc,s Storage) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) { 
		fmt.Println("Calling JWT auth middleware")
		tokenString := r.Header.Get("x-jwt-token")
		token,err := validateJWT(tokenString)
		if err != nil{
			WriteJson(w,http.StatusForbidden,ApiError{Error: "Premission denied"})
			return
		}
		if !token.Valid{
			WriteJson(w,http.StatusForbidden,ApiError{Error: "Premission denied"})
			return
		}
		userID, err := getID(r)
		if err!=nil{
			WriteJson(w,http.StatusForbidden,ApiError{Error: "Premission denied"})
			return
		}
		account, err := s.GetAccountByID(userID)
		if err!=nil{
			WriteJson(w,http.StatusForbidden,ApiError{Error: "Premission denied"})
			return
		}
		claims := token.Claims.(jwt.MapClaims)
		if account.Number != int64(claims["accountNumber"].(float64)) {
			WriteJson(w,http.StatusForbidden,ApiError{Error: "Premission denied"})
			return
		}
		
		handleFunc(w,r)
	}
}


// --------------------/   JWT VALIDATION   /--------------------
func validateJWT(tokentString string)(*jwt.Token,error){
	secret := os.Getenv("JWT_SECRET")
	return jwt.Parse(tokentString,func(token *jwt.Token) (interface{}, error) {
		if _,ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil,fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret),nil
	})
}




//this function is the type defination of apihandlefunc. We are using which returns and error but router.handlefunc doesn't accept this func
type apiFunc func(http.ResponseWriter, *http.Request) error

type ApiError struct{
	Error string
}

func makeHttpHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request)  {
		err:=f(w,r)
		if err != nil{
			WriteJson(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}
