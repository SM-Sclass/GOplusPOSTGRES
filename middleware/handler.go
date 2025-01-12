package middleware

import (
	"PostgreSQLstockbackend/models"
	"database/sql"
	_ "github.com/lib/pq" 
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

type response struct {
	ID      int64  `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
}

func CreateConnection() *sql.DB {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error in loading the env file")
	}

	dbURL := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        os.Getenv("DB_HOST"),
        os.Getenv("DB_PORT"),
        os.Getenv("DB_USER"),
        os.Getenv("DB_PASSWORD"),
        os.Getenv("DB_NAME"))
	db, err := sql.Open("postgres", dbURL)

	if err != nil {
		log.Fatalf("Error in creating the connection, %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Error in pinging the connection. %v", err)
	}

	fmt.Println("Successfully connected!")

	return db
}

func CreateStock(w http.ResponseWriter, r *http.Request) {
	var stock models.Stock
	err := json.NewDecoder(r.Body).Decode(&stock)

	if err != nil {
		log.Fatalf("Error in decoding the request body. %v", err)
	}

	insertID := insertStock(stock)
	res := response{
		ID:      insertID,
		Message: "Stock created succesfully",
	}

	json.NewEncoder(w).Encode(res)
}

func GetStock(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("Error in converting the string to int. %v", err)
	}
	stock, err := getStock(int64(id))
	if err != nil {
		log.Fatalf("Error in getting the stock detail. %v", err)
	}

	json.NewEncoder(w).Encode(stock)
}

func GetAllStocks(w http.ResponseWriter, r *http.Request) {
	stocks, err := getAllStocks()

	if err != nil {
		log.Fatalf("Error in gettign all stock details")
	}

	json.NewEncoder(w).Encode(stocks)
}

func UpdateStock(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Fatalf("Error in converting the string to int. %v", err)
	}
	var stock models.Stock

	json.NewDecoder(r.Body).Decode(&stock)

	if err != nil {
		log.Fatalf("Error in decoding the request body. %v", err)
	}

	updatedRows := updateStock(int64(id), stock)

	msg := fmt.Sprintf("Stock updated successfully %v row(s) affected", updatedRows)

	res := response{
		ID:      int64(id),
		Message: msg,
	}
	json.NewEncoder(w).Encode(res)
}

func DeleteStock(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("Error in converting the string to int. %v", err)
	}

	deletedRows := deleteStock(int64(id))

	msg := fmt.Sprintf("Stock deleted successfully %v row(s) affected", deletedRows)

	res := response{
		ID:      int64(id),
		Message: msg,
	}

	json.NewEncoder(w).Encode(res)
}

func insertStock(stock models.Stock) int64 {
	db := CreateConnection()
	defer db.Close()
	sqlQuery := `INSERT INTO stocks (name, price, company) VALUES($1, $2, $3) RETURNING stockid`
	var id int64

	err := db.QueryRow(sqlQuery, stock.Name, stock.Price, stock.Company).Scan(&id)

	if err != nil {
		log.Fatalf("Error in executing the query. %v", err)
	}
	fmt.Printf("Inserted a single record %v", id)
	return id
}

func getStock(id int64) (models.Stock, error) {
	db := CreateConnection()
	defer db.Close()

	var stock models.Stock

	sqlQuery := `SELECT * FROM stocks WHERE stockid=$1`

	row := db.QueryRow(sqlQuery, id)

	err := row.Scan(&stock.StockID, &stock.Name, &stock.Price, &stock.Company)

	switch err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned")
		return stock, nil

	case nil:
		return stock, nil
	default:
		log.Fatalf("Unable to scan the row . %v", err)
	}

	return stock, err
}

func getAllStocks() ([]models.Stock, error) {
	db := CreateConnection()
	defer db.Close()

	var stocks []models.Stock

	sqlQuery := `SELECT * FROM stocks`

	rows, err := db.Query(sqlQuery)

	if err != nil {
		log.Fatalf("Error in executing the query. %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var stock models.Stock
		err = rows.Scan(&stock.StockID, &stock.Name, &stock.Price, &stock.Company)
		if err != nil {
			log.Fatalf("Error in scanning the row. %v", err)
		}

		stocks = append(stocks, stock)
	}

	return stocks, nil
}

func updateStock(id int64, stock models.Stock) int64 {
	db := CreateConnection()
	defer db.Close()

	sqlQuery := `UPDATE stocks SET name=$1, price=$2, company=$3 WHERE stockid=$4`

	res, err := db.Exec(sqlQuery, stock.Name, stock.Price, stock.Company, id)

	if err != nil {
		log.Fatalf("Error in executing the query. %v", err)
	}

	rowsAffected, err := res.RowsAffected()

	if err != nil {
		log.Fatalf("Error in getting the rows affected. %v", err)
	}
	fmt.Printf("Rows affected %v", rowsAffected)
	return rowsAffected
}

func deleteStock(id int64) int64 {
	db := CreateConnection()
	defer db.Close()

	sqlQuery := `DELETE FROM stocks WHERE stockid=$1`

	res, err := db.Exec(sqlQuery, id)

	if err != nil {
		log.Fatalf("Error in executing the query. %v", err)
	}

	rowsAffected, err := res.RowsAffected()

	if err != nil {
		log.Fatalf("Error in getting the rows affected. %v", err)
	}
	fmt.Printf("Rows affected %v", rowsAffected)
	return rowsAffected
}
