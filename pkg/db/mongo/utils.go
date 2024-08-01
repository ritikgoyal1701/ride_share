package mongo

import "go.mongodb.org/mongo-driver/bson/primitive"

func AppendQueryFilterMap(map1, map2 map[string]QueryFilter) (combinedMap map[string]QueryFilter) {
	combinedMap = make(map[string]QueryFilter)
	for key, value := range map1 {
		combinedMap[key] = value
	}

	for key, value := range map2 {
		combinedMap[key] = value
	}

	return
}

func GetObjectIDFromInterface(data interface{}) (primitive.ObjectID, bool) {
	str, ok := data.(string)
	if !ok {
		return primitive.NilObjectID, false
	}

	if _, err := primitive.ObjectIDFromHex(str); err != nil {
		return primitive.NilObjectID, false
	}

	id, _ := primitive.ObjectIDFromHex(str)
	return id, true
}
