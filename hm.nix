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
    config = lib.mkOption {
      type = lib.types.attr;
      default = {};
      defaultText = "{}";
      description = "JSON attribute set";
    };
  };

  config = lib.mkIf cfg.enable {
    home.packages = lib.mkIf (cfg.package != null) [ cfg.package ];

    programs.zsh.initContent = (
      lib.mkOrder 600 ''
        eval "$(${lib.getExe cfg.package} init zsh)"
      ''
    );

    home.file.".config/integralrc" = builtins.toJSON {
      "$schema" = "${self.packages.${pkgs.system}.default}/share/integral/schema.json";
    } // cfg.config;
  };
}
