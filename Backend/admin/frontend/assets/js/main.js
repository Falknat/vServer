/* ============================================
   vServer Admin Panel - Main Entry Point
   Точка входа приложения
   ============================================ */

import { log, isWailsAvailable, sleep } from './utils/helpers.js';
import { WindowControls } from './ui/window.js';
import { Navigation } from './ui/navigation.js';
import { notification } from './ui/notification.js';
import { modal } from './ui/modal.js';
import { ServicesManager } from './components/services.js';
import { SitesManager } from './components/sites.js';
import { ProxyManager } from './components/proxy.js';
import { VAccessManager } from './components/vaccess.js';
import { SiteCreator } from './components/site-creator.js';
import { api } from './api/wails.js';
import { configAPI } from './api/config.js';
import { initCustomSelects } from './ui/custom-select.js';
import { $ } from './utils/dom.js';

/**
 * Главный класс приложения
 */
class App {
    constructor() {
        this.windowControls = new WindowControls();
        this.navigation = new Navigation();
        this.servicesManager = new ServicesManager();
        this.sitesManager = new SitesManager();
        this.proxyManager = new ProxyManager();
        this.vAccessManager = new VAccessManager();
        this.siteCreator = new SiteCreator();
        
        this.isWails = isWailsAvailable();
        
        log('Приложение инициализировано');
    }

    /**
     * Запустить приложение
     */
    async start() {
        log('Запуск приложения...');

        // Скрываем loader если не в Wails
        if (!this.isWails) {
            notification.hideLoader();
        }

        // Ждём немного перед загрузкой данных
        await sleep(1000);

        if (this.isWails) {
            log('Wails API доступен', 'info');
        } else {
            log('Wails API недоступен (браузерный режим)', 'warn');
        }

        // Загружаем начальные данные
        await this.loadInitialData();

        // Запускаем автообновление
        this.startAutoRefresh();

        // Скрываем loader после загрузки
        if (this.isWails) {
            notification.hideLoader();
        }

        // Настраиваем глобальные функции для совместимости
        this.setupGlobalHandlers();

        // Привязываем кнопки
        this.setupButtons();

        // Инициализируем кастомные select'ы
        initCustomSelects();

        log('Приложение запущено');
    }

    /**
     * Загрузить начальные данные
     */
    async loadInitialData() {
        await Promise.all([
            this.servicesManager.loadStatus(),
            this.sitesManager.load(),
            this.proxyManager.load()
        ]);
    }

    /**
     * Запустить автообновление
     */
    startAutoRefresh() {
        setInterval(async () => {
            await this.loadInitialData();
        }, 5000);
    }

    /**
     * Привязать кнопки
     */
    setupButtons() {
        // Кнопка добавления сайта
        const addSiteBtn = $('addSiteBtn');
        if (addSiteBtn) {
            addSiteBtn.addEventListener('click', () => {
                this.siteCreator.open();
            });
        }

        // Кнопка сохранения настроек
        const saveSettingsBtn = $('saveSettingsBtn');
        if (saveSettingsBtn) {
            saveSettingsBtn.addEventListener('click', async () => {
                await this.saveConfigSettings();
            });
        }

        // Кнопка сохранения vAccess (добавляем обработчик динамически при открытии vAccess)
        // Обработчик будет добавлен в VAccessManager.open()

        // Моментальное переключение Proxy без перезапуска
        const proxyCheckbox = $('proxyEnabled');
        if (proxyCheckbox) {
            proxyCheckbox.addEventListener('change', async (e) => {
                const isEnabled = e.target.checked;
                
                if (isEnabled) {
                    await configAPI.enableProxyService();
                    notification.success('Proxy Manager включен', 1000);
                } else {
                    await configAPI.disableProxyService();
                    notification.success('Proxy Manager отключен', 1000);
                }
            });
        }
    }

