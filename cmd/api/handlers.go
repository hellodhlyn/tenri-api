package main

import (
	"github.com/gin-gonic/gin"
	"github.com/hellodhlyn/tenri-api/internal/models"
	"github.com/hellodhlyn/tenri-api/internal/utils"
	"net/http"
	"time"
)

// GET /q/v1/tasks
func getTasks(serverCtx serverContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		tasks, err := models.GetTasks(c.Request.Context(), serverCtx.rdb)
		if err != nil {
			c.Error(err)
			c.JSON(http.StatusInternalServerError, nil)
			return
		}
		c.JSON(http.StatusOK, tasks)
	}
}

// POST /q/v1/tasks
type postTaskReq struct {
	Text  string `json:"text"`
	DueAt string `json:"dueAt"`
}

func postTask(serverCtx serverContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req postTaskReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, nil)
			return
		}

		dueAt, err := time.Parse(time.RFC3339, req.DueAt)
		if err != nil {
			c.JSON(http.StatusBadRequest, nil)
			return
		}

		task := models.NewTask(req.Text, dueAt)
		err = models.SaveTask(c.Request.Context(), serverCtx.rdb, task)
		if err != nil {
			c.Error(err)
			c.JSON(http.StatusInternalServerError, nil)
			return
		}

		c.JSON(http.StatusCreated, task)
	}
}

// PATCH /q/v1/tasks/:uuid
type patchTaskReq struct {
	Text  string `json:"text"`
	DueAt string `json:"dueAt"`
}

func patchTask(serverCtx serverContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req patchTaskReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, nil)
			return
		}

		tasks, err := models.GetTasks(c.Request.Context(), serverCtx.rdb)
		if err != nil {
			c.Error(err)
			c.JSON(http.StatusInternalServerError, nil)
			return
		}
		task, ok := utils.FindFirstSlice(tasks, func(task models.Task) bool { return task.UUID == c.Param("uuid") })
		if !ok {
			c.JSON(http.StatusNotFound, nil)
			return
		}

		if req.Text != "" {
			task.Text = req.Text
		}
		if req.DueAt != "" {
			dueAt, err := time.Parse(time.RFC3339, req.DueAt)
			if err != nil {
				c.JSON(http.StatusBadRequest, nil)
				return
			}
			task.DueAt = dueAt
		}

		err = models.SaveTask(c.Request.Context(), serverCtx.rdb, task)
		if err != nil {
			c.Error(err)
			c.JSON(http.StatusInternalServerError, nil)
			return
		}
		c.JSON(http.StatusOK, task)
	}
}

// DELETE /q/v1/tasks/:uuid
func deleteTask(serverCtx serverContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		tasks, err := models.GetTasks(c.Request.Context(), serverCtx.rdb)
		if err != nil {
			c.Error(err)
			c.JSON(http.StatusInternalServerError, nil)
			return
		}
		task, ok := utils.FindFirstSlice(tasks, func(task models.Task) bool { return task.UUID == c.Param("uuid") })
		if !ok {
			c.JSON(http.StatusNotFound, nil)
			return
		}

		err = models.DeleteTask(c.Request.Context(), serverCtx.rdb, c.Param("uuid"))
		c.JSON(http.StatusOK, task)
	}
}
