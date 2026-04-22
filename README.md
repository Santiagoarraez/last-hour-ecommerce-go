# Last Hour - E-Commerce

## Descripción general del proyecto

Este proyecto consiste en el desarrollo de una página web responsive tipo e-commerce para una marca ficticia de vapeadores denominada "Last Hour". La aplicación simula una tienda online completa enfocada en la presentación y venta de productos de la marca. El objetivo primario de este proyecto es poner en práctica y consolidar las tecnologías y conceptos básicos del desarrollo front-end (del lado del cliente) adquiridos durante la asignatura, construyendo una interfaz de usuario atractiva, funcional y adaptable a cualquier dispositivo móvil o de escritorio.

## Objetivos del proyecto

*   **Aplicación práctica de tecnologías clave:** Demostrar dominio en el uso conjunto de HTML5 y CSS3 para la creación de interfaces web.
*   **Desarrollo de diseño responsive:** Garantizar una experiencia de usuario (UX) óptima adaptando el contenido visual a múltiples resoluciones de pantalla utilizando técnicas CSS modernas y *media queries*.
*   **Organización y mantenibilidad del código:** Estructurar el código fuente siguiendo buenas prácticas, utilizando modularidad y preprocesadores CSS para facilitar la lectura y el mantenimiento a largo plazo.
*   **Implementación sin dependencias externas:** Desarrollar la interfaz "desde cero", sin depender de frameworks CSS o librerías externas complejas, asegurando una comprensión profunda de las bases del desarrollo web.

## Tecnologías utilizadas

El proyecto ha sido construido utilizando las siguientes tecnologías estándar de la web:

*   **HTML5:** Empleado para la estructura semántica de la página, asegurando accesibilidad y correcta indexación.
*   **CSS3:** Utilizado para el diseño visual, incluyendo el uso de Flexbox y CSS Grid para la creación de *layouts* fluidos y adaptables.
*   **Sass (SCSS):** Preprocesador CSS utilizado para modularizar las hojas de estilo, permitiendo el uso de variables, anidamientos y *mixins*.

## Estructura del proyecto

La organización de carpetas y archivos se ha planificado para separar lógicamente los diferentes recursos del proyecto:

```text
📁 Proyecto
├── 📄 index.html        # Estructura principal de la página de inicio.
├── 📄 products.html     # (Ejemplo) Página de catálogo de productos.
├── 📄 product-detail.html # (Ejemplo) Página de detalle de producto.
├── 📄 contact.html      # (Ejemplo) Página de contacto.
├── 📄 about.html        # (Ejemplo) Página "Sobre nosotros".
└── 📁 assets/           # Contenedor principal de recursos estáticos.
    ├── 📁 images/       # Fotografías, logotipos e iconos utilizados en la interfaz.
    ├── 📁 scss/         # Archivos fuente de Sass, organizados en módulos.
    └── 📁 css/          # Archivos CSS compilados a partir de los archivos SCSS.
```

*   La raíz contiene los documentos HTML principales que definen las distintas vistas de la tienda.
*   La carpeta `assets/` agrupa de forma ordenada todos los complementos estáticos necesarios para el renderizado y funcionamiento visual de la página.

## Explicación de los componentes principales

La interfaz web está construida en base a componentes clave que se repiten a lo largo del sitio:

1.  **Navbar (Barra de navegación):** Componente principal de navegación, implementado con un diseño responsive. En pantallas amplias, muestra el logotipo centrado, enlaces de navegación a la izquierda, y accesos rápidos (cuenta y carrito) a la derecha. En pantallas menores, se colapsa en un menú oculto accesible mediante un botón hamburguesa interactivo, implementado puramente con CSS mediante la técnica del "checkbox hack" (sin necesidad de JavaScript).
2.  **Hero Section:** Área visual destacada en la parte superior de la página principal. Presenta una imagen representativa de "Last Hour" junto con una llamada a la acción (CTA) para invitar al usuario a explorar la tienda.
3.  **Sección de productos destacados:** Una cuadrícula (*grid*) que expone una selección de artículos disponibles. Cada ítem incluye una fotografía de ejemplo, título y precio, permitiendo al usuario previsualizar la oferta.
4.  **Footer (Pie de página):** Cierre informativo de la página que contiene los detalles corporativos de la marca. Incluye contacto directo (correo electrónico y teléfono), ubicación física y un enlace a redes sociales (Instagram), reforzando el *branding* de la marca.

## Cómo se implementó el diseño responsive

La adaptabilidad a dispositivos (móviles, tablets y pantallas de escritorio) se resolvió sin utilizar frameworks externos mediante las siguientes técnicas:

