package rpc

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"sync"
)

func MustMarshalJSON(v interface{}) json.RawMessage {
	r, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}

	return r
}

func PanicToErrorMiddleware() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var err error
			defer func() {
				r := recover()
				if r != nil {
					switch t := r.(type) {
					case string:
						err = errors.New(t)
					case error:
						err = t
					default:
						err = errors.New("Unknown error")
					}

					response := new(struct {
						Error string `json:"error"`
					})
					response.Error = err.Error()

					http.Error(w, string(MustMarshalJSON(response)), http.StatusBadRequest)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}

func CreateMutexMiddleware(m *sync.Mutex) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			m.Lock()
			defer m.Unlock()

			next.ServeHTTP(w, r)

		})
	}
}

func FailIfGenesisNotInitializedMiddleware(ctx *TestkitRPCContext) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if ctx.tg == nil {
				response := new(struct {
					Error string `json:"error"`
				})

				response.Error = "genesis not initialized yet"
				http.Error(w, string(MustMarshalJSON(response)), http.StatusBadRequest)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func MonitorIncomingRequestsMiddleware() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("[mantle/testkit-rpc] incoming requests %s\n", r.RequestURI)

			next.ServeHTTP(w, r)
		})
	}
}
