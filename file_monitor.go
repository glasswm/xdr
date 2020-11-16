package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/fsnotify/fsnotify"
)

type configuration struct {
	Files []string
}

func main() {

	file, err1 := os.Open("files.json")
	if err1 != nil {
		log.Fatal(err1)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	conf := configuration{}
	err := decoder.Decode(&conf)
	if err != nil {
		fmt.Println("Error:", err)
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Println("event:", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("modified file:", event.Name)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	for _, value := range conf.Files {
		err = watcher.Add(value)
		if err != nil {
			log.Println(err)
		}
	}
	<-done

}
