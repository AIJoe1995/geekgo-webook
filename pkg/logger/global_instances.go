package logger

import (
	"sync"
)

// 存放全局变量 类似ZAP的实现， zap里面的global.go 初始化了全局logger _globalL  = NewNop()
// 在生成自己的logger zap.NewDevelopment之后 需要执行zap.ReplaceGlobals 才能正常使用zap.L().Debug打印日志

//var (
//	_globalMu sync.RWMutex
//	_globalL  = NewNop()
//	_globalS  = _globalL.Sugar()
//)
//func L() *Logger {
//	_globalMu.RLock()
//	l := _globalL
//	_globalMu.RUnlock()
//	return l
//}
//func ReplaceGlobals(logger *Logger) func() {
//	_globalMu.Lock()
//	prev := _globalL
//	_globalL = logger
//	_globalS = logger.Sugar()
//	_globalMu.Unlock()
//	return func() { ReplaceGlobals(prev) }
//}

// 在nop.go里面仿照zap下面的代码做一个nop实现
//func NewNop() *Logger {
//	return &Logger{
//		core:        zapcore.NewNopCore(),
//		errorOutput: zapcore.AddSync(ioutil.Discard),
//		addStack:    zapcore.FatalLevel + 1,
//		clock:       zapcore.DefaultClock,
//	}
//}

var gl LoggerV1

var lMutex sync.RWMutex

func SetGlobalLogger(l LoggerV1) {
	lMutex.Lock()
	defer lMutex.Unlock()
	gl = l
}

func L() LoggerV1 {
	lMutex.RLock()
	g := gl
	lMutex.RUnlock()
	return g
}

var GL LoggerV1 = &NopLogger{}
