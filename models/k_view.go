package models

import (
	"bluefoxgo/tools"
	"database/sql"
	"errors"
	"strconv"
	"strings"
)

type K_view struct {
	Id				int64
	Publish_time	int   //发表时间
	Uid				int
	Thumbup_total	int	  //点赞数量
	Tags			string	//话题等于标签，可为空
	Word_total		string	//总字数
	Pic_total		string	//图片多少
	Reading_speed	string  //阅读时间
	Title			string  //标题长度
	Content			string  //内容
	Relative_sku	string	//相关sku，用逗号隔开
	Is_delete		int8    //是否删除，0未删除，1删除
}




//获取所有视角
func(u *K_view) GetAllView(offset int,rows int)(*sql.Rows,error){

	viewlist,err := db.Table("k_view as v").
		Joins("join k_user as u on u.id = v.uid").
		Joins("left join k_view_collection as c on c.view_id = v.id").
		Select("v.id,v.title,u.nick_name,v.word_total,COUNT('c.view_id'),v.publish_time,v.is_delete").
		Group("v.id,v.title,u.nick_name,v.word_total,v.publish_time,v.is_delete").
		Limit(rows).
		Offset(offset).
		Rows()
	if err != nil{
		return nil,err
	}
	return viewlist,nil
}



//查询用户视角
func(u *K_view) Serach(user_name string,offset,rows int)(*sql.Rows,error){

	viewlist,err := db.Table("k_view as v").
		Joins("join k_user as u on u.id = v.uid").
		Joins("left join k_view_collection as c on c.view_id = v.id").
		Where("u.nick_name LIKE ?","%"+user_name+"%").
		Or("u.true_name LIKE ?","%"+user_name+"%").
		Select("v.id,v.title,u.nick_name,v.word_total,COUNT('c.view_id') AS collection_num,v.publish_time,v.is_delete").
		Group("v.id,v.title,u.nick_name,v.word_total,v.publish_time,v.is_delete").
		Limit(rows).
		Offset(offset).
		Rows()
	if err != nil{
		return nil,err
	}
	return viewlist,nil
}




type ViewInfo map[string]interface{}

//获取视角详细信息
//此页面的获取评论将单独设置api
func(view *K_view) Info(vid int64)(ViewInfo,error){
	var BaseInfo,ViewImage,LinkSku	*sql.Rows
	var err1,err2,err3 error
	//获取基本信息
	BaseInfo,err1 = db.Table("k_view AS v").
		Joins("join k_user AS u on u.id = v.uid").
		Joins("left join k_view_collection AS c on c.view_id = v.id").
		Where("v.id = ?",vid).
		Select("v.title,v.thumbup_total,v.publish_time,COUNT('c.view_id') AS collection_num," +
			" u.nick_name,v.is_delete,v.tags,v.content,v.tags").
		Group("v.title,v.thumbup_total,v.publish_time," +
		" u.nick_name,v.is_delete,v.tags,v.content,v.tags").
		Rows()

	//图片uri
	ViewImage,err2 = db.Table("k_view_pic").
		Where("view_id = ?",vid).
		Select("id,image_uri,header_light,footer_light").
		Rows()

	//获取关联产品
	//先获取关联的sku_id
	var sku_id_struct struct{
		Relative_sku string
	}
	res := db.Table("k_view").
			Where("id = ?",vid).
			Select("relative_sku").
			Scan(&sku_id_struct)
	if res.Error != nil{
		return nil,res.Error
	}
	//将sku_id(string)转换为切片
	id := tools.StringToIntArray(sku_id_struct.Relative_sku)
	if len(id)==0{
		LinkSku = nil
	}else{
		//如果sku_id数组不为空，查询数据库
		LinkSku,err3 = db.Table("k_goods_sku AS sku").
			Joins("join k_goods AS gd on gd.id = sku.goods_id").
			Joins("join k_brand AS bd on bd.id = gd.brand_id").
			Select("bd.brand_cn,bd.brand_en," +  //品牌中英文名
				"gd.id,gd.goods_cn,gd.goods_en," + //系列ID，中英文名
				"sku.id,sku.sku_main,sku.sku_cn,sku.sku_en,sku.price"). //sku色系，价格，中英文名
			Where("sku.id IN (?)",id).
			Rows()
	}
	if err1!=nil || err2!=nil || err3!=nil{
		str := tools.ReturnAllError(err1,err2,err3)
		return nil,errors.New(str)
	}
	viewinfo := make(ViewInfo,3)
	viewinfo["BaseInfo"] = tools.RowResult(BaseInfo)
	viewinfo["ViewImage"]= tools.RowResult(ViewImage)
	viewinfo["LinkSku"]  = tools.RowResult(LinkSku)
	return  viewinfo,nil

}




