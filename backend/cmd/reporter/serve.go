package main

import (
	"fmt"
	"os"
	"net/http"
	"path/filepath"

	"github.com/yargevad/filepathx"
	"github.com/gorilla/mux"
	"github.com/urfave/cli/v2"

	amn "github.com/Fantom-foundation/LaScala/backend/utils"
)

var (
	ServePortFlag = cli.IntFlag {
		Name: 	"port",
		Usage: 	"port number to use",
		Aliases: []string{"p"},
		Value: 3333,
	}
	StaticDirFlag = cli.StringFlag {
		Name: 	"static-dir",
		Usage: 	"Directory path to be used as static directory",
		Aliases: []string{"d"},
		Value: "/home/rapolt/dev/static",
	}
	ReportDbFlag = cli.StringFlag {
		Name: 	"report-db",
		Usage: 	"Directory where a db is to be created.",
		Aliases: []string{"db"},
		Value: "/home/rapolt/dev/db",
	}

	RedisAddressFlag = cli.StringFlag {
		Name:   "redis-address",
		Usage:	"address for redis",
		Aliases: []string{"a", "ra"},
		Value: 	"127.0.0.1",
	}
	RedisPortFlag = cli.IntFlag {
		Name:   "redis-port",
		Usage:	"port for redis",
		Aliases: []string{"rp"},
		Value: 	6379,
	}
	RedisTopicFlag = cli.StringFlag {
		Name:   "redis-topic",
		Usage:	"topic for redis",
		Aliases: []string{"t", "rt"},
		Value: 	"task-queue",
	}
)

func ServeReporter(ctx *cli.Context) error {
	worker := amn.NewWorker(
		ctx.String("redis-address"),
		ctx.Int("redis-port"),
		ctx.String("redis-topic"),
	)

	return serveReporter(
		ctx.String("static-dir"),
		ctx.String("report-db"),
		ctx.Int("port"),
		worker,
	)
}

func serveReporter(static string, reportDb string, port int, worker *amn.Worker) error {
	r := mux.NewRouter()

	r.HandleFunc("/report/gen/{mid}/{rid}", generateReport(worker)).Methods("POST")
	r.HandleFunc("/gen/{mid}/{rid}", generateReport(worker)).Methods("POST")

	r.HandleFunc("/report/{mid}/{rid}", serveReport(static)).Methods("GET")
	r.HandleFunc("/run/{rid}", serveReport(static)).Methods("GET")

	r.HandleFunc("/report/{mid}/", serveReports(static)).Methods("GET")
	r.HandleFunc("/report/{mid}", serveReports(static)).Methods("GET")
	r.HandleFunc("/master/{mid}/", serveReports(static)).Methods("GET")
	r.HandleFunc("/master/{mid}", serveReports(static)).Methods("GET")

	r.HandleFunc("/report/", serveMasters).Methods("GET")
	r.HandleFunc("/report", serveMasters).Methods("GET")
	r.HandleFunc("/", serveMasters).Methods("GET")

	s := http.FileServer(http.Dir(static))
	r.PathPrefix("/static").Handler(http.StripPrefix("/static", s))

	return http.ListenAndServe(fmt.Sprintf(":%d", port), r)
}

func serve400(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprint(w, "400 malformed request\n")
	return
}


func serve404(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "404 page not found\n")
	return
}

func serveReport(static string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		rid, ok := vars["rid"]
		if !ok {
			serve400(w, r)
			return
		}

		mid, ok := vars["mid"]
		if !ok {
			pathToGlob := filepath.Join(static, "**", fmt.Sprintf("%s.html", rid))
			matches, err := filepathx.Glob(pathToGlob)
			if err != nil {
				serve400(w, r)
				return
			}
		
			if len(matches) == 0 {
				serve404(w, r)
				return
			}

			mid = filepath.Base(filepath.Dir(matches[0]))
		}

		pathToReport := filepath.Join(static, mid, fmt.Sprintf("%s.html", rid))
		_, err := os.Stat(pathToReport)
		if err != nil {
			serve404(w, r)
			return
		}

		http.ServeFile(w, r, pathToReport)
	}
}

func serveReports(static string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		mid, ok := vars["mid"]
		if !ok {
			serve404(w, r)
			return
		}
		
		pathToReports := filepath.Join(static, mid)
		_, err := os.Stat(pathToReports)
		if err != nil {
			serve404(w, r)
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/static/%s", mid), http.StatusFound)
	}
}

func serveMasters(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/static", http.StatusFound)
}

func generateReport(worker *amn.Worker) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		mid, ok := vars["mid"]
		if !ok {
			serve400(w, r)
			return
		}

		rid, ok := vars["rid"]
		if !ok {
			serve400(w, r)
		}

		fmt.Println("Here!", mid, rid)
		worker.Push(r.Context(), &amn.Task{
			Type: amn.TaskType("hello"),
			MasterId: mid,
			RunId: rid,
		})
	}
}

