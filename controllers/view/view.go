package view

import (
	"bluefoxgo/models"
	"bluefoxgo/tools"
	"github.com/kataras/iris"
)


type Responce struct{
	Msg    string   `josn:"message"`
	Code   int		`josn:"code"`
}



var view models.K_view
var comment models.K_view_comments


//@api [Get]/view/viewlist/{page:int}/{rows:int}"
//获取视角列表
func  GetViewList(ctx iris.Context){
	var myre struct{
		Data interface{}
		Resp Responce
	}

	page,errpage := ctx.Params().GetInt("page")
	rows,errrows := ctx.Params().GetInt("rows")
	if errpage != nil || errrows != nil{
		myre.Data = nil
		myre.Resp = Responce{"参数无效",400}
		ctx.JSON(myre)
		return
	}
	offset := (page-1)*rows
	viewlist,err := view.GetAllView(offset,rows)
	if err != nil{
		myre.Data = nil
		myre.Resp = Responce{err.Error(),500}
		ctx.JSON(myre)
		return
	}
	myre.Data = tools.RowResult(viewlist)
	myre.Resp = Responce{"获取成功",200}
	ctx.JSON(myre)
}



//@api [Get]/view/search/{user_name:string}/{page:int}/{rows:int}
// 查询该用户的视角
func Search(ctx iris.Context){
	myres := struct {
		Data interface{}
		Resp Responce
	}{}
	user_name := ctx.Params().GetString("user_name")
	page,errp := ctx.Params().GetInt("page")
	rows,errw := ctx.Params().GetInt("rows")
	if errp!=nil || errw!=nil{
		myres.Data = nil
		myres.Resp = Responce{"参数错误",400}
		ctx.JSON(myres)
		return
	}
	offset := (page-1)*rows
	searchlist,err := view.Serach(user_name,offset,rows)
	if err!=nil{
		myres.Data = nil
		myres.Resp = Responce{err.Error(),500}
		ctx.JSON(myres)
		return
	}
	myres.Data = tools.RowResult(searchlist)
	myres.Resp = Responce{"获取成功",200}
	ctx.JSON(myres)
}


//@api [Get]/view/info/{vid:int}
// 获取视角详细信息
func Info(ctx iris.Context){
	myres := struct {
		Data interface{}
		Resp Responce
	}{}
	vid,err := ctx.Params().GetInt64("vid")
	if err != nil{
		myres.Data = nil
		myres.Resp = Responce{"参数无效",400}
		ctx.JSON(myres)
		return
	}
	viewinfo,errdb := view.Info(vid)
	if errdb!= nil{
		myres.Data = nil
		myres.Resp = Responce{errdb.Error(),500}
		ctx.JSON(myres)
		return
	}
	myres.Data = viewinfo
	myres.Resp = Responce{"获取成功",200}
	ctx.JSON(myres)
}



type comb struct{
	Comment models.K_view_comments
	Reply	[]models.K_view_comments
}
//@api [Get]/view/comment/{page:int}/{rows:int}
//获取详细视角页面中的评论信息
func Comment(ctx iris.Context){
	var myres struct{
		Data interface{}
		Resp Responce
	}
	vid,errid := ctx.Params().GetInt64("vid")
	page,errp := ctx.Params().GetInt("page")
	rows,errw := ctx.Params().GetInt("rows")
	if errp != nil || errw != nil || errid != nil{
		myres.Data = nil
		myres.Resp = Responce{"参数无效",400}
	}
	offset := (page-1)*rows
	comment,reply,err := comment.Comment(vid,offset,rows)
	if err != nil{
		myres.Data = nil
		myres.Resp = Responce{err.Error(),500}
		ctx.JSON(myres)
		return
	}
	//首先对parent_id进行分类,然后放入评论中
	combine := SortReplyAndCombine(reply,comment)
	myres.Data = combine
	myres.Resp = Responce{"获取成功",200}
	ctx.JSON(myres)
}




//@api [Post]/view/ban
// 禁止此视角
func Ban(ctx iris.Context){
	vid,err := ctx.PostValueInt64("vid")
	if err!=nil{
		ctx.JSON(Responce{"参数错误或未收到参数",400})
		return
	}
	err = view.BanView(vid)
	if err != nil{
		ctx.JSON(Responce{"修改失败",500})
		return
	}
	ctx.JSON(Responce{"修改成功",200})
}






func SortReplyAndCombine(reply []models.K_view_comments,comment []models.K_view_comments)([]comb){
	mapforreply := make(map[int64][]models.K_view_comments,len(reply)) //map的最大长度不会超过reply的长度
	for _,one_reply := range reply { //进行分类
		//_,ok := mapforreply[one_reply.Parent_id]
			mapforreply[one_reply.Parent_id] = append(mapforreply[one_reply.Parent_id],one_reply)
	}
	combine := []comb{}
	for _,one_comment := range comment{
		v,isset := mapforreply[one_comment.Id]
		var one_combine comb
		if isset { //存在
			one_combine = comb{one_comment,v}
		}else{
			one_combine = comb{one_comment,nil}
		}
		combine = append(combine,one_combine)
	}
	return combine
}



//api [Post]/view/update
//修改视角基本信息
func UpdateBaseinfo(ctx iris.Context){
	view := models.K_view{}
	if err := ctx.ReadJSON(view);err!=nil{
		ctx.JSON(Responce{err.Error(),400})
		return
	}
	res := view.UpdataBaseInfo(view)
	if res != nil{
		ctx.JSON(Responce{"修改失败",500})
		return
	}
	ctx.JSON(Responce{"修改成功",200})
}



//@api [Post]/view/sku/add
// 关联sku
func  AddSku(ctx iris.Context){
	myre := struct {
		Data interface{}
		Resp Responce
	}{}
	sku_name := ctx.PostValue("sku_main")
	vid,_		 := ctx.PostValueInt64("vid")

	res,err := view.AddLinkSku(vid,sku_name)
	if err!=nil{
		myre.Data = nil
		myre.Resp = Responce{err.Error(),500}
		ctx.JSON(myre)
		return
	}
	myre.Data = tools.RowResult(res)
	myre.Resp = Responce{"获取成功",200}
	ctx.JSON(myre)
	return
}



//@api [Post]/view/sku/delete
// 关联sku
func DeleteSku(ctx iris.Context){
	sku_id,err1	:= ctx.PostValueInt("sku_id")
	vid,err2	:= ctx.PostValueInt64("vid")
	if err1 != nil || err2 != nil {
		ctx.JSON(Responce{tools.ReturnAllError(err1,err2),400})
		return
	}
	res := view.DeleteLinkSku(sku_id,vid)
	if res!=nil{
		ctx.JSON(Responce{res.Error(),500})
		return
	}
	ctx.JSON(Responce{"删除成功",200})
}




//@api [Post]/view/comment/delete
// 删除此条评论
func DeleteComment(ctx iris.Context){
	cid,err := ctx.PostValueInt64("cid")
	if err != nil{
		ctx.JSON(Responce{"参数无效",400})
		return
	}
	res := comment.DeleteComment(cid)
	if res != nil {
		ctx.JSON(Responce{res.Error(),500})
		return
	}
	ctx.JSON(Responce{"删除成功",200})
}