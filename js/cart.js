/**
 * cart.js
 * Gestión dinámica del carrito usando la API REST.
 */

function cartApiFetch(url, options = {}) {
  const defaults = {
    headers: { 'Content-Type': 'application/json' },
  };

  return fetch(url, Object.assign(defaults, options))
    .then(function (response) {
      if (response.status === 204) return null;
      return response.json().then(function (data) {
        if (!response.ok) throw new Error(data.error || 'Error en el carrito');
        return data;
      });
    });
}

// ── Manejo de Formularios en Portada y Detalle ──────────────────────────────

// Escuchamos todos los formularios de "Añadir al carrito" (clase .quick-add-form o id #add-to-cart-form)
document.addEventListener('submit', function (e) {
  const form = e.target.closest('.quick-add-form') || e.target.closest('#add-to-cart-form');
  if (!form) return;

  e.preventDefault();

  const productId = form.querySelector('[name="product_id"]').value;
  const flavorId = form.querySelector('[name="flavor_id"]')?.value;
  const flavorName = form.querySelector('[name="flavor_name"]')?.value;
  const price = form.querySelector('[name="price"]')?.value;
  const image = form.querySelector('[name="image"]')?.value;

  const quantityInput = form.querySelector('[name="quantity"]');
  const quantity = quantityInput ? parseInt(quantityInput.value) : 1;

  // Recolectamos todos los sabores (para compatibilidad con el sistema anterior si existe)
  const flavorSelects = form.querySelectorAll('select[name^="flavor"]');
  const flavors = Array.from(flavorSelects).map(s => s.value).filter(v => v !== "");
  if (flavorName && flavors.length === 0) flavors.push(flavorName);

  const btn = form.querySelector('button[type="submit"]');
  const originalHtml = btn.innerHTML;

  btn.disabled = true;
  btn.innerHTML = '<i class="fa-solid fa-spinner fa-spin"></i> ...';

  cartApiFetch('/cart/add', {
    method: 'POST',
    body: JSON.stringify({
      product_id: productId,
      flavor_id: flavorId,
      flavor_name: flavorName,
      price: parseFloat(price),
      image: image,
      quantity: quantity,
      flavors: flavors
    })
  })
    .then(function (data) {
      showGlobalNotify('success', '¡Añadido al carrito! 💨');
      animateCartIcon();
    })
    .catch(function (err) {
      showGlobalNotify('error', err.message);
    })
    .finally(function () {
      btn.disabled = false;
      btn.innerHTML = originalHtml;
    });
});

// ── Manejo de la página de Carrito ──────────────────────────────────────────

const cartContainer = document.getElementById('cart-table');
if (cartContainer) {
  cartContainer.addEventListener('click', function (e) {
    // 1. Borrado de items
    const removeBtn = e.target.closest('.btn-remove-item');
    if (removeBtn) {
      handleRemove(removeBtn);
      return;
    }

    // 2. Cambio de cantidad (+ / -)
    const qtyBtn = e.target.closest('.btn-qty');
    if (qtyBtn) {
      handleQtyChange(qtyBtn);
      return;
    }
  });
}

function handleRemove(btn) {
  const productId = btn.dataset.id;
  const row = btn.closest('.app-row');

  btn.disabled = true;
  row.style.opacity = '0.5';

  cartApiFetch(`/api/cart?product_id=${encodeURIComponent(productId)}`, {
    method: 'DELETE'
  })
    .then(function (data) {
      showGlobalNotify('success', 'Producto eliminado');
      // Verificamos si el carrito quedó vacío
      if (!data.cart || !data.cart.items || data.cart.items.length === 0) {
        window.location.reload();
      } else {
        // Animación de salida y eliminación del DOM
        row.style.transform = 'translateX(20px)';
        row.style.opacity = '0';
        setTimeout(() => {
          row.remove();
          updateCartTotals(data.cart);
        }, 300);
      }
    })
    .catch(function (err) {
      showGlobalNotify('error', err.message);
      btn.disabled = false;
      row.style.opacity = '1';
    });
}

function handleQtyChange(btn) {
  const productId = btn.dataset.id;
  const action = btn.dataset.action; // 'inc' o 'dec'
  const row = btn.closest('.app-row');
  const qtySpan = row.querySelector('.qty-value');
  let currentQty = parseInt(qtySpan.textContent);

  const newQty = action === 'inc' ? currentQty + 1 : currentQty - 1;
  if (newQty < 1) return; 

  btn.disabled = true;

  cartApiFetch('/api/cart', {
    method: 'PATCH',
    body: JSON.stringify({ product_id: productId, quantity: newQty })
  })
    .then(function (data) {
      qtySpan.textContent = newQty;
      updateCartTotals(data.cart);
      showGlobalNotify('success', 'Cantidad actualizada');

      // Bug Fix: Calculamos el subtotal leyendo el precio base del DOM 
      // para evitar errores si el backend no devuelve el item esperado.
      const priceText = row.querySelector('.qty-multiplier').textContent; // ej: "x 15.00€"
      const unitPrice = parseFloat(priceText.replace(/[^\d.]/g, ''));
      
      const subtotalEl = row.querySelector('.item-subtotal');
      if (subtotalEl && !isNaN(unitPrice)) {
        const newSubtotal = unitPrice * newQty;
        subtotalEl.textContent = newSubtotal.toFixed(2) + '€';
      }
    })
    .catch(function (err) {
      showGlobalNotify('error', err.message);
    })
    .finally(function () {
      btn.disabled = false;
    });
}

