package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/acme/autocert"
)

var imgs []string
var showid int
var gImgID int

var (
	optDomain  string
	optSecure  bool
	optPort    string
	optImgsDir string
	optSuffix  string
)

var epoch = time.Unix(0, 0).Format(time.RFC1123)
var rd = rand.New(rand.NewSource(time.Now().UnixNano()))

// Taken from https://github.com/mytrile/nocache
var noCacheHeaders = map[string]string{
	"Expires":         epoch,
	"Cache-Control":   "no-cache, private, max-age=0",
	"Pragma":          "no-cache",
	"X-Accel-Expires": "0",
}

var etagHeaders = []string{
	"ETag",
	"If-Modified-Since",
	"If-Match",
	"If-None-Match",
	"If-Range",
	"If-Unmodified-Since",
}

func NoCache(w *http.ResponseWriter, r *http.Request) {
	return
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	// Delete any ETag headers that may have been set
	for _, v := range etagHeaders {
		if r.Header.Get(v) != "" {
			r.Header.Del(v)
		}
	}

	// Set our NoCache headers
	for k, v := range noCacheHeaders {
		w.Header().Set(k, v)
	}

	http.ServeFile(w, r, "form.html")
}

func echoHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(`echo ...`)
	w.Write([]byte("Hello world"))
}

func imgHandler(w http.ResponseWriter, r *http.Request) {

	//imgid := rd.Intn(len(imgs) - 1)

	imgid := gImgID % len(imgs)
	imgDirPath := optImgsDir + `\` + imgs[imgid]
	log.Println(`imgDirPath `, imgDirPath)

	gImgID = gImgID + 1

	// Delete any ETag headers that may have been set
	for _, v := range etagHeaders {
		if r.Header.Get(v) != "" {
			r.Header.Del(v)
		}
	}
	// Set our NoCache headers
	for k, v := range noCacheHeaders {
		w.Header().Set(k, v)
	}

	http.ServeFile(w, r, imgDirPath)
}

func main() {
	flag.BoolVar(&optSecure, "s", false, "run in secure mode (https)")
	flag.StringVar(&optDomain, "n", "z2018168.com", "cookie domain")
	flag.StringVar(&optImgsDir, "d", "D:\\imgs", "-d : 图片的路径.")
	flag.StringVar(&optPort, "p", "80", "listen port")
	flag.StringVar(&optSuffix, "f", "jpg;", "sq img suffix name")
	flag.Parse()

	gImgID = 0
	imgs = LoadImgs(optImgsDir, optSuffix)
	if len(imgs) == 0 {
		fmt.Println(`指定目录下没有二维码文件，退出`)
		return
	}

	for i, img := range imgs {
		fmt.Println(i, img)
	}

	rt := mux.NewRouter()
	rt.HandleFunc("/", helloHandler)
	rt.HandleFunc("/echo", echoHandler)
	rt.HandleFunc("/img", imgHandler)

	if optSecure {
		m := &autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(optDomain, "www."+optDomain),
			Cache:      autocert.DirCache("certcache"),
		}

		/*	go func() {
				log.Fatal(http.ListenAndServe(":http",
					m.HTTPHandler(nil)))
			}()
		*/
		s := &http.Server{
			Addr:           ":443",
			Handler:        rt,
			ReadTimeout:    20 * time.Second,
			WriteTimeout:   20 * time.Second,
			MaxHeaderBytes: 1 << 20,
			TLSConfig: &tls.Config{
				GetCertificate: m.GetCertificate,
			},
		}
		log.Println("Img server started at 443(https) domain", optDomain, " ...")
		err := s.ListenAndServeTLS("", "")
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		s := &http.Server{
			Addr:           ":" + optPort,
			Handler:        rt,
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
			MaxHeaderBytes: 1 << 20,
		}
		log.Println("Img server started at", optPort, ", domain", optDomain, " ...")
		err := s.ListenAndServe()
		if err != nil {
			log.Fatalln(err)
		}
	}
}
