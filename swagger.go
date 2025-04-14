package swagger

import (
	"fmt"
	"html/template"
	"net/http"
	"path"
	"strings"
	"sync"

	"github.com/labstack/echo/v4"
	swaggerFiles "github.com/swaggo/files/v2"
	"github.com/swaggo/swag"
)

const (
	defaultDocURL = "doc.json"
	defaultIndex  = "index.html"
)

// HandlerDefault is the default Swagger handler using default config
var HandlerDefault = New()

// New returns custom Echo handler
func New(config ...Config) echo.HandlerFunc {
	cfg := configDefault(config...)

	// Parse the Swagger UI index template
	index, err := template.New("swagger_index.html").Parse(indexTmpl)
	if err != nil {
		panic(fmt.Errorf("echo: swagger middleware error -> %w", err))
	}

	var (
		prefix string
		once   sync.Once
		fs     = http.FileServer(http.FS(swaggerFiles.FS))
	)

	return func(c echo.Context) error {
		// Initialize prefix and URL only once
		once.Do(func() {
			prefix = strings.TrimSuffix(c.Path(), "*")

			forwardedPrefix := getForwardedPrefix(c)
			if forwardedPrefix != "" {
				prefix = forwardedPrefix + prefix
			}

			if cfg.URL == "" {
				cfg.URL = path.Join(prefix, defaultDocURL)
			}
		})

		// Extract request path
		p := c.Param("*")
		if p == "" {
			p = c.Request().URL.Path
			p = strings.TrimPrefix(p, prefix)
		}

		switch p {
		case "", "/":
			// Redirect to index
			return c.Redirect(http.StatusMovedPermanently, path.Join(prefix, defaultIndex))
		case defaultIndex:
			// Serve Swagger UI
			c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTMLCharsetUTF8)
			return index.Execute(c.Response().Writer, cfg)
		case defaultDocURL:
			// Serve generated swagger docs
			doc, err := swag.ReadDoc(cfg.InstanceName)
			if err != nil {
				return err
			}
			return c.JSONBlob(http.StatusOK, []byte(doc))
		default:
			// Serve static files
			fsHandler := echo.WrapHandler(fs)
			return fsHandler(c)
		}
	}
}

// getForwardedPrefix extracts X-Forwarded-Prefix header if available
func getForwardedPrefix(c echo.Context) string {
	headers := c.Request().Header["X-Forwarded-Prefix"]
	if len(headers) == 0 {
		return ""
	}

	prefix := ""
	for _, rawPrefix := range headers {
		endIndex := len(rawPrefix)
		for endIndex > 1 && rawPrefix[endIndex-1] == '/' {
			endIndex--
		}
		if endIndex != len(rawPrefix) {
			prefix += rawPrefix[:endIndex]
		} else {
			prefix += rawPrefix
		}
	}

	return prefix
}
