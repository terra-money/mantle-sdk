package rpc

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"

	"net/http"
	"sync"
	"time"
)

func StartRPCServer(
	port int,
	tctx *TestkitRPCContext,
) {
	r := mux.NewRouter()

	r.Use(MonitorIncomingRequestsMiddleware())

	// register testkit specific rpc
	RegisterTestkitRPC(r, tctx)

	// install mutex-guard across all endpoints
	r.Use(CreateMutexMiddleware(new(sync.Mutex)))
	r.Use(PanicToErrorMiddleware())

	server := http.Server{
		Handler:      r,
		Addr:         fmt.Sprintf("0.0.0.0:%d", port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Printf("[mantle/testkit-rpc] running on %d...", port)

	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
