package controller

import (
	"isme-go/app/dto"
	"isme-go/app/service"
	"isme-go/app/token"
	"isme-go/common/captcha"
	"isme-go/common/password"
	"isme-go/common/utils"
	"isme-go/framework/response"
	"strings"

	"github.com/gin-gonic/gin"
)

var captchaCache = make(map[string]string)

type Auth struct{}

// 验证码
func (*Auth) Captcha(ctx *gin.Context) {

	svg, code := captcha.New(&captcha.Config{
		Width:      80,
		Height:     40,
		FontSize:   20,
		CharsCount: 4,
	}).Generate()

	captchaCache[strings.ToLower(code)] = code

	// 返回svg文件
	ctx.Data(200, "image/svg+xml", svg)
}

// 登录
func (*Auth) Login(ctx *gin.Context) {

	var param dto.LoginRequest

	if err := ctx.Bind(&param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	if param.Username == "" || param.Password == "" {
		response.NewError().SetMsg("用户名或密码不能为空").Json(ctx)
		return
	}

	defer delete(captchaCache, strings.ToLower(param.Captcha))
	if _, ok := captchaCache[strings.ToLower(param.Captcha)]; !ok {
		response.NewError().SetCode(10003).SetMsg("验证码错误").Json(ctx)
		return
	}

	user := (&service.User{}).GetDetailByUsername(param.Username)
	if !user.Enable {
		response.NewError().SetMsg("用户已禁用").Json(ctx)
		return
	}

	if !password.Verify(user.Password, param.Password) {
		response.NewError().SetMsg("用户名或密码错误").Json(ctx)
		return
	}

	roleIds := (&service.UserRolesRole{}).GetRoleIdsByUserId(user.Id)
	roles := (&service.Role{}).GetListByIds(roleIds, true)
	if len(roles) <= 0 {
		response.NewError().SetMsg("用户未关联角色").Json(ctx)
		return
	}

	roleCodes := make([]string, 0)
	for _, role := range roles {
		roleCodes = append(roleCodes, role.Code)
	}

	accessToken := token.GetClaims(dto.UserTokenResponse{
		Id:              user.Id,
		Username:        user.Username,
		RoleCodes:       roleCodes,
		CurrentRoleCode: roleCodes[0],
	}).GenerateToken()

	response.NewSuccess().SetData("data", map[string]interface{}{
		"accessToken": accessToken,
	}).Json(ctx)
}

// 切换当前角色
func (*Auth) SwitchCurrentRole(ctx *gin.Context) {

	roleCode := ctx.Param("roleCode")

	userId := ctx.GetInt("userId")
	username := ctx.GetString("username")

	roleIds := (&service.UserRolesRole{}).GetRoleIdsByUserId(userId)
	roles := (&service.Role{}).GetListByIds(roleIds, true)

	roleCodes := make([]string, 0)
	for _, role := range roles {
		roleCodes = append(roleCodes, role.Code)
	}

	if !utils.Contains(roleCodes, roleCode) {
		response.NewError().SetMsg("您目前暂无此角色，请联系管理员申请权限").Json(ctx)
		return
	}

	accessToken := token.GetClaims(dto.UserTokenResponse{
		Id:              userId,
		Username:        username,
		RoleCodes:       roleCodes,
		CurrentRoleCode: roleCode,
	}).GenerateToken()

	response.NewSuccess().SetData("data", map[string]interface{}{
		"accessToken": accessToken,
	}).Json(ctx)
}

// 退出登录
func (*Auth) Logout(ctx *gin.Context) {
	response.NewSuccess().SetData("data", true).Json(ctx)
}

// 修改密码
func (*Auth) Password(ctx *gin.Context) {

	var param dto.PasswordRequest

	if err := ctx.Bind(&param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	if param.OldPassword == "" || param.NewPassword == "" {
		response.NewError().SetMsg("旧密码或新密码不能为空").Json(ctx)
		return
	}

	user := (&service.User{}).GetDetailById(ctx.GetInt("userId"))
	if !password.Verify(user.Password, param.OldPassword) {
		response.NewError().SetMsg("旧密码错误").Json(ctx)
		return
	}

	user.Password = password.Generate(param.NewPassword)

	if err := (&service.User{}).Update(dto.UserUpdateRequest{
		Id:       user.Id,
		Password: user.Password,
	}); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	response.NewSuccess().Json(ctx)
}
