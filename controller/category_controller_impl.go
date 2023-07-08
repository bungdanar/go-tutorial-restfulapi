package controller

import (
	"net/http"
	"strconv"
	"tutorial-restfulapi/helper"
	"tutorial-restfulapi/model/web"
	"tutorial-restfulapi/service"

	"github.com/julienschmidt/httprouter"
)

type CategoryControllerImpl struct {
	CategoryService service.CategoryService
}

func NewCategoryController(categoryService service.CategoryService) CategoryController {
	return &CategoryControllerImpl{
		CategoryService: categoryService,
	}
}

func (controller *CategoryControllerImpl) Create(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	categoryCreateReq := web.CategoryCreateRequest{}
	helper.ReadFromReqBody(r, &categoryCreateReq)

	categoryRes := controller.CategoryService.Create(r.Context(), categoryCreateReq)
	webRes := web.WebResponse{
		Code:   201,
		Status: "Created",
		Data:   categoryRes,
	}

	helper.WriteToResBody(w, webRes)
}

func (controller *CategoryControllerImpl) Update(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	categoryUpdateReq := web.CategoryUpdateRequest{}
	helper.ReadFromReqBody(r, &categoryUpdateReq)

	categoryId, err := strconv.Atoi(p.ByName("categoryId"))
	helper.PanicIfError(err)

	categoryUpdateReq.Id = categoryId

	categoryRes := controller.CategoryService.Update(r.Context(), categoryUpdateReq)
	webRes := web.WebResponse{
		Code:   200,
		Status: "Ok",
		Data:   categoryRes,
	}

	helper.WriteToResBody(w, webRes)
}

func (controller *CategoryControllerImpl) Delete(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	categoryId, err := strconv.Atoi(p.ByName("categoryId"))
	helper.PanicIfError(err)

	controller.CategoryService.Delete(r.Context(), categoryId)
	webRes := web.WebResponse{
		Code:   200,
		Status: "Ok",
	}

	helper.WriteToResBody(w, webRes)
}

func (controller *CategoryControllerImpl) FindById(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	categoryId, err := strconv.Atoi(p.ByName("categoryId"))
	helper.PanicIfError(err)

	categoryRes := controller.CategoryService.FindById(r.Context(), categoryId)
	webRes := web.WebResponse{
		Code:   200,
		Status: "Ok",
		Data:   categoryRes,
	}

	helper.WriteToResBody(w, webRes)
}

func (controller *CategoryControllerImpl) FindAll(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	categoryRes := controller.CategoryService.FindAll(r.Context())
	webRes := web.WebResponse{
		Code:   200,
		Status: "Ok",
		Data:   categoryRes,
	}

	helper.WriteToResBody(w, webRes)
}
