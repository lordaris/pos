package main

import "github.com/gin-gonic/gin"

func Router(r *gin.Engine, app *application) {
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello world",
		})
	})

	r.POST("/user", app.createUser)
	r.PUT("/user/:id", app.updateUser)
	r.PUT("/user/:id/role", app.updateUserRole)
	r.GET("/users/:role", app.getUsersByRole)
	r.POST("/product", app.createProduct)
	r.POST("/category", app.createCategory)
	r.GET("/user/:id", app.getUser)
	//	r.GET("/users", app.getUsers)
	r.POST("/roles", app.createRoles)
}
