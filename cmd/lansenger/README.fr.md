[English](README.md) | [简体中文](README.zhHans.md) | [繁體中文](README.zhHant.md) | [繁體中文（香港）](README.zhHantHK.md) | [Français](README.fr.md)

# Lansenger CLI (Go)

Outil en ligne de commande Lansenger — interagissez avec les API Lansenger directement depuis le terminal : envoyez des messages, gérez des groupes, interrogez le personnel/les départements, gérez les calendriers et les tâches, et plus encore.

La syntaxe des commandes est identique aux versions Python et TypeScript. Installez n'importe laquelle.

## Installation

```bash
go install github.com/lansenger-pm/lansenger-sdk-go/cmd/lansenger@latest
```

Ou compiler depuis les sources :

```bash
git clone https://github.com/lansenger-pm/lansenger-sdk-go.git
cd lansenger-sdk-go/cmd/lansenger
go build -o lansenger .
```

Nécessite Go 1.26+.

## Démarrage rapide

### 1. Configurer les identifiants

Sauvegardez les identifiants via `config set` (stockés par profil dans `~/.lansenger/sdk_state.json`, clés masquées, permissions fichier 0600) :

**Identifiants requis** :

```bash
lansenger config set app_id YOUR_APP_ID
lansenger config set app_secret YOUR_APP_SECRET
lansenger config set api_gateway_url https://open.e.lanxin.cn/open/apigw
```

