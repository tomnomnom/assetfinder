package main

import (
	"reflect"
	"runtime"
	"strings"
)

func getFuncName (fn fetchFn) string{
    //get Function name from fetchFn type of format "main.FunctionName"
    //removing "main." using strings.Split returning it
    return strings.Split(runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name(), ".")[1]  
}

func toggleSources (override map[string]*bool, sources []fetchFn) []fetchFn{
    var sourcelist  []fetchFn

    for _,source := range sources{

        if *override[getFuncName(source)]{ //Checking flags for all functions in sources []fetchFn slice
            sourcelist = append(sourcelist,source)
        }

    }
    return sourcelist
}

