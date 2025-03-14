package main

import (
	"html/template"
	"log"
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

var HTMLTemplate = template.Must(template.New("VirtualPage").Parse(`
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8">
    <title>Test Page</title>
  </head>
  <body>
    {{range .}}
		<div><a href="/{{.}}" > /{{.}} </a></div>
	{{end}}
  </body>
</html>
`))

type VirtualPage struct {
	Targets []string
}

type VirtualWorldWideWeb struct {
	numberOfPages int
	seed          int64
	delay         time.Duration
	rand          *rand.Rand
}

func NewVWWW(numberOfPages int, seed int64, delayMS int) *VirtualWorldWideWeb {
	return &VirtualWorldWideWeb{
		numberOfPages: numberOfPages,
		seed:          seed,
		delay:         time.Duration(delayMS) * time.Millisecond,
		rand:          rand.New(rand.NewSource(seed)),
	}
}

func (vwww *VirtualWorldWideWeb) Serve(port int) error {
	http.HandleFunc("/{id}", vwww.renderPage)
	http.HandleFunc("/", vwww.renderIndex)

	log.Printf("─────────────────────────────────────────────\n")
	log.Printf(" Serving requests on http://127.0.0.1:%d\n", port)
	log.Printf(" ↳ numberOfPages=%d\n", vwww.numberOfPages)
	log.Printf(" ↳ seed=%d\n", vwww.seed)
	log.Printf(" ↳ delay=%s\n", vwww.delay)
	log.Printf("─────────────────────────────────────────────\n")
	return http.ListenAndServe(":"+strconv.Itoa(port), nil)
}

func (vwww *VirtualWorldWideWeb) renderIndex(w http.ResponseWriter, req *http.Request) {
	vwww.rand.Seed(vwww.seed)
	seeds := make([]int, 10)
	for i := 0; i < 10; i++ {
		seeds[i] = vwww.rand.Intn(vwww.numberOfPages)
	}
	HTMLTemplate.Execute(w, seeds)
	w.Header().Add("Content-Type", "text/html")
	w.Header().Add("X-Seed", strconv.Itoa(int(vwww.seed)))
	log.Printf("200 OK - %s\n", req.URL)
}

func (vwww *VirtualWorldWideWeb) renderPage(w http.ResponseWriter, req *http.Request) {
	time.Sleep(vwww.delay)
	idStr := req.PathValue("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id < 0 || id >= int64(vwww.numberOfPages) {
		w.WriteHeader(404)
		w.Write([]byte("Not Found"))
		log.Printf("404 Not Found - %s\n", req.URL)
		return
	}

	vwww.rand.Seed(vwww.seed + id)

	var alpha float64
	var max float64
	if vwww.numberOfPages <= 100_000 {
		alpha = 1.5
		max = float64(vwww.numberOfPages) / 5
	} else {
		alpha = 2.0
		max = float64(vwww.numberOfPages) / 10
	}
	nbTargets := randomPowerLaw(alpha, 1, max, *vwww.rand)

	targets := make([]int, nbTargets+1)
	for i := 0; i < nbTargets; i++ {
		targets[i] = vwww.rand.Intn(vwww.numberOfPages)
	}
	targets[nbTargets] = int(id - 1)
	HTMLTemplate.Execute(w, targets)
	w.Header().Add("Content-Type", "text/html")
	w.Header().Add("X-Seed", strconv.Itoa(int(vwww.seed)))
	log.Printf("200 OK - %s\n", req.URL)
}

func randomPowerLaw(alpha float64, min float64, max float64, rand rand.Rand) int {
	if min <= 0 {
		panic("min must be greater than 0 to ensure the power-law is well-defined")
	}
	u := rand.Float64()
	power := 1 - alpha
	x := math.Pow(u*(math.Pow(max, power)-math.Pow(min, power))+math.Pow(min, power), 1/power)
	return int(math.Round(x))
}
