package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

type mainStruct struct {
	disks         []string
	readSize      int
	readSizeTrack chan int
	reportSleep   int
	wg            *sync.WaitGroup
	off           chan int
}

func main() {
	m := new(mainStruct)
	ret := m.main()
	os.Exit(ret)
}

func (m *mainStruct) main() int {
	flag.IntVar(&m.readSize, "read-size", 1048576, "size in bytes of read at a time")
	flag.IntVar(&m.reportSleep, "print-sleep", 1, "seconds between speed sampling and prints")
	flag.Parse()
	m.disks = flag.Args()
	if len(m.disks) != 1 {
		fmt.Println("ERROR: Usage: diskspeed [optional-params] /dev/disk")
		flag.PrintDefaults()
		return 1
	}
	fmt.Printf("[%s] Opening %s for reading\n", time.Now().String(), m.disks[0])
	disk := m.disks[0]
	fh, err := os.Open(disk)
	if err != nil {
		fmt.Println(err)
		return 1
	}
	defer fh.Close()
	no := 0
	var total uint64
	m.readSizeTrack = make(chan int, 1024)
	m.wg = new(sync.WaitGroup)
	m.off = make(chan int, 2)
	m.wg.Add(1)
	go m.reporter()
	startTime := time.Now()
	for err == nil {
		b := make([]byte, m.readSize)
		no, err = fh.Read(b)
		if no >= 0 {
			total = total + uint64(no)
		}
		if err != nil && err != io.EOF {
			fmt.Println(err)
			return 1
		} else if err != nil && err == io.EOF {
			m.off <- 1
			m.readSizeTrack <- 0
			m.wg.Wait()
			m.printSizeTotal(total, startTime)
		}
		m.readSizeTrack <- no
	}
	return 0
}

func (m *mainStruct) reporter() {
	ts := time.Now()
	no := 0
	for {
		no = no + <-m.readSizeTrack
		if ts.Add(time.Duration(m.reportSleep) * time.Second).Before(time.Now()) {
			ts = time.Now()
			m.wg.Add(1)
			go m.printSize(no / m.reportSleep)
			no = 0
		}
		if len(m.off) > 0 {
			m.printSize(no / m.reportSleep)
			return
		}
	}
}

func (m *mainStruct) printSize(size int) {
	if size >= 1048576 {
		fmt.Printf("[%s] %d MB/s for %d s\n", time.Now().String(), size/1048576, m.reportSleep)
	} else if size >= 1024 {
		fmt.Printf("[%s] %d KB/s for %d s\n", time.Now().String(), size/1024, m.reportSleep)
	} else {
		fmt.Printf("[%s] %d B/s for %d s\n", time.Now().String(), size, m.reportSleep)
	}
	if len(m.readSizeTrack)*m.readSize > 1048576 {
		fmt.Printf("[%s] WARN: q_len_bytes_behind=%d\n", time.Now().String(), len(m.readSizeTrack)*m.readSize)
	}
	m.wg.Done()
}

func (m *mainStruct) printSizeTotal(size uint64, startTime time.Time) {
	var total string
	var average string
	if size >= 1048576 {
		total = fmt.Sprintf("%d MB", size/1048576)
	} else if size >= 1024 {
		total = fmt.Sprintf("%d KB", size/1024)
	} else {
		total = fmt.Sprintf("%d B", size)
	}
	timeDelta := uint64(time.Now().Sub(startTime).Seconds())
	if timeDelta == 0 {
		timeDelta = 1
	}
	if size/timeDelta >= 1048576 {
		average = fmt.Sprintf("%d MB/s", size/timeDelta/1048576)
	} else if size/timeDelta >= 1024 {
		average = fmt.Sprintf("%d KB/s", size/timeDelta/1024)
	} else {
		average = fmt.Sprintf("%d B/s", size/timeDelta)
	}
	fmt.Printf("[%s] END: total = %s ; average speed = %s\n", time.Now().String(), total, average)
}
