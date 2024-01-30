package handler

import (
	"database/sql"
	"fmt"
	"go-crud/models"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ErrorRespone struct{
	Message string `json:"message"`
}

type SuccesResponse struct {
	Message string `json:"message"`
}

type ArticleHandler struct {
	DB *sql.DB
}

func InitArticle(db *sql.DB) ArticleHandler {
	return ArticleHandler{
		DB: db,
	}
}


// Function GET data
func (h ArticleHandler) FetchArticles(c echo.Context) (err error) {
	datas := make([]models.Article, 0)
	query := `SELECT id,title,body FROM article`

	rows, err := h.DB.Query(query)
	if err != nil {
		res := ErrorRespone{
			Message: err.Error(),
		}
		return c.JSON(http.StatusInternalServerError, res)
	}
	defer rows.Close()

	for rows.Next(){
		var item models.Article
		err := rows.Scan(
			&item.ID,
			&item.Title,
			&item.Body,
		)

		if err != nil {
			res := ErrorRespone{
				Message: err.Error(),
			}
			return c.JSON(http.StatusInternalServerError, res)
		}

		datas = append(datas, item)
	}

	return c.JSON(http.StatusOK,datas)
}


// Function ADD Data
func (h ArticleHandler) Insert(c echo.Context) (err error)  {
	var item models.Article
	err = c.Bind(&item)
	if err != nil {
		res := ErrorRespone{
			Message: err.Error(),
		}
		return c.JSON(http.StatusUnprocessableEntity, res)
	}

	query := `INSERT INTO article (title, body) VALUES (?,?)`

	dbRes, err := h.DB.Exec(query, item.Title, item.Body)
	if err != nil {
		res := ErrorRespone{
			Message: err.Error(),
		}
		return c.JSON(http.StatusInternalServerError, res)
	}

	insertedID, err := dbRes.LastInsertId()
	if err != nil {
		res := ErrorRespone{
			Message: err.Error(),
		}
		return c.JSON(http.StatusInternalServerError, res)
	}

	item.ID = fmt.Sprintf("%d", insertedID)
	return c.JSON(http.StatusCreated, item)
}


// Function GET by id
func (h ArticleHandler) Get(c echo.Context) (err error)  {
	articleID := c.Param("id")

	query := `SELECT id, title, body FROM article WHERE id=?`
	row := h.DB.QueryRow(query, articleID)
	var res models.Article
	err = row.Scan(
		&res.ID,
		&res.Title,
		&res.Body,
	)
	if err != nil {
		res := ErrorRespone{
			Message: err.Error(),
		}
		if err == sql.ErrNoRows{
			return c.JSON(http.StatusNotFound, res)
		}
		return c.JSON(http.StatusInternalServerError, res)
	}
	return c.JSON(http.StatusCreated, res)
}


// Funtion DELETE by ID
func (h ArticleHandler) Delete(c echo.Context) (err error) {
	articleID := c.Param("id")

	// Buat query DELETE
	query := `DELETE FROM article WHERE id=?`
	result, err := h.DB.Exec(query, articleID)
	if err != nil {
		res := ErrorRespone{
			Message: err.Error(),
		}
		return c.JSON(http.StatusInternalServerError, res)
	}

	// Periksa jumlah baris yang terpengaruh untuk memastikan bahwa artikel dengan ID tersebut ada
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		res := ErrorRespone{
			Message: err.Error(),
		}
		return c.JSON(http.StatusInternalServerError, res)
	}

	if rowsAffected == 0 {
		res := ErrorRespone{
			Message: "Article not found",
		}
		return c.JSON(http.StatusNotFound, res)
	}

	respone := SuccesResponse{
		Message: "Article Succesfully deleted",
	}
	return c.JSON(http.StatusOK, respone)
}

// Function UPDATE by ID 
func (h ArticleHandler) Update(c echo.Context) (err error)  {
	articleID :=c.Param("id")
	var newData models.Article
	if err := c.Bind(&newData); err != nil{
		res := ErrorRespone{
			Message: err.Error(),
		}
		return c.JSON(http.StatusUnprocessableEntity,res)
	}

	// Query Update
	query := `UPDATE article SET title=?, body=? WHERE id=?`
	result, err := h.DB.Exec(query,newData.Title,newData.Body,articleID)
	if err != nil {
		res := ErrorRespone{
			Message: err.Error(),
		}
		return c.JSON(http.StatusInternalServerError, res)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		res := ErrorRespone{
			Message: err.Error(),
		}
		return c.JSON(http.StatusInternalServerError, res)
	}

	if rowsAffected == 0 {
		res := ErrorRespone{
			Message: "Article not found",
		}
		return c.JSON(http.StatusNotFound, res)
	}

	respone := SuccesResponse{
		Message: "Article Succesfully Updated",
	}
	return c.JSON(http.StatusOK, respone)

}