    /**
     * Настроить глобальные обработчики
     */
    setupGlobalHandlers() {
        // Глобальная ссылка на sitesManager
        window.sitesManager = this.sitesManager;
        window.siteCreator = this.siteCreator;

        // Для SiteCreator
        window.backToMainFromAddSite = () => {
            this.siteCreator.backToMain();
        };

        window.toggleCertUpload = () => {
            this.siteCreator.toggleCertUpload();
        };

        window.handleCertFileSelect = (input, certType) => {
            this.siteCreator.handleCertFile(input, certType);
        };

        // Для vAccess
        window.editVAccess = (host, isProxy) => {
            this.vAccessManager.open(host, isProxy);
        };

        window.backToMain = () => {
            this.vAccessManager.backToMain();
        };

        window.switchVAccessTab = (tab) => {
            this.vAccessManager.switchTab(tab);
        };

        window.saveVAccessChanges = async () => {
            await this.vAccessManager.save();
        };

        window.addVAccessRule = () => {
            this.vAccessManager.addRule();
        };

        // Для Settings
        window.loadConfig = async () => {
            await this.loadConfigSettings();
        };

        window.saveSettings = async () => {
            await this.saveConfigSettings();
        };

        // Для модальных окон
        window.editSite = (index) => {
            this.editSite(index);
        };

        window.editProxy = (index) => {
            this.editProxy(index);
        };

        window.setStatus = (status) => {
            this.setModalStatus(status);
        };

        window.setProxyStatus = (status) => {
            this.setModalStatus(status);
        };

        window.addAliasTag = () => {
            this.addAliasTag();
        };

        window.removeAliasTag = (btn) => {
            btn.parentElement.remove();
        };

        window.saveModalData = async () => {
            await this.saveModalData();
        };

        // Drag & Drop для vAccess
        window.dragStart = (event, index) => {
            this.vAccessManager.onDragStart(event);
        };

        window.dragOver = (event) => {
            this.vAccessManager.onDragOver(event);
        };

        window.drop = (event, index) => {
            this.vAccessManager.onDrop(event);
        };

        window.editRuleField = (index, field) => {
            this.vAccessManager.editRuleField(index, field);
        };

        window.removeVAccessRule = (index) => {
            this.vAccessManager.removeRule(index);
        };

        window.closeFieldEditor = () => {
            this.vAccessManager.closeFieldEditor();
        };

        window.addFieldValue = () => {
            this.vAccessManager.addFieldValue();
        };

        window.removeFieldValue = (value) => {
            this.vAccessManager.removeFieldValue(value);
        };

        // Тестовые функции (для браузерного режима)
        window.editTestSite = (index) => {
            const testSites = [
                {name: 'Локальный сайт', host: '127.0.0.1', alias: ['localhost'], status: 'active', root_file: 'index.html', root_file_routing: true},
                {name: 'Тестовый проект', host: 'test.local', alias: ['*.test.local', 'test.com'], status: 'active', root_file: 'index.php', root_file_routing: false},
                {name: 'API сервис', host: 'api.example.com', alias: ['*.api.example.com'], status: 'inactive', root_file: 'index.php', root_file_routing: true}
            ];
            this.sitesManager.sitesData = testSites;
            this.editSite(index);
        };

        window.editTestProxy = (index) => {
            const testProxies = [
                {enable: true, external_domain: 'git.example.ru', local_address: '127.0.0.1', local_port: '3333', service_https_use: false, auto_https: true},
                {enable: true, external_domain: 'api.example.com', local_address: '127.0.0.1', local_port: '8080', service_https_use: true, auto_https: false},
                {enable: false, external_domain: 'test.example.net', local_address: '127.0.0.1', local_port: '5000', service_https_use: false, auto_https: false}
            ];
            this.proxyManager.proxiesData = testProxies;
            this.editProxy(index);
        };

        window.openTestLink = (url) => {
            this.sitesManager.openLink(url);
        };

        window.openSiteFolder = async (host) => {
            await this.sitesManager.handleAction('open-folder', { getAttribute: () => host });
        };

        window.deleteSiteConfirm = async () => {
            await this.deleteSiteConfirm();
        };
    }

    /**
     * Загрузить настройки конфигурации
     */
    async loadConfigSettings() {
        if (!isWailsAvailable()) {
            // Тестовые данные для браузерного режима
            $('mysqlHost').value = '127.0.0.1';
            $('mysqlPort').value = 3306;
            $('phpHost').value = 'localhost';
            $('phpPort').value = 8000;
            $('proxyEnabled').checked = true;
            return;
        }

        const config = await configAPI.getConfig();
        if (!config) return;

        $('mysqlHost').value = config.Soft_Settings?.mysql_host || '127.0.0.1';
        $('mysqlPort').value = config.Soft_Settings?.mysql_port || 3306;
        $('phpHost').value = config.Soft_Settings?.php_host || 'localhost';
        $('phpPort').value = config.Soft_Settings?.php_port || 8000;
        $('proxyEnabled').checked = config.Soft_Settings?.proxy_enabled !== false;
    }

