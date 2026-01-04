package result

type Result[T any] struct {
	value T
	err   error
}

// Ok creates a successful result
func Ok[T any](value T) Result[T] {
	return Result[T]{value: value, err: nil}
}

// Err creates a failed result
func Err[T any](err error) Result[T] {
	var zero T
	return Result[T]{value: zero, err: err}
}

// From creates a Result from a value and an error
func From[T any](value T, err error) Result[T] {
	if err != nil {
		return Err[T](err)
	}
	return Ok(value)
}

// IsSuccess checks if result is successful
func (r Result[T]) IsSuccess() bool {
	return r.err == nil
}

// Value returns the value (panics if error)
func (r Result[T]) Value() T {
	if r.err != nil {
		panic("attempted to get value from error result")
	}
	return r.value
}

// Error returns the error (nil if success)
func (r Result[T]) Error() error {
	return r.err
}

// Map transforms the value if successful
func Map[T any, U any](r Result[T], fn func(T) U) Result[U] {
	if r.err != nil {
		return Err[U](r.err)
	}
	return Ok(fn(r.value))
}

// MapErr transforms the error if failed
func MapErr[T any](r Result[T], fn func(error) error) Result[T] {
	if r.err != nil {
		return Err[T](fn(r.err))
	}
	return Ok(r.value)
}

// Bind chains operations (flatMap)
func Bind[T any, U any](r Result[T], fn func(T) Result[U]) Result[U] {
	if r.err != nil {
		return Err[U](r.err)
	}
	return fn(r.value)
}