// ── Manejo de Promociones (Modal) ───────────────────────────────────────────

const promoModal = document.getElementById('promo-modal');
if (promoModal) {
  const closeBtn = promoModal.querySelector('.modal__close');
  
  // Cerrar al pulsar X
  closeBtn.onclick = () => promoModal.style.display = 'none';
  
  // Cerrar al pulsar fuera del modal
  window.onclick = (e) => {
    if (e.target === promoModal) promoModal.style.display = 'none';
  };

  // Escuchamos clicks en botones de "Ver Pack" o "Añadir Promo"
  document.addEventListener('click', function(e) {
    const btn = e.target.closest('.btn-open-promo');
    if (btn) {
      const promoId = btn.dataset.promoId;
      openPromoModal(promoId);
    }
  });

  function openPromoModal(id) {
    const promo = window.PROMOTIONS.find(p => p.id === id);
    if (!promo) return;

    // Rellenamos datos básicos
    document.getElementById('modal-promo-image').src = promo.image || '/assets/images/placeholder.png';
    document.getElementById('modal-promo-name').textContent = promo.name;
    document.getElementById('modal-promo-price').textContent = promo.price.toFixed(2) + '€';
    
    // Inputs ocultos para el form
    document.getElementById('modal-promo-id').value = promo.id;
    document.getElementById('modal-promo-name-input').value = promo.name;
    document.getElementById('modal-promo-price-input').value = promo.price;
    document.getElementById('modal-promo-image-input').value = promo.image;

    // Generamos selectores de sabores
    const container = document.getElementById('modal-flavor-selectors');
    container.innerHTML = ''; // Limpiamos

    promo.items.forEach((item, index) => {
      const selectorDiv = document.createElement('div');
      selectorDiv.className = 'modal__selector-item';
      
      const label = document.createElement('label');
      label.textContent = `Unidad ${index + 1}: ${item.model_name}`;
      
      const select = document.createElement('select');
      select.name = `flavor_${index}`;
      select.required = true;
      
      // Buscamos sabores de ese modelo
      const flavors = window.VAPE_FLAVORS.filter(f => f.model_id === item.model_id);
      
      if (flavors.length === 0) {
        const opt = document.createElement('option');
        opt.textContent = 'Sin sabores disponibles';
        select.appendChild(opt);
      } else {
        flavors.forEach(f => {
          const opt = document.createElement('option');
          opt.value = f.name;
          opt.textContent = f.name;
          select.appendChild(opt);
        });
      }
      
      selectorDiv.appendChild(label);
      selectorDiv.appendChild(select);
      container.appendChild(selectorDiv);
    });

    promoModal.style.display = 'block';
  }

  // Manejo del submit del modal
  document.getElementById('promo-add-form').onsubmit = function(e) {
    e.preventDefault();
    const form = e.target;
    const btn = form.querySelector('button[type="submit"]');
    const originalHtml = btn.innerHTML;

    const productId = document.getElementById('modal-promo-id').value;
    const promoName = document.getElementById('modal-promo-name-input').value;
    const price = parseFloat(document.getElementById('modal-promo-price-input').value);
    const image = document.getElementById('modal-promo-image-input').value;

    // Recogemos todos los sabores elegidos
    const selects = form.querySelectorAll('select');
    const flavors = Array.from(selects).map(s => s.value);

    btn.disabled = true;
    btn.innerHTML = '<i class="fa-solid fa-spinner fa-spin"></i> Añadiendo...';

    cartApiFetch('/cart/add', {
      method: 'POST',
      body: JSON.stringify({
        product_id: productId,
        flavor_id: 'promo-' + productId, // ID especial para promos
        flavor_name: promoName,
        price: price,
        image: image,
        quantity: 1,
        flavors: flavors
      })
    })
    .then(data => {
      showGlobalNotify('success', '¡Pack añadido al carrito! 🎁');
      promoModal.style.display = 'none';
      animateCartIcon();
    })
    .catch(err => {
      showGlobalNotify('error', err.message);
    })
    .finally(() => {
      btn.disabled = false;
      btn.innerHTML = originalHtml;
    });
  };
}

// ── Utilidades de Notificación (Toasts) ──────────────────────────────────────

function showGlobalNotify(type, message) {
  const container = document.getElementById('toast-container');
  if (!container) return;

  const toast = document.createElement('div');
  toast.className = `toast toast--${type}`;

  const icon = type === 'success' ? 'fa-circle-check' : 'fa-circle-exclamation';

  toast.innerHTML = `
    <i class="fa-solid ${icon}"></i>
    <span class="toast__message">${message}</span>
  `;

  container.appendChild(toast);

  // Animación de salida y eliminación
  setTimeout(() => {
    toast.classList.add('toast--out');
    setTimeout(() => {
      toast.remove();
    }, 500);
  }, 3000);
}

function updateCartTotals(cart) {
  const totalEl = document.getElementById('cart-total-amount');
  if (totalEl) {
    totalEl.textContent = cart.total.toFixed(2) + ' €';
  }
}

function animateCartIcon() {
  const icon = document.querySelector('.fa-cart-shopping');
  if (icon) {
    icon.style.transform = 'scale(1.3) rotate(-10deg)';
    setTimeout(() => { icon.style.transform = ''; }, 200);
  }
}
