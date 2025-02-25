package main

import (
	"example.com/Arifmetik/pkg/handler"
	"net/http"
)

func main() {

	http.HandleFunc("/update", handler.UpdateFunc)
	http.HandleFunc("/update_table_body2_data", handler.UpdateFunc)
	http.HandleFunc("/", handler.MainPage)
	http.HandleFunc("/switch", handler.TimeUpdate)
	http.HandleFunc("/switch2", handler.Switch2)

	http.ListenAndServe(":80", nil)
}
