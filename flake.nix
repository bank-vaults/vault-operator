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

      perSystem = { config, self', inputs', pkgs, system, ... }: {
        devenv.shells.default = {
          languages = {
            go.enable = true;
          };

          services = {
            vault.enable = true;
          };

          packages = with pkgs; [
            gnumake

            golangci-lint

            kind
            kubectl
            kubectl-images
            kustomize
            kubernetes-helm
            helm-docs

            buildah
          ];

          scripts = {
            versions.exec = ''
              go version
              golangci-lint version
              kind version
              kubectl version --client
              echo kustomize $(kustomize version --short)
              echo helm $(helm version --short)
            '';
          };

          enterShell = ''
            versions
          '';

          # https://github.com/cachix/devenv/issues/528#issuecomment-1556108767
          containers = pkgs.lib.mkForce { };
        };
      };
    };
}
