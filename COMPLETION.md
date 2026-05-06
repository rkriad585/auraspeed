# AuraSpeed Shell Completions

AuraSpeed uses [Cobra](https://github.com/spf13/cobra) which automatically provides a `completion` subcommand for generating shell completion scripts.

## Verify Completion Command

Cobra v1.9.1 (current version) includes the `completion` command automatically. Verify it works:
```bash
auraspeed completion --help
```

## Temporary Completion (Current Shell Session)

### Bash
```bash
source <(auraspeed completion bash)
```

### Zsh
```zsh
source <(auraspeed completion zsh)
```

### PowerShell
```powershell
auraspeed completion powershell | Out-String | Invoke-Expression
```

### Fish
```fish
auraspeed completion fish | source
```

## Permanent Completion (All Future Shell Sessions)

### Bash
```bash
# User-local (recommended)
mkdir -p ~/.local/share/bash-completion/completions/
auraspeed completion bash > ~/.local/share/bash-completion/completions/auraspeed

# Or system-wide (requires sudo)
sudo auraspeed completion bash > /etc/bash_completion.d/auraspeed
```

### Zsh
```zsh
# User-local (recommended)
mkdir -p ~/.zsh/completions
auraspeed completion zsh > ~/.zsh/completions/_auraspeed

# Add to your .zshrc if not already present:
# fpath=(~/.zsh/completions $fpath)
# autoload -U compinit && compinit
```

### PowerShell
```powershell
# Append to your PowerShell profile
auraspeed completion powershell | Out-File -Append $PROFILE

# Reload profile
. $PROFILE
```

### Fish
```fish
mkdir -p ~/.config/fish/completions
auraspeed completion fish > ~/.config/fish/completions/auraspeed.fish
```

Fish automatically loads completions from this directory on restart.
