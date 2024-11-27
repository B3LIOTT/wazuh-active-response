package main

import (
	"bufio"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
)

func main() {
	// Lire les données d'entrée via stdin en toute sécurité
	scanner := bufio.NewScanner(os.Stdin)
	var alertData string

	fmt.Println("Reading alert data...")
	for scanner.Scan() {
		alertData += scanner.Text() + "\n"
	}

	// Vérifier les erreurs lors de la lecture de stdin
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading stdin: %v\n", err)
		return
	}

	// Obtenir le chemin du bureau de l'utilisateur actuel
	usr, err := user.Current()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting current user: %v\n", err)
		return
	}

	// Construire le chemin absolu vers le fichier test.txt
	desktopPath := filepath.Join(usr.HomeDir, "", "test.txt")

	// Créer/écrire dans le fichier sur le bureau
	file, err := os.Create(desktopPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating file: %v\n", err)
		return
	}
	defer file.Close()

	// Écrire les données d'alerte dans le fichier
	_, err = file.WriteString(alertData)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing to file: %v\n", err)
		return
	}

	fmt.Println("File test.txt created successfully on Desktop.")
}
