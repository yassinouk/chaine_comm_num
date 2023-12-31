package main

import (
	"encoding/json" // importation du package pour le codage / décodage JSON
	"fmt"           // pour l'impression formattée
	"log"           // pour le journalisation des erreurs
	"math"          // pour les fonctions mathématiques de base
	"math/cmplx"    // pour les fonctions complexes
	"math/rand"     // pour la génération de nombres aléatoires
	"net/http"      // pour la manipulation des requêtes et réponses HTTP

	"github.com/mjibson/go-dsp/fft" // package tiers pour les transformations On utilise ici Fast Fourier Transform (FFT)
)

const (
	noiseFloor = -90.0 // constant utilisée pour définir le plancher de bruit, utilisée dans la fonction de bruit blanc additif gaussien (AWGN)
)

type Data struct { // structure pour stocker les données en y
	Y []int64 `json:"y"` // les données réelles sont de type tableau d'entiers
}

type LTEChannel struct{} // définit un nouveau type pour représenter un canal LTE

func (c *LTEChannel) Modulate(bits []int64) []complex128 { // fonction pour moduler les bits dans les symboles complexes. Chaque bit est converti en symbole complexe
	symbols := make([]complex128, len(bits)) // création d'un tableau pour stocker les symboles

	for i, b := range bits { // iterate over each bit
		symbols[i] = complex(float64(2*b-1), 0) // convert the bit to a complex symbol
	}

	return symbols // retourne les symboles après la modulation
}

func (c *LTEChannel) OFDMModulate(symbols []complex128) []complex128 { // fonction pour la modulation OFDM des symboles
	return fft.IFFT(symbols) // utiliser la transformation inverse de Fourier rapide (IFFT) pour obtenir les symboles modulés OFDM
}

func (c *LTEChannel) OFDMDemodulate(symbols []complex128) []complex128 { // fonction pour la démodulation OFDM des symboles
	return fft.FFT(symbols) // utiliser la transformation de Fourier rapide (FFT) pour obtenir les symboles démodulés OFDM
}

func (c *LTEChannel) Demodulate(symbols []complex128) []int64 { // fonction pour démoduler les symboles en bits
	bits := make([]int64, len(symbols)) // création d'un tableau pour stocker les bits

	for i, s := range symbols { // pour chaque symbole
		if real(s) >= 0 { // si la partie réelle du symbole est supérieure ou égale à zéro
			bits[i] = 1 // le bit est 1
		} else { // sinon
			bits[i] = 0 // le bit est 0
		}
	}

	return bits // retourne les bits après démodulation
}

func (c *LTEChannel) Rayleigh(symbols []complex128) []complex128 { // fonction pour modéliser l'effet de la diffusion de Rayleigh sur les symboles
	for i := range symbols { // pour chaque symbole
		gain := complex(rand.NormFloat64(), rand.NormFloat64()) // génération d'un gain complexe aléatoire suivant une distribution normale
		symbols[i] = symbols[i] * gain                          // multiplication du symbole par le gain pour modéliser l'effet de la diffusion de Rayleigh
	}
	fmt.Println("symbols output from rayleigh channel\n", symbols)

	return symbols // retourne les symboles après l'effet de la diffusion de Rayleigh
}

func (c *LTEChannel) AWGN(symbols []complex128, snr float64) []complex128 { // fonction pour ajouter un bruit blanc additif gaussien (AWGN) aux symboles
	power := math.Pow(10, (noiseFloor-snr)/10)           // calcule la puissance du bruit
	symbolsPlusNoise := make([]complex128, len(symbols)) // créer un tableau pour stocker les symboles plus le bruit

	for i := range symbols { // pour chaque symbole
		noise := complex(rand.Float64()*power, 0) // générer un bruit complexe aléatoire proportionnel à la puissance du bruit
		symbolsPlusNoise[i] = symbols[i] + noise  // ajout du bruit au symbole
	}

	return symbolsPlusNoise // retourne les symboles après ajout du bruit AWGN
}

