package siface

type IRoom interface {
	StartRoom()
	StopRoom()

	GetName() string

	GetUser(name string) (IUser, error)
	AddUser(IUser) error
	RemoveUser(IUser)
	ClearAll()

	GetUserNum() uint32
	GetCapacity() uint32

	BroadCastMsg([]byte)
}
