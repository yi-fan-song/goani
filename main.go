package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"time"
)

const logDir = "./logs/"

func main() {
	// setup
	os.Mkdir(logDir, 0755)
	f, err := createFile(logDir + "latest.log")
	if err != nil {
		log := DefaultLogger()
		log.Error(fmt.Sprint(err))
	}
	defer func() {
		copyLogs(f)

		f.Close()
	}()

	log := &Logger{
		UseColor:    false,
		IncludeTime: true,
		Out:         &CombinedWriter{Writer1: os.Stdout, Writer2: f},
	}

	cfg, err := LoadOrCreateConfig()
	if err != nil {
		log.Warn(fmt.Sprintf("failed to create config: %v", err))
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	go func() {
		for sig := range sigChan {
			if sig == os.Interrupt {
				log.Info("Closing due to ^C signal")
				copyLogs(f)

				os.Exit(0)
			}
		}
	}()

	registerHandlers(cfg.Folders, log)

	// listen for connections
	log.Info(fmt.Sprintf("Listening on :%d", cfg.Port))
	log.Fatal(fmt.Sprint(http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), nil)))
}

func registerHandlers(paths []string, log *Logger) {
	regex, err := regexp.Compile(`(?:^|\/|\.\/)((?:[A-z]| )+)(?:\/|$)`)
	if err != nil {
		panic(err)
	}

	for _, path := range paths {
		if path == "." || path == "./" {
			http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

				lrw := NewLoggingResponseWriter(w)
				http.FileServer(http.Dir(".")).ServeHTTP(lrw, r)

				log.Info(fmt.Sprintf("%s %s %d", r.Method, r.URL, lrw.statusCode))
			}))
		} else {
			uri := fmt.Sprintf("/%s/", regex.FindStringSubmatch(path)[1])
			http.Handle(uri, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				originalURL := *r.URL

				if strings.HasPrefix(r.URL.Path, uri) {
					r.URL.Path = strings.Replace(r.URL.Path, uri, "/", 1)
				}

				lrw := NewLoggingResponseWriter(w)
				http.FileServer(http.Dir(path)).ServeHTTP(lrw, r)

				log.Info(fmt.Sprintf("%s %s %d", r.Method, &originalURL, lrw.statusCode))
			}))
		}
	}
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	// WriteHeader(int) is not called if our response implicitly returns 200 OK, so
	// we default to that status code.
	return &loggingResponseWriter{w, http.StatusOK}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func copyLogs(l *os.File) {
	log := DefaultLogger()
	// copy file
	logfile, err := os.Create(time.Now().Format(logDir+"2006-01-02T15:04:05") + ".log")
	if err != nil {
		log.Error(fmt.Sprint(err))
	}
	_, err = l.Seek(0, 0)
	if err != nil {
		log.Error(fmt.Sprint(err))
	}

	_, err = io.Copy(logfile, l)
	if err != nil {
		log.Error(fmt.Sprint(err))
	}

	logfile.Close()
}
