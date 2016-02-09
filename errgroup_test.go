package errgroup_test

import (
	"errors"
	"testing"

	"github.com/facebookgo/ensure"
	"github.com/facebookgo/errgroup"
)

func TestNada(t *testing.T) {
	t.Parallel()
	var g errgroup.Group
	ensure.Nil(t, g.Wait())
}

func TestOneError(t *testing.T) {
	t.Parallel()
	e := errors.New("")
	var g errgroup.Group
	g.Error(e)
	ensure.True(t, g.Wait() == e)
}

func TestTwoErrors(t *testing.T) {
	t.Parallel()
	e1 := errors.New("e1")
	e2 := errors.New("e2")
	var g errgroup.Group
	g.Error(e1)
	g.Error(e2)
	ensure.DeepEqual(t, g.Wait().Error(), "multiple errors: e1 | e2")
}

func TestInvalidNilError(t *testing.T) {
	defer ensure.PanicDeepEqual(t, "error must not be nil")
	(&errgroup.Group{}).Error(nil)
}

func TestInvalidZeroLengthMultiError(t *testing.T) {
	defer ensure.PanicDeepEqual(t, "MultiError with no errors")
	(errgroup.MultiError{}).Error()
}

func TestInvalidOneLengthMultiError(t *testing.T) {
	defer ensure.PanicDeepEqual(t, "MultiError with only 1 error")
	(errgroup.MultiError{errors.New("")}).Error()
}

func TestAddDone(t *testing.T) {
	t.Parallel()
	var g errgroup.Group
	l := 10
	g.Add(l)
	for i := 0; i < l; i++ {
		go g.Done()
	}
	ensure.Nil(t, g.Wait())
}

func TestNewMultiError(t *testing.T) {
	t.Parallel()

	ensure.Nil(t, errgroup.NewMultiError())
	ensure.Nil(t, errgroup.NewMultiError(nil))
	ensure.Nil(t, errgroup.NewMultiError(nil, nil))

	errA := errors.New("err a")
	errB := errors.New("err b")

	// When only one non-nil error is provided, NewMultiError should return that
	// error without wrapping it in a MultiError.
	ensure.DeepEqual(t, errgroup.NewMultiError(errA), errA)
	ensure.DeepEqual(t, errgroup.NewMultiError(errA, nil), errA)
	ensure.DeepEqual(t, errgroup.NewMultiError(nil, errB), errB)

	// When more than one non-nil error is provided, than NewMultiError should
	// return a MultiError instance.
	multiErrAB := errgroup.NewMultiError(errA, errB).(errgroup.MultiError)
	multiErrNilAB := errgroup.NewMultiError(nil, errA, errB).(errgroup.MultiError)
	multiErrANilB := errgroup.NewMultiError(errA, nil, errB).(errgroup.MultiError)
	multiErrABNil := errgroup.NewMultiError(errA, errB, nil).(errgroup.MultiError)

	expected := errgroup.MultiError{errA, errB}
	ensure.DeepEqual(t, multiErrAB, expected)
	ensure.DeepEqual(t, multiErrNilAB, expected)
	ensure.DeepEqual(t, multiErrANilB, expected)
	ensure.DeepEqual(t, multiErrABNil, expected)
}
