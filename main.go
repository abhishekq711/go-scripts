package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/rest-scripts/utils"
	"go.uber.org/zap"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

const (
	BUCKET_NAME string = "codercom-code-server"
	AWS_REGION  string = "ap-south-1"
)

func main() {

	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	zap.ReplaceGlobals(logger)

	http.HandleFunc("/", handleRequest)
	http.HandleFunc("/favicon.ico", favicon)
	http.ListenAndServe(":8080", nil)

}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	//ScriptApi.zip
	if r.Method != "GET" {
		fmt.Fprintf(w, "Invalid method, expected GET Method, received %v method\n", r.Method)
		zap.L().Error(fmt.Sprintf("Invalid method, expected GET Method, received %v method\n", r.Method))
	}

	item := r.RequestURI[1:]
	fmt.Fprintf(w, "URL GET parameter: %v\n", item)
	zap.L().Info("URL GET parameter: " + item)

	if item == "/" {
		fmt.Fprintf(w, "item/project name in the bucket required, no name specified\n")
		zap.L().Error("item/project name in the bucket required, no name specified\n")
	}

	//TODO-> this check of file is incomplete, need to make it right
	if _, err := os.Stat(item); err == nil {
		zap.L().Info("File is already downloaded, skipping download step...")
		nextStep(item)
	} else {
		err := utils.DownloadObject(BUCKET_NAME, item, AWS_REGION)
		if err != nil {
			zap.L().Error(fmt.Sprintf("Download failed with error: %v", err))
		} else {
			nextStep(item)
		}
	}
}

func nextStep(item string) {
	err := utils.Unzip(item, "myprojects")
	if err != nil {
		zap.L().Fatal("Failed to unzip file")
		return
	}
	zap.L().Info("File unziping successful")

	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube/config"), "/home/abhishek/.kube/config")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "/home/abhishek/.kube/config")
	}
	flag.Parse()

	zap.L().Info("Kubeconfig file path used: " + *kubeconfig)

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		zap.L().Error("Unable to build kube config. Exiting with err " + err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		zap.L().Error("Unable to create kube clientset. Exiting with err " + err.Error())
	}

	utils.ListAllK8Pods(clientset)

	jobName := flag.String("podname", "coder-x", "The name of the pod")
	containerImage := flag.String("image", "codercom/code-server", "Name of the container image")
	entryCommand := flag.String("command", "ls", "The command to run inside the container")

	flag.Parse()

	utils.LaunchK8sPod(clientset, jobName, containerImage, entryCommand)

}

func favicon(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%s\n", r.RequestURI)
	w.Header().Set("Content-Type", "image/x-icon")
	w.Header().Set("Cache-Control", "public, max-age=7776000")
	fmt.Fprintln(w, "data:image/x-icon;base64,iVBORw0KGgoAAAANSUhEUgAAABAAAAAQEAYAAABPYyMiAAAABmJLR0T///////8JWPfcAAAACXBIWXMAAABIAAAASABGyWs+AAAAF0lEQVRIx2NgGAWjYBSMglEwCkbBSAcACBAAAeaR9cIAAAAASUVORK5CYII=")
}