*   **Media Queries:** Interrupciones definidas en dimensiones específicas (`min-width` y `max-width`) para ajustar la disposición de los bloques, ocultar/mostrar elementos y escalar la tipografía según el tamaño del marco de visualización (*viewport*).
*   **Unidades fluidas:** Empleo de porcentajes (`%`) y unidades relativas a la ventana (`vw`, `vh`, `rem`) para dimensionar contenedores, permitiendo que crezcan o se reduzcan de forma proporcional al área disponible.
*   **Flexbox y CSS Grid:** Herramientas nativas de CSS utilizadas exhaustivamente para reubicar elementos horizontal y verticalmente. Esto asegura que la cuadrícula de productos, por ejemplo, adapte el número de columnas de forma orgánica según la resolución.

## Uso de Sass y organización de estilos

Para mantener el código CSS limpio, escalable y evitar redundancias repetitivas, se ha implementado el preprocesador **Sass** con la sintaxis **SCSS**. Las principales estrategias aplicadas incluyen:

*   **Modularidad:** El código de estilos se ha dividido en múltiples archivos (arquitectura de componentes), que luego son agrupados en una única hoja de estilos CSS compilada de forma final.
*   **Variables:** Asignación de colores corporativos, familias tipográficas y tamaños base a variables (`$primary-color`, `$font-base`, etc.) para asegurar la consistencia del manual de identidad visual de "Last Hour" en todo el proyecto.
*   **Anidamiento (Nesting):** Encapsulación lógica de los estilos dentro de sus selectores padre, mejorando notablemente la legibilidad del código SCSS.

## Cómo ejecutar el proyecto localmente

Dado que el proyecto se compone únicamente de lenguajes interpretados en el lado del cliente (diseño estático), su ejecución es sencilla y no requiere backend ni bases de datos.

1.  Clonar el repositorio o descargar el código fuente a su máquina local.
2.  Localizar la carpeta raíz del proyecto (`PEC_1_Last_hour`).
3.  Hacer doble clic sobre el archivo `index.html` para abrirlo directamente en el navegador web predeterminado (Chrome, Firefox, Safari, Edge).
4.  *Sugerencia para desarrolladores:* Para habilitar funcionalidades de auto-recarga (*live reload*), se recomienda abrir el proyecto en un editor como Visual Studio Code y levantar un servidor de desarrollo mediante la extensión *Live Server*.

## Correcciones aplicadas y justificación técnica

A continuación se documentan los cambios implementados para corregir los problemas detectados en la práctica:

* Se ha eliminado el uso de estilos inline en `product-detail.html`, separando completamente estructura y presentación.
* La página de producto se ha reescrito con una semántica HTML correcta: se ha utilizado `<article>`, `<header>`, `<figure>`, `<section>`, `<form>` y etiquetas `<label>` asociadas a sus controles.
* Se ha añadido una hoja de estilos mobile-first para la vista de producto en `scss/pages/_product.scss` y se ha reflejado esa misma lógica en `css/styles.css`.
* El diseño del detalle de producto ahora define estilos base para móvil y amplía el layout con `@media (min-width: ...)`, respetando el enfoque mobile-first.
* Se ha reforzado la consistencia visual mediante la aplicación de clases reutilizables para formularios y botones, en lugar de atributos `style` en el HTML.

Estas correcciones mejoran la accesibilidad, el mantenimiento y el cumplimiento de buenas prácticas profesionales en front-end.

## Uso de IA en el desarrollo

Durante la corrección de la práctica se utilizó un asistente de IA como apoyo técnico para:

* Analizar el enfoque actual de responsive design y confirmar que el diseño debía ser reescrito en mobile-first.
* Proponer la estructura semántica correcta de la página de producto.
* Generar los estilos de `product-view` de forma coherente con la arquitectura SCSS del proyecto.
* Documentar el proceso de cambios con una explicación académica y precisa.

Este uso de IA se ha integrado como apoyo de revisión y no como sustituto de la implementación manual del código.

## Posibles mejoras futuras

A fin de escalar el proyecto y enriquecer la experiencia de usuario general, se plantean posibles iteraciones de mejora:

*   **Funcionalidad de Carrito Dinámico:** Implementar lógica avanzada de JavaScript para permitir a los usuarios añadir y eliminar productos, calcular totales y almacenar temporalmente el progreso a través de `localStorage` o `sessionStorage`.
*   **Integración de Backend y Base de Datos:** Conectar el front-end estático con una API (servidor) para extraer el inventario dinámicamente y procesar interacciones reales, como registros de usuario o procesos de pago ("*checkout*").
*   **Auditoría de Accesibilidad (a11y):** Añadir atributos ARIA, garantizar contrastes de color apropiados y asegurar que la página web sea navegable únicamente mediante teclado, para cumplir así de manera estricta los estándares web modernos.
*   **Micro-interacciones y Animaciones:** Incluir librerías de animación de interfaz o keyframes en CSS para agregar transiciones sutiles cuando los elementos entren en el marco visual (Scroll animations).
