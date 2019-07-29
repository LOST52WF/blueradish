package models

import (
	"errors"
	"github.com/jinzhu/gorm"
)

type K_view_pic struct{
	Id			int
	View_id		int64
	Image_uri	string
	Header_light int8
	Footer_light int8
}


func(pic *K_view_pic) UploadOnePic(vid int64,view_pic K_view_pic)error{
	tx := db.Begin()
	// 注意，一旦你在一个事务中，使用tx作为数据库句柄

	if err := tx.Create(&view_pic).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Model(&K_view{Id:vid}).UpdateColumn("pic_total", gorm.Expr("pic_total + ?", 1)).
	Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}


func(pic *K_view_pic) DeleteOnePic(vid int64,pid int)error {

	tx := db.Begin()
	// 注意，一旦你在一个事务中，使用tx作为数据库句柄

	if err := tx.Where("id = ?", pid).Delete(K_view_pic{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Model(&K_view{Id: vid}).UpdateColumn("pic_total", gorm.Expr("pic_total - ?", 1)).
		Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}






func(pic *K_view_pic) UpdateAllViewPic(vid int64,image_uri []K_view_pic)error{
	//total_pic:= len(image_uri)
	//a.先要进行判断把image_uri不存在的而数据库中存在的进行删除
	//b.其次删选出iamge_uri存在而数据库中不存在的进行 K_view_pic

	// a.1  先获取所有此view_id的图片,此步出问题直接返回err
	all_image := []K_view_pic{}
	res := db.Table("k_view_pic").
		Select("id,image_uri").
		Where("view_id= ?",vid).
		Find(&all_image)
	if res != nil{
		return errors.New("无法获取原本视角图片信息")
	}

	// 筛选出需要删除的图片和需要添加的图片
	need_to_delete_id := needtodeleteid(all_image,image_uri)
	need_to_insert_su := needtoinsert(all_image,image_uri)
	if len(need_to_delete_id)==0 && len(need_to_insert_su)!=0 {
		//需要删除的为0，直接进行添加
		for _,insert_one := range need_to_insert_su{ //gorm不提供多项插入
			db.Create(&insert_one)
		}  //不好使用事务
		db.Model(&K_view{Id:vid}).Update("pic_total",len(need_to_insert_su))

	}else if  len(need_to_delete_id) !=0 && len(need_to_insert_su)==0 {
		//只需要删除
		tx := db.Begin()
		// 注意，一旦你在一个事务中，使用tx作为数据库句柄
		del_res := db.Where("id IN (?)",need_to_delete_id).Delete(K_view_pic{})
		if del_res.Error!=nil{
			tx.Rollback()
		}
		up_res := db.Model(&K_view{Id:vid}).Update("pic_total",len(all_image)-len(need_to_delete_id))
		if up_res.Error!=nil{
			tx.Rollback()
		}
		tx.Commit()

	}else if len(need_to_delete_id) !=0 && len(need_to_insert_su)!=0 {
		//需要删除且添加
		db.Where("id IN (?)",need_to_delete_id).Delete(K_view_pic{})
		for _,insert_one := range need_to_insert_su{ //gorm不提供多项插入
			db.Create(&insert_one)
		}
		db.Model(&K_view{Id:vid}).Update("pic_total",len(image_uri))
	}

	return nil  //不需要添加和删除，直接返回
}


func needtodeleteid(all_image,image_uri []K_view_pic)[]int{
	//copy := make([]K_view_pic,len(all_image))
	//copy = append(copy,all_image...)
	if len(all_image) == 0{   //原本视角并不存在图片 ，可以直接进行添加
		return []int{}
	}
	all_image_map := make(map[int]K_view_pic,len(all_image))
	for _,pic := range all_image{
		all_image_map[pic.Id] = pic
	}
	for  key,v := range all_image_map {  //all_image在外循环，是获取需要删除的图片
		for _,image := range image_uri{
			if v.Image_uri == image.Image_uri {   //不能用ID  存在新上传的图片没有ID等其他字段 ，只有image_uri
				delete(all_image_map,key) //同时存在将map中删除，剩下的就是需要删除的
			}
		}
	}
	need_delete_id := []int{}
	for _,v := range all_image_map{
		need_delete_id = append(need_delete_id,v.Id)
	}
	return need_delete_id
}

func needtoinsert(all_image,image_uri []K_view_pic)[]K_view_pic{
	image_uri_map := make(map[string]K_view_pic,len(image_uri))
	for _,pic := range image_uri{
		image_uri_map[pic.Image_uri] = pic
	}
	for  key,v := range image_uri_map {  //image_uri在外循环，是获取需要加入的图片
		for _,image := range all_image{
			if v.Image_uri == image.Image_uri {
				delete(image_uri_map,key) //同时存在将map中删除，剩下的就是需要插入的
			}
		}
	}
	need_insert := []K_view_pic{}
	for _,v := range image_uri_map{
		need_insert = append(need_insert,v)
	}
	return need_insert
}




