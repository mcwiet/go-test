package encoding

import "encoding/base64"

type CursorEncoder struct{}

func NewCursorEncoder() CursorEncoder {
	return CursorEncoder{}
}

// Encode a string value to a cursor string value
func (c *CursorEncoder) Encode(value string) string {
	return base64.StdEncoding.EncodeToString([]byte(value))
}

// Decode a cursor string value to a string value
func (c *CursorEncoder) Decode(cursor string) (string, error) {
	valBytes, err := base64.StdEncoding.DecodeString(cursor)
	return string(valBytes), err
}
