# Inmobilia — Plataforma de Muebles de Melamina

Plataforma web para diseño paramétrico, visualización 3D, fabricación e instalación de muebles a medida.

## Estructura del monorepo

```txt
frontend/     React + TypeScript + Vite
backend/      Go (Gin, solver, manufacturing, costing)
```

## Arquitectura de capas

| Capa | Responsabilidad | Tipos |
|------|-----------------|-------|
| 1 — Room Geometry | Ambiente físico (paredes, pisos, obstáculos) | `RoomGeometry` |
| 2 — Furniture Definition | Intención de diseño (árbol de volúmenes) | `FurnitureDefinition`, `VolumeNode` |
| 3 — Installation Constraints | Conexión habitación ↔ mueble | `InstallationConstraints` |
| 4 — Constraint Solver | Resolución de medidas reales | `SolveConstraints()` |
| 5 — Resolved Furniture | Geometría calculada | `ResolvedFurniture` |
| 6 — Manufacturing Compiler | Piezas de fabricación | `ManufacturingModel` |
| 7 — Cost Engine | Materiales, desperdicio, costo total | `CostResult` |

**Regla fundamental:** la IA solo genera `FurnitureDefinition`. Nunca genera geometría Three.js ni piezas de fabricación.

## Flujo de API (objetivo)

```txt
POST /ai/parse          → FurnitureDefinition
POST /solver            → ResolvedFurniture
POST /manufacturing     → ManufacturingModel
POST /cost              → CostResult
```

## Orden de desarrollo

1. ✅ Domain Model
2. ✅ Room Geometry (utilidades y validación)
3. ✅ Volume Tree (utilidades y validación)
4. ✅ **Constraint Solver** ← núcleo del sistema
5. ✅ Resolved Furniture (validación y utilidades del output)
6. ✅ Manufacturing Compiler
7. ✅ Cost Engine
8. ✅ React 3D Viewer (solo `ResolvedFurniture`)
9. ✅ AI Integration
10. ✅ PDF / BOM / Cut Plans

## Desarrollo local

### Ambos a la vez (recomendado)

Desde la raíz del repo:

```bash
npm install          # solo la primera vez (instala concurrently)
cd frontend && npm install && cd ..
npm run dev
```

En Windows también podés usar:

```powershell
.\dev.ps1
```

Levanta el API en `http://localhost:8080` y el frontend en `http://localhost:5173` (proxy `/api` → backend).

### Frontend

```bash
cd frontend
npm install
npm run dev
```

### Backend

```bash
cd backend
cp .env.example .env   # opcional: configurar OpenAI
go build ./...
go run ./cmd/server
```

API en `http://localhost:8080` (variable `PORT`).

> **Importante:** las fuentes PDF (`backend/internal/export/fonts/*.ttf`) deben estar en el repo. Sin ellas el backend no compila.

## Domain model

- TypeScript: `frontend/src/domain/`
- Go: `backend/internal/domain/`
- Fixtures de ejemplo: `frontend/src/domain/fixtures/`

## Room Geometry (Fase 2)

Utilidades de geometría y validación de habitaciones. Los ángulos se calculan dinámicamente, nunca se almacenan.

| Capacidad | TypeScript | Go |
|-----------|------------|-----|
| Matemática vectorial / planos | `frontend/src/utils/geometry-math.ts` | `backend/pkg/geom/` |
| Validación de habitación | `validateRoomGeometry()` | `roomgeometry.ValidateRoomGeometry()` |
| Dimensiones derivadas | `getWallHeight()`, `getRoomFloorArea()`, etc. | `roomgeometry.GetWallHeight()`, etc. |
| Marco local de pared | `getWallLocalFrame()` | `roomgeometry.GetWallLocalFrame()` |
| Zócalos | `getSkirtingObstacles()` | `roomgeometry.GetSkirtingObstacles()` |

```bash
cd backend && go test ./internal/roomgeometry/...
```

## Volume Tree (Fase 3)

Utilidades de traversal y validación del árbol de volúmenes (`VolumeNode`).

