package ziface

type IServer interface {
	// start Server
	Start()
	// stop Server
	Stop()
	// run server
	Serve()
	// add rounter
	AddRounter(uint32, IRouter)
}
