package logger

import (
	"github.com/go-resty/resty/v2"
	"log"
	"log/slog"
)

func Info(msg string, args ...any) {
	if len(args) == 0 {
		logger.Error(msg)
		return
	}

	logger.Info(msg, genSlogAttrs(args)...)
}

func ApiError(resp *resty.Response) {
	// Extract request data
	requestData := map[string]interface{}{
		"method":  resp.Request.Method,
		"url":     resp.Request.URL,
		"headers": resp.Request.Header,
		"body":    resp.Request.Body,
	}

	// Log request and response
	logger.Error(resp.String(),
		slog.String("status", resp.Status()),
		slog.Any("request", requestData),
		slog.Any("response", resp.Error()),
	)
}

// Error receives a message as the first parameter, error as the last, and optional args in between
func Error(message string, args ...any) {
	// Check if the last argument is of type error
	if len(args) > 0 {
		if _, ok := args[len(args)-1].(error); ok {
			logger.Error(message, genSlogAttrs(args[:len(args)-1])...)
			return
		}
	}

	// If no error is passed as the last argument
	logger.Error(message, genSlogAttrs(args)...)
}

func genSlogAttrs(args []any) (retval []any) {
	argsLen := len(args)
	const orphanKey = "NO-KEY"

	var orphan any
	if argsLen%2 != 0 {
		log.Println("uneven number of arguments, placing last provided value under the key '" + orphanKey + "'")
		orphan = args[argsLen-1]
		args = args[:argsLen-1] // Truncate the last element
	}

	// Create an empty slice with a capacity
	attrs := make([]slog.Attr, 0, argsLen/2+1) // +1 to accommodate the potential orphan

	// Iterate over the key-value pairs
	for i := 0; i < len(args); i += 2 {
		key, ok := args[i].(string)
		if !ok {
			log.Println("non-string key encountered, skipping:", args[i])
			continue
		}
		attrs = append(attrs, slog.Any(key, args[i+1]))
	}

	// Handle the orphan case
	if orphan != nil {
		attrs = append(attrs, slog.Any(orphanKey, orphan))
	}

	// Flatten the attrs to retval
	for _, attr := range attrs {
		retval = append(retval, attr.Key, attr.Value)
	}

	return
}
