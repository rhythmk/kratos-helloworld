package trace

import (
	"fmt"
	"time"

	opentracing "github.com/opentracing/opentracing-go"
	traceLog "github.com/opentracing/opentracing-go/log"
	"gorm.io/gorm"
	"gorm.io/gorm/utils"
)

const (
	callBackBeforeName = "core:before"
	callBackAfterName  = "core:after"
	startTime          = "_start_time"
)

type TracePlugin struct{}

func (op *TracePlugin) Name() string {
	return "tracePlugin"
}

func (op *TracePlugin) Initialize(db *gorm.DB) (err error) {
	// 开始前
	_ = db.Callback().Create().Before("gorm:before_create").Register(callBackBeforeName, before)
	_ = db.Callback().Query().Before("gorm:query").Register(callBackBeforeName, before)
	_ = db.Callback().Delete().Before("gorm:before_delete").Register(callBackBeforeName, before)
	_ = db.Callback().Update().Before("gorm:setup_reflect_value").Register(callBackBeforeName, before)
	_ = db.Callback().Row().Before("gorm:row").Register(callBackBeforeName, before)
	_ = db.Callback().Raw().Before("gorm:raw").Register(callBackBeforeName, before)

	// 结束后
	_ = db.Callback().Create().After("gorm:after_create").Register(callBackAfterName, after)
	_ = db.Callback().Query().After("gorm:after_query").Register(callBackAfterName, after)
	_ = db.Callback().Delete().After("gorm:after_delete").Register(callBackAfterName, after)
	_ = db.Callback().Update().After("gorm:after_update").Register(callBackAfterName, after)
	_ = db.Callback().Row().After("gorm:row").Register(callBackAfterName, after)
	_ = db.Callback().Raw().After("gorm:raw").Register(callBackAfterName, after)
	return
}

var _ gorm.Plugin = &TracePlugin{}

const gormSpanKey = "___gorm_span"

func before(db *gorm.DB) {
	db.InstanceSet(startTime, time.Now())
	span, _ := opentracing.StartSpanFromContext(db.Statement.Context, "gorm")
	db.InstanceSet(gormSpanKey, span)
	return
}

func after(db *gorm.DB) {
	_ts, isExist := db.InstanceGet(startTime)
	if !isExist {
		return
	}
	ts, ok := _ts.(time.Time)
	if !ok {
		return
	}

	sql := db.Dialector.Explain(db.Statement.SQL.String(), db.Statement.Vars...)
	_traceSapn, ok := db.InstanceGet(gormSpanKey)
	if !ok {
		return
	}
	traceSapn, ok := _traceSapn.(opentracing.Span)
	if !ok {
		return
	}
	defer traceSapn.Finish()
	if db.Error != nil {
		traceSapn.LogFields(traceLog.Error(db.Error))
	}
	// traceSql := []SQL{}
	// getSql := ctx.Value("sql")
	// if getSql != nil {
	// 	traceSql = getSql.([]SQL)
	// }

	var msg string = ""
	if db.Error != nil {
		msg = fmt.Sprintf("%s", db.Error)
	}

	// traceSql = append(traceSql, SQL{
	// 	Timestamp:   time.Now().Format("2006-01-02 15:04:05"),
	// 	SQL:         sql,
	// 	Stack:       utils.FileWithLineNum(),
	// 	Rows:        db.Statement.RowsAffected,
	// 	CostSeconds: time.Since(ts).Seconds(),
	// 	Msg:         msg,
	// })

	//ctx.Set("sql", traceSql)
	//context.WithValue(ctx, "sql", traceSql)

	fmt.Println("抓取到的SQL:", sql)

	traceSapn.LogFields(
		traceLog.String("Timestamp", time.Now().Format("2006-01-02 15:04:05")),
		traceLog.String("SQL", sql),
		traceLog.String("Stack", utils.FileWithLineNum()),
		traceLog.Int64("Rows", db.Statement.RowsAffected),
		traceLog.Float64("CostSeconds", time.Since(ts).Seconds()),
		traceLog.String("Msg", msg),
	)

	return
}

// 告诉编译器 TracePlugin 实现了 gorm.Plugin 接口
var _ gorm.Plugin = &TracePlugin{}
