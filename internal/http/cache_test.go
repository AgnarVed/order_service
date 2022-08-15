package http

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http/httptest"
	"strings"
	"testing"
	"tests2/internal/models"
)

func TestHandler_getOrderFromCacheByID(t *testing.T) {
	app := setup()

	tests := []struct {
		description      string
		route            string
		expectedCode     int
		body             io.Reader
		expectedResponse models.Response
	}{
		// Get 400
		{
			description:      "get 400",
			route:            "/cache/asdf",
			expectedCode:     400,
			body:             nil,
			expectedResponse: models.Response{},
		},
		// Get 200
		{
			description:      "get 200",
			route:            "/cache/b563feb7b2b84b6test5",
			expectedCode:     200,
			body:             nil,
			expectedResponse: models.Response{},
		},
	}

	for _, test := range tests {
		req := httptest.NewRequest("GET", test.route, test.body)
		req.Header.Add("accept", "*/*")
		req.Header.Add("Content-Type", "application/json")

		resp, _ := app.Test(req, -1)
		b, _ := io.ReadAll(resp.Body)

		var got models.Response
		_ = json.Unmarshal(b, &got)
	}
}

func TestHandler_putOrderInCache(t *testing.T) {
	app := setup()

	tests := []struct {
		description      string
		route            string
		expectedCode     int
		body             io.Reader
		expectedResponse models.Response
	}{
		// Get 405
		{
			description:  "get 400 nil body",
			route:        "/order",
			expectedCode: 405,
			body:         nil,
		},
		// Get 400
		{
			description:  "get 400 empty json",
			route:        "/order/post",
			expectedCode: 400,
			body:         strings.NewReader("{}"),
		},
		// Get 200
		{
			description:  "get 200",
			route:        "/order/post",
			expectedCode: 200,
			body: strings.NewReader(`{
  "order_uid": "b563feb7b2b84b6test",
  "track_number": "WBILMTESTTRACK",
  "entry": "WBIL",
  "delivery": {
    "name": "Test Testov",
    "phone": "+9720000000",
    "zip": "2639809",
    "city": "Kiryat Mozkin",
    "address": "Ploshad Mira 15",
    "region": "Kraiot",
    "email": "test@gmail.com"
  },
  "payment": {
    "transaction": "b563feb7b2b84b6test",
    "request_id": "req_id",
    "currency": "USD",
    "provider": "wbpay",
    "amount": 1817,
    "payment_dt": 1637907727,
    "bank": "alpha",
    "delivery_cost": 1500,
    "goods_total": 317,
    "custom_fee": 129
  },
  "items": [
    {
      "chrt_id": 9934930,
      "track_number": "WBILMTESTTRACK",
      "price": 453,
      "rid": "ab4219087a764ae0btest",
      "name": "Mascaras",
      "sale": 30,
      "size": "0",
      "total_price": 317,
      "nm_id": 2389212,
      "brand": "Vivienne Sabo",
      "status": 202
    }
  ],
  "locale": "en",
  "internal_signature": "int signature",
  "customer_id": "test",
  "delivery_service": "meest",
  "shardkey": "9",
  "sm_id": 99,
  "date_created": "2021-11-26T06:22:19Z",
  "oof_shard": "1"
}`),
		},
	}
	for _, test := range tests {
		req := httptest.NewRequest("POST", test.route, test.body)
		req.Header.Add("accept", "*/*")
		req.Header.Add("Content-Type", "application/json")

		resp, _ := app.Test(req, -1)
		b, _ := io.ReadAll(resp.Body)
		assert.Equalf(t, test.expectedCode, resp.StatusCode, test.description+":"+string(b))
	}

}
