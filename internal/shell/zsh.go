package shell

func init() {
	// There is no reason Zsh and Bash won't share the same files
	registerShell("zsh", &Bash{})
}
