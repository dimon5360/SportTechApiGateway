package api

import (
	"app/main/grpc_service"
	"app/main/storage"
	"encoding/json"
	"log"
	"net/http"

	proto "github.com/dimon5360/SportTechProtos/gen/go/proto"
	"github.com/gin-gonic/gin"
)

func GetReport(c *gin.Context) {

	service := grpc_service.ReportServiceInstance()

	ID := c.Params.ByName("user_id")

	info := storage.Redis().Get(ID)

	var user UserInfo

	err := json.Unmarshal(info, &user)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"error": InvalidRequestArgs,
		})
		return
	}

	res, err := service.GetReport(&proto.GetReportRequest{
		UserId: user.Id,
	})

	c.JSON(http.StatusOK, gin.H{
		"user_id":    res.UserId,
		"document":   res.Report,
		"created_at": res.CreatedAt,
		"updated_at": res.UpdatedAt,
	})
}

func CreateReport(c *gin.Context) {

	service := grpc_service.ReportServiceInstance()

	type createUserRequest struct {
		UserId   string `json:"user_id"`
		Document string `json:"report"`
	}

	var req createUserRequest
	if err := c.Bind(&req); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"error": InvalidRequestArgs,
		})
		return
	}

	info := storage.Redis().Get(req.UserId)

	var user UserInfo
	if err := json.Unmarshal(info, &user); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"error": InvalidRequestArgs,
		})
		return
	}

	reportReq := proto.AddReportRequst{
		UserId: user.Id,
		Report: req.Document,
	}

	_, err := service.CreateReport(&reportReq)

	if err != nil {
		log.Printf("Creation report failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Creation report failed",
		})
		return
	}

	c.Status(http.StatusOK)
}
