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
  const quantityInput = form.querySelector('[name="quantity"]');
  const quantity = quantityInput ? parseInt(quantityInput.value) : 1;

  // Recolectamos todos los sabores (para bundles pueden ser varios)
  const flavorSelects = form.querySelectorAll('select[name^="flavor"]');
  const flavors = Array.from(flavorSelects).map(s => s.value).filter(v => v !== "");

  const btn = form.querySelector('button[type="submit"]');
  const originalHtml = btn.innerHTML;

  btn.disabled = true;
  btn.innerHTML = '<i class="fa-solid fa-spinner fa-spin"></i> ...';

  cartApiFetch('/api/cart', {
    method: 'POST',
    body: JSON.stringify({
      product_id: productId,
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
  const row = btn.closest('.app-row'); // FIX: Antes buscaba 'tr'

  btn.disabled = true;
  row.style.opacity = '0.5';

  cartApiFetch(`/api/cart?product_id=${encodeURIComponent(productId)}`, {
    method: 'DELETE'
  })
    .then(function (data) {
      showGlobalNotify('success', 'Producto eliminado');
      if (data.cart.items.length === 0) {
        window.location.reload();
      } else {
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
  const qtySpan = btn.parentElement.querySelector('.qty-value');
  let currentQty = parseInt(qtySpan.textContent);

  const newQty = action === 'inc' ? currentQty + 1 : currentQty - 1;
  if (newQty < 1) return; // No permitimos bajar de 1 aquí (se borra con el botón trash)

  btn.disabled = true;

  cartApiFetch('/api/cart', {
    method: 'PATCH',
    body: JSON.stringify({ product_id: productId, quantity: newQty })
  })
    .then(function (data) {
      qtySpan.textContent = newQty;
      updateCartTotals(data.cart);
      showGlobalNotify('success', 'Cantidad actualizada');
      // Actualizamos el subtotal de la fila
      const item = data.cart.items.find(i => i.product.id === productId);
      const row = btn.closest('.app-row');
      const subtotalEl = row.querySelector('.item-subtotal');
      if (subtotalEl) {
        subtotalEl.textContent = item.subtotal.toFixed(2) + '€';
      }
    })
    .catch(function (err) {
      showGlobalNotify('error', err.message);
    })
    .finally(function () {
      btn.disabled = false;
    });
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
