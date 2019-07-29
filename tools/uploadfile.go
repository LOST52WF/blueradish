package tools

import (
	"errors"
	"fmt"
	"github.com/akkuman/parseConfig"
	"github.com/aliyun-oss-go-sdk/oss"
	"github.com/kataras/iris"
	"mime/multipart"
	"os"
	"path/filepath"
)
//var OC *oss.Client  //实例化对象
var bucketName	string
var uploadPath	string
var osscli *oss.Client

func  NewOSS()(error){
	var config = parseConfig.New("./../config/aliyun_oss.json")

	 endpoint_inf := config.Get("endpoint")	// 此为interface{}格式数据
	 endpoint := endpoint_inf.(string)				// 断言,获取参数endpoint

	 accessKeyId_inf := config.Get("accessKeyId") //获取accessKeyId
	 accessKeyId := accessKeyId_inf.(string)

	 accessKeySecret_inf := config.Get("accessKeySecret_inf")
	 accessKeySecret := accessKeySecret_inf.(string)

	bucketName = config.Get("bucketName").(string)
	uploadPath = config.Get("uploadPath").(string)

	clien, err := oss.New(endpoint,accessKeyId,accessKeySecret)
	if err != nil{
		osscli = nil
		return err
	}
	osscli = clien
	return nil
}

func UploadFileToServer(ctx iris.Context)(int,[]string,error){
	form := ctx.Request().MultipartForm
	files := form.File["files[]"]
	failures := 0
	savepath := []string{}
	for _, file := range files {
		err, savename := saveUploadedFile(file, "./../upload/photoimage")
		if err != nil {
			failures++
			savepath = append(savepath,"")
		}
		savepath = append(savepath,savename)
	}
	if failures != 0{
		return  len(files)-failures,savepath,errors.New(fmt.Sprintf("%d files can't upload",failures))
	}
	return 0,savepath,nil
}

func saveUploadedFile(fh *multipart.FileHeader, destDirectory string) (error,string,) {
	src, err := fh.Open()
	if err != nil {
		return  err,""
	}
	defer src.Close()
	savename := filepath.Join(destDirectory, fh.Filename)
	out, err := os.OpenFile(savename, os.O_WRONLY|os.O_CREATE, os.FileMode(0666))
	if err != nil {
		return err,""
	}
	defer out.Close()
	return err,savename
}


func UploadToOSSApi(ctx iris.Context,filePath [] string, dirName string)([]string,error){

	var fileNames []string
	//var cloudsavpath []string
	imagesChannel := make(chan string)
	if osscli == nil {   //当实例为空时时，创建新实例
		NewOSS()
	}
	// 获取存储空间。
	bucket, err := osscli.Bucket(bucketName)
	if err != nil {
		fmt.Println("Error:", err)   //获取失败
	}
	if len(filePath) > 1 {
		go func() {
			for _, path := range filePath {
				objectName := filepath.Base(path)
				// 上传文件
				err = bucket.PutObjectFromFile(dirName+ "/"+ objectName, path)
				fmt.Println(err)
				if err == nil {
					imagesChannel <- "https://" + bucketName + ".oss-cn-hangzhou.aliyuncs.com/" +dirName+ "/"+ objectName
				}
			}
			close(imagesChannel)
		}()
	} /*else {
		objectName := filepath.Base(filePath[0])
		go func() {
			// 上传文件。
			err = bucket.PutObjectFromFile(dirName+ "/"+objectName,filePath[0])
			//fmt.Println(err)
			if err == nil {
				imagesChannel <- "https://" + bucketName + ".oss-cn-hangzhou.aliyuncs.com/"+dirName+ "/"+objectName
			}
			close(imagesChannel)
		}()
	}*/
	for paths := range imagesChannel {
		//cloudsavpath
		fileNames = append(fileNames, paths)
	}
	fmt.Println(len(fileNames))
	if len(fileNames) > 0{
		return fileNames,nil
		//utils.ReturnResult(ctx,0,"上传OSS成功",fileNames,201)
	}else {
		return fileNames,errors.New("上传失败")
		//utils.ReturnResult(ctx,1,"上传OSS失败","",400)
	}
}
