package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
)

type QueryType string

type QueryFilter struct {
	Query QueryType
	Value interface{}
}

const (
	IDQuery                  QueryType = "id_query"
	BoolQuery                QueryType = "bool_query"
	ExactQuery               QueryType = "exact_query"
	InQuery                  QueryType = "in_query"
	LikeQuery                QueryType = "like_query"
	CaseInsensitiveLikeQuery QueryType = "case_insensitive_like_query"
	CaseInsensitiveFindQuery QueryType = "case_insensitive_find_query"
	CustomQuery              QueryType = "custom_query"
)

func BuildMongoQuery(ctx context.Context, filter map[string]QueryFilter) bson.D {
	result := bson.D{}
	for key, queryFilter := range filter {
		switch queryFilter.Query {
		case ExactQuery:
			result = append(result, bson.E{Key: key, Value: queryFilter.Value})

		case InQuery:
			result = append(result, bson.E{Key: key, Value: bson.M{"$in": queryFilter.Value}})

		case BoolQuery:
			if val, ok := queryFilter.Value.(bool); ok {
				result = append(result, bson.E{Key: key, Value: val})
			}

		case IDQuery:
			if val, ok := GetObjectIDFromInterface(queryFilter.Value); ok {
				result = append(result, bson.E{Key: key, Value: val})
			}

		case LikeQuery:
			if val, ok := queryFilter.Value.(string); ok {
				result = append(result, bson.E{Key: key, Value: bson.M{"$regex": "^" + val}})
			}

		case CaseInsensitiveLikeQuery:
			if val, ok := queryFilter.Value.(string); ok {
				result = append(result, bson.E{Key: key, Value: bson.M{"$regex": "^" + val, "$options": "i"}})
			}

		case CaseInsensitiveFindQuery:
			if val, ok := queryFilter.Value.(string); ok {
				result = append(result, bson.E{Key: key, Value: bson.M{"$regex": "^" + val + "$", "$options": "i"}})
			}

		case CustomQuery:
			result = append(result, bson.E{Key: key, Value: queryFilter.Value})
		}

	}

	return result
}

func BuildMongoSetQuery(setConditions map[string]interface{}) bson.D {
	queries := bson.D{}
	for key, value := range setConditions {
		queries = append(queries, bson.E{Key: key, Value: value})
	}

	result := bson.D{{"$set", queries}}
	return result
}
