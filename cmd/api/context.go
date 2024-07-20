package main

import (
	"github.com/gin-gonic/gin"
	"github.com/lordaris/pos-api/cmd/internal/data"
)

func (app *application) contextSetUser(c *gin.Context, user *data.User) {
	c.Set("user", user)
}

func (app *application) contextGetUser(c *gin.Context) *data.User {
	userInterface, exists := c.Get("user")
	if !exists {
		panic("missing user value in request context")
	}

	user, ok := userInterface.(*data.User)
	if !ok {
		panic("user value in context is not of type *data.User")
	}

	return user
}
