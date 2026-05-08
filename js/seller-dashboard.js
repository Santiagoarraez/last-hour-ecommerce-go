/**
 * seller-dashboard.js
 * PEC 3 — Redes y Sistemas Web
 * Gestión del panel de administración (Modelos, Sabores, Promociones)
 */

document.addEventListener('DOMContentLoaded', () => {
  // ── Constantes de la API ──────────────────────────────────────────────────
  const API_MODELS = '/api/models';
  const API_FLAVORS = '/api/flavors';
  const API_PROMOTIONS = '/api/promotions';

  // ── Referencias al DOM ────────────────────────────────────────────────────
  const tabButtons = document.querySelectorAll('.dashboard__tab-btn');
  const tabContents = document.querySelectorAll('.dashboard__tab-content');

  const modelForm = document.getElementById('model-form');
  const modelTbody = document.getElementById('models-tbody');

  const flavorForm = document.getElementById('flavor-form');
  const flavorModelSelect = document.getElementById('flavor-model-id');
  const flavorTbody = document.getElementById('flavors-tbody');

  const promoForm = document.getElementById('promotion-form');
  const promoItemsContainer = document.getElementById('promo-items-container');
  const btnAddPromoItem = document.getElementById('btn-add-promo-item');
  const promoTbody = document.getElementById('promotions-tbody');

  const notifArea = document.getElementById('dashboard-notification');

  let allModels = [];

  // ── Helpers ───────────────────────────────────────────────────────────────

  function apiFetch(url, options = {}) {
    const defaults = {
      headers: { 'Content-Type': 'application/json' },
    };
    return fetch(url, Object.assign(defaults, options))
      .then(res => {
        if (res.status === 204) return null;
        return res.json().then(data => {
          if (!res.ok) throw new Error(data.error || 'Error en el servidor');
          return data;
        });
      });
  }

  function showNotification(type, message) {
    if (!notifArea) return;
    
    // PEC 3: Uso de templates para notificaciones
    const templateId = type === 'success' ? 'tpl-notify-success' : 'tpl-notify-error';
    const template = document.getElementById(templateId);
    if (!template) return;

    const clone = template.content.cloneNode(true);
    clone.querySelector('.notify-message').textContent = message;

    notifArea.innerHTML = '';
    notifArea.appendChild(clone);
    notifArea.hidden = false;
    
    setTimeout(() => { notifArea.hidden = true; }, 4000);
  }

  function fileToBase64(file) {
    return new Promise((resolve, reject) => {
      if (!file) return resolve("");
      const reader = new FileReader();
      reader.readAsDataURL(file);
      reader.onload = () => resolve(reader.result);
      reader.onerror = error => reject(error);
    });
  }

  // ── Gestión de Tabs ───────────────────────────────────────────────────────

  tabButtons.forEach(btn => {
    btn.addEventListener('click', () => {
      const target = btn.dataset.tab;

      tabButtons.forEach(b => b.classList.remove('active'));
      tabContents.forEach(c => {
        c.classList.remove('active');
        c.hidden = true;
      });

      btn.classList.add('active');
      const content = document.getElementById(target);
      if (content) {
        content.classList.add('active');
        content.hidden = false;
      }
    });
  });

  // ── Lógica de MODELOS ─────────────────────────────────────────────────────

  function loadModels() {
    apiFetch(API_MODELS)
      .then(models => {
        allModels = models || [];
        renderModelsTable(allModels);
        updateFlavorModelSelect(allModels);
      })
      .catch(err => showNotification('error', 'Error al cargar modelos: ' + err.message));
  }

  /**
   * PEC 3: Refactorizada para usar <template> tpl-model-row
   */
  function renderModelsTable(models) {
    if (!modelTbody) return;
    modelTbody.innerHTML = '';
    const template = document.getElementById('tpl-model-row');

    models.forEach(m => {
      const clone = template.content.cloneNode(true);
      
      // Relleno de valores usando querySelector y textContent
      clone.querySelector('.model-name strong').textContent = m.name;
      clone.querySelector('.model-subtitle').textContent = m.subtitle;
      clone.querySelector('.model-description small').textContent = m.description;
      clone.querySelector('.model-price').textContent = m.price.toFixed(2) + ' €';
      
      // Asignación de evento al botón de eliminar
      clone.querySelector('.btn-delete').addEventListener('click', () => deleteModel(m.id));

      modelTbody.appendChild(clone);
    });
  }

  if (modelForm) {
    modelForm.addEventListener('submit', (e) => {
      e.preventDefault();
      const data = {
        name: document.getElementById('model-name').value,
        subtitle: document.getElementById('model-subtitle').value,
        description: document.getElementById('model-description').value,
        price: parseFloat(document.getElementById('model-price').value)
      };

      apiFetch(API_MODELS, { method: 'POST', body: JSON.stringify(data) })
        .then(() => {
          showNotification('success', 'Modelo guardado correctamente');
          modelForm.reset();
          loadModels();
        })
        .catch(err => showNotification('error', err.message));
    });
  }

  function deleteModel(id) {
    if (!confirm('¿Eliminar este modelo?')) return;
    apiFetch(`${API_MODELS}/${id}`, { method: 'DELETE' })
      .then(() => {
        showNotification('success', 'Modelo eliminado');
        loadModels();
      })
      .catch(err => showNotification('error', err.message));
  }

  // ── Lógica de SABORES ─────────────────────────────────────────────────────

  function updateFlavorModelSelect(models) {
    if (!flavorModelSelect) return;
    flavorModelSelect.innerHTML = '<option value="">Selecciona un modelo...</option>';
    models.forEach(m => {
      const opt = document.createElement('option');
      opt.value = m.id;
      opt.dataset.name = m.name;
      opt.textContent = m.name;
      flavorModelSelect.appendChild(opt);
    });
  }

  function loadFlavors() {
    apiFetch(API_FLAVORS)
      .then(flavors => renderFlavorsTable(flavors || []))
      .catch(err => showNotification('error', 'Error al cargar sabores: ' + err.message));
  }

  /**
   * PEC 3: Refactorizada para usar <template> tpl-flavor-row
   */
  function renderFlavorsTable(flavors) {
    if (!flavorTbody) return;
    flavorTbody.innerHTML = '';
    const template = document.getElementById('tpl-flavor-row');

    flavors.forEach(f => {
      const clone = template.content.cloneNode(true);
      
      clone.querySelector('.flavor-model-id').textContent = f.model_id;
      clone.querySelector('.flavor-name strong').textContent = f.name;
      
      const img = clone.querySelector('.flavor-image');
      img.src = f.image || '/assets/images/placeholder.png';
      img.alt = f.name;
      
      clone.querySelector('.btn-delete').addEventListener('click', () => deleteFlavor(f.id));

      flavorTbody.appendChild(clone);
    });
  }

  if (flavorForm) {
    // PEC 3: Eliminado async/await en favor de cadenas .then()
    flavorForm.addEventListener('submit', (e) => {
      e.preventDefault();
      const modelOpt = flavorModelSelect.options[flavorModelSelect.selectedIndex];
      if (!modelOpt || !modelOpt.value) {
        showNotification('error', 'Selecciona un modelo');
        return;
      }
      
      const imageFile = document.getElementById('flavor-image').files[0];
      
      fileToBase64(imageFile)
        .then(imageBase64 => {
          const data = {
            modelID: modelOpt.value,
            modelName: modelOpt.dataset.name,
            name: document.getElementById('flavor-name').value,
            image: imageBase64
          };
          return apiFetch(API_FLAVORS, { method: 'POST', body: JSON.stringify(data) });
        })
        .then(() => {
          showNotification('success', 'Sabor añadido');
          flavorForm.reset();
          loadFlavors();
        })
        .catch(err => showNotification('error', err.message));
    });
  }

  function deleteFlavor(id) {
    if (!confirm('¿Eliminar este sabor?')) return;
    apiFetch(`${API_FLAVORS}/${id}`, { method: 'DELETE' })
      .then(() => {
        showNotification('success', 'Sabor eliminado');
        loadFlavors();
      })
      .catch(err => showNotification('error', err.message));
  }

  // ── Lógica de PROMOCIONES ─────────────────────────────────────────────────

  function createPromoItemRow(selectedId = "") {
    const div = document.createElement('div');
    div.className = 'dashboard__promo-item-row';
    
    let options = '<option value="">Elige un modelo...</option>';
    allModels.forEach(m => {
      options += `<option value="${m.id}" data-name="${m.name}" ${m.id === selectedId ? 'selected' : ''}>${m.name}</option>`;
    });

    div.innerHTML = `
      <select class="form-input promo-item-select" required>${options}</select>
      <button type="button" class="btn btn-modern btn-sm btn-remove-item">
        <i class="fa-solid fa-minus"></i>
      </button>
    `;

    div.querySelector('.btn-remove-item').addEventListener('click', () => div.remove());
    return div;
  }

  if (btnAddPromoItem) {
    btnAddPromoItem.addEventListener('click', () => {
      if (promoItemsContainer) {
        promoItemsContainer.appendChild(createPromoItemRow());
      }
    });
  }

  function loadPromotions() {
    apiFetch(API_PROMOTIONS)
      .then(promos => renderPromotionsTable(promos || []))
      .catch(err => showNotification('error', 'Error al cargar promociones: ' + err.message));
  }

  /**
   * PEC 3: Refactorizada para usar <template> tpl-promotion-row
   */
  function renderPromotionsTable(promos) {
    if (!promoTbody) return;
    promoTbody.innerHTML = '';
    const template = document.getElementById('tpl-promotion-row');

    promos.forEach(p => {
      const clone = template.content.cloneNode(true);
      
      clone.querySelector('.promo-name strong').textContent = p.name;
      clone.querySelector('.promo-description small').textContent = p.description;
      clone.querySelector('.promo-price').textContent = p.price.toFixed(2) + ' €';
      clone.querySelector('.promo-units').textContent = p.units + ' uds.';
      
      clone.querySelector('.btn-edit').addEventListener('click', () => editPromotion(p));
      clone.querySelector('.btn-delete').addEventListener('click', () => deletePromotion(p.id));

      promoTbody.appendChild(clone);
    });
  }

  if (promoForm) {
    // PEC 3: Eliminado async/await en favor de cadenas .then()
    promoForm.addEventListener('submit', (e) => {
      e.preventDefault();
      const id = document.getElementById('promotion-id').value;
      const imageFile = document.getElementById('promo-image').files[0];

      fileToBase64(imageFile)
        .then(imageBase64 => {
          const items = [];
          document.querySelectorAll('.promo-item-select').forEach(sel => {
            const opt = sel.options[sel.selectedIndex];
            if (opt && opt.value) {
              items.push({ model_id: opt.value, model_name: opt.dataset.name });
            }
          });

          const data = {
            name: document.getElementById('promo-name').value,
            description: document.getElementById('promo-description').value,
            price: parseFloat(document.getElementById('promo-price').value),
            units: parseInt(document.getElementById('promo-units').value),
            image: imageBase64 || "",
            items: items
          };

          const method = id ? 'PUT' : 'POST';
          const url = id ? `${API_PROMOTIONS}/${id}` : API_PROMOTIONS;

          return apiFetch(url, { method, body: JSON.stringify(data) });
        })
        .then(() => {
          showNotification('success', 'Promoción guardada');
          promoForm.reset();
          document.getElementById('promotion-id').value = "";
          if (promoItemsContainer) promoItemsContainer.innerHTML = "";
          loadPromotions();
        })
        .catch(err => showNotification('error', err.message));
    });
  }

  function editPromotion(p) {
    const idField = document.getElementById('promotion-id');
    if (idField) idField.value = p.id;
    
    document.getElementById('promo-name').value = p.name;
    document.getElementById('promo-description').value = p.description;
    document.getElementById('promo-price').value = p.price;
    document.getElementById('promo-units').value = p.units;
    
    if (promoItemsContainer) {
      promoItemsContainer.innerHTML = "";
      if (p.items) {
        p.items.forEach(item => {
          promoItemsContainer.appendChild(createPromoItemRow(item.model_id));
        });
      }
    }
    
    const tabPromo = document.getElementById('tab-promotions');
    if (tabPromo) tabPromo.scrollIntoView({ behavior: 'smooth' });
  }

  function deletePromotion(id) {
    if (!confirm('¿Eliminar esta promoción?')) return;
    apiFetch(`${API_PROMOTIONS}/${id}`, { method: 'DELETE' })
      .then(() => {
        showNotification('success', 'Promoción eliminada');
        loadPromotions();
      })
      .catch(err => showNotification('error', err.message));
  }

  // ── Arranque ──────────────────────────────────────────────────────────────

  function init() {
    loadModels();
    loadFlavors();
    loadPromotions();
  }

  init();
});
