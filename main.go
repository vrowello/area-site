package main

import (
  "net/http"
  "html/template"
  "strconv"
  "fmt"
  "os"
  "log"
  "math"
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

  go apothem(sides, length, a)
  go perimeter(sides, length, p)

  result = area(<-a, <-p)

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

func apothem(n float64, l float64, a chan float64) {
	angle := ((180 - (360 / n)) / 2) * 0.01745329252 //finds the angle in degrees and converts to radians
	a <- math.Tan(angle) * (l / 2)
}

func area(apth float64, prmtr float64) float64 {
	return ((apth * prmtr) / 2)
}

func perimeter(n float64, l float64, p chan float64) {
	p <- n * l
}
