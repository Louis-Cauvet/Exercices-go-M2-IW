package main

import "fmt"

// On déclare les constantes en dehors de la fonction
const (
	Nom         = "Louis"
	IMCMaigreur = 18.5
	IMCNormal   = 25.0
	IMCSurpoids = 30.0
)

func main() {
	// On renseigne des données statiques pour le poids et la taille (j'ai mis mes données réelles pour vous prouver le monstre de puissance que je suis)
	poids := 69.0
	taille := 1.81

	// On calcule et on affiche l'IMC
	imc := poids / (taille * taille)
	fmt.Printf("Bonjour %s, votre IMC est : %.2f\n", Nom, imc)

	// On fait matcher la valeur obtenue avec la catégorie correspondante
	var categorie string
	switch {
	case imc < IMCMaigreur:
		categorie = "Maigreur"
	case imc < IMCNormal:
		categorie = "Normal"
	case imc < IMCSurpoids:
		categorie = "Surpoids"
	default:
		categorie = "Obésité"
	}

	fmt.Printf("La catégorie correspondante est : %s\n", categorie)
}
