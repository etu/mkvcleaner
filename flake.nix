{
  description = "etu/mkvcleaner";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, flake-utils, nixpkgs, ... }: flake-utils.lib.eachDefaultSystem (system: let
    pkgs = import nixpkgs { inherit system; };
  in {
    packages = flake-utils.lib.flattenTree {
      default = pkgs.buildGoModule (let
        version = "1.0.0.${nixpkgs.lib.substring 0 8 self.lastModifiedDate}.${self.shortRev or "dirty"}";
      in {
        pname = "mkvcleaner";
        inherit version;

        src = ./.;

        nativeBuildInputs = [ pkgs.makeWrapper ];

        postInstall = ''
          wrapProgram "$out/bin/mkvcleaner" \
            --prefix PATH : "${pkgs.lib.makeBinPath [ pkgs.ffmpeg ]}"
        '';

        vendorSha256 = "0x6n6ijwgsgzyjimc50mkcv99j8fixid6l3jvhwhgrgjbdfpcdlf";
      });
    };

    devShells = flake-utils.lib.flattenTree {
      default = { pkgs, ... }: pkgs.mkShell {
        buildInputs = [
          pkgs.ffmpeg  # for ffprobe and ffmpeg
          pkgs.gnumake # For the Makefile
          pkgs.go      # For building the project
        ];
      };
    };
  });
}
