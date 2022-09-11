package main

import (
	"os"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
)

// go run  .\cmd\logs\
func main() {

	// 输出到控制台
	l := log.DefaultLogger
	l.Log(log.LevelInfo, "stdout_key", "stdout_value")

	l = log.With(l, "ts", log.DefaultTimestamp, "caller", log.DefaultCaller, "trace.id", tracing.TraceID(), "sys", "sys_log")

	h := log.NewHelper(log.NewFilter(l,
		log.FilterKey("key2"),   // 按照key 隐藏 val
		log.FilterValue("val3"), // 按照 val 隐藏 val
		log.FilterFunc(func(level log.Level, keyvals ...interface{}) bool {
			if level == log.LevelWarn {
				return true
			}
			for i := 0; i < len(keyvals); i++ {
				if keyvals[i] == "password" {
					keyvals[i+1] = "*****"
				}
			}
			return false
		}),
	),
	)
	h.Debug("Are you OK!")
	h.Log(log.LevelDebug, "key1", "val1")
	h.Infow("password", "passwordval")
	h.Warnw("key2", "val2")
	h.Warnw("key3", "val3")
	f, err := os.OpenFile("logs.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return
	}
	l = log.NewStdLogger(f)
	l.Log(log.LevelInfo, "file_key", "file_value")
}
