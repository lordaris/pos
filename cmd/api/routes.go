package main

import "github.com/gin-gonic/gin"

func Router(r *gin.Engine, app *application) {
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello world",
		})
	})

	// Users and roles
	r.POST("/user", app.createUser)
	r.GET("/user/:id", app.getUser)
	r.GET("/users/:role", app.getUsersByRole)
	r.PUT("/user/:id", app.updateUser)
	r.PUT("/user/:id/role", app.updateUserRole)
	r.POST("/roles", app.createRoles)

	// Products and categories
	r.POST("/product", app.createProduct)
	r.POST("/category", app.createCategory)
	r.POST("/promotion", app.productPromotion)
}
