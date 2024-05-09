# BW-SETUP-SECRETS

This is an alternative low maintenance and highly reproducible solution to sercets management. I do not think that secrets should be part of your nix config file, neither I think that secrets should be public, however well encrypted.

The script sets up all your secrets from bitwarden, using the bitwarden-cli. The expectation is that you have a note in bitwarden, where the Note part of it is the key value env variable pairs, or anything else that your terminal can source. 

The note can also contain files, like ssh keys or smb credentials that you might need but they need to be stored secretly. These files can be added as attachement to the note.

## Development 

To enter the developmen env on nix you can use `direnv` as we have `.envrc` file prepared.

To generate the template we used: `nix flake init -t github:nix-community/gomod2nix#app`

Then in the project we can use `gomod2nix generate` to generate the shell dependencies.

To run the project we can simply do `nix run` to get the default package from the flake and run it.

## Usage

Note: If you are not using this script throuth the flake package provided, you need to have your system running the `bitwarden/cli` command line tool.

The application will try to read the config file from `~/.config/bw-setup-secrets/conf.toml`.

This file has to contain the following:
```toml
# The note id inside bitwarden where you store your secrets
noteId = "c660a85d-1081-464e-a37f-b396bca432f5"
# The server where your bitwarden is hosted
server = "https://my.custom.domain.com"
# The email address to log into bitwarden
email = "test@test.com"
# The file where the text part of the note will be copied
# This also needs to be sourced (for example in .zshrc)
secretsFile = ".zshsecrets"

# Array of files theat you need to copy
[[files]]
# The filename in bitwarden attachements (on the provided noteId).
srcFile = "config"
# The filename relative to HOME where to copy the file
destFile = ".ssh/config"
# The permissions in octa system for the file
chmod = "0600"

[[files]]
srcFile = "file2"
destFile = ".ssh/file2"
chmod = "0644"
```
