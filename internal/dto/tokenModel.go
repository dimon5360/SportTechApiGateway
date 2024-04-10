package dto

import (
	"log"
	proto "proto/go"
)

type Token struct {
	value string
	age   int
}

func (t *Token) GetValue() string {
	return t.value
}

func (t *Token) GetAge() int {
	return t.age
}

type RestRefreshTokenRequest struct {
	Id           uint64
	RefreshToken string
}

type RestRefreshTokenResponse struct {
	Id           uint64
	RefrestToken Token
	AccessToken  Token
	Error        error
}

func ConvertRest2GrpcRefreshRequest(req *RestRefreshTokenRequest) *proto.RefreshTokenRequest {
	log.Println("rest refresh token request:", req)
	return &proto.RefreshTokenRequest{
		Id:           req.Id,
		RefreshToken: req.RefreshToken,
	}
}

func ConvertGrpc2RestRefreshnResponse(resp *proto.RefreshTokenResponse) *RestRefreshTokenResponse {
	log.Println("protobuf refresh token response:", resp)
	return &RestRefreshTokenResponse{
		Id: resp.Id,
		AccessToken: Token{
			value: resp.AccessToken.GetValue(),
			age:   int(resp.AccessToken.GetAge()),
		},
		RefrestToken: Token{
			value: resp.RefreshToken.GetValue(),
			age:   int(resp.RefreshToken.GetAge()),
		},
		Error: handlerAuthError(resp.Error),
	}
}
