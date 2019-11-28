package swagger

import (
	"net/http"
	"strings"

	"github.com/ddosakura/sola/v2"
	swaggerFiles "github.com/swaggo/files"
	"github.com/swaggo/swag"
)

// Config stores solaSwagger configuration variables.
type Config struct {
	// The url pointing to API definition
	// (normally swagger.json or swagger.yaml).
	// Default is `doc.json`.
	URL string
}

// URL presents the url pointing to API definition
// (normally swagger.json or swagger.yaml).
func URL(url string) func(c *Config) {
	return func(c *Config) {
		c.URL = url
	}
}

// WrapHandler wraps swaggerFiles.Handler and returns sola.Handler
var WrapHandler = SolaWrapHandler()

// SolaWrapHandler wraps `http.Handler` into `sola.Handler`.
func SolaWrapHandler(confs ...func(c *Config)) sola.Handler {
	handler := swaggerFiles.Handler
	config := &Config{
		URL: "doc.json",
	}
	for _, c := range confs {
		c(config)
	}

	return func(c sola.Context) error {
		var matches []string
		uri := c.Request().RequestURI
		if strings.HasSuffix(uri, "/swagger") {
			c.Response().Header().Add("Location", uri+"/index.html")
			return c.String(http.StatusMovedPermanently, "")
		}
		if strings.HasSuffix(uri, "/swagger/") {
			uri += "index.html"
		}
		if matches = re.FindStringSubmatch(uri); len(matches) != 3 {
			return c.Handle(http.StatusNotFound)(c)
		}
		path := matches[2]
		prefix := matches[1]
		handler.Prefix = prefix

		switch path {
		case "index.html":
			return index.Execute(c.Response(), config)
		case "doc.json":
			doc, err := swag.ReadDoc()
			if err != nil {
				return nil
			}
			return c.Blob(http.StatusOK,
				sola.MIMEApplicationJSONCharsetUTF8, []byte(doc))
		}
		handler.ServeHTTP(c.Response(), c.Request())
		return nil
	}
}
