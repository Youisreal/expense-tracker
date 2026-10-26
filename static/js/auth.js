/**
 * auth.js - Logika Autentikasi Frontend (Register, Login, Logout, Sesi)
 */

// Menampilkan pesan alert
function showAlert(elementId, message, isSuccess = false) {
  const el = document.getElementById(elementId);
  if (!el) return;
  el.textContent = message;
  el.className = `alert ${isSuccess ? 'alert-success' : 'alert-danger'}`;
  el.style.display = 'block';
}

// Menyembunyikan pesan alert
function hideAlert(elementId) {
  const el = document.getElementById(elementId);
  if (el) el.style.display = 'none';
}

// Inisialisasi Form Registrasi
function initRegisterForm() {
  const form = document.getElementById('registerForm');
  if (!form) return;

  form.addEventListener('submit', async (e) => {
    e.preventDefault();
    hideAlert('registerAlert');

    const username = document.getElementById('username').value.trim();
    const email = document.getElementById('email').value.trim();
    const password = document.getElementById('password').value;
    const confirmPassword = document.getElementById('confirmPassword').value;
    const submitBtn = document.getElementById('registerBtn');

    // Validasi input klien
    if (username.length < 3) {
      showAlert('registerAlert', 'Username minimal 3 karakter');
      return;
    }
    if (password.length < 6) {
      showAlert('registerAlert', 'Password minimal 6 karakter');
      return;
    }
    if (password !== confirmPassword) {
      showAlert('registerAlert', 'Konfirmasi password tidak cocok');
      return;
    }

    try {
      submitBtn.disabled = true;
      submitBtn.textContent = 'Mendaftar...';

      const response = await fetch('/api/auth/register', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ username, email, password }),
      });

      const result = await response.json();

      if (!response.ok) {
        throw new Error(result.error || 'Gagal melakukan pendaftaran');
      }

      showAlert('registerAlert', 'Pendaftaran berhasil! Mengalihkan ke halaman login...', true);
      setTimeout(() => {
        window.location.href = '/login';
      }, 1500);
    } catch (err) {
      showAlert('registerAlert', err.message);
    } finally {
      submitBtn.disabled = false;
      submitBtn.textContent = 'Daftar Sekarang';
    }
  });
}

// Inisialisasi Form Login
function initLoginForm() {
  const form = document.getElementById('loginForm');
  if (!form) return;

  form.addEventListener('submit', async (e) => {
    e.preventDefault();
    hideAlert('loginAlert');

    const identifier = document.getElementById('identifier').value.trim();
    const password = document.getElementById('password').value;
    const submitBtn = document.getElementById('loginBtn');

    if (!identifier || !password) {
      showAlert('loginAlert', 'Silakan isi Email/Username dan Password');
      return;
    }

    try {
      submitBtn.disabled = true;
      submitBtn.textContent = 'Memproses...';

      const response = await fetch('/api/auth/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ identifier, password }),
      });

      const result = await response.json();

      if (!response.ok) {
        throw new Error(result.error || 'Email/Username atau Password salah');
      }

      showAlert('loginAlert', 'Login berhasil! Mengalihkan ke Dashboard...', true);
      setTimeout(() => {
        window.location.href = '/dashboard';
      }, 800);
    } catch (err) {
      showAlert('loginAlert', err.message);
    } finally {
      submitBtn.disabled = false;
      submitBtn.textContent = 'Masuk';
    }
  });
}

// Logout Pengguna
async function logoutUser() {
  try {
    await fetch('/api/auth/logout', { method: 'POST' });
  } catch (err) {
    console.error('Logout error:', err);
  } finally {
    window.location.href = '/login';
  }
}

// Memeriksa status login pengguna saat membuka halaman
async function checkAuth(options = { redirectIfLoggedIn: false, redirectIfNotLoggedIn: false }) {
  try {
    const res = await fetch('/api/auth/me');
    if (res.ok) {
      const data = await res.json();
      if (options.redirectIfLoggedIn) {
        window.location.href = '/dashboard';
      }
      return data.user;
    } else {
      if (options.redirectIfNotLoggedIn) {
        window.location.href = '/login';
      }
      return null;
    }
  } catch (err) {
    if (options.redirectIfNotLoggedIn) {
      window.location.href = '/login';
    }
    return null;
  }
}

