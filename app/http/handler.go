package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"monobank-processor/config"
	"net/http"
	"os"
	"time"
)

type Handler struct {
	config *config.Config
}

func NewHandler(config *config.Config) *Handler {
	return &Handler{config: config}
}

func (h *Handler) ProcessStatement(writer http.ResponseWriter, request *http.Request) {
	body, _ := ioutil.ReadAll(request.Body)
	for key, values := range request.Header {
		fmt.Println(time.Now().UTC().Format("2006-01-02 15:04:05.999"), key, fmt.Sprintf("%v", values))
	}

	fmt.Println(time.Now().UTC().Format("2006-01-02 15:04:05.999"), "-", "Request Host:", request.Host)
	fmt.Println(time.Now().UTC().Format("2006-01-02 15:04:05.999"), "-", "Request URI:", request.RequestURI)
	fmt.Println(time.Now().UTC().Format("2006-01-02 15:04:05.999"), "-", "Request Remote Addr:", request.RemoteAddr)
	fmt.Println(time.Now().UTC().Format("2006-01-02 15:04:05.999"), "-", "Body:", string(body))

	if err := processBody(body); err != nil {
		fmt.Println(err)
		writer.Write([]byte("err: " + err.Error()))
	}

	writer.Write([]byte("Ok time: " + time.Now().UTC().Format("2006-01-02 15:04:05.999")))

}

func processBody(body []byte) error {
	var si struct {
		Type string
		Data struct {
			StatementItem StatementItem `json:"statementItem"`
			Account       string        `json:"account"`
		}
	}
	if err := json.Unmarshal(body, &si); err != nil {
		return err
	}

	chatID := os.Getenv("CHAT_ID")

	mesW := struct {
		ChatID string `json:"chat_id"`
		Text   string `json:"text"`
	}{chatID, si.Data.StatementItem.ComposeMessage()}

	p, _ := json.Marshal(mesW)
	fmt.Println(string(p))
	var resp *http.Response
	var err error
	if resp, err = http.Post(os.Getenv("SEND_MESSAGE_API"), "application/json", bytes.NewReader(p)); err != nil {
		return fmt.Errorf("GET: %s Error:%s", os.Getenv("SEND_MESSAGE_API"), err.Error())
	}

	fmt.Println("Response status: ", resp.StatusCode)
	return nil
}

type StatementItem struct {
	ID              string `json:"id"`
	Time            int64  `json:"time"`
	Description     string `json:"description"`
	Mcc             int    `json:"mcc"`
	Comment         string `json:"comment"`
	Hold            bool   `json:"hold"`
	Amount          int    `json:"amount"`
	OperationAmount int    `json:"operationAmount"`
	CurrencyCode    int    `json:"currencyCode"`
	CommissionRate  int    `json:"commissionRate"`
	CashbackAmount  int    `json:"cashbackAmount"`
	Balance         int    `json:"balance"`
}

func (si StatementItem) ComposeMessage() string {
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

func Ping(writer http.ResponseWriter, _ *http.Request) {
	writer.WriteHeader(http.StatusOK)
}
