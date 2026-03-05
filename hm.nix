self: { config, lib, pkgs, ... }: let
  cfg = config.programs.integral-prompt;
in {
  options.programs.integral-prompt = {
    enable = lib.mkEnableOption "integral prompt";
    enableZshIntegration = lib.hm.shell.mkZshIntegrationOption { inherit config; };
    enableBashIntegration = lib.hm.shell.mkBashIntegrationOption { inherit config; };
    package = lib.mkOption {
      type = lib.types.package;
      default = self.packages.${pkgs.system}.default;
      defaultText = "integral-prompt.packages.\${system}.default";
      description = "The package used for integral-prompt";
    };
    config = lib.mkOption {
      type = lib.types.attrs;
      default = {};
      defaultText = "{}";
      description = "JSON attribute set";
    };
    configPath = lib.mkOption {
      type = lib.types.enum [
        ".integralrc"
        ".config/integralrc"
        ".config/integralrc.json"
        ".config/integral/rc"
        ".config/integral/rc.json"
      ];
      default = ".config/integralrc.json";
      defaultText = ".config/integralrc.json";
      description = "Path to place your configuration";
    };
  };

  config = lib.mkIf cfg.enable {
    home.packages = lib.mkIf (cfg.package != null) [ cfg.package ];

    programs.zsh.initContent = lib.mkIf cfg.enableZshIntegration (
      lib.mkOrder 600 ''
        eval "$(${lib.getExe cfg.package} init zsh)"
      ''
    );

    programs.bash.initExtra = lib.mkIf cfg.enableBashIntegration ''
      eval "$(${lib.getExe cfg.package} init bash)"
    '';

    home.file.${cfg.configPath}.text = builtins.toJSON ({
      "$schema" = "${self.packages.${pkgs.system}.default}/share/integral/schema.json";
    } // cfg.config);
  };
}
