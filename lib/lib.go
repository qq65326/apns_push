package lib


import(
	"sync"
)

var(
	ServerWaitGroup = sync.WaitGroup{}
	PushWaitGroup = sync.WaitGroup{}
	MainWaitGroup = sync.WaitGroup{}
)