## mpagd_util completion zsh

Generate the autocompletion script for zsh

### Synopsis

Generate the autocompletion script for the zsh shell.

If shell completion is not already enabled in your environment you will need
to enable it.  You can execute the following once:

	echo "autoload -U compinit; compinit" >> ~/.zshrc

To load completions in your current shell session:

	source <(mpagd_util completion zsh)

To load completions for every new session, execute once:

#### Linux:

	mpagd_util completion zsh > "${fpath[1]}/_mpagd_util"

#### macOS:

	mpagd_util completion zsh > $(brew --prefix)/share/zsh/site-functions/_mpagd_util

You will need to start a new shell for this setup to take effect.


```
mpagd_util completion zsh [flags]
```

### Options

```
      --no-descriptions   disable completion descriptions
```

### Options inherited from parent commands

```
      --help   help for this command
```

### SEE ALSO

* [mpagd_util completion](mpagd_util_completion.md)	 - Generate the autocompletion script for the specified shell

###### Auto generated by spf13/cobra on 27-Apr-2025
