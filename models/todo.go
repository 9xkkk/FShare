package models

import (
	"FShare/dao"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// 写死用于监管机构进行展示
type User struct {
	UserID          string `json:"userID" gorm:"primary_key"`
	LoginAccount    string `json:"loginAccount"`
	Password        string `json:"password"`
	UserName        string `json:"userName"`
	UserEmail       string `json:"userEmail"`
	Time            string `json:"time"`
	UserFingerPrint string `json:"userFingerPrint"`
}

type File struct {
	FileID      string `json:"id" gorm:"primary_key"`
	Name        string `json:"name"`
	FileOwner   string `json:"fileOwner" gorm:"primary_key"`
	Description string `json:"description"`
	Size        string `json:"size"`
	Time        string `json:"time"`
	Hash        string `json:"hash"`
	Status      int    `json:"status"` //1:没被申请；2：正在被申请中；3：申请被拒绝；4：可用不可转发；5：可用可转发
}

type Apply struct {
	ApplyOwner    string `json:"applyOwner" gorm:"primary_key"`
	FileOwner     string `json:"fileOwner" gorm:"primary_key"`
	Time          string `json:"time"`
	FileID        string `json:"id" gorm:"primary_key"`
	FileName      string `json:"fileName"`
	Hash          string `json:"txHash"`
	Status        int    `json:"status"`
	FingerPrint   string `json:"fingerPrint"`
	PrivacyBudget string `json:"privacyBudget"`
	IsHandled     bool   `json:"isHandled"`
}

type Applyrecord struct {
	FileID      string `json:"id"`
	Hash        string `json:"txHash" gorm:"primary_key"`
	FingerPrint string `json:"fingerprint"`
	Status      string `json:"status"`
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
	"A": "124.223.171.19",
	"B": "101.43.94.172",
	"C": "124.221.254.11",
	"D": "124.223.210.53",
	"E": "124.222.196.78",
}

var Ip2Node = gin.H{
	"124.223.171.19": "A",
	"101.43.94.172":  "B",
	"124.221.254.11": "C",
	"124.223.210.53": "D",
	"124.222.196.78": "E",
}

var Node string //节点

/*
Todo这个model的增删改查放在这里
*/

func GetHostIp() string {
	//conn, err := net.Dial("udp", "8.8.8.8:53")
	//if err != nil {
	//	fmt.Println("get current host ip err: ", err)
	//	return ""
	//}
	//addr := conn.LocalAddr().(*net.UDPAddr)
	//ip := strings.Split(addr.String(), ":")[0]
	//return ip
	cmd := exec.Command("curl", "myexternalip.com/raw") // 创建一个命令对象，用于获取公网 IP
	output, _ := cmd.Output()
	ip := string(output)
	return ip
}

func DownloadFile(context *gin.Context, node, fileName, applyOwner string) (err error) {
	str := fmt.Sprintf("%v", IP[node])
	address := "http://" + str + ":8080/download/" + fileName + "/" + applyOwner
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

func Download(context *gin.Context, fileName, applyOwner string) (err error) {
	// 构建文件路径
	var fN = strings.Split(fileName, ".")

	dst := fmt.Sprintf("./csvfile/%s_%s_FP.csv", fN[0], applyOwner) // 修改为正确的文件路径

	// 打开文件
	file, err := os.Open(dst)
	if err != nil {
		// 处理错误
		context.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return err
	}
	defer file.Close()

	// 设置响应头
	context.Header("Content-Disposition", "attachment; filename="+fN[0]+"_FP.csv")
	context.Header("Content-Type", "application/octet-stream")
	context.File(dst)
	// 将文件内容传递给客户端
	//_, err = io.Copy(context.Writer, file)
	//if err != nil {
	//	// 处理错误
	//	context.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
	//	return err
	//}

	return nil
}

// 下载生成权限标识后的文件
func DownloadTransformedFile(context *gin.Context, fileName string) (err error) {
	// 构建文件路径
	var fN = strings.Split(fileName, ".")

	dst := fmt.Sprintf("./csvfile/%s_FP.csv", fN[0]) // 修改为正确的文件路径

	// 打开文件
	file, err := os.Open(dst)
	if err != nil {
		// 处理错误
		context.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return err
	}
	defer file.Close()

	// 设置响应头
	context.Header("Content-Disposition", "attachment; filename="+fN[0]+"_FP.csv")
	context.Header("Content-Type", "application/octet-stream")
	context.File(dst)
	return nil
}

// 下载本地文件
func DownloadLocalFile(context *gin.Context, fileName string) (err error) {
	// 构建文件路径
	var fN = strings.Split(fileName, ".")

	dst := fmt.Sprintf("./csvfile/%s.csv", fN[0]) // 修改为正确的文件路径

	// 打开文件
	file, err := os.Open(dst)
	if err != nil {
		// 处理错误
		context.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return err
	}
	defer file.Close()

	// 设置响应头
	context.Header("Content-Disposition", "attachment; filename="+fN[0]+"_FP.csv")
	context.Header("Content-Type", "application/octet-stream")
	context.File(dst)
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
		dst := path.Join("./csvfile", f.Filename) //todo:这里修改文件路径
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
	apply.FingerPrint = ""
	apply.IsHandled = false
	t := time.Now()
	apply.Time = t.Format("2006-01-02 15:04:05")
	if err = dao.DB.Create(&apply).Error; err != nil {
		return err
	}
	return
}

func GetAllFile() (fileList []*File, err error) {
	if err = dao.DB.Where("file_owner <> ?", Node).Find(&fileList).Error; err != nil {
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

func GetFileByID(id, fileOwner string) (file *File, err error) {
	file = new(File)
	if err = dao.DB.Where("file_id = ? and file_owner = ?", id, fileOwner).First(&file).Error; err != nil {
		return nil, err
	}
	return
}

func GetApply(id, owner, fileOwner string) (apply *Apply, err error) {
	apply = new(Apply)
	if err = dao.DB.Where("file_id = ? and apply_owner = ? and file_owner = ?", id, owner, fileOwner).First(&apply).Error; err != nil {
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

func EmbedFingerprint(applyOwner, applyHash, epsilon, fileName string) (string, error) {
	//cmd := exec.Command("venv\\Scripts\\python.exe", "python/embed.py", fileName, applyHash, epsilon, applyOwner)
	cmd := exec.Command("./venv/bin/python", "python/embed.py", fileName, applyHash, epsilon, applyOwner)
	//cmd := exec.Command("/usr/bin/python", "python/embed.py", fileName, applyHash, fingerprint, epsilon)

	output, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
	}
	result := string(output)
	if result == "false" {
		return "", errors.New("embed error")
	}

	print(result)
	return result, nil
}

func GenerateFingerprint(applyHash string) (string, error) {
	//cmd := exec.Command("venv\\Scripts\\python.exe", "python/embed.py", applyHash)
	cmd := exec.Command("./venv/bin/python", "python/generate.py", applyHash)
	//cmd := exec.Command("/usr/bin/python", "python/embed.py", fileName, applyHash)

	output, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
	}
	result := string(output)
	if result == "false" {
		return "", errors.New("generate error")
	}

	print(result)
	return result, nil
}

func UpdateApply(apply *Apply, applyrecord *Applyrecord) (err error) {

	t := time.Now()
	apply.Time = t.Format("2006-01-02 15:04:05")
	applyproperties := apply.ApplyOwner + "#" + apply.FileOwner + "#" + apply.Time + "#" + apply.FileID + "#" + strconv.Itoa(apply.Status)
	apply.Hash = transfer("apply", applyproperties)
	//apply.IsHandled = true

	//todo:根据status=4的时候，可用不可转发
	FP, err := GenerateFingerprint(apply.Hash)
	if err != nil {
		return err
	}
	apply.FingerPrint = FP
	err = dao.DB.Save(apply).Error
	if err != nil {
		return err
	}
	applyrecord.FingerPrint = FP
	applyrecord.FileID = apply.FileID
	applyrecord.Hash = apply.Hash
	applyrecord.Status = strconv.Itoa(apply.Status)
	err = dao.DB.Save(applyrecord).Error
	if err != nil {
		return err
	}
	return

}

func UpdataPrivacyBudget(apply *Apply, epsilon string) (err error) {

	apply.IsHandled = true
	apply.PrivacyBudget = epsilon
	_, err = EmbedFingerprint(apply.ApplyOwner, apply.Hash, apply.PrivacyBudget, apply.FileName)
	if err != nil {
		return err
	}
	err = dao.DB.Save(apply).Error
	if err != nil {
		return err
	}
	return

}

func FileIsExisted(filename string) bool {
	existed := true
	if _, err := os.Stat("csvfile/" + filename); os.IsNotExist(err) {
		existed = false
	}
	return existed
}

func DeleteAFileByID(id, fileOwner string) (err error) {
	var file File
	//删除数据库条目
	err = dao.DB.Where("file_id=? and file_owner = ?", id, fileOwner).Find(&file).Error
	if err != nil {
		return err
	}
	//删除本地源文件
	err = os.Remove("csvfile/" + file.Name)
	if err != nil {
		return err
	}
	//删除嵌入了水印的文件
	fileFP := strings.Split(file.Name, ".")
	Node_delete := [...]string{"A", "B", "C", "D", "E"}
	for _, node := range Node_delete {
		fileFP_Name := fileFP[0] + "_" + node + "_FP.csv"
		if FileIsExisted(fileFP_Name) == true {
			fmt.Println("delete success")
			err = os.Remove("csvfile/" + fileFP_Name)
			if err != nil {
				return err
			}
		}
	}

	err = dao.DB.Where("file_id=? and file_owner=?", id, fileOwner).Delete(&file).Error
	if err != nil {
		return err
	}
	return
}

func DeleteApply(id, owner, fileOwner string) (err error) {
	err = dao.DB.Where("file_id=? and apply_owner=? and file_owner = ?", id, owner, fileOwner).Delete(&Apply{}).Error
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

func ExtractFingerPrint(filePath, epsilon string) (string, string, string, error) {
	fmt.Println(filePath)
	//cmd := exec.Command("D:\\Reaserch\\System development\\project\\FShare\\venv\\Scripts\\python.exe", "python/extract.py", filePath, epsilon)
	cmd := exec.Command("./venv/bin/python", "python/extract.py", filePath, epsilon)
	//cmd := exec.Command("/usr/bin/python", "python/extract.py", filePath)

	output, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
	}
	result := string(output)
	if result == "false" {
		return "", "", "", errors.New("embed error")
	}
	result = strings.Replace(result, "\n", "", -1)
	result = result + "\n"
	fmt.Print(result)
	applyrecord := new(Applyrecord)
	if err = dao.DB.Where("finger_print = ?", result).First(&applyrecord).Error; err != nil {
		println("DB error:", err.Error())
		return applyrecord.Hash, result, applyrecord.FileID, err
	}
	return applyrecord.Hash, result, applyrecord.FileID, nil
}

// 获取上传的核验文件的哈希值
//func FindtxHash(fingerprint string) (file *File, err error) {
//	file = new(File)
//	if err = dao.DB.Where("fingerprint=?", fingerprint).First(&file).Error; err != nil {
//		return nil, err
//	}
//	return
//}

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

func In(FileOwner string, checkNode []string) (err error) {
	for i := range checkNode {
		if FileOwner == checkNode[i] {
			return errors.New("FileOwner made invalid charge")
		}
	}
	return nil
}

func Verify(applydatalist []Hashdata, sourceNode string, FingerprintingNode string) (bool, error) {
	//checkNode := sourceNode
	//NewcheckNode := ""
	//for init := 1; init <= len(applydatalist); init++ {
	//	for j := range applydatalist {
	//		applymessages := strings.Split(applydatalist[j].Result.Tx.Payload.ContentStorage.Value, "#")
	//		applyOwner := applymessages[0]
	//		if applyOwner == checkNode || applyOwner == NewcheckNode {
	//			if applyOwner == NewcheckNode {
	//				checkNode = NewcheckNode
	//			}
	//			switch applymessages[4] {
	//			case "4":
	//				applydatalist[j].Result.Tx.Payload.ContentStorage.Value += "#fail"
	//			case "5":
	//				applydatalist[j].Result.Tx.Payload.ContentStorage.Value += "#success"
	//			default:
	//				applydatalist[j].Result.Tx.Payload.ContentStorage.Value += "#fail"
	//				return applydatalist, errors.New("invalid param")
	//			}
	//			NewcheckNode = applymessages[1]
	//		}
	//	}
	//}
	//var checkNode []string
	//for j := range applydatalist {
	//	applymessages := strings.Split(applydatalist[j].Result.Tx.Payload.ContentStorage.Value, "#")
	//	if applymessages[4] == "4" {
	//		checkNode = append(checkNode, applymessages[0])
	//	}
	//}
	//for k := range applydatalist {
	//	applymessages := strings.Split(applydatalist[k].Result.Tx.Payload.ContentStorage.Value, "#")
	//	err := In(applymessages[1], checkNode)
	//	if err != nil {
	//		applydatalist[k].Result.Tx.Payload.ContentStorage.Value += "#fail"
	//	} else {
	//		applydatalist[k].Result.Tx.Payload.ContentStorage.Value += "#success"
	//	}
	//}
	IsLegal := false
	if sourceNode == FingerprintingNode {
		return true, nil
	}
	for j := range applydatalist {
		applymessages := strings.Split(applydatalist[j].Result.Tx.Payload.ContentStorage.Value, "#")
		if applymessages[0] == sourceNode && applymessages[1] == FingerprintingNode || applymessages[0] == FingerprintingNode && applymessages[1] == sourceNode {
			IsLegal = true
			break
		}
	}
	return IsLegal, nil
}

func TraceBackOnChain(txHash string, sourceNode string) ([]Hashdata, [][]string, string, string, bool, error) {
	//1.通过文件hash，查询文件信息,取出文件ID
	//定义结构体
	var applydatalist []Hashdata

	fileRecord, _ := AnalyzeData(txHash)
	fileRecord2 := strings.Split(fileRecord.Result.Tx.Payload.ContentStorage.Value, "#")
	//文件指纹所有者
	FingerprintingNode := fileRecord2[0]
	//2.取出文件ID,根据文件ID查询申请哈希，得到申请哈希，用一个数组存储循环查询申请记录
	applylist := new([]Applyrecord)
	if err := dao.DB.Where("file_id=?", fileRecord2[3]).Find(&applylist).Error; err != nil {
		return applydatalist, nil, sourceNode, "", false, err
	}
	//循环查询申请记录
	for i := range *applylist {
		apply := (*applylist)[i]
		applydata, _ := AnalyzeData(apply.Hash)
		applydatalist = append(applydatalist, applydata)
	}
	//对申请信息进行核验
	IsLegal, _ := Verify(applydatalist, sourceNode, FingerprintingNode)

	//简化追溯信息
	var checkdata [][]string
	for k := range applydatalist {
		var data []string
		applydatalist_message := applydatalist[k].Result.Tx.Payload.ContentStorage.Value
		applydatalist_messages := strings.Split(applydatalist_message, "#")
		data = append(data, applydatalist_messages[1])
		data = append(data, applydatalist_messages[0])
		data = append(data, applydatalist_messages[2])
		data = append(data, applydatalist_messages[4])
		checkdata = append(checkdata, data)
	}
	return applydatalist, checkdata, sourceNode, FingerprintingNode, IsLegal, nil
}

// 将效用分析文件保存到效用分析缓存区
func SaveFilelocal2(context *gin.Context) (Filetype string, err error) {
	// 将上传的文件取出来
	f, err := context.FormFile("FileAnalysis")
	if err != nil {
		context.JSON(http.StatusOK, gin.H{"message": err.Error()})
		return "", err
	}

	// 获取原始文件名和文件后缀
	originalFileName := f.Filename
	fileExtension := filepath.Ext(originalFileName)

	// 构建新的文件名
	newFileName := "analysisfile" + fileExtension

	// 设置文件名称为新的文件名
	f.Filename = newFileName

	// 将上传的文件保存到指定的本地地址，并返回响应
	log.Println(f.Filename)
	dst := fmt.Sprintf("./analysisfile/%s", f.Filename) // 设置效用分析文件保存的本地地址路径
	if err = context.SaveUploadedFile(f, dst); err != nil {
		return "", err
	}
	return fileExtension, nil
}

func Utility(context *gin.Context, fileid string) error {
	node := string(fileid[0]) // 取出文件所在节点信息

	// 取出文件名并构建对应服务器的路由
	file := new(File)
	if err_1 := dao.DB.Where("file_id=?", fileid).First(file).Error; err_1 != nil {
		println(err_1.Error())
		return err_1
	}
	str := fmt.Sprintf("%v", IP[node])
	address_1 := "http://" + str + ":8080/utility/do_TA/" + file.Name
	address_2 := "http://" + str + ":8080/utility/getfile"

	// 将扰动文件发送给对应的原始文件拥有者的服务器
	if err_2 := sendFileToServerX(address_2); err_2 != nil {
		return err_2
	}
	context.Redirect(http.StatusMovedPermanently, address_1)
	return nil
}

// 将扰动文件发送到原始文件拥有者的函数
func sendFileToServerX(serverBUrl string) error {

	fileLocation := "./analysisfile/analysisfile.csv"
	// 打开文件
	file, err := os.Open(fileLocation)
	if err != nil {
		return fmt.Errorf("打开文件失败: %w", err)
	}
	defer file.Close()

	// 创建 multipart writer
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("perturbed_file", file.Name())
	if err != nil {
		return fmt.Errorf("创建表单文件失败: %w", err)
	}

	// 将文件内容复制到 part
	_, err = io.Copy(part, file)
	if err != nil {
		return fmt.Errorf("复制文件内容失败: %w", err)
	}

	// 关闭 multipart writer
	err = writer.Close()
	if err != nil {
		return fmt.Errorf("关闭 multipart writer 失败: %w", err)
	}

	// 创建 HTTP 客户端和请求
	client := &http.Client{}
	req, err := http.NewRequest("POST", serverBUrl, body)
	if err != nil {
		return fmt.Errorf("创建请求失败: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("接收方返回错误状态码: %s", resp.Status)
	}

	fmt.Println("文件发送成功")
	return nil
}

func DOUtility(context *gin.Context, fileName string) ([][]float64, [][]float64, [][]float64, [][]float64, string, error) {
	// 构建文件路径
	var fN = strings.Split(fileName, ".")

	dst_original := fmt.Sprintf("./csvfile/%s.csv", fN[0]) // 原始文件的文件路径
	dst_perturbed := fmt.Sprintf("./analysisfile/analysisfile_FP.csv")

	// 打开原始文件进行检查
	Original_file, err := os.Open(dst_original)
	if err != nil {
		// 处理错误
		context.JSON(http.StatusNotFound, gin.H{"error": "Original File not found, may be delete"})
		return [][]float64{}, [][]float64{}, [][]float64{}, [][]float64{}, Node, err
	}
	Original_file.Close()

	// 打开扰动文件进行检查
	Perturbed_file, err := os.Open(dst_perturbed)
	if err != nil {
		// 处理错误
		context.JSON(http.StatusNotFound, gin.H{"error": "Perturbed File not found, may be delete"})
		return [][]float64{}, [][]float64{}, [][]float64{}, [][]float64{}, Node, err
	}
	Perturbed_file.Close()

	var Accuracy, W_Precision, W_Recall, W_F1 [][]float64

	// 使用 exec.Command 执行 Python 脚本
	//cmd := exec.Command("venv\\Scripts\\python.exe", "python/Model_Train.py", dst_original, dst_perturbed)
	cmd := exec.Command("./venv/bin/python", "python/Model_Train.py", dst_original, dst_perturbed)
	//cmd := exec.Command("/usr/bin/python", "python/Model_Train.py", dst_original, dst_perturbed)

	// 创建缓冲区来捕获输出和错误
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// 运行命令
	if err := cmd.Run(); err != nil {
		fmt.Println("Error executing Python script:", err)
		fmt.Println("Error Output:", stderr.String())
		return [][]float64{}, [][]float64{}, [][]float64{}, [][]float64{}, Node, err
	}

	// 获取 Python 输出
	result := strings.TrimSpace(stdout.String()) // 获取并去除多余的空白字符
	fmt.Println("Python Output:", result)

	// 将 JSON 字符串解析为三维数组
	var results [][][]float64
	err = json.Unmarshal([]byte(result), &results)
	if err != nil {
		fmt.Println("Error parsing output to [][][]float64:", err)
		return [][]float64{}, [][]float64{}, [][]float64{}, [][]float64{}, Node, err
	}

	// 将解析的结果分别赋值到对应的二维数组中
	if len(results) == 4 {
		Accuracy = results[0]
		W_Precision = results[1]
		W_Recall = results[2]
		W_F1 = results[3]
	}
	return Accuracy, W_Precision, W_Recall, W_F1, Node, nil
}

func DOUtility2(context *gin.Context) ([][]float64, [][]float64, [][]float64, [][]float64, string, error) {
	dst := fmt.Sprintf("./analysisfile/analysisfile.csv")

	file, err := os.Open(dst)
	if err != nil {
		// 处理错误
		context.JSON(http.StatusNotFound, gin.H{"error": "Original File not found, may be delete"})
		return [][]float64{}, [][]float64{}, [][]float64{}, [][]float64{}, "", err
	}
	file.Close()

	var Accuracy, W_Precision, W_Recall, W_F1 [][]float64

	// 使用 exec.Command 执行 Python 脚本
	//cmd := exec.Command("venv\\Scripts\\python.exe", "python/Model_Train2.py", dst)
	cmd := exec.Command("./venv/bin/python", "python/Model_Train2.py", dst)
	//cmd := exec.Command("/usr/bin/python", "python/Model_Train2.py", dst)

	// 创建缓冲区来捕获输出和错误
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// 运行命令
	if err := cmd.Run(); err != nil {
		fmt.Println("Error executing Python script:", err)
		fmt.Println("Error Output:", stderr.String())
		return [][]float64{}, [][]float64{}, [][]float64{}, [][]float64{}, "", err
	}

	// 获取 Python 输出
	result := strings.TrimSpace(stdout.String()) // 获取并去除多余的空白字符
	fmt.Println("Python Output:", result)

	// 将 JSON 字符串解析为三维数组
	var results [][][]float64
	err = json.Unmarshal([]byte(result), &results)
	if err != nil {
		fmt.Println("Error parsing output to [][][]float64:", err)
		return [][]float64{}, [][]float64{}, [][]float64{}, [][]float64{}, "", err
	}

	// 将解析的结果分别赋值到对应的二维数组中
	if len(results) == 4 {
		Accuracy = results[0]
		W_Precision = results[1]
		W_Recall = results[2]
		W_F1 = results[3]
	}
	return Accuracy, W_Precision, W_Recall, W_F1, "Self", nil
}

func DownloadModelFile(context *gin.Context, node, modelname string) (err error) {

	str := fmt.Sprintf("%v", IP[node])
	address := "http://" + str + ":8080/utility/download2/" + modelname + "/perturbed"
	context.Redirect(http.StatusMovedPermanently, address)
	return

}

func DownloadModel(context *gin.Context, ModelName, Type string) (err error) {
	var dst string
	if Type == "" {
		dst = fmt.Sprintf("./Train_Models/original_models/%s.pkl", ModelName) // 修改为正确的文件路径
	} else {
		dst = fmt.Sprintf("./Train_Models/perturbed_models/Perturbed_%s.pkl", ModelName)
	}

	// 打开文件
	file, err := os.Open(dst)
	if err != nil {
		// 处理错误
		context.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return err
	}
	defer file.Close()

	// 设置响应头
	fileName := ModelName
	if Type != "" {
		fileName = "Perturbed_" + ModelName
	}
	context.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s.pkl", fileName))
	context.Header("Content-Type", "application/octet-stream")
	context.File(dst)

	return nil
}

func AddFile(context *gin.Context) (err error) {
	var file File

	//file.FileOwner = context.PostForm("fileOwner") //todo：后面改成Node，这里测试不同节点用
	file.FileOwner = Node
	file.FileID = context.PostForm("file_id")
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
		dst := path.Join("./csvfile", f.Filename) //todo:这里修改文件路径
		_ = context.SaveUploadedFile(f, dst)
		context.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
		//file.Status = 1
		t := time.Now()
		file.Time = t.Format("2006-01-02 15:04:05")
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
