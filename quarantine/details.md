# Quarantine script

This scripts isolate a potentially infected computer and dump the memory with `dumpit.exe` available [here](https://www.magnetforensics.com/resources/magnet-dumpit-for-windows/) for Windows and [here](https://github.com/MagnetForensics/dumpit-linux) for Linux. However, this script is not yet adapted for Linux.

The isolation process for Windows is just the execution of these commands:
```bash
netsh advfirewall set allprofiles firewallpolicy blockinbound,blockoutbound
Get-NetAdapter | Disable-NetAdapter -Confirm:$false
```

The first one applies a full deny rule on the firewall, and the second one disables every network adapter.

To go back at the initial state:
```bash
netsh advfirewall reset
Get-NetAdapter | Enable-NetAdapter -Confirm:$false
```

(NOTE: admin privileges are needed)