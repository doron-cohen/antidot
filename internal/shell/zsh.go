package shell

func init() {
	// There is no reason Zsh and Bash won't share the same files
	registerShell("zsh", &Bash{
		envMapFormat: NewKeyValueMapFormat(
			`^export (?P<key>\w+)="(?P<value>.*)"`,
			"export %s=\"%s\"\n",
		),
		aliasMapFormat: NewKeyValueMapFormat(
			`^alias (?P<alias>\w+)="(?P<command>.*)"`,
			"alias %s=\"%s\"\n",
		),
	})
}
