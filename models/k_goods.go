package models

import (
	"bluefoxgo/tools"
	//goods2 "bluefoxgo/controllers/goods"
	"errors"
	_ "github.com/go-sql-driver/mysql"

	"database/sql"
)

type K_goods struct {
	Id				int
	Brand_id	 	int
	Goods_en		string
	Goods_cn		string
	Market_price	string
	Profile			string
	Goods_pic		string
	Tip				string
	Keywords		string
}





//获取所有系列
func(goods *K_goods) GetAllGoods(page int,rows int)(*sql.Rows,error){

	limit 	:= rows
	var offset	int
	if page > 0 {
		offset = (page-1)*limit
	}else{
		offset = 0
	}

	//var goodslist  []gorose.Map
	goodslist,err := db.Table("k_goods as gd").
		Joins("k_brand as bd","gd.brand_id","=","bd.id").
		Select("gd.id,bd.brand_cn,bd.brand_en,gd.goods_cn,gd.goods_en," + //id,系列，品牌中英文名
			"gd.market_price,gd.goods_pic"). //商品价格，销售状态，图片
		Limit(rows).Offset(offset).
		Rows()
	if err != nil{
		return nil,errors.New(err.Error())
	}

	return goodslist,nil
}

func(goods *K_goods) GetAllGoodsToUseMysql(page int,rows int)(*sql.Rows,error){


	//var goodslist  []K_goods
	res, err1 := db.Table("k_goods as gd").
		Select("gd.id,bd.brand_cn,bd.brand_en,gd.goods_cn,gd.goods_en,gd.market_price,gd.goods_pic").
		Joins("left join k_brand as bd on gd.brand_id = bd.id").
		Limit(rows).Offset((page-1)*rows).
		Rows()
	if err1 != nil{
		return nil,err1
	}
	return res,nil
}


//获取系列详情
//需要系列及商品详情
type AllSkuInfo map[string]interface{}

func(goods *K_goods) GoodsInfo(gid int) (AllSkuInfo,string){
	//获取系列详情
	BaseInfo,err_base := db.Table("k_goods").
					Where("id = ?",gid).
					Select("id,brand_id,goods_cn,goods_en,market_price,profile,goods_pic").
					Rows()
	//获取该系列所有商品信息
	GoodsInfo,err_goods := db.Table("k_goods_sku").
					Where("goods_id = ?",gid).
					Rows()
	allskuinfo := make(AllSkuInfo,2)
	err_string := ""
	if err_base != nil{
		allskuinfo["BaseInfo"]  = nil
		err_string += err_base.Error()
	}else{
		allskuinfo["BaseInfo"]  = tools.RowResult(BaseInfo)
	}
	if err_goods != nil{
		allskuinfo["GoodsInfo"]  = nil
		err_string += " And "+err_goods.Error()
	}else{
		allskuinfo["GoodsInfo"]  = tools.RowResult(GoodsInfo)
	}
	//log.Println(BaseInfo)
	return allskuinfo,err_string
}


//根据商品名查询商品
func(goods *K_goods) SearchGoods(goods_name string,page,rows int)(*sql.Rows,error){

	limit 	:= rows
	var offset	int
	if page > 0 {
		offset = (page-1)*limit
	}else{
		offset = 0
	}

	//大小写不敏感
	goodslist,err := db.Table("k_goods as gd").
		Joins("join k_brand as bd on gd.brand_id = bd.id").
		Select("gd.id,bd.brand_cn,bd.brand_en,gd.goods_cn,gd.goods_en," + //id,系列，品牌中英文名
			"gd.market_price,gd.goods_pic"). //商品价格，销售状态，图片
		Where("gd.goods_en LIKE ?", "%"+goods_name+"%").
		Or("gd.goods_cn LIKE ?", "%"+goods_name+"%").
		Limit(rows).Offset(offset).
		Rows()
	if err != nil{
		return nil,err
	}
	return goodslist,nil

}


