package main

import (
	"context"
	"crypto/tls"
	"github.com/gookit/color"
	"github.com/gorilla/securecookie"
	"github.com/xi2/httpgzip"
	"golang.org/x/crypto/acme/autocert"
	"golang.org/x/net/http2"
	"goyoubbs/cronjob"
	"goyoubbs/goji"
	"goyoubbs/goji/pat"
	"goyoubbs/router"
	"goyoubbs/system"
	"goyoubbs/util"
	"goyoubbs/youdb"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	app := &system.Application{}
	c := system.LoadConfig()
	mcf := &system.MainConf{}
	err := c.UnmarshalKey("Main", mcf)
	if err != nil {

		return
	}
	// check domain
	if strings.HasPrefix(mcf.Domain, "http") {
		dm, err := url.Parse(mcf.Domain)
		if err != nil {
			log.Fatal("domain fmt err", err)
		}
		mcf.Domain = dm.Host
	} else {
		mcf.Domain = strings.Trim(mcf.Domain, "/")
	}

	scf := &system.SiteConf{}
	err2 := c.UnmarshalKey("Site", scf)
	if err2 != nil {
		return
	}
	scf.GoVersion = runtime.Version()
	fMd5, _ := util.HashFileMD5(os.Args[0])
	scf.MD5Sums = fMd5
	scf.MainDomain = strings.Trim(scf.MainDomain, "/")
	log.Println("MainDomain:", scf.MainDomain)
	if scf.TimeZone < -12 || scf.TimeZone > 12 {
		scf.TimeZone = 0
	}
	if scf.UploadMaxSize < 1 {
		scf.UploadMaxSize = 1
	}
	scf.UploadMaxSizeByte = int64(scf.UploadMaxSize) << 20

	app.Cf = &system.AppConf{mcf, scf}
	color.Redln("打开数据库")
	db, err := youdb.Open(mcf.Youdb)
	if err != nil {
		log.Fatalf("Connect Error: %v", err)
	}
	app.Db = db

	defer app.Db.Close()
	// set main node
	db.Hset("keyValue", []byte("main_category"), []byte(scf.MainNodeIds))

	var hashKey []byte
	var blockKey []byte
	if scf.ResetCookieKey {
		hashKey = securecookie.GenerateRandomKey(64)
		blockKey = securecookie.GenerateRandomKey(32)
		_ = db.Hmset("keyValue", []byte("hashKey"), hashKey, []byte("blockKey"), blockKey)
	} else {
		hashKey = append(hashKey, db.Hget("keyValue", []byte("hashKey")).Bytes()...)
		blockKey = append(blockKey, db.Hget("keyValue", []byte("blockKey")).Bytes()...)
		if len(hashKey) == 0 {
			hashKey = securecookie.GenerateRandomKey(64)
			blockKey = securecookie.GenerateRandomKey(32)
			_ = db.Hmset("keyValue", []byte("hashKey"), hashKey, []byte("blockKey"), blockKey)
		}
	}

	app.Sc = securecookie.New(hashKey, blockKey)
	// cron job
	cr := cronjob.BaseHandler{App: app}
	go cr.MainCronJob()

	root := goji.NewMux()

	// static file server
	staticPath := mcf.PubDir
	if len(staticPath) == 0 {
		staticPath = "static"
	}

	root.Handle(pat.New("/.well-known/acme-challenge/*"),
		http.StripPrefix("/.well-known/acme-challenge/", http.FileServer(http.Dir(staticPath))))
	root.Handle(pat.New("/static/*"),
		http.StripPrefix("/static/", http.FileServer(http.Dir(staticPath))))

	root.Handle(pat.New("/*"), router.NewRouter(app))

	// normal http
	// http.ListenAndServe(listenAddr, root)

	var srv *http.Server

	if mcf.HttpsOn {
		// https
		log.Println("Register sll for domain:", mcf.Domain)
		log.Println("TLSCrtFile : ", mcf.TLSCrtFile)
		log.Println("TLSKeyFile : ", mcf.TLSKeyFile)

		root.Use(stlAge)

		tlsCf := &tls.Config{
			NextProtos: []string{http2.NextProtoTLS, "http/1.1"},
		}

		if mcf.Domain != "" && mcf.TLSCrtFile == "" && mcf.TLSKeyFile == "" {

			domains := strings.Split(mcf.Domain, ",")
			certManager := autocert.Manager{
				Prompt:     autocert.AcceptTOS,
				HostPolicy: autocert.HostWhitelist(domains...),
				Cache:      autocert.DirCache("certs"),
				Email:      scf.AdminEmail,
			}
			tlsCf.GetCertificate = certManager.GetCertificate
			//tlsCf.ServerName = domains[0]

			go func() {
				// 必须是 80 端口
				log.Fatal(http.ListenAndServe(":http", certManager.HTTPHandler(nil)))
			}()

		} else {
			// rewrite
			go func() {
				if err := http.ListenAndServe(":"+strconv.Itoa(mcf.HttpPort), http.HandlerFunc(redirectHandler)); err != nil {
					log.Println("Http2https server failed ", err)
				}
			}()
		}

		srv = &http.Server{
			Addr:           ":" + strconv.Itoa(mcf.HttpsPort),
			Handler:        httpgzip.NewHandler(root, nil),
			TLSConfig:      tlsCf,
			MaxHeaderBytes: int(app.Cf.Site.UploadMaxSizeByte),
			ReadTimeout:    5 * time.Second,
			WriteTimeout:   10 * time.Second,
			BaseContext:    func(_ net.Listener) context.Context { return ctx },
		}
		srv.RegisterOnShutdown(cancel)

		go func() {
			// 如何获取 TLSCrtFile、TLSKeyFile 文件参见 https://www.youbbs.org/topic/2169
			if err := srv.ListenAndServeTLS(mcf.TLSCrtFile, mcf.TLSKeyFile); err != http.ErrServerClosed {
				// it is fine to use Fatal here because it is not main gorutine
				log.Fatalf("HTTPS server ListenAndServe: %v", err)
			}
		}()

		log.Println("Web server Listen port", mcf.HttpsPort)
		log.Println("Web server URL", "https://"+mcf.Domain)

	} else {
		// http
		srv = &http.Server{
			Addr:         ":" + strconv.Itoa(mcf.HttpPort),
			Handler:      root,
			ReadTimeout:  6 * time.Second,
			WriteTimeout: 10 * time.Second,
			BaseContext:  func(_ net.Listener) context.Context { return ctx },
		}
		srv.RegisterOnShutdown(cancel)

		go func() {
			if err := srv.ListenAndServe(); err != http.ErrServerClosed {
				// it is fine to use Fatal here because it is not main gorutine
				log.Fatalf("HTTP server ListenAndServe: %v", err)
			}
		}()

		log.Println("Web server Listen port", mcf.HttpPort)
	}

	// graceful stop
	// subscribe to SIGINT signals
	signalChan := make(chan os.Signal, 1)

	signal.Notify(
		signalChan,
		syscall.SIGHUP,  // kill -SIGHUP XXXX
		syscall.SIGINT,  // kill -SIGINT XXXX or Ctrl+c
		syscall.SIGTERM, // kill -SIGTERM XXXX
		syscall.SIGQUIT, // kill -SIGQUIT XXXX
		//syscall.SIGUSR2,
	)

	<-signalChan // wait for SIGINT
	log.Print("os.Interrupt - shutting down...\n")

	go func() {
		<-signalChan
		log.Fatal("os.Kill - terminating...\n")
	}()

	// 等待 30秒 ，等请求结束，同时不允许新的请求
	gracefulCtx, cancelShutdown := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelShutdown()

	srv.SetKeepAlivesEnabled(false)
	if err := srv.Shutdown(gracefulCtx); err != nil {
		log.Printf("shutdown error: %v\n", err)
		defer os.Exit(1)
		return
	} else {
		app.Db.Close() // !important 留意上下文位置
		log.Printf("gracefully stopped\n")
	}

	// manually cancel context if not using httpServer.RegisterOnShutdown(cancel)
	// cancel()

	defer os.Exit(0)
	return
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	target := "https://" + r.Host + r.URL.Path
	if len(r.URL.RawQuery) > 0 {
		target += "?" + r.URL.RawQuery
	}
	// consider HSTS if your clients are browsers
	w.Header().Set("Connection", "close")
	http.Redirect(w, r, target, 302)
}

func stlAge(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// add max-age to get A+
		w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
