# Integral Prompt

<img src="https://raw.githubusercontent.com/Readf0x/integral-prompt/refs/heads/main/screenshots/indev_v0.2.1.png">

## About
This prompt theme was created mostly due to my dissatisfaction with existing ones. My main issues were:
- overcomplicated configuration
- feature bloat
- wrapping issues

I have successfully solved these issues (in my opinion).

## Usage
Load with your favorite plugin loader (only officially supports [antidote](https://antidote.sh/)), or source `init.zsh` in your `.zshrc`.
To configure, add a `~/.integralrc` file. It can also be placed at:
- `$XDG_CONFIG_HOME/integralrc`
- `$XDG_CONFIG_HOME/integral/rc`
- `$XDG_CONFIG_HOME/integral/rc.zsh`
- `~/.config/integralrc`
- `~/.config/integral/rc`
- `~/.config/integral/rc.zsh`

## Planned Features
- [ ] Configuration files
- [x] Transient Prompt
- [x] Multi-line prompt
- [x] Rerender on terminal resize
- [ ] Module loader
    - [x] Colors
    - [x] Formatting
    - [x] Icons
    - [x] Positions
    - [x] Order
    - [x] Seperators
    - [x] Dynamic triggers
    - [ ] Right prompt
- [ ] Plugin Support
- [ ] Documentation

### Planned Modules
- [x] Background Jobs
- [x] Battery
- [x] CPU
- [x] CWD
- [x] Clock
- [x] [Direnv](https://github.com/direnv/direnv)
- [x] Error Codes
- [ ] Git
- [x] Nix Shell
- [x] SSH
- [ ] Uptime
- [x] Vim
