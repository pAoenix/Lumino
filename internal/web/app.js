const state = {
  datasets: [],
  selectedId: null,
  loading: false,
};

const uploadForm = document.querySelector("#uploadForm");
const uploadStatus = document.querySelector("#uploadStatus");
const uploadButton = document.querySelector("#uploadButton");
const fileInput = document.querySelector("#fileInput");
const selectedFileName = document.querySelector("#selectedFileName");
const refreshButton = document.querySelector("#refreshButton");
const datasetList = document.querySelector("#datasetList");
const datasetCount = document.querySelector("#datasetCount");
const summaryDatasets = document.querySelector("#summaryDatasets");
const summaryRows = document.querySelector("#summaryRows");
const summaryColumns = document.querySelector("#summaryColumns");
const summarySize = document.querySelector("#summarySize");
const detailPanel = document.querySelector("#detailPanel");
const detailTitle = document.querySelector("#detailTitle");
const detailMeta = document.querySelector("#detailMeta");
const downloadLink = document.querySelector("#downloadLink");
const profileGrid = document.querySelector("#profileGrid");
const previewTable = document.querySelector("#previewTable");

fileInput.addEventListener("change", () => {
  const file = fileInput.files[0];
  selectedFileName.textContent = file ? `${file.name} · ${formatSize(file.size)}` : "单文件最大 100MB，CSV 支持预览";
});

uploadForm.addEventListener("submit", async (event) => {
  event.preventDefault();
  setUploadState(true, "正在上传...", "");

  try {
    const response = await fetch("/api/datasets", {
      method: "POST",
      body: new FormData(uploadForm),
    });
    const body = await response.json();
    if (!response.ok) {
      throw new Error(body.error || "上传失败");
    }
    uploadForm.reset();
    selectedFileName.textContent = "单文件最大 100MB，CSV 支持预览";
    setUploadState(false, "上传完成", "success");
    await loadDatasets(body.id);
  } catch (error) {
    setUploadState(false, error.message, "error");
  }
});

refreshButton.addEventListener("click", () => loadDatasets(state.selectedId));

async function loadDatasets(selectId) {
  state.loading = true;
  refreshButton.disabled = true;
  datasetList.innerHTML = '<p class="empty">正在读取...</p>';

  try {
    const response = await fetch("/api/datasets");
    const datasets = await response.json();
    if (!response.ok) {
      throw new Error(datasets.error || "读取数据失败");
    }
    state.datasets = datasets;
    state.selectedId = selectId || state.datasets[0]?.id || null;
    renderSummary();
    renderList();
    if (state.selectedId) {
      await loadDetail(state.selectedId);
    } else {
      detailPanel.classList.add("hidden");
      clearPreview();
    }
  } catch (error) {
    datasetList.innerHTML = `<p class="empty">${escapeHTML(error.message)}</p>`;
  } finally {
    state.loading = false;
    refreshButton.disabled = false;
  }
}

function renderSummary() {
  const totals = state.datasets.reduce(
    (acc, item) => {
      acc.rows += item.rows || 0;
      acc.columns += item.columns?.length || 0;
      acc.size += item.size || 0;
      return acc;
    },
    { rows: 0, columns: 0, size: 0 },
  );

  datasetCount.textContent = `${state.datasets.length} 个数据集`;
  summaryDatasets.textContent = formatNumber(state.datasets.length);
  summaryRows.textContent = formatNumber(totals.rows);
  summaryColumns.textContent = formatNumber(totals.columns);
  summarySize.textContent = formatSize(totals.size);
}

function renderList() {
  datasetList.innerHTML = "";

  if (state.datasets.length === 0) {
    datasetList.innerHTML = '<p class="empty">暂无数据，先上传一个 CSV 试试。</p>';
    return;
  }

  for (const item of state.datasets) {
    const row = document.createElement("article");
    row.className = `dataset-item ${item.id === state.selectedId ? "active" : ""}`;
    row.innerHTML = `
      <button class="dataset-main" type="button">
        <span class="dataset-title">${escapeHTML(item.name)}</span>
        <span class="dataset-meta">${formatNumber(item.rows || 0)} 行 · ${formatNumber(item.columns?.length || 0)} 列 · ${formatDate(item.createdAt)}</span>
        <span class="dataset-desc">${escapeHTML(item.description || item.fileName)}</span>
      </button>
      <div class="dataset-actions">
        <span class="dataset-size">${formatSize(item.size)}</span>
        <button class="delete-button" type="button" aria-label="删除 ${escapeHTML(item.name)}">删除</button>
      </div>
    `;
    row.querySelector(".dataset-main").addEventListener("click", () => loadDetail(item.id));
    row.querySelector(".delete-button").addEventListener("click", (event) => deleteDataset(item, event.currentTarget));
    datasetList.appendChild(row);
  }
}

