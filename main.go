package main

import (
	"flag"
	"path/filepath"

	"github.com/rest-scripts/utils"
	"go.uber.org/zap"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {

	// var region string = "ap-south-1"

	// if len(os.Args) != 3 {
	// 	exitErrorf("Bucket and item names required\nUsage: %s bucket_name item_name",
	// 		os.Args[0])
	// }

	// bucket := os.Args[1]
	// item := os.Args[2]

	// DownloadObject(bucket, item, region)

	// fmt.Println(Unzip("./my-project.zip", "./myproject"))

	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	sugar := logger.Sugar()

	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube/config"), "/home/abhishek/.kube/config")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "/home/abhishek/.kube/config")
	}
	flag.Parse()

	sugar.Infof("Kubeconfig file path used: %s", *kubeconfig)

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		sugar.Errorf("Unable to create kube clientset. Exiting with err %v", err.Error())
	}

	utils.ListAllK8Pods(clientset)

	// jobName := flag.String("podname", "coder-x", "The name of the pod")
	// containerImage := flag.String("image", "codercom/code-server", "Name of the container image")
	// entryCommand := flag.String("command", "ls", "The command to run inside the container")

	// flag.Parse()

	// utils.LaunchK8sPod(clientset, jobName, containerImage, entryCommand)

}
