package main

import "aitu.edu.kz/aitu-news/pkg/models"

type templateData struct {
	Articles []*models.Article
	User     *models.User
}
