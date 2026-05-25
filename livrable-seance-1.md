# Projet supervisé B1 — Séance 1
# PLAYZONE — Plateforme communautaire gaming

**Équipe :** AMIR Mohand Arezki · MATHÉO · ADIB Houssine
**Date :** Mai 2026

🔗 **Dépôt GitHub :** https://github.com/Mohand-123/projet-forum
📋 **Tableau Trello :** https://trello.com/b/0C55Nrs2/projet-forum

---

## 1. Présentation générale du projet

### 1.1 Concept

**PLAYZONE** est une plateforme communautaire en ligne dédiée aux **joueurs et passionnés d'esport**. Elle se présente sous la forme d'un forum moderne où les gamers peuvent se retrouver autour de leurs jeux favoris, discuter, partager leurs expériences et faire connaissance avec d'autres joueurs.

Le slogan retenu est : *"Where players meet."*

### 1.2 Cible

| Profil | Description |
|---|---|
| Joueurs casuals | 16–25 ans, jouent quelques heures par semaine, cherchent à discuter et trouver des amis |
| Joueurs compétitifs | 18–30 ans, MOBA / FPS / fighting games, échangent stratégies et coéquipiers |
| Streamers / créateurs | Partagent leurs clips et highlights, cherchent une audience |
| Fans d'esport | Suivent les compétitions, discutent des équipes et des joueurs pros |

### 1.3 Inspirations

- **Reddit** — structure en catégories thématiques, fils et commentaires
- **Discord** — communautés dédiées par jeu, dark mode
- **Twitch** — codes visuels gaming, ambiance esport
- **Steam Community** — profils joueurs, partage de guides

### 1.4 Identité visuelle

- **Mode sombre** par défaut (norme dans le gaming)
- **Palette** : violet néon `#7c3aed`, cyan électrique `#06b6d4`, fond noir profond `#0a0e1a`
- **Typographies** : **Orbitron** (titres, look sci-fi) et **Inter** (corps de texte)

---

## 2. Utilisateurs identifiés, rôles et actions

La plateforme distingue **deux types d'utilisateurs** :

### 2.1 Visiteur (utilisateur non authentifié)

**Actions possibles :**

| Code | Action |
|---|---|
| V1 | Consulter la liste des catégories de jeux |
| V2 | Parcourir les fils de discussion d'une catégorie |
| V3 | Lire un fil et toutes ses réponses |
| V4 | S'inscrire (créer un compte) |
| V5 | Se connecter à un compte existant |

### 2.2 Utilisateur authentifié

> Hérite de toutes les actions du Visiteur, et en plus :

- Créer un nouveau fil dans une catégorie
- Répondre dans un fil existant
- Réagir aux messages (like / dislike)
- Modifier / supprimer son propre contenu
- Se déconnecter

### 2.3 Administrateur

> Hérite des actions de l'Utilisateur, et en plus :

- Accéder au tableau de bord d'administration
- Modifier l'état d'un fil (ouvert / fermé / archivé)
- Supprimer n'importe quel fil ou message
- Bannir / débannir un utilisateur

---

## 3. Conception de la base de données — MLD

### 3.1 Tables principales

#### Table `users`
| Champ | Type | Contraintes |
|---|---|---|
| **id** | INTEGER | PRIMARY KEY, AUTO_INCREMENT |
| username | TEXT | UNIQUE, NOT NULL |
| email | TEXT | UNIQUE, NOT NULL |
| password_hash | TEXT | NOT NULL (SHA-512 + salt) |
| salt | TEXT | NOT NULL |
| role | TEXT | DEFAULT 'user' (user / admin) |
| is_banned | INTEGER | DEFAULT 0 |
| created_at | DATETIME | DEFAULT now |

#### Table `categories`
| Champ | Type | Contraintes |
|---|---|---|
| **id** | INTEGER | PRIMARY KEY |
| name | TEXT | UNIQUE, NOT NULL |
| slug | TEXT | UNIQUE, NOT NULL |

#### Table `tags`
| Champ | Type | Contraintes |
|---|---|---|
| **id** | INTEGER | PRIMARY KEY |
| name | TEXT | UNIQUE, NOT NULL |
| color | TEXT | DEFAULT '#7c3aed' |

#### Table `threads`
| Champ | Type | Contraintes |
|---|---|---|
| **id** | INTEGER | PRIMARY KEY |
| title | TEXT | NOT NULL |
| content | TEXT | NOT NULL |
| **author_id** | INTEGER | FK → users(id) |
| **category_id** | INTEGER | FK → categories(id) |
| state | TEXT | ouvert / ferme / archive |
| created_at | DATETIME | DEFAULT now |

#### Table `thread_tags` (jointure)
| Champ | Type |
|---|---|
| **thread_id** | FK → threads(id) |
| **tag_id** | FK → tags(id) |

#### Table `messages`
| Champ | Type | Contraintes |
|---|---|---|
| **id** | INTEGER | PRIMARY KEY |
| **thread_id** | INTEGER | FK → threads(id) |
| **author_id** | INTEGER | FK → users(id) |
| content | TEXT | NOT NULL |
| created_at | DATETIME | DEFAULT now |

#### Table `reactions`
| Champ | Type | Contraintes |
|---|---|---|
| **id** | INTEGER | PRIMARY KEY |
| **user_id** | INTEGER | FK → users(id) |
| **message_id** | INTEGER | FK → messages(id) |
| type | TEXT | like / dislike |
| | | UNIQUE (user_id, message_id) |

### 3.2 Relations principales

| De | Cardinalité | Vers | Sens |
|---|---|---|---|
| `users` | 1 — N | `threads` | Un utilisateur crée plusieurs fils |
| `users` | 1 — N | `messages` | Un utilisateur poste plusieurs messages |
| `users` | 1 — N | `reactions` | Un utilisateur réagit à plusieurs messages |
| `categories` | 1 — N | `threads` | Une catégorie contient plusieurs fils |
| `threads` | 1 — N | `messages` | Un fil contient plusieurs messages |
| `threads` | N — N | `tags` | Via la table de jointure `thread_tags` |
| `messages` | 1 — N | `reactions` | Un message reçoit plusieurs réactions |

---

## 4. Liens utiles

### 4.1 Dépôt Git public

🔗 **https://github.com/Mohand-123/projet-forum**

Le dépôt est public sur GitHub. Tout le code source ainsi que la documentation (MLD, use cases, README) y sont versionnés.

### 4.2 Tableau Trello

🔗 **https://trello.com/b/0C55Nrs2/projet-forum**

Tableau Trello partagé de l'équipe avec les colonnes "À faire", "En cours" et "Terminé".

### 4.3 Répartition globale du travail

| Membre | Domaine principal | Responsabilités |
|---|---|---|
| **AMIR Mohand Arezki** | Back-end / BDD | Conception BDD, routes Go, authentification, sécurité |
| **MATHÉO** | Front-end / UX | Design UI, templates HTML, charte graphique, intégration |
| **ADIB Houssine** | Organisation / Tests / Docs | Suivi Trello, tests manuels, documentation technique |

---

*Document rendu dans le cadre du projet supervisé B1 — Séance 1.*
