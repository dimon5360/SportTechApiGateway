package endpoint

import (
	constants "app/main/internal/endpoint/common"
	model "app/main/internal/model"
	repository "app/main/internal/repository"
	"app/main/internal/storage"
	"encoding/json"
	"log"
	"net/http"

	proto "github.com/dimon5360/SportTechProtos/gen/go/proto"
	"github.com/gin-gonic/gin"
)

type reportEndpoint struct {
	repo repository.Interface
}

func NewReportEndpoint(repo repository.Interface) *reportEndpoint {
	return &reportEndpoint{
		repo: repo,
	}
}

func (e *reportEndpoint) Get(c *gin.Context) {

	ID := c.Params.ByName("user_id")

	info := storage.Redis().Get(ID)

	var user model.UserInfo

	err := json.Unmarshal(info, &user)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"error": constants.InvalidRequestArgs,
		})
		return
	}

	res, err := e.repo.Get(&proto.GetReportRequest{
		UserId: user.Id,
	})

	if err != nil {
		log.Printf("Creation report failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Creation report failed",
		})
		return
	}

	if val, ok := res.(*proto.ReportResponse); ok {
		c.JSON(http.StatusOK, gin.H{
			"user_id":    val.UserId,
			"document":   val.Report,
			"created_at": val.CreatedAt,
			"updated_at": val.UpdatedAt,
		})
		return
	}

	c.Status(http.StatusInternalServerError)
	log.Println("invalid repository response")
}

func (e *reportEndpoint) Post(c *gin.Context) {

	type createUserRequest struct {
		UserId   string `json:"user_id"`
		Document string `json:"report"`
	}

	var req createUserRequest
	if err := c.Bind(&req); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"error": constants.InvalidRequestArgs,
		})
		return
	}

	info := storage.Redis().Get(req.UserId)

	var user model.UserInfo
	if err := json.Unmarshal(info, &user); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"error": constants.InvalidRequestArgs,
		})
		return
	}

	reportReq := proto.AddReportRequst{
		UserId: user.Id,
		Report: req.Document,
	}

	res, err := e.repo.Add(&reportReq)

	if err != nil {
		log.Printf("Creation report failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Creation report failed",
		})
		return
	}

	if val, ok := res.(*proto.ReportResponse); ok {
		c.JSON(http.StatusOK, gin.H{
			"user_id": val.UserId, // TODO: #1 further replace to report ID
		})

		c.Redirect(http.StatusFound, "/profile/get/"+req.UserId)
		return
	}

	c.Status(http.StatusInternalServerError)
	log.Println("invalid repository response")
}
