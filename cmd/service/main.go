package main

import (
	"flag"
	"github.com/dgraph-io/badger/v3"
	"github.com/itsabgr/omp/internal/db"
	"github.com/itsabgr/omp/internal/handler"
	"github.com/itsabgr/omp/internal/utils"
	"log"
	"net/http"
	"time"
)

var flagAddr = flag.String("addr", ":4444", "service address")
var flagDir = flag.String("dir", "", "data directory")

func init() {
	flag.Parse()
	log.Println("addr", *flagAddr)
	log.Println("dir", *flagDir)

}
func main() {
	defer func() {
		if rec := recover(); rec != nil {
			log.Panicln(rec)
		}
	}()
	badger1 := utils.Must(badger.Open(badger.DefaultOptions(*flagDir).WithLogger(nil).WithInMemory(*flagDir == "")))
	defer badger1.Close()
	utils.Throw((&http.Server{
		Addr:              *flagAddr,
		Handler:           handler.New(db.New(badger1)),
		IdleTimeout:       time.Second * 5,
		WriteTimeout:      time.Second * 5,
		ReadHeaderTimeout: time.Second * 2,
		ReadTimeout:       time.Second * 5,
		MaxHeaderBytes:    1000,
		ErrorLog:          nil,
	}).ListenAndServe())
}
