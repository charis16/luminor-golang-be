package middleware

import (
	"net/http"

	"github.com/charis16/luminor-golang-be/src/utils"
	"github.com/gin-gonic/gin"
)

func AdminRequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr, err := c.Cookie("admin_access_token")

		if err != nil || tokenStr == "" {
			utils.RespondError(c, http.StatusUnauthorized, "missing access token cookie")
			return
		}

		_, claims, err := utils.ValidateAccessToken(tokenStr)
		if err != nil {
			utils.RespondError(c, http.StatusUnauthorized, "invalid or expired token")
			return
		}

		// Inject ke context
		c.Set("user_id", claims.UserID)
		c.Set("role", claims.Role)

		c.Next()
	}
}

func RequireRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("role")
		if !exists || userRole != role {
			utils.RespondError(c, http.StatusForbidden, "forbidden: insufficient role")
			return
		}
		c.Next()
	}
}