async function deleteDataset(item, button) {
  const confirmed = window.confirm(`确定删除“${item.name}”吗？此操作会同时删除本地文件，无法恢复。`);
  if (!confirmed) {
    return;
  }

  button.disabled = true;
  button.textContent = "删除中";

  try {
    const response = await fetch(`/api/datasets/${item.id}`, { method: "DELETE" });
    if (!response.ok) {
      let message = "删除失败";
      try {
        const body = await response.json();
        message = body.error || message;
      } catch {
        message = response.statusText || message;
      }
      throw new Error(message);
    }

    const nextSelectedId =
      state.selectedId === item.id ? state.datasets.find((dataset) => dataset.id !== item.id)?.id || null : state.selectedId;
    await loadDatasets(nextSelectedId);
  } catch (error) {
    button.disabled = false;
    button.textContent = "删除";
    window.alert(error.message);
  }
}

async function loadDetail(id) {
  state.selectedId = id;
  renderList();
  clearPreview("正在读取数据详情...");

  const response = await fetch(`/api/datasets/${id}`);
  const dataset = await response.json();
  if (!response.ok) {
    detailPanel.classList.add("hidden");
    return;
  }

  detailPanel.classList.remove("hidden");
  detailTitle.textContent = dataset.name;
  detailMeta.textContent = `${dataset.fileName} · ${formatSize(dataset.size)} · ${formatNumber(dataset.rows || 0)} 行 · ${formatNumber(dataset.columns?.length || 0)} 列`;
  downloadLink.href = `/api/datasets/${id}/download`;

  renderProfile(dataset);
  renderPreview(dataset);
}

function renderProfile(dataset) {
  const metrics = [
    ["文件大小", formatSize(dataset.size), false],
    ["总行数", formatNumber(dataset.rows || 0), false],
    ["字段数", formatNumber(dataset.columns?.length || 0), false],
    ["数值字段", formatNumber(dataset.profile?.numeric?.length || 0), false],
  ];

  for (const item of dataset.profile?.numeric || []) {
    metrics.push([
      item.column,
      `min ${round(item.min)} / avg ${round(item.avg)} / max ${round(item.max)}`,
      true,
    ]);
  }

  profileGrid.innerHTML = metrics
    .map(
      ([label, value, numeric]) =>
        `<div class="metric ${numeric ? "numeric" : ""}"><strong title="${escapeHTML(label)}">${escapeHTML(label)}</strong><span>${escapeHTML(value)}</span></div>`,
    )
    .join("");
}

function renderPreview(dataset) {
  const columns = dataset.columns || [];
  const rows = dataset.preview || [];
  if (columns.length === 0 || rows.length === 0) {
    previewTable.innerHTML = "<tbody><tr><td>暂无可预览数据</td></tr></tbody>";
    return;
  }

  previewTable.innerHTML = `
    <thead><tr>${columns.map((column) => `<th>${escapeHTML(column)}</th>`).join("")}</tr></thead>
    <tbody>
      ${rows
        .map((row) => `<tr>${columns.map((column) => `<td>${escapeHTML(row[column] || "")}</td>`).join("")}</tr>`)
        .join("")}
    </tbody>
  `;
}

function clearPreview(message = "") {
  profileGrid.innerHTML = "";
  previewTable.innerHTML = message ? `<tbody><tr><td>${escapeHTML(message)}</td></tr></tbody>` : "";
}

function setUploadState(isUploading, message, tone) {
  uploadButton.disabled = isUploading;
  uploadButton.textContent = isUploading ? "上传中..." : "上传数据";
  uploadStatus.textContent = message;
  uploadStatus.className = `status-message ${tone}`;
}

function formatSize(bytes) {
  if (!bytes) return "0 B";
  const units = ["B", "KB", "MB", "GB"];
  const index = Math.min(Math.floor(Math.log(bytes) / Math.log(1024)), units.length - 1);
  return `${(bytes / 1024 ** index).toFixed(index === 0 ? 0 : 1)} ${units[index]}`;
}

function formatNumber(value) {
  return new Intl.NumberFormat("zh-CN").format(value);
}

function formatDate(value) {
  return new Intl.DateTimeFormat("zh-CN", {
    month: "2-digit",
    day: "2-digit",
    hour: "2-digit",
    minute: "2-digit",
  }).format(new Date(value));
}

function round(value) {
  return Number(value).toFixed(2);
}

function escapeHTML(value) {
  return String(value)
    .replaceAll("&", "&amp;")
    .replaceAll("<", "&lt;")
    .replaceAll(">", "&gt;")
    .replaceAll('"', "&quot;")
    .replaceAll("'", "&#039;");
}

loadDatasets();
