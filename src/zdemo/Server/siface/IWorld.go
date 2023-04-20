package siface

type IWorld interface {
	CreateRoom(string, uint32) error
	GetRoom(string) (IRoom, error)
	RemoveRoom(IRoom)
	UserSwitchRoom(IUser, string, string) error
	GetAllRoomMsg() string
}
