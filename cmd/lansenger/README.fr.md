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
| `config` | Gérer les identifiants | `set`, `show`, `clear`, `list-profiles` |
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

# Envoyer dans un groupe avec @all (user_token facultatif, affiché comme bot sans)
lansenger message send-text group123 "Annonce" --group --mention-all

# @mention de personnes spécifiques dans un groupe
lansenger message send-text group123 "Veuillez vérifier" --group --mention staff001

# Diffusion via le canal bot
lansenger message send-bot-message text '{"content":"Avis"}' --chat-id user001 --chat-id user002
```

### Gestion des groupes

```bash
# Créer un groupe
lansenger group create "Groupe Projet" org001 --staff staff001 --staff staff002

# Voir les infos du groupe
lansenger group info group123

# Voir les membres du groupe
lansenger group members group123

# Voir la liste des groupes (bot peut lister ses groupes)
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

# Télécharger un média application/bot (utilisé pour send-text / send-file)
lansenger media upload-app /path/to/file.pdf --media-type file

# Télécharger un média vers un fichier local
lansenger media download-to-file MEDIA_ID --output /path/to/save.pdf
```

## Options globales

| Option | Description |
|------|------|
| `--json` / `-j` | Sortie JSON brute au lieu de tableaux formatés |
| `--profile` / `-P` | Utiliser un profil d'identifiants spécifique (défaut : `default`) |

## Profils multi-applications / multi-bots

Le CLI prend en charge plusieurs profils, chacun correspondant à un appID, avec des identifiants isolés :

```bash
# Configurer la première application (bot personnel)
lansenger config set app_id xxx1 --profile my-bot
lansenger config set app_secret xxx1 --profile my-bot

# Configurer la deuxième application (bot d'organisation)
lansenger config set app_id xxx2 --profile org-bot
lansenger config set app_secret xxx2 --profile org-bot

# Utiliser un profil spécifique
lansenger --profile org-bot staff basic-info STAFF_ID
```

## Sécurité

- Identifiants stockés dans `~/.lansenger/sdk_state.json` avec permissions `0600`
- `config show` masque tous les champs secrets (`***`), seuls `api_gateway_url` et `passport_url` sont affichés en clair
- Variables d'environnement `LANSENGER_APP_ID` / `LANSENGER_APP_SECRET` / `LANSENGER_ENCODING_KEY` / `LANSENGER_CALLBACK_TOKEN` supportées pour CI/CD

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
