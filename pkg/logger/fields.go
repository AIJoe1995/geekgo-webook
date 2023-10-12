package logger

//仿照zap.String
//zap.String() 其他方法还有 zap.Int() zap.Any()等
// // String constructs a field with the given key and value.

func String(key, val string) Field {
	return Field{
		Key:   key,
		Value: val,
	}
}
