# GoproToGpxAPI

GoproToGpxAPI est une API qui permet d'extraire les données de télémétrie GPS d'une vidéo GoPro et de les convertir en un fichier GPX.

## Comment ça marche

L'API reçoit le contenu d'un fichier via une requête HTTP POST multipart/form-data. Elle utilise ensuite `exiftool` pour extraire les métadonnées du fichier. Si le fichier contient des données de télémétrie GPS, ces données sont extraites et converties en un fichier GPX.

## Comment l'utiliser

1. Envoyez une requête POST à l'URL `/video` avec un formulaire multipart contenant un champ de fichier nommé "file". Le fichier doit être une image de vidéo GoPro (.bin) contenant des données de télémétrie GPS.

2. L'API répondra avec l'UUID de la vidéo en base lié au GPX généré à partir des données de télémétrie GPS de la vidéo.

## Installation

1. Clonez ce dépôt.
2. Installez `exiftool` sur votre système.
3. Exécutez `make run` pour run le projet.
4. Exécutez l'exécutable généré pour démarrer le serveur.

## Tests

Pour exécuter les tests, utilisez la commande `go test -v ./...`.

## Nettoyage

Pour supprimer l'exécutable généré, utilisez la commande `make clean`.

## Licence

Ce projet est sous licence MIT.