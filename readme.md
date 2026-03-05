# Integral Prompt

<img src="https://raw.githubusercontent.com/Readf0x/integral-prompt/refs/heads/main/screenshots/indev_v0.3.0.png">

## About
This prompt theme was created mostly due to my dissatisfaction with existing ones. My main issues were:
- overcomplicated configuration
- feature bloat
- wrapping issues

## Installation
<!--Load with your favorite plugin loader (only officially supports [antidote](https://antidote.sh/)), or source `init.zsh` in your `.zshrc`.-->
Grab [`integral.deb`](https://github.com/Readf0x/integral-prompt/releases/latest/download/integral.deb) for Debian based distros, or [`integral.tar.gz`](https://github.com/Readf0x/integral-prompt/releases/latest/download/integral.tar.gz) for others.

### Flake install (Home Manager)
Add to your inputs and add `integral-prompt.homeManagerModules.default` to your home manager imports
```nix
# flake.nix
{
  inputs = {
    integral-prompt.url = "github:readf0x/integral-prompt";
  };
}
```

```nix
# home.nix
{ inputs, ... }: {
  imports = [ inputs.integral-prompt.homeManagerModules.default ];

  programs.integral-prompt.enable = true;
}
```

### Manual build
Install golang and run `./build.sh` to generate a tarball for generic linux.
For Debian based distros, run `./build.sh deb`.

## Usage

### Zsh
Add the following to your `.zshrc`
```sh
source <(integral init zsh)
```

### Bash
Add the following to your `.bashrc`
```sh
eval "$(integral init bash)"
```

### Configuration
To configure, add a `.integralrc` file.
```sh
integral config > ~/.integralrc
```

> [!NOTE]
> It can also be placed at:
> - `$XDG_CONFIG_HOME/integralrc`
> - `$XDG_CONFIG_HOME/integralrc.json`
> - `$XDG_CONFIG_HOME/integral/rc`
> - `$XDG_CONFIG_HOME/integral/rc.json`
> If `$XDG_CONFIG_HOME` is undefined, it will fall back to `~/.config`

The configuration options aren't yet documented, but if you have a JSON language server simply add
```json
"$schema": "/usr/share/integral/schema.json"
```
to the top of your configuration, and the LSP can list all available options. I know that's not ideal, but I
haven't added [jsonschema](https://github.com/invopop/jsonschema) description fields to the config types. Once that's done, I should be able to generate
actual documentation as well.

#### Home Manager Configuration
```nix
programs.integral-prompt = {
  enable = true;
  # enable shell integration here or with 'home.shell.enable<Shell>Integration'
  # for zsh
  enableZshIntegration = true;
  # for bash
  enableBashIntegration = true;
  config = {
    # JSON config
  };
}
```

## Planned Features
- [x] Configuration files
- [x] Transient Prompt
- [x] Multi-line prompt
- [x] Rerender on terminal resize
- [x] Module loader
    - [x] Colors
    - [x] Formatting
    - [x] Icons
    - [x] Positions
    - [x] Order
    - [x] Seperators
    - [x] Dynamic triggers
    - [x] Right prompt
- [ ] Plugin Support
- [ ] Documentation
- [ ] Module timeouts

### Planned Modules
- [x] Background Jobs
- [x] Battery
- [x] CPU
- [x] CWD
- [x] Clock
- [x] [Direnv](https://github.com/direnv/direnv)
- [x] Error Codes
- [x] Git
- [x] Nix Shell
- [x] SSH
- [x] Uptime
- [x] Vim
