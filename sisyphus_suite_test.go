package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestSisyphus(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Sisyphus Suite")
}
