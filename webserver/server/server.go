package server

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"click-analytics/config"

	"github.com/gorilla/mux"
)

type LogResponse struct {
	http.ResponseWriter
	StatusCode int
	BytesSent  int
}

func NewLogResponseWriter(w http.ResponseWriter) *LogResponse {
	return &LogResponse{ResponseWriter: w}
}

func (w *LogResponse) WriteHeader(code int) {
	w.StatusCode = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *LogResponse) Write(body []byte) (int, error) {
	w.BytesSent += len(body)
	return w.ResponseWriter.Write(body)
}

type KafkaHandler struct {
	handler http.Handler
}

//ServeHTTP handles the request by passing it to the real handler and kafka
func (l *KafkaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	logWriter := NewLogResponseWriter(w)
	l.handler.ServeHTTP(logWriter, r)
	fmt.Printf("status: %d , bytes_sent:%d, elapsed: %d\n", logWriter.StatusCode, logWriter.BytesSent, time.Since(start))
}

//NewKafkaHandler constructs a new KafkaHandler middleware handler
func NewKafkaHandler(handlerToWrap http.Handler) *KafkaHandler {
	return &KafkaHandler{handler: handlerToWrap}
}

func Start() {
	conf := config.Get()

	r := mux.NewRouter().StrictSlash(true)

	// Web Content
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./server/web/")))

	rKafka := NewKafkaHandler(r)
	fmt.Printf("Listen in %s\n", conf.Server.Port)
	log.Fatal(http.ListenAndServe(conf.Server.Port, rKafka))

}
