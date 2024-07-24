package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Router(r *gin.Engine, app *application) {
	// r.GET("/", func(c *gin.Context) {
	// 	c.JSON(200, gin.H{
	// 		"message": "Hello world",
	// 	})
	// })

	// Users and roles
	r.POST("/user", app.createUser)
	r.GET("/user/:id", app.getUser)
	r.GET("/users/:role", app.getUsersByRole)
	r.PUT("/user/:id", app.updateUser)
	r.PUT("/user/:id/role", app.updateUserRole)
	r.POST("/roles", app.createRoles)
	r.DELETE("/user/:id", app.deleteUser)
	r.POST("/tokens/authentication", app.createAuthenticationToken)
	// Products and categories
	r.POST("/product", app.createProduct)
	r.POST("/category", app.createCategory)
	r.GET("/categories", app.getCategories)
	r.PUT("/category/:id", app.updateCategory)
	r.POST("/promotion", app.productPromotion)
	r.GET("/product/:barcode", app.getProduct)
	r.DELETE("/product/:barcode", app.deleteProduct)
	r.PUT("/product/:barcode", app.updateProduct)

	// Invoices
	r.POST("/invoice", app.createInvoice)

	// TODO: Delete route. Used for testing purposes only.
	{
		r.GET("/test", app.authenticate(), func(c *gin.Context) {
			user := app.contextGetUser(c)
			c.JSON(http.StatusOK, gin.H{"user": user})
		})
	}
}
