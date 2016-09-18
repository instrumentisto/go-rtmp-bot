package main

import (
	"flag"
	"github.com/instrumentisto/go-rtmp-bot"
	"github.com/instrumentisto/go-rtmp-bot/model"
	"github.com/instrumentisto/go-rtmp-bot/prometheus"
	"log"

	"github.com/instrumentisto/go-rtmp-bot/controller"
	"github.com/instrumentisto/go-rtmp-bot/redis"
	"github.com/instrumentisto/go-rtmp-bot/utils"
	"os"
)

var (
	test_launcher *rtmp_bot.Launcher // The application is a stress tester for rtmp media servers.
	report        *model.Report      // Test report value object.
	listenAddress = flag.String(
		"web.listen-address",
		":9132",
		"Address to listen on for web interface and telemetry.")
	metricPath = flag.String(
		"web.telemetry-path",
		"/metrics",
		"Path under which to expose metrics.")
	redis_url = flag.String("redis", "localhost:6379", "redis url")
	flvPath   = flag.String("flv_file", "", "Test flv file path")
	server    = flag.String("server", "stress_test", "Media server name")
	rtmp_url  = flag.String("rtmp_url",
		"rtmp://rtmp_server:1935/live",
		"RTMP Server application URL")
)

// Starts web interface for run stress tests.
// Man can open the web interface in browser with url "http://host:8082"
func main() {
	flag.Parse()
	defer os.Exit(1)
	if *flvPath == "" {
		log.Fatal("flv file not specified!")
		return
	}
	report = model.NewReport(*server)
	prometheus_client := prometheus.NewReportExportClient(
		*listenAddress, *metricPath, report)
	go prometheus_client.Run()
	log.Printf("listen redis: %v", *redis_url)
	app_handler := controller.AppHandler{
		Signal_chan: make(chan *model.Signal),
	}
	listener := redis.NewRedisListener(*redis_url, "", 0, app_handler)
	err := listener.WriteToMap("stress_test:status", *server, "ready")
	if err != nil {
		log.Printf("ERROR write status: %s", err.Error())
	}
	log.Printf("server: %s is ready", *server)
	defer listener.Close()
	defer listener.WriteToMap("stress_test:status", *server, "down")
	go listener.Listen()

	for {
		select {
		case signal := <-app_handler.Signal_chan:
			if signal.SignalType == redis.START_COMMAND {
				log.Println("HANDLE start test!!!")
				model_count, err := listener.Read("stress-test:model_count")
				if err != nil {
					log.Printf("Can not read models from readis: %v", err)
				}
				client_count, err := listener.Read("stress-test:client_count")
				if err != nil {
					log.Printf("Can not read clients from readis: %v", err)
				}

				start_request := new(model.StartRequest)
				start_request.ServerURL = *rtmp_url
				start_request.ModelCount = int(model_count)
				start_request.ClientCount = int(client_count)

				report.ResetReport(
					utils.GetUUID(),
					start_request.ModelCount, start_request.ClientCount)
				test_launcher = rtmp_bot.NewLauncher(
					start_request, report, *flvPath)
				go test_launcher.Start()
				listener.WriteToMap("stress_test:status", *server, "started")

			} else if signal.SignalType == redis.STOP_COMMAND {
				log.Println("HANDLE stop test")
				if test_launcher != nil {
					test_launcher.Stop()
					test_launcher = nil
					log.Println("test launcer stopped")
				}
				report.ResetReport("", 0, 0)
				listener.WriteToMap("stress_test:status", *server, "ready")
			}
		}
	}

}
