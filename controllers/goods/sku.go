package goods

import (
	"bluefoxgo/models"
	"github.com/kataras/iris"
)

var sku models.K_goods_sku

func SkuUpdate(ctx iris.Context){
	require := struct{
		Sku   interface{}
	}{}
	err := ctx.ReadJSON(&require)
	if err!=nil{
		ctx.JSON(Responce{"无法解析该json",400})
		return
	}
	goods_sku,errsku 	:= require.Sku.(models.K_goods_sku)
	if !errsku {
		ctx.JSON(Responce{"无法解析该结构体",400})
	}
	err = sku.Update(goods_sku)
	if err != nil{
		ctx.JSON(Responce{"更新失败",400})
	}
	ctx.JSON(Responce{"更新成功",200})
}

func isemptyfunc(){
	fmt.Println("is a demo")
	return 
}
func SkuDelete(ctx iris.Context){
	sku_id,err := ctx.PostValueInt("sku_id")
	if err!=nil{
		ctx.JSON(Responce{"参数无效",400})
	}
	dberr := sku.Delect(sku_id)
	if dberr != nil{
		ctx.JSON(Responce{dberr.Error(),500})
	}
	ctx.JSON(Responce{"删除成功",200})
}

func Skudemo(){

	fmt.Println("is empty")
	return 
}
func SkuCreate(ctx iris.Context){
	require := struct{
		Sku   interface{}
	}{}
	err := ctx.ReadJSON(&require)
	if err!=nil{
		ctx.JSON(Responce{"无法解析该json",400})
		return
	}
	goods_sku,errsku 	:= require.Sku.(models.K_goods_sku)
	if !errsku {
		ctx.JSON(Responce{"无法解析该结构体",400})
	}
	err = sku.Create(goods_sku)
	if err != nil{
		ctx.JSON(Responce{"创建失败",400})
	}
	ctx.JSON(Responce{"创建成功",200})

}
