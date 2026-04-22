# Last Hour - E-Commerce Dinámico (PEC 2)

## Descripción general del proyecto
Este proyecto es la evolución de la tienda "Last Hour" hacia una **aplicación web dinámica** desarrollada con el lenguaje de programación **Go (Golang)**. Se trata de un e-commerce de vapeadores que ahora cuenta con lógica de servidor, gestión de datos persistentes y medidas de seguridad profesionales.

El objetivo de esta fase (PEC 2) ha sido migrar la estructura estática previa a una arquitectura de capas (MVC), permitiendo interacciones reales como el registro de usuarios, la edición de perfiles y la gestión de pedidos.

## Nuevas Funcionalidades (PEC 2)
*   **Seguridad con Bcrypt:** Implementación de encriptación de contraseñas para todos los usuarios.
*   **Gestión de Perfiles:** Los usuarios pueden ver y editar su información personal (Nombre, Email, Teléfono) desde un panel dedicado.
*   **Checkout por WhatsApp:** Sistema de pedidos que genera un mensaje automático detallado con los productos del carrito y el total, redirigiendo al usuario para confirmar la compra por WhatsApp.
*   **Interfaz Enriquecida:** Integración de **Font Awesome** para una experiencia de usuario más profesional mediante el uso de iconos en la navegación, formularios y botones.
*   **Servidor Web en Go:** Servidor robusto que gestiona rutas, plantillas HTML dinámicas y persistencia en archivos JSON.

## Tecnologías utilizadas
*   **Backend:** Go (Golang) 1.22+.
*   **Seguridad:** Bcrypt (hashing de contraseñas).
*   **Frontend:** HTML5, CSS3, SCSS.
*   **Iconografía:** Font Awesome 6.
*   **Arquitectura:** Diseño por capas (Handlers, Services, Storage, Models).

## Estructura del proyecto
```text
📁 PEC_1_Last_hour
├── 📁 bin/              # Binarios compilados (Linux/Windows)
├── 📁 cmd/server/       # Punto de entrada del servidor (main.go)
├── 📁 internal/         # Lógica de la aplicación
│   ├── 📁 handlers/     # Controladores de rutas web (MVC - Controllers)
│   ├── 📁 models/       # Estructuras de datos (MVC - Models)
│   ├── 📁 services/     # Lógica de negocio (Seguridad, Carrito)
│   └── 📁 storage/      # Persistencia de datos (JSON)
├── 📁 data/             # Archivos JSON (Usuarios, Productos, Carritos)
├── 📁 templates/        # Plantillas HTML dinámicas
├── 📁 assets/           # Imágenes y recursos estáticos
└── 📁 css/              # Hojas de estilo compiladas
```

## Guía de Ejecución (Windows/WSL)
Debido a restricciones de seguridad de Windows con binarios no firmados de Go, se recomienda ejecutar el proyecto utilizando el entorno **WSL (Windows Subsystem for Linux)**.

### Requisitos previos
1.  Tener instalado **WSL** (Ubuntu u otra distribución).
2.  Tener instalado **Go** en Windows.

### Pasos para ejecutar:
1.  **Compilar para Linux** (desde PowerShell):
    ```powershell
    $env:GOOS='linux'; $env:GOARCH='amd64'; go build -o bin/lasthour-linux ./cmd/server
    ```
2.  **Ejecutar mediante WSL**:
    ```powershell
    wsl ./bin/lasthour-linux
    ```
3.  **Acceder a la aplicación**:
    Abre tu navegador en: [http://localhost:8080](http://localhost:8080)

## Usuarios de Prueba (Demo)
*   **Vendedor:** `seller@lasthour.com` / clave: `seller123`
*   **Cliente:** `customer@lasthour.com` / clave: `customer123`

## Uso de IA y Ética
Este proyecto ha contado con el apoyo de IA para la refactorización de código, implementación de seguridad y documentación, asegurando que el estudiante comprenda cada proceso y manteniendo la integridad académica de la PEC.

---
*PEC 2 - Redes y Sistemas Web - Grado en Ingeniería Informática*
