# yassinouk.github.io
# Readme
![Hero Image](ensate.png) <!-- hero image -->
## Description
Ce code implémente une simulation basique d'un canal LTE (Long-Term Evolution) en utilisant le langage de programmation Go. Le canal LTE est modélisé avec des fonctionnalités telles que la modulation/démodulation BPSK, la modulation OFDM (Orthogonal Frequency-Division Multiplexing), l'effet de diffusion de Rayleigh, et l'ajout de bruit blanc additif gaussien (AWGN). Le résultat de la simulation est affiché, et un serveur web est également lancé pour fournir des données en réponse à des requêtes HTTP.

## Fonctionnalités principales
- **Modulation BPSK**: Les bits générés sont modulés en symboles complexes bpsk.
- **Modulation OFDM**: Les symboles modulés sont soumis à une modulation OFDM à l'aide de la transformée de Fourier rapide (FFT/IFFT).
- **Effet de diffusion de Rayleigh**: Modélisation de l'effet de la diffusion de Rayleigh sur les symboles transmis.
- **Bruit blanc additif gaussien (AWGN)**: Ajout d'un niveau de bruit gaussien aux symboles transmis pour simuler des conditions réelles.

## Utilisation
Pour exécuter la simulation et le serveur web, assurez-vous d'avoir Go installé sur votre machine. Exécutez ensuite le fichier principal `go run main.go`. Le serveur web sera disponible à l'adresse [http://localhost:8080](http://localhost:8080).
ou bien cliquer sur le fichier executable avec index.html dans le meme dossier, finalement vous pouvez visiter le test deployment via:[vercel.app](http://localhost:8080)
Le serveur web expose trois endpoints API pour obtenir des données sous forme de réponses JSON :
- `/api/endpoint1`: Les bits démodulés après la simulation.
- `/api/endpoint2`: Les symboles transmis modélisés comme des entiers après l'effet de diffusion de Rayleigh.
- `/api/endpoint3`: Les bits d'origine générés pour la simulation.

## Exemple d'exécution
1. Générer 100 bits aléatoires.
2. Moduler les bits en symboles complexes.
3. Appliquer la modulation OFDM.
4. Modéliser l'effet de diffusion de Rayleigh.
5. Ajouter du bruit blanc additif gaussien (AWGN).
6. Démoduler les symboles transmis.
7. Comparer les bits d'origine et les bits démodulés.
8. Le ratio de correspondance des bits est affiché à la fin de l'exécution.

## Dépendances tierces
Ce projet utilise le package tiers `github.com/mjibson/go-dsp/fft` pour les transformations FFT/IFFT.

## Auteur
[OUAKKA Yassin]

## Licence
Ce projet est sous licence MIT.
