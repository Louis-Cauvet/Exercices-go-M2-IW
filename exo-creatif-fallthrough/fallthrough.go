// Exo créatif — Évaluateur de super-héros
//
// Ce programme utilise 3 nouvelles notions du langage Go : iota, slice et fallthrough.
//
// Je suis parti sur l'idée d'une classification de super-héros (étant fan de My Hero Academia)
// Le principe : chaque héros possède un niveau d'expérience.
// Ainsi, selon ce niveau, on peut déduire les pouvoirs qu'il possède.

package main

import "fmt"

// Grâce au iota, on génère automatiquement des entiers croissants à partir de 0.
// Chaque constante suivante prend automatiquement la valeur précédente + 1.
const (
	Novice        = iota + 1 // 1
	Confirmé                 // 2
	Professionnel            // 3
	Legende                  // 4
)

// On définit le type "Héros" afin d'indiquer quelles propriétés il possède
type Heros struct {
	nom        string
	expérience int
}

func main() {
	// Avec le slice, on définit une collection dynamique d'éléments similaires (les niveaux d'exépriences des héros).
	// Celle-ci peut être amenée à évoluer.
	// On réutilise ici les constantes iota définies ci-dessus
	heros := []Heros{
		{"Novice", Novice},
		{"Confirmé", Confirmé},
		{"Professionnel", Professionnel},
		{"Légende", Legende},
	}

	fmt.Println("Voici la classification des pouvoirs de chaque héros selon son expérience :")

	// On utilise une boucle "for" pour parcourir le slice de héros.
	// "_" permet d'ignorer l'index, "h" désigne chaque itération de héros
	for _, h := range heros {
		fmt.Printf("=== %s (niveau %d) ===\n", h.nom, h.expérience)

		// Avec un fallthrough, si le niveau d'exépérience d'un héros est rencontré,
		// alors on lui attribut également tous les pouvoirs des niveau précédents
		// Par exemple, un héros Légende exécute tous les case pour accumuler des autres pouvoirs.
		switch h.expérience {
		case Legende:
			fmt.Println("- Immortalité")
			fallthrough
		case Professionnel:
			fmt.Println("- Rayon laser")
			fallthrough
		case Confirmé:
			fmt.Println("- Super force")
			fallthrough
		case Novice:
			fmt.Println("- Vol")
		}
		fmt.Println()
	}
}
