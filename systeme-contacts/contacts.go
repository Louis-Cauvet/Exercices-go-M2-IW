package main

import "fmt"

// On définit la structure de base "Personne", qui est commune à tous les types de contacts
// On précise des tags `json:"..."` afin de sépcifier de la métadonnée pour les différentes propriétés.
// `omitempty` sur Email indique que le champ omis du JSON s'il est vide.
type Personne struct {
	Prenom string `json:"prenom"`
	Nom    string `json:"nom"`
	Age    int    `json:"age"`
	Email  string `json:"email,omitempty"`
}

// Cette méthode utilise un receiver par valeur,
// et donc fait une copie de la structure originale car la méthode ne modifie rien
func (p Personne) NomComplet() string {
	return p.Prenom + " " + p.Nom
}

// Même chose ici, on utilise un receiver par valeur pour construire la phrase de présentation du contact.
func (p Personne) Presentation() string {
	return fmt.Sprintf("%s, %d ans - %s", p.NomComplet(), p.Age, p.Email)
}

// On définit la structure Adresse de manière indépendante, pour qu'elle soit réutilisable par composition.
// Ici aussi, des tags Json sont spécifiés
type Adresse struct {
	Rue        string `json:"rue"`
	Ville      string `json:"ville"`
	CodePostal string `json:"code_postal"`
}

// On définit une fonction explicite qui retourne l'adresse complète
func (a Adresse) Format() string {
	return fmt.Sprintf("%s, %s %s", a.Rue, a.CodePostal, a.Ville)
}

// La structure "Employe" utilise par composition les structures "Personne" et "Adresse".
// L'embedding nous permet d'appeler directement "e.NomComplet()",
// sans devoir passer par "e.Personne.NomComplet()" (bien que ca reste possible).
// Cette portion de code démontre donc bien que le combo "Structs + Méthodes + embedding"
// est une alternative viable à l'héritage, qui n'existe pas dans Go
// On a fait en sorte que la propriété "Salaire" d'un employé ne figure pas dans les métadonnée du JSON produit.
type Employe struct {
	Personne
	Adresse
	Poste   string  `json:"poste"`
	Salaire float64 `json:"-"`
}

// On déclare une fonction qui affiche toutes les informations d'un employé
func (e Employe) FicheEmploye() string {
	// e.Adresse.Format() est explicite ici car Personne et Adresse ont toutes les deux
	// des méthodes. Si on ne le précisait pas, Go ne saurait pas laquelle promouvoir.
	return fmt.Sprintf("Employé : %s\nAdresse : %s\nPoste : %s | Salaire : %.2f€",
		e.Presentation(), e.Adresse.Format(), e.Poste, e.Salaire)
}

// On déclare cette fonction qui utilise un receiver par pointeur (*Employe) et modifie donc la structure originale
func (e *Employe) AugmenterSalaire(pct float64) {
	e.Salaire += e.Salaire * pct / 100
}

// On déclare Etudiant qui est une structure composée uniquement de la structure Personne,
// car les étudiants n'ont pas d'adresse.
// Toutes les propriétés sont renseignées dans le JSON produit
type Etudiant struct {
	Personne
	Promo   string  `json:"promo"`
	Moyenne float64 `json:"moyenne"`
}

// Cette fonction permet d'obtenir la mention de l'étudiant selon sa moyenne
func (et Etudiant) MentionObtenue() string {
	switch {
	case et.Moyenne >= 16:
		return "TB"
	case et.Moyenne >= 14:
		return "B"
	case et.Moyenne >= 12:
		return "AB"
	default:
		return "P"
	}
}

// On déclare une fonction qui affiche toutes les informations d'un étudiant
func (et Etudiant) FicheEtudiant() string {
	return fmt.Sprintf("Étudiant : %s\nPromo : %s | Moyenne : %.2f (%s)",
		et.Presentation(), et.Promo, et.Moyenne, et.MentionObtenue())
}

func main() {
	// On définit une slice de structure, pour renseigner en dur les données des employés
	employes := []Employe{
		{
			Personne: Personne{"Raphaêl", "Canetti", 39, "rcanetti@myges.fr"},
			Adresse:  Adresse{"125 Rue Baraban", "Lyon", "69003"},
			Poste:    "Enseignant",
			Salaire:  3500,
		},
		{
			Personne: Personne{"Bénédicte", "Maurin-Crépet", 50, "bmaurin@myges.fr"},
			Adresse:  Adresse{"53 Cours Albert Thomas", "Lyon", "69003"},
			Poste:    "Responsable pédagogique",
			Salaire:  4100,
		},
	}

	// On définit une slice de structure, pour renseigner en dur les données des étudiants
	etudiants := []Etudiant{
		{
			Personne: Personne{"Louis", "Cauvet", 24, "lcauvet@myges.fr"},
			Promo:    "M2 Ingénierie du web",
			Moyenne:  12.4,
		},
		{
			Personne: Personne{"Duperthuy", "Hugo", 24, "hduperthuy@myges.fr"},
			Promo:    "M2 Ingénierie du web",
			Moyenne:  14.26,
		},
	}

	// On peut ajouter dynamiquement un élément à la slice,
	// puisque celle-ci n'a pas de taille fixe
	etudiants = append(etudiants, Etudiant{
		Personne: Personne{"Lisa", "Michallon", 28, "lmichallont@myges.fr"},
		Promo:    "M2 Ingénierie du web",
		Moyenne:  17.0,
	})

	// A l'aide d'un map, on associe un nom complet à la fiche contact correspondante.
	annuaire := make(map[string]string)

	fmt.Println("=== Employés ===")
	// On effectue une boucle par index (et non par valeur) sur les employés
	// pour que AugmenterSalaire() modifie bien la structure dans la slice
	for i := range employes {
		// On augmente le salaire de 10% pour chaque employé
		employes[i].AugmenterSalaire(10)
		// On récupère sa fiche
		fiche := employes[i].FicheEmploye()
		annuaire[employes[i].NomComplet()] = fiche
		// On affiche la fiche de l'employé
		fmt.Println(fiche)
		fmt.Println()
	}

	// On récupère une variable (le premier employé) grâce à son adresse mémoire avec "&"
	premier := &employes[0]
	// On modifie le poste de l'employé en pointant sur sa valeur
	premier.Poste = "Développeur Go"
	// On accède explicitement à la valeur du pointeur avec "*"
	fmt.Printf("Poste modifié via pointeur : %s de %s\n", (*premier).Poste, (*premier).NomComplet())
	fmt.Println()

	fmt.Println("=== Étudiants ===")
	// Dans cette boucle, il n'y a pas de modification à réaliser donc on peut itérer sur "_"
	for _, et := range etudiants {
		// ON récupère la fiche étudiante
		fiche := et.FicheEtudiant()
		annuaire[et.NomComplet()] = fiche
		// On affiche la fiche de l'étudiant
		fmt.Println(fiche)
		fmt.Println()
	}

	// On met en place un idiome pour lire dans la map des étudiants,
	// et renvoyer la valeur si elle est trouvée,
	// ou rien du tout dans le cas contraire (car ok vaut alors "false")
	if fiche, ok := annuaire["Louis Cauvet"]; ok {
		fmt.Println("=== Fiche recherchée dans l'annuaire ===")
		fmt.Println(fiche)
		fmt.Println()
	}

	// On affiche le total du nombre de contacts
	fmt.Printf("Annuaire : %d contacts enregistrés\n", len(annuaire))
}
