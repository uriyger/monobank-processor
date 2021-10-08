package http

import (
	"bytes"
	"monobank-processor/config"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
)

func initHandler() http.Handler {

	c := &config.Config{
		ChatID:         "chat_id",
		HTTPPort:       80,
		SendMessageURL: "http://telegram/sendMes",
	}
	h := NewHandler(c)
	os.Setenv("SEND_MESSAGE_API", c.SendMessageURL)
	return NewRouter(h)

}

func initTestReqBody() []byte {
	tReqBody := []byte(`{
        "type": "StatementItem",
        "data": {
            "account": "nqG07yl-YivPP_btoAzgYA",
            "statementItem": {
                "id": "yR0rQlSW2M3r6pI8",
                "time": 1604406086,
                "description": "Uber",
                "mcc": 4121,
                "amount": -9724,
                "operationAmount": -9724,
                "currencyCode": 980,
                "commissionRate": 0,
                "cashbackAmount": 2917,
                "balance": 2403378,
                "hold": true
            }
        }
    }`)

	return tReqBody
}

func TestHandler_ProcessStatement(t *testing.T) {

	h := initHandler()
	server := httptest.NewServer(h)
	host := server.URL
	defer server.Close()

	gock.New("http://telegram").URL("/sendMes").Reply(http.StatusOK)

	resp, err := server.Client().Post(host+"/statement", "application/json", bytes.NewBuffer(initTestReqBody()))

	assert.NotNil(t, resp)
	assert.NoError(t, err)
}
