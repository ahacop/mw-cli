{
  description = "Merriam-Webster CLI dictionary tool";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in
      {
        packages.default = pkgs.buildGoModule {
          pname = "mw-cli";
          version = "0.1.0";
          src = ./.;

          vendorHash = "sha256-/Gsqc8rEptMBItqeb/N/gE4V3iUGZa8k1GqUR1+togY=";

          ldflags = [ "-s" "-w" ];

          meta = with pkgs.lib; {
            description = "Merriam-Webster CLI dictionary tool";
            homepage = "https://github.com/ahacop/mw-cli";
            license = licenses.gpl3;
            maintainers = [ ];
            mainProgram = "mw-cli";
          };
        };

        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            go
            gopls
            gotools
            go-tools
            golangci-lint
            just
          ];

          shellHook = ''
            echo "Merriam-Webster CLI Development Environment"
            echo "Go version: $(go version)"
            echo ""
            echo "Don't forget to set your DICTIONARY_KEY environment variable!"
            echo "  export DICTIONARY_KEY=your-api-key-here"
          '';
        };
      }
    );
}
