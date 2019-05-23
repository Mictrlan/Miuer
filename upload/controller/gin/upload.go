package gin

import (
	"database/sql"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/Mictrlan/Miuer/upload/utility"

	"github.com/Mictrlan/Miuer/upload/model/mysql"
	"github.com/gin-gonic/gin"
)

var (
	errServerNotExists = errors.New("[RegisterRouter]: server is nil")
	errRequest         = errors.New("Request is not post method")
	erruserID          = errors.New("userID invalid")
)

// UploadController -
type UploadController struct {
	db  *sql.DB
	URL string
	UID func(c *gin.Context) (uint32, error)
}

// New create new uploadcontroller
func New(db *sql.DB, URL string, UID func(c *gin.Context) (uint32, error)) *UploadController {
	return &UploadController{
		db:  db,
		URL: URL,
		UID: UID,
	}
}

// Register register router
func (uc *UploadController) Register(r gin.IRouter) {
	if r == nil {
		log.Fatal(errServerNotExists)
	}

	err := mysql.CreateTable(uc.db)
	if err != nil {
		log.Fatal(err)
	}

	err = checkDir(utility.PictureDir, utility.VideoDir, utility.OtherDir)
	if err != nil {
		log.Fatal(err)
	}

	r.POST("/api/v1/user/upload", uc.Upload)

}

// checkDir Verify directory existence， if directory dosen't exists then create it
// Stat returns a FileInfo describing the named file.
// IsNotExist if file or directory not exists return true
// MkdirAll creates a directory named path
func checkDir(path ...string) error {
	for _, name := range path {
		_, err := os.Stat(utility.FileUploadDir + "/" + name)
		if err != nil {
			if os.IsNotExist(err) {
				err = os.MkdirAll(utility.FileUploadDir+"/"+name, 0777)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// Upload -
func (uc *UploadController) Upload(ctx *gin.Context) {

	if ctx.Request.Method != "POST" {
		ctx.Error(errRequest)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	userID, err := uc.UID(ctx)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	if userID == utility.InvalidUID {
		ctx.Error(erruserID)
		ctx.JSON(http.StatusExpectationFailed, gin.H{"status": http.StatusExpectationFailed})
		return
	}

	// FormFile returns the first file for the provided form key.
	// RemoveAll removes any temporary files associated with a Form.
	// get file and header
	file, header, err := ctx.Request.FormFile(utility.FileKey)
	defer func() {
		file.Close()
		ctx.Request.MultipartForm.RemoveAll()
	}()
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusForbidden, gin.H{"status": http.StatusForbidden})
		return
	}

	// file 只可以别读取一次
	fileNew, _ := ioutil.ReadAll(file)

	MD5Str, err := utility.MD5(fileNew)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusForbidden, gin.H{"status": http.StatusForbidden})
		return
	}

	filePath, err := mysql.QueryPathByMD5(uc.db, MD5Str)
	if err == nil {
		ctx.JSON(http.StatusOK, gin.H{
			"status":   http.StatusOK,
			"filePath": uc.URL + filePath,
		})

		return
	}

	if err != mysql.ErrNoRows {
		ctx.Error(err)
		ctx.JSON(http.StatusNotAcceptable, gin.H{"status": http.StatusNotAcceptable})
		return
	}

	// Ext returns the file name extension used by path.
	fileSuffix := path.Ext(header.Filename)
	filePath = utility.FileUploadDir + "/" + utility.ClassifyBySuffix(fileSuffix) + "/" + MD5Str + fileSuffix

	err = utility.CopyFile(filePath, fileNew)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusNotModified, gin.H{"status": http.StatusNotModified})
		return
	}

	err = mysql.Insert(uc.db, userID, filePath, MD5Str)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusPreconditionFailed, gin.H{"status": http.StatusPreconditionFailed})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"status":   http.StatusCreated,
		"filePath": uc.URL + filePath,
	})
}
