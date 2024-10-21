package routers

import (
	"FShare/controller"
	"FShare/middlewares"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {

	r := gin.Default()
	r.Use(middlewares.Cors())
	r.Static("/dist", "./dist")
	//r.Static("/static", "static")
	//r.LoadHTMLGlob("templates/*")
	r.GET("/myfile", controller.IndexHandler)
	r.GET("/download/:fileName", controller.Download)

	//页面一路由
	v1Group := r.Group("myfile")
	{

		//下载
		v1Group.GET("/download/:fileName/:destNode", controller.DownloadFile)
		//查看一个文件信息
		v1Group.GET("/onefile/:id", controller.GetFileByID)
		//按钮触发更改状态
		v1Group.PUT("/applied/:id", controller.UpdateStatus)
		//更改申请状态
		v1Group.PUT("/update/:id/:applyOwner", controller.UpdateApplyStatus)
		//删除
		v1Group.DELETE("/deleteapply/:id/:applyOwner", controller.DeleteApply)
		v1Group.DELETE("/deletefile/:id", controller.DeleteAFile)
		//添加
		v1Group.POST("/addfile/")
	}

	//页面二路由
	v2Group := r.Group("/upload")
	{
		v2Group.POST("/confirm", controller.UploadFile)
	}

	//页面三路由
	r.GET("/browse", controller.GetFileList)
	v3Group := r.Group("/browse")
	{
		v3Group.POST("/apply", controller.CreateApply)
	}

	//页面四路由
	v4Group := r.Group("verify")
	{
		//返回页面的HTML静态文件
		v4Group.GET("/index", controller.IndexHandlerv4)
		//上传需要进行追溯的文件
		v4Group.POST("/upload2", controller.UploadFileLocal)
		//查询所上传的文件的水印信息
		v4Group.GET("/fingerprint/:type", controller.GetFingerPrint)
		//查询所上传文件在区块链上所保存的信息
		v4Group.GET("/traceback/:txHash/:sourceNode", controller.TraceBackOnChain)
	}

	//页面五路由
	v5Group := r.Group("utility")
	{
		//返回页面的HTML静态页面
		v5Group.GET("/index", controller.IndexHandlerv5)
		//上传需要进行效用分析的文件
		v5Group.POST("/upload3", controller.UploadFileLocal2)
		//对上传的文件进行模型训练和效用分析
		v5Group.GET("/train_analysis", controller.Utility)
		//接受扰动文件并进行模型训练和效用分析
		v5Group.POST("/getfile", controller.GetFile)
		//对扰动文件和原始文件进行模型训练和效用分析
		v5Group.GET("/do_TA/:filename", controller.DOUtility)
		//进行模型下载
		v5Group.GET("/download/:ModelName/:Node", controller.DownloadModelFile)
		//实际的模型下载
		v5Group.GET("download2/:ModelName/:Type", controller.DownloadModel)
	}

	return r
}
