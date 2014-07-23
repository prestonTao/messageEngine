package messageEngine

import (
	// "fmt"
	"time"
)

type TimeOut struct {
	isTimeOutChan chan bool
	duration      time.Duration
	f             func()
}

func (this *TimeOut) Do(duration time.Duration) bool {
	this.duration = duration
	go this.run()

	select {
	case <-this.isTimeOutChan:
		close(this.isTimeOutChan)
		return false
	case <-time.After(this.duration):
		return true
	}

	// return <-this.isTimeOutChan
}

func (this *TimeOut) run() {
	// defer func() {
	// 	if r := recover(); r != nil {
	// 		//超时了
	// 		// fmt.Println("out超时了")
	// 	} else {
	// 		//没超时
	// 		// fmt.Println("out没超时")
	// 	}

	// }()

	// go func() {
	// 	defer func() {
	// 		if r := recover(); r != nil {
	// 			//没超时
	// 			// fmt.Println("in没超时")
	// 		} else {
	// 			//超时了
	// 			// fmt.Println("in超时了")
	// 		}

	// 	}()
	// 	time.Sleep(this.duration)
	// 	this.isTimeOutChan <- true
	// 	close(this.isTimeOutChan)
	// }()
	this.f()
	this.isTimeOutChan <- false

}

func NewTimeOut(f func()) *TimeOut {
	to := TimeOut{
		isTimeOutChan: make(chan bool),
		f:             f,
	}
	return &to
}
