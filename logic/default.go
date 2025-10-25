package logic

import (
	"asapgiri/golib/logger"
)

var log = logger.Logger {
    Color: logger.Colors.Cyan,
    Pretext: "logic",
}

func filter[T any](s []T, keep func(T) bool) []T {
    var result []T
    for _, v := range(s) {
        if keep(v) {
            result = append(result, v)
        }
    }
    return result
}
