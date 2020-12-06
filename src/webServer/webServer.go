package webServer

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"qrcode-server/src/appConst"
	"qrcode-server/src/qrCode"
	"runtime"
	"time"
)

//type QrResponse struct {
//	Id int
//	QrBase64 []byte
//}

type WebServer struct {
	app   *appConst.App
	srv   *http.Server
	qrGen *qrCode.QrRender
	//qrChan				chan *bytes.Buffer
	qrChan chan *qrCode.QrData
}

func toJson(obj interface{}) []byte {
	str, _ := json.Marshal(obj)
	return str
}

type qrF struct {
	Id       int    `json:"id"`
	QrBase64 string `json:"qrBase64"`
}

func (p *WebServer) StartServer() error {
	addr := p.app.GetServerAddr()

	serverMux := mux.NewRouter()
	serverMux.PathPrefix("/ping").HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		runtime.GC()
		fmt.Fprintf(writer, "pong")
	})
	serverMux.PathPrefix("/version").HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		runtime.GC()
		fmt.Fprintf(writer, p.app.GetAppVersion())
	})
	//serverMux.PathPrefix("/qrCode").HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
	//	head := writer.Header()
	//	head.Add("Content-Type", "image/png")
	//	io.Copy(writer, <- p.qrChan)
	//})
	serverMux.PathPrefix("/qrCode").HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		head := writer.Header()
		head.Add("Content-Type", "application/json")
		//head.Add("Content-Type", "image/png")
		qr := <-p.qrChan
		//qrResp := qrCode.QrData{
		//	Id:       qr.Id,
		//	//QrBase64: t.Bytes(),
		//	QrBase64: []byte("!!!-test-!!!"),//[]byte(fmt.Sprintf("data:image/png;base64,%s", base64.StdEncoding.EncodeToString(qr.QrBase64))),
		//}
		qrF := qrF{
			Id:       qr.Id,
			QrBase64: fmt.Sprintf("data:image/png;base64,%s", base64.StdEncoding.EncodeToString(qr.QrBase64)),
			//QrBase64: string(qr.QrBase64),
		}
		jsonResponse := toJson(qrF)
		//jsonResponse := toJson(qrResp)
		//io.Copy(writer, bytes.NewReader(qrResp.QrBase64))
		io.Copy(writer, bytes.NewReader(jsonResponse))
		//(writer, getErrorJson(err.Error()), http.StatusNotFound)
	})

	http.Handle("/", serverMux)

	srv := &http.Server{
		Addr:         addr,
		ReadTimeout:  90 * time.Second,
		WriteTimeout: 90 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	p.srv = srv

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.WithFields(log.Fields{
			"package":  "proxy",
			"function": "StartServer",
			"error":    err,
		}).Info("Server start error")
		return err
	}
	return nil
}

func TimeTrack(start time.Time) {
	timeRun := time.Since(start)
	fmt.Print("run time: ", timeRun, "\n")
}

func (p *WebServer) Shutdown() error {
	if err := p.srv.Shutdown(context.Background()); err != nil {
		log.WithFields(log.Fields{
			"package":  "proxy",
			"function": "Shutdown",
			"error":    err,
		}).Info("Server stop error")
		return err
	}
	return nil
}

func New() (*WebServer, error) {
	app, err := appConst.GetApp("")

	if err != nil {
		log.WithFields(log.Fields{
			"package":  "proxy",
			"function": "New",
			"error":    err,
		}).Info("Error getting configuration instance")
		return nil, err
	}

	qr, err := qrCode.New(&qrCode.QrConf{
		ColorCenter: app.Qr.ColorCenter,
		ColorEdge:   app.Qr.ColorEdge,
		LogoPath:    app.Qr.LogoPath,
		PixRadius:   app.Qr.PixRadius,
		Template:    qrCode.GetQrText(app.Qr.Template, app.Qr.IdPrefix, app.Qr.IdMax),
	})

	qrChan := qrCode.NewBuf(qr, qrCode.QrBufConf{
		IdMax:          app.Qr.IdMax,
		IdPrefix:       app.Qr.IdPrefix,
		Lifetime:       app.Qr.Lifetime,
		SizeBuffer:     app.Qr.SizeBuffer,
		NumberOfTreads: app.Qr.NumberOfThreads,
		Template:       app.Qr.Template,
	})

	return &WebServer{
		app:    app,
		qrGen:  qr,
		qrChan: qrChan,
	}, nil
}
