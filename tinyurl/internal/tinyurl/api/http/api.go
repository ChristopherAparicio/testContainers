package http

import (
	"encoding/json"
	"net/http"

	"github.com/christapa/tinyurl/internal/tinyurl/usecases"
	tinyError "github.com/christapa/tinyurl/pkg/error"
	"github.com/christapa/tinyurl/pkg/logger"
	"github.com/labstack/echo/v4"
)

//
//go:generate oapi-codegen -package http -o openapi.gen.server.go -generate server openapi.json
//go:generate oapi-codegen -package http -o openapi.gen.types.go -generate types openapi.json
//go:generate oapi-codegen -package http -o openapi.gen.spec.go -generate spec openapi.json

type HttpHandler struct {
	Service usecases.URL
}

func NewHttpHandler(service usecases.URL) *HttpHandler {
	return &HttpHandler{
		Service: service,
	}
}

// POST :
func (h HttpHandler) PostCreate(c echo.Context) error {
	var body PostCreateJSONBody
	err := c.Bind(&body)
	if err != nil {
		httpError(c, tinyError.New(tinyError.InvalidArgument, err.Error()))
	}

	url, err := h.Service.CreateShortenUrl(c.Request().Context(), body.OriginalUrl, apiToDomainExpiration(body.ExpirationDate))
	if err != nil {
		logger.Errorf("Failed to create shorten URL: %v", err)
		return httpError(c, err)

	}

	return c.JSON(http.StatusCreated, domainUrlToApi(url))

}

// GET : /:<shortUrl>
func (h HttpHandler) GetSlug(c echo.Context, slug string) error {
	fullUrl, err := h.Service.GetOriginalUrl(c.Request().Context(), slug)
	if err != nil {
		logger.Errorf("Failed to get original URL: %v", err)
		return httpError(c, err)
	}

	return c.Redirect(301, fullUrl)
}

type ApplicationJsonErrorBody struct {
	Message  string `json:"message"`
	Title    string `json:"title"`
	Status   int    `json:"status"`
	Details  string `json:"details"`
	Instance string `json:"instance"`
}

// httpError : Return an error in the format of application/problem+json
func httpError(c echo.Context, err error) error {
	c.Response().Header().Set(echo.HeaderContentType, "application/problem+json")

	tinyErr := tinyError.NewErrorFromDomain(err)

	var body ApplicationJsonErrorBody
	if tinyErr == nil {
		c.Response().WriteHeader(500)
		body = newDefaultApplicationJsonProblem()

	} else {
		c.Response().WriteHeader(GetHttpCode(tinyErr))
		body = newApplicationJsonErrorBodyFromTinyError(tinyErr, c.Request().URL.Path)
	}

	return json.NewEncoder(c.Response()).Encode(body)
}

func newApplicationJsonErrorBodyFromTinyError(err *tinyError.Error, path string) ApplicationJsonErrorBody {
	return ApplicationJsonErrorBody{
		Message:  GetUserFriendlyMessage(err),
		Title:    GetUserFriendlyMessage(err),
		Status:   GetHttpCode(err),
		Details:  GetUserFriendlyMessage(err),
		Instance: path,
	}
}

func newDefaultApplicationJsonProblem() ApplicationJsonErrorBody {
	return ApplicationJsonErrorBody{
		Message:  "An error occurred",
		Title:    "An error occurred",
		Status:   500,
		Details:  "An error occurred",
		Instance: "/",
	}
}
