self: { config, lib, pkgs, ... }: let
  cfg = config.programs.zsh.integral-prompt;
in {
  options.programs.zsh.integral-prompt = {
    enable = lib.mkEnableOption "integral prompt";
    package = lib.mkOption {
      type = lib.types.package;
      default = self.packages.${pkgs.system}.default;
      defaultText = "integral-prompt.packages.\${system}.default";
      description = "The package used for integral-prompt";
    };
    enableZshIntegration = lib.hm.shell.mkZshIntegrationOption { inherit config; };
  };

  config = lib.mkIf cfg.enable {
    home.packages = lib.mkIf (cfg.package != null) [ cfg.package ];

    programs.zsh.initContent = lib.mkIf cfg.enableZshIntegration (
      lib.mkOrder 600 ''
        eval "$(${lib.getExe cfg.package} init zsh)"
      ''
    );
  };
}
