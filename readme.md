# Integral Prompt

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
    - [ ] Colors
    - [ ] Positions
    - [x] Order
    - [ ] Seperators
    - [ ] Dynamic triggers
- [x] Plugin Support
- [ ] Documentation

### Planned Modules
- [ ] Background Jobs
- [ ] Battery
- [ ] CPU
- [x] CWD
- [ ] Clock
- [x] Error Codes
- [x] Git
- [x] Nix Shell
- [ ] SSH
- [ ] Uptime
- [x] Vim
