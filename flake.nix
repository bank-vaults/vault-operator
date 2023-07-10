{
  description = "Kubernetes Operator for Hashicorp Vault";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    flake-parts.url = "github:hercules-ci/flake-parts";
    devenv.url = "github:cachix/devenv";
  };

  outputs = inputs@{ flake-parts, ... }:
    flake-parts.lib.mkFlake { inherit inputs; } {
      imports = [
        inputs.devenv.flakeModule
      ];

      systems = [ "x86_64-linux" "x86_64-darwin" "aarch64-darwin" ];

      perSystem = { config, self', inputs', pkgs, system, ... }: rec {
        devenv.shells = {
          default = {
            languages = {
              go.enable = true;
            };

            services = {
              vault.enable = true;
            };

            pre-commit.hooks = {
              nixpkgs-fmt.enable = true;
              yamllint.enable = true;
              hadolint.enable = true;
            };

            packages = with pkgs; [
              gnumake

              kubectl
              kubectl-images

              yamllint
              hadolint
            ] ++ [
              self'.packages.kurun
              self'.packages.envtpl
              self'.packages.cidr
            ];

            scripts = {
              versions.exec = ''
                go version
                kubectl version --client
              '';
            };

            enterShell = ''
              versions
            '';

            # https://github.com/cachix/devenv/issues/528#issuecomment-1556108767
            containers = pkgs.lib.mkForce { };
          };

          ci = devenv.shells.default;
        };

        packages = {
          # TODO: create flake in source repo
          kurun = pkgs.buildGoModule rec {
            pname = "kurun";
            version = "0.7.0";

            src = pkgs.fetchFromGitHub {
              owner = "banzaicloud";
              repo = "kurun";
              rev = "${version}";
              sha256 = "sha256-b7ucOpTv+JON1yYxb1OhxBTZhyppKssOP7GNkmaCI5s=";
            };

            vendorSha256 = "sha256-kbdYDzPSNU3s4E4OwEGG9nbg66EwX18t+SVB4GejsNA=";

            subPackages = [ "." ];

            ldflags = [
              "-w"
              "-s"
              "-X main.version=v${version}"
            ];
          };

          envtpl = pkgs.buildGoModule rec {
            pname = "envtpl";
            version = "428c2d7";

            src = pkgs.fetchFromGitHub {
              owner = "subfuzion";
              repo = "envtpl";
              rev = "428c2d7";
              sha256 = "sha256-w1HaBB7M+yQyslFk+hHHxkz9kcniKFkS7CbD6ABrgU8=";
            };

            vendorSha256 = null;

            subPackages = [ "cmd/envtpl" ];
          };

          cidr = pkgs.buildGoPackage rec {
            pname = "cidr";
            version = "9c69a7cbc86a584f29cb8492b245e17b0267237d";

            goPackagePath = "github.com/hankjacobs/cidr";

            src = pkgs.fetchFromGitHub {
              owner = "hankjacobs";
              repo = "cidr";
              rev = "9c69a7cbc86a584f29cb8492b245e17b0267237d";
              sha256 = "sha256-kdPTGjXcna/Khdcvn+IWjoqCeWoQnYvXdEAy0bqKb24=";
            };

            subPackages = [ "." ];
          };
        };
      };
    };
}
