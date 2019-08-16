package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"go.opencensus.io/trace"
	"gocloud.dev/server"
	"gocloud.dev/server/sdserver"
)

type config struct {
	port int
}

func parseConfig() config {
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil || port == 0 {
		port = 8080
	}
	return config{
		port: port,
	}
}

// Handler just writes some bytes to a file
func Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", "inline; filename=\"test.dat\"")

	u, err := url.Parse(r.RequestURI)
	if err != nil {
		http.Error(w, "failed to parse query", http.StatusInternalServerError)
	}
	size, err := strconv.Atoi(u.Query().Get("size"))
	if err != nil {
		http.Error(w, "failed to parse size", http.StatusInternalServerError)
	}
	if size == 0 {
		size = 100
	}

	// TODO: Replace "A" with poop emoji and divide by 4
	for i := 0; i < size*1024*1024; i++ {
		w.Write([]byte("A"))
	}
	return
}

func listenVanilla() {
	args := parseConfig()
	fmt.Printf("Vanilla listening on port %d...\n", args.port)
	http.HandleFunc("/", Handler)
	err := http.ListenAndServe(
		fmt.Sprintf(":%d", args.port),
		nil,
	)

	if err != nil {
		log.Fatal(err)
	}
}

func listenCloud() {
	args := parseConfig()
	// Enable tracing for a small percentage of our requests by default
	var samplingPolicy trace.Sampler
	var exporter trace.Exporter
	samplingPolicy = trace.ProbabilitySampler(0.05)

	options := &server.Options{
		RequestLogger:         sdserver.NewRequestLogger(),
		TraceExporter:         exporter,
		DefaultSamplingPolicy: samplingPolicy,
		Driver:                &server.DefaultDriver{},
	}
	srv := server.New(http.DefaultServeMux, options)

	http.HandleFunc("/", Handler)

	fmt.Printf("Cloud listening on port %d...\n", args.port)
	log.Fatal(srv.ListenAndServe(fmt.Sprintf(":%d", args.port)))
}

func main() {
	listenVanilla()
}
