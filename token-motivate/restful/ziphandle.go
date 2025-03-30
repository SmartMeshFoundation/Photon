package restful

import (
	"archive/zip"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/karlseguin/ccache"
	"go.cryptoscope.co/ssb/restful/params"
)

var cache = ccache.New(ccache.Configure().MaxSize(50).ItemsToPrune(5).OnDelete(func(item *ccache.Item) {
	item.Value().(*os.File).Close()
}))

func NewZipHandle(f io.ReaderAt, size int64) (*Ziphandle, error) {
	res := &Ziphandle{
		f:        f,
		fileSize: size,
		File:     make(map[string][]*ZipFileStruct),
	}
	err := res.ZipParse()
	if err != nil {
		return nil, err
	}
	return res, nil
}

// 解析接受的文件
type Ziphandle struct {
	f         io.ReaderAt
	fileSize  int64
	File      map[string][]*ZipFileStruct
	callbacks []callback //注入解析函数
}

type callback func(handle *Ziphandle) (err error)

func CreateFile(z *Ziphandle) (err error) {
	if z == nil {
		return errors.New("Uninitialized Ziphandle")
	}
	if z.File == nil {
		return errors.New("no file")
	}
	for _, tf := range z.File {
		for _, fi := range tf {
			if !IsFileExist(fi.Path) {
				err = os.MkdirAll(fi.Path, os.ModePerm)
			}
			f := fi.File
			srcFile, err := f.Open()
			defer srcFile.Close()
			if err != nil {
				return err

			}

			newFile, err := os.Create(filepath.Join(params.GameUserFilePath, f.Name))
			defer newFile.Close()
			if err != nil {
				return err
			}
			_, err = io.Copy(newFile, srcFile)
			if err != nil {
				return err
			}

		}
	}
	return
}

type ZipFileStruct struct {
	FileName string
	Path     string
	FileType string
	File     *zip.File
}

func (z *Ziphandle) ZipParse() error {
	zipFile, err := zip.NewReader(z.f, z.fileSize)
	if err != nil {
		return err
	}
	for _, zf := range zipFile.File {
		info := zf.FileInfo()
		if info.IsDir() {
			continue
		} else {
			resFile := &ZipFileStruct{
				FileName: info.Name(),
				Path:     filepath.Join(params.GameUserFilePath, strings.TrimRight(zf.FileHeader.Name, info.Name())),
				//Path:strings.TrimRight(zf.FileHeader.Name,info.Name()),
				File:     zf,
				FileType: "",
			}
			//fmt.Println(resFile.Path)
			ext := strings.TrimLeft(filepath.Ext(info.Name()), ".")
			if ext == "" {
				continue
			}
			resFile.FileType = ext
			z.File[ext] = append(z.File[ext], resFile)
		}
	}
	return nil
}

func (z *Ziphandle) OnHandle(cb ...callback) {
	if z.callbacks == nil {
		z.callbacks = make([]callback, 0)
	}
	if len(cb) > 0 {
		z.callbacks = append(z.callbacks, cb...)
	}
}

func (z *Ziphandle) Handle() (err error) {
	if z == nil || z.f == nil || z.File == nil {
		return errors.New("parameter error")
	}
	for _, cb := range z.callbacks {
		if err = cb(z); err != nil {
			return err
		}
	}
	return
}

func IsFileExist(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}
