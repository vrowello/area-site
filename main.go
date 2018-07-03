package main

import (
  "net/http"
  "html/template"
  "strconv"
  "fmt"
  "os"
  "log"
  "github.com/vrowello/reg-area/apothem"
  "github.com/vrowello/reg-area/area"
  "github.com/vrowello/reg-area/perimeter"
)

var (
 result float64
)

func main() {

  http.HandleFunc("/", server)
  http.ListenAndServe(GetPort(), nil)
}


type Area struct {
	Message string
	Ans  float64
}

type AreaData struct {
  Area_nums []Area
  Success bool
}

func server(w http.ResponseWriter, r *http.Request) {
  page := template.Must(template.ParseFiles("area-site.html"))
  if r.Method != http.MethodPost {
    page.Execute(w, nil)
      return
  }
  r.ParseForm()

  sides, err := strconv.ParseFloat(r.FormValue("SIDES"), 64)
  if err != nil {
    log.Fatal(err)
  }

  length, err := strconv.ParseFloat(r.FormValue("LENGTH"), 64)
  if err != nil {
    log.Fatal(err)
  }

  a := make(chan float64)
  p := make(chan float64)

  go apothem.Apothem(sides, length, a)
  go perimeter.Perimeter(sides, length, p)

  result = area.Area(<-a, <-p)

  output := AreaData{
    Area_nums: []Area{
		  {Message: "Area of Polygon:", Ans: result},
      },
          Success: true,
    }

  page.Execute(w, output)
}

func GetPort() string {
	var port = os.Getenv("PORT")
 	// Set a default port if there is nothing in the environment
 	if port == "" {
 		port = "4747"
 		fmt.Println("INFO: No PORT environment variable detected, defaulting to " + port)
 	}
 	return ":" + port
}
