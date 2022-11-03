package delivery

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"moonlay-test/domain"
	"moonlay-test/utils"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type ListHandler struct {
	ListUsecase domain.ListUsecase
}

type StoreListRequest struct {
	ParentID    uint64                `json:"parent_id,omitempty"`
	Title       string                `json:"title" validate:"required,max=100"`
	Description string                `json:"description" validate:"required,max=1000"`
	File        *multipart.FileHeader `json:"file,omitempty"`
}

type UpdateListRequest struct {
	ID          uint64                `json:"id" validate:"required"`
	ParentID    uint64                `json:"parent_id,omitempty"`
	Title       string                `json:"title" validate:"required,max=100"`
	Description string                `json:"description" validate:"required,max=1000"`
	File        *multipart.FileHeader `json:"file,omitempty"`
}

func NewListHttpHandler(e *echo.Echo, ListUsecase domain.ListUsecase) {
	handler := &ListHandler{ListUsecase: ListUsecase}

	listGroup := e.Group("list")
	listGroup.GET("", handler.Fetch)
	listGroup.POST("", handler.Store)
	listGroup.GET("/:id", handler.getById)
	listGroup.PUT("/:id", handler.Update)
	listGroup.DELETE("/:id", handler.Delete)

	listGroup.GET("/:parent_id/sub_list", handler.FetchSublist)
	listGroup.POST("/:parent_id/sub_list", handler.StoreSublist)

	listGroup.GET("/:parent_id/sub_list/:id", handler.GetSublistById)
	listGroup.PUT("/:parent_id/sub_list/:id", handler.UpdateSublist)
	listGroup.DELETE("/:parent_id/sub_list/:id", handler.DeleteSublist)
}

func validateStoreRequest(m *StoreListRequest) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)

	if err != nil {
		return false, err
	}

	return true, nil
}

func validateUpdateRequest(m *UpdateListRequest) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)

	if err != nil {
		return false, err
	}

	return true, nil
}

func uploadFileHandler(file *multipart.FileHeader, defaultName string) (fileName string, err error) {
	src, err := file.Open()
	if err != nil {
		return defaultName, err
	}

	defer src.Close()

	fileType := strings.Split(file.Filename, ".")[1]

	if fileType != "pdf" && fileType != "txt" {
		err = errors.New("Key: 'StoreListRequest.File' Error:Field validation for 'File' failed on the 'mimes:pdf,txt' tag")
		return defaultName, err
	}

	cwd, _ := os.Getwd()
	err = os.MkdirAll("uploads", os.ModePerm)
	if err != nil {
		return defaultName, err
	}

	checkDefaultName := strings.Split(defaultName, ".")
	if len(checkDefaultName) > 1 {
		defaultName = checkDefaultName[0]
	} else {
		defaultName = time.Now().Format("20060102150405_999999")
	}

	fileName = fmt.Sprintf("%s.%s", defaultName, fileType)
	path := filepath.Join(cwd, "uploads", fileName)

	dst, err := os.Create(filepath.FromSlash(path))
	if _, err = io.Copy(dst, src); err != nil {
		return defaultName, err
	}
	defer dst.Close()

	return fileName, err
}

func (h *ListHandler) Fetch(c echo.Context) (err error) {
	ctx := c.Request().Context()
	lists, err := h.ListUsecase.Fetch(ctx)

	if err != nil {
		return c.JSON(utils.GetStatusCode(err), utils.ResponseMessage{Message: c.Param("data")})
	}

	return c.JSON(http.StatusOK, lists)
}

