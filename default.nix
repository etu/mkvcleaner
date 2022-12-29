{ pkgs ? import <nixpkgs> {}, ... }:

pkgs.buildGoModule {
  pname = "mkvcleaner";
  version = "20221229";

  src = ./.;

  nativeBuildInputs = [ pkgs.makeWrapper ];

  postInstall = ''
    wrapProgram "$out/bin/mkvcleaner" \
      --prefix PATH : "${pkgs.lib.makeBinPath [ pkgs.ffmpeg ]}"
  '';

  vendorSha256 = "0x6n6ijwgsgzyjimc50mkcv99j8fixid6l3jvhwhgrgjbdfpcdlf";
}
