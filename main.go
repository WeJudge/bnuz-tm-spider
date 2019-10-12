package main

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/recover"
	"net/http/cookiejar"
)

const (
	PortalToken string = "12345"
	PortalListen string = "0.0.0.0:7777"
)

type UserInfo struct {
	TeamName string `json:"team_name"`
	RealName string `json:"real_name"`
	SchoolName string `json:"school_name"`
	CoachName string `json:"coach_name"`
}

func main() {
	app := iris.New()
	app.Logger().SetLevel("debug")
	app.Use(recover.New())
	app.Use(logger.New())

	app.Handle("POST", "/check", func(ctx iris.Context) {

		success := CheckSignaturesFromIrisContext(ctx)
		if !success {
			return
		}

		username := ctx.PostValueDefault("username", "")
		password := ctx.PostValueDefault("password", "")

		cookies, err := cookiejar.New(nil)
		if err != nil {
			return
		}

		xsrftoken := initLogin(cookies)
		ret := postLogin(cookies, xsrftoken, username, password)

		if ret {
			userInfo := getUserInfo(cookies)
			if userInfo != nil {
				ui := UserInfo {
					TeamName: userInfo.User.Name,
					RealName: userInfo.User.Id,
				}
				_, _ = ctx.JSON(RESTfulAPIResult{
					Status: true,
					Data: ui,
				})
			} else {
				_, _ = ctx.JSON(RESTfulAPIResult{
					Status: false,
					Message: "解析信息失败",
				})
			}
		} else {
			_, _ = ctx.JSON(RESTfulAPIResult{
				Status: false,
				Message: "登录认证失败",
			})
		}
	})

	_ = app.Run(iris.Addr(PortalListen), iris.WithoutServerError(iris.ErrServerClosed))
}