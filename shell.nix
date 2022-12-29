{ pkgs ? (import <nixpkgs> {}) }:

pkgs.mkShell {
  buildInputs = with pkgs; [
    ffmpeg-full # for ffprobe and ffmpeg
    go          # For building the project
  ];
}
