/**
 * account.js
 * Gestión del panel de cuenta: actualización de perfil y cambio de contraseña.
 * Usa la API REST con Promesas (PEC 3).
 */

const accNotif = document.getElementById('account-notification');

function showAccNotify(type, message) {
  if (!accNotif) return;
  accNotif.className = 'dashboard__notification alert alert--' + (type === 'success' ? 'success' : 'error');
  accNotif.textContent = message;
  accNotif.hidden = false;
  window.scrollTo({ top: 0, behavior: 'smooth' });
  setTimeout(function () { accNotif.hidden = true; }, 4000);
}

function setButtonLoading(btn, loading, originalHTML) {
  btn.disabled = loading;
  btn.innerHTML = loading
    ? '<i class="fa-solid fa-spinner fa-spin"></i> Saving...'
    : originalHTML;
}

// ── Formulario de perfil ─────────────────────────────────────────────────────
var accForm = document.getElementById('account-form');
var accBtn  = document.getElementById('btn-save-account');

if (accForm && accBtn) {
  var accBtnOriginal = accBtn.innerHTML;

  accForm.addEventListener('submit', function (e) {
    e.preventDefault();

    var payload = {
      name:  document.getElementById('acc-name').value.trim(),
      email: document.getElementById('acc-email').value.trim(),
      phone: document.getElementById('acc-phone').value.trim()
    };

    setButtonLoading(accBtn, true);

    fetch('/api/account', {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload)
    })
    .then(function (res) {
      return res.json().then(function (data) {
        if (!res.ok) throw new Error(data.error || 'Error al actualizar');
        return data;
      });
    })
    .then(function () {
      showAccNotify('success', 'Profile updated successfully!');
    })
    .catch(function (err) {
      showAccNotify('error', err.message);
    })
    .finally(function () {
      setButtonLoading(accBtn, false, accBtnOriginal);
    });
  });
}

// ── Formulario de contraseña ─────────────────────────────────────────────────
var pwdForm = document.getElementById('password-form');
var pwdBtn  = document.getElementById('btn-save-password');

if (pwdForm && pwdBtn) {
  var pwdBtnOriginal = pwdBtn.innerHTML;

  pwdForm.addEventListener('submit', function (e) {
    e.preventDefault();

    var newPwd     = document.getElementById('pwd-new').value;
    var confirmPwd = document.getElementById('pwd-confirm').value;

    if (newPwd !== confirmPwd) {
      showAccNotify('error', 'New passwords do not match.');
      return;
    }

    var payload = {
      current_password: document.getElementById('pwd-current').value,
      new_password:     newPwd
    };

    setButtonLoading(pwdBtn, true);

    fetch('/api/account/password', {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload)
    })
    .then(function (res) {
      return res.json().then(function (data) {
        if (!res.ok) throw new Error(data.error || 'Error al cambiar contraseña');
        return data;
      });
    })
    .then(function () {
      showAccNotify('success', 'Password updated successfully!');
      pwdForm.reset();
    })
    .catch(function (err) {
      showAccNotify('error', err.message);
    })
    .finally(function () {
      setButtonLoading(pwdBtn, false, pwdBtnOriginal);
    });
  });
}
