package test

import (
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"
	"tutorial-restfulapi/app"
	"tutorial-restfulapi/controller"
	"tutorial-restfulapi/helper"
	"tutorial-restfulapi/middleware"
	"tutorial-restfulapi/model/domain"
	"tutorial-restfulapi/repository"
	"tutorial-restfulapi/service"

	"github.com/go-playground/validator/v10"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

func setupTestDB() *sql.DB {
	db, err := sql.Open("mysql", "root:brightshield!23@tcp(localhost:3306)/go-tutorial-restfulapi")
	helper.PanicIfError(err)

	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(20)
	db.SetConnMaxIdleTime(10 * time.Minute)
	db.SetConnMaxLifetime(60 * time.Minute)

	return db
}

func setupRouter(db *sql.DB) http.Handler {
	validate := validator.New()

	categoryRepository := repository.NewCategoryRepositoryImpl()
	categoryService := service.NewCategoryServiceImpl(categoryRepository, db, validate)
	categoryController := controller.NewCategoryControllerImpl(categoryService)

	router := app.NewRouter(categoryController)

	return middleware.NewAuthMiddleware(router)
}

func truncateCategory(db *sql.DB) {
	db.Exec("TRUNCATE category")
}

func TestCreateCategorySuccess(t *testing.T) {
	db := setupTestDB()
	truncateCategory(db)
	router := setupRouter(db)

	reqBody := strings.NewReader(`{"name": "Gadget"}`)
	req := httptest.NewRequest(http.MethodPost, "http://localhost:5500/api/categories", reqBody)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-API-Key", "mySecureKey")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	res := recorder.Result()
	assert.Equal(t, 201, res.StatusCode)

	body, _ := io.ReadAll(res.Body)
	var data map[string]interface{}
	json.Unmarshal(body, &data)

	assert.Equal(t, 201, int(data["code"].(float64)))
	assert.Equal(t, "Created", data["status"])
	assert.Equal(t, "Gadget", data["data"].(map[string]interface{})["name"])
}

func TestCreateCategoryFailed(t *testing.T) {
	db := setupTestDB()
	truncateCategory(db)
	router := setupRouter(db)

	reqBody := strings.NewReader(`{"name": ""}`)
	req := httptest.NewRequest(http.MethodPost, "http://localhost:5500/api/categories", reqBody)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-API-Key", "mySecureKey")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	res := recorder.Result()
	assert.Equal(t, 400, res.StatusCode)

	body, _ := io.ReadAll(res.Body)
	var data map[string]interface{}
	json.Unmarshal(body, &data)

	assert.Equal(t, 400, int(data["code"].(float64)))
	assert.Equal(t, "Bad Request", data["status"])
}

func TestUpdateCategorySuccess(t *testing.T) {
	db := setupTestDB()
	truncateCategory(db)

	tx, _ := db.Begin()
	categoryRepository := repository.NewCategoryRepositoryImpl()
	category := categoryRepository.Save(context.Background(), tx, domain.Category{
		Name: "Gadget",
	})
	tx.Commit()

	router := setupRouter(db)

	reqBody := strings.NewReader(`{"name": "Gadget and Style"}`)
	req := httptest.NewRequest(http.MethodPut, "http://localhost:5500/api/categories/"+strconv.Itoa(category.Id), reqBody)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-API-Key", "mySecureKey")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	res := recorder.Result()
	assert.Equal(t, 200, res.StatusCode)

	body, _ := io.ReadAll(res.Body)
	var data map[string]interface{}
	json.Unmarshal(body, &data)

	assert.Equal(t, 200, int(data["code"].(float64)))
	assert.Equal(t, "Ok", data["status"])
	assert.Equal(t, category.Id, int(data["data"].(map[string]interface{})["id"].(float64)))
	assert.Equal(t, "Gadget and Style", data["data"].(map[string]interface{})["name"])
}

func TestUpdateCategoryFailed(t *testing.T) {
	db := setupTestDB()
	truncateCategory(db)

	tx, _ := db.Begin()
	categoryRepository := repository.NewCategoryRepositoryImpl()
	category := categoryRepository.Save(context.Background(), tx, domain.Category{
		Name: "Gadget",
	})
	tx.Commit()

	router := setupRouter(db)

	reqBody := strings.NewReader(`{"name": ""}`)
	req := httptest.NewRequest(http.MethodPut, "http://localhost:5500/api/categories/"+strconv.Itoa(category.Id), reqBody)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-API-Key", "mySecureKey")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	res := recorder.Result()
	assert.Equal(t, 400, res.StatusCode)

	body, _ := io.ReadAll(res.Body)
	var data map[string]interface{}
	json.Unmarshal(body, &data)

	assert.Equal(t, 400, int(data["code"].(float64)))
	assert.Equal(t, "Bad Request", data["status"])
}

func TestGetCategorySuccess(t *testing.T) {
	db := setupTestDB()
	truncateCategory(db)

	tx, _ := db.Begin()
	categoryRepository := repository.NewCategoryRepositoryImpl()
	category := categoryRepository.Save(context.Background(), tx, domain.Category{
		Name: "Gadget",
	})
	tx.Commit()

	router := setupRouter(db)

	req := httptest.NewRequest(http.MethodGet, "http://localhost:5500/api/categories/"+strconv.Itoa(category.Id), nil)
	req.Header.Add("X-API-Key", "mySecureKey")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	res := recorder.Result()
	assert.Equal(t, 200, res.StatusCode)

	body, _ := io.ReadAll(res.Body)
	var data map[string]interface{}
	json.Unmarshal(body, &data)

	assert.Equal(t, 200, int(data["code"].(float64)))
	assert.Equal(t, "Ok", data["status"])
	assert.Equal(t, category.Id, int(data["data"].(map[string]interface{})["id"].(float64)))
	assert.Equal(t, category.Name, data["data"].(map[string]interface{})["name"])
}

func TestGetCategoryFailed(t *testing.T) {
	db := setupTestDB()
	truncateCategory(db)

	router := setupRouter(db)

	req := httptest.NewRequest(http.MethodGet, "http://localhost:5500/api/categories/1", nil)
	req.Header.Add("X-API-Key", "mySecureKey")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	res := recorder.Result()
	assert.Equal(t, 404, res.StatusCode)

	body, _ := io.ReadAll(res.Body)
	var data map[string]interface{}
	json.Unmarshal(body, &data)

	assert.Equal(t, 404, int(data["code"].(float64)))
	assert.Equal(t, "Not Found", data["status"])
}

func TestDeleteCategorySuccess(t *testing.T) {
	db := setupTestDB()
	truncateCategory(db)

	tx, _ := db.Begin()
	categoryRepository := repository.NewCategoryRepositoryImpl()
	category := categoryRepository.Save(context.Background(), tx, domain.Category{
		Name: "Gadget",
	})
	tx.Commit()

	router := setupRouter(db)

	req := httptest.NewRequest(http.MethodDelete, "http://localhost:5500/api/categories/"+strconv.Itoa(category.Id), nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-API-Key", "mySecureKey")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	res := recorder.Result()
	assert.Equal(t, 200, res.StatusCode)

	body, _ := io.ReadAll(res.Body)
	var data map[string]interface{}
	json.Unmarshal(body, &data)

	assert.Equal(t, 200, int(data["code"].(float64)))
	assert.Equal(t, "Ok", data["status"])
}

func TestDeleteCategoryFailed(t *testing.T) {
	db := setupTestDB()
	truncateCategory(db)

	router := setupRouter(db)

	req := httptest.NewRequest(http.MethodDelete, "http://localhost:5500/api/categories/1", nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-API-Key", "mySecureKey")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	res := recorder.Result()
	assert.Equal(t, 404, res.StatusCode)

	body, _ := io.ReadAll(res.Body)
	var data map[string]interface{}
	json.Unmarshal(body, &data)

	assert.Equal(t, 404, int(data["code"].(float64)))
	assert.Equal(t, "Not Found", data["status"])
}

func TestListCategoriesSuccess(t *testing.T) {
	db := setupTestDB()
	truncateCategory(db)

	tx, _ := db.Begin()
	categoryRepository := repository.NewCategoryRepositoryImpl()
	category1 := categoryRepository.Save(context.Background(), tx, domain.Category{
		Name: "Gadget",
	})
	category2 := categoryRepository.Save(context.Background(), tx, domain.Category{
		Name: "Computer",
	})
	tx.Commit()

	router := setupRouter(db)

	req := httptest.NewRequest(http.MethodGet, "http://localhost:5500/api/categories", nil)
	req.Header.Add("X-API-Key", "mySecureKey")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	res := recorder.Result()
	assert.Equal(t, 200, res.StatusCode)

	body, _ := io.ReadAll(res.Body)
	var data map[string]interface{}
	json.Unmarshal(body, &data)

	assert.Equal(t, 200, int(data["code"].(float64)))
	assert.Equal(t, "Ok", data["status"])

	var categories = data["data"].([]interface{})
	categoryRes1 := categories[0].(map[string]interface{})
	categoryRes2 := categories[1].(map[string]interface{})

	assert.Equal(t, category1.Id, int(categoryRes1["id"].(float64)))
	assert.Equal(t, category1.Name, categoryRes1["name"])

	assert.Equal(t, category2.Id, int(categoryRes2["id"].(float64)))
	assert.Equal(t, category2.Name, categoryRes2["name"])
}

func TestUnauthorized(t *testing.T) {
	db := setupTestDB()
	truncateCategory(db)

	router := setupRouter(db)

	req := httptest.NewRequest(http.MethodGet, "http://localhost:5500/api/categories", nil)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	res := recorder.Result()
	assert.Equal(t, 401, res.StatusCode)

	body, _ := io.ReadAll(res.Body)
	var data map[string]interface{}
	json.Unmarshal(body, &data)

	assert.Equal(t, 401, int(data["code"].(float64)))
	assert.Equal(t, "Unauthorized", data["status"])
}
