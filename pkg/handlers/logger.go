package handlers

import (
	"log"

	"github.com/a8uhnf/container-h/pkg/event"
	api_v1 "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type Logger struct {
	Client kubernetes.Interface
}

func (l *Logger) ObjectCreated(obj interface{}) {
	nses, err := l.Client.CoreV1().Namespaces().List(meta_v1.ListOptions{})
	if err != nil {
		panic(err)
	}
	switch obj.(type) {
	case *api_v1.Secret:
		for _, n := range nses.Items {
			err := l.createSecret(obj.(*api_v1.Secret), n.ObjectMeta.Name)
			if err != nil {
				log.Println(err)
				recover()
			}
		}
	case *api_v1.Namespace:
		err = l.createSecretNewNamespace(obj.(*api_v1.Namespace).ObjectMeta.Name)
		if err != nil {
			log.Println(err)
			recover()
		}
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
		case *api_v1.Secret:
			log.Printf("Resource %s: %s\n", obj.(*api_v1.Secret).ObjectMeta.Name, ev)
		case *api_v1.Namespace:
			log.Printf("Resource %s: %s\n", obj.(*api_v1.Namespace).ObjectMeta.Name, ev)
		}
		return
	}

	log.Printf("Resource %s: %s\n", obj.(event.Event).Name, ev)
}

func (l *Logger) createSecret(s *api_v1.Secret, ns string) error {
	if s.ObjectMeta.GetNamespace() == ns {
		log.Println("secret already present in this ns.")
		return nil
	}
	s.ObjectMeta.SetNamespace(ns)
	s.ObjectMeta.ResourceVersion = ""
	_, err := l.Client.CoreV1().Secrets(ns).Create(s)
	if err != nil {
		return err
	}
	return nil
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