//修改系列信息
//跟上传图片不在同一个函数中，图片上传在控制器中 且先获取goods_pic原字段信息，字符串合并后修改指定ID字段
func(goods *K_goods) UpdateGoods(gd K_goods)(bool,error){
	//dataforup := map[string]interface{}{
	//	"id":				gd.Id,
	//	"brand_id":	 	gd.Brand_id,
	//	"goods_en"	:		gd.Goods_en,  //默认修改英文名
	//	//"goods_cn":		string
	//	"market_price":	gd.Market_price,
	//	"profile":			gd.Profile,
	//	//"goods_pic":		gd.Goods_pic,
	//}
	res := db.Save(gd)
	if res.Error != nil{
		return false,res.Error
	}else{
		return true,nil
	}
}


//修改商品信息
//跟上传图片在同一个函数中，图片上传在控制器中 先从前端获取当前图片字段，字符串合并后修改指定ID的sku_pic字段
func(goods *K_goods) UpdateGoodsSku(sku K_goods_sku)(bool,error){
	dataforup := map[string]interface{}{
		"id":			sku.Id,
		"sku_main":	sku.Sku_main,
		//Sku_en		string  默认修改中文名
		"sku_cn":		sku.Sku_cn,
		"sku_pic":		sku.Sku_pic,
		"price"	:		sku.Price,
		"texture":		sku.Texture,
		"net_weight":	sku.Net_weight,
		"net_volume":	sku.Net_volume,
		"color_hex":	sku.Color_hex,
		"color_rgb":	sku.Color_rgb,
		"image_uri":	sku.Image_uri,
		"status_sale":	sku.Status_sale,
		"flag_new":	sku.Flag_new,
		"flag_hot":	sku.Flag_hot,
		"official_tag":	sku.Official_tag,
		"total_weight":	sku.Total_weight,
	}

	err := db.Save(&dataforup)
	if err != nil{
		return false,err.Error
	}else{
		return true,nil
	}
}


//上传系列信息，商品信息
//goods和sku都现在数据库中insert后 再进行图片上传 且goods不在同一个函数中，sku在一个函数中
func(goods *K_goods) CreateGoods(args ...interface{})(error){

	if len(args)==2{  //双参数,插入goods和sku
		db.Begin()
		gd,_ := args[0].(K_goods)
		res1 := db.Create(gd)
		if res1.Error !=nil {
			db.Rollback()
		}
		sku,_ := args[1].(K_goods_sku)
		res2 := db.Create(sku)
		if res2.Error != nil{
			db.Rollback()
		}
		db.Commit()
		if res1.Error !=nil || res2.Error != nil {
			errstr := res1.Error.Error()+"||"+res2.Error.Error()
			return errors.New(errstr)
		}
		return nil
	}
	//单参数，插入sku
	sku,_ := args[0].(K_goods_sku)
	res2 := db.Create(sku)
	if res2.Error != nil{
		return  res2.Error
	}
	return nil

	//datagoods_forins := map[string]interface{}{
	//	//id				int
	//	"brand_id":	 	gd.Brand_id,
	//	"goods_en"	:		gd.Goods_en,  //默认修改英文名
	//	//"goods_cn":		string
	//	"market_price":	gd.Market_price,
	//	"profile":			gd.Profile,
	//	//"goods_pic":		gd.Goods_pic,
	//}
	//datasku_forins := map[string]interface{}{
	//	//id			int
	//	"sku_main":	sku.Sku_main,
	//	//Sku_en		string  默认修改中文名
	//	"sku_cn":		sku.Sku_cn,
	//	"sku_pic":		sku.Sku_pic,
	//	"price"	:		sku.Price,
	//	"texture":		sku.Texture,
	//	"net_weight":	sku.Net_weight,
	//	"net_volume":	sku.Net_volume,
	//	"color_hex":	sku.Color_hex,
	//	"color_rgb":	sku.Color_rgb,
	//	"image_uri":	sku.Image_uri,
	//	"status_sale":	sku.Status_sale,
	//	"flag_new":	sku.Flag_new,
	//	"flag_hot":	sku.Flag_hot,
	//	"official_tag":	sku.Official_tag,
	//	"total_weight":	sku.Total_weight,
	//}
}



//判断是否已经存在该商品
func(goods K_goods)Isset(goods_name string)(int,bool){
	goodss := K_goods{}
	db.Where("goods_cn = ?", goods_name).
		Or("goods_en = ?",goods_name).
		Select("id").
		First(&goodss)
	if goodss.Id == 0 {
		return -1,false
	}
	return goodss.Id,true
}

