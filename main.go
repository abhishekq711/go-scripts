package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/rest-scripts/utils"
	"go.uber.org/zap"
)

const (
	BUCKET_NAME string = "codercom-code-server"
	AWS_REGION  string = "ap-south-1"
)

func main() {

	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	zap.ReplaceGlobals(logger)

	// var kubeconfig *string
	// if home := homedir.HomeDir(); home != "" {
	// 	kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube/config"), "/home/abhishek/.kube/config")
	// } else {
	// 	kubeconfig = flag.String("kubeconfig", "", "/home/abhishek/.kube/config")
	// }
	// flag.Parse()

	// sugar.Infof("Kubeconfig file path used: %s", *kubeconfig)

	// // use the current context in kubeconfig
	// config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	// if err != nil {
	// 	sugar.Errorf("Unable to build kube config. Exiting with err %v", err.Error())
	// }

	// // create the clientset
	// clientset, err := kubernetes.NewForConfig(config)
	// if err != nil {
	// 	sugar.Errorf("Unable to create kube clientset. Exiting with err %v", err.Error())
	// }

	// utils.ListAllK8Pods(clientset)

	http.HandleFunc("/", greet)
	http.ListenAndServe(":8080", nil)

	// jobName := flag.String("podname", "coder-x", "The name of the pod")
	// containerImage := flag.String("image", "codercom/code-server", "Name of the container image")
	// entryCommand := flag.String("command", "ls", "The command to run inside the container")

	// flag.Parse()

	// utils.LaunchK8sPod(clientset, jobName, containerImage, entryCommand)

}

func greet(w http.ResponseWriter, r *http.Request) {
	//ScriptApi.zip
	if r.Method != "GET" {
		fmt.Fprintf(w, "Invalid method, expected GET Method, received %v method\n", r.Method)
		zap.L().Error("Invalid method, expected GET Method\n")
	}

	item := r.RequestURI[1:]
	fmt.Fprintf(w, "URL GET parameter: %v\n", item)
	zap.L().Info("URL GET parameter: " + item)

	if item == "/" {
		fmt.Fprintf(w, "Bucket and item names required, no bucket name specified\n")
		zap.L().Error("Bucket and item names required, no bucket name specified\n")
	}

	if _, err := os.Stat(item); err == nil {
		zap.L().Info("File is already downloaded, skipping download step...")

	} else {
		utils.DownloadObject(BUCKET_NAME, item, AWS_REGION)
		zap.L().Info("File download from aws s3 bucket successful")
	}

	err := utils.Unzip(item, "myproject")
	if err != nil {
		zap.L().Error("Failed unzip file")
	} else {
		zap.L().Info("File unziping successful")
	}
}
