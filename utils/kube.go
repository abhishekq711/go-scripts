package utils

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func ListAllK8Pods(clientset *kubernetes.Clientset) {

	nodeList, err := clientset.CoreV1().Pods("").List(context.Background(), metav1.ListOptions{})

	if err != nil {
		zap.L().Error("Unable to list nodes, " + err.Error())
	}

	zap.L().Info(fmt.Sprintf("Total number of nodes: %d", len(nodeList.Items)))

	for m, n := range nodeList.Items {
		zap.L().Info(fmt.Sprintf("%d %s", m, n.Name))
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
						// NFS: &v1.NFSVolumeSource{
						// 	Server: "nfs-server-ip-address",
						// 	Path:   "/dir/path/on/nfs/server",
						// },
						// PersistentVolumeClaim: &v1.PersistentVolumeClaimVolumeSource{
						// 	ClaimName: "efs-claim",
						// },
					},
				},
			},
			RestartPolicy: v1.RestartPolicyAlways,
		},
	}

	_, err := pods.Create(context.Background(), podSpec, metav1.CreateOptions{})
	if err != nil {
		zap.L().Fatal(fmt.Sprintf("Failed to create K8s pod, %v", err))
	}

	zap.L().Info("Created K8s pod successfully")
}
