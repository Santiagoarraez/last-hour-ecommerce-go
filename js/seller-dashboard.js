/**
 * seller-dashboard.js
 * PEC 3 — Redes y Sistemas Web
 *
 * Módulo JavaScript puro (sin librerías externas) que gestiona el
 * panel de administración de productos usando:
 *   - Promesas (fetch API con .then()/.catch())
 *   - Etiquetas <template> de HTML para renderizar el DOM
 *   - CRUD completo sobre la API REST en /api/products
 */

// ── Constantes ────────────────────────────────────────────────────────────────

const API_BASE = '/api/products';

// ── Referencias al DOM ────────────────────────────────────────────────────────

const tbodyEl        = document.getElementById('products-tbody');
const tableWrapper   = document.getElementById('table-wrapper');
const tableLoading   = document.getElementById('table-loading');
const tableEmpty     = document.getElementById('table-empty');
const formPanel      = document.getElementById('form-panel');
const formPanelTitle = document.getElementById('form-panel-title');
const productForm    = document.getElementById('product-form');
const productIdInput = document.getElementById('product-id');
const notifArea      = document.getElementById('dashboard-notification');
const searchInput    = document.getElementById('search-input');

// Botones principales
const btnNewProduct = document.getElementById('btn-new-product');
const btnCancelForm = document.getElementById('btn-cancel-form');
const btnCancelForm2 = document.getElementById('btn-cancel-form-2');

// Campos del formulario
const fieldName        = document.getElementById('field-name');
const fieldSubtitle    = document.getElementById('field-subtitle');
const fieldDescription = document.getElementById('field-description');
const fieldPrice       = document.getElementById('field-price');
const fieldImage       = document.getElementById('field-image');
const fieldAlt         = document.getElementById('field-alt');
const fieldFlavors     = document.getElementById('field-flavors');
const fieldFeatured    = document.getElementById('field-featured');

// Plantillas HTML (<template>)
const tplRow            = document.getElementById('tpl-product-row');
const tplNotifySuccess  = document.getElementById('tpl-notify-success');
const tplNotifyError    = document.getElementById('tpl-notify-error');

// Cache en memoria para búsqueda local sin re-peticiones
let allProducts = [];

// ── Funciones de la API (Promesas) ────────────────────────────────────────────

/**
 * apiFetch — función base para todas las peticiones a la API REST.
 * Devuelve una Promesa que resuelve con el JSON de la respuesta,
 * o rechaza con un Error si el servidor devuelve un código de error.
 *
 * @param {string} url     - URL de la API
 * @param {object} options - opciones de fetch (method, body, headers)
 * @returns {Promise<any>}
 */
function apiFetch(url, options = {}) {
  const defaults = {
    headers: { 'Content-Type': 'application/json' },
  };

  return fetch(url, Object.assign(defaults, options))
    .then(function (response) {
      // 204 No Content (DELETE exitoso) no tiene cuerpo JSON
      if (response.status === 204) {
        return null;
      }
      // Para cualquier otro código no-OK, leer el error del JSON
      return response.json().then(function (data) {
        if (!response.ok) {
          throw new Error(data.error || 'Error desconocido en el servidor');
        }
        return data;
      });
    });
}

/**
 * fetchProducts — GET /api/products
 * Obtiene todos los productos del servidor.
 * @returns {Promise<Array>}
 */
function fetchProducts() {
  return apiFetch(API_BASE);
}

/**
 * createProduct — POST /api/products
 * Envía un nuevo producto al servidor.
 * @param {object} productData
 * @returns {Promise<object>}
 */
function createProduct(productData) {
  return apiFetch(API_BASE, {
    method: 'POST',
    body: JSON.stringify(productData),
  });
}

/**
 * updateProduct — PUT /api/products/{id}
 * Actualiza un producto existente.
 * @param {string} id
 * @param {object} productData
 * @returns {Promise<object>}
 */
function updateProduct(id, productData) {
  return apiFetch(API_BASE + '/' + id, {
    method: 'PUT',
    body: JSON.stringify(productData),
  });
}

/**
 * deleteProduct — DELETE /api/products/{id}
 * Elimina un producto permanentemente.
 * @param {string} id
 * @returns {Promise<null>}
 */
function deleteProduct(id) {
  return apiFetch(API_BASE + '/' + id, {
    method: 'DELETE',
  });
}

