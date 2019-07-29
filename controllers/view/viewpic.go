package view

import (
	"bluefoxgo/models"
	"bluefoxgo/tools"
	"github.com/kataras/iris"
)

var pic models.K_view_pic


type AllImageRequire struct {
	Vid		int64
	Data	[]models.K_view_pic
}



//@api [Post]/view/image/upload
// 上传图片
func UploadImage(ctx iris.Context){
	require := struct {
		vid 		int64
		view_pic	models.K_view_pic
	}{}
	if err := ctx.ReadJSON(require); err != nil{
		ctx.JSON(Responce{"无法解析请求json",400})
		return
	}
	res := pic.UploadOnePic(require.vid,require.view_pic)
	if res!= nil{
		ctx.JSON(Responce{"上传失败",500})
		return
	}
	ctx.JSON(Responce{"上传成功",200})
}



//@api [Post]/view/image/delete
// 删除图片
func DeleteImage(ctx iris.Context){
	vid,err1 := ctx.PostValueInt64("vid")
	pid,err2 := ctx.PostValueInt("pid")
	if err1 != nil || err2 !=nil{
		err_str := tools.ReturnAllError(err1,err2)
		ctx.JSON(Responce{err_str,400})
		return
	}
	res := pic.DeleteOnePic(vid,pid)
	if res != nil{
		ctx.JSON(Responce{"删除失败",500})
		return
	}
	ctx.JSON(Responce{"删除成功",200})
}


func UpDateAllImage(ctx iris.Context){
	requ := AllImageRequire{}
	if err := ctx.ReadJSON(requ); err != nil{
		ctx.JSON(Responce{"无法解析请求json",400})
		return
	}
	//同时需要修改两张表 view_pic和view表
	res := pic.UpdateAllViewPic(requ.Vid,requ.Data)
	ctx.JSON(Responce{res.Error(),200})
}

