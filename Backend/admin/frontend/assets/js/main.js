/* ============================================
   vServer Admin Panel - Main Entry Point
   Точка входа приложения
   ============================================ */

import { isWailsAvailable, sleep } from './utils/helpers.js';
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

// Главный класс приложения
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
    }

    // Загрузить шаблоны из templates.html
    async loadTemplates() {
        try {
            const response = await fetch('templates.html');
            const html = await response.text();
            document.getElementById('templates-container').innerHTML = html;
        } catch (error) {
            // Игнорируем ошибку
        }
    }

    // Получить шаблон по ID
    getTemplate(templateId) {
        const template = document.getElementById(templateId);
        return template ? template.content.cloneNode(true) : null;
    }

    // Запустить приложение
    async start() {
        // Загружаем шаблоны
        await this.loadTemplates();

        // Скрываем loader если не в Wails
        if (!this.isWails) {
            notification.hideLoader();
        }

        // Ждём немного перед загрузкой данных
        await sleep(1000);

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
    }

    // Загрузить начальные данные
    async loadInitialData() {
        await Promise.all([
            this.servicesManager.loadStatus(),
            this.sitesManager.load(),
            this.proxyManager.load()
        ]);
    }

    // Запустить автообновление
    startAutoRefresh() {
        setInterval(async () => {
            await this.loadInitialData();
        }, 5000);
    }

    // Привязать кнопки
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

    // Настроить глобальные обработчики
    setupGlobalHandlers() {
        Object.assign(window, {
            // Ссылки на менеджеры
            sitesManager: this.sitesManager,
            siteCreator: this.siteCreator,
            
            // SiteCreator
            backToMainFromAddSite: () => this.siteCreator.backToMain(),
            toggleCertUpload: () => this.siteCreator.toggleCertUpload(),
            handleCertFileSelect: (input, certType) => this.siteCreator.handleCertFile(input, certType),
            
            // vAccess
            editVAccess: (host, isProxy) => this.vAccessManager.open(host, isProxy),
            backToMain: () => this.vAccessManager.backToMain(),
            switchVAccessTab: (tab) => this.vAccessManager.switchTab(tab),
            saveVAccessChanges: async () => await this.vAccessManager.save(),
            addVAccessRule: () => this.vAccessManager.addRule(),
            dragStart: (e) => this.vAccessManager.onDragStart(e),
            dragOver: (e) => this.vAccessManager.onDragOver(e),
            drop: (e) => this.vAccessManager.onDrop(e),
            editRuleField: (i, f) => this.vAccessManager.editRuleField(i, f),
            removeVAccessRule: (i) => this.vAccessManager.removeRule(i),
            closeFieldEditor: () => this.vAccessManager.closeFieldEditor(),
            addFieldValue: () => this.vAccessManager.addFieldValue(),
            removeFieldValue: (v) => this.vAccessManager.removeFieldValue(v),
            
            // Settings
            loadConfig: async () => await this.loadConfigSettings(),
            saveSettings: async () => await this.saveConfigSettings(),
            
            // Модальные окна
            editSite: (i) => this.editSite(i),
            editProxy: (i) => this.editProxy(i),
            setStatus: (s) => this.setModalStatus(s),
            setProxyStatus: (s) => this.setModalStatus(s),
            addAliasTag: () => this.addAliasTag(),
            removeAliasTag: (btn) => btn.parentElement.remove(),
            saveModalData: async () => await this.saveModalData(),
            deleteSiteConfirm: async () => await this.deleteSiteConfirm(),
            
            // Тестовые функции
            editTestSite: (i) => {
                this.sitesManager.sitesData = [
                    {name: 'Локальный сайт', host: '127.0.0.1', alias: ['localhost'], status: 'active', root_file: 'index.html', root_file_routing: true},
                    {name: 'Тестовый проект', host: 'test.local', alias: ['*.test.local', 'test.com'], status: 'active', root_file: 'index.php', root_file_routing: false},
                    {name: 'API сервис', host: 'api.example.com', alias: ['*.api.example.com'], status: 'inactive', root_file: 'index.php', root_file_routing: true}
                ];
                this.editSite(i);
            },
            editTestProxy: (i) => {
                this.proxyManager.proxiesData = [
                    {enable: true, external_domain: 'git.example.ru', local_address: '127.0.0.1', local_port: '3333', service_https_use: false, auto_https: true},
                    {enable: true, external_domain: 'api.example.com', local_address: '127.0.0.1', local_port: '8080', service_https_use: true, auto_https: false},
                    {enable: false, external_domain: 'test.example.net', local_address: '127.0.0.1', local_port: '5000', service_https_use: false, auto_https: false}
                ];
                this.editProxy(i);
            },
            openTestLink: (url) => this.sitesManager.openLink(url),
            openSiteFolder: async (host) => await this.sitesManager.handleAction('open-folder', { getAttribute: () => host })
        });
    }

    // Загрузить настройки конфигурации
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

    // Сохранить настройки конфигурации
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

    // Редактировать сайт
    editSite(index) {
        const site = this.sitesManager.sitesData[index];
        if (!site) return;

        const template = this.getTemplate('edit-site-template');
        if (!template) return;

        const container = document.createElement('div');
        container.appendChild(template);

        // Открываем модальное окно с шаблоном
        modal.open('Редактировать сайт', container.innerHTML);
        window.currentEditType = 'site';
        window.currentEditIndex = index;

        // Заполняем данные ПОСЛЕ открытия модального окна
        setTimeout(() => {
            const statusBtn = document.querySelector(`[data-status="${site.status}"]`);
            if (statusBtn) statusBtn.classList.add('active');
            
            const editName = $('editName');
            const editHost = $('editHost');
            const editRootFile = $('editRootFile');
            const editRouting = $('editRouting');
            
            if (editName) editName.value = site.name;
            if (editHost) editHost.value = site.host;
            if (editRootFile) editRootFile.value = site.root_file;
            if (editRouting) editRouting.checked = site.root_file_routing;

            // Добавляем alias теги
            const aliasContainer = $('aliasTagsContainer');
            if (aliasContainer) {
                site.alias.forEach(alias => {
                    const tag = document.createElement('span');
                    tag.className = 'tag';
                    tag.innerHTML = `${alias}<button class="tag-remove" onclick="removeAliasTag(this)"><i class="fas fa-times"></i></button>`;
                    aliasContainer.appendChild(tag);
                });
            }

            // Привязываем обработчик кнопок статуса
            document.querySelectorAll('.status-btn').forEach(btn => {
                btn.onclick = () => this.setModalStatus(btn.dataset.value);
            });
        }, 50);

        this.addDeleteButtonToModal();
    }

    // Редактировать прокси
    editProxy(index) {
        const proxy = this.proxyManager.proxiesData[index];
        if (!proxy) return;

        const template = this.getTemplate('edit-proxy-template');
        if (!template) return;

        const container = document.createElement('div');
        container.appendChild(template);

        // Открываем модальное окно с шаблоном
        modal.open('Редактировать прокси', container.innerHTML);
        window.currentEditType = 'proxy';
        window.currentEditIndex = index;

        // Заполняем данные ПОСЛЕ открытия модального окна
        setTimeout(() => {
            const status = proxy.enable ? 'enable' : 'disable';
            const statusBtn = document.querySelector(`[data-status="${status}"]`);
            if (statusBtn) statusBtn.classList.add('active');

            const editDomain = $('editDomain');
            const editLocalAddr = $('editLocalAddr');
            const editLocalPort = $('editLocalPort');
            const editServiceHTTPS = $('editServiceHTTPS');
            const editAutoHTTPS = $('editAutoHTTPS');

            if (editDomain) editDomain.value = proxy.external_domain;
            if (editLocalAddr) editLocalAddr.value = proxy.local_address;
            if (editLocalPort) editLocalPort.value = proxy.local_port;
            if (editServiceHTTPS) editServiceHTTPS.checked = proxy.service_https_use;
            if (editAutoHTTPS) editAutoHTTPS.checked = proxy.auto_https;

            // Привязываем обработчик кнопок статуса
            document.querySelectorAll('.status-btn').forEach(btn => {
                btn.onclick = () => this.setModalStatus(btn.dataset.value);
            });
        }, 50);

        this.removeDeleteButtonFromModal();
    }

    // Установить статус в модальном окне
    setModalStatus(status) {
        const buttons = document.querySelectorAll('.status-btn');
        buttons.forEach(btn => {
            btn.classList.remove('active');
            if (btn.dataset.value === status) {
                btn.classList.add('active');
            }
        });
    }

    // Добавить alias tag
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

    // Сохранить данные модального окна
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

    // Перезапустить HTTP/HTTPS сервисы
    async restartHttpServices() {
        notification.show('Перезапуск HTTP/HTTPS...', 'success', 800);
        await configAPI.stopHTTPService();
        await configAPI.stopHTTPSService();
        await sleep(500);
        await configAPI.startHTTPService();
        await configAPI.startHTTPSService();
    }

    // Сохранить данные сайта
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

        const result = await configAPI.saveConfig(JSON.stringify(config, null, 4));
        if (result.startsWith('Error')) {
            notification.error(result);
            return;
        }

        await this.restartHttpServices();
        notification.success('Изменения сохранены и применены!', 1000);
        await this.sitesManager.load();
        modal.close();
    }

    // Сохранить данные прокси
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

        const result = await configAPI.saveConfig(JSON.stringify(config, null, 4));
        if (result.startsWith('Error')) {
            notification.error(result);
            return;
        }

        await this.restartHttpServices();
        notification.success('Изменения сохранены и применены!', 1000);
        await this.proxyManager.load();
        modal.close();
    }

    // Добавить кнопку удаления в модальное окно
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

    // Удалить кнопку удаления из модального окна
    removeDeleteButtonFromModal() {
        const deleteBtn = document.querySelector('#modalDeleteBtn');
        if (deleteBtn) deleteBtn.remove();
    }

    // Подтверждение удаления сайта
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
            await this.restartHttpServices();
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

