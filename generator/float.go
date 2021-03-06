package generator

const writeFloatMethod = `
func writeFloat(r float32, w io.Writer) error {
	bits := uint64(math.Float32bits(r))
	const byteCount = 4
	return encodeFloat(w, byteCount, bits)
}
`

const encodeFloatMethod = `
func encodeFloat(w io.Writer, byteCount int, bits uint64) error {
	var err error
	var bb []byte
	bw, ok := w.(ByteWriter)
	if ok {
		bw.Grow(byteCount)
	} else {
		bb = make([]byte, 0, byteCount)
	}
	for i := 0; i < byteCount; i++ {
		if bw != nil {
			err = bw.WriteByte(byte(bits & 255))
			if err != nil {
				return err
			}
		} else {
			bb = append(bb, byte(bits&255))
		}
		bits = bits >> 8
	}
	if bw == nil {
		_, err = w.Write(bb)
		return err
	}
	return nil
}
`
const readFloatMethod = `
func readFloat(r io.Reader) (float32, error) {
	buf := make([]byte, 4)
	_, err := io.ReadFull(r, buf)
	if err != nil {
		return 0, err
	}
	bits := binary.LittleEndian.Uint32(buf)
	val := math.Float32frombits(bits)
	return val, nil

}
`

type floatField struct {
	name         string
	defaultValue float32
	hasDefault   bool
}

func (s *floatField) Name() string {
	return toPublicName(s.name)
}

func (s *floatField) FieldType() string {
	return "Float"
}

func (s *floatField) GoType() string {
	return "float32"
}

func (s *floatField) SerializerMethod() string {
	return "writeFloat"
}

func (s *floatField) DeserializerMethod() string {
	return "readFloat"
}

func (e *floatField) AddStruct(p *Package) {}

func (e *floatField) AddSerializer(p *Package) {
	p.addStruct(UTIL_FILE, "ByteWriter", byteWriterInterface)
	p.addFunction(UTIL_FILE, "", "writeFloat", writeFloatMethod)
	p.addFunction(UTIL_FILE, "", "encodeFloat", encodeFloatMethod)
	p.addImport(UTIL_FILE, "math")
	p.addImport(UTIL_FILE, "io")
}

func (e *floatField) AddDeserializer(p *Package) {
	p.addFunction(UTIL_FILE, "", "readFloat", readFloatMethod)
	p.addImport(UTIL_FILE, "math")
	p.addImport(UTIL_FILE, "encoding/binary")
	p.addImport(UTIL_FILE, "io")
}
