package kubernetes

import (
	"fmt"
	"os"

	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// CreateNamespaceForUser creates a new namespace and makes the specified user
// an admin of that namespace
func CreateNamespaceForUser(name, user string) (*corev1.Namespace, error) {
	clientset, err := getClientSet()
	if err != nil {
		return nil, err
	}

	ns, err := clientset.CoreV1().Namespaces().Create(&corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	})
	if err != nil {
		return nil, err
	}

	_, err = clientset.RbacV1().RoleBindings(name).Create(&rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: name,
			Name:      fmt.Sprintf("nsmaker:%s", user),
		},
		RoleRef: rbacv1.RoleRef{
			Kind: "ClusterRole",
			Name: "admin",
		},
		Subjects: []rbacv1.Subject{
			rbacv1.Subject{
				Kind: rbacv1.UserKind,
				Name: user,
			},
		},
	})
	if err != nil {
		return nil, err
	}

	return ns, nil
}

func getClientSet() (*kubernetes.Clientset, error) {
	var config *rest.Config
	var err error

	if os.Getenv("KUBERNETES_SERVICE_HOST") != "" && os.Getenv("KUBERNETES_SERVICE_PORT") != "" {
		config, err = rest.InClusterConfig()
		if err != nil {
			return nil, err
		}
	} else if os.Getenv("KUBECONFIG") != "" {
		config, err = clientcmd.BuildConfigFromFlags("", os.Getenv("KUBECONFIG"))
		if err != nil {
			return nil, err
		}
	} else {
		config, err = clientcmd.BuildConfigFromFlags("", os.Getenv("HOME")+"/.kube/config")
		if err != nil {
			return nil, err
		}
	}

	return kubernetes.NewForConfig(config)
}
