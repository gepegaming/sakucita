package utils

func Int32Ptr(v int) *int32 {
	i := int32(v)
	return &i
}

func Int64Ptr(v int64) *int64 {
	return &v
}