    /**
     * Сохранить настройки конфигурации
     */
    async saveConfigSettings() {
        const saveBtn = $('saveSettingsBtn');
        const originalText = saveBtn.querySelector('span').textContent;

        if (!isWailsAvailable()) {
            notification.success('Настройки сохранены (тестовый режим)', 1000);
            return;
        }

        try {
            saveBtn.disabled = true;
            saveBtn.querySelector('span').textContent = 'Сохранение...';

            const config = await configAPI.getConfig();
            config.Soft_Settings.mysql_host = $('mysqlHost').value;
            config.Soft_Settings.mysql_port = parseInt($('mysqlPort').value);
            config.Soft_Settings.php_host = $('phpHost').value;
            config.Soft_Settings.php_port = parseInt($('phpPort').value);
            config.Soft_Settings.proxy_enabled = $('proxyEnabled').checked;

            const configJSON = JSON.stringify(config, null, 4);
            const result = await configAPI.saveConfig(configJSON);

            if (result.startsWith('Error')) {
                notification.error(result);
                return;
            }

            saveBtn.querySelector('span').textContent = 'Перезапуск сервисов...';
            await configAPI.restartAllServices();

            notification.success('Настройки сохранены и сервисы перезапущены!', 1500);
        } catch (error) {
            notification.error('Ошибка: ' + error.message);
        } finally {
            saveBtn.disabled = false;
            saveBtn.querySelector('span').textContent = originalText;
        }
    }

    /**
     * Редактировать сайт
     */
    editSite(index) {
        const site = this.sitesManager.sitesData[index];
        if (!site) return;

        const content = `
            <div class="settings-form">
                <div class="form-group">
                    <label class="form-label">Статус сайта:</label>
                    <div class="status-toggle">
                        <button class="status-btn ${site.status === 'active' ? 'active' : ''}" onclick="setStatus('active')" data-value="active">
                            <i class="fas fa-check-circle"></i> Active
                        </button>
                        <button class="status-btn ${site.status === 'inactive' ? 'active' : ''}" onclick="setStatus('inactive')" data-value="inactive">
                            <i class="fas fa-times-circle"></i> Inactive
                        </button>
                    </div>
                </div>
                <div class="form-group">
                    <label class="form-label">Название сайта:</label>
                    <input type="text" class="form-input" id="editName" value="${site.name}">
                </div>
                <div class="form-group">
                    <label class="form-label">Host:</label>
                    <input type="text" class="form-input" id="editHost" value="${site.host}">
                </div>
                <div class="form-group">
                    <label class="form-label">Alias:</label>
                    <div class="tag-input-wrapper">
                        <input type="text" class="form-input" id="editAliasInput" placeholder="Введите alias и нажмите Добавить...">
                        <button class="action-btn" onclick="addAliasTag()"><i class="fas fa-plus"></i> Добавить</button>
                    </div>
                    <div class="tags-container" id="aliasTagsContainer">
                        ${site.alias.map(alias => `
                            <span class="tag">
                                ${alias}
                                <button class="tag-remove" onclick="removeAliasTag(this)"><i class="fas fa-times"></i></button>
                            </span>
                        `).join('')}
                    </div>
                </div>
                <div class="form-row">
                    <div class="form-group">
                        <label class="form-label">Root файл:</label>
                        <input type="text" class="form-input" id="editRootFile" value="${site.root_file}">
                    </div>
                    <div class="form-group">
                        <label class="form-label">Роутинг:</label>
                        <div class="toggle-wrapper">
                            <label class="toggle-switch">
                                <input type="checkbox" id="editRouting" ${site.root_file_routing ? 'checked' : ''}>
                                <span class="toggle-slider"></span>
                            </label>
                            <span class="toggle-label">Включён</span>
                        </div>
                    </div>
                </div>
            </div>
        `;

        modal.open('Редактировать сайт', content);
        window.currentEditType = 'site';
        window.currentEditIndex = index;
        
        // Добавляем кнопку удаления в футер модального окна
        this.addDeleteButtonToModal();
    }

