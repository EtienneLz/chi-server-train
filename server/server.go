package server

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"io"
	"log"
	"net/http"
	"strconv"
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
		w.Write([]byte("Hello World !"))
	})

	r.Get("/order", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Order here"))
	})
	r.Post("/order", retrieveSend)
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
			return
		}

		marshal, err := json.Marshal(refundStruct)
		if err != nil {
			return
		}
		_, err = w.Write(marshal)
		if err != nil {
			return
		}
	}
}
