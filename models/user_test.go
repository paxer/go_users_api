package models

import (
	"reflect"
	"testing"
)

func TestNewUser(t *testing.T) {
	type args struct {
		name  string
		email string
	}
	tests := []struct {
		name    string
		args    args
		want    *User
		wantErr bool
	}{
		{name: "Email is empty", args: args{name: "Darth", email: ""}, want: nil, wantErr: true},
		{name: "Name is empty", args: args{name: "", email: "darth@vader.com"}, want: nil, wantErr: true},
		{name: "Name and Email not empty", args: args{name: "Darth Vader", email: "darth@vader.com"}, want: &User{Name: "Darth Vader", Email: "darth@vader.com"}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewUser(tt.args.name, tt.args.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUser() = %v, want %v", got, tt.want)
			}
		})
	}
}
