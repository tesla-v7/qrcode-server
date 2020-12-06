package qrCode

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"math"
	"math/rand"

	"time"
)

type qrBuf struct {
	qr *QrRender
	//qrChan chan *bytes.Buffer
	qrChan chan *QrData

	conf QrBufConf
	//qrHistory map[int]int64
	//qrHistory map[int] *QrData
	//mux sync.Mutex
	//qrTreadHistory []map[int]int64
	rangeMaxId int
}

type QrBufConf struct {
	IdMax          int
	IdPrefix       int
	Lifetime       int64
	SizeBuffer     int
	NumberOfTreads int
	Template       string
}

type QrData struct {
	Id int
	//timeOfUsed int64
	QrBase64 []byte
}

//func (qr *qrBuf) add(qrId int) bool {
//	qr.mux.Lock()
//	timeStart, ok := qr.qrHistory[qrId]
//	timeNow := time.Now().Unix()
//	if ok && timeStart+qr.conf.Lifetime < timeNow {
//		//fmt.Println("dubl", qrId)
//		qr.mux.Unlock()
//		return false
//	}
//	qr.qrHistory[qrId] = timeNow
//	qr.mux.Unlock()
//	return true
//}

func (qb *qrBuf) start() {
	bitNumber := len(fmt.Sprintf("%d", qb.conf.IdMax))
	qb.conf.IdPrefix *= int(math.Pow(10, float64(bitNumber)))
	for i := 0; i < qb.conf.NumberOfTreads; i++ {
		go qb.renderQrInBuf_(i * qb.rangeMaxId)
		//go qb.renderQrInBuf()
	}
}

func GetQrText(template string, idPrefix int, idMax int) string {
	bitNumber := len(fmt.Sprintf("%d", idMax))
	idPrefix *= int(math.Pow(10, float64(bitNumber)))
	return fmt.Sprintf(template, idPrefix+idMax)
}

//func (qb *qrBuf) getGr() {
//
//}

//func (qb *qrBuf) renderQr() (*image.RGBA, int) {
//	//uid, _ := uuid.NewV4()
//	//text := "{\"requestId\": \""+ uid.String() +"\"}"
//	var id int
//	for id = rand.Intn(qb.conf.IdMax); !qb.add(id); id = rand.Intn(qb.conf.IdMax) {
//	}
//
//	id += qb.conf.IdPrefix
//	text := fmt.Sprintf(qb.conf.Template, id)
//	qr, _ := qb.qr.GetQrCodeImage(&text)
//	return qr, id
//}

func (qb *qrBuf) renderQr_(id int) *image.RGBA {
	text := fmt.Sprintf(qb.conf.Template, id)
	qr, _ := qb.qr.GetQrCodeImage(&text)
	return qr
}

//func test(qr *qrBuf) {
//	lineMapSize := int(unsafe.Sizeof(int(0)) + unsafe.Sizeof(int64(0)))
//	for {
//
//		if len(qr.qrChan) >= 8000 {
//			fmt.Println("qrMap", len(qr.qrHistory), len(qr.qrHistory)*lineMapSize)
//		}
//		time.Sleep(time.Second)
//	}
//}

//func (qb *qrBuf)renderQrInBuf(){
//	for{
//		var t bytes.Buffer
//		qr, _ := qb.renderQr()
//		png.Encode(&t, qr)
//		qb.qrChan <- &t
//	}
//}

//func (qb *qrBuf) renderQrInBuf() {
//	for {
//		var t bytes.Buffer
//		qr, id := qb.renderQr()
//		png.Encode(&t, qr)
//		j := &QrData{
//			Id:       id,
//			QrBase64: t.Bytes(),
//			//QrBase64: fmt.Sprintf("data:image/png;base64,%s", base64.StdEncoding.EncodeToString(t.Bytes())),
//		}
//		qb.qrChan <- j
//	}
//}

func (qb *qrBuf) renderQrInBuf_(rangeSegmentStart int) {
	history := make(map[int]int64)
	rangeMaxId := qb.rangeMaxId - 1
	for {
		var t bytes.Buffer
		var id int
		for id = rand.Intn(rangeMaxId); !qb.addId(id, &history); id = rand.Intn(rangeMaxId) {
		}
		id = qb.conf.IdPrefix + rangeSegmentStart + id
		qr := qb.renderQr_(id)
		png.Encode(&t, qr)
		j := &QrData{
			Id:       id,
			QrBase64: t.Bytes(),
		}
		qb.qrChan <- j
	}
}

func (qb *qrBuf) addId(id int, history *map[int]int64) bool {
	timeStart, ok := (*history)[id]
	timeNow := time.Now().Unix()
	if ok && timeStart+qb.conf.Lifetime < timeNow {
		fmt.Print("duble ", ok, id, timeNow-timeStart, "\n")
		return false
	}
	(*history)[id] = timeNow
	return true
}

func NewBuf(qrRender *QrRender, qrConf QrBufConf) chan *QrData {
	qrBuf := &qrBuf{
		qr:     qrRender,
		qrChan: make(chan *QrData, qrConf.SizeBuffer),
		//qrChan:  make(chan *bytes.Buffer, qrConf.SizeBuffer),
		conf: qrConf,
		//qrHistory: make(map[int]int64),
		//mux:       sync.Mutex{},
		rangeMaxId: int(qrConf.IdMax / qrConf.NumberOfTreads),
	}
	qrBuf.start()

	return qrBuf.qrChan
}
