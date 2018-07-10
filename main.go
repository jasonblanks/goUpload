package main

import (
	"net/http"
	"fmt"
	"time"
	"cloud.google.com/go/storage"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"bytes"
	"mime/multipart"
	"os"
	"io/ioutil"
	"math/rand"
	"crypto/sha1"
	"encoding/hex"
	"strconv"
	"gopkg.in/yaml.v2"
)

const projectID = "INSERT-GCS-PROJECT-NAME"

type upload struct {
	Sha1 string
	CreationTime time.Time
	ExpireDate time.Time
	AccessedTime time.Time
	RequestedBy string
	RequesterInput
	UploaderInput
}
type RequesterInput struct{
	Details string
	Firstname string
	Lastname string
	CaseId string
}
type UploaderInput struct{
	Details string
	Firstname string
	Email string
	Lastname string
}

var RegesteredUploads map[string]*upload

func registerGet(c *gin.Context){
	c.HTML(http.StatusOK, "register.gtpl", gin.H{},)
}

func registerPost(c *gin.Context) {
	fmt.Println(c.PostForm("firstname"))
	t := time.Now()
	rand.Seed(time.Now().UTC().UnixNano())
	text := t.String() + string(rand.Intn(999999999999-0))
	hasher := sha1.New()
	hasher.Write([]byte(text))
	newHash := hex.EncodeToString(hasher.Sum(nil))
	timeValue, err := strconv.Atoi(c.PostForm("timeValue"))
	if err != nil {
		fmt.Println(err)
		return
	}
	var timein time.Time
	if c.PostForm("timeMeasure") == "minute" {
		timein = time.Now().Local().Add(time.Minute * time.Duration(timeValue))
	} else if c.PostForm("timeMeasure") == "hour" {
		timein = time.Now().Local().Add(time.Hour * time.Duration(timeValue))
	} else if c.PostForm("timeMeasure") == "day" {
		timein = time.Now().Local().Add(time.Hour * 24 * time.Duration(timeValue))
	}
	RegesteredUploads[newHash] = &upload{
		Sha1: newHash,
		CreationTime: time.Now(),
		ExpireDate: timein,
		RequestedBy: c.PostForm("firstname"),
		RequesterInput: RequesterInput{
			Details:   c.PostForm("reason"),
			Firstname: c.PostForm("firstname"),
			CaseId:    c.PostForm("caseNumber"),
		},
	}
	fmt.Println(timein)
	fmt.Println(RegesteredUploads)
	fmt.Println(RegesteredUploads[newHash])
	c.String(http.StatusOK, fmt.Sprintf("New upload link registered: http://127.0.0.1:8080/upload/%s", newHash))
	//URLs = append(URLs,newUrl)
}
func createSignedURL(filename string, ExpireTime time.Time) (string){
	bucket := "INSERT-BUCKET-NAME-HERE"
	method := "PUT"
	//expires := time.Now().Add(time.Minute * 60)
	expires := ExpireTime

	url, err := storage.SignedURL(bucket, filename, &storage.SignedURLOptions{
		GoogleAccessID: "INSERT-SERVICE-ACCOUNT-HERE",
		PrivateKey:     []byte("INSERT-PRIVATE-KEY-STRING-HERE"),
		Method:         method,
		Expires:        expires,
	})
	if err != nil {
		fmt.Println("Error " + err.Error())
	}
	//fmt.Println("URL = " + url)
	return url
}
func uploadGet(c *gin.Context) {
	url := c.Param("url")
	if _, ok := RegesteredUploads[url]; ok {
		if time.Now().Before(RegesteredUploads[url].ExpireDate){
			c.HTML(http.StatusOK, "uploadForm.gtpl", gin.H{"token": url},)
		}
	}
}
func uploadPost(c *gin.Context) {
	url := c.Param("url")
	RegesteredUploads[url].UploaderInput.Details = c.PostForm("userReason")
	RegesteredUploads[url].UploaderInput.Email = c.PostForm("userEmail")
	RegesteredUploads[url].UploaderInput.Lastname = c.PostForm("userLastname")
	RegesteredUploads[url].UploaderInput.Firstname = c.PostForm("userFirstname")

	caseNumber := RegesteredUploads[url].RequesterInput.CaseId
	file, header , err := c.Request.FormFile("uploadfile")
	filename := header.Filename
	pathName := caseNumber + "/" + filename
	infoFile := caseNumber + "/info.yml"

	if err != nil {
		fmt.Println(err)
		return
	}
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPut, createSignedURL(pathName, RegesteredUploads[url].ExpireDate), file)
	if err != nil {
		// handle error
		log.Fatal(err)
		fmt.Println(err)
	}
	_, err = client.Do(req)
	if err != nil {
			// handle error
			log.Fatal(err)
			fmt.Println(err)
	}
	d, err := yaml.Marshal(&RegesteredUploads)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	dd := bytes.NewBuffer(d)
	req, err = http.NewRequest(http.MethodPut, createSignedURL(infoFile, RegesteredUploads[url].ExpireDate), dd)
	if err != nil {
		// handle error
		log.Fatal(err)
		fmt.Println(err)
	}
	_, err = client.Do(req)
	if err != nil {
		// handle error
		log.Fatal(err)
		fmt.Println(err)
	}
	c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", filename))
	//fmt.Printf("--- m dump:\n%s\n\n", string(d))

	}
func postFile(filename string, targetUrl string) error {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	// this step is very important
	fileWriter, err := bodyWriter.CreateFormFile("uploadfile", filename)
	if err != nil {
		fmt.Println("error writing to buffer")
		return err
	}

	// open file handle
	fh, err := os.Open(filename)
	if err != nil {
		fmt.Println("error opening file")
		return err
	}
	defer fh.Close()

	//iocopy
	_, err = io.Copy(fileWriter, fh)
	if err != nil {
		return err
	}

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	resp, err := http.Post(targetUrl, contentType, bodyBuf)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	resp_body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println(resp.Status)
	fmt.Println(string(resp_body))
	return nil
}
func main() {
	RegesteredUploads = make(map[string]*upload)
	router := gin.Default()
	router.LoadHTMLFiles("register.gtpl", "uploadForm.gtpl")

	router.GET("/register", registerGet)
	router.POST("/register", registerPost)
	//router.GET("/registered", registered)
	router.GET("/upload/:url", uploadGet)
	router.POST("/upload/:url", uploadPost)

	router.Run(":8080")
}