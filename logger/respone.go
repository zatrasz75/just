package logger

import (
	"github.com/gofrs/uuid/v5"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

type ResponLogger struct {
	consoleLogger *log.Logger
}

// NewResponLogger создает новый логгер
func NewResponLogger(output *log.Logger) *ResponLogger {
	return &ResponLogger{consoleLogger: output}
}

func (rl *ResponLogger) LoggingResponse(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		reqID := req.URL.Query().Get("request_id")
		if reqID == "" {
			rID, _ := uuid.NewV4()
			reqID = rID.String()
		}
		req.Header.Set("X-Request-ID", reqID)
		w.Header().Set("X-Request-ID", reqID)
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		lrw := NewLoggingResponseWriter(w)
		next.ServeHTTP(lrw, req)

		methodColor := getMethodColor(req.Method)
		statusColor := getStatusCodeColor(lrw.statusCode)

		// Формируем строку лога с цветами
		str := strings.Builder{}
		str.WriteString(colorGray)
		str.WriteString(">-- ip: ")
		str.WriteString(colorReset)
		str.WriteString(req.RemoteAddr)
		str.WriteString(colorGray)
		str.WriteString(", host: ")
		str.WriteString(colorReset)
		str.WriteString(req.Host)
		str.WriteString(colorGray)
		str.WriteString(" url: ")
		str.WriteString(colorReset)
		str.WriteString(req.URL.Path)
		str.WriteString(colorGray)
		str.WriteString(", method: ")
		str.WriteString(methodColor)
		str.WriteString(req.Method)
		str.WriteString(colorReset)
		str.WriteString(colorGray)
		str.WriteString(" status code: ")
		str.WriteString(statusColor)
		str.WriteString(strconv.Itoa(lrw.statusCode))
		str.WriteString(" ")
		str.WriteString(http.StatusText(lrw.statusCode))
		str.WriteString(colorReset)
		str.WriteString(colorGray)
		str.WriteString(", trace id: ")
		str.WriteString(colorReset)
		str.WriteString(reqID)

		rl.consoleLogger.Println(str.String())
	})
}

func getMethodColor(method string) string {
	switch method {
	case "GET":
		return colorGreen // Зеленый
	case "POST":
		return colorBlue // Синий
	case "PUT":
		return colorYellow // Желтый
	case "PATCH":
		return colorYellow // Желтый
	case "DELETE":
		return colorRed // Красный
	default:
		return colorReset // Обычный цвет
	}
}

func getStatusCodeColor(code int) string {
	switch {
	case code >= 200 && code < 300:
		return colorGreen // Зеленый для успешных ответов
	case code >= 300 && code < 400:
		return colorYellow // Желтый для перенаправлений
	case code >= 400:
		return colorRed // Красный для ошибок
	default:
		return colorReset // Сброс цвета для остальных случаев
	}
}
