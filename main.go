package main

import (
	"bluefoxgo/controllers/user"
	"bluefoxgo/controllers/goods"
	"bluefoxgo/controllers/view"
	_ "bluefoxgo/router"
	"github.com/kataras/iris"
)

func main() {
	app := iris.New()

	//用户板块
	app.PartyFunc("/user", func(users iris.Party) {
		//users.Use(myAuthMiddlewareHandler)
		users.Get("/search/{true_name:string}/{page:int}/{rows:int}",user.Search) //查询用户
		users.Get("/userlist/{page:int}/{rows:int}",user.GetUserList ) //获取用户列表
		users.Get("/info/{uid:int64}",user.Info)	// 获取用户详细信息
		users.Post("/ban",user.Ban)        			//禁用或者解禁用户
		users.Post("/ban/list",user.Banlist)     	//批量禁用或者解禁用户
		users.Post("/create",user.Create)     		//创建用户

	})

	//商品板块
	app.PartyFunc("/goods", func(goodss iris.Party) {
		goodss.Get("/search/{goods_name:string}/{page:int}/{rows:int}",goods.SearchGoods) //查询系列
		goodss.Get("/goodslist/{page:int}/{rows:int}",goods.GetGoodsList2 ) //获取系列列表
		goodss.Get("/info/{gid:int}",goods.GoodsInfo)	// 获取系列详细信息
		goodss.Post("/create",goods.Create)     		    //创建(goods&sku)商品
		goodss.Post("/sku/update",goods.SkuUpdate)						//修改商品(sku)信息
		goodss.Post("/sku/delete",goods.SkuDelete)						//删除商品(sku)信息
		goodss.Post("/sku/create",goods.SkuCreate)						//添加商品(sku)信息
	})

	//视角板块
	app.PartyFunc("/view", func(views iris.Party) {
		views.Get("/search/{user_name:string}/{page:int}/{rows:int}",view.Search) //查询该用户的视角
		views.Get("/viewlist/{page:int}/{rows:int}", view.GetViewList) //获取视角列表
		views.Get("/info/{vid:int64}",view.Info)						 // 获取视角详细信息
		views.Get("/comment/{vid:int64}/{page:int}/{rows:int}",view.Comment)  //获取详细视角页面中的评论信息
		views.Post("/ban}",view.Ban)										 //禁止此视角
		views.Post("/comment/delete",view.DeleteComment)				 //删除此条评论
		views.Post("/image/upload}",view.UploadImage)					 //上传图片
		views.Post("/image/delete}",view.DeleteImage)					 //删除图片
		views.Post("/update}",view.UpdateBaseinfo)						 //修改视角基本信息
		views.Post("/sku/add",view.AddSku)								 //关联sku
		views.Post("/sku/delete",view.DeleteSku)							 //删除关联sku
	})

	//视角模块

	app.Run(iris.Addr(":8080"))

}