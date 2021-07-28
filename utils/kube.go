package utils

import (
	"context"
	"fmt"
	"log"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func ListAllK8Pods(clientset *kubernetes.Clientset) {

	nodeList, err := clientset.CoreV1().Pods("").List(context.Background(), metav1.ListOptions{})

	if err != nil {
		ExitErrorf("Unable to list nodes, %v", err)
	}

	fmt.Println(len(nodeList.Items))

	for m, n := range nodeList.Items {
		fmt.Println(m, n.Name)
	}
}

func LaunchK8sPod(clientset *kubernetes.Clientset, podName *string, image *string, cmd *string) {
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
							MountPath: "/home/abhishek/Downloads/rest-scripts/myprojects",
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
							Type: (*v1.HostPathType)(Strptr("DirectoryOrCreate")),
						},
					},
				},
			},
			RestartPolicy: v1.RestartPolicyAlways,
		},
	}

	_, err := pods.Create(context.Background(), podSpec, metav1.CreateOptions{})
	if err != nil {
		log.Fatalln("Failed to create K8s pod, ", err)
	}

	log.Println("Created K8s pod successfully")
}
