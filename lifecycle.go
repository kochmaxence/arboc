package arboc

func (cmd *Cmd) CallPersistentPreRun(cmdCtx CmdCtx) error {
	if cmd.PersistentPreRun == nil {
		return nil
	}

	return cmd.PersistentPreRun(cmdCtx.Cmd.cobraCmd, cmdCtx.Args)
}

func (cmd *Cmd) CallPreRun(cmdCtx CmdCtx) error {
	if cmd.PreRun == nil {
		return nil
	}

	return cmd.PreRun(cmdCtx.Cmd.cobraCmd, cmdCtx.Args)
}

func (cmd *Cmd) CallRun(cmdCtx CmdCtx) error {
	if cmd.Run == nil {
		return nil
	}

	return cmd.Run(cmdCtx.Cmd.cobraCmd, cmdCtx.Args)
}

func (cmd *Cmd) CallPersistentPostRun(cmdCtx CmdCtx) error {
	if cmd.PersistentPostRun == nil {
		return nil
	}

	return cmd.PersistentPostRun(cmdCtx.Cmd.cobraCmd, cmdCtx.Args)
}

func (cmd *Cmd) CallPostRun(cmdCtx CmdCtx) error {
	if cmd.PostRun == nil {
		return nil
	}

	return cmd.PostRun(cmdCtx.Cmd.cobraCmd, cmdCtx.Args)
}
