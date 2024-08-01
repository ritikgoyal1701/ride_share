package mongo

import (
	"go.mongodb.org/mongo-driver/mongo/options"
)

type FindOptions struct {
	Option        *options.FindOptions    // Used When Multi True
	FindOneOption *options.FindOneOptions // Used When Multi False
	Multi         bool                    // Default True
}

func NewFindOptions(opts ...func(o *FindOptions)) *FindOptions {
	mongoFindOptions := &FindOptions{
		Option:        options.Find(),
		FindOneOption: options.FindOne(),
		Multi:         true,
	}

	for _, opt := range opts {
		if opt == nil {
			continue
		}

		opt(mongoFindOptions)
	}

	return mongoFindOptions
}

func SetMulti(multi bool) func(o *FindOptions) {
	return func(o *FindOptions) {
		o.Multi = multi
	}
}

func (o *FindOptions) GetMongoFindOptions() *options.FindOptions {
	return o.Option
}

func (o *FindOptions) GetMongoFindOneOptions() *options.FindOneOptions {
	return o.FindOneOption
}

func (o *FindOptions) SetProjection(fields map[string]interface{}) *FindOptions {
	if o == nil {
		return nil
	}

	if len(fields) == 0 {
		return o
	}

	if !o.Multi {
		o.FindOneOption.SetProjection(fields)

		return o
	}

	o.Option.SetProjection(fields)

	return o
}
