package tadapter

import (
	"context"
	"domofonEmulator/client/models"
	"domofonEmulator/client/web/views/components"
	"fmt"
	"io"

	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
)

func Render(c *fiber.Ctx, component templ.Component, code int) error {
	return adaptor.HTTPHandler(templ.Handler(component, templ.WithStatus(code)))(c)
}

func RenderIntercomAndNotificationResponse(message string, intercomData *models.Intercom) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		if _, err := fmt.Fprintf(w, `<div id="notification-area" hx-swap-oob="true">`); err != nil {
			return err
		}
		if err := components.IntercomControlResponse(message, "success").Render(ctx, w); err != nil {
			return err
		}
		if _, err := fmt.Fprintf(w, `</div>`); err != nil {
			return err
		}

		if _, err := fmt.Fprintf(w, `<div id="domofon-card-%d" hx-swap-oob="true">`, intercomData.ID); err != nil {
			return err
		}
		if err := components.SingleDomofonCard(*intercomData).Render(ctx, w); err != nil {
			return err
		}
		if _, err := fmt.Fprintf(w, `</div>`); err != nil {
			return err
		}

		return nil
	})
}

func RenderDoorControlResponse(message, status string, intercom *models.Intercom) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		fmt.Fprintf(w, `<div id="notification-area" hx-swap-oob="true">`)
		components.IntercomControlResponse(message, status).Render(ctx, w)
		fmt.Fprintf(w, `</div>`)

		fmt.Fprintf(w, `<div id="domofon-card-%d" hx-swap-oob="true">`, intercom.ID)
		components.SingleDomofonCard(*intercom).Render(ctx, w)
		fmt.Fprintf(w, `</div>`)

		return nil
	})
}
