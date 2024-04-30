# System Configurator
An operating system agnostic tool installing packages, running setup scripts, and sharing these configurations across machines.

Currently supported package managers are: apk, apt, brew, dnf, snap, and pacman.

# Installation
1. Ensure `go` is installed on your system
> Installation steps are outlined [here](https://go.dev/doc/install).
2. Run `go install drew-english/system-configurator@latest`

# Usage
The system configurator CLI bases its configuration on a yaml file located at `$HOME/.config/system-configurator/config.yml`. The CLI has CRUD operations that will manage the configuration file (or the system directly) for the user, but they are by no means required for use. 

Additionally, the CLI has 3 different operation modes: configuration, system, and hybrid. Each mode will tell system configurator how it should perform the command given. See the help output for specific behavior of a given command.

Example configuration file:
```yaml
packages:
  - fzf
  - zoxide
  - name: stow
    version: 3.1.2
    alternates:
      apk: something-else
      apt:
        name: something-else
        version: 1.2.3
```

## Common Commands
`scfg help {command}` will show you a relevant description and help for the command or subcommand you are attempting to run.

### Package Management
#### Sync
`scfg package sync`

By default this will sync the packages from the configuration file to the system, installing only packages that are missing. See `scfg help package sync` for use with other modes.

#### Add
`scfg package add fzf`

By default this will add specified packages to the configuration. See `scfg help package add` for use with other modes.

#### Remove
`scfg package rm fzf`

By default this will remove the specified packages from the configuration. See `scfg help package rm` for use with other modes.
