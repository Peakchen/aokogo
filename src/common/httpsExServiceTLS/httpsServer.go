// add by stefan

package httpsExServiceTLS

import (
	"common/define"
	"fmt"
	"log"
	"net/http"
	"time"
)

func StartHttpsServiceTLS(port int, certfile string, keyfile string, handler define.HandlerFunc) {
	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", port),
		Handler:        handler,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		MaxHeaderBytes: 1 << 16,
	}

	err := s.ListenAndServeTLS(certfile, keyfile)
	if err != nil {
		log.Fatal("https ListenAndServe: ", err)
		return
	}
}
