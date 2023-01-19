package arboc

import "github.com/spf13/cobra"

type CobraHandler func(cobraCmd *cobra.Command, args []string) error
type CmdHandler func(cmdCtx CmdCtx) error

func wrapHandler(cmd *Cmd, fn CmdHandler) CobraHandler {
	return func(cobraCmd *cobra.Command, args []string) error {
		return fn(CmdCtx{Cmd: cmd, Args: args})
	}
}

func (cmd *Cmd) OnPersistentPreRun(fn CmdHandler) *Cmd {
	cmd.PersistentPreRun = func(cobraCmd *cobra.Command, args []string) error {
		cmdCtx := CmdCtx{Cmd: cmd, Args: args}

		if cmd.ChainPersistentPreRun {
			if cmd.parent != nil {
				if err := cmd.parent.CallPersistentPreRun(cmdCtx); err != nil {
					return err
				}
			}
		}

		return fn(cmdCtx)
	}

	return cmd
}

func (cmd *Cmd) OnPreRun(fn CmdHandler) *Cmd {
	cmd.PreRun = wrapHandler(cmd, fn)

	return cmd
}

func (cmd *Cmd) OnRun(fn CmdHandler) *Cmd {
	cmd.Run = wrapHandler(cmd, fn)

	return cmd
}

func (cmd *Cmd) OnPersistentPostRun(fn CmdHandler) *Cmd {
	cmd.PersistentPostRun = func(cobraCmd *cobra.Command, args []string) error {
		cmdCtx := CmdCtx{Cmd: cmd, Args: args}

		if cmd.ChainPersistentPostRun {
			if cmd.parent != nil {
				if err := cmd.parent.CallPersistentPostRun(cmdCtx); err != nil {
					return err
				}
			}
		}

		return fn(cmdCtx)
	}

	return cmd
}

func (cmd *Cmd) OnPostRun(fn CmdHandler) *Cmd {
	cmd.PostRun = wrapHandler(cmd, fn)

	return cmd
}
