package rest

import (
	"github.com/cs471-buffetpos/buffet-pos-backend/domain/exceptions"
	"github.com/cs471-buffetpos/buffet-pos-backend/domain/requests"
	"github.com/cs471-buffetpos/buffet-pos-backend/domain/usecases"
	"github.com/cs471-buffetpos/buffet-pos-backend/utils"
	"github.com/gofiber/fiber/v2"
)

type TableHandler interface {
	AddTable(c *fiber.Ctx) error
	FindTableByID(c *fiber.Ctx) error
}

type tableHandler struct {
	service usecases.TableUseCase
}

func NewTableHandler(service usecases.TableUseCase) TableHandler {
	return &tableHandler{
		service: service,
	}
}

// Add Table
// @Summary Add Table
// @Description Add a new table.
// @Tags Manage
// @Accept json
// @Produce json
// @Param request body requests.AddTableRequest true "Add Table request"
// @Success 200 {object} responses.AddTableResponse
// @Failure 400
// @Failure 401
// @Failure 500
// @Router /manage/tables [post]
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
func (t *tableHandler) AddTable(c *fiber.Ctx) error {
	var req *requests.AddTableRequest
	if err := c.BodyParser(&req); err != nil {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if err := utils.ValidateStruct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}

	if err := t.service.AddTable(c.Context(), req); err != nil {
		switch err {
		case exceptions.ErrDuplicatedTableName:
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Table name already exists",
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Table added successfully",
	})
}

// Find Table By ID
// @Summary Find Table By ID
// @Description Find table by ID.
// @Tags Manage
// @Accept json
// @Produce json
// @Param id path string true "Table ID"
// @Success 200 {object} responses.FindTableResponse
// @Failure 400
// @Failure 401
// @Failure 500
// @Router /manage/tables/{id} [get]
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
func (t *tableHandler) FindTableByID(c *fiber.Ctx) error {
	id, err := utils.ValidateUUID(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid UUID",
		})
	}

	table, err := t.service.FindTableByID(c.Context(), *id)
	if err != nil {
		switch err {
		case exceptions.ErrTableNotFound:
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Table not found",
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(table)
}
