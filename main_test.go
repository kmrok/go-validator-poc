package main

import (
	"go-validator-poc/validator1"
	"go-validator-poc/validator2"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	log.SetOutput(ioutil.Discard)
	os.Exit(m.Run())
}

func Benchmark_Validator1(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		validator1.Validate()
	}
}

func Benchmark_Validator2(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		validator2.Validate()
	}
}
