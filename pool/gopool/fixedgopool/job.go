package fixedgopool

type JobI interface {
	Run()
}

type Func func()
