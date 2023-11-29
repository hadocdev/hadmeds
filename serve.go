package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	//"log"
	//"net/http"
	//"github.com/gorilla/mux"
)


func main() {
	router := gin.Default()
	router.GET("/api/search",queryData)

	router.Run("localhost:8000")
}

func queryData(c *gin.Context){
	query := c.Query("query")

	db, err := sql.Open("sqlite3","../hadoc-data/drugs.db")
	checkErr(err)

	s := fmt.Sprintf(`
			SELECT mfg, brand, generic, strength, dosage, price
			FROM Drugs
			WHERE
			 brand LIKE "%%%s%%" OR
			 generic LIKE "%%%s%%"
			ORDER BY brand
			LIMIT 20	
		`, query)
	rows, err := db.Query(s)
	checkErr(err)

	var mfg, brand, generic, strength, dosage, price string
	var response []map[string] any

	for rows.Next(){
		err = rows.Scan(&mfg, &brand, &generic, &strength, &dosage, &price)
		checkErr(err)
		response = append(response, map[string]any{"mfg":mfg, "brand":brand, "generic":generic, "strength":strength, "dosage":dosage, "price":price})
	}

	rows.Close()
	db.Close()

	c.IndentedJSON(200,response)
}

func checkErr (err error){
	if err != nil {
		panic(err)
	}
}

/*func main(){
	r := mux.NewRouter()
	r.HandleFunc("/{path}", indexHandler)
	log.Fatal(http.ListenAndServe(":8080",r))
}

func indexHandler(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	path := vars["path"]
	fmt.Fprintf(w, "You're looking for the URL %s\n", path)
}
*/



















