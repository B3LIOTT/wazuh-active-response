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

var dumpPath = "C:/path/to/dumpit.exe"

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
	cmd := exec.Command("Get-NetAdapter", "|", "Disable-NetAdapter", "-Confirm:$false")
	_, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("erreur lors de la désactivation des interfaces réseau: %w", err)
	}
	return nil
}

// Memory dump
func generateFullMemoryDump() error {
	// Commande pour effectuer un dump mémoire
	// dumpit.exe /a /o memory.raw /q
	cmd := exec.Command(dumpPath, "/a")
	_, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("erreur lors de la génération du dump mémoire: %w", err)
	}
	return nil
}

// TODO
func Delete() error {
	return nil
}
