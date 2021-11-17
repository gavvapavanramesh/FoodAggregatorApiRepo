package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/gorilla/mux"
)

type grains struct {
	ItemId   string
	ItemName string
	Quantity int
	Price    string
}
type vegetables struct {
	ProductId   string
	ProductName string
	Quantity    int
	Price       string
}
type fruits struct {
	Id       string
	Name     string
	Quantity int
	Price    string
}

const grainsUrl = "https://run.mocky.io/v3/e6c77e5c-aec9-403f-821b-e14114220148"
const vegetableUrl = "https://run.mocky.io/v3/4ec58fbc-e9e5-4ace-9ff0-4e893ef9663c"
const fruitsUrl = "https://run.mocky.io/v3/c51441de-5c1a-4dc2-a44e-aab4f619926b"

var wg sync.WaitGroup
var grain []grains
var vegetable []vegetables
var fruit []fruits

var cache1 []grains
var cache2 []vegetables
var cache3 []fruits

var summary []interface{}

func main() {
	fmt.Println("Server Started on port 4000..... ")

	r := mux.NewRouter()
	//routing
	r.HandleFunc("/", Home).Methods("GET")
	r.HandleFunc("/buy-item/{id}", GetByItemName).Methods("GET")
	r.HandleFunc("/buy-item-qty/{nam}/{quant}", GetByItemNameAndQuantity).Methods("GET")
	r.HandleFunc("/buy-item-qty-price/{nam}/{quant}/{pric}", GetByItemNameAndQuantityAndPrice).Methods("GET")
	r.HandleFunc("/show-summary/", GetSummary).Methods("GET")
	r.HandleFunc("/fast-buy-item/{id}", GetItemFast).Methods("GET")
	//listen to a port 4000
	log.Fatal(http.ListenAndServe(":4000", r))

	defer wg.Wait()
}

func Home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`<h1 > Food Aggregator </h1>`))
}

//Challenge 1

func GetByItemName(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get By item Name")

	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)

	getgrains(0)

	for _, g := range grain {

		if g.ItemName == params["id"] {
			json.NewEncoder(w).Encode(g)
			return
		}
	}

	getvegetables(0)

	for _, v := range vegetable {
		if v.ProductName == params["id"] {
			json.NewEncoder(w).Encode(v)
			return
		}
	}

	getfruits(0)

	for _, f := range fruit {
		if f.Name == params["id"] {
			json.NewEncoder(w).Encode(f)
			return
		}
	}

	json.NewEncoder(w).Encode("NOT_FOUND")

}

//Challenge 2

func GetByItemNameAndQuantity(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get By item Name And Quanity")
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	i, _ := strconv.Atoi(params["quant"])

	getgrains(0)

	for _, g := range grain {

		if g.ItemName == params["nam"] && i <= g.Quantity {
			json.NewEncoder(w).Encode(g)
			return
		}

	}

	getvegetables(0)

	for _, v := range vegetable {
		if v.ProductName == params["nam"] && i <= v.Quantity {
			json.NewEncoder(w).Encode(v)
			return
		}
	}

	getfruits(0)

	for _, f := range fruit {
		if f.Name == params["nam"] && i <= f.Quantity {
			json.NewEncoder(w).Encode(f)
			return
		}
	}

	json.NewEncoder(w).Encode(" NOT_FOUND")

}

//Challenge 3

func GetByItemNameAndQuantityAndPrice(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get By item Name And Quanity And Price")
	w.Header().Set("Content-Type", "application/json")
	var val string
	params := mux.Vars(r)
	i, _ := strconv.Atoi(params["quant"])
	if !strings.Contains(params["pric"], "$") {
		val = "$" + params["pric"]
	} else if strings.Index(params["pric"], "$") > 1 {
		k := strings.TrimRight(params["pric"], "$")
		val = "$" + k

	} else {
		val = params["pric"]
	}
	for _, c1 := range cache1 {
		if val == c1.Price && c1.ItemName == params["nam"] && i <= c1.Quantity {
			fmt.Println("from cache data structure")
			json.NewEncoder(w).Encode(c1)
			return
		}
	}
	for _, c2 := range cache2 {
		if val == c2.Price && c2.ProductName == params["nam"] && i <= c2.Quantity {
			fmt.Println("from cache data structure")
			json.NewEncoder(w).Encode(c2)
			return
		}
	}
	for _, c3 := range cache3 {
		if val == c3.Price && c3.Name == params["nam"] && i <= c3.Quantity {
			fmt.Println("from cache data structure")
			json.NewEncoder(w).Encode(c3)
			return
		}
	}
	getgrains(0)

	for _, g := range grain {
		if val == g.Price && g.ItemName == params["nam"] && i <= g.Quantity {
			cache1 = append(cache1, g)
			summary = append(summary, g)
			json.NewEncoder(w).Encode(g)
			return
		}

	}

	getvegetables(0)

	for _, v := range vegetable {
		if val == v.Price && v.ProductName == params["nam"] && i <= v.Quantity {
			cache2 = append(cache2, v)
			summary = append(summary, v)
			json.NewEncoder(w).Encode(v)
			return
		}
	}

	getfruits(0)
	for _, f := range fruit {
		if val == f.Price && f.Name == params["nam"] && i <= f.Quantity {
			cache3 = append(cache3, f)
			summary = append(summary, f)
			json.NewEncoder(w).Encode(f)
			return
		}
	}

	json.NewEncoder(w).Encode("NOT_FOUND")

}

//Challenge 4

func GetSummary(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Summary")

	json.NewEncoder(w).Encode(summary)

}

//challenge 5

func GetItemFast(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get item Name Quickly")
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)

	go getfruits(1)
	go getvegetables(1)
	go getgrains(1)
	wg.Add(3)

	for _, g := range grain {

		if g.ItemName == params["id"] {
			json.NewEncoder(w).Encode(g)
			return
		}

	}
	for _, v := range vegetable {
		if v.ProductName == params["id"] {
			json.NewEncoder(w).Encode(v)
			return
		}
	}
	for _, f := range fruit {
		if f.Name == params["id"] {
			json.NewEncoder(w).Encode(f)
			return
		}
	}
	json.NewEncoder(w).Encode("!! Not Found")

}

//helper functions

func getgrains(waitt int) {
	response, er := http.Get(grainsUrl)
	HandleErr(er)

	content, _ := ioutil.ReadAll(response.Body)

	err := json.Unmarshal(content, &grain)
	HandleErr(err)
	if waitt == 1 {
		wg.Done()
	}
	defer response.Body.Close()
}

func getvegetables(waitt int) {
	response1, er := http.Get(vegetableUrl)
	HandleErr(er)
	content1, _ := ioutil.ReadAll(response1.Body)

	err1 := json.Unmarshal(content1, &vegetable)
	HandleErr(err1)
	if waitt == 1 {
		wg.Done()
	}
	defer response1.Body.Close()
}

func getfruits(waitt int) {
	response2, er := http.Get(fruitsUrl)
	HandleErr(er)
	content2, _ := ioutil.ReadAll(response2.Body)

	err2 := json.Unmarshal(content2, &fruit)
	HandleErr(err2)
	if waitt == 1 {
		wg.Done()
	}
	defer response2.Body.Close()
}

//Error Handling Function

func HandleErr(err error) {
	if err != nil {
		panic(err)
	}
}
