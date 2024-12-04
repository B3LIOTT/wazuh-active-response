# Wazuh Active Response

Custom GoLang scripts for Wazuh Active Response.

## How to build active response binaries ?

First of all, main-to-modify is the go code which has to be modified in order to give him the active response script you want. It handles interactions between the Wazuh server and the Agent when a active response script is called.
To do so, you have to use the `makeFile` with the following parameters:
```bash
make build IMPORT=your/import OUTPUT=binary-name
```

Here is a example:
```bash
make build IMPORT=quarantine OUTPUT=quarantine_active_response
```

You can use `make clean` in order to delete the old `main.go`

## How to add custom scripts ?

If you want to add a custom active response script you have to respect the following structure in your go script:
```go
...

func Add() error {
    return error or nil
}

func Delete() error {
    return error or nil
}
```

You have to define the `Delete` function even if timeout is not allowed for your active response (just leave it empty)

## TODO
- Modify the makefile in order to allow Linux target compilation
- Create more active response scripts
- Improve the structure of custom scripts
