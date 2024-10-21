package controller

import (
	"FShare/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
)

/*
url     --> controller --> logic   --> model
请求来了 --> 控制器       --> 业务逻辑 --> 模型层的增删改查
*/

func IndexHandler(context *gin.Context) {
	//context.HTML(http.StatusOK, "index.html", nil)
	//这里更改节点名称，服务器分别为ABCDE，用models.Node

	myApply, err := models.GetMyApply()
	if err != nil {
		context.JSON(http.StatusOK, gin.H{"get my apply error": err.Error()})
	}
	applyList, err := models.GetApplyList()
	if err != nil {
		context.JSON(http.StatusOK, gin.H{"get apply list error": err.Error()})
	}
	fileList, err := models.GetFileList()
	if err != nil {
		context.JSON(http.StatusOK, gin.H{"get file list error": err.Error()})
	}
	context.JSON(http.StatusOK, gin.H{
		"myapply":         myApply,
		"applylist":       applyList,
		"filelist":        fileList,
		"NodeInformation": models.Node,
	})
}

func DownloadFile(context *gin.Context) {
	filename, _ := context.Params.Get("fileName")
	node, _ := context.Params.Get("destNode")
	if err := models.DownloadFile(context, node, filename); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		context.JSON(http.StatusOK, gin.H{"file": "success download"})
	}
}

func Download(context *gin.Context) {
	filename, _ := context.Params.Get("fileName")
	if err := models.Download(context, filename); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		context.JSON(http.StatusOK, gin.H{"file": "success download"})
	}
}

func IndexHandlerv4(context *gin.Context) {
	context.HTML(http.StatusOK, "verify.html", nil)
}

func IndexHandlerv5(context *gin.Context) {
	context.HTML(http.StatusOK, "utility.html", nil)
}

func UploadFile(context *gin.Context) {
	//var file models.File
	//context.BindJSON(&file)
	if err := models.UploadFiles(context); err != nil {
		context.JSON(http.StatusOK, gin.H{"error": err.Error()})
	} else {
		context.JSON(http.StatusOK, gin.H{"file": "success"})
	}
}

func CreateApply(context *gin.Context) {
	//1.从请求中把数据取出来
	//var apply models.Apply
	var file models.File
	_ = context.BindJSON(&file)
	//2.存入数据库
	//3.返回响应
	if err := models.CreateApply(&file); err != nil {
		context.JSON(http.StatusOK, gin.H{"error": err.Error()})
	} else {
		context.JSON(http.StatusOK, file)
	}
}

func GetFileList(context *gin.Context) {
	fileList, err := models.GetAllFile()
	if err != nil {
		context.JSON(http.StatusOK, gin.H{"error": err.Error()})
	} else {
		context.JSON(http.StatusOK, fileList)
	}
}

func GetFileByID(context *gin.Context) {
	id, ok := context.Params.Get("id")
	if !ok {
		context.JSON(http.StatusOK, gin.H{"error": "id not exist"})
		return
	}
	//查询数据库是否有这个id
	file, err := models.GetFileByID(id)
	if err != nil {
		context.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	} else {
		context.JSON(http.StatusOK, file)
	}
}

func UpdateStatus(context *gin.Context) {
	//拿到请求里的id
	id, ok := context.Params.Get("id")
	if !ok {
		context.JSON(http.StatusOK, gin.H{"error": "id not exist"})
		return
	}
	//查询数据库是否有这个id
	file, err := models.GetFileByID(id)
	if err != nil {
		context.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	//放入todo变量
	_ = context.BindJSON(&file)
	//新信息保存到数据库
	err = models.UpdateFile(file)
	if err != nil {
		context.JSON(http.StatusOK, gin.H{"error": err.Error()})
	} else {
		context.JSON(http.StatusOK, file)
	}
}

func UpdateApplyStatus(context *gin.Context) {
	//拿到请求里的id
	id, ok := context.Params.Get("id")
	if !ok {
		context.JSON(http.StatusOK, gin.H{"error": "id not exist"})
		return
	}
	owner, ok := context.Params.Get("applyOwner")
	if !ok {
		context.JSON(http.StatusOK, gin.H{"error": "owner not exist"})
		return
	}
	//查询数据库是否有这个id
	apply, err := models.GetApply(id, owner)
	if err != nil {
		context.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	//放入变量
	_ = context.BindJSON(&apply)
	applyrecord := new(models.Applyrecord)
	_ = context.BindJSON(&applyrecord)
	//新信息保存到数据库
	err = models.UpdateApply(apply, applyrecord)
	if err != nil {
		context.JSON(http.StatusOK, gin.H{"error": err.Error()})
	} else {
		context.JSON(http.StatusOK, apply)
	}
}

func DeleteApply(context *gin.Context) {
	id, ok := context.Params.Get("id")
	if !ok {
		context.JSON(http.StatusOK, gin.H{"error": "id not exist"})
		return
	}
	owner, ok := context.Params.Get("applyOwner")
	if !ok {
		context.JSON(http.StatusOK, gin.H{"error": "id not exist"})
		return
	}
	//查询数据库是否有这个id并删除
	if err := models.DeleteApply(id, owner); err != nil {
		context.JSON(http.StatusOK, gin.H{"error": err.Error()})
	} else {
		context.JSON(http.StatusOK, gin.H{id: "deleted"})
	}
}

func DeleteAFile(context *gin.Context) {
	id, ok := context.Params.Get("id")
	if !ok {
		context.JSON(http.StatusOK, gin.H{"error": "id not exist"})
		return
	}
	//查询数据库是否有这个id并删除
	if err := models.DeleteAFileByID(id); err != nil {
		context.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
			"file":  id,
		})
	} else {
		context.JSON(http.StatusOK, gin.H{id: "deleted"})
	}
}

