package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	v1beta1 "k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

// FIXME: start with using logging library, like zap

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}

// FIXME: we don't need all s3 bucket, infact, start with a config file through the same you can get the configured s3 bucket content and the path
// func getS3Buckets(region string) {
// 	sess, err := session.NewSession(&aws.Config{
// 		Region: aws.String(region)},
// 	)

// 	if err != nil {
// 		exitErrorf("Unable to create a new session, %v", err)
// 	}

// FIXME: use a struct for the client, the clients are singleton use the calls as `func (s3client *S3) getProdjectCode(bucket, project)
// 	// Create S3 service client
// 	svc := s3.New(sess)

// 	result, err := svc.ListBuckets(nil)
// 	if err != nil {
// 		exitErrorf("Unable to list buckets, %v", err)
// 	}

// FIXME: cleanup bucket discovery
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
// FIXME, use session as struct outside method call
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

// TODO: the below code snippet looks good
// func launchK8sPod(clientset *kubernetes.Clientset, podName *string, image *string, cmd *string) {
// 	pods := clientset.CoreV1().Pods("kube-system")

// 	podSpec := &v1.Pod{
// 		ObjectMeta: metav1.ObjectMeta{
// 			Name:      *podName,
// 			Namespace: "kube-system",
// 		},
// 		Spec: v1.PodSpec{
// 			Containers: []v1.Container{
// 				{
// 					Name:  *podName,
// 					Image: *image,
// 					Ports: []v1.ContainerPort{
// 						{
// 							ContainerPort: 8080,
// 						},
// 					},
// 					VolumeMounts: []v1.VolumeMount{
// 						{
// 							Name:      "project",
// 							MountPath: "/home/abhishek/Downloads/my-project",
// 						},
// 					},
// 				},
// 			},
// 			Volumes: []v1.Volume{
// 				{
// 					Name: "project",
// 					VolumeSource: {
// 						v1.HostPathVolumeSource: {
// 							Type: *v1.HostPathDirectoryOrCreate,
// 							Path: "/home/abhishek/Downloads/rest-scripts",
// 						},
// 					},
// 				},
// 			},
// 			RestartPolicy: v1.RestartPolicyAlways,
// 		},
// 	}

// 	_, err := pods.Create(context.Background(), podSpec, metav1.CreateOptions{})
// 	if err != nil {
// 		exitErrorf("Failed to create K8s pod, %v", err)
// 		log.Fatalln("Failed to create K8s pod, ", err)
// 	}

// 	//print job details
// 	log.Println("Created K8s pod successfully")
// }

// FIXME:// need pod, if can't use above use atleast Replicaset
func launchK8sReplicaSet(clientset *kubernetes.Clientset, deploymentName *string, image *string, cmd *string) {

	client := clientset.AppsV1().ReplicaSet("infra")

	pod := &appsv1.ReplicaSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      *deploymentName,
			Namespace: "infra",
			Labels: map[string]string{
				"app": *deploymentName,
			},
		},
		Spec: appsv1.ReplicaSetSpec{
			Replicas: int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app":       *deploymentName,
					"component": *deploymentName,
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app":       *deploymentName,
						"component": *deploymentName,
					},
					Annotations: map[string]string{
						"sidecar.istio.io/inject": "false",
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:            *deploymentName,
							Image:           *image,
							ImagePullPolicy: apiv1.PullAlways,
							Ports: []apiv1.ContainerPort{
								{
									Name:          *deploymentName,
									ContainerPort: 8080,
								},
							},
							// TODO: Add volume mount and volumes
							Resources: apiv1.ResourceRequirements{
								Requests: apiv1.ResourceList{
									apiv1.ResourceCPU:    resource.MustParse("500m"),
									apiv1.ResourceMemory: resource.MustParse("512Mi"),
								},
								Limits: apiv1.ResourceList{
									apiv1.ResourceCPU:    resource.MustParse("500m"),
									apiv1.ResourceMemory: resource.MustParse("512Mi"),
								},
							},
						},
					},
				},
			},
		},
	}

	// Create Deployment
	fmt.Println("Creating deployment...")
	result, err := client.Create(context.TODO(), pod, metav1.CreateOptions{})
	if err != nil {
		exitErrorf("Failed to create K8s deployment, %v", err)
		log.Fatalln("Failed to create K8s deployment, ", err)
	}
	fmt.Printf("Created deployment %q.\n", result.GetObjectMeta().GetName())

	//print job details
	log.Println("Created K8s deployment successfully")
}

// TODO: att httpserver which can take a path like /{project}

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
	entryCommand := flag.String("command", "", "The command to run inside the container")

	flag.Parse()

	// TODO: http.serve should call launchPod function, by passing the projectName

	launchK8sDeployment(clientset, jobName, containerImage, entryCommand)

	//launchK8sIngress(clientset, "coder-x-ingress", entryCommand)
}

// FIXME: remove
func launchK8sIngress(clientset *kubernetes.Clientset, ingressName string, cmd *string) {

	ingressClient := clientset.ExtensionsV1beta1().Ingresses("infra")

	ingress := &v1beta1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name: ingressName,
			Annotations: map[string]string{
				"kubernetes.io/ingress.class":               "alb",
				"alb.ingress.kubernetes.io/security-groups": "sg-091a5b18c9896ce28",
				"alb.ingress.kubernetes.io/scheme":          "internet-facing",
				"alb.ingress.kubernetes.io/target-type":     "ip",
				"alb.ingress.kubernetes.io/group.name":      "default",
				"alb.ingress.kubernetes.io/listen-ports":    "[{\"HTTPS\":443}]",
			},
		},
		Spec: v1beta1.IngressSpec{
			Rules: []v1beta1.IngressRule{
				{
					Host: "coder-x.dev.nslhub.click",
				},
			},
		},
	}

	// Create Deployment
	fmt.Println("Creating Ingress...")
	result, err := ingressClient.Create(context.TODO(), ingress, metav1.CreateOptions{})
	if err != nil {
		exitErrorf("Failed to create K8s ingress, %v", err)
		log.Fatalln("Failed to create K8s ingress, ", err)
	}
	fmt.Printf("Created ingress %q.\n", result.GetObjectMeta().GetName())

	//print job details
	log.Println("Created K8s ingress successfully")
}

func int32Ptr(i int32) *int32 { return &i }
