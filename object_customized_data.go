package caskin

import (
	"bytes"
	"encoding/json"
)

type CustomizedData interface {
	GetName() string
	GetObjectType() ObjectType
}

func CustomizedData2Object(customized CustomizedData, object Object) {
	object.SetName(customized.GetName())
	object.SetObjectType(customized.GetObjectType())
	b, _ := json.Marshal(customized)
	object.SetCustomizedData(b)
}

func CustomizedDataEqualObject(customized CustomizedData, object Object) bool {
	if customized.GetName() != object.GetName() ||
		customized.GetObjectType() != object.GetObjectType() {
		return false
	}
	b1, _ := object.GetCustomizedData().MarshalJSON()
	b2, _ := json.Marshal(customized)
	return bytes.Compare(b1, b2) == 0
}

func Object2CustomizedData(object Object, factory func() CustomizedData) (CustomizedData, error) {
	from := object.GetCustomizedData()
	to := factory()
	return to, json.Unmarshal(from, to)
}

func ObjectArray2CustomizedDataArray(objects []Object, factory func() CustomizedData) ([]CustomizedData, error) {
	var customized []CustomizedData
	for _, v := range objects {
		to, err := Object2CustomizedData(v, factory)
		if err != nil {
			return nil, err
		}
		customized = append(customized, to)
	}
	return customized, nil
}

func ObjectArray2CustomizedDataPair(objects []Object, factory func() CustomizedData) ([]*CustomizedDataPair, error) {
	customized, err := ObjectArray2CustomizedDataArray(objects, factory)
	if err != nil {
		return nil, err
	}
	var pair []*CustomizedDataPair
	for i, v := range objects {
		pair = append(pair, &CustomizedDataPair{
			Object:               v,
			ObjectCustomizedData: customized[i],
		})
	}
	return pair, nil
}