func (h *ListHandler) getById(c echo.Context) (err error) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)

	if err != nil {
		return c.JSON(utils.GetStatusCode(err), utils.ResponseMessage{Message: err.Error()})
	}

	ctx := c.Request().Context()
	list, err := h.ListUsecase.GetById(ctx, id)

	if err != nil {
		return c.JSON(utils.GetStatusCode(err), utils.ResponseMessage{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, list)
}

func (h *ListHandler) Store(c echo.Context) (err error) {
	var fileName string
	req := StoreListRequest{
		Title:       c.FormValue("title"),
		Description: c.FormValue("description"),
	}

	file, err := c.FormFile("file")
	if err == nil {
		fileName, err = uploadFileHandler(file, "")
		if err != nil {
			return c.JSON(utils.GetStatusCode(err), utils.ResponseMessage{Message: err.Error()})
		}
	}

	var validated bool
	if validated, err = validateStoreRequest(&req); !validated {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	ctx := c.Request().Context()
	list := domain.List{Title: req.Title, Description: req.Description, File: fileName}
	err = h.ListUsecase.Store(ctx, &list)

	if err != nil {
		return c.JSON(utils.GetStatusCode(err), utils.ResponseMessage{Message: err.Error()})

	}

	return c.JSON(http.StatusCreated, list)
}

func (h *ListHandler) Update(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)

	if err != nil {
		return c.JSON(utils.GetStatusCode(err), utils.ResponseMessage{Message: err.Error()})
	}

	ctx := c.Request().Context()
	list, err := h.ListUsecase.GetById(ctx, id)

	if err != nil {
		return c.JSON(utils.GetStatusCode(err), utils.ResponseMessage{Message: err.Error()})
	}

	req := UpdateListRequest{
		ID:          id,
		Title:       c.FormValue("title"),
		Description: c.FormValue("description"),
	}

	var validated bool
	if validated, err = validateUpdateRequest(&req); !validated {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	var fileName string
	file, err := c.FormFile("file")

	if err == nil {
		fileName, err = uploadFileHandler(file, list.File)
		if err != nil {
			return c.JSON(utils.GetStatusCode(err), utils.ResponseMessage{Message: err.Error()})
		}
	}

	list.Title = req.Title
	list.Description = req.Description
	list.File = fileName

	err = h.ListUsecase.Update(ctx, &list)

	if err != nil {
		return c.JSON(utils.GetStatusCode(err), utils.ResponseMessage{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, list)
}

func (h *ListHandler) Delete(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(utils.GetStatusCode(err), utils.ResponseMessage{Message: err.Error()})
	}

	ctx := c.Request().Context()
	err = h.ListUsecase.Delete(ctx, id)
	if err != nil {
		return c.JSON(utils.GetStatusCode(err), utils.ResponseMessage{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, utils.ResponseMessage{Message: "Deleted"})
}

// Sublist Operations
func (h *ListHandler) FetchSublist(c echo.Context) (err error) {
	parentId, err := strconv.ParseUint(c.Param("parent_id"), 10, 64)
	ctx := c.Request().Context()
	lists, err := h.ListUsecase.FetchSublist(ctx, parentId)

	if err != nil {
		return c.JSON(utils.GetStatusCode(err), utils.ResponseMessage{Message: c.Param("data")})
	}

	return c.JSON(http.StatusOK, lists)
}

func (h *ListHandler) GetSublistById(c echo.Context) (err error) {
	parentId, err := strconv.ParseUint(c.Param("parent_id"), 10, 64)
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)

	if err != nil {
		return c.JSON(utils.GetStatusCode(err), utils.ResponseMessage{Message: err.Error()})
	}

	ctx := c.Request().Context()
	list, err := h.ListUsecase.GetSublistById(ctx, parentId, id)

	if err != nil {
		return c.JSON(utils.GetStatusCode(err), utils.ResponseMessage{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, list)
}

func (h *ListHandler) StoreSublist(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("parent_id"), 10, 64)

	if err != nil {
		return c.JSON(utils.GetStatusCode(err), utils.ResponseMessage{Message: err.Error()})
	}

	ctx := c.Request().Context()
	parent, err := h.ListUsecase.GetById(ctx, id)

	if err != nil {
		return c.JSON(utils.GetStatusCode(err), utils.ResponseMessage{Message: err.Error()})
	}

	var fileName string
	req := StoreListRequest{
		Title:       c.FormValue("title"),
		Description: c.FormValue("description"),
	}

	file, err := c.FormFile("file")
	if err == nil {
		fileName, err = uploadFileHandler(file, "")
		if err != nil {
			return c.JSON(utils.GetStatusCode(err), utils.ResponseMessage{Message: err.Error()})
		}
	}

	var validated bool
	if validated, err = validateStoreRequest(&req); !validated {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	list := domain.List{
		ParentID:    &parent.ID,
		Title:       req.Title,
		Description: req.Description,
		File:        fileName,
	}
	err = h.ListUsecase.Store(ctx, &list)

	if err != nil {
		return c.JSON(utils.GetStatusCode(err), utils.ResponseMessage{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, list)
}

func (h *ListHandler) UpdateSublist(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("parent_id"), 10, 64)

	if err != nil {
		return c.JSON(utils.GetStatusCode(err), utils.ResponseMessage{Message: err.Error()})
	}

	ctx := c.Request().Context()
	_, err = h.ListUsecase.GetById(ctx, id)

	if err != nil {
		return c.JSON(utils.GetStatusCode(err), utils.ResponseMessage{Message: err.Error()})
	}

	var fileName string
	req := UpdateListRequest{
		Title:       c.FormValue("title"),
		Description: c.FormValue("description"),
	}

	file, err := c.FormFile("file")
	if err == nil {
		fileName, err = uploadFileHandler(file, "")
		if err != nil {
			return c.JSON(utils.GetStatusCode(err), utils.ResponseMessage{Message: err.Error()})
		}
	}

	var validated bool
	if validated, err = validateUpdateRequest(&req); !validated {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	list := domain.List{
		ID:          id,
		ParentID:    &req.ParentID,
		Title:       req.Title,
		Description: req.Description,
		File:        fileName,
	}
	err = h.ListUsecase.Update(ctx, &list)

	if err != nil {
		return c.JSON(utils.GetStatusCode(err), utils.ResponseMessage{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, list)
}

func (h *ListHandler) DeleteSublist(c echo.Context) error {
	parentId, err := strconv.ParseUint(c.Param("parent_id"), 10, 64)
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)

	if err != nil {
		return c.JSON(utils.GetStatusCode(err), utils.ResponseMessage{Message: err.Error()})
	}

	ctx := c.Request().Context()
	err = h.ListUsecase.DeleteSublist(ctx, parentId, id)
	if err != nil {
		return c.JSON(utils.GetStatusCode(err), utils.ResponseMessage{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, utils.ResponseMessage{Message: "Deleted"})
}
