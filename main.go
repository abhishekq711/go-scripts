package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/go-code-server/utils"
	"go.uber.org/zap"
)

const (
	S3_BUCKET  string = "codercom-code-server"
	AWS_REGION string = "ap-south-1"
)

func main() {

	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	zap.ReplaceGlobals(logger)

	//e.g. /ScriptApi.zip
	http.HandleFunc("/", handleRequest)
	http.HandleFunc("/favicon.ico", favicon)
	http.ListenAndServe(":8080", nil)

}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		fmt.Fprintf(w, "Invalid method, expected GET Method, received %v method\n", r.Method)
		zap.L().Error(fmt.Sprintf("Invalid method, expected GET Method, received %v method\n", r.Method))
	}

	item := r.RequestURI[1:]
	fmt.Fprintf(w, "URL GET parameter (this is also name of file): %s\n", item)
	zap.L().Info("URL GET parameter (this is also name of file): " + item)

	//Downloading the project from s3 bucket and then move to nextStep function
	if _, err := os.Stat(item); err == nil {
		//TODO-> this check of file is incomplete, need to make it right
		zap.L().Info("File appears to be already downloaded, skipping download step...")
		nextStep(item)
	} else {
		aws_sess, _ := session.NewSession(&aws.Config{
			Region: aws.String(AWS_REGION)},
		)
		err := utils.DownloadObject(S3_BUCKET, item, aws_sess)
		if err != nil {
			fmt.Fprintf(w, "Cannot download the specified file (filename: %s), make sure filename/path is correct and check logs for more info\n", item)
			zap.L().Error(fmt.Sprintf("Download (of filename: %s) failed with error: %v", item, err))
		} else {
			nextStep(item)
		}
	}
}

// Unziping and pod launching
func nextStep(item string) {
	err := utils.Unzip(item, "myprojects")
	if err != nil {
		zap.L().Fatal("Failed to unzip file")
		return
	}
	zap.L().Info("File unziping successful")

	// build and create the clientset
	clientset, err := utils.GetKubeClient()
	if err != nil {
		zap.L().Error("Unable to create kube clientset. Exiting with err " + err.Error())
		return
	}
	zap.L().Info("clientset created successfully")

	// utils.ListAllK8Pods(clientset)

	jobName := flag.String("podname", "coder-x", "The name of the pod")
	containerImage := flag.String("image", "codercom/code-server", "Name of the container image")
	entryCommand := flag.String("command", "ls", "The command to run inside the container")

	flag.Parse()

	utils.LaunchK8sPod(clientset, jobName, containerImage, entryCommand)

}

//handling favicon uri
func favicon(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%s\n", r.RequestURI)
	w.Header().Set("Content-Type", "image/x-icon")
	w.Header().Set("Cache-Control", "public, max-age=7776000")
	fmt.Fprintln(w, "data:image/x-icon;base64,iVBORw0KGgoAAAANSUhEUgAAABAAAAAQEAYAAABPYyMiAAAABmJLR0T///////8JWPfcAAAACXBIWXMAAABIAAAASABGyWs+AAAAF0lEQVRIx2NgGAWjYBSMglEwCkbBSAcACBAAAeaR9cIAAAAASUVORK5CYII=")
}
