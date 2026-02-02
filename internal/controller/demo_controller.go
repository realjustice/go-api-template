package controller

import (
	"strconv"

	"go-api-template/internal/model"
	"go-api-template/internal/service"
	"go-api-template/pkg/errors"
	"go-api-template/pkg/web"
)

// DemoController Demo 控制器
type DemoController struct {
	demoService *service.DemoService
}

// NewDemoController 创建 Demo Controller
func NewDemoController(demoService *service.DemoService) *DemoController {
	return &DemoController{
		demoService: demoService,
	}
}

// GetByID 根据 ID 获取
// @Summary 获取单个 Demo
// @Tags Demo
// @Param id path int true "Demo ID"
// @Success 200 {object} model.Demo
// @Router /api/v1/demos/{id} [get]
func (c *DemoController) GetByID(ctx *web.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		web.BadRequest(ctx, "invalid id")
		return
	}

	demo, err := c.demoService.GetByID(ctx.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			web.NotFound(ctx, "demo not found")
			return
		}
		web.InternalError(ctx, "get demo failed")
		return
	}

	web.Success(ctx, demo)
}

// GetAll 获取所有
// @Summary 获取所有 Demo
// @Tags Demo
// @Success 200 {array} model.Demo
// @Router /api/v1/demos [get]
func (c *DemoController) GetAll(ctx *web.Context) {
	demos, err := c.demoService.GetAll(ctx.Request.Context())
	if err != nil {
		web.InternalError(ctx, "get demos failed")
		return
	}

	web.Success(ctx, demos)
}

// CreateRequest 创建请求
type CreateRequest struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content"`
	Status  int    `json:"status"`
}

// Create 创建
// @Summary 创建 Demo
// @Tags Demo
// @Param request body CreateRequest true "创建参数"
// @Success 200 {object} model.Demo
// @Router /api/v1/demos [post]
func (c *DemoController) Create(ctx *web.Context) {
	var req CreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		web.BadRequest(ctx, "invalid request: "+err.Error())
		return
	}

	demo := &model.Demo{
		Title:   req.Title,
		Content: req.Content,
		Status:  req.Status,
	}

	err := c.demoService.Create(ctx.Request.Context(), demo)
	if err != nil {
		web.InternalError(ctx, "create demo failed")
		return
	}

	web.SuccessWithMessage(ctx, "demo created successfully", demo)
}

// UpdateRequest 更新请求
type UpdateRequest struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content"`
	Status  int    `json:"status"`
}

// Update 更新
// @Summary 更新 Demo
// @Tags Demo
// @Param id path int true "Demo ID"
// @Param request body UpdateRequest true "更新参数"
// @Success 200
// @Router /api/v1/demos/{id} [put]
func (c *DemoController) Update(ctx *web.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		web.BadRequest(ctx, "invalid id")
		return
	}

	var req UpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		web.BadRequest(ctx, "invalid request: "+err.Error())
		return
	}

	demo := &model.Demo{
		Title:   req.Title,
		Content: req.Content,
		Status:  req.Status,
	}

	err = c.demoService.Update(ctx.Request.Context(), uint(id), demo)
	if err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			web.NotFound(ctx, "demo not found")
			return
		}
		web.InternalError(ctx, "update demo failed")
		return
	}

	web.SuccessWithMessage(ctx, "demo updated successfully", nil)
}

// Delete 删除
// @Summary 删除 Demo
// @Tags Demo
// @Param id path int true "Demo ID"
// @Success 200
// @Router /api/v1/demos/{id} [delete]
func (c *DemoController) Delete(ctx *web.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		web.BadRequest(ctx, "invalid id")
		return
	}

	err = c.demoService.Delete(ctx.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			web.NotFound(ctx, "demo not found")
			return
		}
		web.InternalError(ctx, "delete demo failed")
		return
	}

	web.SuccessWithMessage(ctx, "demo deleted successfully", nil)
}
