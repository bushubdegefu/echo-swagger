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

var HandlerDefault = New()

// New returns custom Echo middleware handler
func New(config ...Config) echo.HandlerFunc {
	cfg := configDefault(config...)

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
		once.Do(func() {
			// This is the registered route path pattern, e.g., /swagger/*
			prefix = strings.TrimSuffix(c.Path(), "*")

			forwardedPrefix := getForwardedPrefix(c)
			if forwardedPrefix != "" {
				prefix = forwardedPrefix + prefix
			}

			if cfg.URL == "" {
				cfg.URL = path.Join(prefix, defaultDocURL)
			}
		})

		// Extract the actual path being requested
		requestPath := strings.TrimPrefix(c.Request().URL.Path, prefix)

		switch requestPath {
		case defaultIndex:
			c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTMLCharsetUTF8)
			return index.Execute(c.Response().Writer, cfg)

		case defaultDocURL:
			doc, err := swag.ReadDoc(cfg.InstanceName)
			if err != nil {
				return err
			}
			return c.JSONBlob(http.StatusOK, []byte(doc))

		case "", "/":
			return c.Redirect(http.StatusMovedPermanently, path.Join(prefix, defaultIndex))

		default:
			// Reset URL path so http.FileServer can find the correct file
			c.Request().URL.Path = requestPath
			fs.ServeHTTP(c.Response().Writer, c.Request())
			return nil
		}
	}
}

// getForwardedPrefix extracts X-Forwarded-Prefix header
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