    /**
     * Редактировать прокси
     */
    editProxy(index) {
        const proxy = this.proxyManager.proxiesData[index];
        if (!proxy) return;

        const content = `
            <div class="settings-form">
                <div class="form-group">
                    <label class="form-label">Статус прокси:</label>
                    <div class="status-toggle">
                        <button class="status-btn ${proxy.enable ? 'active' : ''}" onclick="setProxyStatus('enable')" data-value="enable">
                            <i class="fas fa-check-circle"></i> Включён
                        </button>
                        <button class="status-btn ${!proxy.enable ? 'active' : ''}" onclick="setProxyStatus('disable')" data-value="disable">
                            <i class="fas fa-times-circle"></i> Отключён
                        </button>
                    </div>
                </div>
                <div class="form-group">
                    <label class="form-label">Внешний домен:</label>
                    <input type="text" class="form-input" id="editDomain" value="${proxy.external_domain}">
                </div>
                <div class="form-row">
                    <div class="form-group">
                        <label class="form-label">Локальный адрес:</label>
                        <input type="text" class="form-input" id="editLocalAddr" value="${proxy.local_address}">
                    </div>
                    <div class="form-group">
                        <label class="form-label">Локальный порт:</label>
                        <input type="text" class="form-input" id="editLocalPort" value="${proxy.local_port}">
                    </div>
                </div>
                <div class="form-row">
                    <div class="form-group">
                        <label class="form-label">HTTPS к сервису:</label>
                        <div class="toggle-wrapper">
                            <label class="toggle-switch">
                                <input type="checkbox" id="editServiceHTTPS" ${proxy.service_https_use ? 'checked' : ''}>
                                <span class="toggle-slider"></span>
                            </label>
                            <span class="toggle-label">Включён</span>
                        </div>
                    </div>
                    <div class="form-group">
                        <label class="form-label">Авто HTTPS:</label>
                        <div class="toggle-wrapper">
                            <label class="toggle-switch">
                                <input type="checkbox" id="editAutoHTTPS" ${proxy.auto_https ? 'checked' : ''}>
                                <span class="toggle-slider"></span>
                            </label>
                            <span class="toggle-label">Включён</span>
                        </div>
                    </div>
                </div>
            </div>
        `;

        modal.open('Редактировать прокси', content);
        window.currentEditType = 'proxy';
        window.currentEditIndex = index;
        
        // Убираем кнопку удаления (для прокси не нужна)
        this.removeDeleteButtonFromModal();
    }

    /**
     * Установить статус в модальном окне
     */
    setModalStatus(status) {
        const buttons = document.querySelectorAll('.status-btn');
        buttons.forEach(btn => {
            btn.classList.remove('active');
            if (btn.dataset.value === status) {
                btn.classList.add('active');
            }
        });
    }

    /**
     * Добавить alias tag
     */
    addAliasTag() {
        const input = $('editAliasInput');
        const value = input?.value.trim();

        if (value) {
            const container = $('aliasTagsContainer');
            const tag = document.createElement('span');
            tag.className = 'tag';
            tag.innerHTML = `
                ${value}
                <button class="tag-remove" onclick="removeAliasTag(this)"><i class="fas fa-times"></i></button>
            `;
            container.appendChild(tag);
            input.value = '';
        }
    }

    /**
     * Сохранить данные модального окна
     */
    async saveModalData() {
        if (!isWailsAvailable()) {
            notification.success('Данные сохранены (тестовый режим)', 1000);
            modal.close();
            return;
        }

        if (window.currentEditType === 'site') {
            await this.saveSiteData();
        } else if (window.currentEditType === 'proxy') {
            await this.saveProxyData();
        }
    }

    /**
     * Сохранить данные сайта
     */
    async saveSiteData() {
        const index = window.currentEditIndex;
        const tags = document.querySelectorAll('#aliasTagsContainer .tag');
        const aliases = Array.from(tags).map(tag => tag.textContent.trim());
        const statusBtn = document.querySelector('.status-btn.active');

        const config = await configAPI.getConfig();
        config.Site_www[index] = {
            name: $('editName').value,
            host: $('editHost').value,
            alias: aliases,
            status: statusBtn ? statusBtn.dataset.value : 'active',
            root_file: $('editRootFile').value,
            root_file_routing: $('editRouting').checked
        };

        const configJSON = JSON.stringify(config, null, 4);
        const result = await configAPI.saveConfig(configJSON);

        if (result.startsWith('Error')) {
            notification.error(result);
        } else {
            notification.show('Перезапуск HTTP/HTTPS...', 'success', 800);

            await configAPI.stopHTTPService();
            await configAPI.stopHTTPSService();
            await sleep(500);
            await configAPI.startHTTPService();
            await configAPI.startHTTPSService();

            notification.success('Изменения сохранены и применены!', 1000);
            await this.sitesManager.load();
            modal.close();
        }
    }

