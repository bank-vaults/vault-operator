{
  description = "Kubernetes Operator for Hashicorp Vault";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs { inherit system; };
      in
      {
        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            git
            gnumake

            go_1_20
            golangci-lint

            kind
            kubectl
            kubectl-images
            kustomize
            kubernetes-helm

            buildah

            vault
          ];
        };
      }
    );
}
