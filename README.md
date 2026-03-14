# Gateway Service

El **Gateway Service** (*API Gateway*) actúa como el punto de entrada único (Single Entry Point) para todas las peticiones externas dirigidas al **School Tracking System**. Este servicio es responsable de enrutar el tráfico HTTP de forma segura hacia los microservicios internos correspondientes, implementar políticas de CORS, y validar la autenticidad de los usuarios.

## Funcionalidades Principales

- **Reverse Proxying**: Redirige dinámicamente las solicitudes hacia otros microservicios (ej. `auth`, `fleet`, `tracking`).
- **Autenticación Unificada (JWT)**: Intercepta y valida tokens JWT. Si el token es válido, extrae información fundamental (como `UserID` y `Role`) y la inyecta como *Headers* HTTP (`X-User-ID`, `X-User-Role`) para que los servicios internos puedan consumirla sin necesidad de volver a procesar el token.
- **Seguridad Perimetral**: Configuración base segura (Timeouts, CORS configs).
- *(En el futuro)*: Terminación SSL, Rate Limiting, y enrutamiento WebSockets gRPC.

## Requisitos Previos

- Go 1.24.0+
- (Recomendado) Asegúrate de que el [Auth Service](../auth) esté corriendo localmente en el puerto `8080` para poder rutearle el tráfico de autenticación.

## Endpoints y Ruteo (Proxy Config)

Toda solicitud que empiece por `/api/v1` pasa a través del Gateway.
Actualmente, redirige las siguientes rutas al servicio `auth`:

```
[Gateway] GET/POST/PUT /api/v1/auth/*  --->  [Auth Service] GET/POST/PUT /api/v1/auth/*
```

Por ejemplo, un login en tu entorno de desarrollo sería:
`POST http://localhost:8000/api/v1/auth/login`

## Variables de Entorno

El servicio utiliza el paquete `github.com/caarlos0/env/v10` para leer la configuración. Si existe un archivo `.env` en la raíz del servicio, lo cargará automáticamente al arrancar.

> **IMPORTANTE**: Existe un archivo `.env.template` en la raíz de este proyecto. Cópialo y renómbralo a `.env` en tu máquina local.

| Variable | Tipo | Default | Descripción |
| :--- | :--- | :--- | :--- |
| `PORT` | int | `8000` | Puerto en el que iniciará el servidor (evita conflicto con `8080`). |
| `ENVIRONMENT` | string | `development` | Entorno de ejecución (`development`, `production`). |
| `AUTH_SERVICE_URL` | string | `http://localhost:8080` | URL base del microservicio Auth para el proxy inverso. |
| `JWT_SECRET` | string | `dev-secret-change-in-prod` | Secreto asimétrico/simétrico para validar la firma de los tokens JWT de autorización emitidos por `Auth`. |

### Configuración de Seguridad para el JWT
El `JWT_SECRET` configurado en el Gateway **debe ser exactamente el mismo** que utilizó el servicio `Auth` para firmar y emitir el token.

## Setup & Run Local

1. Navega al directorio del servicio.
   ```bash
   cd services/gateway
   ```
2. Crea tu archivo de entorno a partir del template.
   ```bash
   cp .env.template .env
   ```
3. Ejecuta el servicio localmente.
   ```bash
   go run cmd/api/main.go cmd/api/module.go
   ```

El Gateway estará corriendo y escuchando en el puerto definido (por defecto `8000`), listo para enrutar el tráfico.
