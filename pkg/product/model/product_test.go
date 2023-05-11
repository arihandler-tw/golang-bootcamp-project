package model

import (
	"errors"
	"reflect"
	"testing"
	"time"
)

func TestProduct(t *testing.T) {
	type args struct {
		id          string
		price       float32
		description string
		creation    time.Time
	}
	tests := []struct {
		name  string
		args  args
		want  *Product
		error error
	}{
		{
			name: "Valid product is created",
			args: args{
				id:          "id-1",
				price:       1.0,
				description: "valid description",
				creation:    time.Time{},
			},
			want: &Product{
				ID:          "id-1",
				Price:       1.0,
				Creation:    time.Time{},
				Description: "valid description",
			},
		},
		{
			name: "Product with description larger than 50 character returns an error",
			args: args{
				id:          "id-1",
				price:       1.0,
				description: "this description is a little more than 50 character long",
				creation:    time.Time{},
			},
			want:  nil,
			error: errors.New("description should be less than 50 characters long"),
		},
		{
			name: "Product with non-ASCII characters present its ID returns an error",
			args: args{
				id:          "invalid-id-千",
				price:       1.0,
				description: "valid description",
				creation:    time.Time{},
			},
			want:  nil,
			error: errors.New("ID should contain only ASCII characters"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewProduct(tt.args.id, tt.args.price, tt.args.description, tt.args.creation)
			if err != nil && !reflect.DeepEqual(err, tt.error) {
				t.Errorf("NewProduct got unexpected error: %v", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewProduct got %v, wanted %v", got, tt.want)
			}
		})
	}
}
