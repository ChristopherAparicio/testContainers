package http

import (
	"time"

	"github.com/christapa/tinyurl/internal/tinyurl/domain"
	tinyError "github.com/christapa/tinyurl/pkg/error"
)

func apiToDomainExpiration(expiration *int) time.Time {
	if expiration == nil {
		return time.Time{}
	}

	return time.Unix(int64(*expiration), 0)
}

func domainUrlToApi(url domain.Url) URL {
	var expiration *int
	if !url.Expiration.IsZero() {
		expirationTimestamp := int(url.Expiration.Unix())
		expiration = &expirationTimestamp
	}

	return URL{
		OriginalUrl:    url.OriginalURL,
		ShortenedUrl:   url.ShortenURL,
		ExpirationDate: expiration,
	}
}

func GetHttpCode(e *tinyError.Error) int {
	return CodeToHTTP(e.Code)
}

func CodeToHTTP(code tinyError.Code) int {
	switch code {
	case tinyError.OK:
		return 200
	case tinyError.InvalidArgument:
		return 400
	case tinyError.NotFound:
		return 404
	case tinyError.AlreadyExists:
		return 409
	case tinyError.PermissionDenied:
		return 403
	case tinyError.Unauthenticated:
		return 401
	case tinyError.DeadlineExceeded:
		return 504
	case tinyError.Internal:
		return 500
	default:
		return 500
	}
}
func GetUserFriendlyMessage(e *tinyError.Error) string {
	switch e.Code {
	case tinyError.OK:
		return "OK"
	case tinyError.InvalidArgument:
		return "Invalid Argument"
	case tinyError.NotFound:
		return "Not Found"
	case tinyError.AlreadyExists:
		return "Already Exists"
	case tinyError.PermissionDenied:
		return "Permission Denied"
	case tinyError.Unauthenticated:
		return "Unauthenticated"
	case tinyError.DeadlineExceeded:
		return "Deadline Exceeded"
	case tinyError.Internal:
		return "Internal Server Error"
	default:
		return "Internal Server Error"
	}
}
