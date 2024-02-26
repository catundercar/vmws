package vm

import (
	"errors"
	"github.com/catundercar/vmws-go/internal/pkg/vm"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

func init() {
	VMCmd.AddCommand(
		list,
		start,
	)
}

var VMCmd = &cobra.Command{
	Use:   "vm",
	Short: "vms management",
}

var list = &cobra.Command{
	Use:     "list",
	Short:   "list vms",
	PreRunE: initManagement,
	RunE: func(cmd *cobra.Command, args []string) error {
		vms, err := m.List(cmd.Context())
		if err != nil {
			return err
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Name", "DisplayName", "Status", "Path"})
		for _, v := range vms {
			table.Append([]string{v.Name, v.DisplayName, v.Status.String(), v.Path})
		}
		table.Render()
		return nil
	},
}

var start = &cobra.Command{
	Use:     "start [Name]",
	Short:   "start a vm by name",
	PreRunE: initManagement,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires at least one arg")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		err := m.Run(cmd.Context(), args[0])
		if err != nil {
			return err
		}
		return nil
	},
}

var m vm.Management

func initManagement(cmd *cobra.Command, args []string) (err error) {
	m, err = vm.NewManagement(vm.Config{
		UserName: viper.GetString("username"),
		Password: viper.GetString("password"),
	})
	return
}
