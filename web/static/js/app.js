function showToast(text, type = 'success') {
    const el = document.getElementById('toast');
    el.textContent = text;
    el.className = `toast show ${type}`;
    setTimeout(() => el.classList.remove('show'), 3000);
}

function formatDate(dateStr) {
    const date = new Date(dateStr);
    const today = new Date();
    const isToday = date.toDateString() === today.toDateString();
    if (isToday) {
        return 'Сегодня, ' + date.toLocaleTimeString('ru-RU', { hour: '2-digit', minute: '2-digit' });
    }
    return date.toLocaleString('ru-RU', { day: 'numeric', month: 'short', hour: '2-digit', minute: '2-digit' });
}

function toInputFormat(isoStr) {
    if (!isoStr) return '';
    const date = new Date(isoStr);
    date.setMinutes(date.getMinutes() - date.getTimezoneOffset());
    return date.toISOString().slice(0, 16);
}

function getUTCDateString(localDateString) {
    if (!localDateString) return "";
    const d = new Date(localDateString);
    return d.toISOString().slice(0, 19);
}

function setDefaultDate(elementId) {
    const now = new Date();
    now.setMinutes(now.getMinutes() - now.getTimezoneOffset());
    const tomorrow = new Date(now);
    tomorrow.setDate(tomorrow.getDate() + 1);
    document.getElementById(elementId).value = tomorrow.toISOString().slice(0, 16);
}

function esc(unsafe) {
    if (!unsafe) return '';
    return unsafe.replace(/&/g, "&amp;").replace(/</g, "&lt;").replace(/>/g, "&gt;").replace(/"/g, "&quot;");
}

const createModal = document.getElementById('createModal');
const editModal = document.getElementById('editModal');

function openCreateModal() {
    createModal.classList.add('open');
    setDefaultDate('createDate');
    document.getElementById('createTitle').focus();
}
function closeCreateModal() { createModal.classList.remove('open'); }

function openEditModal(task, type) {
    document.getElementById('editTaskId').value = task.id;
    document.getElementById('editTaskType').value = type;
    document.getElementById('editTitle').value = task.title;
    document.getElementById('editDescription').value = task.description || '';

    const activeFields = document.getElementById('activeTaskFields');
    const modalTitle = document.getElementById('editModalTitle');

    if (type === 'active') {
        modalTitle.textContent = "Редактирование";
        activeFields.style.display = 'block';
        document.getElementById('editCreatedAt').value = toInputFormat(task.created_at);
        document.getElementById('editNextReview').value = toInputFormat(task.next_review_date);
    } else {
        modalTitle.textContent = "Редактирование (Архив)";
        activeFields.style.display = 'none';
    }
    editModal.classList.add('open');
}
function closeEditModal() { editModal.classList.remove('open'); }

window.onclick = function (event) {
    if (event.target == createModal) closeCreateModal();
    if (event.target == editModal) closeEditModal();
}

document.addEventListener('keydown', function (event) {
    if (event.key === 'Escape') {
        closeCreateModal();
        closeEditModal();
    }
});

async function loadTasks() {
    try {
        const res = await fetch('/tasks');
        const data = await res.json();
        if (res.ok) {
            renderActiveTasks(data.tasks || []);
            renderSucceededTasks(data.succeeded_tasks || []);
            document.getElementById('activeCount').textContent = (data.tasks || []).length;
        }
    } catch (e) { console.error(e); }
}

function renderActiveTasks(tasks) {
    const list = document.getElementById('tasksList');
    if (!tasks.length) {
        list.innerHTML = `<div class="empty-state"><span class="empty-icon">🎉</span><p>Всё чисто! Задач нет.</p></div>`;
        return;
    }
    list.innerHTML = tasks.map(t => `
                <div class="task-card" onclick='openEditModal(${JSON.stringify(t)}, "active")'>
                    <div>
                        <h3>${esc(t.title)}</h3>
                        <div class="task-desc">${esc(t.description || 'Нет описания')}</div>
                    </div>
                    <div class="task-meta">
                        <div class="date-badge"><span>Срок:</span><span>${formatDate(t.next_review_date)}</span></div>
                        <button class="btn-check-circle" onclick="event.stopPropagation(); completeTask(${t.id})" title="Завершить">✓</button>
                    </div>
                </div>
            `).join('');
}

function renderSucceededTasks(tasks) {
    const list = document.getElementById('succeededTasksList');
    if (!tasks.length) {
        list.innerHTML = `<div class="empty-state"><p>История пуста.</p></div>`;
        return;
    }
    list.innerHTML = tasks.map(t => `
                <div class="succeeded-card" onclick='openEditModal(${JSON.stringify(t)}, "succeeded")'>
                    <h4>${esc(t.title)}</h4>
                    ${t.description ? `<p style="font-size:13px; color:#86868B">${esc(t.description)}</p>` : ''}
                    <small>Завершено: ${formatDate(t.completed_at)}</small>
                </div>
            `).join('');
}

document.getElementById('createTaskForm').addEventListener('submit', async (e) => {
    e.preventDefault();
    const btn = document.getElementById('createSubmitBtn');
    const originalText = btn.textContent;
    btn.textContent = "Создание...";
    btn.disabled = true;

    const data = {
        title: document.getElementById('createTitle').value.trim(),
        description: document.getElementById('createDescription').value.trim() || null,
        next_review_date: getUTCDateString(document.getElementById('createDate').value)
    };

    try {
        const res = await fetch('/tasks', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(data)
        });
        if (res.ok) {
            showToast('Задача создана!');
            document.getElementById('createTaskForm').reset();
            closeCreateModal();
            loadTasks();
        } else {
            const err = await res.json();
            showToast(err.error, 'error');
        }
    } catch (e) { showToast('Ошибка сети', 'error'); }

    btn.textContent = originalText;
    btn.disabled = false;
});

async function saveTaskChanges() {
    const id = document.getElementById('editTaskId').value;
    const type = document.getElementById('editTaskType').value;

    const data = {
        title: document.getElementById('editTitle').value,
        description: document.getElementById('editDescription').value || null,
    };

    let url = `/tasks/${id}`;

    if (type === 'active') {
        data.created_at = getUTCDateString(document.getElementById('editCreatedAt').value);
        data.next_review_date = getUTCDateString(document.getElementById('editNextReview').value);
    } else {
        url = `/tasks/succeeded/${id}`;
    }

    try {
        const res = await fetch(url, {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(data)
        });
        if (res.ok) {
            closeEditModal();
            loadTasks();
            showToast('Сохранено');
        } else {
            showToast('Ошибка обновления', 'error');
        }
    } catch (e) { showToast('Ошибка сети', 'error'); }
}

async function deleteCurrentTask() {
    const id = document.getElementById('editTaskId').value;
    const type = document.getElementById('editTaskType').value;
    if (!confirm('Удалить эту задачу безвозвратно?')) return;

    let url = `/tasks/${id}`;
    if (type === 'succeeded') url = `/tasks/succeeded/${id}`;

    try {
        const res = await fetch(url, { method: 'DELETE' });
        if (res.ok) {
            closeEditModal();
            loadTasks();
            showToast('Удалено');
        } else {
            showToast('Ошибка при удалении', 'error');
        }
    } catch (e) { showToast('Ошибка сети', 'error'); }
}

async function completeTask(id) {
    try {
        const res = await fetch(`/tasks/${id}/complete`, { method: 'POST' });
        if (res.ok) {
            loadTasks();
            showToast('Задача выполнена!');
        }
    } catch (e) { console.error(e); }
}

document.addEventListener('DOMContentLoaded', loadTasks);