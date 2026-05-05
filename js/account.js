/**
 * account.js
 * Actualización dinámica del perfil de usuario vía API.
 */

const accForm = document.getElementById('account-form');
const accBtn  = document.getElementById('btn-save-account');
const accNotif = document.getElementById('account-notification');

if (accForm) {
  accForm.addEventListener('submit', function (e) {
    e.preventDefault();

    const name  = document.getElementById('acc-name').value.trim();
    const email = document.getElementById('acc-email').value.trim();
    const phone = document.getElementById('acc-phone').value.trim();

    accBtn.disabled = true;
    accBtn.innerHTML = '<i class="fa-solid fa-spinner fa-spin"></i> Saving...';

    fetch('/api/account', {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ name: name, email: email, phone: phone })
    })
    .then(function (response) {
      return response.json().then(function (data) {
        if (!response.ok) throw new Error(data.error || 'Error al actualizar');
        return data;
      });
    })
    .then(function (data) {
      showAccNotify('success', '¡Perfil actualizado correctamente!');
    })
    .catch(function (err) {
      showAccNotify('error', err.message);
    })
    .finally(function () {
      accBtn.disabled = false;
      accBtn.innerHTML = '<i class="fa-solid fa-save"></i> Save Changes';
    });
  });
}

function showAccNotify(type, message) {
  if (!accNotif) return;
  accNotif.className = 'dashboard__notification alert alert--' + (type === 'success' ? 'success' : 'error');
  accNotif.textContent = message;
  accNotif.hidden = false;
  setTimeout(() => { accNotif.hidden = true; }, 3000);
}
