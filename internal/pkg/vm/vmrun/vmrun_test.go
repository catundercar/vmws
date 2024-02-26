package vmrun

import (
	"context"
	"fmt"
	"os"
	"path"
	"testing"
)

func TestVMRun_List(t *testing.T) {
	vr := &VMRun{}
	paths, err := vr.List(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	for _, p := range paths {
		fmt.Println(p)
	}
}

func TestVMRun_Start(t *testing.T) {
	vr := &VMRun{}
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatal(err)
	}
	err = vr.Start(context.Background(), path.Join(home, "vmware/worker2/worker2.vmx"))
	if err != nil {
		// TODO: asset
		t.Fatal(err)
	}
	TestVMRun_List(t)
}
