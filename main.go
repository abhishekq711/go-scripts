package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}

// func getS3Buckets(region string) {
// 	sess, err := session.NewSession(&aws.Config{
// 		Region: aws.String(region)},
// 	)

// 	if err != nil {
// 		exitErrorf("Unable to create a new session, %v", err)
// 	}

// 	// Create S3 service client
// 	svc := s3.New(sess)

// 	result, err := svc.ListBuckets(nil)
// 	if err != nil {
// 		exitErrorf("Unable to list buckets, %v", err)
// 	}

// 	fmt.Println("Buckets:")

// 	for _, b := range result.Buckets {
// 		fmt.Printf("* %s created on %s\n",
// 			aws.StringValue(b.Name), aws.TimeValue(b.CreationDate))
// 	}
// }

// func DownloadObject(bucket, item, region string) {

// 	file, err := os.Create(item)
// 	if err != nil {
// 		exitErrorf("Unable to open file %q, %v", item, err)
// 	}

// 	defer file.Close()

// 	sess, _ := session.NewSession(&aws.Config{
// 		Region: aws.String(region)},
// 	)

// 	downloader := s3manager.NewDownloader(sess)

// 	numBytes, err := downloader.Download(file,
// 		&s3.GetObjectInput{
// 			Bucket: aws.String(bucket),
// 			Key:    aws.String(item),
// 		})
// 	if err != nil {
// 		exitErrorf("Unable to download item %q, %v", item, err)
// 	}

// 	fmt.Println("Downloaded", file.Name(), numBytes, "bytes")
// }

// func Unzip(src, dest string) error {
// 	r, err := zip.OpenReader(src)
// 	if err != nil {
// 		return err
// 	}
// 	defer func() {
// 		if err := r.Close(); err != nil {
// 			panic(err)
// 		}
// 	}()

// 	os.MkdirAll(dest, 0755)

// 	// Closure to address file descriptors issue with all the deferred .Close() methods
// 	extractAndWriteFile := func(f *zip.File) error {
// 		rc, err := f.Open()
// 		if err != nil {
// 			return err
// 		}
// 		defer func() {
// 			if err := rc.Close(); err != nil {
// 				panic(err)
// 			}
// 		}()

// 		path := filepath.Join(dest, f.Name)

// 		// Check for ZipSlip (Directory traversal)
// 		if !strings.HasPrefix(path, filepath.Clean(dest)+string(os.PathSeparator)) {
// 			return fmt.Errorf("illegal file path: %s", path)
// 		}

// 		if f.FileInfo().IsDir() {
// 			os.MkdirAll(path, f.Mode())
// 		} else {
// 			os.MkdirAll(filepath.Dir(path), f.Mode())
// 			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
// 			if err != nil {
// 				return err
// 			}
// 			defer func() {
// 				if err := f.Close(); err != nil {
// 					panic(err)
// 				}
// 			}()

// 			_, err = io.Copy(f, rc)
// 			if err != nil {
// 				return err
// 			}
// 		}
// 		return nil
// 	}

// 	for _, f := range r.File {
// 		err := extractAndWriteFile(f)
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }

func launchK8sPod(clientset *kubernetes.Clientset, podName *string, image *string, cmd *string) {
	pods := clientset.CoreV1().Pods("infra")

	podSpec := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      *podName,
			Namespace: "infra",
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:  *podName,
					Image: *image,
					Ports: []v1.ContainerPort{
						{
							ContainerPort: 8080,
						},
					},
					Env: []v1.EnvVar{
						{
							Name:  "PASSWORD",
							Value: "mypass",
						},
					},
					VolumeMounts: []v1.VolumeMount{
						{
							Name:      "project",
							MountPath: "/home/abhishek/Downloads/my-project",
						},
					},
				},
			},
			Volumes: []v1.Volume{
				{
					Name: "project",
					VolumeSource: v1.VolumeSource{
						HostPath: &v1.HostPathVolumeSource{
							Path: "/home/abhishek/Downloads/rest-scripts",
						},
					},
				},
			},
			RestartPolicy: v1.RestartPolicyAlways,
		},
	}

	_, err := pods.Create(context.Background(), podSpec, metav1.CreateOptions{})
	if err != nil {
		exitErrorf("Failed to create K8s pod, %v", err)
		log.Fatalln("Failed to create K8s pod, ", err)
	}

	log.Println("Created K8s pod successfully")
}

func main() {

	// var region string = "ap-south-1"

	// getS3Buckets(region)

	// if len(os.Args) != 3 {
	// 	exitErrorf("Bucket and item names required\nUsage: %s bucket_name item_name",
	// 		os.Args[0])
	// }

	// bucket := os.Args[1]
	// item := os.Args[2]

	// DownloadObject(bucket, item, region)

	// fmt.Println(Unzip("./my-project.zip", "./myproject"))

	// Kube-client

	// rules := clientcmd.NewDefaultClientConfigLoadingRules()
	// kubeconfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(rules, &clientcmd.ConfigOverrides{})
	// config, err := kubeconfig.ClientConfig()

	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube/config"), "/home/abhishek/.kube/config")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "/home/abhishek/.kube/config")
	}
	flag.Parse()

	fmt.Println("hello1", *kubeconfig)

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	nodeList, err := clientset.CoreV1().Pods("").List(context.Background(), metav1.ListOptions{})

	if err != nil {
		// exitErrorf("Unable to list nodes, %v", err)
		panic(err)
	}

	fmt.Println(len(nodeList.Items))

	for m, n := range nodeList.Items {
		fmt.Println(m, n.Name)
	}

	jobName := flag.String("podname", "coder-x", "The name of the pod")
	containerImage := flag.String("image", "codercom/code-server", "Name of the container image")
	entryCommand := flag.String("command", "ls", "The command to run inside the container")

	flag.Parse()

	launchK8sPod(clientset, jobName, containerImage, entryCommand)

}
