/**
 * app.js - Logika Dashboard Keuangan, Filter, Transaksi, & Visualisasi Chart.js
 */

let currentMonth = new Date().getMonth() + 1; // 1-12
let currentYear = new Date().getFullYear();
let activeTransactionType = 'expense';
let comparisonChart = null;
let categoryChart = null;

// Format Rupiah
function formatRupiah(amount) {
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    maximumFractionDigits: 0,
  }).format(amount || 0);
}

// Format Tanggal
function formatDate(dateStr) {
  if (!dateStr) return '-';
  const d = new Date(dateStr);
  return d.toLocaleDateString('id-ID', {
    day: 'numeric',
    month: 'short',
    year: 'numeric',
  });
}

// Inisialisasi Dashboard
document.addEventListener('DOMContentLoaded', async () => {
  // Verifikasi login pengguna
  const user = await checkAuth({ redirectIfNotLoggedIn: true });
  if (!user) return;

  const navUserEl = document.getElementById('navUser');
  if (navUserEl) {
    navUserEl.innerHTML = `👤 <strong>${escapeHtml(user.username)}</strong>`;
  }

  // Setup dropdown filter periode (Bulan & Tahun)
  setupPeriodSelectors();

  // Load awal
  await fetchCategories();
  populateCategoryFilter();
  await refreshDashboard();

  // Event listener form transaksi
  const transForm = document.getElementById('transactionForm');
  if (transForm) {
    transForm.addEventListener('submit', handleTransactionSubmit);
  }

  // Event listener form kategori
  const catForm = document.getElementById('categoryForm');
  if (catForm) {
    catForm.addEventListener('submit', handleCategorySubmit);
  }
});

// Setup Selector Bulan & Tahun
function setupPeriodSelectors() {
  const monthSelect = document.getElementById('filterMonth');
  const yearSelect = document.getElementById('filterYear');

  if (monthSelect) {
    monthSelect.value = currentMonth;
    monthSelect.addEventListener('change', (e) => {
      currentMonth = e.target.value;
      refreshDashboard();
    });
  }

  if (yearSelect) {
    // Isi rentang tahun: 2 tahun lalu hingga 1 tahun ke depan
    yearSelect.innerHTML = '';
    const nowY = new Date().getFullYear();
    for (let y = nowY - 2; y <= nowY + 1; y++) {
      const opt = document.createElement('option');
      opt.value = y;
      opt.textContent = y;
      if (y === currentYear) opt.selected = true;
      yearSelect.appendChild(opt);
    }
    yearSelect.addEventListener('change', (e) => {
      currentYear = e.target.value;
      refreshDashboard();
    });
  }
}

// Refresh seluruh data dashboard
async function refreshDashboard() {
  await Promise.all([
    loadSummary(),
    loadTransactions()
  ]);
}

// Mengambil dan Menampilkan Ringkasan Keuangan (Stats & Grafik)
async function loadSummary() {
  try {
    let url = `/api/transactions/summary`;
    if (currentMonth && currentYear) {
      url += `?month=${currentMonth}&year=${currentYear}`;
    }

    const res = await fetch(url);
    if (!res.ok) throw new Error('Gagal memuat ringkasan');
    const result = await res.json();
    const data = result.data;

    // Update Kartu Statistik
    document.getElementById('totalIncome').textContent = formatRupiah(data.total_income);
    document.getElementById('totalExpense').textContent = formatRupiah(data.total_expense);
    
    const balanceEl = document.getElementById('netBalance');
    balanceEl.textContent = formatRupiah(data.net_balance);
    if (data.net_balance < 0) {
      balanceEl.style.color = 'var(--danger)';
    } else {
      balanceEl.style.color = 'var(--primary)';
    }

    // Render Grafik Chart.js
    renderComparisonChart(data.total_income, data.total_expense);
    renderCategoryBreakdownChart(data.expense_categories || []);

  } catch (err) {
    console.error('Error loading summary:', err);
  }
}

// Render Grafik Perbandingan Pemasukan vs Pengeluaran (Bar Chart)
function renderComparisonChart(income, expense) {
  const ctx = document.getElementById('comparisonChart');
  if (!ctx) return;

  if (comparisonChart) {
    comparisonChart.destroy();
  }

  comparisonChart = new Chart(ctx, {
    type: 'bar',
    data: {
      labels: ['Pemasukan', 'Pengeluaran'],
      datasets: [{
        label: 'Total (Rp)',
        data: [income, expense],
        backgroundColor: [
          'rgba(16, 185, 129, 0.85)', // Success green
          'rgba(239, 68, 68, 0.85)',  // Danger red
        ],
        borderColor: [
          '#10b981',
          '#ef4444',
        ],
        borderWidth: 1,
        borderRadius: 8,
      }]
    },
    options: {
      responsive: true,
      maintainAspectRatio: false,
      plugins: {
        legend: { display: false },
        tooltip: {
          callbacks: {
            label: (ctx) => ` ${formatRupiah(ctx.raw)}`
          }
        }
      },
      scales: {
        y: {
          beginAtZero: true,
          ticks: {
            callback: (val) => 'Rp ' + (val >= 1000000 ? (val / 1000000) + ' Jt' : (val / 1000) + ' Rb')
          }
        }
      }
    }
  });
}

