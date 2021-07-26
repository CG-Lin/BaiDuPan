package handler

import (
	"BaiDuPan/meta"
	"BaiDuPan/util"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)
//上传文件
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		data, err := ioutil.ReadFile("./static/view/index.html")
		if err != nil {
			io.WriteString(w, "internel server error")
			return
		}
		io.WriteString(w, string(data))
	} else if r.Method == "POST" {
		file, head, err := r.FormFile("file")
		if err != nil {
			fmt.Printf("Failed to start server,err:%s\n", err.Error())
			return
		}
		defer file.Close()

		fileMeta := meta.FileMeta{
			FileName: head.Filename,
			Location: "/tmp/" + head.Filename,
			UploadAt: time.Now().Format("2006-01-02 15:04:05"),
		}

		newFile, err := os.Create(fileMeta.Location)
		if err != nil {
			fmt.Printf("Failed to start server,err:%s\n", err.Error())
			return
		}
		defer newFile.Close()

		fileMeta.FileSize, err = io.Copy(newFile, file)
		if err != nil {
			fmt.Printf("Failed to start server,err:%s\n", err.Error())
			return
		}
		newFile.Seek(0,0)
		fileMeta.FileShal = util.FileSha1(newFile)
		meta.UpdateFileMeta(fileMeta)

		http.Redirect(w,r,"/file/upload/suc",http.StatusFound)
	}
}

func UploadSucHandler(w http.ResponseWriter,r *http.Request)  {
	io.WriteString(w,"Upload finished!")
}

func FileQueryHandler(w http.ResponseWriter,r *http.Request)  {
	//解析url传递的参数，对于POST则解析响应包的主体
	r.ParseForm()
	//注意:如果没有调用ParseForm方法，下面无法获取表单的数据
	filehash:=r.Form["filehash"][0]
	fMeta :=meta.GetFileMeta(filehash)
	data,err:=json.Marshal(fMeta)
	if err!=nil{
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}