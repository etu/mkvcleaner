{ pkgs ? (import <nixpkgs> {}) }:

pkgs.mkShell {
  buildInputs = with pkgs; [
    ffmpeg-full    # for ffprobe and ffmpeg
    jq             # for jq
    coreutils-full # for printf, cut, wc
    findutils      # for find
    gnugrep        # for grep
    gnused         # for sed
  ];
}
