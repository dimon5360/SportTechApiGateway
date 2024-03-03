package endpoint

import (
	"app/main/internal/endpoint"
	"app/main/internal/repository"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type authEndpoint struct {
	user  repository.Interface
	token repository.Interface
}

func New(repo ...repository.Interface) endpoint.Interface {
	if len(repo) != 2 {
		log.Fatalln()
		return nil
	}

	e := &authEndpoint{
		user:  repo[0],
		token: repo[1],
	}

	if err := e.user.Init(); err != nil {
		log.Fatal(err)
	}

	if err := e.token.Init(); err != nil {
		log.Fatal(err)
	}

	return e
}

func (e *authEndpoint) Get(c *gin.Context) {

	log.Println("unimplemented method")
	c.Status(http.StatusOK)
}

func (e *authEndpoint) Post(c *gin.Context) {

	log.Println("unimplemented method")
	c.Status(http.StatusOK)
}

// func ProcessingFailed(c *gin.Context, err error, message string, status int) {

// 	log.Println(err.Error())

// 	c.JSON(status, gin.H{
// 		"error": api.InvalidRequestArgs,
// 	})
// }

// func generateToken(id uint64) (string, error) {

// 	payload := jwt.MapClaims{
// 		"sub": id,
// 		"iat": time.Now().Add(time.Hour * 72).Unix(),
// 	}

// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
// 	var jwtSecretKey = []byte("very-secret-key")

// 	t, err := token.SignedString(jwtSecretKey)
// 	if err != nil {
// 		return "", err
// 	}

// 	return t, nil
// }

// func AuthenticateUser(c *gin.Context) {

// 	service := grpc_service.AuthServiceInstance()

// 	type authUserRequest struct {
// 		Email    string `json:"email"`
// 		Password string `json:"password"`
// 	}

// 	var req authUserRequest
// 	err := c.Bind(&req)
// 	if err != nil {
// 		ProcessingFailed(c, err, api.InvalidRequestArgs, http.StatusBadRequest)
// 		return
// 	}

// 	fmt.Printf("%s: %v\n", "Auth user request", req)

// 	res, err := service.Auth(&proto.AuthUserRequest{
// 		Email:    req.Email,
// 		Password: req.Password,
// 	})

// 	if err != nil {
// 		ProcessingFailed(c, err, "Authentication failed", http.StatusUnauthorized)
// 		return
// 	}

// 	var id uint64 = res.Id
// 	token, err := generateToken(id)
// 	if err != nil {
// 		ProcessingFailed(c, err, "JWT token signing failed", http.StatusInternalServerError)
// 		return
// 	}

// 	expireIn := time.Hour * 24

// 	c.SetCookie("access_token", token, int(time.Now().Add(expireIn).Unix()), "", "", true, false)
// 	ck, err := c.Cookie("access_token")
// 	if err != nil {
// 		log.Println(ck)
// 	}

// 	bytes, err := json.Marshal(model.UserInfo{
// 		Id:          id,
// 		AccessToken: token,
// 	})

// 	if err != nil {
// 		ProcessingFailed(c, err, "JWT token handling failed", http.StatusInternalServerError)
// 		return
// 	}

// 	storage.Redis().Store(fmt.Sprintf("%d", id), bytes, expireIn)

// 	if err = profile.VerifyProfile(id); err != nil {
// 		c.Redirect(http.StatusFound, "/profile/create")
// 	}

// 	userInfo := gin.H{
// 		"user_id":       id,
// 		"refresh_token": "default refresh token",
// 		"access_token":  token,
// 	}

// 	c.JSON(http.StatusOK, userInfo)
// }
