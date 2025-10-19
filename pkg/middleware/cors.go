package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SetupCORS sets up CORS middleware for the given gin engine.
func SetupCORS(allowedOrigins []string) gin.HandlerFunc {
	// Log once when middleware is created
	if len(allowedOrigins) == 0 {
		log.Println("[CORS] ERROR: No allowed origins configured! Check your ALLOWED_ORIGINS env var.")
	} else {
		log.Printf("[CORS] Middleware initialized with allowed origins: %v", allowedOrigins)
	}

	return func(c *gin.Context) {
		origin := c.GetHeader("Origin") // exact string Chrome sends

		// Log for debugging
		log.Printf("[CORS] Request from origin: '%s', Method: %s, Path: %s", origin, c.Request.Method, c.Request.URL.Path)

		// Always set common headers first
		c.Header("Vary", "Origin")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")

		// Check if origin is in the allowed list
		isAllowed := false
		if origin != "" {
			for _, allowed := range allowedOrigins {
				log.Printf("[CORS] Comparing origin '%s' with allowed '%s'", origin, allowed)
				if origin == allowed {
					isAllowed = true
					c.Header("Access-Control-Allow-Origin", origin)
					c.Header("Access-Control-Allow-Credentials", "true")
					log.Printf("[CORS] ✓ Origin %s is ALLOWED", origin)
					break
				}
			}
		}

		if !isAllowed && origin != "" {
			log.Printf("[CORS] ✗ Origin '%s' is NOT in allowed list: %v", origin, allowedOrigins)
			log.Printf("[CORS] WARNING: Request will likely fail due to CORS")
		}

		// Handle preflight
		if c.Request.Method == http.MethodOptions {
			log.Printf("[CORS] Handling OPTIONS preflight request")
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}
