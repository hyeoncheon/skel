package actions

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gobuffalo/buffalo"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/cloudfoundry"
)

func init() {
	gothic.Store = App().SessionStore

	uartProvider := cloudfoundry.New(
		os.Getenv("UART_URL"),
		os.Getenv("UART_KEY"),
		os.Getenv("UART_SECRET"),
		fmt.Sprintf("%s%s", App().Host, "/auth/uart/callback"),
		"profile",
	)
	uartProvider.SetName("uart")

	goth.UseProviders(
		uartProvider,
	)
}

func AuthCallback(c buffalo.Context) error {
	user, err := gothic.CompleteUserAuth(c.Response(), c.Request())
	if err != nil {
		return c.Error(http.StatusUnauthorized, err)
	}
	// Do something with the user, maybe register them/sign them in
	return c.Render(http.StatusOK, r.JSON(user))
}
