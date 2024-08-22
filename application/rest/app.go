package rest

import (
	"antrein/bc-dashboard/application/common/resource"
	"antrein/bc-dashboard/application/common/usecase"
	"antrein/bc-dashboard/internal/handler/grpc/analytic"
	"antrein/bc-dashboard/internal/handler/rest/auth"
	"antrein/bc-dashboard/internal/handler/rest/project"
	"antrein/bc-dashboard/model/config"
	"compress/gzip"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

func setupCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE, PATCH")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
}

type gzipResponseWriter struct {
	http.ResponseWriter
	Writer *gzip.Writer
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func compressHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}
		w.Header().Set("Content-Encoding", "gzip")
		gz, _ := gzip.NewWriterLevel(w, gzip.BestSpeed)
		defer gz.Close()

		gzw := gzipResponseWriter{ResponseWriter: w, Writer: gz}
		next.ServeHTTP(gzw, r)
	})
}

func ApplicationDelegate(cfg *config.Config, uc *usecase.CommonUsecase, rsc *resource.CommonResource) (http.Handler, error) {
	router := mux.NewRouter()

	router.HandleFunc("/bc/dashboard/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Makan nasi pagi-pagi, ngapain kamu disini?")
	})
	router.HandleFunc("/bc/dashboard/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "pong!")
	})

	// routes

	// auth
	authRoute := auth.New(cfg, uc.AuthUsecase, rsc.Vld)
	authRoute.RegisterRoute(router)

	// project
	projectRoute := project.New(cfg, uc.ProjectUsecase, uc.ConfigUsecase, rsc.Vld)
	projectRoute.RegisterRoute(router)

	// analytic
	analyticRouter := analytic.New(cfg, rsc.GRPC)
	analyticRouter.RegisterRoute(router)

	handlerWithMiddleware := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		setupCORS(w)

		if r.Method == "OPTIONS" {
			return
		}

		compressHandler(router).ServeHTTP(w, r)
	})

	return handlerWithMiddleware, nil
}
