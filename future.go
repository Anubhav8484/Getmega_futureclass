package gofuture

import (
	"reflect"
	"time"
	"errors"
)

type Future struct {
	Cancelled         bool
	Done              bool
	Running           bool
	Result            interface{}
	InterfaceChannel <-chan interface{}
}


func (f *Future) cancel() bool{
      if f.Running{
	     return false
	  }
	  f.Cancelled= true;
	  return true;
}


func (f *Future) cancelled() bool{
      if f.Running{
	     return false
	  }
	  return f.Cancelled
}

func (f *Future) running() bool{
      return f.Running
}

func (f *Future) done() bool{
      return (f.Cancelled || f.Done)
	      
}

//result function when timeout is not given
func (f *Future) result() interface{} {
	if (f.Done || f.Cancelled) {
		return f.result
	}
	f.Result = <-f.InterfaceChannel
	return f.Result
}

//result function when timeout is  given
func result (f *Future,timeout time.Duration)  interface{} {
	if (f.Done || f.Cancelled) {
		return f.Result
	}
	timeoutChannel := time.After(timeout)
	select {
	case res := <-f.InterfaceChannel:
		f.Result = res
		f.Done = true
	case <-timeoutChannel:
		f.Result = nil
		f.Done = true
	}
	return f.Result
}

//Exception function when timeout is  given
func exception(f *Future,timeout time.Duration) error{
        if f.Cancelled {
	   return errors.New("Cancelled")
	}
	if f.Done{
	   return nil
	}
        timeoutChannel := time.After(timeout)
	select {
	case <-f.InterfaceChannel:
		return nil
	case <-timeoutChannel:
		return errors.New("Timeout")
	}

}

//Exception function when timeout is not given
func (f *Future) exception() error{
        if f.Cancelled {
	    return errors.New("Cancelled")
	}
        if(f.Running){
	   return errors.New("Invalid state")
	}
	return nil

}

//submit function
func ThreadPoolExecutorSubmit(implem interface{}, args ...interface{}) *Future {
    var f Future
	
    f.Running = false
	f.Cancelled = false
	f.Done = false
	f.Result =nil
	
	valIn := make([]reflect.Value, len(args), len(args))    

	func_name := reflect.ValueOf(implem)         //function name to be called
    
	for idx, elt := range args {
		valIn[idx] = reflect.ValueOf(elt)        //function arguments
	}
	
	interfaceChannel := make(chan interface{}, 1)

	go func() {
	    f.Running = true
		res := func_name.Call(valIn)
		for idx, _ := range res {
			interfaceChannel <- res[idx].Interface()
		}
		f.Running = false
		f.Done = true
	}()
    f.InterfaceChannel= interfaceChannel
	return &f
}
