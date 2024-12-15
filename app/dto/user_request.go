package dto

type UserPageRequest struct {
	PageRequest
	Username string `query:"username" form:"username"`
	Gender   *int   `query:"gender" form:"gender"`
	Enable   *int   `query:"enable" form:"enable"`
}

type UserAddRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Enable   bool   `json:"enable"`
	RoleIds  []int  `json:"roleIds"`
}

type UserUpdateRequest struct {
	Id       int    `json:"id"`
	Password string `json:"password"`
	Enable   bool   `json:"enable"`
	RoleIds  []int  `json:"roleIds"`
}

type UserProfileUpdateRequest struct {
	Id       int    `json:"id"`
	Nickname string `json:"nickName"`
	Gender   int    `json:"gender"`
	Avatar   string `json:"avatar"`
	Email    string `json:"email"`
	Address  string `json:"address"`
}