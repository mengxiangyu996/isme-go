package controller

import (
	"isme-go/app/dto"
	"isme-go/app/service"
	"isme-go/common/password"
	"isme-go/framework/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

type User struct{}

// 获取当前登录用户的详情信息
func (*User) Detail(ctx *gin.Context) {

	userId := ctx.GetInt("userId")
	roleCode := ctx.GetString("roleCode")

	user := (&service.User{}).GetDetailById(userId)

	userResponse := dto.UserDetailResponse{
		UserResponse: user,
	}

	roleIds := (&service.UserRolesRole{}).GetRoleIdsByUserId(user.Id)
	roles := (&service.Role{}).GetListByIds(roleIds, true)

	userResponse.Roles = roles

	profile := (&service.Profile{}).GetDetailByUserId(user.Id)
	userResponse.Profile = profile

	if len(roles) > 0 {
		for _, role := range roles {
			if role.Code == roleCode {
				userResponse.CurrentRole = role
				break
			}
		}
	}

	response.NewSuccess().SetData("data", userResponse).Json(ctx)
}

// 获取用户列表
func (*User) Page(ctx *gin.Context) {

	var param dto.UserPageRequest

	if err := ctx.BindQuery(&param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	userPages := make([]dto.UserPageResponse, 0)

	users, count := (&service.User{}).Page(param)
	if len(users) > 0 {
		for _, user := range users {
			var userPage dto.UserPageResponse
			userPage.UserResponse = user
			profile := (&service.Profile{}).GetDetailByUserId(user.Id)
			userPage.Gender = profile.Gender
			userPage.Avatar = profile.Avatar
			userPage.Address = profile.Address
			userPage.Email = profile.Email
			roleIds := (&service.UserRolesRole{}).GetRoleIdsByUserId(user.Id)
			roles := (&service.Role{}).GetListByIds(roleIds, true)
			userPage.Roles = roles
			userPages = append(userPages, userPage)
		}
	}

	response.NewSuccess().SetData("data", map[string]interface{}{
		"pageData": userPages,
		"total":    count,
	}).Json(ctx)
}

// 删除用户
func (*User) Delete(ctx *gin.Context) {

	id, _ := strconv.Atoi(ctx.Param("id"))

	if id == 1 {
		response.NewError().SetMsg("不能删除根用户").Json(ctx)
		return
	}

	if id == ctx.GetInt("userId") {
		response.NewError().SetMsg("非法操作，不能删除自己").Json(ctx)
		return
	}

	if err := (&service.User{}).Delete(id); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	response.NewSuccess().Json(ctx)
}

// 添加用户
func (*User) Add(ctx *gin.Context) {

	var param dto.UserAddRequest

	if err := ctx.Bind(&param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	if param.Username == "" {
		response.NewError().SetMsg("用户名不能为空").Json(ctx)
		return
	}

	if param.Password == "" {
		response.NewError().SetMsg("密码不能为空").Json(ctx)
		return
	}

	user := (&service.User{}).GetDetailByUsername(param.Username)
	if user.Id > 0 {
		response.NewError().SetMsg("用户名已存在").Json(ctx)
		return
	}

	param.Password = password.Generate(param.Password)

	if err := (&service.User{}).Insert(param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	response.NewSuccess().Json(ctx)
}

// 修改资料
func (*User) ProfileUpdate(ctx *gin.Context) {

	id, _ := strconv.Atoi(ctx.Param("id"))

	var param dto.UserProfileUpdateRequest

	if err := ctx.Bind(&param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	userId := ctx.GetInt("userId")

	if id != userId {
		response.NewError().SetMsg("越权操作，用户资料只能本人修改").Json(ctx)
		return
	}

	param.Id = id

	if err := (&service.Profile{}).Update(param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	response.NewSuccess().Json(ctx)
}

// 修改密码
func (*User) PasswordReset(ctx *gin.Context) {

	var param dto.UserUpdateRequest

	if err := ctx.Bind(&param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	userId, _ := strconv.Atoi(ctx.Param("id"))
	if userId <= 0 {
		userId = ctx.GetInt("userId")
	}

	if err := (&service.User{}).Update(dto.UserUpdateRequest{
		Id:       userId,
		Password: password.Generate(param.Password),
	}); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	response.NewSuccess().Json(ctx)
}

// 修改用户
func (*User) Update(ctx *gin.Context) {

	var param dto.UserUpdateRequest

	if err := ctx.Bind(&param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	param.Id, _ = strconv.Atoi(ctx.Param("id"))

	if err := (&service.User{}).Update(param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	response.NewSuccess().Json(ctx)
}
