package middleware

import (
	"log"
	"net/http"
	"strings"
	"time"

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
		origin := c.GetHeader("Origin") // exact string browsers send
		requestID, _ := c.Get("request_id")

		// Rich request-level logging
		log.Printf("[CORS] %s Request from origin='%s' method=%s path=%s remote_ip=%s request_id=%v",
			time.Now().Format(time.RFC3339), origin, c.Request.Method, c.Request.URL.Path, c.ClientIP(), requestID)

		// Snapshot a few headers useful for debugging
		log.Printf("[CORS] Headers: User-Agent='%s' Referer='%s' Content-Type='%s'",
			c.GetHeader("User-Agent"), c.GetHeader("Referer"), c.GetHeader("Content-Type"))

		// Always set common headers first
		c.Header("Vary", "Origin")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With, X-Request-ID")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Expose-Headers", "Content-Length, X-Request-ID")
		c.Header("Access-Control-Max-Age", "600") // seconds

		// Check if origin is in the allowed list
		isAllowed := false
		if origin != "" {
			for _, allowed := range allowedOrigins {
				log.Printf("[CORS] Comparing origin '%s' with allowed '%s'", origin, allowed)
				// allow exact match or scheme-insensitive match (strip scheme)
				if origin == allowed || stripScheme(origin) == stripScheme(allowed) {
					isAllowed = true
					c.Header("Access-Control-Allow-Origin", origin)
					// Only set credentials when origin is explicit
					c.Header("Access-Control-Allow-Credentials", "true")
					log.Printf("[CORS] ✓ Origin %s is ALLOWED (matched %s)", origin, allowed)
					break
				}
			}
		}

		if !isAllowed && origin != "" {
			log.Printf("[CORS] ✗ Origin '%s' is NOT in allowed list: %v", origin, allowedOrigins)
		}

		// Handle preflight
		if c.Request.Method == http.MethodOptions {
			log.Printf("[CORS] Handling OPTIONS preflight request origin=%s allowed=%t", origin, isAllowed)
			if isAllowed {
				c.AbortWithStatus(http.StatusNoContent)
			} else {
				// When not allowed, return Forbidden to make the failure explicit in logs/tools
				c.AbortWithStatus(http.StatusForbidden)
			}
			return
		}

		c.Next()
	}
}

// stripScheme removes http(s) scheme and trailing slashes for simple comparison
func stripScheme(u string) string {
	u = strings.TrimSpace(u)
	if strings.HasPrefix(u, "http://") {
		u = strings.TrimPrefix(u, "http://")
	} else if strings.HasPrefix(u, "https://") {
		u = strings.TrimPrefix(u, "https://")
	}
	u = strings.TrimSuffix(u, "/")
	return u
}