| Capacidad | TypeScript | Go |
|-----------|------------|-----|
| Validación | `validateFurnitureDefinition()` | `volumetree.ValidateFurnitureDefinition()` |
| Traversal | `walkVolumeTree()`, `findVolumeNodeById()` | `volumetree.WalkVolumeTree()`, `FindVolumeNodeByID()` |
| Análisis | `getNodeCount()`, `collectFeatures()` | `volumetree.GetNodeCount()`, `CollectFeatures()` |
| Splits | `sumSplitRatios()`, `axisToDimension()` | `volumetree.SumSplitRatios()`, `AxisToDimension()` |

Reglas de validación:
- IDs únicos en nodos, features y fronts
- Nodos con hijos deben tener `split`; splits requieren hijos
- Ratios deben sumar 1; `fixed` debe ser positivo
- Hijos deben alinear constraints con el eje del split del padre

```bash
cd backend && go test ./internal/volumetree/...
```

## Constraint Solver (Fase 4)

Resuelve `RoomGeometry + FurnitureDefinition + InstallationConstraints` → `ResolvedFurniture`.

```go
resolved, err := solver.SolveConstraints(room, furniture, installation)
```

| Etapa | Responsabilidad |
|-------|-----------------|
| Validación | Rechaza room/furniture inválidos antes de resolver |
| Install space | Calcula zona útil con holguras, tolerancias, zócalos y offsets |
| Volume tree | Reparte espacio por splits (ratio/fixed) en ejes x/y/z |
| Constraints | `fill`, `fixed`, `ratio`, `min`, `max` por eje |
| Features/Fronts | Heredan el bounding box del compartimento resuelto |

Fixture de instalación: `frontend/src/domain/fixtures/example-installation.json`

```bash
cd backend && go test ./internal/solver/...
```

## Resolved Furniture (Fase 5)

Contrato y utilidades sobre la salida del solver. Solo medidas reales — sin `ratio`, `fill` ni constraints abstractos.

| Capacidad | TypeScript | Go |
|-----------|------------|-----|
| Validación | `validateResolvedFurniture()` | `resolvedfurniture.ValidateResolvedFurniture()` |
| Traversal | `walkResolvedTree()`, `findResolvedVolumeById()` | `WalkResolvedTree()`, `FindResolvedVolumeByID()` |
| Geometría | `volumeContains()`, `volumeOverlap()` | `VolumeContains()`, `VolumeOverlap()` |
| Métricas | `externalDimensions()`, `totalLeafVolumeMm3()` | `ExternalDimensions()`, `TotalLeafVolumeMm3()` |

El solver valida automáticamente su output antes de retornarlo.

```bash
cd backend && go test ./internal/resolvedfurniture/...
```

## Manufacturing Compiler (Fase 6)

Transforma `ResolvedFurniture` en piezas de fabricación.

```go
model, err := manufacturing.CompileManufacturing(resolved)
```

| Piezas generadas | Fuente |
|------------------|--------|
| Laterales, base, techo, trasera | Carcasa exterior (nodo root) |
| Divisiones | Splits entre hijos (X/Y/Z) |
| Repisas | Feature `shelf_set` |
| Cajones | Feature `drawer_stack` + frentes |
| Puertas | Front `door` / `drawer_front` |
| Herrajes | Bisagras, correderas, barral |
| Tapacantos | Bordes visibles automáticos |

```bash
cd backend && go test ./internal/manufacturing/...
```

## Cost Engine (Fase 7)

Calcula costos desde `ManufacturingModel`.

```go
cost, err := costing.CalculateCost(model)
```

| Componente | Cálculo |
|------------|---------|
| Materiales | m² de piezas agrupadas por material × tarifa |
| Tapacantos | metros lineales × tarifa por metro |
| Herrajes | cantidad × tarifa unitaria (bisagras, correderas, barral) |
| Mano de obra | horas estimadas × tarifa/hora |
| Desperdicio | % sobre costo de tableros |
| Total | subtotal + desperdicio |

Tarifas configurables con `CalculateCostWithOptions()`.

```bash
cd backend && go test ./internal/costing/...
```

## React 3D Viewer (Fase 8)

Renderiza exclusivamente `ResolvedFurniture` con React Three Fiber.

```bash
cd frontend
npm run dev
# Abrir http://localhost:5173/viewer
```

