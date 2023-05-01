package repository

import (
	"app/internal/core/domain"
	"context"
	"reflect"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repo struct {
	Collection *mongo.Collection
	Mux        sync.Mutex
}

func (r *Repo) Create(i interface{}) error {
	rt := reflect.TypeOf(i)
	switch rt.Kind() {
	case reflect.Slice, reflect.Array:
		return r.createMany(i)
	default:
		return r.createOne(i)
	}
}

func (r *Repo) createOne(i interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if m, ok := i.(domain.ModelInterface); ok {
		m.Stamp()
		m.SetID(primitive.NewObjectID())
	}
	if _, err := r.Collection.InsertOne(ctx, i); err != nil {
		return err
	}
	return nil
}

func (r *Repo) createMany(i interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	v := reflect.ValueOf(i)
	is := []interface{}{}
	for i := 0; i < v.Len(); i++ {
		if m, ok := v.Index(i).Interface().(domain.ModelInterface); ok {
			m.Stamp()
			m.SetID(primitive.NewObjectID())
			is = append(is, m)
		}
	}
	if _, err := r.Collection.InsertMany(ctx, is); err != nil {
		return err
	}
	return nil
}

// Update update
func (r *Repo) Update(i interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	var id primitive.ObjectID
	if m, ok := i.(domain.ModelInterface); ok {
		m.UpdateStamp()
		id = m.GetID()
	}
	r.Mux.Lock()
	err := r.Collection.FindOneAndReplace(ctx, primitive.M{"_id": id}, i).Err()
	r.Mux.Unlock()
	if err != nil {
		return err
	}
	return nil
}

// HardDelete delete
func (r *Repo) HardDelete(i interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	var id primitive.ObjectID
	if m, ok := i.(domain.ModelInterface); ok {
		id = m.GetID()
	}
	r.Mux.Lock()
	_, err := r.Collection.DeleteOne(ctx,
		primitive.D{
			primitive.E{
				Key:   "_id",
				Value: id,
			},
		})
	r.Mux.Unlock()
	if err != nil {
		return err
	}
	return nil
}
