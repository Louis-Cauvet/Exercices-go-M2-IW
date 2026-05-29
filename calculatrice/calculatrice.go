package main

import "fmt"

func operer(a, b float64, op string) (float64, error) {
	// On gère l'erreur en cas de division par 0
	if op == "/" && b == 0 {
		return 0, fmt.Errorf("division par zéro")
	}

	// On effectue l'opération pour obtenir le résultat
	fn := creerOperation(op)

	// On gère l'erreur en cas d'opération inconnue
	if fn == nil {
		return 0, fmt.Errorf("opération inconnue")
	}

	// On renvoie le résultat et l'erreur
	return fn(a, b), nil
}

func creerOperation(op string) func(float64, float64) float64 {
	// On effectue le calcul selon le type d'opération renseigné
	switch op {
	case "+":
		return func(a, b float64) float64 { return a + b }
	case "-":
		return func(a, b float64) float64 { return a - b }
	case "*":
		return func(a, b float64) float64 { return a * b }
	case "/":
		return func(a, b float64) float64 { return a / b }
	default:
		return nil
	}
}

func main() {
	// On demande à l'utilisateur de rentrer son calcul
	fmt.Println("Saisissez : <nombre> <opération> (+, -, *, /) <nombre>")
	fmt.Println("Saisissez '0 quit 0' pour quitter.")

	// On parcourt chaque calcul renseigné par l'utilisateur
	for {
		var a, b float64
		var op string
		fmt.Print("Entrée : ")
		// On attend les données dans le format "<nombre> <opération> <nombre>"		fmt.Scan(&a, &op, &b)

		// Si l'utilisateur demande à quitter, on stoppe la boucle for
		if op == "quit" {
			fmt.Println("Au revoir !")
			break
		} else {
			// On procède au calcul
			resultat, err := operer(a, b, op)

			// Si une erreur est rencontrée, on l'affiche
			if err != nil {
				fmt.Println("Erreur :", err)
			} else {
				// Sinon, on affiche le résultat du calcul
				fmt.Printf("%.2f %s %.2f = %.2f\n", a, op, b, resultat)
			}
		}
	}
}
