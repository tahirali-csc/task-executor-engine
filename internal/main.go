package main

import (
	"context"
	"github.com/tahirali-csc/task-executor-engine/engine"
	"github.com/tahirali-csc/task-executor-engine/engine/kube"

	"log"
)

func main() {
	kubeEngine, err := kube.NewFile("", "", "")
	if err != nil {
		log.Println(err)
	}

	spec := &engine.Spec{
		Metadata: engine.Metadata{
			Namespace: "default",
			UID:       "java",
		},
		Image:   "alpine:latest",
		Command: []string{"/bin/sh", "-c", "date && ls -al"},
	}

	log.Println("Starting pod....")
	kubeEngine.Start(context.Background(), spec)

	log.Println("Waiting pod....")
	kubeEngine.Wait(context.Background(), spec)
	log.Println("Done")

	if err != nil {
		log.Println(err)
	}
}