func (c *LTEChannel) Transmit(bits []int64, snr float64) []complex128 { // fonction pour transmettre des bits sur le canal LTE
	symbols := c.Modulate(bits)       // modulation des bits en symboles
	symbols = c.OFDMModulate(symbols) // modulation OFDM des symboles
	symbols = c.Rayleigh(symbols)     // ajout de l'effet de la diffusion de Rayleigh aux symboles
	return c.AWGN(symbols, snr)       // ajout du bruit AWGN aux symboles et retour des symboles transmis
}

func CalculateNorm(c complex128) int64 { // fonction pour calculer la norme d'un nombre complexe
	return int64(1e4 * cmplx.Abs(c)) // retourne la norme du nombre complexe
}

func main() { // fonction principale
	lteChannel := &LTEChannel{} // crée un nouveau canal LTE
	bits := make([]int64, 1000) // génère un tableau de 100 bits
	for i := range bits {       // pour chaque bit
		bits[i] = int64(rand.Intn(2)) // génère un bit aléatoire (0 ou 1)
	}
	fmt.Println("bernouilli generated bits\n", bits) // affiche les bits générés

	transmitted := lteChannel.Transmit(bits, 20.0) // transmet les bits sur le canal LTE avec un SNR de 20 dB
	dem := lteChannel.OFDMDemodulate(transmitted)  // démodule les symboles transmis
	outbits := lteChannel.Demodulate(dem)          // démodule les symboles en bits

	fmt.Println("original bits\n", bits)
	fmt.Println("+++++++++++++++++++++++++++++++++++++++++++") // affiche les bits d'origine
	fmt.Println("+++++++++++++++++++++++++++++++++++++++++++") // affiche les bits d'origine
	fmt.Println("transmitted\n", transmitted)                  // affiche les symboles transmis
	fmt.Println("+++++++++++++++++++++++++++++++++++++++++++")
	fmt.Println("+++++++++++++++++++++++++++++++++++++++++++")
	fmt.Println("ofdm demodulation:\n", dem) // affiche les symboles démodulés
	fmt.Println("+++++++++++++++++++++++++++++++++++++++++++")
	fmt.Println("+++++++++++++++++++++++++++++++++++++++++++")
	fmt.Println("outbits demodulated bpsk (or finale demodulated):\n", outbits) // affiche les bits de sortie
	fmt.Println("+++++++++++++++++++++++++++++++++++++++++++")
	fmt.Println("+++++++++++++++++++++++++++++++++++++++++++")
	bitMatch := 0         // initialise le compteur de correspondance de bits
	for i := range bits { // pour chaque bit
		if bits[i] == outbits[i] { // si le bit d'origine correspond au bit de sortie
			bitMatch++ // incrémenter le compteur de correspondance de bits
		}
	}

	fmt.Printf("Bit match ratio: %f \n", float64(bitMatch)/float64(len(bits)))
	// affiche le ratio de correspondance des bits

	// Créer un gestionnaire de serveur de fichiers
	fs := http.FileServer(http.Dir("."))

	// Affecte toutes les URL à servir en tant que fichiers statiques à partir du répertoire courant
	http.Handle("/", fs)

	// Enregistre dans le journal si le serveur ne parvient pas à démarrer. Démarre le serveur sur le port 8080

	http.HandleFunc("/api/endpoint1", func(w http.ResponseWriter, r *http.Request) {
		ServeData(w, r, func() []int64 { return outbits })
	})
	http.HandleFunc("/api/endpoint2", func(w http.ResponseWriter, r *http.Request) {
		ServeData(w, r, func() []int64 {
			intBits := make([]int64, len(transmitted))
			fmt.Println("transmitted", transmitted)
			for i, t := range transmitted {
				intBits[i] = CalculateNorm(t)
			}
			return intBits
		})
	})
	http.HandleFunc("/api/endpoint3", func(w http.ResponseWriter, r *http.Request) {
		ServeData(w, r, func() []int64 { return bits })
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
	// démarre le serveur web sur le port 8080 et enregistre dans le journal en cas d'erreur
}

func ServeData(w http.ResponseWriter, r *http.Request, genData func() []int64) {
	// fonction pour servir les données en tant que réponse JSON à une requête HTTP
	data := Data{Y: genData()} // génère les données

	w.Header().Set("Content-type", "application/json") // définir le type de contenu de la réponse
	json.NewEncoder(w).Encode(data)                    // encoder les données en JSON et les écrire dans la réponse
}