// Render Grafik Proporsi Kategori Pengeluaran (Doughnut Chart)
function renderCategoryBreakdownChart(categories) {
  const ctx = document.getElementById('categoryChart');
  if (!ctx) return;

  if (categoryChart) {
    categoryChart.destroy();
  }

  const noDataEl = document.getElementById('categoryChartNoData');

  if (!categories || categories.length === 0) {
    ctx.style.display = 'none';
    if (noDataEl) noDataEl.style.display = 'flex';
    return;
  }

  ctx.style.display = 'block';
  if (noDataEl) noDataEl.style.display = 'none';

  const labels = categories.map(c => c.name);
  const data = categories.map(c => c.total);
  const colors = [
    '#3b82f6', '#ef4444', '#10b981', '#f59e0b', '#8b5cf6',
    '#ec4899', '#14b8a6', '#f97316', '#6366f1', '#84cc16'
  ];

  categoryChart = new Chart(ctx, {
    type: 'doughnut',
    data: {
      labels: labels,
      datasets: [{
        data: data,
        backgroundColor: colors.slice(0, labels.length),
        borderWidth: 2,
        borderColor: '#ffffff',
      }]
    },
    options: {
      responsive: true,
      maintainAspectRatio: false,
      plugins: {
        legend: {
          position: 'right',
          labels: { boxWidth: 12, font: { size: 11 } }
        },
        tooltip: {
          callbacks: {
            label: (ctx) => ` ${ctx.label}: ${formatRupiah(ctx.raw)}`
          }
        }
      },
      cutout: '65%'
    }
  });
}

// Mengambil dan Menampilkan Riwayat Transaksi dengan Filter
async function loadTransactions() {
  const tbody = document.getElementById('transactionTableBody');
  const emptyState = document.getElementById('transactionEmptyState');
  if (!tbody) return;

  const typeFilter = document.getElementById('filterType').value;
  const categoryFilter = document.getElementById('filterCategory').value;

  let url = `/api/transactions?`;
  const params = [];
  if (currentMonth && currentYear) {
    params.push(`month=${currentMonth}`, `year=${currentYear}`);
  }
  if (typeFilter) params.push(`type=${typeFilter}`);
  if (categoryFilter) params.push(`category_id=${categoryFilter}`);

  url += params.join('&');

  try {
    tbody.innerHTML = '<tr><td colspan="6" style="text-align: center; color: var(--text-muted); padding: 2rem;">Memuat transaksi...</td></tr>';
    
    const res = await fetch(url);
    if (!res.ok) throw new Error('Gagal memuat transaksi');
    const result = await res.json();
    const transactions = result.data || [];

    if (transactions.length === 0) {
      tbody.innerHTML = '';
      if (emptyState) emptyState.style.display = 'block';
      return;
    }

    if (emptyState) emptyState.style.display = 'none';

    let html = '';
    transactions.forEach(t => {
      const isExpense = t.type === 'expense';
      const badgeClass = isExpense ? 'badge-expense' : 'badge-income';
      const amountClass = isExpense ? 'amount-expense' : 'amount-income';
      const typeLabel = isExpense ? 'Pengeluaran' : 'Pemasukan';
      const sign = isExpense ? '-' : '+';
      const categoryName = t.category ? t.category.name : '-';
      const dateDisplay = formatDate(t.date);

      html += `
        <tr>
          <td><strong>${dateDisplay}</strong></td>
          <td>
            <span class="badge ${badgeClass}" style="margin-right: 0.35rem;">${typeLabel}</span>
            ${escapeHtml(categoryName)}
          </td>
          <td style="color: var(--text-muted);">${escapeHtml(t.description || '-')}</td>
          <td class="${amountClass}">
            ${sign} ${formatRupiah(t.amount)}
          </td>
          <td>
            <div style="display: flex; gap: 0.4rem;">
              <button class="btn btn-outline btn-sm" onclick="openEditTransaction(${t.id})" title="Edit">✏️ Edit</button>
              <button class="btn btn-danger-outline btn-sm" onclick="deleteTransaction(${t.id})" title="Hapus">🗑️ Hapus</button>
            </div>
          </td>
        </tr>
      `;
    });

    tbody.innerHTML = html;
  } catch (err) {
    console.error('Error loading transactions:', err);
    tbody.innerHTML = `<tr><td colspan="6" style="text-align: center; color: var(--danger); padding: 1.5rem;">Terjadi kesalahan memuat data.</td></tr>`;
  }
}

// Buka Modal Tambah Transaksi
function openAddTransactionModal() {
  document.getElementById('transactionModalTitle').textContent = 'Tambah Transaksi';
  document.getElementById('transactionForm').reset();
  document.getElementById('editTransactionId').value = '';
  document.getElementById('transactionDate').valueAsDate = new Date();
  setTransactionType('expense');
  hideAlert('transactionAlert');
  document.getElementById('transactionModal').classList.add('show');
}

