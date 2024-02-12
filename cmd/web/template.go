package main

import (
	models "asik1/pkg"
)

type TemplateData struct {
	Flash           string
	OneNews         *models.News
	News            []*models.News
	User            *models.User
	IsAuthenticated bool
}
