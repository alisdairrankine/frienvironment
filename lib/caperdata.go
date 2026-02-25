package lib

type U8 byte
type U16 uint16
type Bool bool
type CapString string
type CapPrimitive interface {
	U8 | U16 | CapString
}

type CapStruct map[string]any
type List[T CapPrimitive] []T

const (
	IdentifierNil         = 0x00 // (identifier)
	IdentifierU8          = 0x01 // (identifier, data)
	IdentifierU16         = 0x02 // (identifier, data)
	IdentifierTrue        = 0x03 // (identifier)
	IdentifierFalse       = 0x04 // (identifier)
	IdentifierString      = 0x05 // (identifier,length, data)
	IdentifierStruct      = 0x06 // (identifier, length, fields)
	IdentifierStructField = 0x07 // (string, primitive)
	IdentifierList        = 0x08 // (identifier,indentifier, length, data)
)

func ParseCaperData(raw []byte) any {

	return nil
}
