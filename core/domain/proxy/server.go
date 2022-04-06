package proxy

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/vite-cloud/vite/core/domain/datadir"
	"github.com/vite-cloud/vite/core/domain/deployment"
	"github.com/vite-cloud/vite/core/domain/log"
	"github.com/vite-cloud/vite/core/domain/plane"
	"golang.org/x/crypto/acme/autocert"
	"io"
	"net"
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
	Logger      *Logger
}

func New(stdout io.Writer, deployment *deployment.Deployment) (*Proxy, error) {
	dir, err := Store.Dir()
	if err != nil {
		return nil, err
	}

	logFile, err := log.Store.Open("proxy.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	l := &Logger{
		writer: &log.CompositeWriter{
			Writers: []log.Writer{
				&log.FileWriter{File: logFile},
				&log.FileWriter{File: stdout},
			},
		},
	}
	router := &Router{deployment: deployment, logger: l, API: plane.New()}

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
		Logger: l,
	}, nil
}

func (p *Proxy) Run(HTTP string, HTTPS string, unsecure bool) {
	handlerHTTP := newServer(HTTP, func(w http.ResponseWriter, r *http.Request) {
		p.CertManager.HTTPHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "GET" && r.Method != "HEAD" {
				http.Error(w, "Use HTTPS", http.StatusBadRequest)
				return
			}

			target := "https://" + replacePort(r.Host, HTTPS) + r.URL.RequestURI()
			http.Redirect(w, r, target, http.StatusFound)
		})).ServeHTTP(w, r)

		p.Logger.LogR(r, log.DebugLevel, "redirect to https")
	})

	handler := newServer(HTTPS, p.Router.Proxy)
	handler.TLSConfig = &tls.Config{
		MinVersion: tls.VersionTLS13,
		GetCertificate: func(info *tls.ClientHelloInfo) (*tls.Certificate, error) {
			if unsecure {
				return GetAutoCert()
			}
			return p.CertManager.GetCertificate(info)
		},
	}

	finisher := Finisher{
		Keepers: []*Keeper{
			{"http", handlerHTTP, time.Second * 10},
			{"https", handler, time.Second * 10},
		},
		logger: p.Logger,
	}

	go func() {
		err := handlerHTTP.ListenAndServe()
		if err != nil {
			if err == http.ErrServerClosed {
				p.Logger.Log(log.InfoLevel, "server shutdown", log.Fields{
					"name": "http",
				})
				return
			}
			p.Logger.Log(log.ErrorLevel, "http server error", log.Fields{
				"error": err,
			})
			os.Exit(1)
		}
	}()

	go func() {
		err := handler.ListenAndServeTLS("", "")
		if err != nil {
			if err == http.ErrServerClosed {
				p.Logger.Log(log.InfoLevel, "server shutdown", log.Fields{
					"name": "https",
				})
				return
			}

			p.Logger.Log(log.ErrorLevel, "https server error", log.Fields{
				"error": err,
			})
			os.Exit(1)
		}
	}()

	finisher.Wait()
}

type Logger struct {
	writer log.Writer
}

func (l Logger) Log(level log.Level, message string, fields log.Fields) {
	err := l.writer.Write(level, message, fields)
	if err != nil {
		panic(err)
	}
}

func (l *Logger) LogR(r *http.Request, level log.Level, message string) {
	l.Log(level, message, log.Fields{
		"host":   r.Host,
		"method": r.Method,
		"path":   r.URL.Path,
	})
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

func replacePort(url string, newPort string) string {
	host, _, err := net.SplitHostPort(url)
	if err != nil {
		return url
	}
	return net.JoinHostPort(host, newPort)
}
