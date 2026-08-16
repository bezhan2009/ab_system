package middlewares

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"github.com/gin-gonic/gin"
)

func DebugBody() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Body == nil {
			c.Next()
			return
		}

		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			fmt.Println("[DEBUG] failed to read body:", err)
			c.Next()
			return
		}

		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		if len(bodyBytes) == 0 {
			c.Next()
			return
		}

		fmt.Println("========== HTTP REQUEST ==========")
		fmt.Println("Method:", c.Request.Method)
		fmt.Println("Path:", c.Request.URL.Path)

		var prettyJSON bytes.Buffer
		if err := json.Indent(&prettyJSON, bodyBytes, "", "  "); err == nil {
			fmt.Println("Body (JSON):")
			fmt.Println(prettyJSON.String())
		} else {
			fmt.Println("Body (RAW):")
			fmt.Println(string(bodyBytes))
		}

		fmt.Println("==================================")

		c.Next()
	}
}
