package user

import (
	"bluefoxgo/models"
	"bluefoxgo/tools"
	"github.com/kataras/iris"
	"log"
)

var user models.K_user
type Responce struct{
	Msg    string   `josn:"message"`
	Code   int		`josn:"code"`
}
type Require struct {
	Data interface{}
}


//@api {get} /user/userlist/:page/:rows获取用户列表
func GetUserList(ctx iris.Context) {   //获取用户列表
	page, _ := ctx.Params().GetInt("page")
	rows, _ := ctx.Params().GetInt("rows")
	userlist,err := user.GetAllUser(page,rows)
	var myresp struct {
		Data interface{}
		Resp Responce
	}
	if err != nil{
		myresp.Data = nil
		myresp.Resp.Msg = err.Error()
		myresp.Resp.Code= 500
	}else{
		myresp.Data = tools.RowResult(userlist)
		myresp.Resp.Msg = "获取成功"
		myresp.Resp.Code = 200

	}
	////log.Println("v:userlist:",userlist)
	////log.Println(reflect.TypeOf(userlist))
	//log.Println("v:myrespformap:",myrespformap)
	//log.Println(reflect.TypeOf(myrespformap))
	ctx.JSON(myresp)
}



//@api {post} /user/ban 禁用或者解禁用户
//对接时需要修改post内容  重新获取post中的数据
func Ban(ctx iris.Context) {
	require  := Require{}
	responce := Responce{}
	if err := ctx.ReadJSON(&require); err != nil {
		//ctx.StatusCode(iris.StatusUnauthorized)
		responce.Msg  = err.Error()
		responce.Code = 400
	}
	uid,err := require.Data.(int64)
	if uid < 0 || !err  {
		responce.Msg  = "用户ID无效"
		responce.Code = 400
	}
	if res,err := user.DelectUser(uid);!res||err!=nil{ //出错
		responce.Code = 500
		responce.Msg  = err.Error()
	}else{
		responce.Msg  = "修改成功"
		responce.Code = 200
	}
	ctx.JSON(responce)
}


//@api {post} /user/ban/list 批量禁用或者解禁用户
//对接时需要修改post内容  重新获取post中的数据
func Banlist(ctx iris.Context){
	var useridlist []int64
	require  := Require{}
	responce := Responce{}
	if err := ctx.ReadJSON(&require); err != nil {
		responce.Msg  = err.Error()
		responce.Code = 400
		ctx.JSON(responce)
		return
	}
	//可以先断言
	//switch require.Data.(type) {
		//case []int64 :
	useridlist,err := require.Data.([]int64)
	//}
	if  !err || len(useridlist) == 0{
		responce.Code = 400
		responce.Msg  = "数据无效或者无数据"
		ctx.JSON(responce)
		return
	}
	res,err1 := user.BanorDebanUser(useridlist)
	if len(res)!=0 || err1 != nil{
		responce.Code = 500
		responce.Msg  = "内部错误"
	} else{
		responce.Code = 200
		responce.Msg  = "修改成功"
	}
	ctx.JSON(responce)
}



//@api {post} /user/create 创建用户
//对接时用户的数据需要修改
func Create(ctx iris.Context){
	//var insectuser models.K_user
	var responce Responce
	//var require  Require
	nick_name := ctx.PostValue("nick_name")
	gender,err:= ctx.PostValueInt("gender")
	if err != nil{  //默认为男
		gender = 1
	}
	profile := ctx.PostValue("profile")
	photo_uri := ctx.PostValue("photo_uri")
	//num,savepath,err:=tools.UploadFileToServer(ctx)
	//if num != 0 || err != nil{
	//	responce.Code = 500
	//	responce.Msg  = err.Error()
	//	//ctx.StopExecution()
	//	ctx.JSON(responce)
	//	return
	//}
	//fileNames,erross := tools.UploadToOSSApi(ctx,savepath,"profile_photo")
	//if err != nil{
	//	responce.Code = 500
	//	responce.Msg  = erross.Error()
	//}
	user := models.K_user{Nick_name:nick_name,Gender:int8(gender),Profile:profile,Photo_uri:photo_uri}
	ok := user.CreateUser(user)
	if ok != nil{
		responce.Code = 500
		responce.Msg  = "创建失败"
		//ctx.StopExecution()
		ctx.JSON(responce)
		return
	}
	responce.Code = 200
	responce.Msg  = "创建成功"
	ctx.JSON(responce)
}



//@api {get} /user/search/:true_name/:page/:rows 查询用户
//对接时用户的数据需要修改
func Search(ctx iris.Context){

	var  myresp  struct {
		Data []interface{}
		Resp Responce
	}
	true_name := ctx.Params().GetString("true_name")
	page,err1 := ctx.Params().GetInt("page")
	rows,err2 := ctx.Params().GetInt("rows")
	//log.Println(true_name)
	if err1 != nil || err2 !=nil{
		myresp.Data = nil
		myresp.Resp.Code = 400
		myresp.Resp.Msg = "参数错误"
		ctx.JSON(myresp)
		return
	}
	userlist,err := user.SearchUser(true_name,page,rows)
	if err != nil{
		myresp.Data = nil
		myresp.Resp.Code = 500
		myresp.Resp.Msg = err.Error()
	}else{
		myresp.Data = tools.RowResult(userlist)
		myresp.Resp.Code = 200
		myresp.Resp.Msg  = "获取成功"
	}
	//log.Println(userlist)
	ctx.JSON(myresp)
}


//@api {get} /user/info/:uid 获取用户详细信息
func Info(ctx iris.Context){
	var myresp  struct {
		Data models.Allinfo
		Resp Responce
	}
	uid,err := ctx.Params().GetInt64("uid")
	log.Println(uid)
	if err != nil{
		myresp.Data = models.Allinfo{}
		myresp.Resp.Code = 400
		myresp.Resp.Msg  = "数据无效"
		ctx.JSON(myresp)
		return
	}
	if Userinfo,err := user.Information(uid);err!=nil{
		myresp.Data = models.Allinfo{}
		myresp.Resp.Code = 500
		myresp.Resp.Msg  = "暂无数据"
	}else{
		myresp.Data = Userinfo
		myresp.Resp.Code = 200
		myresp.Resp.Msg  = "获取成功"
	}
	log.Println(myresp)
	//uc.Data["json"] = myresponce
	ctx.JSON(myresp)
}

func isdemo()
{
	return 
}
