package devices

type Graphics struct {
	surfaces map[byte]Surface
}

type Surface struct {
	id byte
}

func ParseDrawCommands(data []byte) {

}