| Componente | Rol |
|------------|-----|
| `ResolvedFurnitureViewer` | Canvas R3F + luces + controles |
| `ResolvedVolumeMesh` | Volúmenes recursivos como cajas 3D |
| `FeatureMesh` / `FrontMesh` | Overlays de features y frentes |
| `viewer-store` | Selección y capas visibles (Zustand) |

Fixture resuelto: `frontend/src/domain/fixtures/example-resolved-closet.json`  
Regenerar: `cd backend/cmd/export-fixture && go run .`

## AI Integration (Fase 9)

`POST /ai/parse` convierte lenguaje natural en `FurnitureDefinition`. La IA **nunca** genera geometría Three.js ni piezas de fabricación.

| Endpoint | Input | Output |
|----------|-------|--------|
| `POST /ai/parse` | `{ description, name? }` | `FurnitureDefinition` |
| `POST /solver` | room + furniture + installation | `ResolvedFurniture` |
| `POST /manufacturing` | `{ resolved }` | `ManufacturingModel` |
| `POST /cost` | `{ model }` | `CostResult` |

### Variables de entorno (backend)

Copia `backend/.env.example` → `backend/.env`. El servidor carga `.env` automáticamente al arrancar desde `backend/`.

| Variable | Default | Descripción |
|----------|---------|-------------|
| `PORT` | `8080` | Puerto HTTP |
| `GIN_MODE` | `debug` | Usa `release` en producción |
| `CORS_ALLOWED_ORIGIN` | `*` | Origen permitido para CORS |
| `TRUSTED_PROXIES` | `127.0.0.1,::1` | IPs de proxy inverso (`none` = desactivar) |
| `AI_PROVIDER` | `mock` (o `openai` si hay API key) | Proveedor IA |
| `OPENAI_API_KEY` | — | API key OpenAI |
| `OPENAI_MODEL` | `gpt-4o-mini` | Modelo |
| `OPENAI_BASE_URL` | OpenAI oficial | URL compatible OpenAI |

Sin API key usa **mock parser** (closet/ropero/cajones por keywords).

### Variables de entorno (frontend)

Copia `frontend/.env.example` → `frontend/.env.local` si necesitas overrides.

| Variable | Default | Descripción |
|----------|---------|-------------|
| `VITE_API_BASE` | `/api` | Base URL del API (proxy Vite en dev) |

### Frontend

```bash
# Terminal 1
cd backend && go run ./cmd/server

# Terminal 2
cd frontend && npm run dev
# http://localhost:5173/assistant
```

## PDF / BOM / Cut Plans (Fase 10)

| Endpoint | Output |
|----------|--------|
| `POST /export/bom` | Bill of Materials JSON |
| `POST /export/cut-plans` | Planos de corte agrupados por material |
| `POST /export/pdf` | PDF con BOM, cortes, herrajes y costos |

```bash
cd backend && go test ./internal/export/...
```

Desde el AI Assistant (con pipeline completo) puedes descargar PDF, BOM JSON y planos de corte.

## Producción

### Checklist de despliegue

1. **Fuentes PDF** — confirmar que `backend/internal/export/fonts/DejaVuSans.ttf` y `DejaVuSans-Bold.ttf` están en el repositorio.
2. **Backend** — `GIN_MODE=release`, `CORS_ALLOWED_ORIGIN` apuntando al dominio del frontend, `TRUSTED_PROXIES` con la IP del reverse proxy.
3. **OpenAI** — `AI_PROVIDER=openai` y `OPENAI_API_KEY` en el entorno del servidor (nunca commitear `.env`).
4. **Frontend** — build estático con la URL real del API:

```bash
cd frontend
VITE_API_BASE=https://api.tu-dominio.com npm run build
```

Servir `frontend/dist/` detrás de nginx, Caddy o similar.

### Health check

`GET /health` devuelve:

```json
{ "status": "ok", "service": "inmobilia-api", "ai_provider": "openai (gpt-4o-mini)" }
```

### Auth (futuro)

JWT / multi-usuario no está implementado aún. El API es abierto en desarrollo; en producción restringe acceso con firewall, VPN o API gateway hasta añadir autenticación.
