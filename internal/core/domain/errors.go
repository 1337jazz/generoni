package domain

type Kind string

const (
	KindInternal        Kind = "internal_error"
	KindInvalidArgument Kind = "invalid_request"
	KindNotFound        Kind = "not_found"
)

type Error struct {
	Code string
	Kind Kind
}

func NewError(code string, kind Kind) *Error {
	return &Error{Code: code, Kind: kind}
}

func (e *Error) Error() string {
	return e.Code
}

var (
	ErrItemNotFound  = NewError("item_not_found", KindNotFound)
	ErrInvalidItemID = NewError("invalid_item_id", KindInvalidArgument)
	ErrTitleRequired = NewError("title_required", KindInvalidArgument)
	ErrTitleTooLong  = NewError("title_too_long", KindInvalidArgument)
)
