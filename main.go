package main

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

type RawTransaction struct {
	ValueNate string
	ValueSim  string
	ValueSun  string
	Total     string
}

type Transaction struct {
	ValueNate float64
	ValueSim  float64
	ValueSun  float64
	Total     float64
}

func main() {
	conn, dbErr := pgx.Connect(context.Background(), "postgresql://postgres:postgres@localhost/split")
	if dbErr != nil {
		log.Printf("Splt: DB Failed to Connect %s", dbErr)
	}
	defer conn.Close(context.Background())

	tmpl := template.Must(template.ParseFiles("forms.html"))

	router := mux.NewRouter()

	router.HandleFunc("/split", func(w http.ResponseWriter, r *http.Request) {
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
		natevalue, _ := strconv.ParseFloat(rt.ValueNate, 64)
		simvalue, _ := strconv.ParseFloat(rt.ValueSim, 64)
		sunvalue, _ := strconv.ParseFloat(rt.ValueSun, 64)
		total, _ := strconv.ParseFloat(rt.Total, 64)
		transaction := Transaction{
			ValueNate: natevalue,
			ValueSim:  simvalue,
			ValueSun:  sunvalue,
			Total:     total,
		}
		// do something with details
		splits := []float64{total * (transaction.ValueNate / 100), total * (transaction.ValueSim / 100), total * (transaction.ValueSun / 100)}
		log.Printf("%.2f %.2f %.2f", splits[0], splits[1], splits[2])

		nateFV := total * (transaction.ValueNate / 100)
		_, err := conn.Exec(context.Background(), "INSERT INTO transactions(name, value, memo) values('nate', $1, $2)", nateFV, r.FormValue("memo"))
		if err != nil {
			log.Printf("err:%s\n", err)
			return
		}
		simFV := total * (transaction.ValueSim / 100)
		_, err = conn.Exec(context.Background(), "INSERT INTO transactions(name, value, memo) values('simran', $1, $2)", simFV, r.FormValue("memo"))
		if err != nil {
			log.Printf("err:%s\n", err)
			return
		}

		sunFV := total * (transaction.ValueSun / 100)
		_, err = conn.Exec(context.Background(), "INSERT INTO transactions(name, value, memo) values('sunjana', $1, $2)", sunFV, r.FormValue("memo"))
		if err != nil {
			log.Printf("err:%s\n", err)
			return
		}

		_ = tmpl.Execute(w, struct{ Success bool }{true})
	})

	router.HandleFunc("/split", Split)

	err := http.ListenAndServe(":8021", router)
	if err != nil {
		log.Println(err)
	}
}

func Split(w http.ResponseWriter, r *http.Request) {

}
