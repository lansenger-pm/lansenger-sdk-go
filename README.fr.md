[English](README.md) | [简体中文](README.zhHans.md) | [繁體中文](README.zhHant.md) | [繁體中文（香港）](README.zhHantHK.md) | [Français](README.fr.md)

# lansenger-sdk-go

SDK Go pour la plateforme Lansenger (蓝信) — prend en charge les applications Lansenger, les robots d'organisation et les robots personnels.

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Version: 0.9.18](https://img.shields.io/badge/Version-0.9.18-blue)](https://github.com/lansenger-pm/lansenger-sdk-go)
[![Go 1.26+](https://img.shields.io/badge/Go-1.26%2B-blue)](https://go.dev/)
[![Tests: 148](https://img.shields.io/badge/Tests-148-green)](https://github.com/lansenger-pm/lansenger-sdk-go)

> La bibliothèque SDK n'a aucune dépendance externe (stdlib uniquement). Le CLI (`cmd/lansenger`) utilise [spf13/cobra](https://github.com/spf13/cobra).

## Types de robots pris en charge

| Type de robot | Authentification | WebSocket inbound | Toutes les API |
|----------|------|-------------------|----------|
| **Application Lansenger** | appToken + userToken | ✗ (utilise webhook) | ✓ |
| **Bot d'organisation** | appToken + userToken | ✗ (utilise webhook) | ✓ |
| **Bot personnel** | appToken | ✓ (WebSocket) | ✓ (limité pour les API non-robot) |

Les trois types de robots utilisent le même mécanisme d'authentification : `appToken` est requis pour chaque appel API ; `userToken` n'est nécessaire que pour des opérations spécifiques au niveau utilisateur (infos utilisateur, recherche de personnel, calendrier, etc.).

## Fonctionnalités

- **Client unique** — `LansengerClient` avec `context.Context` pour tous les appels API
- **Persistance des identifiants et tokens** — `CredentialStore` sauvegarde les identifiants et tokens dans un fichier JSON (survit aux redémarrages)
- **Authentification utilisateur OAuth2** — URL d'autorisation, échange de code, renouvellement de token
- **Organisation & départements** — infos organisation, détails/sous-departements/personnel du département
- **Personnel & contacts** — infos basiques/détaillées, mappage d'ID, ancêtres de département, recherche
- **Messagerie** — 3 canaux de chat privé (robot, compte officiel, impersonnation utilisateur) + chat de groupe, tous types de messages, @mention, identité d'émetteur humain/robot
- **Cartes enrichies** — appCard (avec mises à jour dynamiques), oacard, linkCard, appArticles
- **Messages en streaming** — diffusion temps réel basée sur SSE pour agents IA
- **Upload/download de médias** — fichiers, images, vidéos avec détection automatique du type, récupération du chemin de téléchargement, upload app/robot
- **Gestion des messages** — révoquer, mise à jour dynamique de carte, rappel urgent
- **Groupes V2** — créer, infos, membres, liste, vérification de membership, mise à jour des paramètres & membres, dissoudre
- **Calendrier & planification** — calendrier principal, CRUD de planification, gestion des participants, mise à jour des métadonnées des participants
- **Todo unifié** — créer, mettre à jour, supprimer, interroger, gestion des exécutants, comptes de statut
- **Événements de callback** — 24 types d'événements, analyse de données structurées, vérification de signature

## Installation rapide

**SDK (bibliothèque)**:
```bash
go get github.com/lansenger-pm/lansenger-sdk-go
```

**CLI (pour agents IA & débogage)**:
```bash
go install github.com/lansenger-pm/lansenger-sdk-go/cmd/lansenger@latest
lansenger version
```

Le CLI partage les identifiants avec le SDK via `~/.lansenger/sdk_state.json`. Après l'installation, configurez les identifiants :
```bash
lansenger config set app_id YOUR_APP_ID
lansenger config set app_secret YOUR_APP_SECRET
```

## 1. Authentification

### appToken — Requis pour tous les appels API

Chaque méthode du SDK requiert `appToken`. Le client l'obtient et le renouvelle automatiquement en utilisant votre `appID` + `appSecret`. Vous n'avez jamais besoin de gérer appToken manuellement — le `TokenManager` gère le cycle de vie :

1. **Premier appel** → `GET /v1/apptoken/create` avec appID + appSecret → retourne `appToken` (valide 2 heures)
2. **Appels suivants** → réutilisation de l'appToken en cache jusqu'à expiration
3. **Token expiré** → renouvellement automatique via le même endpoint

```go
// appToken est géré automatiquement — configurez simplement appID + appSecret
client := lansenger.NewClient("your-appid", "your-secret")

// Vous pouvez aussi obtenir/invalider le token manuellement
token, err := client.GetToken(ctx)
client.InvalidateToken() // force le renouvellement au prochain appel
```

### userToken — Nécessaire uniquement pour certains endpoints

`userToken` représente l'autorisation d'un utilisateur Lansenger spécifique (obtenu via OAuth2). Il n'est requis que pour :
- Informations au niveau utilisateur (FetchUserInfo, FetchStaffDetail, SearchStaff)
- Opérations de calendrier & planification (FetchPrimaryCalendar, CreateSchedule, etc.)
- Opérations de groupe en tant qu'émetteur humain

### Obtenir les identifiants

| Type de robot | Comment obtenir appID + appSecret |
|----------|----------------------------|
| **Bot personnel** | Client desktop Lansenger → Contacts → Bots intelligents → Bots personnels → cliquer sur l'icône ℹ️ (le client mobile NE montre PAS les identifiants) |
| **Application Lansenger** | Créer sur le Centre développeur Lansenger — peut nécessiter l'approbation de l'administrateur d'organisation |
| **Bot d'organisation** | Créer sur le Centre développeur Lansenger — peut nécessiter l'approbation de l'administrateur d'organisation |

### Authentification utilisateur OAuth2

```go
// Construire l'URL d'autorisation — rediriger l'utilisateur vers le passeport Lansenger
url := client.BuildAuthorizeURL("https://myapp.com/callback", "", "state123")

// Après autorisation de l'utilisateur, échanger le code contre userToken + refreshToken
tokenResult, err := client.ExchangeCode(ctx, "auth_code_from_callback", "https://myapp.com/callback")

// Renouveler un userToken expiré
newToken, err := client.RefreshUserToken(ctx, tokenResult.RefreshToken, "")

// Obtenir les infos utilisateur
userInfo, err := client.FetchUserInfo(ctx, tokenResult.UserToken)
```

## 2. Organisation & Départements

```go
// Informations organisation
org, err := client.FetchOrgInfo(ctx, "orgId", "")

// Hiérarchie des départements
detail, err := client.FetchDepartmentDetail(ctx, "deptId", "", "")
children, err := client.FetchDepartmentChildren(ctx, "deptId", "")
staffs, err := client.FetchDepartmentStaffs(ctx, "deptId", "", 1, 100)
```

## 3. Personnel & Contacts

```go
// Infos basiques du personnel
staff, err := client.FetchStaffBasicInfo(ctx, "staffOpenId", "")

// Profil détaillé (userToken recommandé)
detail, err := client.FetchStaffDetail(ctx, "staffOpenId", "ut")

// Mappage téléphone → staffId
mapping, err := client.FetchStaffIdMapping(ctx, "orgId", "mobile", "13800138000", "")

// Ancêtres de département pour un membre du personnel
ancestors, err := client.FetchDepartmentAncestors(ctx, "staffOpenId", "")

// Recherche de personnel (requiert userToken ou userID)
results, err := client.SearchStaff(ctx, "Zhang San", "ut", "", true, nil, 1, 10)

// IDs de champs extra organisation
fields, err := client.FetchOrgExtraFieldIDs(ctx, "orgId", "", 1, 1000)
```

## 4. Messagerie & Médias

#### Chat privé robot — le plus courant

```go
result, err := client.SendText(ctx, "staff123", "Bonjour !", "", 0, "", false, nil, false, "", "")
result, err := client.SendMarkdown(ctx, "staff123", "**Gras**", false, nil, false, "", "")
result, err := client.SendFile(ctx, "staff123", "/path/to/report.pdf", "", 0, "", false, "", "")
```

#### Canal compte officiel

```go
result, err := client.SendAccountMessage(ctx, "text",
    map[string]interface{}{"content": "Notice système"},
    []string{"staff1", "staff2"}, nil, "524288-xxxx", "", "", "")
```

#### Canal impersonnation utilisateur (requiert userToken)

```go
result, err := client.SendUserMessage(ctx, "staff456", "text",
    map[string]interface{}{"content": "Bonjour"}, "ut", "")
```

#### Chat de groupe

```go
// Bot → groupe
result, err := client.SendText(ctx, "group123", "Annonce", "", 0, "", false, nil, true, "", "")

// Humain → groupe (avec userToken)
result, err := client.SendGroupMessage(ctx, "group123", "text",
    map[string]interface{}{"content": "Je m'en charge"}, "ut", "", false, nil, "", "", "")

// @mention dans un groupe
result, err := client.SendText(ctx, "group123", "Important !", "", 0, "", true, nil, true, "", "")
```

#### Cartes enrichies

```go
// appCard
params := &lansenger.AppCardParams{
    ChatID: "staff123", BodyTitle: "Approbation", IsDynamic: true,
}
result, err := client.SendAppCardWithParams(ctx, params)

// linkCard
params := &lansenger.LinkCardParams{
    ChatID: "staff123", Title: "Article", Link: "https://...",
}
result, err := client.SendLinkCardWithParams(ctx, params)

// Mettre à jour le statut d'une carte dynamique
updateParams := &lansenger.DynamicCardUpdateParams{
    MsgID: "msg123", IsLastUpdate: true,
}
result, err := client.UpdateDynamicCard(ctx, updateParams)
```

#### Messages en streaming (pour agents IA)

```go
result, err := client.CreateStreamMessage(ctx, "staff1", "staff", "stream1")
result, err := client.FetchStreamMessage(ctx, "msg123")
```

#### Médias

```go
// Upload (service principal — type numérique)
upload, err := client.UploadAppMedia(ctx, "/path/to/file.pdf", lansenger.AppMediaTypeFile, 0, 0, 0)

// Upload (app/robot — type chaîne, supporte width/height/duration)
upload, err := client.UploadAppMedia(ctx, "/path/to/video.mp4",
    lansenger.AppMediaTypeVideo, 680, 480, 300)

// Download
download, err := client.DownloadMedia(ctx, "media123")

// Télécharger et sauvegarder dans un fichier
path, err := client.DownloadMediaToFile(ctx, "media123", "/path/to/save.pdf")

// Récupérer les infos du chemin de téléchargement
pathInfo, err := client.FetchMediaPath(ctx, "media123", "ut")

// Révoquer des messages
result, err := client.RevokeMessage(ctx, []string{"msg1", "msg2"}, "robot", "")

// Envoyer un rappel urgent
result, err := client.SendReminder(ctx, "msg123", []int{1, 2}, []string{"staff1", "staff2"})
```

## 5. Groupes

```go
// Créer un groupe
info := &lansenger.GroupCreateInfo{
    Name: "Chat Projet", OrgID: 1, StaffIDList: []string{"s1", "s2", "s3"},
}
group, err := client.CreateGroup(ctx, info, "")

// Obtenir infos & membres
info, err := client.FetchGroupInfo(ctx, "groupOpenId", "")
members, err := client.FetchGroupMembers(ctx, "groupOpenId", "", 0, 100)
groups, err := client.FetchGroupList(ctx, "", 0, 100)

// Vérifier le membership
result, err := client.CheckIsInGroup(ctx, "groupOpenId", "", "staff1")

// Mettre à jour les paramètres
result, err := client.UpdateGroupInfo(ctx, "groupId", map[string]interface{}{"name": "Nouveau Nom"}, "")

// Ajouter/supprimer des membres
result, err := client.UpdateGroupMembers(ctx, "groupId",
    []string{"staff4"}, []string{"staff3"}, nil, "")

// Dissoudre un groupe
result, err := client.DissolveGroup(ctx, "groupId", "ut")
```

## 6. Calendrier & Planification

```go
// Obtenir le calendrier principal (requiert userToken ou userID)
cal, err := client.FetchPrimaryCalendar(ctx, "ut", "uid1")

// Créer une planification (startTime/endTime sont des objets map, allDay est "yes"/"no")
schedule, err := client.CreateSchedule(ctx, cal.CalendarID, "Réunion d'équipe",
    map[string]interface{}{"time": "2024-01-15T09:00"},
    map[string]interface{}{"time": "2024-01-15T10:00"},
    nil, "", "no", "", nil, "", "", "", "ut", "")

// Obtenir/supprimer une planification
info, err := client.FetchSchedule(ctx, "cal1", "sch1", "ut", "")
result, err := client.DeleteSchedule(ctx, "cal1", "sch1", "", "", "", "ut", "")

// Mettre à jour une planification
result, err := client.UpdateSchedule(ctx, "cal1", "sch1",
    map[string]interface{}{"summary": "Réunion mise à jour"}, "ut", "")

// Liste des planifications dans un intervalle de temps
schedules, err := client.FetchScheduleList(ctx, "cal1",
    map[string]interface{}{"time": "2024-01-15T00:00"},
    map[string]interface{}{"time": "2024-01-15T23:59"}, "ut", "")

// Gestion des participants (participants sont []string)
attendees, err := client.FetchScheduleAttendees(ctx, "cal1", "sch1", 1, 10, "ut", "")
result, err := client.AddScheduleAttendees(ctx, "cal1", "sch1",
    []string{"staff2"}, "", "", "", "ut", "")
result, err := client.DeleteScheduleAttendees(ctx, "cal1", "sch1",
    []string{"staff2"}, "", "", "", "ut", "")

// Mettre à jour les métadonnées des participants
result, err := client.UpdateScheduleAttendeeMeta(ctx, "cal1", "sch1",
    map[string]interface{}{"rsvpStatus": "accepted"}, "ut", "")
```

## 7. Todo unifié

```go
// Créer une tâche todo
todo, err := client.CreateTodoTask(ctx, "Demande d'approbation", lansenger.TodoTypeApproval,
    "https://app.com/a/1", "https://pc.app.com/a/1", []string{"staff1"}, "org1", "", "", "", "")

// Mettre à jour le statut (11=en attente de lecture, 12=lu, 21=en attente de faire, 22=terminé)
result, err := client.UpdateTodoTaskStatus(ctx, "taskId", lansenger.TodoStatusDone, "org1", "", "")

// Mettre à jour le contenu
result, err := client.UpdateTodoTask(ctx, "taskId", "Mis à jour", "l", "p", "org1", "", "")

// Supprimer (émetteur uniquement)
result, err := client.DeleteTodoTask(ctx, "taskId", "org1", "", "")

// Interroger
list, err := client.FetchTodoTaskList(ctx, "org1", nil, "", nil, "")
task, err := client.FetchTodoTaskByID(ctx, "taskId", "org1", "", "")
task, err := client.FetchTodoTaskBySourceID(ctx, "src1", "org1", "", "")
counts, err := client.FetchTodoTaskStatusCounts(ctx, "staff1", "org1", "", "", "")

// Gestion des exécutants
result, err := client.AddExecutors(ctx, []string{"staff2"}, "org1", "taskId", "")
result, err := client.DeleteExecutors(ctx, []string{"staff2"}, "org1", "taskId", "")
executors, err := client.FetchExecutorList(ctx, "taskId", "org1", "", nil, "")
```

## 8. Événements de callback

```go
// Analyser le payload webhook en clair (non chiffré) — chaîne de requête URL ou JSON
events, err := lansenger.ParseCallbackPayload("eventType=staff_modify&staffId=s001&orgId=org1")

// Analyser le callback JSON en clair
events, err = lansenger.ParseCallbackPayload(`{"events":[{"eventType":"staff_modify","data":{"staffId":"s001"}}],"orgId":"org1","appId":"app1"}`)

// Déchiffrer le payload de callback chiffré (AES-256-CBC)
result, err := lansenger.DecryptCallbackPayload(encryptedData, encodingKey, knownAppID)
fmt.Println(result.OrgID, result.AppID, result.Events)

// Vérifier la signature (SHA1, conforme au protocole Lansenger)
valid := lansenger.VerifyCallbackSignature(timestamp, nonce, signature, encodingKey, dataEncrypt, callbackToken)

// Types d'événements disponibles (24 types, mappage de champs structuré)
types := lansenger.GetCallbackEventTypes()
```

## 9. Lecture de chats

```go
// Obtenir la liste de chats de l'utilisateur (privé + groupe)
chats, err := client.FetchChatList(ctx, "ut", "private", "", "", "")

// Obtenir les messages de chat privé avec une personne spécifique
msgs, err := client.FetchChatMessages(ctx, "ut", 10, "", "s001", "", "", "", "")

// Obtenir les messages de chat de groupe
msgs, err := client.FetchChatMessages(ctx, "ut", 10, "", "", "g001", "", "", "")
```

## Matrice de capacités des types de messages

| msgType | Markdown | @mention | Attachements | Canaux privés | Chat de groupe | Remarques |
|---------|----------|----------|-------------|------------------|------------|-------|
| `text` | ✗ | ✓ (groupe) | ✓ | Bot, Compte officiel, Impersonnation utilisateur | ✓ | Maximum 6000 octets |
| `formatText` | ✓ | ✗ | ✗ | Impersonnation utilisateur uniquement | ✓ | Markdown via formatType=1 |
| `oacard` | ✗ | ✗ | ✗ | Bot, Compte officiel, Impersonnation utilisateur | ✓ | Carte simple avec champs |
| `appCard` | ✓ (tags div) | ✗ | ✗ | Bot, Compte officiel, Impersonnation utilisateur | ✓ | Carte enrichie, mises à jour dynamiques |
| `linkCard` | ✗ | ✗ | ✗ | Bot, Compte officiel | ✓ | Carte de prévisualisation de lien |
| `appArticles` | ✗ | ✗ | ✗ | Bot privé uniquement | ✓ | Liste d'articles (1+ articles) |

**Chat de groupe** prend en charge tous les types de messages. Seul le chat de groupe prend en charge @mention.

## Configuration

### Vue d'ensemble des identifiants

Tous les identifiants sont persistés par profil dans `~/.lansenger/sdk_state.json` (permissions 0600) :

| Identifiant | Requis | Clé CLI | Description |
|-------------|--------|---------|-------------|
| App ID | ✓ | `app_id` | ID application/robot Lansenger |
| App Secret | ✓ | `app_secret` | Secret application/robot Lansenger |
| API Gateway URL | ✓ | `api_gateway_url` | Point d'accès API (défaut : `https://open.e.lanxin.cn/open/apigw`) |
| Passport URL | OAuth2 uniquement | `passport_url` | URL page d'autorisation OAuth2 |
| Redirect URI | OAuth2 uniquement | `redirect_uri` | URI de redirection OAuth2 (défaut : `http://localhost:8765`) |
| Encoding Key | Callbacks uniquement | `encoding_key` | Clé de déchiffrement AES-256-CBC |
| Callback Token | Callbacks uniquement | `callback_token` | Token de vérification de signature callback |

### Configuration CLI

```bash
# Étape 1 : Configurer les identifiants requis
lansenger config set app_id YOUR_APP_ID
lansenger config set app_secret YOUR_APP_SECRET
lansenger config set api_gateway_url https://open.e.lanxin.cn/open/apigw

# Étape 2 (optionnel) : URL Passport pour OAuth2 (nécessaire pour userToken)
lansenger config set passport_url https://passport.lx.qianxin.com
lansenger config set redirect_uri http://localhost:8765   # URI de redirection OAuth2 (défaut)

# Étape 3 (optionnel) : Identifiants callback (nécessaire pour Webhook)
lansenger config set encoding_key YOUR_ENCODING_KEY
lansenger config set callback_token YOUR_CALLBACK_TOKEN

# Vérifier la configuration
lansenger config show

# Support multi-profil (ex. organisations/applications séparées)
lansenger config set app_id APP2_ID --profile org2
lansenger config set app_secret APP2_SECRET --profile org2
lansenger --profile org2 staff basic-info STAFF_ID
```

### Configuration SDK

**Depuis le code** (direct) :
```go
client := lansenger.NewClient("app_id", "app_secret")
// URL gateway personnalisée si nécessaire
cfg := lansenger.NewConfig("app_id", "app_secret")
cfg.APIGatewayURL = "https://custom-gateway.example.com"
cfg.PassportURL = "https://passport.example.com"
cfg.RedirectURI = "https://myapp.com/callback"
cfg.EncodingKey = "your_encoding_key"
cfg.CallbackToken = "your_callback_token"
client := lansenger.NewClientWithConfig(cfg)
```

**Depuis l'environnement** (auto-détection) :

| Variable | Requis | Description | Défaut |
|----------|--------|-------------|--------|
| `LANSENGER_APP_ID` | ✓ | ID App/Bot | — |
| `LANSENGER_APP_SECRET` | ✓ | Secret App/Bot | — |
| `LANSENGER_API_GATEWAY_URL` | ✗ | URL Gateway API | `https://open.e.lanxin.cn/open/apigw` |
| `LANSENGER_PASSPORT_URL` | ✗ | URL Passport (OAuth2) | — |
| `LANSENGER_REDIRECT_URI` | ✗ | URI de redirection OAuth2 | `http://localhost:8765` |
| `LANSENGER_ENCODING_KEY` | ✗ | Clé de déchiffrement callback | — |
| `LANSENGER_CALLBACK_TOKEN` | ✗ | Token callback (défaut = encoding_key) | — |
| `LANSENGER_HTTP_TIMEOUT` | ✗ | Timeout HTTP (secondes) | `30` |

```go
client, err := lansenger.NewClientFromEnv()
```

### Persistance des identifiants et tokens

Par défaut, les identifiants et tokens restent uniquement en mémoire (perdus à la fin du processus). Activez la persistance fichier avec `CredentialStore` :

```go
// Auto-persist vers ~/.lansenger/sdk_state.json (permissions 0600)
store := lansenger.NewCredentialStore("", "default")
store.SaveCredentials("app_id", "app_secret", "https://apigw.lx.qianxin.com", "https://passport.lx.qianxin.com")
store.SaveCallbackConfig("encoding_key", "callback_token")

// Sauvegarder les tokens
store.SaveAppToken("token123", 7200)
store.SaveUserToken("ut123", "rt123", 7200)

// Charger les tokens (retourne chaîne vide si expiré)
token, err := store.LoadAppToken()

// Identifiants partagés avec le SDK Python (même format ~/.lansenger/sdk_state.json)
```

Avec la persistance activée :
- **appToken** peut être sauvegardé et restauré au redémarrage (évite les appels API redondants)
- **userToken + refreshToken** peuvent être sauvegardés après l'échange OAuth2
- **Identifiants + URLs** sont sauvegardés ensemble pour une récupération complète de la configuration

## Identité et permissions

### Matrice de capacités d'identité

La plateforme Lansenger propose trois types d'identité avec différents accès API :

| Domaine de commande | Bot personnel | App d'organisation (auto-hébergée) | App d'organisation + Bot | Remarques |
|--------|:---:|:---:|:---:|------|
| `message send-text/markdown/file/...` (bot DM) | **Y** | N | **Y** | Seuls les bots peuvent envoyer des DM bot |
| `message send-text --group` (chat de groupe) | **Y** | N | **Y** | Le robot personnel prend désormais en charge la messagerie de groupe |
| `message send-group-message` | **Y** | N | **Y** | Idem ci-dessus |
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

## Structure du projet

```
lansenger-sdk-go/
├── client.go            # LansengerClient — client principal avec helpers HTTP
├── config.go            # Config — configuration + variables d'environnement
├── constants.go         # Endpoints API, types de médias, types d'événements de callback
├── errors.go            # Hiérarchie LansengerError (Auth/Config/API/Network/File)
├── models.go            # 50+ types de structs résultat/params
├── auth.go              # TokenManager — cycle de vie appToken avec renouvellement auto
├── user_token_manager.go # UserTokenManager — cycle de vie userToken avec renouvellement auto
├── url_helpers.go       # BuildAPIURL — pattern Options pour construction d'URL
├── oauth.go             # OAuth2 URL d'autorisation, échange de code, renouvellement de token
├── contacts.go          # API Personnel & infos organisation
├── users.go             # API Profil utilisateur
├── departments.go       # API Départements
├── groups.go            # API Groupes V2
├── chats.go             # API Liste de chats & messages
├── account_messages.go  # Canal compte officiel (4.6.1)
├── user_messages.go     # Canal impersonnation utilisateur (4.6.3)
├── group_messages.go    # Canal chat de groupe (4.6.2)
├── robot_messages.go      # Canal robot (4.6.12)
├── messaging.go         # Méthodes de convenance + révoquer + mise à jour dynamique
├── streaming.go         # Messages en streaming SSE
├── media.go             # Upload/download de fichiers & images
├── todos.go             # Todo unifié (4.33) — 12 endpoints
├── calendars.go         # Calendrier & planification (4.23) — 10 endpoints
├── callbacks.go         # Analyse d'événements de callback + déchiffrement AES-256-CBC + vérification de signature SHA1
├── persistence.go       # CredentialStore — persistance dans fichier JSON
├── version.go           # Constante de version SDK
├── *_test.go            # 136 tests unitaires + 10 tests d'intégration
├── cmd/lansenger/       # Outil CLI (config, oauth, messaging, staff, etc.)
├── go.mod
└── README.md
```

## Développement

```bash
go test . -v                        # tous les tests (136 unitaires + 10 intégration)
go test . -run TestIntegration      # tests d'intégration uniquement (requiert ~/.lansenger/sdk_state.json + réseau)
```

## Licence

MIT — voir [LICENSE](LICENSE).