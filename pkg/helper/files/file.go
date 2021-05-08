package files

import (
	"fmt"
	"ginpro/config"
	util "ginpro/pkg/security"
	"io"
	"io/ioutil"
	"mime/multipart"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func LogFile() (io.Writer, string) {
	// 目录路径
	logDirPath := config.Conf.App.LogDir
	if ok := CheckSavePath(logDirPath); ok {
		if err := os.MkdirAll(logDirPath, os.ModePerm); err != nil {
			panic(err)
		}
	}
	// 文件名
	logFileName := config.Conf.App.AppName
	// 文件路径
	logFilePath := path.Join(logDirPath, logFileName)
	// 获取文件句柄
	f, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		fmt.Println("Open Src File err", err)
	}
	return f, logFilePath
}

type FileType int

const TypeImage FileType = iota + 1

func GetFileName(name string) string {
	ext := GetFileExt(name)
	fileName := strings.TrimSuffix(name, ext)
	fileName = util.EncodeMD5(fileName)

	return fileName + ext
}

func GetFileExt(name string) string {
	return path.Ext(name)
}

func GetSavePath() string {
	return config.Conf.App.UploadSavePath
}

func CheckSavePath(dst string) bool {
	_, err := os.Stat(dst)
	return os.IsNotExist(err)
}

func CheckContainExt(t FileType, name string) bool {
	ext := GetFileExt(name)
	ext = strings.ToUpper(ext)
	switch t {
	case TypeImage:
		for _, allowExt := range config.Conf.App.UploadImageAllowExts {
			if strings.ToUpper(allowExt) == ext {
				return true
			}
		}

	}

	return false
}

func CheckMaxSize(t FileType, f multipart.File) bool {
	content, _ := ioutil.ReadAll(f)
	size := len(content)
	switch t {
	case TypeImage:
		if size >= config.Conf.App.UploadImageMaxSize*1024*1024 {
			return true
		}
	}

	return false
}

func CheckPermission(dst string) bool {
	_, err := os.Stat(dst)
	return os.IsPermission(err)
}

func CreateSavePath(dst string, perm os.FileMode) error {
	err := os.MkdirAll(dst, perm)
	if err != nil {
		return err
	}

	return nil
}

func SaveFile(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}

func Basename(path string) string {
	return filepath.Base(path)
}
