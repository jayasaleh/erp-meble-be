package middleware

import (
	"strconv"

	"real-erp-mebel/be/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	limiterMemory "github.com/ulule/limiter/v3/drivers/store/memory"
	"go.uber.org/zap"
)

var rateLimiter *limiter.Limiter

// InitRateLimiter menginisialisasi rate limiter
func InitRateLimiter() error {
	// Create rate limiter: 100 requests per minute per IP
	rate, err := limiter.NewRateFromFormatted("100-M")
	if err != nil {
		return err
	}

	store := limiterMemory.NewStore()
	rateLimiter = limiter.New(store, rate)

	utils.GetLogger().Info("Rate limiter initialized",
		zap.String("rate", "100 requests per minute"),
	)

	return nil
}

// RateLimitMiddleware adalah middleware untuk rate limiting
func RateLimitMiddleware() gin.HandlerFunc {
	logger := utils.GetLogger()

	return func(c *gin.Context) {
		// Get IP address
		ip := c.ClientIP()

		// Get context for limiter
		context, err := rateLimiter.Get(c, ip)
		if err != nil {
			logger.Error("Rate limiter error",
				zap.String("ip", ip),
				zap.Error(err),
			)
			c.Next()
			return
		}

		// Set rate limit headers
		c.Header("X-RateLimit-Limit", strconv.FormatInt(context.Limit, 10))
		c.Header("X-RateLimit-Remaining", strconv.FormatInt(context.Remaining, 10))
		c.Header("X-RateLimit-Reset", strconv.FormatInt(context.Reset, 10))

		// Check if limit exceeded
		if context.Reached {
			logger.Warn("Rate limit exceeded",
				zap.String("ip", ip),
				zap.String("path", c.Request.URL.Path),
			)
		utils.ErrorResponse(c, 429, "Too many requests", map[string]string{
			"error": "Rate limit exceeded. Please try again later.",
		})
			c.Abort()
			return
		}

		c.Next()
	}
}

