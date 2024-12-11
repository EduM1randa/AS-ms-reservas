package controllers_test

import (
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"gitlab.com/los-pelagatos-al-dia/pids-productos/controllers"
	db "gitlab.com/los-pelagatos-al-dia/pids-productos/database"
	"gitlab.com/los-pelagatos-al-dia/pids-productos/dto"
	"gitlab.com/los-pelagatos-al-dia/pids-productos/models"
)

func TestMain(m *testing.M) {
	fmt.Println("Ejecutando antes de los tests")
	db.ConnectMock()
	os.Exit(m.Run())
}

func TestNoProduct(t *testing.T) {
	var productDto dto.ProductDto
	productDto.ProductId = 0

	_, err := controllers.GetProduct(productDto)
	if err == nil {
		t.Errorf("Devolvio un producto")
	}
}

func TestInsertProduct(t *testing.T) {
	var categoryDto dto.CreateCategoryDto
	categoryDto.Name = "Categoria1"

	db.Mock.ExpectBegin()
	//
	query := regexp.QuoteMeta("INSERT INTO \"categories\" (\"created_at\",\"updated_at\",\"deleted_at\",\"name\")")
	rows := sqlmock.
		NewRows([]string{"id"}).
		AddRow(1)

	db.Mock.ExpectQuery(query).WithArgs(time.Now().Round(0), time.Now().Round(0), nil, categoryDto.Name).WillReturnRows(rows)
	db.Mock.ExpectCommit()

	responseCat, errCat1 := controllers.PostCategory(categoryDto)

	category := responseCat.Data.(models.Category)
	var catDto dto.CategoryDto
	catDto.CategoryId = strconv.Itoa(1)

	fmt.Println("Error: ", errCat1)
	log.Println(responseCat)

	fmt.Println("Id categoria: ", category.ID)

	rows2 := sqlmock.
		NewRows([]string{"id", "name"}).
		AddRow(category.ID, category.Name)

	query2 := regexp.QuoteMeta("SELECT id, name FROM \"categories\" WHERE id")

	db.Mock.ExpectQuery(query2).WithArgs(strconv.Itoa(1)).WillReturnRows(rows2)
	db.Mock.ExpectCommit()

	log.Println("Hol")

	responseCat2, errCat := controllers.GetCategory(catDto)

	fmt.Println("Encontre categoria", responseCat2)
	log.Println("Error busqueda categoria", errCat)

	var productDto dto.CreateProductDto
	productDto.Name = "Name"

	_, err := controllers.PostProduct(productDto)
	if err == nil {
		t.Errorf("Devolvio un producto")
	}

}

func TestPostCategoryError(t *testing.T) {
	var categoryDto dto.CreateCategoryDto
	categoryDto.Name = "Categoria1"

	db.Mock.ExpectBegin()

	query := regexp.QuoteMeta("INSERT INTO \"categories\" (\"created_at\",\"updated_at\",\"deleted_at\",\"name\")")

	db.Mock.ExpectQuery(query).WithArgs(time.Now().Round(0), time.Now().Round(0), nil, categoryDto.Name).WillReturnError(errors.New("error"))
	db.Mock.ExpectCommit()

	responseCat, errCat1 := controllers.PostCategory(categoryDto)

	fmt.Println("Error: ", errCat1)
	log.Println(responseCat)
}

func TestGetCategoryNotFound(t *testing.T) {
	var catDto dto.CategoryDto
	catDto.CategoryId = strconv.Itoa(1)

	query2 := regexp.QuoteMeta("SELECT id, name FROM \"categories\" WHERE id")

	db.Mock.ExpectBegin()
	db.Mock.ExpectQuery(query2).WithArgs(strconv.Itoa(1)).WillReturnError(errors.New("error"))
	db.Mock.ExpectCommit()

	responseCat2, errCat := controllers.GetCategory(catDto)

	fmt.Println("Error: ", errCat)
	log.Println(responseCat2)
}

func TestGetCategories(t *testing.T) {
	var categoriesDto dto.CategoriesDto
	categoriesDto.Limit = "15"

	db.Mock.ExpectBegin()
	// query := regexp.QuoteMeta("SELECT id, name FROM \"categories\" LIMIT \"limit\"")
	query := "SELECT id, name FROM categories WHERE categories.deleted_at IS NULL LIMIT 15"
	rows := sqlmock.
		NewRows([]string{"id", "name"}).
		AddRow("1", "perro")

	// db.Mock.ExpectQuery(query).WithArgs(strconv.Itoa(15)).WillReturnRows(rows)
	db.Mock.ExpectQuery(query).WithArgs().WillReturnRows(rows)
	db.Mock.ExpectCommit()

	responseCat, errCat := controllers.GetCategories(categoriesDto)

	fmt.Println("Error: ", errCat)
	log.Println(responseCat)
}

func TestGetCategory(t *testing.T) {
	var catDto dto.CategoryDto
	catDto.CategoryId = strconv.Itoa(1)

	query2 := regexp.QuoteMeta("SELECT id, name FROM \"categories\" WHERE id")

	db.Mock.ExpectBegin()
	db.Mock.ExpectQuery(query2).WithArgs(strconv.Itoa(1)).WillReturnError(errors.New("error"))
	db.Mock.ExpectCommit()

	responseCat2, errCat := controllers.GetCategory(catDto)

	fmt.Println("Error: ", errCat)
	log.Println(responseCat2)
}

func TestGetCategoriesError(t *testing.T) {
	var categoriesDto dto.CategoriesDto
	categoriesDto.Limit = "15"

	query := regexp.QuoteMeta("SELECT id, name FROM \"categories\" LIMIT limit")

	db.Mock.ExpectBegin()
	db.Mock.ExpectQuery(query).WithArgs(strconv.Itoa(15)).WillReturnError(errors.New("error"))
	db.Mock.ExpectCommit()

	responseCat2, errCat := controllers.GetCategories(categoriesDto)

	fmt.Println("Error: ", errCat)
	log.Println(responseCat2)
}

func TestObtenerLimit(t *testing.T) {
	limit := controllers.ObtenerLimit("20", 10)
	log.Println(limit)
}

func TestObtenerLimitUseDefault(t *testing.T) {
	limit := controllers.ObtenerLimit("-20", 10)
	log.Println(limit)
}

func TestGetTagsError(t *testing.T) {
	var tagsDto dto.TagsDto
	tagsDto.Limit = "15"

	db.Mock.ExpectBegin()
	query := regexp.QuoteMeta("SELECT id, name FROM \"tags\" LIMIT \"limit\"")
	// query := "SELECT id, name FROM categories WHERE categories.deleted_at IS NULL LIMIT 15"

	// db.Mock.ExpectQuery(query).WithArgs(strconv.Itoa(15)).WillReturnRows(rows)
	db.Mock.ExpectQuery(query).WithArgs("15").WillReturnError(errors.New("error"))
	db.Mock.ExpectCommit()

	response, err := controllers.GetTags(tagsDto)

	fmt.Println("Error: ", err)
	log.Println(response)
}

func ExampleHello() {
	fmt.Println("hello")
	// Output: hello
}

func ExampleSalutations() {
	fmt.Println("hello, and")
	fmt.Println("goodbye")
	// Output:
	// hello, and
	// goodbye
}
