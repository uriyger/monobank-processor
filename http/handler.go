package http

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"monobank-processor/model"
	"monobank-processor/processor"
	"net/http"
	"time"

	"go.uber.org/zap"

	"monobank-processor/config"
)

// Handler provides http handler methods
type Handler struct {
	config *config.Config
	logger *zap.Logger
	p      processor.Processor
}

// NewHandler inits http handler
func NewHandler(config *config.Config, logger *zap.Logger, p processor.Processor) *Handler {
	return &Handler{config: config, logger: logger, p: p}
}

func (h *Handler) ProcessStatement(writer http.ResponseWriter, request *http.Request) {
	body, _ := ioutil.ReadAll(request.Body)

	h.logger.Info("request", zap.Any("body", string(body)))

	if err := h.processBody(body); err != nil {
		h.logger.Error("process statement", zap.Error(err))
		writer.Write([]byte("err: " + err.Error()))
	}

	writer.Write([]byte("Ok time: " + time.Now().UTC().Format("2006-01-02 15:04:05.999")))

}

func (h *Handler) processBody(body []byte) error {

	var ms model.MonoStatement
	if err := json.Unmarshal(body, &ms); err != nil {
		return err
	}

	if err := h.p.ProcessStmt(ms.Data.StatementItem); err != nil {
		return fmt.Errorf("process statement failed: %w", err)
	}

	return nil
}

func Ping(writer http.ResponseWriter, _ *http.Request) {
	writer.WriteHeader(http.StatusOK)
}
