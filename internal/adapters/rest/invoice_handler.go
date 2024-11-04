package rest

import (
	"github.com/cs471-buffetpos/buffet-pos-backend/domain/requests"
	"github.com/cs471-buffetpos/buffet-pos-backend/domain/usecases"
	"github.com/cs471-buffetpos/buffet-pos-backend/utils"
	"github.com/gofiber/fiber/v2"
)

type InvoiceHandler interface {
	GetAllUnpaidInvoices(c *fiber.Ctx) error
	SetInvoicePaid(c *fiber.Ctx) error
	CancelInvoice(c *fiber.Ctx) error
}

type invoiceHandler struct {
	service usecases.InvoiceUseCase
}

func NewInvoiceHandler(service usecases.InvoiceUseCase) InvoiceHandler {
	return &invoiceHandler{
		service: service,
	}
}

// Get All Unpaid Invoices
// @Summary Get All Unpaid
// @Description Get all unpaid invoices.
// @Tags Manage
// @Accept json
// @Produce json
// @Success 200 {array} responses.InvoiceDetail
// @Router /manage/invoices/unpaid [get]
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
func (i *invoiceHandler) GetAllUnpaidInvoices(c *fiber.Ctx) error {
	invoices, err := i.service.GetAllUnpaidInvoices(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(invoices)
}

// Set Invoice Paid
// @Summary Set Invoice Paid
// @Description Set invoice as paid.
// @Tags Manage
// @Accept json
// @Produce json
// @Param request body requests.UpdateInvoiceStatusRequest true "Update Invoice Status Request"
// @Success 200 {object} responses.SuccessResponse
// @Router /manage/invoices/set-paid [put]
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
func (i *invoiceHandler) SetInvoicePaid(c *fiber.Ctx) error {
	var req *requests.UpdateInvoiceStatusRequest
	if err := c.BodyParser(&req); err != nil {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if err := utils.ValidateStruct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}
	err := i.service.SetPaidInvoice(c.Context(), req.InvoiceID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"message": "Invoice paid successfully",
	})
}

// Cancel Invoice
// @Summary Cancel Invoice
// @Description Cancel invoice.
// @Tags Manage
// @Accept json
// @Produce json
// @Param request body requests.UpdateInvoiceStatusRequest true "Update Invoice Status Request"
// @Success 200 {object} responses.SuccessResponse
// @Router /manage/invoices/cancel [delete]
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
func (i *invoiceHandler) CancelInvoice(c *fiber.Ctx) error {
	var req *requests.UpdateInvoiceStatusRequest
	err := i.service.DeleteInvoice(c.Context(), req.InvoiceID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"message": "Invoice cancelled successfully",
	})
}