// ── Renderizado con <template> ────────────────────────────────────────────────

/**
 * renderProductRow — clona el <template id="tpl-product-row"> y
 * rellena sus campos con los datos del producto.
 * @param {object} product
 * @returns {HTMLElement} - el nodo <tr> listo para insertar
 */
function renderProductRow(product) {
  // Clonamos el contenido del <template> (sin activarlo aún)
  const clone = tplRow.content.cloneNode(true);
  const row = clone.querySelector('tr');

  // Establecemos el ID del producto en el atributo data del <tr>
  row.dataset.id = product.id;

  // Imagen del producto
  const img = row.querySelector('img');
  img.src = product.image || '/assets/images/hqd-catalog-new.png';
  img.alt = product.alt || product.name;

  // Rellenamos los campos de texto
  row.querySelector('.dashboard__product-name').textContent     = product.name;
  row.querySelector('.dashboard__product-subtitle').textContent = product.subtitle;
  row.querySelector('.dashboard__product-price').textContent    = product.price.toFixed(2) + ' €';

  // Badge de "Destacado"
  const badge = row.querySelector('.dashboard__badge');
  if (product.featured) {
    badge.textContent = '★ Destacado';
    badge.classList.add('dashboard__badge--featured');
  } else {
    badge.textContent = 'Normal';
    badge.classList.add('dashboard__badge--normal');
  }

  // Listeners de los botones de acción
  row.querySelector('.btn-edit').addEventListener('click', function () {
    openEditForm(product);
  });

  row.querySelector('.btn-delete').addEventListener('click', function () {
    confirmDelete(product);
  });

  return row;
}

/**
 * renderTable — dibuja la tabla completa con la lista de productos.
 * Gestiona también los estados vacío/cargando.
 * @param {Array} products
 */
function renderTable(products) {
  // Limpiamos las filas actuales del tbody
  tbodyEl.innerHTML = '';

  // Ocultamos el spinner de carga
  tableLoading.hidden = true;

  if (!products || products.length === 0) {
    tableWrapper.hidden = true;
    tableEmpty.hidden = false;
    return;
  }

  tableEmpty.hidden = true;
  tableWrapper.hidden = false;

  // Añadimos cada fila clonando el <template>
  products.forEach(function (product) {
    tbodyEl.appendChild(renderProductRow(product));
  });
}

// ── Notificaciones ────────────────────────────────────────────────────────────

/**
 * showNotification — muestra un mensaje temporal de éxito o error
 * usando las plantillas <template> correspondientes.
 * @param {'success'|'error'} type
 * @param {string} message
 */
function showNotification(type, message) {
  const tpl = type === 'success' ? tplNotifySuccess : tplNotifyError;
  const clone = tpl.content.cloneNode(true);
  clone.querySelector('.notify-message').textContent = message;

  notifArea.innerHTML = '';
  notifArea.appendChild(clone);
  notifArea.hidden = false;

  // Auto-ocultar después de 4 segundos
  setTimeout(function () {
    notifArea.hidden = true;
    notifArea.innerHTML = '';
  }, 4000);
}

// ── Gestión del Formulario ────────────────────────────────────────────────────

/** Abre el panel con el formulario vacío para crear un producto. */
function openCreateForm() {
  productForm.reset();
  productIdInput.value = '';
  formPanelTitle.textContent = 'Nuevo Producto';
  document.getElementById('btn-submit-form').innerHTML =
    '<i class="fa-solid fa-plus"></i> Crear Producto';
  formPanel.hidden = false;
  fieldName.focus();
  // Scroll suave hasta el formulario
  formPanel.scrollIntoView({ behavior: 'smooth', block: 'start' });
}

/**
 * Abre el panel con el formulario relleno con los datos del producto a editar.
 * @param {object} product
 */
function openEditForm(product) {
  productIdInput.value    = product.id;
  fieldName.value         = product.name;
  fieldSubtitle.value     = product.subtitle;
  fieldDescription.value  = product.description;
  fieldPrice.value        = product.price;
  fieldImage.value        = product.image;
  fieldAlt.value          = product.alt || '';
  fieldFlavors.value      = (product.flavors || []).join(', ');
  fieldFeatured.checked   = product.featured;

  formPanelTitle.textContent = 'Editar Producto';
  document.getElementById('btn-submit-form').innerHTML =
    '<i class="fa-solid fa-cloud-arrow-up"></i> Guardar Cambios';
  formPanel.hidden = false;
  fieldName.focus();
  formPanel.scrollIntoView({ behavior: 'smooth', block: 'start' });
}

