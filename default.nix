{ lib, pkgs ? (
    let
      inherit (builtins) fetchTree fromJSON readFile;
      inherit ((fromJSON (readFile ./flake.lock)).nodes) nixpkgs gomod2nix;
    in
    import (fetchTree nixpkgs.locked) {
      overlays = [
        (import "${fetchTree gomod2nix.locked}/overlay.nix")
      ];
    }
  )
, buildGoApplication ? pkgs.buildGoApplication
}:

buildGoApplication {
  pname = "bw-setup-secrets";
  version = "0.1";
  pwd = ./.;
  src = ./.;
  modules = ./gomod2nix.toml;
  nativeBuildInputs = [ 
    pkgs.makeWrapper
  ];

  postFixup = ''
    wrapProgram "$out/bin/bw-setup-secrets" --set PATH ${lib.makeBinPath [
      pkgs.bitwarden-cli
    ]}
  '';
}
