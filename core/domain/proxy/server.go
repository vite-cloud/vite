package proxy

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/vite-cloud/vite/core/domain/datadir"
	"github.com/vite-cloud/vite/core/domain/deployment"
	"github.com/vite-cloud/vite/core/domain/log"
	"golang.org/x/crypto/acme/autocert"
	"net/http"
	"os"
	"time"
)

const (
	Store = datadir.Store("certs")
)

type Proxy struct {
	Router      *Router
	CertManager *autocert.Manager
}

func New(deployment *deployment.Deployment) (*Proxy, error) {
	dir, err := Store.Dir()
	if err != nil {
		return nil, err
	}

	router := &Router{deployment: deployment}
	return &Proxy{
		Router: router,
		CertManager: &autocert.Manager{
			Prompt: autocert.AcceptTOS,
			HostPolicy: func(ctx context.Context, host string) error {
				ok, err := router.Accepts(host)
				if err != nil {
					return err
				}

				if !ok {
					return fmt.Errorf("%s is not allowed", host)
				}

				return nil
			},
			Cache: autocert.DirCache(dir),
		},
	}, nil
}

func (p *Proxy) Run(HTTP string, HTTPS string) {
	handlerHTTP := newServer(HTTP, func(w http.ResponseWriter, r *http.Request) {
		p.CertManager.HTTPHandler(nil).ServeHTTP(w, r)

		LogR(r, log.DebugLevel, "redirect to https")
	})

	handler := newServer(HTTPS, p.Router.Proxy)
	handler.TLSConfig = &tls.Config{
		MinVersion: tls.VersionTLS13,
		GetCertificate: func(info *tls.ClientHelloInfo) (*tls.Certificate, error) {
			return p.CertManager.GetCertificate(info)
		},
	}

	finisher := Finisher{
		Keepers: []*Keeper{
			{
				Name:    "http",
				Timeout: time.Second * 10,
				Server:  handlerHTTP,
			},
			{
				Name:    "https",
				Timeout: time.Second * 10,
				Server:  handler,
			},
		},
	}

	go func() {
		err := handlerHTTP.ListenAndServe()
		if err != nil {
			if err == http.ErrServerClosed {
				Log(log.InfoLevel, "server shutdown", log.Fields{
					"name": "http",
				})
				return
			}
			Log(log.ErrorLevel, "http server error", log.Fields{
				"error": err,
			})

			// todo: use cli.Out()
			fmt.Println(err)
			os.Exit(1)
		}
	}()

	go func() {
		err := handler.ListenAndServeTLS("", "")
		if err != nil {
			if err == http.ErrServerClosed {
				Log(log.InfoLevel, "server shutdown", log.Fields{
					"name": "https",
				})
				return
			}

			Log(log.ErrorLevel, "https server error", log.Fields{
				"error": err,
			})

			// todo: use cli.Out()
			fmt.Println(err)
			os.Exit(1)
		}
	}()

	finisher.Wait()
}

func newServer(port string, handler http.HandlerFunc) *http.Server {
	return &http.Server{
		Addr:           ":" + port,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
		Handler:        handler,
		// todo: compat with std log library
		//ErrorLog:       GetLogger(),
	}
}
