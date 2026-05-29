package main

import (
	"errors"
	"fmt"
	"strings"
)

// On définit un type "Categorie".
// Grâce au iota, chaque constante de ce type reçoit automatiquement la valeur entière suivante (0, 1, 2...).
type Categorie int

const (
	Smartphone Categorie = iota // 0
	Ordinateur                  // 1
	Tablette                    // 2
	Accessoire                  // 3
	Autre                       // 4
)

// On crée une fonction qui retourne le nom lisible d'une catégorie.
// On utilise le fallthrough pour que Smartphone et Tablette partagent le même label "Mobile" :
// quand l'élément observé vaut Smartphone, on tombe dans le case Tablette et on retourne "Mobile" dans les deux cas.
func labelCategorie(c Categorie) string {
	switch c {
	case Smartphone:
		fallthrough
	case Tablette:
		return "Mobile"
	case Ordinateur:
		return "Ordinateur"
	case Accessoire:
		return "Accessoire"
	default:
		return "Autre"
	}
}

// On centralise les erreurs métier une seule fois ici,
// pour pouvoir les comparer précisément plus tard avec errors.Is().
var (
	ErrIDDuplique       = errors.New("ID déjà existant")
	ErrStockInsuffisant = errors.New("stock insuffisant")
	ErrProduitInconnu   = errors.New("produit introuvable")
)

// On définit la structure Produit avec des tags JSON sur chaque propriété (chaque propriété est exportable car elles commencent par des majuscules).
// Le champ "actif" est en minuscule : il n'est pas exporté et n'apparaît donc pas en JSON.
type Produit struct {
	ID        int       `json:"id"`
	Nom       string    `json:"nom"`
	Marque    string    `json:"marque"`
	Prix      float64   `json:"prix"`
	Stock     int       `json:"stock"`
	Categorie Categorie `json:"categorie"`
	actif     bool
}

// On créé une fonction qui initialise un produit actif, avec toutes ses propriétés.
func nouveauProduit(id int, nom, marque string, prix float64, stock int, cat Categorie) Produit {
	return Produit{ID: id, Nom: nom, Marque: marque, Prix: prix, Stock: stock, Categorie: cat, actif: true}
}

// On crée une méthode d'affichage par défaut d'un produit.
// En implémentant cette méthode, fmt.Println(p) appellera automatiquement ce format.
func (p Produit) String() string {
	return fmt.Sprintf("[%d] %s %s — %.2f€ | Stock: %d | %s",
		p.ID, p.Marque, p.Nom, p.Prix, p.Stock, labelCategorie(p.Categorie))
}

// On définit la structure Catalogue qui contient une slice de Produit.
// Pour éviter que de la duplication d'ID se produise, on définit la propriété "produits" en minuscule,
// ainsi, on force l'accès via les méthodes implémentées dans ce but.
// La prorpiété "prochainID" a pour vocation d'être maintenu automatiquement pour éviter que l'utilisateur ne saisisse les id à la main.
type Catalogue struct {
	produits   []Produit
	prochainID int
}

// On crée une fonction qui ajoute un produit dans le catalogue.
// On utilise un receiver par pointeur (*Catalogue) pour modifier la slice originale.
func (c *Catalogue) AjouterProduit(p Produit) error {
	for _, existant := range c.produits {
		// Si l'ID existe déjà, on retourne une erreur
		if existant.ID == p.ID {
			return fmt.Errorf("AjouterProduit id=%d : %w", p.ID, ErrIDDuplique)
		}
	}
	// append() ajoute dynamiquement un élément à la slice puisqu'elle n'a pas de taille fixe
	c.produits = append(c.produits, p)
	// On maintient prochainID à jour pour qu'il reste toujours supérieur à tous les IDs existants
	if p.ID >= c.prochainID {
		c.prochainID = p.ID + 1
	}
	return nil
}

