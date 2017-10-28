package main

import (
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
)

func Test_main(t *testing.T) {
	tests := []struct {
		name string
	}{
	// TODO: Add test cases.
	}
	for range tests {
		t.Run(tt.name, func(t *testing.T) {
			main()
		})
	}
}

func Test_startPage(t *testing.T) {
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name string
		args args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			startPage(tt.args.c)
		})
	}
}

func Test_apiCepabertoJSON(t *testing.T) {
	type args struct {
		resp *CepAbertoResult
	}
	tests := []struct {
		name string
		args args
		want string
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := apiCepabertoJSON(tt.args.resp); got != tt.want {
				t.Errorf("apiCepabertoJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_apiCep(t *testing.T) {
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name string
		args args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiCep(tt.args.c)
		})
	}
}

func Test_apiViacepJSON(t *testing.T) {
	type args struct {
		resp *ViaCepResult
	}
	tests := []struct {
		name string
		args args
		want string
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := apiViacepJSON(tt.args.resp); got != tt.want {
				t.Errorf("apiViacepJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_apiWriteJSON(t *testing.T) {
	type args struct {
		resp *brcepResult
	}
	tests := []struct {
		name string
		args args
		want string
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := apiWriteJSON(tt.args.resp); got != tt.want {
				t.Errorf("apiWriteJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_viacep(t *testing.T) {
	type args struct {
		cep string
	}
	tests := []struct {
		name string
		args args
		want *ViaCepResult
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := viacep(tt.args.cep); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("viacep() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_cepaberto(t *testing.T) {
	type args struct {
		cep string
	}
	tests := []struct {
		name string
		args args
		want *CepAbertoResult
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cepaberto(tt.args.cep); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("cepaberto() = %v, want %v", got, tt.want)
			}
		})
	}
}
