/**
 * categories.js - Pengelolaan Kategori Bawaan & Custom Kategori
 */

let allCategories = [];

// Mengambil seluruh kategori dari backend
async function fetchCategories(type = '') {
  try {
    const url = type ? `/api/categories?type=${type}` : '/api/categories';
    const res = await fetch(url);
    if (!res.ok) throw new Error('Gagal mengambil data kategori');
    const result = await res.json();
    allCategories = result.data || [];
    return allCategories;
  } catch (err) {
    console.error('Error fetching categories:', err);
    return [];
  }
}

// Mengisi dropdown kategori berdasarkan tipe transaksi (income/expense)
function populateCategorySelect(selectId, type, selectedId = null) {
  const select = document.getElementById(selectId);
  if (!select) return;

  select.innerHTML = '<option value="">-- Pilih Kategori --</option>';

  const filtered = allCategories.filter(c => !type || c.type === type);

  // Pisahkan bawaan dan kustom untuk kejelasan pengguna
  const defaultCats = filtered.filter(c => c.is_default);
  const customCats = filtered.filter(c => !c.is_default);

  if (defaultCats.length > 0) {
    const defaultGroup = document.createElement('optgroup');
    defaultGroup.label = 'Kategori Bawaan';
    defaultCats.forEach(cat => {
      const opt = document.createElement('option');
      opt.value = cat.id;
      opt.textContent = cat.name;
      if (selectedId && Number(selectedId) === cat.id) opt.selected = true;
      defaultGroup.appendChild(opt);
    });
    select.appendChild(defaultGroup);
  }

  if (customCats.length > 0) {
    const customGroup = document.createElement('optgroup');
    customGroup.label = 'Kategori Custom Saya';
    customCats.forEach(cat => {
      const opt = document.createElement('option');
      opt.value = cat.id;
      opt.textContent = cat.name;
      if (selectedId && Number(selectedId) === cat.id) opt.selected = true;
      customGroup.appendChild(opt);
    });
    select.appendChild(customGroup);
  }
}

// Buka Modal Kelola Kategori
function openCategoryModal() {
  document.getElementById('categoryModal').classList.add('show');
  renderCategoryList();
}

// Tutup Modal Kelola Kategori
function closeCategoryModal() {
  document.getElementById('categoryModal').classList.remove('show');
  document.getElementById('categoryForm').reset();
  document.getElementById('editCategoryId').value = '';
  document.getElementById('categoryFormSubmitBtn').textContent = 'Tambah Kategori';
  hideAlert('categoryAlert');
}

// Render daftar kategori di dalam modal
async function renderCategoryList() {
  const container = document.getElementById('categoryListContainer');
  if (!container) return;

  container.innerHTML = '<p style="text-align: center; color: var(--text-muted); font-size: 0.9rem;">Memuat kategori...</p>';
  await fetchCategories();

  if (allCategories.length === 0) {
    container.innerHTML = '<p style="text-align: center; color: var(--text-muted);">Belum ada kategori.</p>';
    return;
  }

  let html = '';
  allCategories.forEach(cat => {
    const typeBadge = cat.type === 'income' 
      ? '<span class="badge badge-income">Pemasukan</span>' 
      : '<span class="badge badge-expense">Pengeluaran</span>';

    const defaultBadge = cat.is_default 
      ? '<span class="badge badge-default" title="Kategori bawaan sistem">Bawaan</span>' 
      : '<span class="badge" style="background: #e0e7ff; color: #3730a3;">Custom</span>';

    let actionButtons = '';
    if (!cat.is_default) {
      actionButtons = `
        <div style="display: flex; gap: 0.4rem;">
          <button class="btn btn-outline btn-sm" onclick="startEditCategory(${cat.id}, '${escapeHtml(cat.name)}', '${cat.type}')">Edit</button>
          <button class="btn btn-danger-outline btn-sm" onclick="deleteCategory(${cat.id})">Hapus</button>
        </div>
      `;
    } else {
      actionButtons = `<span style="font-size: 0.8rem; color: var(--text-muted);">Read-only</span>`;
    }

    html += `
      <div class="category-item">
        <div style="display: flex; align-items: center; gap: 0.6rem;">
          <strong>${escapeHtml(cat.name)}</strong>
          ${typeBadge}
          ${defaultBadge}
        </div>
        <div>
          ${actionButtons}
        </div>
      </div>
    `;
  });

  container.innerHTML = html;
}

// Mulai proses edit kategori custom
function startEditCategory(id, name, type) {
  document.getElementById('editCategoryId').value = id;
  document.getElementById('categoryName').value = name;
  document.getElementById('categoryType').value = type;
  document.getElementById('categoryFormSubmitBtn').textContent = 'Perbarui Kategori';
  document.getElementById('categoryName').focus();
}

// Batalkan mode edit
function cancelEditCategory() {
  document.getElementById('editCategoryId').value = '';
  document.getElementById('categoryForm').reset();
  document.getElementById('categoryFormSubmitBtn').textContent = 'Tambah Kategori';
}

// Simpan kategori (Tambah / Update)
async function handleCategorySubmit(e) {
  e.preventDefault();
  hideAlert('categoryAlert');

  const id = document.getElementById('editCategoryId').value;
  const name = document.getElementById('categoryName').value.trim();
  const type = document.getElementById('categoryType').value;
  const submitBtn = document.getElementById('categoryFormSubmitBtn');

  if (!name) {
    showAlert('categoryAlert', 'Nama kategori wajib diisi');
    return;
  }

  try {
    submitBtn.disabled = true;
    let url = '/api/categories';
    let method = 'POST';

    if (id) {
      url = `/api/categories/${id}`;
      method = 'PUT';
    }

    const res = await fetch(url, {
      method: method,
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ name, type }),
    });

    const result = await res.json();
    if (!res.ok) {
      throw new Error(result.error || 'Gagal menyimpan kategori');
    }

    cancelEditCategory();
    await renderCategoryList();
    // Perbarui dropdown form transaksi jika sedang terbuka
    if (typeof refreshTransactionCategories === 'function') {
      refreshTransactionCategories();
    }
  } catch (err) {
    showAlert('categoryAlert', err.message);
  } finally {
    submitBtn.disabled = false;
  }
}

// Hapus kategori custom
async function deleteCategory(id) {
  if (!confirm('Apakah Anda yakin ingin menghapus kategori custom ini?')) return;

  try {
    const res = await fetch(`/api/categories/${id}`, { method: 'DELETE' });
    const result = await res.json();
    if (!res.ok) {
      throw new Error(result.error || 'Gagal menghapus kategori');
    }
    await renderCategoryList();
    if (typeof refreshTransactionCategories === 'function') {
      refreshTransactionCategories();
    }
  } catch (err) {
    alert(err.message);
  }
}

// Helper escape HTML
function escapeHtml(str) {
  if (!str) return '';
  return str
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
    .replace(/'/g, '&#039;');
}

