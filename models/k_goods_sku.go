package models

type K_goods_sku struct {
	Id				int
	Color_pic	 	string
	Goods_id		int
	Sort_id			string
	Sku_main		string
	Sku_en			string
	Sku_cn			string
	Color_cn		string
	Sku_pic			string
	Price			float64
	Texture			string
	Net_weight		float64
	Net_volume		float64
	Color_hex		string
	Color_rgb		string
	Image_uri		string
	Hot_level		int
	Collect_total	int
	Sales_total		int
	Status_sale		int8
	Flag_new		int8
	Flag_hot		int8
	Official_tag	string
	Total_weight	float64
	Moisture		int8
	Glossiness		int8
	Persistance		int8
	Chromaticity	int8
	Coverage		int8
	Impress_tags	string
	Record_no		string
}

func(sku *K_goods_sku)	Delect(sku_id int)error{
	//goods表和 sku表 都没有软输出  所以直接删除
	goods_sku := K_goods_sku{Id:sku_id}
	res := db.Delete(&goods_sku)
	return res.Error
}

func(sku *K_goods_sku)	Create(goods_sku K_goods_sku)error{
	res := db.Create(&goods_sku)
	if res.Error != nil{
		return res.Error
	}
	return nil
}

func(sku *K_goods_sku)	Update(goods_sku K_goods_sku)error{
	sku_id := K_goods_sku{Id:goods_sku.Id}
	res := db.Model(&sku_id).
		Updates(map[string]interface{}{
		"sku_cn"	:	goods_sku.Sku_cn,  		//商品名
		"sku_main"	:   goods_sku.Sku_main, 	//色系
		"price"		:	goods_sku.Price,	  	//价格
		"texture"	:	goods_sku.Texture,  	//质地
		"net_weight":	goods_sku.Net_weight,	//净含量
		"net_volume":	goods_sku.Net_volume,	//净容量
		"total_weight":goods_sku.Total_weight,	//总重量
		"tolor_hex":	goods_sku.Color_hex, 	//颜色hex
		"color_rgb":	goods_sku.Color_rgb, 	//颜色RGB
		"flag_new":	goods_sku.Flag_new,		//是否为新品
		"flag_hot"	:	goods_sku.Flag_hot, 	//平台以外的来源
		"tatus_sale":	goods_sku.Status_sale, 	//商品状态
		"official_tag":goods_sku.Official_tag, //官方标签
		})
	if res.Error != nil{
		return res.Error
	}
	return nil
}