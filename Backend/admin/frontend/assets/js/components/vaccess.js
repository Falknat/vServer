/* ============================================
   vAccess Component
   Управление правилами доступа
   ============================================ */

import { api } from '../api/wails.js';
import { $, hide, show } from '../utils/dom.js';
import { notification } from '../ui/notification.js';
import { modal } from '../ui/modal.js';
import { isWailsAvailable } from '../utils/helpers.js';

// Класс для управления vAccess правилами
export class VAccessManager {
    constructor() {
        this.vAccessHost = '';
        this.vAccessIsProxy = false;
        this.vAccessRules = [];
        this.vAccessReturnSection = 'sectionSites';
        this.draggedIndex = null;
        this.editingField = null;
    }

    // Открыть редактор vAccess
    async open(host, isProxy) {
        this.vAccessHost = host;
        this.vAccessIsProxy = isProxy;

        // Запоминаем откуда пришли
        if ($('sectionSites').style.display !== 'none') {
            this.vAccessReturnSection = 'sectionSites';
        } else if ($('sectionProxy').style.display !== 'none') {
            this.vAccessReturnSection = 'sectionProxy';
        }

        // Загружаем правила
        if (isWailsAvailable()) {
            const config = await api.getVAccessRules(host, isProxy);
            this.vAccessRules = config.rules || [];
        } else {
            // Тестовые данные для браузерного режима
            this.vAccessRules = [
                {
                    type: 'Disable',
                    type_file: ['*.php'],
                    path_access: ['/uploads/*'],
                    ip_list: [],
                    exceptions_dir: [],
                    url_error: '404'
                }
            ];
        }

        // Обновляем UI
        const subtitle = isProxy 
            ? 'Управление правилами доступа для прокси-сервиса' 
            : 'Управление правилами доступа для сайта';

        $('breadcrumbHost').textContent = host;
        $('vAccessSubtitle').textContent = subtitle;

        // Переключаем на страницу редактора
        this.hideAllSections();
        show($('sectionVAccessEditor'));

        // Рендерим правила и показываем правильную вкладку
        this.renderRulesList();
        this.switchTab('rules');

        // Привязываем кнопку сохранения
        const saveBtn = $('saveVAccessBtn');
        if (saveBtn) {
            saveBtn.onclick = async () => await this.save();
        }
    }

    // Скрыть все секции
    hideAllSections() {
        hide($('sectionServices'));
        hide($('sectionSites'));
        hide($('sectionProxy'));
        hide($('sectionSettings'));
        hide($('sectionVAccessEditor'));
    }

    // Вернуться на главную
    backToMain() {
        this.hideAllSections();
        show($('sectionServices'));
        show($('sectionSites'));
        show($('sectionProxy'));
    }

    // Переключить вкладку
    switchTab(tab) {
        const tabs = document.querySelectorAll('.vaccess-tab[data-tab]');
        tabs.forEach(t => {
            if (t.dataset.tab === tab) {
                t.classList.add('active');
            } else {
                t.classList.remove('active');
            }
        });

        if (tab === 'rules') {
            show($('vAccessRulesTab'));
            hide($('vAccessHelpTab'));
        } else {
            hide($('vAccessRulesTab'));
            show($('vAccessHelpTab'));
        }
    }

    // Сохранить изменения
    async save() {
        if (isWailsAvailable()) {
            const config = { rules: this.vAccessRules };
            const configJSON = JSON.stringify(config);
            const result = await api.saveVAccessRules(this.vAccessHost, this.vAccessIsProxy, configJSON);

            if (result.startsWith('Error')) {
                notification.error(result, 2000);
            } else {
                notification.success('✅ Правила vAccess успешно сохранены', 1000);
            }
        } else {
            // Браузерный режим - просто показываем уведомление
            notification.success('Данные сохранены (тестовый режим)');
        }
    }

