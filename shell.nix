{ pkgs ? import <nixpkgs> { } }:
with pkgs;
mkShell {
  buildInputs = [
    go
    gotools
    esptool
    protobuf
    protoc-gen-go
  ];

  shellHook = ''
    PATH=$PATH:$HOME/stuff/tinygo/bin
  '';
}
