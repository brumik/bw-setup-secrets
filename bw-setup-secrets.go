package main

import (
	"fmt"
  "os/exec"
  "os"
  "github.com/pelletier/go-toml/v2"
  "syscall"
  "golang.org/x/term"
  "io/fs"
  "strconv"
)

type ConfigFiles struct {
  SrcFile   string
  DestFile  string
  Chmod     string
}

type Config struct { 
  NoteId      string
  Server      string
  Email       string
  SecretsFile string
  Files       []ConfigFiles
}

type Globals struct {
  Toml      Config
  HomePath  string
  Password  string
  Session   string
}

func parseToml(homePath string) Config {
	data, err := os.ReadFile(homePath + ".config/bw-setup-secrets/conf.toml")
	if err != nil {
		panic(err)
	}

  var cfg Config
  err2 := toml.Unmarshal(data, &cfg)
  if err2 != nil {
    panic(err2)
  }

  return cfg;
}

func handleCmd(cmd *exec.Cmd) string {
  stdout, err := cmd.CombinedOutput()
  if err != nil {
    fmt.Println("Error happened")
    fmt.Println(string(stdout))
    fmt.Println(err.Error())
    panic(err.Error())
  }
  return string(stdout)
}

func readPassword() string {
  fmt.Println("Enter your password to unlock the vault")
  passwd, err := term.ReadPassword(int(syscall.Stdin))
  if err != nil {
    fmt.Println("Error: could not read the password")  
    panic(err)
  }
  return string(passwd)
}

func writeToFile(name string, content []byte, perms fs.FileMode) {
  if err := os.WriteFile(name, content, perms); err != nil {
    panic(err)
  }
}

func createSecretsFile(cfg Globals) {
  content, err := exec.Command(
    "bw", "get", "notes", cfg.Toml.NoteId, "--session", cfg.Session,
  ).Output()
  if err != nil {
    fmt.Println("Could not write the file", cfg.Toml.SecretsFile)
    panic(err)
  }

  writeToFile(cfg.HomePath + cfg.Toml.SecretsFile, content, 0644)
  fmt.Println(
    "Updated file:",
    cfg.HomePath + cfg.Toml.SecretsFile,
    "Reload your terminal to take effect",
  )
}

func main() {
  var cfg Globals
  cfg.HomePath = os.Getenv("HOME") + "/"
  cfg.Toml = parseToml(cfg.HomePath)
  cfg.Password = readPassword()

  var cmd *exec.Cmd
  cmd = exec.Command("bw", "config", "server", cfg.Toml.Server)
  fmt.Println(handleCmd(cmd))

  out, err := exec.Command("bw", "login", "--check").CombinedOutput()
  if err == nil {
    fmt.Println(string(out))
    cmd = exec.Command("bw", "unlock", cfg.Password, "--raw")
    cfg.Session = handleCmd(cmd)
    if cfg.Session == "" {
      panic("Could not set session key from unlock")
    }
  } else {
    cmd = exec.Command("bw", "login", cfg.Toml.Email, cfg.Password, "--raw")
    cfg.Session = handleCmd(cmd)
    if cfg.Session == "" {
      panic("Could not set session key from login")
    }
  }

  cmd = exec.Command("bw", "sync", "--session", cfg.Session)
  fmt.Println(handleCmd(cmd))

  createSecretsFile(cfg)

  for _, value := range cfg.Toml.Files {
    fileMods, _ := strconv.ParseUint(value.Chmod, 8, 32)

    // TODO: Provede the --output flag with paht and set the chmod later
    // This will be hopefully faster
    content, err := exec.Command(
      "bw", "get", "attachment", value.SrcFile,
      "--itemid", cfg.Toml.NoteId,
      "--raw", "--session", cfg.Session,
    ).Output()
    if err != nil {
      fmt.Println("Could not write the file", value.DestFile)
      panic(err)
    }

    writeToFile(cfg.HomePath + value.DestFile, content, os.FileMode(fileMods))
    fmt.Println("Written " + cfg.HomePath + value.DestFile)
  }

  cmd = exec.Command("bw", "lock", "--session", cfg.Session)
  fmt.Println(handleCmd(cmd))
}
