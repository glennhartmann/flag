{
  inputs = {
    nixpkgs.url = github:NixOS/nixpkgs;
    flake-compat.url = "https://flakehub.com/f/edolstra/flake-compat/1.tar.gz";
    flake-utils.url = "github:numtide/flake-utils";
  };
  outputs = { self, nixpkgs, flake-compat, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs { inherit system; };
        flag = pkgs.buildGoModule {
          pname = "flag";
          version = "v0.6.0";
          src = builtins.path { path = ./.; name = "flag"; };
          vendorHash = "sha256-3PnXB8AfZtgmYEPJuh0fwvG38dtngoS/lxyx3H+rvFs=";
        };
      in
      {
        packages = {
          inherit flag;
          default = flag;
        };
      }
    );
}
