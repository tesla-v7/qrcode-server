package appConst

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

type App struct {
	scheme string
	addr   string

	appVersion string
	lifetime   int64

	Qr qrConf
}

type qrConf struct {
	Template        string
	IdMax           int
	IdPrefix        int
	LogoPath        string
	ColorCenter     string
	ColorEdge       string
	PixRadius       int
	Lifetime        int64
	SizeBuffer      int
	NumberOfThreads int
}

var app *App

func (app *App) SetAddr(addr string) {
	app.addr = addr
}

func (app *App) GetServerAddr() string {
	return app.addr
}

func (app *App) GetServerFullAddr() string {
	return fmt.Sprintf("%s://%s", app.scheme, app.addr)
}

func (app *App) GetAppVersion() string {
	return app.appVersion
}

func new(appVersion string) (*App, error) {
	filePathFull, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.WithFields(log.Fields{
			"package":  "appConst",
			"function": "new",
			"error":    err,
			//"notifier": viper,
		}).Fatal("Failed to get file location")
		return nil, err
	}

	viper.SetConfigFile(fmt.Sprintf("%s/config.toml", filePathFull))
	viper.SetDefault("qr.template", "{\"requestId\": %d}")
	viper.SetDefault("qr.idMax", 99999999)
	viper.SetDefault("qr.idPrefix", 1)
	viper.SetDefault("qr.logoPath", "logo.png")
	viper.SetDefault("qr.colorCenter", "05ffb8")
	viper.SetDefault("qr.colorEdge", "007aff")
	viper.SetDefault("qr.pixRadius", 3)
	viper.SetDefault("qr.lifetime", 120)
	viper.SetDefault("qr.sizeBuffer", 10)
	viper.SetDefault("qr.numberOfThreads", 1)

	if err := viper.ReadInConfig(); err != nil {
		log.WithFields(log.Fields{
			"package":  "appConst",
			"function": "new",
			"error":    err,
			//"notifier": viper,
		}).Error("Failed to read config file")
		return nil, err
	}

	logoPath := filePathFull + "/" + viper.GetString("qr.logoPath")
	if _, err := os.Stat(logoPath); os.IsNotExist(err) {
		return nil, errors.New("logo file note found: " + logoPath)
	}

	return &App{
		scheme:     "http",
		appVersion: appVersion,
		Qr: qrConf{
			Template:        viper.GetString("qr.template"),
			IdMax:           viper.GetInt("qr.idMax"),
			IdPrefix:        viper.GetInt("qr.idPrefix"),
			LogoPath:        logoPath,
			ColorCenter:     viper.GetString("qr.colorCenter"),
			ColorEdge:       viper.GetString("qr.colorEdge"),
			PixRadius:       viper.GetInt("qr.pixRadius"),
			SizeBuffer:      viper.GetInt("qr.sizeBuffer"),
			NumberOfThreads: viper.GetInt("qr.numberOfThreads"),
		},
	}, nil
}

func GetApp(appVersion string) (*App, error) {
	if app == nil {
		var err error
		if app, err = new(appVersion); err != nil {
			return nil, err
		}
	}
	return app, nil
}
