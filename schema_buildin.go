package caskin

// NamedObject build in Object for name encode/decode
type NamedObject struct {
	Name string
}

func (o *NamedObject) GetID() uint64 {
	return 0
}

func (o *NamedObject) SetID(id uint64) {
}

func (o *NamedObject) Encode() string {
	return o.Name
}

func (o *NamedObject) Decode(code string) error {
	o.Name = code
	return nil
}

func (o *NamedObject) GetParentID() uint64 {
	return 0
}

func (o *NamedObject) SetParentID(pid uint64) {
}

func (o *NamedObject) GetDomainID() uint64 {
	return 0
}

func (o *NamedObject) SetDomainID(did uint64) {
}

func (o *NamedObject) GetName() string {
	return o.Name
}

func (o *NamedObject) SetName(name string) {
	o.Name = name
}

func (o *NamedObject) GetObjectType() ObjectType {
	return ""
}
