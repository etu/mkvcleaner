{ pkgs ? (import <nixpkgs> {}) }:

pkgs.mkShell {
  buildInputs = with pkgs; [
    mkvtoolnix-cli # for mkvmerge
    coreutils-full # for printf, cut, wc
    findutils      # for find
    gnugrep        # for grep
    gnused         # for sed
  ];
}
