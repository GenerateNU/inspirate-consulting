package models

import (
	"bytes"
	"encoding/json"
	"reflect"

	"github.com/danielgtaylor/huma/v2"
)

// Nullable is a PATCH field that distinguishes three request states:
//   - omitted:     Set == false           → keep the stored value
//   - null:        Set == true, Null      → clear the stored value
//   - a value:     Set == true, !Null     → store Value
//
// A plain pointer cannot tell "omitted" from "null", which is why a nullable
// column could not be cleared through a COALESCE-based update.
type Nullable[T any] struct {
	Set   bool
	Null  bool
	Value T
}

// NewNullable returns a Nullable holding v.
func NewNullable[T any](v T) Nullable[T] { return Nullable[T]{Set: true, Value: v} }

// Null returns a Nullable that explicitly clears the field.
func Null[T any]() Nullable[T] { return Nullable[T]{Set: true, Null: true} }

// UnmarshalJSON is only called when the key is present in the body, so being
// called at all means the field was sent.
func (n *Nullable[T]) UnmarshalJSON(b []byte) error {
	n.Set = true
	if bytes.Equal(bytes.TrimSpace(b), []byte("null")) {
		n.Null = true
		return nil
	}
	return json.Unmarshal(b, &n.Value)
}

// Schema documents the field as the inner type, nullable. Field tags such as
// minimum, maximum and format still apply on top of it.
func (n Nullable[T]) Schema(r huma.Registry) *huma.Schema {
	s := *r.Schema(reflect.TypeOf(n.Value), true, "")
	s.Nullable = true
	return &s
}

// Ptr returns the value to write to the database: nil when cleared or omitted.
func (n Nullable[T]) Ptr() *T {
	if !n.Set || n.Null {
		return nil
	}
	v := n.Value
	return &v
}
