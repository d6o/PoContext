package middleware

import (
	"net/http"

	"github.com/disiqueira/PoContext/context"
)

func TimeoutRecover() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {

			done := make(chan bool)

			go func() {
				next.ServeHTTP(w, r)
				select {
				case <-r.Context().Done():

				default:
					done <- true
				}

				close(done)
			}()

			log := context.Logger(r.Context())
			select {
			case <-done:
				return
			case <-r.Context().Done():
				log.Warn("context timeout elapsed")
				return
			}

		}
		return http.HandlerFunc(fn)
	}
}