    /**
     * Сохранить данные прокси
     */
    async saveProxyData() {
        const index = window.currentEditIndex;
        const statusBtn = document.querySelector('.status-btn.active');
        const isEnabled = statusBtn && statusBtn.dataset.value === 'enable';

        const config = await configAPI.getConfig();
        config.Proxy_Service[index] = {
            Enable: isEnabled,
            ExternalDomain: $('editDomain').value,
            LocalAddress: $('editLocalAddr').value,
            LocalPort: $('editLocalPort').value,
            ServiceHTTPSuse: $('editServiceHTTPS').checked,
            AutoHTTPS: $('editAutoHTTPS').checked
        };

        const configJSON = JSON.stringify(config, null, 4);
        const result = await configAPI.saveConfig(configJSON);

        if (result.startsWith('Error')) {
            notification.error(result);
        } else {
            notification.show('Перезапуск HTTP/HTTPS...', 'success', 800);

            await configAPI.stopHTTPService();
            await configAPI.stopHTTPSService();
            await sleep(500);
            await configAPI.startHTTPService();
            await configAPI.startHTTPSService();

            notification.success('Изменения сохранены и применены!', 1000);
            await this.proxyManager.load();
            modal.close();
        }
    }

    /**
     * Добавить кнопку удаления в модальное окно
     */
    addDeleteButtonToModal() {
        const footer = document.querySelector('.modal-footer');
        if (!footer) return;

        // Удаляем старую кнопку удаления если есть
        const oldDeleteBtn = footer.querySelector('#modalDeleteBtn');
        if (oldDeleteBtn) oldDeleteBtn.remove();

        // Создаём кнопку удаления
        const deleteBtn = document.createElement('button');
        deleteBtn.className = 'action-btn delete-btn';
        deleteBtn.id = 'modalDeleteBtn';
        deleteBtn.innerHTML = `
            <i class="fas fa-trash"></i>
            <span>Удалить сайт</span>
        `;
        deleteBtn.onclick = () => this.deleteSiteConfirm();

        // Вставляем перед кнопкой "Отмена"
        const cancelBtn = footer.querySelector('#modalCancelBtn');
        if (cancelBtn) {
            footer.insertBefore(deleteBtn, cancelBtn);
        }
    }

    /**
     * Удалить кнопку удаления из модального окна
     */
    removeDeleteButtonFromModal() {
        const deleteBtn = document.querySelector('#modalDeleteBtn');
        if (deleteBtn) deleteBtn.remove();
    }

    /**
     * Подтверждение удаления сайта
     */
    async deleteSiteConfirm() {
        const index = window.currentEditIndex;
        const site = this.sitesManager.sitesData[index];
        if (!site) return;

        // Подтверждение
        const confirmed = confirm(
            `⚠️ ВНИМАНИЕ!\n\n` +
            `Вы действительно хотите удалить сайт "${site.name}" (${site.host})?\n\n` +
            `Будут удалены:\n` +
            `• Папка сайта: WebServer/www/${site.host}/\n` +
            `• SSL сертификаты (если есть)\n` +
            `• Запись в конфигурации\n\n` +
            `Это действие НЕОБРАТИМО!`
        );

        if (!confirmed) return;

        try {
            notification.show('Удаление сайта...', 'info', 1000);

            const result = await api.deleteSite(site.host);

            if (result.startsWith('Error')) {
                notification.error(result, 3000);
                return;
            }

            notification.success('✅ Сайт успешно удалён!', 1500);

            // Перезапускаем HTTP/HTTPS
            notification.show('Перезапуск серверов...', 'success', 800);
            await configAPI.stopHTTPService();
            await configAPI.stopHTTPSService();
            await sleep(500);
            await configAPI.startHTTPService();
            await configAPI.startHTTPSService();

            notification.success('🚀 Серверы перезапущены!', 1000);

            // Закрываем модальное окно и обновляем список
            modal.close();
            await this.sitesManager.load();

        } catch (error) {
            notification.error('Ошибка: ' + error.message, 3000);
        }
    }
}

// Инициализация приложения при загрузке DOM
document.addEventListener('DOMContentLoaded', () => {
    const app = new App();
    app.start();
});

log('vServer Admin Panel загружен');

