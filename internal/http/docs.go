package http

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"tests2/internal/models"
	"tests2/internal/utils"
)

func (h *Handler) initRoutesDocs(v1 fiber.Router) {
	v1.Get("/docs", h.getDocs)
	v1.Post("/", h.insertDoc)
	docGroup := v1.Group("/docs")

	docGroup.Get("/:docID", h.getDocByName)
	docGroup.Delete("/:docID", h.deleteDocByName)
}

func (h *Handler) insertDoc(ctx *fiber.Ctx) error {
	multipart, err := ctx.FormFile("file")
	if err != nil {
		return h.Response(ctx, http.StatusBadRequest, nil, err)
	}
	buffer, err := multipart.Open()
	if err != nil {
		return h.Response(ctx, http.StatusBadRequest, nil, err)
	}
	defer buffer.Close()

	file := &models.FileInfo{
		ObjectName:  multipart.Filename,
		FileBuffer:  buffer,
		ContentType: multipart.Header["Content-Type"][0],
		FileSize:    multipart.Size,
	}

	if err := h.services.Doc.UploadDoc(ctx.UserContext(), file); err != nil {
		if errors.Is(err, utils.ErrFile) {
			return h.Response(ctx, http.StatusUnsupportedMediaType, nil, err)
		} else {
			return h.Response(ctx, http.StatusInternalServerError, nil, err)
		}
	}

	return h.Response(ctx, http.StatusCreated, nil, nil)
}

func (h *Handler) getDocs(ctx *fiber.Ctx) error {

	return nil
}

func (h *Handler) getDocByName(ctx *fiber.Ctx) error {
	buffer, err := h.services.Doc.DownloadDoc(ctx.Context(), ctx.Params("filename"))
	if err != nil {
		return h.Response(ctx, http.StatusInternalServerError, nil, err)
	}

	return h.Response(ctx, http.StatusOK, buffer, nil)
}

func (h *Handler) deleteDocByName(ctx *fiber.Ctx) error {
	return nil
}
