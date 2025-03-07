package quanrantine

import (
	"fmt"
	"os/exec"
)

// =============================================
// Author: b3liott
// Date: 2024-12-04
// License: MIT
// Description: This active response puts in
//				quarantine the concerned computer
// 				It applies a full DENY rule on
//				the firewall, disable network
//				interfaces and dump the RAM.
// =============================================

var (
	dumperPath = `C:\Users\win10\Desktop\Comae-Toolkit-v20230117\x64\DumpIt.exe`
	dumpOutput = `C:\Users\win10\Desktop\mem_dump.dmp`
)

func Add(keys []interface{}) error {
	if err := blockAllTraffic(); err != nil {
		return err
	}
	if err := disableAllNetworkAdapters(); err != nil {
		return err
	}
	if err := generateFullMemoryDump(); err != nil {
		return err
	}
	return nil
}

// Full DENY applied on the windows firewall
func blockAllTraffic() error {
	cmd := exec.Command("netsh", "advfirewall", "set", "allprofiles", "firewallpolicy", "blockinbound,blockoutbound")
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("erreur lors du blocage du trafic réseau: %w", err)
	}
	return nil
}

// Disable every network interface
func disableAllNetworkAdapters() error {
	command := `Get-NetAdapter | Disable-NetAdapter -Confirm:$false`
	cmd := exec.Command("powershell", "-NoProfile", "-Command", command)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("erreur lors de la désactivation des interfaces réseau: %w", err)
	}
	return nil
}

// Memory dump
func generateFullMemoryDump() error {
	// Commande pour effectuer un dump mémoire
	cmd := exec.Command(dumperPath, "/QUIET", "/OUTPUT", dumpOutput)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("erreur lors de la génération du dump mémoire: %w", err)
	}
	return nil
}

// TODO
func Delete() error {
	return nil
}
