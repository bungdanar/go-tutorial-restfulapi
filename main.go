package main

import (
	"net/http"
	"tutorial-restfulapi/helper"
	"tutorial-restfulapi/middleware"

	_ "github.com/go-sql-driver/mysql"
)

func NewServer(authMiddleware *middleware.AuthMiddleware) *http.Server {
	return &http.Server{
		Addr:    "localhost:5500",
		Handler: authMiddleware,
	}
}

func main() {
	// db := app.NewDB()
	// validate := validator.New()

	// categoryRepository := repository.NewCategoryRepositoryImpl()
	// categoryService := service.NewCategoryServiceImpl(categoryRepository, db, validate)
	// categoryController := controller.NewCategoryControllerImpl(categoryService)
	// router := app.NewRouter(categoryController)
	// authMiddleware := middleware.NewAuthMiddleware(router)

	// server := NewServer(authMiddleware)

	// err := server.ListenAndServe()
	// helper.PanicIfError(err)

	server := InitializedServer()

	err := server.ListenAndServe()
	helper.PanicIfError(err)
}
