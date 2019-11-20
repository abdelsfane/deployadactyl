package interfaces

import "github.com/gin-gonic/gin"

// Endpoints interface.
type Endpoints interface {
	Deploy(c *gin.Context)
}
