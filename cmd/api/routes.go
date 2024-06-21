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
	r.POST("/product", app.createProduct)
	r.POST("/category", app.createCategory)
	// 	r.GET("/user/:id", app.getUserByID)
	//	r.GET("/users", app.getUsers)

	// TODO: Think of a better solution for creating or updating the roles
	/*
	* Executing this path will create the default roles and its permissions
	* (which now are hardcoded in the defaultRoles method)
	 */
	r.GET("/defaultroles", app.defaultRoles)
}
