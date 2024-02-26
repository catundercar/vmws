package vmrest

import (
	"context"
	"encoding/json"
	"fmt"
	vm_v1 "github.com/catundercar/vmws-go/api/vm/v1"
	"io"
	"net/http"
	"path"
	"strings"
)

type VMRest struct {
	UserName string
	Password string
}

// ListVMs list all VMs.
// A Get HTTP Request to vmrest.
// GET /vms
func (vmrest *VMRest) ListVMs(ctx context.Context) ([]*vm_v1.VM, error) {
	// Note: only support 127.0.0.1.
	url := "http://127.0.0.1:8697" + "/api/vms"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(vmrest.UserName, vmrest.Password)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpext response status: %d %s", resp.StatusCode, resp.Status)
	}
	byt, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	vms := make([]vm, 0)
	if err = json.Unmarshal(byt, &vms); err != nil {
		return nil, err
	}

	vmList := make([]*vm_v1.VM, 0, len(vms))
	for _, vm := range vms {
		vmList = append(vmList, &vm_v1.VM{
			Name:        vm.ID,
			DisplayName: strings.TrimSuffix(path.Base(vm.Path), path.Ext(vm.Path)),
			Path:        vm.Path,
		})
	}
	return vmList, nil
}

type vm struct {
	ID   string `json:"id"`
	Path string `json:"path"`
}
