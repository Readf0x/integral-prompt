self: { config, lib, ... }: let
  cfg = config.programs.zsh.integral-prompt;
in {
  options.programs.zsh.integral-prompt = {
    enable = lib.mkEnableOption "integral prompt";
    package = lib.mkPackageOption self.packages "integral";
    enableZshIntegration = lib.hm.shell.mkZshIntegrationOption { inherit config; };
  };

  config = lib.mkIf cfg.enable {
    programs.zsh.initContent = lib.mkIf cfg.enableZshIntegration (
      lib.mkOrder 600 ''
        eval "$(${lib.getExe cfg.package} init zsh)"
      ''
    );
  };
}