// On crée une fonction variadique, qui permet de passer
// autant de produits que l'on veut en paramètres.
func (c *Catalogue) AjouterProduits(produits ...Produit) []error {
	var errs []error
	for _, p := range produits {
		// On retourne une slice des erreurs rencontrées (une par produit problématique).
		if err := c.AjouterProduit(p); err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}

// On crée une fonction qui recherche un produit par son identifiant.
func (c *Catalogue) TrouverParID(id int) (Produit, error) {
	for _, p := range c.produits {
		if p.ID == id {
			return p, nil
		}
	}
	// On retourne une erreur wrappée si aucun produit ne correspond.
	return Produit{}, fmt.Errorf("TrouverParID id=%d : %w", id, ErrProduitInconnu)
}

// On définit une fonction qui utilise une closure pour filtrer les produits.
// La closure "correspond" capture la variable "cat" de la fonction parente
// et peut être appelée sur chaque produit comme un filtre réutilisable.
func (c *Catalogue) TrouverParCategorie(cat string) []Produit {
	// strings.EqualFold permet une comparaison insensible à la casse ("mobile" == "Mobile")
	correspond := func(p Produit) bool {
		return strings.EqualFold(labelCategorie(p.Categorie), cat)
	}
	var resultats []Produit
	for _, p := range c.produits {
		// Si le produit recherché est trouvé, on le renvoie
		if correspond(p) {
			resultats = append(resultats, p)
		}
	}
	return resultats
}

// On crée une fonction qui applique un pourcentage de réduction sur tous les produits d'une catégorie.
// On itère par index pour modifier directement la structure dans la slice (on ne travaille pas avec une copie ici).
func (c *Catalogue) AppliquerReduction(categorie string, pct float64) int {
	count := 0
	for i, p := range c.produits {
		if strings.EqualFold(labelCategorie(p.Categorie), categorie) {
			c.produits[i].Prix -= p.Prix * pct / 100
			count++
		}
	}
	return count
}

// On crée une fonction qui réduit le stock d'un produit de la quantité demandée.
// Si le stock tombe à 0, le produit est automatiquement marqué inactif (champ non exporté "actif").
func (c *Catalogue) Vendre(id int, qte int) error {
	for i, p := range c.produits {
		if p.ID == id {
			if p.Stock < qte {
				return fmt.Errorf("Vendre id=%d qte=%d : %w", id, qte, ErrStockInsuffisant)
			}
			c.produits[i].Stock -= qte
			if c.produits[i].Stock == 0 {
				// actif est accessible ici bien qu'il soit "privé", car on est dans le même package
				c.produits[i].actif = false
			}
			return nil
		}
	}
	return fmt.Errorf("Vendre id=%d : %w", id, ErrProduitInconnu)
}

// On créé une fonction qui construit un résumé du catalogue.
// On utilise une map pour compter le nombre de produits par catégorie :
// la clé est le nom de la catégorie (string), la valeur est le compteur (int).
func (c *Catalogue) Rapport() string {
	valeurTotale := 0.0
	parCategorie := make(map[string]int)

	for _, p := range c.produits {
		valeurTotale += p.Prix * float64(p.Stock)
		parCategorie[labelCategorie(p.Categorie)]++
	}

	rapport := fmt.Sprintf("Nb produits : %d | Valeur totale du stock : %.2f€\n", len(c.produits), valeurTotale)
	// On parcourt la map pour afficher le détail par catégorie
	for cat, nb := range parCategorie {
		rapport += fmt.Sprintf("  %-12s : %d produit(s)\n", cat, nb)
	}
	return rapport
}

func main() {
	// defer enregistre une fonction à exécuter à la fin de main, quoi qu'il arrive.
	// Ainsi, même si main se termine par un return anticipé, ce message sera toujours affiché.
	defer fmt.Println("\nMerci d'avoir utilisé TechShop. À bientôt !")

	catalogue := &Catalogue{}

	// On pré-remplit le catalogue via la fonction "AjouterProduits" qui accepte N produits en un seul appel
	errs := catalogue.AjouterProduits(
		nouveauProduit(1, "iPhone 15 Pro", "Apple", 1329.00, 12, Smartphone),
		nouveauProduit(2, "MacBook Pro M3", "Apple", 2499.00, 5, Ordinateur),
		nouveauProduit(3, "Galaxy S24", "Samsung", 899.00, 20, Smartphone),
		nouveauProduit(4, "iPad Air", "Apple", 799.00, 8, Tablette),
		nouveauProduit(5, "AirPods Pro", "Apple", 279.00, 30, Accessoire),
	)
	// On anticipe les éventuelles erreurs d'initialisation
	for _, err := range errs {
		fmt.Println("Erreur init :", err)
	}

	// On entre dans la boucle principale du menu CLI
	for {
		fmt.Println("\n=== TechShop — Boutique CLI ===")
		fmt.Println("[1] Ajouter  [2] Chercher  [3] Soldes  [4] Vendre  [5] Rapport  [0] Quitter")
		fmt.Print("> ")

		var choix int
		fmt.Scan(&choix)

		switch choix {
		case 0:
			// dans ce cas, cela déclenche le defer défini en début de main
			return
		case 1:
			ajouterProduitCLI(catalogue)
		case 2:
			chercherProduitCLI(catalogue)
		case 3:
			soldesCLI(catalogue)
		case 4:
			vendreCLI(catalogue)
		case 5:
			fmt.Println(catalogue.Rapport())
		default:
			fmt.Println("Choix invalide.")
		}
	}
}

// On crée une fonction qui lit les informations d'un produit depuis le terminal
// et l'ajoute au catalogue. L'ID est assigné automatiquement via prochainID, afin d'éviter les doublons et les erreurs de saisie de l'utilisateur
func ajouterProduitCLI(c *Catalogue) {
	var stock, cat int
	var nom, marque string
	var prix float64

	fmt.Print("Nom : ")
	fmt.Scan(&nom)
	fmt.Print("Marque : ")
	fmt.Scan(&marque)
	fmt.Print("Prix : ")
	fmt.Scan(&prix)
	fmt.Print("Stock : ")
	fmt.Scan(&stock)
	fmt.Println("Catégorie : [0] Smartphone  [1] Ordinateur  [2] Tablette  [3] Accessoire  [4] Autre")
	fmt.Print("> ")
	fmt.Scan(&cat)

	// L'ID est attribué automatiquement depuis le compteur interne du catalogue
	id := c.prochainID

	// Si une erreur à été rencontrée durant l'ajout du produit, on indique une erreur
	if err := c.AjouterProduit(nouveauProduit(id, nom, marque, prix, stock, Categorie(cat))); err != nil {
		fmt.Println("Erreur :", err)
		return
	}
	fmt.Printf("Produit ajouté avec l'ID %d.\n", id)
}

// On créé une fonction qui affiche tous les produits d'une catégorie donnée.
// La comparaison est insensible à la casse grâce à strings.EqualFold dans la fonction "TrouverParCategorie".
func chercherProduitCLI(c *Catalogue) {
	fmt.Print("Catégorie (Mobile / Ordinateur / Accessoire / Autre) : ")
	var cat string
	fmt.Scan(&cat)

	produits := c.TrouverParCategorie(cat)
	// On gère le cas où aucun produit n'est trouvé
	if len(produits) == 0 {
		fmt.Println("Aucun produit trouvé.")
		return
	}
	// fmt.Println appelle automatiquement la méthode String() de chaque Produit
	for _, p := range produits {
		fmt.Println(p)
	}
}

// On crée une fonction qui applique une réduction en pourcentage sur une catégorie entière.
func soldesCLI(c *Catalogue) {
	fmt.Print("Catégorie à solder : ")
	var cat string
	fmt.Scan(&cat)
	fmt.Print("Réduction (%) : ")
	var pct float64
	fmt.Scan(&pct)

	n := c.AppliquerReduction(cat, pct)
	fmt.Printf("%d produit(s) mis en solde avec -%.0f%%.\n", n, pct)
}

// On crée une fonction qui enregistre une vente en réduisant le stock du produit ciblé.
func vendreCLI(c *Catalogue) {
	fmt.Print("ID produit : ")
	var id int
	fmt.Scan(&id)
	fmt.Print("Quantité : ")
	var qte int
	fmt.Scan(&qte)

	if err := c.Vendre(id, qte); err != nil {
		if errors.Is(err, ErrStockInsuffisant) {
			// On gère le cas où le stock du produit est insuffisant
			fmt.Println("Impossible : stock insuffisant.")
		} else if errors.Is(err, ErrProduitInconnu) {
			// On gère le cas où le stock du produit n'est pas reconnu
			fmt.Println("Impossible : produit introuvable.")
		} else {
			fmt.Println("Erreur :", err)
		}
		return
	}
	fmt.Println("Vente enregistrée.")
}
