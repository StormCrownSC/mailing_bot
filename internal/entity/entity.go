package entity

type Admin struct {
	TelegramID int64
	Name       string
	IsDeleted  bool
}

type User struct {
	TelegramID int64
	Name       string
	IsDeleted  bool
}

type Text struct {
	ID        uint32
	Text      string
	IsDeleted bool
}
