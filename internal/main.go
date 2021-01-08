package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path"
	"sync"
	"time"

	"github.com/tahirali-csc/task-executor-engine/engine"
	"github.com/tahirali-csc/task-executor-engine/engine/kube"

	"log"
)

func main() {
	kubeEngine, err := kube.NewFile("", "", "")
	if err != nil {
		log.Println(err)
	}

	mountPath := os.Getenv("MOUNT_PATH")
	claimName := os.Getenv("CLAIM_NAME")

	spec := &engine.Spec{
		Metadata: engine.Metadata{
			Namespace: "default",
			UID:       "java",
		},
		Image:   "alpine:latest",
		Command: []string{"/bin/sh", "-c", "date && sleep 10s && ls -al"},
		Volumes: []engine.VolumeMount{
			{
				Name:      "logs-drive",
				ClaimName: claimName,
				MountPath: mountPath,
			},
		},
	}

	err = kubeEngine.Setup(context.Background(), spec)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("Starting pod....")
	err = kubeEngine.Start(context.Background(), spec)
	if err != nil {
		log.Println(err)
		return
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()

		log.Println("Tailing...")
		r, err := kubeEngine.Tail(context.Background(), spec)
		if err != nil {
			log.Println("Tail::::", err)
			return
		}

		currTime := time.Now().Unix()
		f, err := os.Create(path.Join(mountPath, fmt.Sprintf("%d.log", currTime)))
		if err != nil {
			log.Println("Can not create file::::", err)
			return
		}

		f.Chmod(777)

		bw := bufio.NewWriter(f)
		br := bufio.NewReader(r)
		for {
			line, _, err := br.ReadLine()
			if err != nil {

				log.Println(err)
				return
			}

			// log.Println(string(line))
			bw.WriteString(string(line) + "\n")
			bw.Flush()
		}
	}()

	log.Println("Waiting pod....")
	kubeEngine.Wait(context.Background(), spec)
	log.Println("Done")
	wg.Wait()

	if err != nil {
		log.Println(err)
	}
}
