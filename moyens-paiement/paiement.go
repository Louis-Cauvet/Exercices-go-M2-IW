package main

import (
	"fmt"
	"math"
)

// Interface commune à tous les moyens de paiement
type Payeur interface {
	Payer(montant float64) (string, error)
}

// Carte bancaire
type CarteCredit struct {
	Numero    string
	Titulaire string
	Solde     float64
}

func (c *CarteCredit) Payer(montant float64) (string, error) {
	if montant > c.Solde { // refuse le paiement si le solde ne couvre pas le montant
		return "", fmt.Errorf("solde insuffisant sur la carte de %s", c.Titulaire)
	}
	c.Solde -= montant // débite le solde (pointeur requis pour modifier la struct d'origine)
	return fmt.Sprintf("Transaction CB #%s confirmée", c.Numero), nil
}

// PayPal
type PayPal struct {
	Email string
	Solde float64
}

func (p *PayPal) Payer(montant float64) (string, error) {
	if montant > p.Solde { // même vérification que pour la carte
		return "", fmt.Errorf("solde insuffisant pour %s", p.Email)
	}
	p.Solde -= montant
	return fmt.Sprintf("Paiement PayPal de %.2f€ vers %s", montant, p.Email), nil
}

// Crypto
type Crypto struct {
	Adresse string
	Solde   float64
	Monnaie string
}

func (c *Crypto) Payer(montant float64) (string, error) {
	if montant > c.Solde {
		return "", fmt.Errorf("solde insuffisant sur le wallet %s", c.Adresse)
	}
	c.Solde -= montant
	// conversion € → BTC (1 BTC = 50000€), arrondie à 3 décimales
	quantite := math.Round(montant/50000*1000) / 1000
	return fmt.Sprintf("Paiement de %.3f %s vers %s", quantite, c.Monnaie, c.Adresse), nil
}

// Traite le panier avec le moyen de paiement fourni
func ProcesserPanier(payeur Payeur, articles []float64) {
	total := 0.0
	for _, prix := range articles {
		total += prix
	}
	fmt.Printf("Total du panier : %.2f€\n", total)

	// type switch : identifie le type concret caché derrière l'interface Payeur
	switch payeur.(type) {
	case *CarteCredit:
		fmt.Println("Mode de paiement : Carte bancaire")
	case *PayPal:
		fmt.Println("Mode de paiement : PayPal")
	case *Crypto:
		fmt.Println("Mode de paiement : Crypto-monnaie")
	}

	// Chaque type exécute sa propre logique de Payer
	resultat, err := payeur.Payer(total)
	if err != nil {
		fmt.Println("Erreur :", err)
		return
	}
	fmt.Println(resultat)
}

func main() {
	articles := []float64{29.99, 15.50, 4.99}

	cb := &CarteCredit{Numero: "1234", Titulaire: "Alice", Solde: 100}
	ProcesserPanier(cb, articles)

	fmt.Println()

	pp := &PayPal{Email: "bob@mail.com", Solde: 100}
	ProcesserPanier(pp, articles)

	fmt.Println()

	crypto := &Crypto{Adresse: "0xABC", Solde: 100, Monnaie: "BTC"}
	ProcesserPanier(crypto, articles)
}