// UploadFileLocal 将上传的文件保存到本地用于校验
func UploadFileLocal(context *gin.Context) {
	if Filetype, err := models.SaveFilelocal(context); err != nil {
		context.JSON(http.StatusOK, gin.H{"error": err.Error()})
	} else {
		context.JSON(http.StatusOK, gin.H{
			"status": "Upload success",
			"type":   Filetype,
		})
	}
}

// GetFingerPrint 提取文件水印
func GetFingerPrint(context *gin.Context) {
	filetype, ok := context.Params.Get("type")
	if !ok {
		context.JSON(http.StatusOK, gin.H{"error": "type of file is not exist"})
		return
	}
	//根据文件类型查找到文件的具体路径
	filepath, err := models.GetVerifyFile(filetype)
	if err != nil {
		context.JSON(http.StatusOK, gin.H{
			"error":    err.Error(),
			"filetype": filetype,
		})
	} else {
		TxHash, fingerprint, fileId, err := models.ExtractFingerPrint(filepath)
		fmt.Println(filepath, fingerprint, fileId)
		if err != nil {
			context.JSON(http.StatusOK, gin.H{"error": err.Error()})
		}
		context.JSON(http.StatusOK, gin.H{
			"fingerprint": fingerprint,
			"filetxHash":  TxHash,
		})
	}
}

// TraceBackOnChain 提取区块链上文件的信息
func TraceBackOnChain(context *gin.Context) {
	//首先需要提取出哈希值
	txHash, ok := context.Params.Get("txHash")
	if !ok {
		context.JSON(http.StatusOK, gin.H{"error": "txHash is not exist"})
		return
	}
	//取出获取该文件的节点信息
	sourceNode, ok := context.Params.Get("sourceNode")
	if !ok {
		context.JSON(http.StatusOK, gin.H{"error": "source error"})
		return
	}
	sourceNode = "G"
	//传入文件区块链哈希
	if applydatalist, filedata, checkdata, err := models.TraceBackOnChain(txHash, sourceNode); err != nil {
		context.JSON(http.StatusOK, gin.H{"error": err.Error()})
	} else {
		context.JSON(http.StatusOK, gin.H{
			"status":        "Trace back success",
			"checkdata":     checkdata,
			"filedata":      filedata,
			"applydatalist": applydatalist,
		})
	}
}

// UploadFileLocal 将上传的文件保存到本地用于效用分析
func UploadFileLocal2(context *gin.Context) {
	if _, err := models.SaveFilelocal2(context); err != nil {
		context.JSON(http.StatusOK, gin.H{"error": err.Error()})
	} else {
		context.JSON(http.StatusOK, gin.H{
			"status": "Upload success",
		})
	}
}

// Utility 进行效用分析
func Utility(context *gin.Context) {
	dirPth := fmt.Sprintf("./analysisfile/analysisfile.csv")
	_, _, fileid, err := models.ExtractFingerPrint(dirPth)
	//if err != nil {
	//	context.JSON(http.StatusOK, gin.H{"error": err.Error()})
	//}

	// 验证文件是否含有水印, 如果有水印则是扰动数据, 如果没有水印就是原始数据
	if fileid != "" {
		if err = models.Utility(context, fileid); err != nil {
			context.JSON(http.StatusOK, gin.H{"error": err.Error()})
		} else {
			context.JSON(http.StatusOK, gin.H{"file": "success utility"})
		}
	} else {
		if Accuracy, W_Precision, W_Recall, W_F1, Node, err := models.DOUtility2(context); err != nil {
			context.JSON(http.StatusOK, gin.H{"error": err.Error()})
		} else {
			context.JSON(http.StatusOK, gin.H{
				"Accuracy":    Accuracy,
				"W_Precision": W_Precision,
				"W_Recall":    W_Recall,
				"W_F1":        W_F1,
				"Node":        Node,
			})
		}
	}
}

// 对扰动文件和原始文件进行模型训练和精度对比
func DOUtility(context *gin.Context) {
	filename, _ := context.Params.Get("filename")
	if accuracy, w_precision, w_recall, w_f1, node, err := models.DOUtility(context, filename); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		context.JSON(http.StatusOK, gin.H{
			"Accuracy":    accuracy,
			"W_Precision": w_precision,
			"W_Recall":    w_recall,
			"W_F1":        w_f1,
			"Node":        node,
		})
	}
}

func GetFile(context *gin.Context) {
	// 从表单中获取文件
	file, err := context.FormFile("perturbed_file")
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "无法获取文件"})
		return
	}

	// 获取原始文件名和文件后缀
	originalFileName := file.Filename
	fileExtension := filepath.Ext(originalFileName)

	// 构建新的文件名
	newFileName := "analysisfile_FP" + fileExtension

	// 设置保存路径
	dst := fmt.Sprintf("./analysisfile/%s", newFileName)

	// 保存文件
	if err := context.SaveUploadedFile(file, dst); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "文件写入失败"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "文件已成功上传", "filename": newFileName})
}

func DownloadModelFile(context *gin.Context) {
	modelname, _ := context.Params.Get("ModelName")
	node, _ := context.Params.Get("Node")
	if node == "Self" {
		if err := models.DownloadModel(context, modelname, ""); err != nil {
			context.JSON(http.StatusOK, gin.H{"error": err.Error()})
		}
	} else {
		if err := models.DownloadModelFile(context, node, modelname); err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
	}

	context.JSON(http.StatusOK, gin.H{"file": "success download"})
}

func DownloadModel(context *gin.Context) {
	modelname, _ := context.Params.Get("ModelName")
	Type, _ := context.Params.Get("Type")
	if err := models.DownloadModel(context, modelname, Type); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		context.JSON(http.StatusOK, gin.H{"file": "success download"})
	}
}
