package v1

import (
	"github.com/gin-gonic/gin"
)

func ApiV1(r *gin.Engine) {
	v1Group := r.Group("/api/v1")
	{
		v1GroupAuth := v1Group.Group("/auth")
		{
			v1GroupAuth.POST("/register", Register)
		}
	}
}
