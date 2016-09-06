package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"log"
	"net/http"
	"github.com/instrumentisto/go-rtmp-bot"
	"github.com/instrumentisto/go-rtmp-bot/prometheus"
	"github.com/instrumentisto/go-rtmp-bot/model"
	"github.com/instrumentisto/go-rtmp-bot/utils"
	"os"
)

var (
	test_launcher *rtmp_bot.Launcher // The application is a stress tester for rtmp media servers.
	report        *model.Report  // Test report value object.
	listenAddress = flag.String(
		"web.listen-address",
		":9132",
		"Address to listen on for web interface and telemetry.")
	metricPath = flag.String(
		"web.telemetry-path",
		"/metrics",
		"Path under which to expose metrics.")
	api_addrs = flag.String("api.addrs",":8083",
		"Address to listen http requests for API")
	flvPath = flag.String("flv_file", "","Test flv file path")
)

// Starts web interface for run stress tests.
// Man can open the web interface in browser with url "http://host:8082"
func main() {
	report = model.NewReport()
	flag.Parse()
	defer os.Exit(1)
	if(*flvPath == ""){
		log.Fatal("flv file not specified!")
		return
	}
	prometheus_client := prometheus.NewReportExportClient(
		*listenAddress, *metricPath, report)
	go prometheus_client.Run()
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/start_test", startTest)
	router.HandleFunc("/stop_test", stopTest)
	router.HandleFunc("/status", getStatus)
	log.Fatal(http.ListenAndServe(*api_addrs, router))
}

// Start test request handler.
// Runs stress test with request parameters.
// Before run test stops tests if the test is running.
func startTest(w http.ResponseWriter, r *http.Request) {
	writeAccessHeaders(w)
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
		fmt.Fprintln(w, model.GetResponse(2))
		return
	}
	decoder := schema.NewDecoder()
	start_request := new(model.StartRequest)
	err = decoder.Decode(start_request, r.PostForm)

	if err != nil {
		log.Fatal(err)
		fmt.Fprintln(w, model.GetResponse(2))
		return
	}
	report.ResetReport(
		utils.GetUUID(), start_request.ModelCount, start_request.ClientCount)
	if test_launcher != nil {
		fmt.Fprintln(w, model.GetResponse(2))
		return
	}
	test_launcher = rtmp_bot.NewLauncher(start_request, report,*flvPath)
	go test_launcher.Start()
	fmt.Fprintln(w, model.GetResponse(1))
}

// Stop test request handler.
// Stops all tests.
// Clean report.
func stopTest(w http.ResponseWriter, r *http.Request) {
	if test_launcher != nil {
		go test_launcher.Stop()
		test_launcher = nil
	}
	report.ResetReport("", 0, 0)
	writeAccessHeaders(w)
	fmt.Fprintln(w, model.GetResponse(0))
}

// Returns status of current test.
func getStatus(w http.ResponseWriter, r *http.Request) {
	writeAccessHeaders(w)
	if test_launcher != nil {
		fmt.Fprintln(w, model.GetResponse(1))
		return
	}
	fmt.Fprintln(w, model.GetResponse(0))
}

// Writes CORS headers to HTTP response.
func writeAccessHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods",
		"POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers",
		"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}
