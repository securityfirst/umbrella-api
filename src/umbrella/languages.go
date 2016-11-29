package main

import (
	"umbrella/models"
	"github.com/gin-gonic/gin"
)

func (um *Umbrella) getLanguages(c *gin.Context) {
	languages := []models.Language{
		{
			Name:  "English",
			Label: "en-gb",
		},
	}
	c.JSON(200, languages)
}