func(view *K_view) BanView(vid int64)error{
	//先查询该用户为ban还是非ban
	v := K_view{}
	err := 	db.Select("id,is_delete").First(&v,vid)
	if err !=  nil {
		return err.Error
	}
	//0未删除  1已删除
	if v.Is_delete == 1{  //当前为1,修改为0
		err = db.Model(&v).Update("is_delete", 0)
	}else{   ////当前为0，修改为1
		err = db.Model(&v).Update("is_delete", 1)
	}

	if  err != nil{
		return err.Error
	}
	return nil
}

func(view *K_view) UpdataBaseInfo(baseinfo K_view)error{
	res := db.Model(&view).Updates(baseinfo)
	if res.Error != nil{
		return res.Error
	}
	return nil
}


func(view *K_view) AddLinkSku(vid int64,sku_main string)(*sql.Rows,error){
	//先获取sku的ID  所以首先需要查询 sku表获取sku_id
	sku := K_goods_sku{}
	res := db.Table("k_goods_sku").
		Select("id").
		Where("sku_main = ?",sku_main).First(&sku)
	//需要返回的值
	sku_info,err  := db.Table("k_goods_sku AS sku").
		Joins("join k_goods AS gd on sku.goods_id = gd.id").
		Joins("join k_brand AS bd on gd.brand_id = bd.id").
		Select("sku.id,sku.sku_main,sku.price," +
			"gd.id,gd.goods_cn,gd.goods_en," +
			"bd.brand_cn,bd.brand_en").
		Where("sku.id = ?",sku.Id).
		Rows()
	if err!=nil{
		return nil,err
	}
	//将sku_id加入到view中
	//tx := db.Begin()
	view_one  := K_view{}
	if resv := db.Table("k_view").Where("id = ?",vid).
		Select("id,relative_sku").First(&view_one);resv.Error!=nil{
			return nil,resv.Error
	}
	string_link_sku := view_one.Relative_sku+","+strconv.Itoa(sku.Id)
	res = db.Table("k_view").Where("id = ?", vid).Update("relative_sku",string_link_sku)
	if res.Error != nil{
		return nil,res.Error
	}
	return sku_info,nil
}

func(view *K_view) DeleteLinkSku(sku_id int,vid int64)error{

	view_one  := K_view{}
	if resv := db.Table("k_view").Where("id = ?",vid).
		Select("id,relative_sku").First(&view_one);resv.Error!=nil{
		return resv.Error
	}
	link_sku_arr := strings.Split(view_one.Relative_sku,",")
	tag := -1  //标记 标记的ID将会从link_sku_arr删除
	sku_id_for_str := strconv.Itoa(sku_id)
	for index,v := range link_sku_arr{
		if v == sku_id_for_str {
			tag = index
			break
		}
	}
	if tag >=0 {
		link_sku_arr = append(link_sku_arr[:tag],link_sku_arr[tag+1:]...)
	}else{
		return errors.New("无法删除,未找到此view下的关联SKU有此SKU")
	}
	string_link_sku := strings.Join(link_sku_arr,",")
	res := db.Table("k_view").Where("id = ?", vid).Update("relative_sku",string_link_sku)
	if res.Error != nil{
		return res.Error
	}
	return nil
}

