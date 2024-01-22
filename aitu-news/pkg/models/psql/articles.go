package psql

import (
	"aitu.edu.kz/aitu-news/pkg/models"
	"database/sql"
	"errors"
	"strconv"
)

type ArticleModel struct {
	DB *sql.DB
}

func (m *ArticleModel) Insert(title, text, categoryId string) (int, error) {
	stmt := `INSERT INTO articles (title, text, category_id, publication_date, img_path)
		VALUES($1, $2, $3, NOW(), 'articleImg1.jpg')  RETURNING id`
	insertedId := 0
	err := m.DB.QueryRow(stmt, title, text, categoryId).Scan(&insertedId)
	if err != nil {
		return 0, err
	}
	return int(insertedId), nil
}

func (m *ArticleModel) Get(id int) (*models.Article, error) {
	a := &models.Article{}
	stmt := `SELECT id, title, text, category_id, publication_date, img_path FROM articles WHERE id = $1`
	row := m.DB.QueryRow(stmt, id)
	err := row.Scan(&a.ID, &a.Title, &a.Text, &a.CategoryId, &a.PublicationDate, &a.ImgPath)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}
	return a, nil
}

func (m *ArticleModel) Latest() ([]*models.Article, error) {
	stmt := `SELECT * FROM articles ORDER BY publication_date DESC limit 5`
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	articles := []*models.Article{}

	for rows.Next() {
		article := &models.Article{}
		err = rows.Scan(&article.ID, &article.Title, &article.Text, &article.CategoryId, &article.PublicationDate, &article.ImgPath)
		article.FormattedDate = article.PublicationDate.Format("2006-01-02")
		if err != nil {
			return nil, err
		}
		articles = append(articles, article)
	}
	return articles, nil
}

func (m *ArticleModel) GetAll(categoryId int) ([]*models.Article, error) {
	stmt := `SELECT * FROM articles ORDER BY publication_date DESC`
	if categoryId > 0 && categoryId < 5 {
		stmt = `SELECT * FROM articles WHERE category_id = ` + strconv.Itoa(categoryId) + ` ORDER BY publication_date DESC`
	}
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	articles := []*models.Article{}

	for rows.Next() {
		article := &models.Article{}
		err = rows.Scan(&article.ID, &article.Title, &article.Text, &article.CategoryId, &article.PublicationDate, &article.ImgPath)
		article.FormattedDate = article.PublicationDate.Format("2006-01-02")
		if err != nil {
			return nil, err
		}
		articles = append(articles, article)
	}
	return articles, nil
}
