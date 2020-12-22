package rpc

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"math/rand"
	"net/http"
	"runtime/debug"
	"sync"
	"time"
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

					debug.PrintStack()
					fmt.Println(response.Error)

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

func MonitorIncomingRequestsMiddleware() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("[mantle/testkit-rpc] incoming request %s %s\n", r.Method, r.RequestURI)

			next.ServeHTTP(w, r)
		})
	}
}

func GenerateTestkitIdentifier(chainId string) string {
	rand.Seed(time.Now().Unix())
	return fmt.Sprintf("%s_%d", chainId, rand.Intn(int(time.Now().Unix())))
}
