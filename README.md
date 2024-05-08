# BW-SETUP-SECRETS

To enter the developmen env on nix you can use `direnv` as we have `.envrc` file prepared.

To generate the template we used: `nix flake init -t github:nix-community/gomod2nix#app`

Then in the project we can use `gomod2nix generate` to generate the shell dependencies.

To run the project we can simply do `nix run` to get the default package from the flake and run it.
