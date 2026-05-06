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
    const div = document.createElement('div');
    div.className = `alert alert-${type === 'success' ? 'success' : 'danger'}`;
    div.innerHTML = `<i class="fa-solid fa-${type === 'success' ? 'circle-check' : 'circle-exclamation'}"></i> <span>${message}</span>`;
    notifArea.innerHTML = '';
    notifArea.appendChild(div);
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

  function renderModelsTable(models) {
    if (!modelTbody) return;
    modelTbody.innerHTML = '';
    models.forEach(m => {
      const tr = document.createElement('tr');
      tr.innerHTML = `
        <td><strong>${m.name}</strong></td>
        <td>${m.subtitle}</td>
        <td><small>${m.description}</small></td>
        <td class="dashboard__product-price">${m.price.toFixed(2)} €</td>
        <td>
          <button class="btn btn-modern btn-sm" onclick="deleteModel('${m.id}')">
            <i class="fa-solid fa-trash"></i> Eliminar
          </button>
        </td>
      `;
      modelTbody.appendChild(tr);
    });
  }

  if (modelForm) {
    modelForm.addEventListener('submit', (e) => {
      e.preventDefault();
      const data = {
        name: document.getElementById('model-name').value,
        subtitle: document.getElementById('model-subtitle').value,
        description: document.getElementById('model-description').value,
        price: document.getElementById('model-price').value
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

  window.deleteModel = (id) => {
    if (!confirm('¿Eliminar este modelo?')) return;
    apiFetch(`${API_MODELS}/${id}`, { method: 'DELETE' })
      .then(() => {
        showNotification('success', 'Modelo eliminado');
        loadModels();
      })
      .catch(err => showNotification('error', err.message));
  };

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

  function renderFlavorsTable(flavors) {
    if (!flavorTbody) return;
    flavorTbody.innerHTML = '';
    flavors.forEach(f => {
      const tr = document.createElement('tr');
      tr.innerHTML = `
        <td><span class="dashboard__badge dashboard__badge--normal">${f.model_id}</span></td>
        <td><strong>${f.name}</strong></td>
        <td><img src="${f.image || '/assets/images/placeholder.png'}" class="dashboard__product-thumb"></td>
        <td>
          <button class="btn btn-modern btn-sm" onclick="deleteFlavor('${f.id}')">
            <i class="fa-solid fa-trash"></i> Eliminar
          </button>
        </td>
      `;
      flavorTbody.appendChild(tr);
    });
  }

  if (flavorForm) {
    flavorForm.addEventListener('submit', async (e) => {
      e.preventDefault();
      const modelOpt = flavorModelSelect.options[flavorModelSelect.selectedIndex];
      if (!modelOpt || !modelOpt.value) {
        showNotification('error', 'Selecciona un modelo');
        return;
      }
      const imageFile = document.getElementById('flavor-image').files[0];
      const imageBase64 = await fileToBase64(imageFile);

      const data = {
        modelID: modelOpt.value,
        modelName: modelOpt.dataset.name,
        name: document.getElementById('flavor-name').value,
        image: imageBase64
      };

      apiFetch(API_FLAVORS, { method: 'POST', body: JSON.stringify(data) })
        .then(() => {
          showNotification('success', 'Sabor añadido');
          flavorForm.reset();
          loadFlavors();
        })
        .catch(err => showNotification('error', err.message));
    });
  }

  window.deleteFlavor = (id) => {
    if (!confirm('¿Eliminar este sabor?')) return;
    apiFetch(`${API_FLAVORS}/${id}`, { method: 'DELETE' })
      .then(() => {
        showNotification('success', 'Sabor eliminado');
        loadFlavors();
      })
      .catch(err => showNotification('error', err.message));
  };

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

  function renderPromotionsTable(promos) {
    if (!promoTbody) return;
    promoTbody.innerHTML = '';
    promos.forEach(p => {
      const tr = document.createElement('tr');
      tr.innerHTML = `
        <td><strong>${p.name}</strong></td>
        <td><small>${p.description}</small></td>
        <td class="dashboard__product-price">${p.price.toFixed(2)} €</td>
        <td>${p.units} uds.</td>
        <td>
          <div class="dashboard__cell--actions">
            <button class="btn btn-secondary btn-sm" onclick='editPromotion(${JSON.stringify(p)})'>
              <i class="fa-solid fa-pen"></i>
            </button>
            <button class="btn btn-modern btn-sm" onclick="deletePromotion('${p.id}')">
              <i class="fa-solid fa-trash"></i>
            </button>
          </div>
        </td>
      `;
      promoTbody.appendChild(tr);
    });
  }

  if (promoForm) {
    promoForm.addEventListener('submit', async (e) => {
      e.preventDefault();
      const id = document.getElementById('promotion-id').value;
      const imageFile = document.getElementById('promo-image').files[0];
      const imageBase64 = await fileToBase64(imageFile);

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

      apiFetch(url, { method, body: JSON.stringify(data) })
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

  window.editPromotion = (p) => {
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
  };

  window.deletePromotion = (id) => {
    if (!confirm('¿Eliminar esta promoción?')) return;
    apiFetch(`${API_PROMOTIONS}/${id}`, { method: 'DELETE' })
      .then(() => {
        showNotification('success', 'Promoción eliminada');
        loadPromotions();
      })
      .catch(err => showNotification('error', err.message));
  };

  // ── Arranque ──────────────────────────────────────────────────────────────

  function init() {
    loadModels();
    loadFlavors();
    loadPromotions();
  }

  init();
});
