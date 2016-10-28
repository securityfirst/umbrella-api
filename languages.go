package main

import "github.com/gin-gonic/gin"

func (um *Umbrella) getLanguages(c *gin.Context) {
	languages := []Language{
		Language{
			Name:  "English",
			Label: "en-gb",
		},
	}
	c.JSON(200, languages)
}
