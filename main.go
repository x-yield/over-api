package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/rakyll/statik/fs"
	"github.com/sirupsen/logrus"

	"github.com/utrack/clay/v2/log"
	"github.com/utrack/clay/v2/transport/middlewares/mwgrpc"
	"github.com/utrack/clay/v2/transport/server"
	// We're using statik-compiled files of Swagger UI
	// for the sake of example.
	_ "github.com/utrack/clay/doc/example/static/statik"

	"github.com/x-yield/over-api/internal/app/overload-service"
	"github.com/x-yield/over-api/service"
	"github.com/x-yield/over-api/tools"
)

func main() {
	// Wire up our bundled Swagger UI
	staticFS, err := fs.New()
	if err != nil {
		logrus.Fatal(err)
	}
	hmux := chi.NewRouter()
	hmux.Mount("/", http.FileServer(staticFS))

	db := tools.NewDbConnector()
	defer db.Close()
	influxdb := tools.NewInfluxDbConnector()
	defer influxdb.Close()
	s3 := tools.NewS3Service()
	defer s3.Close()

	overloadSrv := service.NewOverloadService(db, influxdb, s3)
	overloadImpl := overload.NewOverloadService(overloadSrv)

	srv := server.NewServer(
		7000,
		// Pass our mux with Swagger UI
		server.WithHTTPMux(hmux),
		// Recover from both HTTP and gRPC panics and use our own middleware
		server.WithGRPCUnaryMiddlewares(mwgrpc.UnaryPanicHandler(log.Default)),
	)
	err = srv.Run(overloadImpl)
	if err != nil {
		logrus.Fatal(err)
	}
}
