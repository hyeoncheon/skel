package actions

import (
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo-pop/pop/popmw"
	"github.com/gobuffalo/envy"
	csrf "github.com/gobuffalo/mw-csrf"
	forcessl "github.com/gobuffalo/mw-forcessl"
	i18n "github.com/gobuffalo/mw-i18n"
	paramlogger "github.com/gobuffalo/mw-paramlogger"
	"github.com/gobuffalo/packr"
	"github.com/markbates/goth/gothic"
	"github.com/unrolled/secure"

	"github.com/hyeoncheon/skel/models"
)

// global variables
var (
	ENV = envy.Get("GO_ENV", "development")
	app *buffalo.App
	T   *i18n.Translator
)

// App is the root of the application
func App() *buffalo.App {
	if app == nil {
		app = buffalo.New(buffalo.Options{
			Env:         ENV,
			SessionName: "_skel_session",
		})

		app.Use(forceSSL())
		app.Use(paramlogger.ParameterLogger)
		app.Use(csrf.New)
		app.Use(popmw.Transaction(models.DB))
		app.Use(translations())

		app.GET("/", HomeHandler)
		app.GET("/logout", LogoutHandler)

		// authorization with uart, placed on different app stack
		auth := app.Group("/auth")
		authHandler := buffalo.WrapHandlerFunc(gothic.BeginAuthHandler)
		auth.GET("/{provider}", authHandler)
		auth.GET("/{provider}/callback", AuthCallback)

		// application wide middlewares
		app.Use(authorizeKeeper)
		app.Middleware.Skip(authorizeKeeper, HomeHandler)
		app.Use(contextMapper)

		app.GET("/profile", ProfileShow)

		dr := DocsResource{}
		d := app.Group("/docs")
		d.Resource("/", dr)
		d.GET("/{lang}/{permalink}/", dr.ShowByLang).Name("docLangPath")

		// resources for administrators
		u := app.Resource("/users", UsersResource{})
		u.Use(adminKeeper)

		app.ServeFiles("/", assetsBox)
	}

	return app
}

// translations will load locale files, set up the translator `actions.T`,
// and will return a middleware to use to load the correct locale for each
// request.
// for more information: https://gobuffalo.io/en/docs/localization
func translations() buffalo.MiddlewareFunc {
	var err error
	if T, err = i18n.New(packr.NewBox("../locales"), "en-US"); err != nil {
		app.Stop(err)
	}
	return T.Middleware()
}

func t(c buffalo.Context, str string, args ...interface{}) string {
	s := T.Translate(c, str, args...)
	if s == str {
		c.Logger().WithField("category", "i18n").Warnf("UNTRANSLATED: %v", str)
	}
	return s
}

// forceSSL will return a middleware that will redirect an incoming request
// if it is not HTTPS. "http://example.com" => "https://example.com".
// This middleware does **not** enable SSL. for your application. To do that
// we recommend using a proxy: https://gobuffalo.io/en/docs/proxy
// for more information: https://github.com/unrolled/secure/
func forceSSL() buffalo.MiddlewareFunc {
	return forcessl.Middleware(secure.Options{
		SSLRedirect:     ENV == "production",
		SSLProxyHeaders: map[string]string{"X-Forwarded-Proto": "https"},
	})
}
