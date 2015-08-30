package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
)

/*
* curl -u dev:dev -v 'http://localhost:8080/addresses' -G --data-urlencode "per=1" --data-urlencode "page=1" --data-urlencode "q=united" --data-urlencode 'cols=["country", "line1"]' --data-urlencode 'condOp=OR'
 */

func main() {
	path := flag.String("path", "", "The path to be queried: 'http://host:port/[path]?params'")
	per := flag.Int("per", 10, "The max number of results to display.")
	page := flag.Int("page", 1, "The page number to paginate.")
	q := flag.String("q", "", "The search string")
	cols := flag.String("cols", "", "A comma delimited list of column names to search into.")
	condOp := flag.String("condOp", "OR", "Conditional operator, AND or OR")
	matchOp := flag.String("matchOp", "=", "Match operator, = or LIKE or ILIKE")
	flag.Parse()

	urlBase := "http://localhost:8080/api/"

	user := os.Getenv("CRUD_USER")
	pw := os.Getenv("CRUD_PW")

	if user == "" {
		log.Fatal("Missing env var: CRUD_USER")
	}
	if pw == "" {
		log.Fatal("Missing env var: CRUD_PW")
	}

	params := fmt.Sprintf("?per=%d&page=%d", *per, *page)

	if *q != "" {
		*q = url.QueryEscape(*q)
		params += "&q=" + *q
	}

	if *cols != "" {
		params += "&cols=" + *cols
	}

	if *condOp != "AND" && *condOp != "OR" {
		*condOp = "OR"
	}
	params += "&condOp=" + *condOp + "&matchOp=" + url.QueryEscape(*matchOp)

	u := urlBase + *path + params

	client := &http.Client{}
	req, err := http.NewRequest("GET", u, nil)
	req.SetBasicAuth(user, pw)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	fmt.Println("--------------------------------------------------------------------------------------------")
	fmt.Println(u)
	fmt.Println("--------------------------------------------------------------------------------------------")

	fmt.Println(string(body))
}
