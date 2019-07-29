package goods

import (
	"bluefoxgo/models"
	"bluefoxgo/tools"
	"github.com/kataras/iris"
	"log"
)
var goods models.K_goods

type Responce struct{
	Msg    string   `josn:"message"`
	Code   int		`josn:"code"`
}

type Require struct {
	Data interface{} `json:"data"`
}


//@api {get} /goods/goodslist/{page:int}/{rows:int}获取商品列表
func GetGoodsList(ctx iris.Context){
	page, _ := ctx.Params().GetInt("page")
	rows, _ := ctx.Params().GetInt("rows")
	goodslist,err := goods.GetAllGoods(page,rows)
	var myresp struct {
		Data interface{}
		Resp Responce
	}
	if err != nil{
		myresp.Data = nil
		myresp.Resp.Msg = err.Error()
		myresp.Resp.Code= 500
	}else{
		myresp.Data = tools.RowResult(goodslist)
		myresp.Resp.Msg = "获取成功"
		myresp.Resp.Code = 200

	}
	//log.Println("err:",err)
	//log.Println(goodslist)
	ctx.JSON(myresp)

}



//@api {get} /goods/search/:goods_name/:page/:rows 查询系列
func SearchGoods(ctx iris.Context){

	var  myresp  struct {
		Data []interface{}
		Resp Responce
	}
	true_name := ctx.Params().GetString("goods_name")
	page,err1 := ctx.Params().GetInt("page")
	rows,err2 := ctx.Params().GetInt("rows")
	log.Println(true_name)
	if err1 != nil || err2 !=nil{
		myresp.Data = nil
		myresp.Resp.Code = 400
		myresp.Resp.Msg = "参数错误"
		ctx.JSON(myresp)
		return
	}
	goodslist,err := goods.SearchGoods(true_name,page,rows)
	if err != nil{
		myresp.Data = nil
		myresp.Resp.Code = 500
		myresp.Resp.Msg = err.Error()
	}else{
		myresp.Data = tools.RowResult(goodslist)
		myresp.Resp.Code = 200
		myresp.Resp.Msg  = "获取成功"
	}
	//log.Println(goodslist)
	ctx.JSON(myresp)
}



//@api {get} /goods/info/:gid 获取用户详细信息
func GoodsInfo(ctx iris.Context){
	var myresp  struct {
		Data models.AllSkuInfo
		Resp Responce
	}
	gid,err := ctx.Params().GetInt("gid")
	//log.Println(uid)
	if err != nil{
		myresp.Data = models.AllSkuInfo{}
		myresp.Resp.Code = 400
		myresp.Resp.Msg  = "数据无效"
		ctx.JSON(myresp)
		return
	}
	if goodsinfo,err_string := goods.GoodsInfo(gid);len(err_string)!=0{
		myresp.Data = models.AllSkuInfo{}
		myresp.Resp.Code = 500
		myresp.Resp.Msg  = "暂无数据"
	}else{
		myresp.Data = goodsinfo
		myresp.Resp.Code = 200
		myresp.Resp.Msg  = "获取成功"
	}
	ctx.JSON(myresp)
}


func GetGoodsList2(ctx iris.Context) {
	page, _ := ctx.Params().GetInt("page")
	rows, _ := ctx.Params().GetInt("rows")
	goodslist, err := goods.GetAllGoodsToUseMysql(page, rows)
	var myresp struct {
		Data interface{}
		Resp Responce
	}
	if err != nil {
		myresp.Data = nil
		myresp.Resp.Msg = err.Error()
		myresp.Resp.Code = 500
	} else {
		myresp.Data = *goodslist
		myresp.Resp.Msg = "获取成功,form2"
		myresp.Resp.Code = 200
	}
	//log.Println("err:", err)
	//log.Println(*goodslist)
	send := tools.RowResult(goodslist)
	//result,err := scanner.ScanMap(goodslist)
	if err != nil{
		ctx.WriteString("出错infrom2")
	}
	myresp.Data = send
	ctx.JSON(myresp)

}

func Create(ctx iris.Context){
	var goods models.K_goods
	var sku   models.K_goods_sku
	require := struct{
		Goods interface{}
		Sku   interface{}
	}{}
	err := ctx.ReadJSON(&require)
	if err!=nil{
		ctx.JSON(Responce{"无法解析该json",400})
		return
	}
	goods,errgoods 	:= require.Goods.(models.K_goods)
	sku,errsku		:= require.Sku.(models.K_goods_sku)
	if !errgoods || !errsku {
		ctx.JSON(Responce{"非指定结构体",400})
	}
	var err_db error
	if goods_id,ok := goods.Isset(goods.Goods_cn);goods_id<0 || !ok{
		//不存在该商品，直接创建，将goods和sku传入CreateGoods()中
		err_db = goods.CreateGoods(goods,sku)
	}else{  //存在该商品
		err_db = goods.CreateGoods(sku)
	}
	if err_db==nil{
		ctx.JSON(Responce{"创建成功",200})
		return
	}
	ctx.JSON(Responce{err_db.Error(),500})
	//error_num,savenames,err := tools.UploadFileToServer(ctx)
	//goods.Goods_cn		= goods_name //默认为中文名
	//goods.Brand_id		= brand_id
	//goods.Market_price  = market_price
	//goods.Profile		= profile
	/*系列信息*/
	//goods_name  := ctx.PostValue("nick_name")
	//brand_id,_  := ctx.PostValueInt("brand_id")
	//market_price:= ctx.PostValue("market_price")
	//profile     := ctx.PostValue("profile")
	/*商品信息*/
	//sku.Sku_cn	    = ctx.PostValue("sku_name")  //商品名
	//sku.Sku_main	= ctx.PostValue("sku_main")  //色系
	//sku.Price,_		= ctx.PostValueFloat64("price")	  //价格
	//sku.Texture		= ctx.PostValue("texture")	  //质地
	//sku.Net_weight,_= ctx.PostValueFloat64("net_weight")//净含量
	//sku.Net_volume,_= ctx.PostValueFloat64("net_volume")//净容量
	//sku.Total_weight,_= ctx.PostValueFloat64("total_weight")//总重量
	//sku.Color_hex	= ctx.PostValue("color_hex") //颜色
	//sku.Color_rgb	= ctx.PostValue("color_rgb") //颜色RGB
	//Flag_new_int,_ := ctx.PostValueInt("flag_new")//1 新品 0  不是
	//sku.Flag_new	=int8(Flag_new_int)
	//Flag_hot_int,_	:=ctx.PostValueInt("flag_hot") //平台以外的来源
	//sku.Flag_hot	=int8(Flag_hot_int)
	//Status_sale_int,_:=ctx.PostValueInt("tatus_sale") //商品状态
	//sku.Status_sale = int8(Status_sale_int)
	//sku.Official_tag =ctx.PostValue("official_tag")  //官方标签
}


