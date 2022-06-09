package app

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"monobank-processor/config"
	phttp "monobank-processor/http"
	"monobank-processor/processor"
)

type App struct {
	logger  *zap.Logger
	config  *config.Config
	router  *mux.Router
	handler *phttp.Handler
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
			a.logger.Error(err.Error())
		} else {
			fmt.Println(err)
		}
		return 1
	}

	a.logger.Sugar().Info("Monobank processor started on %d", a.config.HTTPPort)
	if err := http.ListenAndServe(":"+strconv.Itoa(a.config.HTTPPort), a.router); err != nil {
		fmt.Println(time.Now().UTC().Format("2006-01-02 15:04:05.999"), "-", "Error:", err.Error())
		return 1
	}
	return 0
}

func (a *App) initLogger() error {
	//conf := zap.NewDevelopmentConfig()
	conf := zap.NewProductionConfig()
	conf.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	conf.Encoding = "json"
	conf.Sampling = nil
	conf.DisableStacktrace = true
	conf.DisableCaller = false
	//conf.Level.SetLevel(zapcore.InfoLevel)

	logger, err := conf.Build()
	if err != nil {
		return fmt.Errorf("zap logger: %w", err)
	}
	a.logger = logger

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
	p := processor.NewProcessor(a.config)
	a.handler = phttp.NewHandler(a.config, a.logger, p)
	return nil
}

func (a *App) initRouter() error {
	if a.handler == nil {
		return errors.New("router initialization: handler is not initialized")
	}
	a.router = phttp.NewRouter(a.handler)

	return nil
}