/** Cierra y resetea el panel del formulario. */
function closeForm() {
  formPanel.hidden = true;
  productForm.reset();
  productIdInput.value = '';
}

/**
 * Recoge los datos del formulario y los devuelve como objeto.
 * @returns {object}
 */
function getFormData() {
  return {
    name:        fieldName.value.trim(),
    subtitle:    fieldSubtitle.value.trim(),
    description: fieldDescription.value.trim(),
    price:       parseFloat(fieldPrice.value),
    image:       fieldImage.value.trim(),
    alt:         fieldAlt.value.trim(),
    flavors:     fieldFlavors.value.trim(),
    featured:    fieldFeatured.checked,
  };
}

// ── Confirmación de borrado ───────────────────────────────────────────────────

/**
 * Pide confirmación y elimina el producto.
 * Usa una Promesa encadenada: deleteProduct → loadProducts.
 * @param {object} product
 */
function confirmDelete(product) {
  if (!window.confirm('¿Eliminar "' + product.name + '"? Esta acción no se puede deshacer.')) {
    return;
  }

  deleteProduct(product.id)
    .then(function () {
      showNotification('success', '"' + product.name + '" eliminado correctamente.');
      return loadProducts(); // Recargamos la tabla tras el borrado
    })
    .catch(function (err) {
      showNotification('error', 'No se pudo eliminar: ' + err.message);
    });
}

// ── Carga inicial ─────────────────────────────────────────────────────────────

/**
 * loadProducts — carga los productos de la API y re-renderiza la tabla.
 * Promesa principal que se encadena con renderTable.
 * @returns {Promise<void>}
 */
function loadProducts() {
  tableLoading.hidden = false;
  tableWrapper.hidden = true;
  tableEmpty.hidden = true;

  return fetchProducts()
    .then(function (products) {
      allProducts = products || [];
      renderTable(allProducts);
    })
    .catch(function (err) {
      tableLoading.hidden = true;
      showNotification('error', 'Error al cargar los productos: ' + err.message);
    });
}

// ── Búsqueda local ────────────────────────────────────────────────────────────

/** Filtra la tabla según el texto de búsqueda (sin llamada a la API). */
function filterTable(query) {
  const q = query.toLowerCase().trim();
  if (!q) {
    renderTable(allProducts);
    return;
  }
  const filtered = allProducts.filter(function (p) {
    return p.name.toLowerCase().includes(q) ||
           p.subtitle.toLowerCase().includes(q);
  });
  renderTable(filtered);
}

// ── Event Listeners ───────────────────────────────────────────────────────────

// Botón "Nuevo Producto"
btnNewProduct.addEventListener('click', openCreateForm);

// Botones "Cancelar"
btnCancelForm.addEventListener('click', closeForm);
btnCancelForm2.addEventListener('click', closeForm);

// Búsqueda en tiempo real (filtra la caché local, no llama a la API)
searchInput.addEventListener('input', function () {
  filterTable(this.value);
});

// Envío del formulario (crear o actualizar según si hay ID)
productForm.addEventListener('submit', function (event) {
  // Prevenimos el POST tradicional del formulario
  event.preventDefault();

  const data = getFormData();
  const id   = productIdInput.value;

  // Deshabilitamos el botón para evitar doble envío
  const submitBtn = document.getElementById('btn-submit-form');
  submitBtn.disabled = true;

  // Decidimos si es creación (sin ID) o edición (con ID)
  const operation = id
    ? updateProduct(id, data).then(function (updated) {
        showNotification('success', '"' + updated.name + '" actualizado correctamente.');
      })
    : createProduct(data).then(function (created) {
        showNotification('success', '"' + created.name + '" creado correctamente.');
      });

  operation
    .then(function () {
      closeForm();
      return loadProducts(); // Recargamos la tabla
    })
    .catch(function (err) {
      showNotification('error', 'Error: ' + err.message);
    })
    .finally(function () {
      submitBtn.disabled = false;
    });
});

// ── Arranque ──────────────────────────────────────────────────────────────────

// Cargamos los productos al cargar la página
loadProducts();
