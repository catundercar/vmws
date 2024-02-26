package vm

import (
	"context"
	"fmt"
	vm_v1 "github.com/catundercar/vmws-go/api/vm/v1"
	"github.com/catundercar/vmws-go/internal/pkg/vm/vmrest"
	"github.com/catundercar/vmws-go/internal/pkg/vm/vmrun"
	"log"
	"sort"
	"sync"
	"time"
)

var VMNotFound = fmt.Errorf("not found")

type Management interface {
	List(ctx context.Context) ([]*vm_v1.VM, error)
	Get(ctx context.Context, name string) (*vm_v1.VM, error)
	Run(ctx context.Context, name string) error
}

type Config struct {
	UserName string
	Password string
}

type management struct {
	mux sync.RWMutex
	vms map[string]*vm_v1.VM

	vmrest *vmrest.VMRest
	vmrun  *vmrun.VMRun
}

func NewManagement(cfg Config) (Management, error) {
	m := &management{
		mux:   sync.RWMutex{},
		vms:   make(map[string]*vm_v1.VM),
		vmrun: &vmrun.VMRun{},
		vmrest: &vmrest.VMRest{
			UserName: cfg.UserName,
			Password: cfg.Password,
		},
	}
	return m, m.load()
}

func (m *management) load() error {
	m.mux.Lock()
	defer m.mux.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	vmList, err := m.vmrest.ListVMs(ctx)
	if err != nil {
		return err
	}
	for _, v := range vmList {
		displayName, err := GetDisplayNameFromPath(v.Path)
		if err != nil {
			return err
		}
		v.DisplayName = displayName
	}

	ps, err := m.vmrun.List(ctx)
	if err != nil {
		return err
	}
	runningMap := make(map[string]struct{})
	for _, p := range ps {
		runningMap[p] = struct{}{}
	}

	for _, vm := range vmList {
		if _, ok := runningMap[vm.Path]; ok {
			vm.Status = vm_v1.VM_Running
		}
		m.vms[vm.Name] = vm
	}
	return nil
}

// List all VMs.
func (m *management) List(_ context.Context) ([]*vm_v1.VM, error) {
	m.mux.Lock()
	defer m.mux.Unlock()

	vms := make([]*vm_v1.VM, 0, len(m.vms))
	for _, vm := range m.vms {
		vms = append(vms, vm)
	}
	sort.SliceStable(vms, func(i, j int) bool {
		return vms[i].DisplayName < vms[j].DisplayName
	})
	sort.SliceStable(vms, func(i, j int) bool {
		return vms[i].Status > vms[j].Status
	})
	return vms, nil
}

func (m *management) Get(ctx context.Context, name string) (*vm_v1.VM, error) {
	//TODO implement me
	panic("implement me")
}

func (m *management) Run(ctx context.Context, name string) error {
	m.mux.RLock()
	defer m.mux.RUnlock()

	vm, ok := m.vms[name]
	if !ok {
		return fmt.Errorf("vm:%s %w", name, VMNotFound)
	}
	if vm.Status == vm_v1.VM_Running {
		return nil
	}
	if err := m.vmrun.Start(ctx, vm.Path); err != nil {
		return fmt.Errorf("start vm: %s %w", vm.Name, err)
	}
	log.Println("success start ", vm.Path)
	return nil
}
