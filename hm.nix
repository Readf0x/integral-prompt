self: { config, lib, ... }: let
  cfg = config.integral-prompt;
in {
  options.integral = {
    enable = lib.mkEnableOption "integral";
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
