package models

import (
	"fmt"
	"log"
	proto "proto/go"
)

type RestLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RestLoginResponse struct {
	Id           uint64
	RefrestToken Token
	AccessToken  Token
	ProfileId    uint64
	IsValidated  bool
	Error        error
}

func handlerAuthError(err proto.AuthError) error {
	switch err {
	case proto.AuthError_OK:
		log.Println("success error code")
		return nil
	case proto.AuthError_NOT_FOUND:
		return fmt.Errorf("user not found")
	case proto.AuthError_ALREADY_EXIST:
		return fmt.Errorf("user already exists")
	default:
		log.Println("unhandled error code")
		return nil
	}
}

func ConvertRest2GrpcLoginRequest(req *RestLoginRequest) *proto.LoginUserRequest {
	log.Println("rest login request:", req)
	return &proto.LoginUserRequest{
		Email:    req.Email,
		Password: req.Password,
	}
}

func ConvertGrpc2RestLoginResponse(resp *proto.LoginUserResponse) *RestLoginResponse {
	log.Println("protobuf login response:", resp)
	return &RestLoginResponse{
		Id: resp.Id,
		AccessToken: Token{
			value: resp.AccessToken.GetValue(),
			age:   int(resp.AccessToken.GetAge()),
		},
		RefrestToken: Token{
			value: resp.RefreshToken.GetValue(),
			age:   int(resp.RefreshToken.GetAge()),
		},
		ProfileId:   resp.ProfileId,
		IsValidated: resp.IsValidated,
		Error:       handlerAuthError(resp.Error),
	}
}

type RestRegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func ConvertRest2GrpcRegisterRequest(req *RestRegisterRequest) *proto.RegisterUserRequest {
	log.Println("rest register request:", req)
	return &proto.RegisterUserRequest{
		Email:    req.Email,
		Password: req.Password,
	}
}

func ConvertGrpc2RestRegisterResponse(resp *proto.RegisterUserResponse) error {
	log.Println("protobuf register response:", resp)
	return handlerAuthError(resp.Error)
}
