package vm

import (
	"context"
	"fmt"
	"testing"
)

func Test_management_Get(t *testing.T) {

}

func Test_management_List(t *testing.T) {
	m, err := NewManagement(Config{
		UserName: "admin",
		Password: "A123456!a",
	})
	if err != nil {
		t.Fatal(err)
	}
	vms, err := m.List(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	for _, vm := range vms {
		fmt.Println(vm.DisplayName, vm.Status, vm.Path)
	}
}

func Test_management_Run(t *testing.T) {

}