    // Отрисовать список правил
    renderRulesList() {
        const tbody = $('vAccessTableBody');
        const emptyState = $('vAccessEmpty');
        const table = document.querySelector('.vaccess-table');

        if (!tbody) return;

        // Показываем/скрываем пустое состояние
        if (this.vAccessRules.length === 0) {
            if (table) hide(table);
            if (emptyState) show(emptyState);
            return;
        } else {
            if (table) show(table);
            if (emptyState) hide(emptyState);
        }

        tbody.innerHTML = this.vAccessRules.map((rule, index) => `
            <tr draggable="true" data-index="${index}">
                <td class="drag-handle"><i class="fas fa-grip-vertical"></i></td>
                <td data-field="type" data-index="${index}">
                    <span class="badge ${rule.type === 'Allow' ? 'badge-yes' : 'badge-no'}">${rule.type}</span>
                </td>
                <td data-field="type_file" data-index="${index}">
                    ${(rule.type_file || []).length > 0 ? (rule.type_file || []).map(f => `<code class="mini-tag">${f}</code>`).join(' ') : '<span class="empty-field">-</span>'}
                </td>
                <td data-field="path_access" data-index="${index}">
                    ${(rule.path_access || []).length > 0 ? (rule.path_access || []).map(p => `<code class="mini-tag">${p}</code>`).join(' ') : '<span class="empty-field">-</span>'}
                </td>
                <td data-field="ip_list" data-index="${index}">
                    ${(rule.ip_list || []).length > 0 ? (rule.ip_list || []).map(ip => `<code class="mini-tag">${ip}</code>`).join(' ') : '<span class="empty-field">-</span>'}
                </td>
                <td data-field="exceptions_dir" data-index="${index}">
                    ${(rule.exceptions_dir || []).length > 0 ? (rule.exceptions_dir || []).map(e => `<code class="mini-tag">${e}</code>`).join(' ') : '<span class="empty-field">-</span>'}
                </td>
                <td data-field="url_error" data-index="${index}">
                    <code class="mini-tag">${rule.url_error || '404'}</code>
                </td>
                <td>
                    <button class="icon-btn-small" data-action="remove-rule" data-index="${index}" title="Удалить"><i class="fas fa-trash"></i></button>
                </td>
            </tr>
        `).join('');

        // Добавляем обработчики
        this.attachRulesEventListeners();
    }

    // Добавить обработчики событий для правил
    attachRulesEventListeners() {
        // Drag & Drop
        const rows = document.querySelectorAll('#vAccessTableBody tr[draggable]');
        rows.forEach(row => {
            row.addEventListener('dragstart', (e) => this.onDragStart(e));
            row.addEventListener('dragover', (e) => this.onDragOver(e));
            row.addEventListener('drop', (e) => this.onDrop(e));
        });

        // Клик по ячейкам для редактирования
        const cells = document.querySelectorAll('#vAccessTableBody td[data-field]');
        cells.forEach(cell => {
            cell.addEventListener('click', () => {
                const field = cell.getAttribute('data-field');
                const index = parseInt(cell.getAttribute('data-index'));
                this.editRuleField(index, field);
            });
        });

        // Кнопки удаления
        const removeButtons = document.querySelectorAll('[data-action="remove-rule"]');
        removeButtons.forEach(btn => {
            btn.addEventListener('click', () => {
                const index = parseInt(btn.getAttribute('data-index'));
                this.removeRule(index);
            });
        });
    }

    // Добавить новое правило
    addRule() {
        this.vAccessRules.push({
            type: 'Disable',
            type_file: [],
            path_access: [],
            ip_list: [],
            exceptions_dir: [],
            url_error: '404'
        });
        
        this.switchTab('rules');
        this.renderRulesList();
    }

    // Удалить правило
    removeRule(index) {
        this.vAccessRules.splice(index, 1);
        this.renderRulesList();
    }

