package responsewriter

import (
	"bytes"

	"github.com/gin-gonic/gin"
)

// 重写gin的ResponseWriter，用户接收响应体
type ResponseWriter struct {
	gin.ResponseWriter
	Body *bytes.Buffer
}

func (rw *ResponseWriter) Write(b []byte) (int, error) {

	rw.Body.Write(b)

	return rw.ResponseWriter.Write(b)
}

func (rw *ResponseWriter) WriteString(s string) (int, error) {

	rw.Body.WriteString(s)

	return rw.ResponseWriter.WriteString(s)
}
