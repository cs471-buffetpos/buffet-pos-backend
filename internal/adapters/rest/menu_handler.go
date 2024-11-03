package rest

import (
	"github.com/cs471-buffetpos/buffet-pos-backend/domain/exceptions"
	"github.com/cs471-buffetpos/buffet-pos-backend/domain/requests"
	"github.com/cs471-buffetpos/buffet-pos-backend/domain/usecases"
	"github.com/cs471-buffetpos/buffet-pos-backend/utils"
	"github.com/gofiber/fiber/v2"
)

type MenuHandler interface {
	Create(c *fiber.Ctx) error
	FindAll(c *fiber.Ctx) error
	FindByID(c *fiber.Ctx) error
}

type menuHandler struct {
	service usecases.MenuUseCase
}

func NewMenuHandler(service usecases.MenuUseCase) MenuHandler {
	return &menuHandler{
		service: service,
	}
}

// Add Menu
// @Summary Add Menu
// @Description Add a new menu.
// @Tags Manage
// @Accept mpfd
// @Param request formData requests.AddMenuRequest true "Add Menu request"
// @Param image formData file true "Menu Image"
// @Success 200 {object} responses.SuccessResponse
// @Router /manage/menus [post]
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
func (m *menuHandler) Create(c *fiber.Ctx) error {
	var req requests.AddMenuRequest
	if err := utils.PopulateStructFromForm(c, &req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse form data",
		})
	}

	if validationErr := utils.ValidateStruct(req); validationErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": validationErr.Message,
		})
	}

	file, err := utils.OpenFile(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if err := m.service.Create(c.Context(), &req, file); err != nil {
		switch err {
		case exceptions.ErrDuplicatedMenuName:
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Menu name already exist",
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Create Menu successfully",
	})
}

// Find All Menus
// @Summary Find All Menus
// @Description Find all menus.
// @Tags Manage
// @Accept json
// @Produce json
// @Success 200 {array} responses.MenuDetail
// @Router /manage/menus [get]
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
func (m *menuHandler) FindAll(c *fiber.Ctx) error {
	res, err := m.service.FindAll(c.Context())
	if err != nil {
		switch err {
		case exceptions.ErrMenuNotFound:
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Menu not have",
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Get All Menus",
		"data":    res,
	})
}

// Find Menu By ID
// @Summary Find Menu By ID
// @Description Find menu by ID.
// @Tags Manage
// @Accept json
// @Produce json
// @Param id path string true "Menu ID"
// @Success 200 {object} responses.MenuDetail
// @Router /manage/menus/{id} [get]
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
func (m *menuHandler) FindByID(c *fiber.Ctx) error {
	id, err := utils.ValidateUUID(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid UUID",
		})
	}

	res, err := m.service.FindByID(c.Context(), *id)
	if err != nil {
		switch err {
		case exceptions.ErrMenuNotFound:
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Menu not found",
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(res)
}