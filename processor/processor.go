package processor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"monobank-processor/config"
	"monobank-processor/messenger"
	"monobank-processor/model"
	"net/http"
	"os"
	"time"
)

type Processor interface {
	ProcessStmt(si model.StatementItem) error
}

type processor struct {
	config *config.Config
}

// NewProcessor inits processor
func NewProcessor(c *config.Config) Processor {
	return processor{config: c}
}

func (p processor) ProcessStmt(si model.StatementItem) error {

	mesW := messenger.TeleMsg{
		ChatID: p.config.ChatID,
		Text:   ComposeMessage(si),
	}

	ms, _ := json.Marshal(mesW)

	var resp *http.Response
	var err error
	if resp, err = http.Post(os.Getenv("SEND_MESSAGE_API"), "application/json", bytes.NewReader(ms)); err != nil {
		return fmt.Errorf("GET: %s Error:%s", os.Getenv("SEND_MESSAGE_API"), err.Error())
	}

	fmt.Println("Response status: ", resp.StatusCode)

	return nil
}

func ComposeMessage(si model.StatementItem) string {
	buf := bytes.Buffer{}

	kiev, _ := time.LoadLocation("Europe/Kiev")

	cur := "UAH"
	buf.WriteString(time.Unix(si.Time, 0).In(kiev).Format("02/01/2006 15:04") + "\n")
	buf.WriteString(si.Description + "\n")
	if si.Comment != "" {
		buf.WriteString(si.Comment + "\n")
	}
	buf.WriteString(intToPrice(si.Amount) + cur + "\n")
	buf.WriteString("Остаток:" + intToPrice(si.Balance) + cur + "\n")

	return buf.String()
}

func intToPrice(v int) string {
	return fmt.Sprintf("%v", float32(v)/float32(100))
}
