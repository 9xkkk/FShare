package models

import (
	"FShare/dao"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type File struct {
	FileID      string `json:"id" gorm:"primary_key"`
	Name        string `json:"name"`
	FileOwner   string `json:"fileOwner"`
	Description string `json:"description"`
	Size        string `json:"size"`
	Time        string `json:"time"`
	Hash        string `json:"hash"`
	//Fingerprint string `json:"fingerprint"` //todo: 后续需要将二者加上
	Status int `json:"status"` //1:没被申请；2：正在被申请中；3：申请被拒绝；4：可用不可转发；5：可用可转发
}

type Apply struct {
	ApplyOwner string `json:"applyOwner" gorm:"primary_key"`
	FileOwner  string `json:"fileOwner"`
	Time       string `json:"time"`
	FileID     string `json:"id" gorm:"primary_key"`
	FileName   string `json:"fileName"`
	Hash       string `json:"txHash"`
	Status     int    `json:"status"`
}

type Hashdata struct {
	Result struct {
		Tx struct {
			Execer  string `json:"execer"`
			Payload struct {
				ContentStorage struct {
					Content interface{} `json:"content"`
					Key     string      `json:"key"`
					Op      int         `json:"op"`
					Value   string      `json:"value"`
				} `json:"contentStorage"`
				Ty int `json:"ty"`
			} `json:"payload"`
			RawPayload string `json:"rawPayload"`
			Signature  struct {
				Ty        int    `json:"ty"`
				Pubkey    string `json:"pubkey"`
				Signature string `json:"signature"`
			} `json:"signature"`
			Fee    int    `json:"fee"`
			Feefmt string `json:"feefmt"`
			Expire int    `json:"expire"`
			Nonce  int    `json:"nonce"`
			From   string `json:"from"`
			To     string `json:"to"`
			Hash   string `json:"hash"`
		} `json:"tx"`
	} `json:"result"`
}

var IP = gin.H{
	"A": "124.223.171.19", //王钺程
	"B": "101.43.94.172",  //李炳翰
	"C": "124.221.254.11", //金严
	"D": "124.223.210.53", //叶克炉
	"E": "124.222.196.78", //唐聪
	"F": "10.96.92.7",     //kxq
	"G": "10.96.208.18",   //wyc
}

var Ip2Node = gin.H{
	"124.223.171.19": "A",
	"101.43.94.172":  "B",
	"124.221.254.11": "C",
	"124.223.210.53": "D", //叶克炉
	"124.222.196.78": "E", //唐聪
	"10.96.92.7":     "F", //kxq
	"10.96.208.18":   "G", //wyc
}

var Node string //节点

/*
Todo这个model的增删改查放在这里
*/

func GetHostIp() string {
	conn, err := net.Dial("udp", "8.8.8.8:53")
	if err != nil {
		fmt.Println("get current host ip err: ", err)
		return ""
	}
	addr := conn.LocalAddr().(*net.UDPAddr)
	ip := strings.Split(addr.String(), ":")[0]
	return ip
}

func DownloadFile(context *gin.Context, node, fileName string) (err error) {
	str := fmt.Sprintf("%v", IP[node])
	address := "http://" + str + ":8080/download/" + fileName
	//fmt.Println(str)
	context.Redirect(http.StatusMovedPermanently, address)
	return

}

//func Download(context *gin.Context, fileName string) (err error) {
//	dst := fmt.Sprintf("./%s", fileName) //todo:这里修改文件路径
//	context.Header("Content-Type", "application/octet-stream")
//	context.Header("Content-Disposition", "attachment; filename="+fileName)
//	//context.Header("Content-Transfer-Encoding", "binary")
//	context.File(dst)
//	return
//}

func Download(context *gin.Context, fileName string) (err error) {
	// 构建文件路径
	dst := fmt.Sprintf("./%s", fileName) // 修改为正确的文件路径

	// 打开文件
	file, err := os.Open(dst)
	if err != nil {
		// 处理错误
		context.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return err
	}
	defer file.Close()

	// 设置响应头
	context.Header("Content-Disposition", "attachment; filename="+fileName)
	context.Header("Content-Type", "application/octet-stream")

	// 将文件内容传递给客户端
	_, err = io.Copy(context.Writer, file)
	if err != nil {
		// 处理错误
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return err
	}

	return nil
}

func UploadFiles(context *gin.Context) (err error) {
	var file File

	//file.FileOwner = context.PostForm("fileOwner") //todo：后面改成Node，这里测试不同节点用
	file.FileOwner = Node
	file.Name = context.PostForm("name")
	file.Description = context.PostForm("description")
	file.Size = context.PostForm("size")
	file.Status, _ = strconv.Atoi(context.PostForm("status"))

	f, err := context.FormFile("f1")
	f.Filename = file.Name //+".xml" todo:这里更改文件名可以加xml后缀
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	} else {
		//保存读取的文件到本地服务器
		dst := path.Join("./", f.Filename) //todo:这里修改文件路径
		_ = context.SaveUploadedFile(f, dst)
		context.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
		//生成文件ID
		t := time.Now()
		file.Time = t.Format("2006-01-02 15:04:05")
		timestamp := strconv.FormatInt(t.UTC().UnixNano(), 10)
		randnum := fmt.Sprintf("%04v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(10000))
		file.FileID = Node + timestamp + randnum
		//file.Status = 1
		fileproperties := file.FileID + "#" + file.Name + "#" + file.FileOwner + "#" + file.Description + "#" + file.Size + "#" + file.Time + "#" + strconv.Itoa(file.Status)
		fmt.Println(fileproperties)
		file.Hash = transfer("file", fileproperties)
		//file.Fingerprint = GenertaeFingerPrint(file)
		if err = dao.DB.Create(&file).Error; err != nil {
			return err
		}
	}
	return
}

func CreateApply(file *File) (err error) {
	var apply Apply
	apply.ApplyOwner = Node
	apply.FileOwner = file.FileOwner
	apply.FileID = file.FileID
	apply.FileName = file.Name
	apply.Status = 2
	t := time.Now()
	apply.Time = t.Format("2006-01-02 15:04:05")
	if err = dao.DB.Create(&apply).Error; err != nil {
		return err
	}
	return
}

func GetAllFile() (fileList []*File, err error) {
	if err = dao.DB.Find(&fileList).Error; err != nil {
		return nil, err
	}
	return
}

func GetFileList() (fileList []*File, err error) {
	if err = dao.DB.Where("file_owner = ?", Node).Find(&fileList).Error; err != nil {
		return nil, err
	}
	return
}

func GetMyApply() (applyList []*Apply, err error) {
	if err = dao.DB.Where("apply_owner = ?", Node).Find(&applyList).Error; err != nil {
		return nil, err
	}
	return
}

func GetApplyList() (applyList []*Apply, err error) {
	if err = dao.DB.Where("file_owner = ?", Node).Find(&applyList).Error; err != nil {
		return nil, err
	}
	return
}

func GetFileByID(id string) (file *File, err error) {
	file = new(File)
	if err = dao.DB.Where("file_id = ?", id).First(&file).Error; err != nil {
		return nil, err
	}
	return
}

func GetApply(id, owner string) (apply *Apply, err error) {
	apply = new(Apply)
	if err = dao.DB.Where("file_id = ? and apply_owner = ?", id, owner).First(&apply).Error; err != nil {
		return nil, err
	}
	return
}

func UpdateFile(file *File) (err error) {
	err = dao.DB.Save(file).Error
	if err != nil {
		return err
	}
	return

}

func UpdateApply(apply *Apply) (err error) {
	//todo:需要将更改后的申请记录存入数据库
	t := time.Now()
	apply.Time = t.Format("2006-01-02 15:04:05")
	applyproperties := apply.ApplyOwner + "#" + apply.FileOwner + "#" + apply.Time + "#" + apply.FileID + "#" + strconv.Itoa(apply.Status)
	apply.Hash = transfer("apply", applyproperties)
	err = dao.DB.Save(apply).Error
	if err != nil {
		return err
	}
	return

}

func DeleteAFileByID(id string) (err error) {
	err = dao.DB.Where("file_id=?", id).Delete(&File{}).Error
	if err != nil {
		return err
	}
	return
}

func DeleteApply(id, owner string) (err error) {
	err = dao.DB.Where("file_id=? and apply_owner=?", id, owner).Delete(&Apply{}).Error
	if err != nil {
		return err
	}
	return
}

// 将核验文件保存到核验缓存区
func SaveFilelocal(context *gin.Context) (Filetype string, err error) {
	// 将上传的文件取出来
	f, err := context.FormFile("FileVerify")
	if err != nil {
		context.JSON(http.StatusOK, gin.H{"message": err.Error()})
		return "", err
	}

	// 获取原始文件名和文件后缀
	originalFileName := f.Filename
	fileExtension := filepath.Ext(originalFileName)

	// 构建新的文件名
	newFileName := "verifyfile" + fileExtension

	// 设置文件名称为新的文件名
	f.Filename = newFileName

	// 将上传的文件保存到指定的本地地址，并返回响应
	log.Println(f.Filename)
	dst := fmt.Sprintf("./verifyfile/%s", f.Filename) // 设置核验文件保存的本地地址路径
	if err = context.SaveUploadedFile(f, dst); err != nil {
		return "", err
	}
	return fileExtension, nil
}

// 获取上传的核验文件的路径
func GetVerifyFile(filetype string) (FilePath string, err error) {
	dirPth := fmt.Sprintf("./verifyfile")
	fis, err := ioutil.ReadDir(filepath.Clean(filepath.ToSlash(dirPth)))
	if err != nil {
		return "", err
	}

	for _, f := range fis {
		_path := filepath.Join(dirPth, f.Name())
		// 指定格式
		if filepath.Ext(f.Name()) == filetype {
			FilePath = _path
			break // 一旦找到匹配的文件，就退出循环
		}
	}
	return FilePath, nil
}

// 获取上传的核验文件的哈希值
func FindtxHash(fingerprint string) (file *File, err error) {
	file = new(File)
	if err = dao.DB.Where("fingerprint=?", fingerprint).First(&file).Error; err != nil {
		return nil, err
	}
	return
}

func AnalyzeData(txHash string) (Hashdata, error) {
	encodeddata := queryTx(txHash)
	//定义结构体以匹配JSON数据
	var decodedata Hashdata
	err := json.Unmarshal(encodeddata, &decodedata)
	if err != nil {
		return decodedata, err
	}
	return decodedata, nil
}

func TraceBackOnChain(txHash string) ([]Hashdata, Hashdata, error) {
	//todo: 传hash值，进行查询文件信息，进行错误验证。根据文件id查询申请哈希，得到申请hash，用一个数组存储循环查询申请记录
	//1.通过文件hash，查询文件信息,取出文件ID
	//定义结构体
	var filedata Hashdata
	var applydatalist []Hashdata
	//进行查询文件信息
	filedata, err := AnalyzeData(txHash)
	if err != nil {
		fmt.Println("JSON解析错误:", err)
		return applydatalist, filedata, err
	}
	//2.取出文件ID,根据文件ID查询申请哈希，得到申请哈希，用一个数组存储循环查询申请记录
	filemessage := filedata.Result.Tx.Payload.ContentStorage.Value
	filemessages := strings.Split(filemessage, "#")
	fileID := filemessages[6]
	applylist := new([]Apply)
	if err = dao.DB.Where("file_id=?", fileID).First(&applylist).Error; err != nil {
		return applydatalist, filedata, err
	}
	//循环查询申请记录
	for i := range *applylist {
		apply := (*applylist)[i]
		applydata, _ := AnalyzeData(apply.Hash)
		applydatalist = append(applydatalist, applydata)
	}
	return applydatalist, filedata, err
}
