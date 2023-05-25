package service

import (
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	type args struct {
		r repository
	}
	tests := []struct {
		name string
		args args
		want service
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.r); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_service_AllSessions(t *testing.T) {
	type fields struct {
		r repository
	}
	type args struct {
		date   string
		offset int
		limit  int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []entity.CinemaSession
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := service{
				r: tt.fields.r,
			}
			got, err := s.AllSessions(tt.args.date, tt.args.offset, tt.args.limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("AllSessions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AllSessions() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_service_SessionsForHall(t *testing.T) {
	type fields struct {
		r repository
	}
	type args struct {
		hallId int
		date   string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []entity.CinemaSession
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := service{
				r: tt.fields.r,
			}
			got, err := s.SessionsForHall(tt.args.hallId, tt.args.date)
			if (err != nil) != tt.wantErr {
				t.Errorf("SessionsForHall() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SessionsForHall() got = %v, want %v", got, tt.want)
			}
		})
	}
}
