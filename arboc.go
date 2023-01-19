package arboc

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type Cmd struct {
	cobraCmd *cobra.Command

	parent      *Cmd
	built       bool
	initialized bool

	Use     string
	Aliases []string
	Short   string
	Long    string

	Args cobra.PositionalArgs

	SubCmds []*Cmd

	Configs []interface{}

	PersistentPreRun  CobraHandler
	PreRun            CobraHandler
	Run               CobraHandler
	PersistentPostRun CobraHandler
	PostRun           CobraHandler

	ChainPersistentPreRun  bool
	ChainPersistentPostRun bool
}

func NewCmd(use, short string) *Cmd {
	logrus.Debugf("#NewCmd(%s, %s)", use, short)
	cmd := &Cmd{Use: use, Short: short}

	return cmd
}

func (cmd *Cmd) CobraCmd() *cobra.Command {
	return cmd.cobraCmd
}

func (cmd *Cmd) Built() bool {
	return cmd.built
}

func (cmd *Cmd) Initialized() bool {
	return cmd.initialized
}

func (cmd *Cmd) SetUse(use string) *Cmd {
	cmd.Use = use
	return cmd
}

func (cmd *Cmd) SetAliases(aliases []string) *Cmd {
	cmd.Aliases = aliases
	return cmd
}

func (cmd *Cmd) SetShort(short string) *Cmd {
	cmd.Short = short
	return cmd
}

func (cmd *Cmd) SetLong(long string) *Cmd {
	cmd.Long = long
	return cmd
}

func (cmd *Cmd) SetPositionalArgs(args cobra.PositionalArgs) *Cmd {
	cmd.Args = args
	return cmd
}

func (cmd *Cmd) SetParent(parent *Cmd) *Cmd {
	cmd.parent = parent

	return cmd
}

func (cmd *Cmd) AddSubCmd(subCmd *Cmd) *Cmd {
	if !subCmd.Built() {
		subCmd.Build()
	}

	cmd.SubCmds = append(cmd.SubCmds, subCmd)

	subCmd.SetParent(cmd)

	return cmd
}

func (cmd *Cmd) AddSubCmds(subCmds ...*Cmd) *Cmd {
	for _, subCmd := range subCmds {
		cmd.AddSubCmd(subCmd)
	}

	return cmd
}

func (cmd *Cmd) SetFlagsConfig(ptrs ...interface{}) *Cmd {
	if cmd.Configs == nil {
		cmd.Configs = []interface{}{}
	}

	cmd.Configs = append(cmd.Configs, ptrs...)

	return cmd
}

func (cmd *Cmd) Init() *Cmd {
	if cmd.Configs != nil {
		for _, cfg := range cmd.Configs {
			GenerateFlags(cmd.CobraCmd(), cfg)
		}
	}

	cmd.initialized = true

	return cmd
}

func (cmd *Cmd) InitChain() *Cmd {
	if !cmd.Initialized() {
		cmd.Init()
	}

	for _, subCmd := range cmd.SubCmds {
		subCmd.InitChain()
	}

	return cmd
}

func (cmd *Cmd) Build() (*Cmd, error) {
	cmd.cobraCmd = &cobra.Command{
		Use:                cmd.Use,
		Aliases:            cmd.Aliases,
		Short:              cmd.Short,
		Long:               cmd.Long,
		Args:               cmd.Args,
		PersistentPreRunE:  cmd.PersistentPreRun,
		RunE:               cmd.Run,
		PostRunE:           cmd.PostRun,
		PersistentPostRunE: cmd.PersistentPostRun,
	}

	cmd.built = true

	return cmd, nil
}

func (cmd *Cmd) BuildChain() *Cmd {
	if !cmd.Built() {
		cmd.Build()
	}

	for _, subCmd := range cmd.SubCmds {
		if !subCmd.Built() {
			subCmd.Build()
		}

		cmd.CobraCmd().AddCommand(subCmd.CobraCmd())
		subCmd.BuildChain()
	}

	return cmd
}

func (cmd *Cmd) Execute() error {
	if !cmd.Built() {
		cmd.Build()
	}

	if !cmd.Initialized() {
		cmd.Init()
	}

	return cmd.CobraCmd().Execute()
}
