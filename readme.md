# Integral Prompt

<img src="https://raw.githubusercontent.com/Readf0x/integral-prompt/refs/heads/main/screenshots/indev_v0.2.1.png" width="100%">

**Potential users be warned!** This plugin is very new, and has yet to reach feature maturity. Issues are expected, but should hopefully be minor.

## About
This prompt theme was created mostly due to my dissatisfaction with existing ones. My main issues were:
- overcomplicated configuration
- feature bloat
- wrapping issues

I have successfully solved these issues (in my opinion).

## Usage
Load with your favorite plugin loader (only officially supports [antidote](https://github.com/zsh-users/antidote)), or source `init.zsh` in your `.zshrc`.
To configure, add a `~/.integralrc` file. It can also be placed at:
- `$XDG_CONFIG_HOME/integralrc`
- `$XDG_CONFIG_HOME/integral/rc`
- `$XDG_CONFIG_HOME/integral/rc.zsh`
- `~/.config/integralrc`
- `~/.config/integral/rc`
- `~/.config/integral/rc.zsh`

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
    - [ ] Seperators
    - [x] Dynamic triggers
    - [ ] Right prompt
- [x] Plugin Support
- [ ] Documentation

### Planned Modules
- [x] Background Jobs
- [x] Battery
- [ ] CPU
- [x] CWD
- [x] Clock
- [x] Error Codes
- [x] Git
- [x] Nix Shell
- [x] SSH
- [x] Uptime
- [x] Vim
