# RideNow - Microservices

Application de microservices pour la gestion de courses (RideNow) avec des services pour les utilisateurs (chauffeurs et passagers) et les courses.

## Architecture

L'application est composée de plusieurs microservices :

- **Users Service** : Gestion des chauffeurs et des passagers
- **Rides Service** : Gestion des courses

Chaque service possède sa propre base de données MongoDB et communique avec les autres services via HTTP.

## Démarrage rapide

### Prérequis

- Docker et Docker Compose
- Go 1.24+ (pour le développement local)

### Lancement avec Docker Compose

```bash
docker compose up -d
```

Les services seront disponibles sur :

- **Users Service** : `http://localhost:3000`
- **Rides Service** : `http://localhost:8080`
- **Users Database** : `localhost:27019`
- **Rides Database** : `localhost:27020`

---

## Service Users

Ce service gère les chauffeurs et les passagers pour l'application RideNow.

### URL de base

- Docker : `http://localhost:3000`
- Dans le réseau Docker : `http://users-service:3000`

### Endpoints Chauffeurs

#### Créer un chauffeur

Crée un nouveau chauffeur. Le chauffeur sera défini comme disponible par défaut.

```bash
curl -X POST http://localhost:3000/drivers \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Jean Dupont"
  }'
```

#### Obtenir tous les chauffeurs

Récupère tous les chauffeurs, optionnellement filtrés par disponibilité.

```bash
# Obtenir tous les chauffeurs
curl -X GET http://localhost:3000/drivers

# Obtenir uniquement les chauffeurs disponibles
curl -X GET "http://localhost:3000/drivers?available=true"
```

#### Mettre à jour le statut d'un chauffeur

Met à jour le statut de disponibilité d'un chauffeur.

```bash
curl -X PATCH http://localhost:3000/drivers/{driver_id}/status \
  -H "Content-Type: application/json" \
  -d '{
    "is_available": false
  }'
```

**Exemple :**

```bash
curl -X PATCH http://localhost:3000/drivers/507f1f77bcf86cd799439011/status \
  -H "Content-Type: application/json" \
  -d '{
    "is_available": false
  }'
```

### Endpoints Passagers

#### Créer un passager

Crée un nouveau passager.

```bash
curl -X POST http://localhost:3000/passengers \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Alice Martin"
  }'
```

#### Obtenir tous les passagers

Récupère tous les passagers.

```bash
curl -X GET http://localhost:3000/passengers
```

#### Obtenir un passager par ID

Récupère un passager spécifique par son ID.

```bash
curl -X GET http://localhost:3000/passengers/{passenger_id}
```

**Exemple :**

```bash
curl -X GET http://localhost:3000/passengers/507f1f77bcf86cd799439011
```

#### Mettre à jour un passager

Met à jour les informations d'un passager.

```bash
curl -X PUT http://localhost:3000/passengers/{passenger_id} \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Alice Martin"
  }'
```

**Exemple :**

```bash
curl -X PUT http://localhost:3000/passengers/507f1f77bcf86cd799439011 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Alice Martin"
  }'
```

#### Supprimer un passager

Supprime un passager par son ID.

```bash
curl -X DELETE http://localhost:3000/passengers/{passenger_id}
```

**Exemple :**

```bash
curl -X DELETE http://localhost:3000/passengers/507f1f77bcf86cd799439011
```

### Exemples de réponses

#### Réponse Chauffeur

```json
{
  "id": "507f1f77bcf86cd799439011",
  "name": "Jean Dupont",
  "is_available": true
}
```

#### Réponse Passager

```json
{
  "id": "507f1f77bcf86cd799439011",
  "name": "Alice Martin",
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:00Z"
}
```

---

## Service Rides

Ce service gère les courses pour l'application RideNow. Il communique avec le service Users pour obtenir les chauffeurs disponibles et mettre à jour leur statut.

### URL de base

- Docker : `http://localhost:8080`
- Dans le réseau Docker : `http://rides-service:8080`

### Endpoints Courses

#### Créer une course

Crée une nouvelle course. Le service sélectionne automatiquement un chauffeur disponible et le marque comme indisponible.

```bash
curl -X POST http://localhost:8080/rides \
  -H "Content-Type: application/json" \
  -d '{
    "passengerId": "passenger-001",
    "from_zone": "Centre-ville",
    "to_zone": "Aéroport"
  }'
```

**Réponse :**