// Buka Modal Edit Transaksi
async function openEditTransaction(id) {
  try {
    const res = await fetch(`/api/transactions/${id}`);
    if (!res.ok) throw new Error('Gagal mengambil data transaksi');
    const result = await res.json();
    const t = result.data;

    document.getElementById('transactionModalTitle').textContent = 'Edit Transaksi';
    document.getElementById('editTransactionId').value = t.id;
    document.getElementById('transactionAmount').value = t.amount;
    document.getElementById('transactionDescription').value = t.description || '';
    
    // Set tanggal format YYYY-MM-DD
    const dateFormatted = t.date ? t.date.split('T')[0] : '';
    document.getElementById('transactionDate').value = dateFormatted;

    setTransactionType(t.type);
    populateCategorySelect('transactionCategory', t.type, t.category_id);

    hideAlert('transactionAlert');
    document.getElementById('transactionModal').classList.add('show');
  } catch (err) {
    alert(err.message);
  }
}

// Tutup Modal Transaksi
function closeTransactionModal() {
  document.getElementById('transactionModal').classList.remove('show');
  document.getElementById('transactionForm').reset();
}

// Ubah Tipe Transaksi di Form Modal (Pemasukan / Pengeluaran)
function setTransactionType(type) {
  activeTransactionType = type;
  document.getElementById('transactionType').value = type;

  const btnExpense = document.getElementById('typeBtnExpense');
  const btnIncome = document.getElementById('typeBtnIncome');

  if (type === 'expense') {
    btnExpense.className = 'type-toggle-btn active-expense';
    btnIncome.className = 'type-toggle-btn';
  } else {
    btnIncome.className = 'type-toggle-btn active-income';
    btnExpense.className = 'type-toggle-btn';
  }

  populateCategorySelect('transactionCategory', type);
}

// Submit Form Transaksi (Tambah atau Edit)
async function handleTransactionSubmit(e) {
  e.preventDefault();
  hideAlert('transactionAlert');

  const id = document.getElementById('editTransactionId').value;
  const amount = parseInt(document.getElementById('transactionAmount').value, 10);
  const type = document.getElementById('transactionType').value;
  const categoryId = parseInt(document.getElementById('transactionCategory').value, 10);
  const date = document.getElementById('transactionDate').value;
  const description = document.getElementById('transactionDescription').value.trim();
  const submitBtn = document.getElementById('transactionSubmitBtn');

  if (!amount || amount <= 0) {
    showAlert('transactionAlert', 'Nominal harus lebih dari 0');
    return;
  }
  if (!categoryId) {
    showAlert('transactionAlert', 'Silakan pilih kategori transaksi');
    return;
  }
  if (!date) {
    showAlert('transactionAlert', 'Silakan pilih tanggal transaksi');
    return;
  }

  try {
    submitBtn.disabled = true;
    submitBtn.textContent = 'Menyimpan...';

    const payload = {
      amount,
      type,
      category_id: categoryId,
      date,
      description
    };

    let url = '/api/transactions';
    let method = 'POST';

    if (id) {
      url = `/api/transactions/${id}`;
      method = 'PUT';
    }

    const res = await fetch(url, {
      method: method,
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload),
    });

    const result = await res.json();
    if (!res.ok) {
      throw new Error(result.error || 'Gagal menyimpan transaksi');
    }

    closeTransactionModal();
    await refreshDashboard();
  } catch (err) {
    showAlert('transactionAlert', err.message);
  } finally {
    submitBtn.disabled = false;
    submitBtn.textContent = 'Simpan Transaksi';
  }
}

// Hapus Transaksi
async function deleteTransaction(id) {
  if (!confirm('Apakah Anda yakin ingin menghapus catatan transaksi ini?')) return;

  try {
    const res = await fetch(`/api/transactions/${id}`, { method: 'DELETE' });
    const result = await res.json();
    if (!res.ok) {
      throw new Error(result.error || 'Gagal menghapus transaksi');
    }
    await refreshDashboard();
  } catch (err) {
    alert(err.message);
  }
}

// Mengisi dropdown filter kategori di atas tabel riwayat
function populateCategoryFilter() {
  const filterCat = document.getElementById('filterCategory');
  if (!filterCat) return;

  filterCat.innerHTML = '<option value="">Semua Kategori</option>';
  allCategories.forEach(cat => {
    const opt = document.createElement('option');
    opt.value = cat.id;
    opt.textContent = `${cat.name} (${cat.type === 'income' ? 'Pemasukan' : 'Pengeluaran'})`;
    filterCat.appendChild(opt);
  });
}

// Callback saat kategori bertambah/diubah di modal kategori
function refreshTransactionCategories() {
  populateCategoryFilter();
  if (document.getElementById('transactionModal').classList.contains('show')) {
    populateCategorySelect('transactionCategory', activeTransactionType);
  }
}

