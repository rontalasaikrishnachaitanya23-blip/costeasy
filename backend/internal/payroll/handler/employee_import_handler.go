package handler

import (
	"net/http"
	"strconv"

	"github.com/chaitu35/costeasy/backend/internal/payroll/imports"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type EmployeeImportHandler struct {
	importSvc *imports.EmployeeImportService
}

func NewEmployeeImportHandler(svc *imports.EmployeeImportService) *EmployeeImportHandler {
	return &EmployeeImportHandler{importSvc: svc}
}

// GET /payroll/employees/import/template
func (h *EmployeeImportHandler) DownloadTemplate(c *gin.Context) {
	f, err := imports.GenerateEmployeeTemplate()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create template"})
		return
	}
	defer f.Close()

	c.Header("Content-Disposition", "attachment; filename=employee_import_template.xlsx")
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	_ = f.Write(c.Writer)
}

// POST /payroll/employees/import
func (h *EmployeeImportHandler) Import(c *gin.Context) {
	// Assume orgID & userID from auth middleware/claims
	orgIDStr := c.GetHeader("X-Org-ID")
	userIDStr := c.GetHeader("X-User-ID")
	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid org id"})
		return
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing file"})
		return
	}
	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot open upload"})
		return
	}
	defer src.Close()

	res, err := h.importSvc.Import(c.Request.Context(), orgID, userID, file.Filename, src)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Header("X-Import-Total", strconv.Itoa(res.TotalRows))
	c.Header("X-Import-Valid", strconv.Itoa(res.ValidRows))
	c.Header("X-Import-Invalid", strconv.Itoa(res.InvalidRows))
	c.JSON(http.StatusOK, res)
}
