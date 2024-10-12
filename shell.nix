{ pkgs ? import <nixpkgs> {} }:
  pkgs.mkShell {
    nativeBuildInputs = with pkgs.buildPackages; [ go_1_21 go-task ];
    hardeningDisable = [ "fortify" ];
}
