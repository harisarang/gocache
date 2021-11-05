package cache

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	cache      int
	mutex      sync.Mutex
	cacheTimer *time.Timer
	stop       chan bool
)

func Client() {
	var wg sync.WaitGroup
	mutex = sync.Mutex{}
	cont := "y"
	cache = 0
	stop = make(chan bool)
	cacheTimer = time.NewTimer(time.Second * 5)
	wg.Add(2)
	go InitializeCache()
	go func() {
		for cont == "y" {
			if cache == 0 {
				mutex.Lock()
				cache = MakeHit()
				cacheTimer.Stop()
				stop <- true
				cacheTimer = time.NewTimer(5 * time.Second)
				mutex.Unlock()
			} else {
				log.Printf("Cached hit. Hits: %d\n", cache)
			}
			fmt.Print("Continue? (y/n): ")
			fmt.Scan(&cont)
		}
		if cont == "n" {
			stop <- false
			wg.Done()
			wg.Done()
		}
	}()
	wg.Wait()
}

func MakeHit() int {
	resp, err := http.Get("http://localhost:2811/")
	if err != nil {
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	sb := string(body)
	log.Println(sb)
	i, err := strconv.Atoi(strings.Split(sb, ": ")[1])
	if err != nil {
		log.Fatal(err)
	}
	return i
}

func InitializeCache() {
	for {
		select {
		case <-cacheTimer.C:
			mutex.Lock()
			cache = 0
			mutex.Unlock()
		case x := <-stop:
			if x == false {
				return
			} else {
				continue
			}
		}
	}
}
