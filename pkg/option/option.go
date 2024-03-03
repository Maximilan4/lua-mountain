package option

type (
	Option[Arg any]    func(a Arg)
	ErrOption[Arg any] func(a Arg) error
)
