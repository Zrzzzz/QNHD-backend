package safety

import "github.com/gin-gonic/gin"

func Safety() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1")
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Download-Options", "noopen")
		c.Header("Content-Security-Policy", "'none'")
		c.Header("Referrer-Policy", "no-referrer")
		c.Header("X-Permitted-Cross-Domain-Policies", "master-only")

		c.Next()
	}
}
