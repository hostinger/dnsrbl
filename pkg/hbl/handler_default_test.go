package hbl

import (
	"context"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

var (
	h Handler
	s Service
	r Repository
)

func Test_handler_HandleAddressesGetAll(t *testing.T) {
	e := echo.New()

	r.CreateAddress(context.Background(), &Address{IP: "127.0.0.1", Author: "Test", Comment: "Test", Action: "Block"})
	r.CreateAddress(context.Background(), &Address{IP: "127.0.0.2", Author: "Test", Comment: "Test", Action: "Block"})
	r.CreateAddress(context.Background(), &Address{IP: "127.0.0.3", Author: "Test", Comment: "Test", Action: "Block"})

	req := httptest.NewRequest("GET", "/api/v1/addresses", nil)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	ctx := e.NewContext(req, rec)
	ctx.SetPath("/api/v1/addresses")

	if assert.NoError(t, h.HandleAddressesGetAll(ctx)) {
		var addresses []Address
		if err := json.Unmarshal(rec.Body.Bytes(), &addresses); err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, 3, len(addresses))
		assert.Equal(t, 200, rec.Code)
	}
}

func Test_handler_HandleHealth(t *testing.T) {
	e := echo.New()

	req := httptest.NewRequest("GET", "/health", nil)
	rec := httptest.NewRecorder()

	ctx := e.NewContext(req, rec)
	ctx.SetPath("/health")

	if assert.NoError(t, h.HandleHealth(ctx)) {
		assert.Equal(t, "OK", rec.Body.String())
		assert.Equal(t, 200, rec.Code)
	}
}

func Test_handler_HandleVersion(t *testing.T) {
	e := echo.New()

	req := httptest.NewRequest("GET", "/version", nil)
	rec := httptest.NewRecorder()

	ctx := e.NewContext(req, rec)
	ctx.SetPath("/version")

	if assert.NoError(t, h.HandleVersion(ctx)) {
		assert.Equal(t, "1.0.0", rec.Body.String())
		assert.Equal(t, 200, rec.Code)
	}
}

func init() {
	r = NewMockRepository()

	s = &service{
		repository: r,
	}

	h = &handler{
		service: s,
	}
}
