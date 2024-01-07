package repository


type (
	ContextKey string
)

var (
	RequestIdContextKey = ContextKey("id")
)


func (c ContextKey) GetKey() string {
	return string(c)
}

