# PLAYZONE

Forum communautaire gaming - projet supervise B1.

## Equipe

- AMIR Mohand Arezki
- Matheo
- ADIB Houssine

## Idee

Une plateforme de discussion type Reddit / Stack Overflow, orientee gaming. Les utilisateurs peuvent creer des fils par jeu, poster des messages, reagir avec like/dislike, et un admin modere la plateforme.

## Technologies

- **Golang** (backend, architecture MVC en couches)
- **SQLite** (base de donnees relationnelle)
- **go-chi** (routeur HTTP)
- **JWT** (authentification, golang-jwt/v5)
- **html/template** (vues)
- HTML / CSS pur

## Architecture

Le projet suit un decoupage en couches inspire du modele MVC :

```
playzone/
├── main.go              -> point d'entree
├── config/              -> constantes (port, JWT secret, BDD)
├── database/            -> connexion BDD + schema + seed
├── router/              -> definition des routes
├── controllers/         -> handlers HTTP, validation entrees
├── services/            -> logique metier
├── repositories/        -> acces aux donnees
├── models/              -> structs / entites
├── middleware/          -> authentification JWT, role admin
├── views/               -> templates HTML
├── static/              -> CSS, JS, images
└── data/                -> base SQLite (auto-creee)
```

## Prerequis

- **Go** version 1.21 ou plus recente
- Connexion internet pour les polices Google (Orbitron, Inter)

## Installation et lancement

```bash
git clone https://github.com/Mohand-123/projet-forum.git
cd projet-forum/playzone
go mod download
go run .
```

Le serveur demarre sur `http://localhost:3000`.

## Comptes de demo

Au premier lancement, la base est seedee avec deux comptes :

| Pseudo | Mot de passe   | Role  |
|--------|----------------|-------|
| admin  | Admin1234567!  | admin |
| demo   | Demo1234567!   | user  |

## Fonctionnalites implementees

- **FT-1** Inscription (mdp 12 chars min, 1 maj, 1 special, hashe SHA-512 + salt)
- **FT-2** Connexion par pseudo ou email, token JWT en cookie
- **FT-3** Creation de fils avec etats (ouvert / ferme / archive)
- **FT-4** Consultation des fils (archive masque)
- **FT-5** Publication de messages dans un fil ouvert
- **FT-6** Like / Dislike avec score (unique par user/message)
- **FT-7** Edition / suppression de son contenu, admin peut tout supprimer
- **FT-8** Tri des messages par date (recent / ancien) ou popularite
- **FT-9** Pagination par lots de 10 / 20 / 30 / tout
- **FT-10** Affichage des fils par tag / categorie
- **FT-11** Recherche par titre, tag ou categorie (utilisateur authentifie)
- **FT-12** Tableau de bord admin avec ban / debannissement

## Liens utiles

- Depot GitHub : https://github.com/Mohand-123/projet-forum
- Tableau Trello : https://trello.com/invite/b/6a102eea2bdd75bde684a5dd/ATTIe56f5810b83ee8f9c5215eaaebfa303fED55119A/mon-tableau-trello
