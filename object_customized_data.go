package caskin

import (
	"bytes"
	"encoding/json"
)

type ObjectCustomizedData interface {
	GetName() string
	GetObjectType() ObjectType
}

func ObjectCustomizedData2Object(customized ObjectCustomizedData, object Object) {
	object.SetName(customized.GetName())
	object.SetObjectType(customized.GetObjectType())
	b, _ := json.Marshal(customized)
	object.SetCustomizedData(b)
}

func ObjectCustomizedDataEqualObject(customized ObjectCustomizedData, object Object) bool {
	if customized.GetName() != object.GetName() ||
		customized.GetObjectType() != object.GetObjectType() {
		return false
	}
	b1, _ := object.GetCustomizedData().MarshalJSON()
	b2, _ := json.Marshal(customized)
	return bytes.Compare(b1, b2) == 0
}

func ObjectArray2ObjectCustomizedDataArray(objects []Object, factory func() ObjectCustomizedData) ([]ObjectCustomizedData, error) {
	var customized []ObjectCustomizedData
	for _, v := range objects {
		from := v.GetCustomizedData()
		to := factory()
		if err := json.Unmarshal(from, to); err != nil {
			return nil, err
		}
		customized = append(customized, to)
	}
	return customized, nil
}

func ObjectArray2Pair(objects []Object, factory func() ObjectCustomizedData) ([]*CustomizedDataPair, error) {
	customized, err := ObjectArray2ObjectCustomizedDataArray(objects, factory)
	if err != nil {
		return nil, err
	}
	var pair []*CustomizedDataPair
	for i, v := range objects {
		pair = append(pair, &CustomizedDataPair{
			Object: v,
			ObjectCustomizedData: customized[i],
		})
	}
	return pair, nil
}