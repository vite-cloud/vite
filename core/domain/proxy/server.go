package proxy

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/vite-cloud/go-zoup"
	"github.com/vite-cloud/grace"
	"github.com/vite-cloud/vite/core/domain/config"
	"github.com/vite-cloud/vite/core/domain/datadir"
	"github.com/vite-cloud/vite/core/domain/deployment"
	"github.com/vite-cloud/vite/core/domain/log"
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
		writer: &zoup.CompositeWriter{
			Writers: []zoup.Writer{
				&zoup.FileWriter{File: logFile},
				&zoup.FileWriter{File: stdout},
			},
		},
	}

	conf, err := config.Get(deployment.Locator)
	if err != nil {
		return nil, err
	}

	router := &Router{deployment: deployment, logger: l, API: NewAPI(), config: conf}

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
	httpServer := &http.Server{
		Addr:           ":" + HTTP,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			p.CertManager.HTTPHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "GET" && r.Method != "HEAD" {
					http.Error(w, "Use HTTPS", http.StatusBadRequest)
					return
				}

				target := "https://" + replacePort(r.Host, HTTPS) + r.URL.RequestURI()
				http.Redirect(w, r, target, http.StatusFound)
			})).ServeHTTP(w, r)

			p.Logger.LogR(r, zoup.DebugLevel, "redirect to https")
		}),
	}

	httpsServer := &http.Server{
		Addr:           ":" + HTTPS,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
		Handler:        http.HandlerFunc(p.Router.Proxy),
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS13,
			GetCertificate: func(info *tls.ClientHelloInfo) (*tls.Certificate, error) {
				if unsecure {
					return GetAutoCert()
				}
				return p.CertManager.GetCertificate(info)
			},
		},
	}

	finisher := grace.New(
		grace.WithServer("http", httpServer, time.Second*10),
		grace.WithServer("https", httpsServer, time.Second*10),
	)

	go p.startServer(httpServer)
	go p.startServer(httpsServer)

	finisher.Wait()
}

type Logger struct {
	writer zoup.Writer
}

func (l Logger) Log(level zoup.Level, message string, fields zoup.Fields) {
	err := l.writer.Write(level, message, fields)
	if err != nil {
		panic(err)
	}
}

func (l *Logger) LogR(r *http.Request, level zoup.Level, message string) {
	l.Log(level, message, zoup.Fields{
		"host":   r.Host,
		"method": r.Method,
		"path":   r.URL.Path,
	})
}
func replacePort(url string, newPort string) string {
	host, _, err := net.SplitHostPort(url)
	if err != nil {
		return url
	}
	return net.JoinHostPort(host, newPort)
}

func (p *Proxy) startServer(server *http.Server) {
	var err error
	if server.TLSConfig == nil {
		err = server.ListenAndServe()
	} else {
		err = server.ListenAndServeTLS("", "")
	}

	if err != nil {
		if err == http.ErrServerClosed {
			p.Logger.Log(zoup.InfoLevel, "server shutdown", zoup.Fields{
				"port": server.Addr,
			})
			return
		}

		p.Logger.Log(zoup.ErrorLevel, "https server error", zoup.Fields{
			"port": err,
		})
		os.Exit(1)
	}
}
