package dnsrbl

import (
	"context"
	"net/http"

	"github.com/hostinger/dnsrbl/common"
	"github.com/labstack/echo/v4"
)

// Handler ...
type Handler struct {
	addressStore AddressStore
}

// NewHandler ...
func NewHandler(addressStore AddressStore) *Handler {
	return &Handler{
		addressStore: addressStore,
	}
}

// BlockHandler ...
// @Summary Block an IP addresss
// @Description Block an IP address
// @Accept  json
// @Produce  json
// @Success 200 {object} common.SimpleResponse
// @Failure 422 {object} common.SimpleResponse
// @Failure 500 {object} common.SimpleResponse
// @Router /block [POST]
func (h *Handler) BlockHandler(c echo.Context) error {
	var address Address
	if err := c.Bind(&address); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, common.SimpleResponse{Status: "Failure", Message: err.Error()})
	}
	if err := h.addressStore.Create(context.Background(), address); err != nil {
		return c.JSON(http.StatusInternalServerError, common.SimpleResponse{Status: "Failure", Message: err.Error()})
	}
	return c.JSON(http.StatusOK, common.SimpleResponse{Status: "Success"})
}

// UnblockHandler ...
// @Summary Unblock an IP addresss
// @Description Unblock an IP address
// @Accept  json
// @Produce  json
// @Param address path string true "IP address to unblock"
// @Success 200 {object} common.SimpleResponse
// @Failure 422 {object} common.SimpleResponse
// @Failure 500 {object} common.SimpleResponse
// @Router /unblock/{address} [DELETE]
func (h *Handler) UnblockHandler(c echo.Context) error {
	address := c.Param("address")
	if err := h.addressStore.Delete(context.Background(), address); err != nil {
		return c.JSON(http.StatusInternalServerError, common.SimpleResponse{Status: "Failure", Message: err.Error()})
	}
	return c.JSON(http.StatusOK, common.SimpleResponse{Status: "Success"})
}

// SearchHandler ...
// @Summary Search for an IP address
// @Description Search for an IP address
// @Accept  json
// @Produce  json
// @Param address path string true "IP address to search for"
// @Success 200 {object} common.SimpleResponse
// @Failure 422 {object} common.SimpleResponse
// @Failure 500 {object} common.SimpleResponse
// @Router /search/{address} [GET]
func (h *Handler) SearchHandler(c echo.Context) error {
	ip := c.Param("address")
	address, err := h.addressStore.Get(context.Background(), ip)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, common.SimpleResponse{Status: "Failure", Message: err.Error()})
	}
	return c.JSON(http.StatusOK, common.SimpleResponse{Status: "Success", Data: address})
}
