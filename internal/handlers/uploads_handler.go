package handlers

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

type UploadsHandler struct {
	Folder string
}

const maxUploadChenin int64 = 1024 * 1024

func (h UploadsHandler) SubirImagenChenin(c *fiber.Ctx) error {
	file, err := c.FormFile("imagen")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "manda una imagen en el campo 'imagen'"})
	}
	if file.Size > maxUploadChenin {
		return c.Status(fiber.StatusRequestEntityTooLarge).JSON(fiber.Map{"error": "la imagen no debe pesar mas de 1MB"})
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	permitidas := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".webp": true, ".gif": true}
	if !permitidas[ext] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "solo se permiten imagenes jpg, png, webp o gif"})
	}

	nombreBase := strings.TrimSuffix(filepath.Base(file.Filename), ext)
	nombreBase = strings.ReplaceAll(strings.ToLower(nombreBase), " ", "-")
	nombreFinal := fmt.Sprintf("%d-%s%s", time.Now().Unix(), nombreBase, ext)
	rutaFisica := filepath.Join(h.Folder, nombreFinal)

	if err := c.SaveFile(file, rutaFisica); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"mensaje": "imagen subida",
		"ruta":    "/uploads/" + nombreFinal,
	})
}
