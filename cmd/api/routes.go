package main

import (
	"github.com/gin-gonic/gin"
)

func Router(r *gin.Engine, app *application) {
	// Create a base router group that applies authentication and logging to all routes
	base := r.Group("/")
	base.Use(app.authenticate())

	// Users and roles routes
	users := base.Group("/user")
	{
		users.POST("", app.allowRole("manager", "admin", "hr"), app.loggerMiddleware(), app.createUser)             // Create a new user
		users.GET("/:id", app.allowRole("manager"), app.getUser)                                                    // Get a specific user by ID
		users.GET("/role/:role", app.getUsersByRole)                                                                // Get users by role
		users.PUT("/:id", app.allowRole("manager", "admin", "hr"), app.loggerMiddleware(), app.updateUser)          // Update a user's information
		users.PUT("/:id/role", app.allowRole("manager", "admin", "hr"), app.loggerMiddleware(), app.updateUserRole) // Update a user's role
		users.DELETE("/:id", app.allowRole("manager", "admin", "hr"), app.loggerMiddleware(), app.deleteUser)       // Delete a user
	}

	// Roles routes
	roles := base.Group("/roles")
	{
		roles.POST("", app.allowRole("manager", "admin", "hr"), app.loggerMiddleware(), app.createRoles) // Create new roles
		roles.GET("", app.getAllRoles)                                                                   // Get all roles
	}

	// Authentication routes
	r.POST("/auth/token", app.createAuthenticationToken) // Create an stateless authentication token

	// Products routes
	products := base.Group("/products")
	{
		products.POST("", app.allowRole("admin", "manager", "price_checker"), app.loggerMiddleware(), app.createProduct)            // Create a new product
		products.PUT("/:barcode", app.allowRole("admin", "manager", "price_checker"), app.loggerMiddleware(), app.updateProduct)    // Update a product
		products.DELETE("/:barcode", app.allowRole("admin", "manager", "price_checker"), app.loggerMiddleware(), app.deleteProduct) // Delete a product
	}

	r.GET("/products/:barcode", app.getProduct) // Get a product by barcode

	// Categories routes
	categories := base.Group("/categories")
	{
		categories.POST("", app.allowRole("admin", "manager", "price_checker"), app.loggerMiddleware(), app.createCategory)       // Create a new category
		categories.PUT("/:id", app.allowRole("admin", "manager", "price_checker"), app.loggerMiddleware(), app.updateCategory)    // Update a category
		categories.DELETE("/:id", app.allowRole("admin", "manager", "price_checker"), app.loggerMiddleware(), app.deleteCategory) // Delete a category
	}

	r.GET("/categories", app.getCategories) // Get all categories

	// Promotions routes
	promotions := base.Group("/promotions")
	{
		promotions.POST("", app.allowRole("admin", "manager", "price_checker"), app.loggerMiddleware(), app.productPromotion)        // Create a new promotion
		promotions.PUT("/:barcode", app.allowRole("admin", "manager", "price_checker"), app.loggerMiddleware(), app.updatePromotion) // Update a promotion
	}

	r.GET("/promotions/:barcode", app.getPromotion) // Get a promotion by product barcode

	// Invoices routes
	invoices := base.Group("/invoices")
	{
		invoices.POST("", app.createInvoice)           // Create a new invoice
		invoices.GET("/:ticketnumber", app.getInvoice) // Get an invoice by ticket number
	}

	// Inventory routes
	inventory := base.Group("/inventory")
	{
		// Apply additional middleware for admin/manager roles
		inventory.Use(app.allowRole("admin", "manager"))
		inventory.POST("", app.createInventoryMovement) // Create an inventory movement
	}

	r.GET("/sales/:startdate/:enddate", app.getSales)
}
