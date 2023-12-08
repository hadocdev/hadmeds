package main

import (
	 //"fmt"
	"github.com/gin-gonic/gin"
	"database/sql"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	//"log"
	"math"
	"net/http"
	"slices"
	//"github.com/gorilla/mux"
)


func main() {
	router := gin.Default()
	router.GET("/api/search",queryData)

	router.Run("localhost:8000")
}

func queryData(c *gin.Context){
	query := c.Query("query")

	db, err := sqlx.Open("sqlite3","../hadoc-data/drugs.db")
	checkErr(err)

	total_row := db.QueryRow("SELECT COUNT(DISTINCT brand) FROM Drugs;")
	var total int
	_ = total_row.Scan(&total)
	
	s := `SELECT mfg, brand, generic, strength, dosage, price 
		  FROM Drugs
		  LIMIT ?`
	rows, err := db.Query(s,total)
	checkErr(err)

	var mfg, brand, generic, strength, dosage, price sql.NullString
	var response []map[string] any

	for rows.Next(){
		err = rows.Scan(&mfg, &brand, &generic, &strength, &dosage, &price)
		checkErr(err)
		response = append(response, map[string]any{"mfg":mfg.String, "brand":brand.String, "generic":generic.String, "strength":strength.String, "dosage":dosage.String, "price":price.String})
	}
	
	slices.SortFunc(response, levenshteinCompare(query))

	if len(response) > 20{
		response = response[:20]
	}
	
	rows.Close()
	db.Close()

	c.IndentedJSON(http.StatusOK,response)
}

func checkErr (err error){
	if err != nil {
		panic(err)
	}
}

func editDist (s1, s2 string) int{
	// initialising the data
	prev_line := make([]int, len(s1)+1)
	curr_line := make([]int, len(s1)+1)
	for i:=0; i<len(prev_line); i++{
		prev_line[i] = i
	}
	var to_left, to_up, to_corner, curr_inequal int

	// running the DP algorithm
	for i:=1; i<len(s2)+1; i++{
		curr_line[0] = i
		for j:=1; j<len(s1)+1; j++{
			to_left = curr_line[j-1]
			to_up = prev_line[j]
			to_corner = prev_line[j-1]

			if s1[j-1] == s2[i-1] {
				curr_inequal = 0
			} else {
				curr_inequal = 1
			}

			curr_line[j] = minimum([]int{to_corner+curr_inequal, to_left+1, to_up+1})
		}
			copy(prev_line, curr_line)
	}

	return curr_line[len(s1)]
}

func minimum (nums []int) int{
	min_num := math.MaxInt64
	for _, n := range nums{
		if n < min_num{
			min_num = n
		}
	}
	return min_num
}

func levenshteinCompare (q string) func (map[string] any, map[string] any) int{
	cmpFunc := func (e1, e2 map[string] any) int{
		dist_e1_brand := editDist(q, e1["brand"].(string))
		dist_e2_brand := editDist(q, e2["brand"].(string))
		dist_e1_generic := editDist(q, e1["generic"].(string))
		dist_e2_generic := editDist(q, e2["generic"].(string))

		min_n := minimum([]int{dist_e1_brand - dist_e2_brand, dist_e1_generic - dist_e2_brand, dist_e1_brand - dist_e2_generic, dist_e1_generic - dist_e2_generic})
		return min_n
	}

	return cmpFunc
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
