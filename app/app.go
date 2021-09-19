package app

import (
	"fmt"
	"monobank-processor/config"
	"net/http"
	"strconv"
	"time"

	"go.uber.org/zap/zapcore"

	"github.com/gorilla/mux"
	"go.uber.org/zap"

	apphttp "monobank-processor/app/http"
)

type App struct {
	logger  *zap.SugaredLogger
	config  *config.Config
	router  *mux.Router
	handler *apphttp.Handler
}

func (a *App) Init() error {

	for _, init := range []func() error{
		a.initLogger,
		a.initConfig,
		a.initHTTPHandler,
		a.initRouter,
	} {
		if err := init(); err != nil {
			return fmt.Errorf("initialization: %w", err)
		}
	}

	return nil
}

func (a *App) Run() int {
	if err := a.Init(); err != nil {
		if a.logger != nil {
			a.logger.Error(err)
		} else {
			fmt.Println(err)
		}
		return 1
	}

	a.logger.Infof("Monobank processor started on %d", a.config.HTTPPort)
	if err := http.ListenAndServe(":"+strconv.Itoa(a.config.HTTPPort), a.router); err != nil {
		fmt.Println(time.Now().UTC().Format("2006-01-02 15:04:05.999"), "-", "Error:", err.Error())
		return 1
	}
	return 0
}

func (a *App) initLogger() error {
	conf := zap.NewProductionConfig()
	conf.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339Nano)
	conf.DisableCaller = true

	logger, err := conf.Build()
	if err != nil {
		return fmt.Errorf("zap logger: %w", err)
	}
	a.logger = logger.Sugar()

	a.logger.Debug("Debug on Prod")

	//conf.Development = true
	conf.Level.SetLevel(zapcore.DebugLevel)

	//logger, _ = conf.Build()
	//a.logger = logger.Sugar()

	a.logger.Debug("Debug on Dev")

	return nil
}

func (a *App) initConfig() error {
	var err error
	a.config, err = config.NewConfig()
	if err != nil {
		return fmt.Errorf("initConfig: %w", err)
	}

	return nil
}

func (a *App) initHTTPHandler() error {
	a.handler = apphttp.NewHandler(a.config)
	return nil
}
