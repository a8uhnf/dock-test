package handlers

import (
	"fmt"
	"log"

	"github.com/a8uhnf/container-h/pkg/event"
	api_v1 "k8s.io/api/extensions/v1beta1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type Logger struct {
	Client kubernetes.Interface
}

func (l *Logger) ObjectCreated(obj interface{}) {
	switch obj.(type) {
	case *api_v1.Deployment:
		fmt.Println("New Deployment")
	}
	logEvent(l, obj, "created")
}

func (l *Logger) ObjectDeleted(obj interface{}) {
	logEvent(l, obj, "deleted")
}

func (l *Logger) ObjectUpdated(oldObj, newObj interface{}) {
	logEvent(l, newObj, "updated")
}

func logEvent(l *Logger, obj interface{}, ev string) {
	if ev == "created" {
		switch obj.(type) {
		case *api_v1.Deployment:
			log.Printf("Resource %s: %s\n", obj.(*api_v1.Deployment).ObjectMeta.Name, ev)
		}
		return
	}

	log.Printf("Resource %s: %s\n", obj.(event.Event).Name, ev)
}

func (l *Logger) createSecretNewNamespace(ns string) error {
	labelSlector := meta_v1.ListOptions{
		LabelSelector: "app=secret-sync",
	}
	secrets, err := l.Client.CoreV1().Secrets("default").List(labelSlector)
	if err != nil {
		return err
	}

	for _, s := range secrets.Items {
		s.ObjectMeta.SetNamespace(ns)
		s.ObjectMeta.ResourceVersion = ""
		_, err := l.Client.CoreV1().Secrets(ns).Create(&s)
		if err != nil {
			return err
		}
	}
	return nil
}
