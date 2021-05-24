package main

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

type RawTransaction struct {
	ValueNate string
	ValueSim string
	ValueSun string
	Total string
}

type Transaction struct {
	ValueNate float64
	ValueSim float64
	ValueSun float64
	Total float64
}

func main(){
	conn, dbErr := pgx.Connect(context.Background(), "postgresql://postgres:postgres@localhost/split")
	if dbErr != nil {
		log.Printf("Splt: DB Failed to Connect %s", dbErr)
	}
	defer conn.Close(context.Background())


	tmpl := template.Must(template.ParseFiles("forms.html"))


	router := mux.NewRouter()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			tmpl.Execute(w, nil)
			return
		}

		rt := RawTransaction{
			ValueNate: r.FormValue("natevalue"),
			ValueSim:  r.FormValue("simvalue"),
			ValueSun:  r.FormValue("sunvalue"),
			Total:     r.FormValue("total"),
		}
		natevalue, _  :=  strconv.ParseFloat(rt.ValueNate, 64)
		simvalue, _  :=  strconv.ParseFloat(rt.ValueSim, 64)
		sunvalue, _  :=  strconv.ParseFloat(rt.ValueSun, 64)
		total, _  :=  strconv.ParseFloat(rt.Total, 64)
		transaction := Transaction{
			ValueNate: natevalue,
			ValueSim:  simvalue,
			ValueSun:  sunvalue,
			Total:     total,
		}
		// do something with details
		splits := []float64{total*(transaction.ValueNate/100), total*(transaction.ValueSim/100), total*(transaction.ValueSun/100)}
		fmt.Printf("%.2f %.2f %.2f", splits[0], splits[1], splits[2])

		_ = tmpl.Execute(w, struct{ Success bool }{true})
	})

	router.HandleFunc("/split", Split)

	_ = http.ListenAndServe(":8080", router)
}

func Split(w http.ResponseWriter, r *http.Request){

}
