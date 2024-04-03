package endpoint

import (
	"app/main/internal/endpoint"
	repository "app/main/internal/repository"
	"log"
	"net/http"
	"strconv"

	proto "proto/go"

	"github.com/gin-gonic/gin"
)

type reportEndpoint struct {
	repo repository.ReportInterface
}

func New(reportRepository repository.ReportInterface) (endpoint.Report, error) {
	e := &reportEndpoint{
		repo: reportRepository,
	}

	if err := e.repo.Init(); err != nil {
		return nil, err
	}
	return e, nil
}

func (e *reportEndpoint) Get(c *gin.Context) {

	ID := c.Params.ByName("user_id")
	userId, err := strconv.ParseUint(ID, 10, 64)
	if err != nil {
		log.Println("Invalid conversion from string to uint64")
		c.Status(http.StatusInternalServerError)
		return
	}

	res, err := e.repo.Read(&proto.GetReportRequest{
		UserId: userId,
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
			"error": endpoint.InvalidRequestArgs,
		})
		return
	}

	userId, err := strconv.ParseUint(req.UserId, 10, 64)
	if err != nil {
		log.Println("Invalid conversion from string to uint64")
		c.Status(http.StatusInternalServerError)
		return
	}

	reportReq := proto.AddReportRequst{
		UserId: userId,
		Report: req.Document,
	}

	res, err := e.repo.Create(&reportReq)

	if err != nil {
		log.Printf("Creation report failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Creation report failed",
		})
		return
	}

	if val, ok := res.(*proto.ReportResponse); ok {
		c.JSON(http.StatusOK, gin.H{
			"report_id": val.UserId,
		})

		c.Redirect(http.StatusFound, "/profile/"+req.UserId)
		return
	}

	c.Status(http.StatusInternalServerError)
	log.Println("invalid repository response")
}