```json
{
  "id": "507f1f77bcf86cd799439011",
  "passengerId": "passenger-001",
  "driverId": "507f1f77bcf86cd799439012",
  "from_zone": "Centre-ville",
  "to_zone": "Aéroport",
  "price": 25.5,
  "status": "ASSIGNED",
  "paymentStatus": "PENDING",
  "createdAt": "2024-01-15T10:30:00Z",
  "updatedAt": "2024-01-15T10:30:00Z"
}
```

#### Obtenir une course par ID

Récupère une course spécifique par son ID.

```bash
curl -X GET http://localhost:8080/rides/{ride_id}
```

**Exemple :**

```bash
curl -X GET http://localhost:8080/rides/507f1f77bcf86cd799439011
```

**Réponse :**

```json
{
  "id": "507f1f77bcf86cd799439011",
  "passengerId": "passenger-001",
  "driverId": "507f1f77bcf86cd799439012",
  "from_zone": "Centre-ville",
  "to_zone": "Aéroport",
  "price": 25.5,
  "status": "ASSIGNED",
  "paymentStatus": "PENDING",
  "createdAt": "2024-01-15T10:30:00Z",
  "updatedAt": "2024-01-15T10:30:00Z"
}
```

#### Mettre à jour le statut d'une course

Met à jour le statut d'une course. Les statuts possibles sont :

- `ASSIGNED` : Course assignée à un chauffeur
- `IN_PROGRESS` : Course en cours
- `COMPLETED` : Course terminée (le paiement est automatiquement capturé et le chauffeur redevient disponible)
- `CANCELLED` : Course annulée

```bash
curl -X PATCH http://localhost:8080/rides/{ride_id}/status \
  -H "Content-Type: application/json" \
  -d '{
    "status": "COMPLETED"
  }'
```

**Exemple :**

```bash
curl -X PATCH http://localhost:8080/rides/507f1f77bcf86cd799439011/status \
  -H "Content-Type: application/json" \
  -d '{
    "status": "COMPLETED"
  }'
```

**Réponse :**

```json
{
  "id": "507f1f77bcf86cd799439011",
  "passengerId": "passenger-001",
  "driverId": "507f1f77bcf86cd799439012",
  "from_zone": "Centre-ville",
  "to_zone": "Aéroport",
  "price": 25.5,
  "status": "COMPLETED",
  "paymentStatus": "CAPTURED",
  "createdAt": "2024-01-15T10:30:00Z",
  "updatedAt": "2024-01-15T10:35:00Z"
}
```

### Comportement automatique

- **Lors de la création d'une course** :

  - Un chauffeur disponible est automatiquement sélectionné
  - Le chauffeur est marqué comme indisponible (`is_available: false`)

- **Lors de la complétion d'une course** (`status: "COMPLETED"`) :
  - Le paiement est automatiquement capturé (`paymentStatus: "CAPTURED"`)
  - Le chauffeur redevient disponible (`is_available: true`)

---

## Notes générales

### Identifiants

- Tous les IDs sont des ObjectIDs MongoDB (chaînes hexadécimales de 24 caractères)
- Les timestamps sont au format ISO 8601 (UTC)

### Bases de données

- **Users Database** : `ridenow_users`

  - Collections : `drivers`, `passengers`
  - Port externe : `27019`

- **Rides Database** : `ridenow_rides`
  - Collection : `rides`
  - Port externe : `27020`

### Variables d'environnement

#### Users Service

- `MONGO_URI` : URI de connexion MongoDB (par défaut : `mongodb://user_app:strong_app_password@users-service-database:27017/user_db?authSource=ridenow_users`)
- `PORT` : Port d'écoute du service (par défaut : `3000`)

#### Rides Service

- `MONGO_URI` : URI de connexion MongoDB (par défaut : `mongodb://user_app:strong_app_password@rides-service-database:27017/user_db?authSource=ridenow_rides`)
- `USERS_SERVICE_URL` : URL du service Users (par défaut : `http://users-service:3000`)
- `PORT` : Port d'écoute du service (par défaut : `8080`)

### Initialisation des bases de données

Les bases de données sont automatiquement initialisées avec des données d'exemple lors du premier démarrage via les scripts `init-mongo.js` :

- `services/users/init-mongo.js` : Crée des chauffeurs et passagers d'exemple
- `services/rides/init-mongo.js` : Crée des courses d'exemple
