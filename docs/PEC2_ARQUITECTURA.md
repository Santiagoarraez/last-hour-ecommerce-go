# PEC 2 - Arquitectura del servidor web en Go

## Interpretacion del enunciado

La PEC 2 pide migrar la web estatica de la PEC 1 a una aplicacion web dinamica. Por tanto, el navegador ya no accede directamente a documentos HTML independientes, sino que realiza peticiones HTTP a un servidor desarrollado en Go.

El servidor recibe cada peticion mediante `net/http`, selecciona un handler en funcion de la ruta, ejecuta la logica necesaria, accede a datos si corresponde y genera una respuesta HTML con plantillas.

## Conceptos del Tema 3 aplicados

- Modelo cliente-servidor: el navegador es el cliente y el programa Go actua como servidor web.
- Protocolo HTTP: cada interaccion se modela como una peticion y una respuesta.
- Rutas y handlers: cada URL se asocia a una funcion controladora.
- Metodos HTTP: `GET` se usa para consultar paginas y `POST` para enviar formularios, login, carrito y gestion de productos.
- Procesamiento en servidor: los datos de formularios se validan, se transforman y se guardan en el servidor.
- Generacion dinamica de HTML: las vistas se generan con `html/template`.
- Persistencia: los productos, mensajes, usuarios y carritos se almacenan en archivos JSON.
- Arquitectura por capas: se separan presentacion, logica de negocio y datos.

## Arquitectura elegida

Se aplica una arquitectura por capas similar a MVC:

```text
Cliente
  |
  | HTTP request
  v
Servidor Go con net/http
  |
  v
Handlers
  |
  v
Services
  |
  v
Storage
  |
  v
Archivos JSON
```

La respuesta vuelve al cliente mediante plantillas HTML:

```text
Datos procesados -> Templates -> HTML -> HTTP response
```

## Responsabilidades

- `cmd/server/main.go`: punto de entrada. Crea dependencias, registra rutas y arranca el servidor.
- `internal/handlers`: capa de presentacion/controladores HTTP. Lee peticiones y genera respuestas.
- `internal/services`: capa de logica de negocio. Valida datos y coordina operaciones.
- `internal/models`: estructuras de datos de la aplicacion.
- `internal/storage`: capa de persistencia. Lee y escribe archivos JSON.
- `templates`: vistas HTML dinamicas.
- `data`: datos persistentes en formato JSON.
- `css` y `assets`: recursos estaticos servidos por el servidor.

## Rutas HTTP

| Metodo | Ruta | Handler | Funcion |
| --- | --- | --- | --- |
| GET | `/` | `Home` | Muestra la pagina principal con productos destacados |
| GET | `/products` | `Products` | Muestra el catalogo desde JSON |
| GET | `/products/{id}` | `ProductDetail` | Muestra un producto concreto |
| GET | `/about` | `About` | Muestra informacion del proyecto |
| GET | `/contact` | `Contact` | Muestra el formulario |
| POST | `/contact` | `Contact` | Procesa y guarda el formulario |
| GET | `/login` | `Login` | Muestra el formulario de acceso |
| POST | `/login` | `Login` | Procesa credenciales y crea cookie de sesion |
| GET | `/register` | `Register` | Muestra el formulario de registro |
| POST | `/register` | `Register` | Crea un comprador en `users.json` |
| POST | `/logout` | `Logout` | Elimina la cookie de sesion |
| GET | `/account` | `Account` | Muestra los datos del usuario autenticado |
| GET | `/cart` | `Cart` | Muestra el carrito del comprador |
| POST | `/cart/add` | `CartAdd` | Agrega un producto al carrito |
| POST | `/cart/remove` | `CartRemove` | Elimina un producto del carrito |
| POST | `/cart/checkout` | `CartCheckout` | Procesa la compra simulada y vacia el carrito |
| GET | `/seller/products` | `SellerProducts` | Lista productos para vendedor |
| GET | `/seller/products/new` | `SellerProductNew` | Formulario para crear producto |
| POST | `/seller/products/new` | `SellerProductNew` | Crea producto en `products.json` |
| GET | `/seller/products/edit/{id}` | `SellerProductEdit` | Formulario de edicion |
| POST | `/seller/products/edit/{id}` | `SellerProductEdit` | Actualiza producto |
| POST | `/seller/products/delete/{id}` | `SellerProductDelete` | Elimina producto |

## Flujo del formulario

1. El cliente solicita `GET /contact`.
2. El servidor devuelve el formulario HTML.
3. El cliente envia los datos mediante `POST /contact`.
4. El handler lee los campos del `request`.
5. El servicio valida que los campos obligatorios no esten vacios.
6. La capa de persistencia guarda el mensaje en `data/messages.json`.
7. El servidor genera una respuesta HTML de confirmacion.

## Flujo del carrito

1. El cliente autenticado solicita `GET /products/{id}`.
2. El servidor genera una pagina de detalle con un formulario `POST /cart/add`.
3. El cliente envia `product_id` y `quantity`.
4. El handler comprueba la cookie de sesion y obtiene el usuario.
5. `CartService` valida que el producto existe mediante `ProductService`.
6. `CartStorage` actualiza `data/carts.json`.
7. El servidor redirige a `GET /cart` y renderiza el carrito actualizado.

## Flujo del vendedor

1. El vendedor inicia sesion mediante `POST /login`.
2. El servidor valida credenciales y guarda una cookie de sesion.
3. El vendedor accede a `GET /seller/products`.
4. Los handlers de vendedor comprueban que el rol del usuario sea `seller`.
5. Las altas, ediciones y borrados de productos se procesan con metodos `POST`.
6. `ProductService` aplica la logica y `ProductStorage` modifica `data/products.json`.

## Justificacion tecnica

Se usa `net/http` porque forma parte de la biblioteca estandar de Go y permite implementar directamente el modelo HTTP peticion-respuesta estudiado en el Tema 3.

Se usa `html/template` para separar el HTML de la logica de negocio. Los handlers no contienen codigo HTML; solo seleccionan la plantilla y le pasan datos.

Se usa JSON como persistencia porque el enunciado permite guardar la informacion en JSON o JSONL. Para esta web, JSON es suficiente y evita introducir una base de datos no necesaria para el alcance de la practica.

La carpeta `internal` se usa porque es una convencion habitual en proyectos Go para codigo propio de la aplicacion que no debe importarse desde otros proyectos.

La autenticacion se implementa de forma didactica mediante una cookie `user_id` y usuarios persistidos en JSON. Esta decision permite demostrar sesiones, perfiles y control de acceso sin introducir librerias externas ni bases de datos, manteniendo el foco de la PEC en el modelo cliente-servidor y el procesamiento HTTP en servidor.