**Authentification OAuth2 (remplissez si vous avez besoin d'un userToken)** :

```bash
lansenger config set passport_url https://passport.lx.qianxin.com
lansenger config set redirect_uri http://localhost:8765   # URI de redirection OAuth2 (défaut)
```

**Réception des callbacks (remplissez si vous devez analyser/vérifier les webhooks)** :

```bash
lansenger config set encoding_key YOUR_ENCODING_KEY
lansenger config set callback_token YOUR_CALLBACK_TOKEN
```

Vous pouvez également configurer via les variables d'environnement (compatible CI/CD) :

```bash
export LANSENGER_APP_ID=YOUR_APP_ID
export LANSENGER_APP_SECRET=YOUR_APP_SECRET
export LANSENGER_REDIRECT_URI=http://localhost:8765
```

### 2. Voir la configuration

```bash
lansenger config show
```

### 3. Envoyez votre premier message

```bash
lansenger message send-text staff001 "Hello from CLI!"
```

## Aperçu des commandes

| Groupe | Description | Sous-commandes |
|--------|------|--------|
| `config` | Gérer les identifiants | `set`, `show`, `clear`, `list-profiles`, `delete-profile`, `list-users` |
| `message` | Envoyer et gérer les messages | `send-text`, `send-markdown`, `send-file`, `send-image-url`, `send-link-card`, `send-app-articles`, `send-app-card`, `send-oacard`, `send-bot-message`, `send-group-message`, `send-account-message`, `send-user-message`, `update-dynamic-card`, `revoke`, `query-groups`, `send-reminder` |
| `group` | Gérer les groupes | `create`, `info`, `members`, `list`, `check`, `update`, `update-members`, `dismiss` |
| `staff` | Interroger les infos du personnel | `basic-info`, `detail`, `ancestors`, `id-mapping`, `org-extra-fields`, `search`, `org-info` |
| `department` | Interroger les départements | `detail`, `children`, `staffs` |
| `calendar` | Calendrier et planification | `primary`, `create-schedule`, `fetch-schedule`, `delete-schedule`, `list-schedules`, `attendees`, `add-attendees`, `delete-attendees`, `update-schedule`, `attendee-meta` |
| `todo` | Gestion des tâches | `create`, `update`, `update-status`, `delete`, `list`, `fetch-by-id`, `fetch-by-source`, `status-counts`, `executor-status`, `add-executors`, `delete-executors`, `executor-list` |
| `oauth` | Authentification OAuth2 | `authorize-url`, `exchange-code`, `refresh-token`, `user-info`, `parse-callback`, `validate-state` |
| `callback` | Analyse des événements callback | `parse-payload`, `decrypt-payload`, `verify-signature`, `event-types` |
| `media` | Opérations sur les fichiers média | `upload`, `upload-app`, `download`, `download-to-file`, `path` |
| `streaming` | Messages en streaming (IA) | `create`, `fetch` |
| `chat` | Conversations et messages | `list`, `messages` |
| `health` | Vérification de connexion | `check` |

## Exemples courants

### Messagerie

```bash
# Envoyer un message texte
lansenger message send-text chat123 "Bonjour !"

# Envoyer un message Markdown
lansenger message send-markdown chat123 "**Gras** *italique*"

# Envoyer un fichier
lansenger message send-file chat123 /path/to/report.pdf

# Envoyer une image depuis une URL
lansenger message send-image-url chat123 https://example.com/photo.jpg

# Envoyer une carte lien
lansenger message send-link-card chat123 "Documentation" "Lire ceci" https://docs.example.com

# Envoyer une carte applicative
lansenger message send-app-card chat123 "Titre de la carte" --content "Texte" --card-link https://example.com

# Envoyer plusieurs articles
lansenger message send-app-articles chat123 '{"title":"Article 1","url":"https://a.com"}' '{"title":"Article 2","url":"https://b.com"}'

# Envoyer une carte d'approbation OA
lansenger message send-oacard chat123 "Titre approbation" --head "Notification" --field '{"key":"Demandeur","value":"Jean"}'

# Envoyer dans un groupe avec @all (user_token facultatif, affiché comme robot sans)
lansenger message send-text group123 "Annonce" --group --mention-all

# @mention de personnes spécifiques dans un groupe
lansenger message send-text group123 "Veuillez vérifier" --group --mention staff001

# @mention de bots spécifiques dans le groupe
lansenger message send-text group123 "Bot check" --group --mention-bot bot001 --mention-bot bot002

# Répondre à un message (référence de message)
lansenger message send-text group123 "Got it" --group --ref-msg-id 524288-xxx

# Diffusion via le canal robot
lansenger message send-bot-message text '{"content":"Avis"}' --chat-id user001 --chat-id user002

# Réponse du canal bot (référence de message)
lansenger message send-bot-message text '{"content":"Reply"}' --chat-id user001 --ref-msg-id 524288-xxx

# Rechercher la liste des identifiants de groupe
lansenger message query-groups --page 0 --size 100
```

### Gestion des groupes

```bash
# Créer un groupe
lansenger group create "Groupe Projet" org001 --staff staff001 --staff staff002

# Voir les infos du groupe
lansenger group info group123

# Voir les membres du groupe
lansenger group members group123

# Voir la liste des groupes (robot peut lister ses groupes)
lansenger group list

# Voir la liste des groupes en tant qu'utilisateur (nécessite user_token)
lansenger group list --user-token YOUR_USER_TOKEN

# Vérifier l'appartenance au groupe
lansenger group check group123 --staff-id staff001

# Mettre à jour les infos du groupe
lansenger group update group123 --name "Nouveau nom" --desc "Nouvelle description"

# Ajouter/supprimer des membres
lansenger group update-members group123 --add staff003 --remove staff001
```

### Interrogation du personnel

```bash
# Infos de base du personnel
lansenger staff basic-info staff001

# Infos détaillées du personnel
lansenger staff detail staff001

# Rechercher du personnel
lansenger staff search "Zhang San" --user-token YOUR_USER_TOKEN

# Mapping d'ID (téléphone → staffId)
lansenger staff id-mapping org001 mobile 13800138000

# Ancêtres du département
lansenger staff ancestors staff001
```

### Fichiers média

```bash
# Télécharger un fichier plateforme principale
lansenger media upload /path/to/file.pdf --media-type 3

# Télécharger un média application/robot (utilisé pour send-text / send-file)
lansenger media upload-app /path/to/file.pdf --media-type file

# Télécharger un média vers un fichier local
lansenger media download-to-file MEDIA_ID --output /path/to/save.pdf
```

## Options globales

| Option | Description |
|------|------|
| `--json` / `-j` | Sortie JSON brute au lieu de tableaux formatés |
| `--profile` / `-P` | Utiliser un profil d'identifiants spécifique (défaut : `default`) |
| `--as <staff_id>` | Charge et rafraîchit automatiquement le jeton utilisateur pour le staff_id spécifié depuis le stockage des identifiants |

## Profils multi-applications / multi-robots

Le CLI prend en charge plusieurs profils, chacun correspondant à un appID, avec des identifiants isolés :

```bash
# Configurer la première application (robot personnel)
lansenger config set app_id xxx1 --profile my-bot
lansenger config set app_secret xxx1 --profile my-bot

# Configurer la deuxième application (robot d'organisation)
lansenger config set app_id xxx2 --profile org-bot
lansenger config set app_secret xxx2 --profile org-bot

# Supprimer un profil (bascule automatiquement vers default si actif)
lansenger config delete-profile my-bot

# Utiliser un profil spécifique
lansenger --profile org-bot staff basic-info STAFF_ID
```

## Sécurité

- Identifiants stockés dans `~/.lansenger/sdk_state.json` avec permissions `0600`
- `config show` masque tous les champs secrets (`***`), seuls `api_gateway_url` et `passport_url` sont affichés en clair
- Variables d'environnement `LANSENGER_APP_ID` / `LANSENGER_APP_SECRET` / `LANSENGER_ENCODING_KEY` / `LANSENGER_CALLBACK_TOKEN` supportées pour CI/CD

## Identité et permissions

### Matrice de capacités d'identité

La plateforme Lansenger propose trois types d'identité avec différents accès API :

| Domaine de commande | Bot personnel | App d'organisation (auto-hébergée) | App d'organisation + Bot | Remarques |
|--------|:---:|:---:|:---:|------|
| `message send-text/markdown/file/...` (bot DM) | **Y** | N | **Y** | Seuls les bots peuvent envoyer des DM bot |
| `message send-text --group` (chat de groupe) | N* | N | **Y** | L'API bot personnel le supporte mais pas encore de fonction rejoindre-groupe |
| `message send-group-message` | N* | N | **Y** | Idem ci-dessus |
| `message send-account-message` (compte officiel) | N | **Y** | **Y** | Nécessite la capacité compte officiel |
| `message send-user-message` (utilisateur à utilisateur) | N | **Y** | **Y** | Nécessite userToken + OAuth2 |
| `message revoke` | **Y** | **Y** | **Y** | Révoquer ses propres messages |
| `staff *` (contacts lecture seule) | N | **Y** | **Y** | `search` nécessite en plus userToken |
| `department *` | N | **Y** | **Y** | Applications niveau organisation uniquement |
| `calendar *` | N | **Y** | **Y** | Avec userToken = identité utilisateur ; sans = identité bot |
| `todo *` | N | **Y** | **Y** | Applications niveau organisation uniquement |
| `chat list/messages` | N | **Y** | **Y** | Applications niveau organisation uniquement |
| `group *` (gestion de groupes V2) | N | N | **Y** | Nécessite que le bot soit dans le groupe |
| `media upload` | **Y** | **Y** | **Y** | Upload général |
| `media upload-app` | **Y** | **Y** | **Y** | Apps auto-hébergées uniquement (pas ISV) |
| `media download/path` | **Y** | **Y** | **Y** | Téléchargement général |
| `oauth *` | N | **Y** | **Y** | Applications niveau organisation uniquement |
| `streaming *` | N | **Y** | **Y** | Applications niveau organisation uniquement |
| `callback *` (analyse d'événements) | N/A | N/A | N/A | Opération pure de données, aucune identité requise |

> \* **N\*** = La capacité API existe, mais la fonction rejoindre-groupe n'est pas encore disponible.

> **Bot personnel** peut uniquement envoyer/recevoir des messages et uploader/télécharger des fichiers. Impossible d'accéder aux contacts, groupes, calendriers ou OAuth2.
>
> **App d'organisation vs App d'organisation + Bot** : Même appID/appSecret. La seule différence concerne les canaux de messagerie — seuls les bots peuvent envoyer des DM bot et des messages de groupe (car seuls les bots peuvent rejoindre les groupes). Toutes les autres API (contacts, calendrier, todo, chat, OAuth2, streaming) fonctionnent de manière identique pour les deux. Actuellement, seules les apps auto-hébergées supportent la capacité bot.

### Permissions du Centre développeur

Au-delà du type d'identité, certains appels API dépendent également des permissions activées dans le Centre développeur Lansenger. L'organisation peut restreindre l'accès développeur, nécessitant l'assistance de l'administrateur.

**Permissions de base (activées par défaut) :**

| Permission | Description |
|------|------|
| Get basic user info | Obtenir les infos de base du personnel pour la connexion système/app |
| Send notification messages | Obtenir les canaux de message d'organisation pour envoyer des messages aux personnes/groupes |

**Permissions avancées (désactivées par défaut, doivent être activées manuellement) :**

| Permission | Description |
|------|------|
| Contacts read-only | Accès en lecture aux contacts |
| Contacts edit | Accès en édition aux contacts (créer/mettre à jour/supprimer du personnel) |
| Sensitive info - Phone | Accéder aux numéros de téléphone des utilisateurs |
| Sensitive info - Email | Accéder aux emails des utilisateurs |
| Sensitive info - ID number | Accéder aux numéros d'identité des utilisateurs |
| Sensitive info - Employee ID | Accéder aux identifiants employé des utilisateurs |
| Map unique attribute to staff ID | Mapper téléphone/email/ID employé vers l'ID personnel |
| App edit | Créer et mettre à jour des applications |
| Groups read-only | Accès en lecture aux groupes |
| Groups edit | Accès en édition aux groupes |
| Calendar read-only | Accès en lecture au calendrier et aux planifications |
| Calendar edit | Accès en édition au calendrier et aux planifications |
| Upload media | Permission d'upload de fichiers média |
| Workbench template read | Accès en lecture aux modèles de workbench |
| Workbench template write | Accès en écriture aux modèles de workbench |

En cas d'erreur de permission, vérifiez d'abord que le type d'identité supporte l'opération, puis invitez l'utilisateur à activer la permission avancée correspondante dans le Centre développeur (contactez l'administrateur d'organisation si l'accès est impossible).

## Compatibilité CLI

Ce CLI Go partage la même syntaxe de commande que les versions Python et TypeScript :

```bash
# Python CLI
pip install lansenger-cli

# Go CLI
go install github.com/lansenger-pm/lansenger-sdk-go/cmd/lansenger@latest

# TypeScript CLI
npm install -g lansenger-cli
```

## Licence

Licence MIT
