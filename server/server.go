package server

import (
	"database/sql"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/go-sql-driver/mysql"
	"go_train/data"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

type orderInfo struct {
	RemoteCustomerReference string   `json:"remote-customer-reference"`
	PurchaseList            []string `json:"purchase-list"`
	Amount                  string   `json:"amount,omitempty"`
}

type refund struct {
	NumberOfPurchase int `json:"number-of-purchase"`
	RefundAmount     int `json:"refund-amount,omitempty"`
}

type errorStruct struct {
	ErrorCode    int    `json:"error-code"`
	ErrorMessage string `json:"error-message"`
}

var error400 = errorStruct{ErrorCode: 400, ErrorMessage: "Bad Request"}

func Init() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("Hello World !"))
		if err != nil {
			return
		}
	})

	r.Get("/order", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("Order here"))
		if err != nil {
			return
		}
	})
	r.Post("/order", retrieveSend)
	println("Server started")
	err := http.ListenAndServe(":3000", r)
	if err != nil {
		return
	}
}

func retrieveSend(w http.ResponseWriter, r *http.Request) {
	var order orderInfo

	err := r.ParseForm()
	if err != nil {
		return
	}
	b, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Read failed")
		return
	}
	//log.Println(b)
	err = json.Unmarshal(b, &order)
	if err != nil {
		log.Println("Unmarshall failed", err)
		return
	}
	amountCheck, err := strconv.Atoi(order.Amount)
	if err != nil || amountCheck < 0 {
		order.Amount = ""
		println("Wrong amount format")
	}
	if order.RemoteCustomerReference == "" {
		//http.Error(w, error400.ErrorMessage, error400.ErrorCode)
		w.WriteHeader(error400.ErrorCode)
		_, err := w.Write([]byte(error400.ErrorMessage + "\n"))
		if err != nil {
			return
		}

		marshal, err := json.Marshal(error400)
		if err != nil {
			return
		}
		_, err = w.Write(marshal)
		if err != nil {
			return
		}
		return
	}
	_, err = w.Write([]byte("Customer: " + order.RemoteCustomerReference + "\nPurchase list: \n"))
	if err != nil {
		return
	}
	i := 0
	for _, str := range order.PurchaseList {
		i++
		_, err := w.Write([]byte(str + "\n"))
		if err != nil {
			return
		}
	}
	_, err = w.Write([]byte("Amount :" + order.Amount + "\n"))
	if err != nil {
		return
	}
	var refundStruct refund
	refundStruct.NumberOfPurchase = i
	if order.Amount != "" {
		amount, err := strconv.Atoi(order.Amount)
		if err != nil {
			println("Wrong amount format")
		} else {
			refundStruct.RefundAmount = amount / 2.0
		}
	}

	addToDb(order.RemoteCustomerReference)

	marshal, err := json.Marshal(refundStruct)
	if err != nil {
		return
	}
	_, err = w.Write(marshal)
	if err != nil {
		return
	}
}

func addToDb(customerRef string) {
	currentTime := time.Now()
	dateFormatted := currentTime.Format("02-01-2006")
	_, err := data.G_DB.Exec("INSERT INTO user_update (remote_customer_reference, last_update_date) VALUES (?, ?) "+
		"ON DUPLICATE KEY UPDATE last_update_date = ?", customerRef, dateFormatted, dateFormatted)
	if err != nil {
		println(err.Error())
		return
	}

}

func DbConnect(name string) *sql.DB {
	db, err := sql.Open("mysql", "root:pass@tcp(0.0.0.0:3308)/"+name)
	if err != nil {
		panic(err.Error())
	}
	//defer db.Close()

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS " + name + ".user_update(remote_customer_reference varchar(50), last_update_date varchar(10), PRIMARY KEY (remote_customer_reference))")
	if err != nil {
		panic(err.Error())
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS " + name + ".user_purchase(remote_customer_reference varchar(50), purchase_lists varchar(10000))")
	if err != nil {
		panic(err.Error())
	}

	println("Connection to Database successful")
	return db
}
