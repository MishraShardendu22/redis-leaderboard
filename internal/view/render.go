package view

import (
	"bytes"

	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
)

func Render(c *fiber.Ctx, component templ.Component) error {
	var buffer bytes.Buffer
	if err := component.Render(c.UserContext(), &buffer); err != nil {
		return err
	}

	c.Set(fiber.HeaderContentType, fiber.MIMETextHTMLCharsetUTF8)
	_, err := c.Write(buffer.Bytes())
	return err
}
