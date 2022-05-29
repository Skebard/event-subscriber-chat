package chat

type RoomError struct {
	RoomIdentifier string
	customMessage  string
}

func (err *RoomError) Error() string {
	if len(err.customMessage) > 0 {
		return err.customMessage + ". Identifier: " + err.RoomIdentifier
	}
	return "Invalid room. Identifier: " + err.RoomIdentifier
}

func NewRoomError(roomIdentifier string) error {
	return &RoomError{roomIdentifier, ""}
}

type RoomNotFoundError struct {
	RoomError
}

func NewRoomNotFoundError(roomIdentifier string) error {
	return &RoomNotFoundError{RoomError{roomIdentifier, "Room not found"}}
}

type RoomIdentifierAlreadyInUse struct {
	RoomError
}

func NewRoomIdentifierAlreadyInUse(roomIdentifier string) error {
	return &RoomIdentifierAlreadyInUse{RoomError{roomIdentifier, "Room identifier already in use"}}

}
