package models

import (
	"bluefoxgo/tools"
	"database/sql"
	"errors"
	"fmt"
	//"github.com/gohouse/gorose"

	//"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)
type K_user struct {

	Id				int64
	Mini_openid 	string
	True_name		string
	Union_id		string
	Nick_name		string
	Photo_uri		string
	Avatar_uri		string
	Status			int8
	Tel				string
	Last_login_time	int
	Create_time		int
	Gender			int8
	Platform_name	string
	Profile			string
}



//获取所有用户
func(u *K_user) GetAllUser(page int,rows int)(*sql.Rows,error){

	var offset	int
	if page > 0 {
		offset = (page-1)*rows
	}else{
		offset = 0
	}

	userlist,err := db.Table("k_user").
			Select("id,photo_uri,nick_name,true_name,tel,gender,last_login_time,create_time,status").
			Limit(rows).
			Offset(offset).
			Rows()
	if err != nil{
		return nil,err
	}
	return userlist,nil
}




//ban和解ban 用户 对应update修改
func(u *K_user ) DelectUser(id int64)(bool,error){
	//先查询该用户为ban还是非ban
	user := K_user{}
	err := 	db.Select("id,status").First(&user,id)
	if err !=  nil {
		return false,errors.New("修改失败")
	}
	//var errdate error
	//1 禁止 10 正常
	if user.Status == 1{  //当前为1,修改为10
		 err = db.Model(&user).Update("status", 10)
	}else{   ////当前为10，修改为1
		 err = db.Model(&user).Update("status", 1)
	}

	if  err != nil{
		return false,errors.New("修改失败")
	}
	return true,nil
}



//批量ban和解ban 用户
func(u *K_user) BanorDebanUser(useridlist []int64)([]int64,error){
	result := []int64{}
	for _,id := range useridlist {
		if res,err := u.DelectUser(id);!res||err!=nil{
			result = append(result,id)
		}
	}
	if len(result)==0 {
		return []int64{},nil
	}
	return result,errors.New("存在未被执行用户")
}


//创建用户
//到时候last_login_time,union_id,create_time的获取方式可能需要修改
func(u *K_user) CreateUser(user K_user) error{
	// 一条数据
	//var insterdata = map[string]interface{}{
	//	"mini_openid" :	user.Mini_openid,
	//	"true_name":		user.True_name,
	//	"union_id"	:		user.Union_id,
	//	"nick_name":		user.Nick_name,
	//	"photo_uri":		user.Photo_uri,
	//	"avatar_uri":		user.Avatar_uri,
	//	"status"	:		user.Status,
	//	"tel"	:			user.Tel,
	//	"last_login_time":user.Last_login_time,
	//	"create_time":		user.Create_time,
	//	"gender":			user.Gender,
	//	"platform_name	string": user.Platform_name,
	//	"profile":			user.Profile,
	//}
	// insert into user (age, job) values (17, 'it3')
	ok := db.NewRecord(user)
	if !ok  {
		return errors.New("未创建成功")
	}

	return nil
}


//查询用户
func(u *K_user) SearchUser(true_name string,page,rows int)(*sql.Rows,error){

	limit := rows
	var offset	int
	if page > 0 {
		offset = (page-1)*limit
	}else{
		offset = 0
	}
	//var userlist  []K_user
	////大小写不敏感
	userlist,err := db.Table("k_user").
		Where("true_name LIKE ?", "%"+true_name+"%").
		Select("id,photo_uri,nick_name,true_name,tel,gender,last_login_time,create_time,status").
		Offset(offset).
		Limit(rows).
		Rows()

	if err  != nil{
		return nil,err
	}
	return userlist,nil

}

//获取用户详细信息
//分为四块：用户基本信息，地址，愿望单，其他信息
type Allinfo map[string]interface{}

func(u *K_user) Information(uid int64)(Allinfo,error){
	//前端分块显示 所以这里用多次单标查询  不使用连表查询
	//db := DB()
	//获取用户基本信息
	allinfo := make(Allinfo,4)
	userinfo,err1 := db.Table("k_user").
		Select("id,mini_openid,photo_uri,avatar_uri,nick_name,true_name,platform_name,tel," +
			"gender,last_login_time,create_time,status,profile").
		Where("id = ?", uid).
		Rows()
	if err1!=nil{
		allinfo["Baseinfo"] = nil
	}else{
		allinfo["Baseinfo"] = tools.RowResult(userinfo)
		fmt.Println("Baseinfo:",allinfo["Baseinfo"])
	}

	//收货地址
	address,err2  := db.Table("k_user_addresses as ad").
		Joins("left join k_cpr as cpr on ad.city_id = cpr.id").
		//Joins("left join emails on emails.user_id = users.id").Scan(&results)
		Select("ad.details,ad.tel,cpr.name,ad.true_name,ad.address_cpr").
		Where("ad.uid = ?",uid).
		Rows()
	if err2!=nil{
		allinfo["Address"] = nil
	}else{
		allinfo["Address"] = tools.RowResult(address)
		fmt.Println("Address:",allinfo["Address"])
	}
	//愿望单
	goods_col,err3 := db.Table("k_goods_collection as gc").
		Joins("join k_goods_sku on gc.sku_id = k_goods_sku.id").
		Joins("join k_goods on k_goods_sku.goods_id = k_goods.id").
		Joins("join k_brand on k_goods.brand_id = k_brand.id").
		Select("k_goods_sku.sku_pic,k_goods_sku.color_hex," +
			"K_goods.goods_cn,k_brand.brand_cn,gc.status").
		Where("gc.uid = ?",uid).
		Rows()
	if err3!=nil{
		allinfo["Goods_col"] = nil
	}else{
		allinfo["Goods_col"] = tools.RowResult(goods_col)
		fmt.Println("Goods_col:",allinfo["Goods_col"])
	}
	//其他信息
	//视角
	otherinfo := make(map[string]interface{},3)
	view 	 := 0  // 视角
	impscount:= 0  // 印象
	invcount := 0  // 是否参加调查

	//视角
	db.Table("k_view").Where("uid=?",uid).Count(&view)
	//log.Println(viewcount)
	//log.Println(reflect.TypeOf(viewcount))
	if view == 0{
		otherinfo["view"] = 0
	}else{
		otherinfo["view"] = view
	}
	//印象
	db.Table("k_impression").Where("uid=?",uid).Count(&impscount)
	if impscount == 0{
		otherinfo["impression"] = 0
	}else{
		otherinfo["impression"] = impscount
	}
	//是否参加问卷调查
	db.Table("k_user_inv").Where("uid=?",uid).Count(&invcount)
	if invcount == 0{
		otherinfo["inv"] = "未参与问卷"
	}else{
		otherinfo["inv"] = "已参与问卷"
	}
	allinfo["OtherInfo"] = otherinfo

	return allinfo,nil

}