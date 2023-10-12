package zapx

import "go.uber.org/zap/zapcore"

// 对于sensitive的信息 写向日志文件的时候做脱敏处理
// // Core is a minimal, fast logger interface. It's designed for library authors
//// to wrap in a more user-friendly API.

type MyCore struct {
	zapcore.Core // 组合了zapcore.Core接口 修改Write方法

}

//	Write(Entry, []Field) error
//
// 装饰器模式， 但是怎么让日志模块使用MyCore 而不是zapcore.Core ???
func (c MyCore) Write(entry zapcore.Entry, fds []zapcore.Field) error {
	for _, fd := range fds {
		if fd.Key == "phone" {
			phone := fd.String
			fd.String = phone[:3] + "****" + phone[7:]
		}
	}
	return c.Core.Write(entry, fds)
}
