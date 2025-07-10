package util

type RollbackFunc func()

func NewRollback() ([]RollbackFunc, func()) {
	var rollbackFuncs []RollbackFunc

	fail := func() {
		for i := len(rollbackFuncs) - 1; i >= 0; i-- {
			func(f RollbackFunc) {
				defer func() {
					_ = recover() // silently ignore panic
				}()
				f()
			}(rollbackFuncs[i])
		}
	}

	return rollbackFuncs, fail
}
