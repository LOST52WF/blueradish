package models

import (
	"bluefoxgo/tools"
	"errors"
)

type K_view_comments struct{
	Id				int64
	View_id			int64
	Comment_content	string
	From_uid		int64
	From_uid_photo	string
	From_name		string
	To_uid			int64
	To_uid_name		string
	Comment_time	int
	Parent_id		int64
	Thumb_up		int
	Up_master		int8
	Layer_master	int8
}




//获取评论及其评论回复
/*
1.首先获取视角的评论(通过view_id,to_uid,parent_id进行判断)
2.再获取评论的回复
 */
 func(comment *K_view_comments)  Comment(vid int64,offset,rows int)([]K_view_comments,[]K_view_comments,error){
	 commentslist := []K_view_comments{}
	 reply 		  := []K_view_comments{}
	 //先获取评论视角的评论
	 res := db.Table("k_view_comments").
			Where("view_id = ? AND parent_id = ?",vid,0).//0代表默认是对视角进行评论
			Select("id,from_uid,from_uid_photo,from_name,comment_content,comment_time,thumb_up").
			Offset(offset).
			Limit(rows).
			Find(&commentslist)
	//再获取回复评论的评论
	 ans := db.Table("k_view_comments").
			Where("view_id = ? AND parent_id <> ?",vid,0). //0代表默认是对视角进行评论
			Select("id,from_uid,from_uid_photo,from_name,comment_content,comment_time,thumb_up,parent_id").
	 		//Group("view_id").
			Find(&reply)
	 if res.Error!=nil || ans.Error!=nil{
	 	return nil,nil,errors.New(tools.ReturnAllError(res.Error,ans.Error))
	 }
	 return commentslist,reply,nil
 }



func(comment *K_view_comments)  DeleteComment(cid int64)error{
	res := db.Delete(&K_view_comments{Id:cid})
	if res.Error!=nil{
		return res.Error
	}
	return nil

}