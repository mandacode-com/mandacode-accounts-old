package clienthandlerdto

import "mandacode.com/accounts/role/internal/domain/model"

type ClientData struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

type EnrollGroupUserRequest struct {
	*ClientData `json:"client_data"`
}
type EnrollGroupUserResponse struct {
	*model.GroupUser `json:"group_user"`
}

type GetAllGroupUsersRequest struct {
	*ClientData `json:"client_data"`
}
type GetAllGroupUsersResponse struct {
	GroupUsers []*model.GroupUser `json:"group_users"`
}

type CheckGroupUserRequest struct {
	*ClientData `json:"client_data"`
}
type CheckGroupUserResponse struct {
	IsEnrolled bool `json:"is_enrolled"`
}

type DeleteGroupUserRequest struct {
	*ClientData `json:"client_data"`
}
type DeleteGroupUserResponse struct {
	Success bool `json:"success"`
}

type DeleteGroupUserByGroupIDRequest struct {
	*ClientData `json:"client_data"`
}
type DeleteGroupUserByGroupIDResponse struct {
	Success bool `json:"success"`
}
