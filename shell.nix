with (import <nixpkgs> {});

mkShell {
  buildInputs = [
    gnumake
    go
    gocode

    # Deps needed for goav
    pkgconfig
    ffmpeg-full
  ];
}
