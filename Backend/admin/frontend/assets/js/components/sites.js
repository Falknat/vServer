/* ============================================
   Sites Component
   Управление сайтами
   ============================================ */

import { api } from '../api/wails.js';
import { isWailsAvailable } from '../utils/helpers.js';
import { $ } from '../utils/dom.js';

// Класс для управления сайтами
export class SitesManager {
    constructor() {
        this.sitesData = [];
        this.mockData = [
            {
                name: 'Локальный сайт',
                host: '127.0.0.1',
                alias: ['localhost'],
                status: 'active',
                root_file: 'index.html',
                root_file_routing: true
            },
            {
                name: 'Тестовый проект',
                host: 'test.local',
                alias: ['*.test.local', 'test.com'],
                status: 'active',
                root_file: 'index.php',
                root_file_routing: false
            },
            {
                name: 'API сервис',
                host: 'api.example.com',
                alias: ['*.api.example.com'],
                status: 'inactive',
                root_file: 'index.php',
                root_file_routing: true
            }
        ];
    }

    // Загрузить список сайтов
    async load() {
        if (isWailsAvailable()) {
            this.sitesData = await api.getSitesList();
        } else {
            // Используем тестовые данные если Wails недоступен
            this.sitesData = this.mockData;
        }
        this.render();
    }

    // Отрисовать список сайтов
    render() {
        const tbody = $('sitesTable')?.querySelector('tbody');
        if (!tbody) return;

        tbody.innerHTML = '';

        this.sitesData.forEach((site, index) => {
            const row = document.createElement('tr');
            const statusBadge = site.status === 'active' ? 'badge-online' : 'badge-offline';
            const aliases = site.alias.join(', ');

            row.innerHTML = `
                <td>${site.name}</td>
                <td><code class="clickable-link" data-url="http://${site.host}">${site.host} <i class="fas fa-external-link-alt"></i></code></td>
                <td><code>${aliases}</code></td>
                <td><span class="badge ${statusBadge}">${site.status}</span></td>
                <td><code>${site.root_file}</code></td>
                <td>
                    <button class="icon-btn" data-action="open-folder" data-host="${site.host}" title="Открыть папку"><i class="fas fa-folder-open"></i></button>
                    <button class="icon-btn" data-action="edit-vaccess" data-host="${site.host}" data-is-proxy="false" title="vAccess"><i class="fas fa-shield-alt"></i></button>
                    <button class="icon-btn" data-action="edit-site" data-index="${index}" title="Редактировать"><i class="fas fa-edit"></i></button>
                </td>
            `;

            tbody.appendChild(row);
        });

        // Добавляем обработчики событий
        this.attachEventListeners();
    }

    // Добавить обработчики событий
    attachEventListeners() {
        // Кликабельные ссылки
        const links = document.querySelectorAll('.clickable-link[data-url]');
        links.forEach(link => {
            link.addEventListener('click', () => {
                const url = link.getAttribute('data-url');
                this.openLink(url);
            });
        });

        // Кнопки действий
        const buttons = document.querySelectorAll('[data-action]');
        buttons.forEach(btn => {
            btn.addEventListener('click', () => {
                const action = btn.getAttribute('data-action');
                this.handleAction(action, btn);
            });
        });
    }

    // Обработчик действий
    async handleAction(action, btn) {
        const host = btn.getAttribute('data-host');
        const index = parseInt(btn.getAttribute('data-index'));
        const isProxy = btn.getAttribute('data-is-proxy') === 'true';

        switch (action) {
            case 'open-folder':
                await api.openSiteFolder(host);
                break;
            case 'edit-vaccess':
                if (window.editVAccess) {
                    window.editVAccess(host, isProxy);
                }
                break;
            case 'edit-site':
                if (window.editSite) {
                    window.editSite(index);
                }
                break;
        }
    }

    // Открыть ссылку
    openLink(url) {
        if (window.runtime?.BrowserOpenURL) {
            window.runtime.BrowserOpenURL(url);
        } else {
            window.open(url, '_blank');
        }
    }
}