    // Редактировать поле правила
    editRuleField(index, field) {
        const rule = this.vAccessRules[index];

        if (field === 'type') {
            // Переключаем тип
            rule.type = rule.type === 'Allow' ? 'Disable' : 'Allow';
            this.renderRulesList();
        } else if (field === 'url_error') {
            // Простой prompt для ошибки
            const value = prompt('Страница ошибки:', rule.url_error || '404');
            if (value !== null) {
                rule.url_error = value;
                this.renderRulesList();
            }
        } else {
            // Для массивов - показываем форму редактирования
            this.showFieldEditor(index, field);
        }
    }

    // Показать редактор поля
    showFieldEditor(index, field) {
        const rule = this.vAccessRules[index];
        const fieldNames = {
            'type_file': 'Расширения файлов',
            'path_access': 'Пути доступа',
            'ip_list': 'IP адреса',
            'exceptions_dir': 'Исключения'
        };

        const placeholders = {
            'type_file': '*.php',
            'path_access': '/admin/*',
            'ip_list': '127.0.0.1',
            'exceptions_dir': '/public/*'
        };

        const content = `
            <div class="field-editor">
                <div class="tag-input-wrapper" style="margin-bottom: 16px;">
                    <input type="text" class="form-input" id="fieldInput" placeholder="${placeholders[field]}">
                    <button class="action-btn" id="addFieldValueBtn"><i class="fas fa-plus"></i> Добавить</button>
                </div>
                <div class="tags-container" id="fieldTags">
                    ${(rule[field] || []).map(value => `
                        <span class="tag">
                            ${value}
                            <button class="tag-remove" data-value="${value}"><i class="fas fa-times"></i></button>
                        </span>
                    `).join('')}
                </div>
                <button class="action-btn" id="closeFieldEditorBtn" style="margin-top: 20px;">
                    <i class="fas fa-check"></i> Готово
                </button>
            </div>
        `;

        this.editingField = { index, field };
        modal.openFieldEditor(fieldNames[field], content);

        // Добавляем обработчики
        setTimeout(() => {
            $('addFieldValueBtn')?.addEventListener('click', () => this.addFieldValue());
            $('closeFieldEditorBtn')?.addEventListener('click', () => this.closeFieldEditor());
            
            const removeButtons = document.querySelectorAll('#fieldTags .tag-remove');
            removeButtons.forEach(btn => {
                btn.addEventListener('click', () => {
                    const value = btn.getAttribute('data-value');
                    this.removeFieldValue(value);
                });
            });
        }, 100);
    }

    // Добавить значение в поле
    addFieldValue() {
        const input = $('fieldInput');
        const value = input?.value.trim();

        if (value && this.editingField) {
            const { index, field } = this.editingField;
            if (!this.vAccessRules[index][field]) {
                this.vAccessRules[index][field] = [];
            }
            this.vAccessRules[index][field].push(value);
            input.value = '';
            this.showFieldEditor(index, field);
        }
    }

    // Удалить значение из поля
    removeFieldValue(value) {
        if (this.editingField) {
            const { index, field } = this.editingField;
            const arr = this.vAccessRules[index][field];
            const idx = arr.indexOf(value);
            if (idx > -1) {
                arr.splice(idx, 1);
                this.showFieldEditor(index, field);
            }
        }
    }

    // Закрыть редактор поля
    closeFieldEditor() {
        modal.closeFieldEditor();
        this.renderRulesList();
    }

    // Drag & Drop handlers
    onDragStart(event) {
        this.draggedIndex = parseInt(event.target.getAttribute('data-index'));
        event.target.style.opacity = '0.5';
    }

    onDragOver(event) {
        event.preventDefault();
    }

    onDrop(event) {
        event.preventDefault();
        const dropIndex = parseInt(event.target.closest('tr').getAttribute('data-index'));

        if (this.draggedIndex === null || this.draggedIndex === dropIndex) return;

        const draggedRule = this.vAccessRules[this.draggedIndex];
        this.vAccessRules.splice(this.draggedIndex, 1);
        this.vAccessRules.splice(dropIndex, 0, draggedRule);

        this.draggedIndex = null;
        this.renderRulesList();
    }
}

