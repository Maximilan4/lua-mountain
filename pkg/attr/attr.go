package attr


func GetTyped[RValue any](m map[string]any, key string) (val RValue, ok bool) {
	if _, ok = m[key]; ok {
		val, ok = m[key].(RValue)
	}
	return
}
