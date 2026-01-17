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
        this.certsCache = {};
        this.mockData = [
            {
                name: 'Home Voxsel',
                host: 'home.voxsel.ru',
                alias: ['home.voxsel.com'],
                status: 'active',
                root_file: 'index.html',
                root_file_routing: true,
                auto_create_ssl: false
            },
            {
                name: 'Finance',
                host: 'finance.voxsel.ru',
                alias: [],
                status: 'active',
                root_file: 'index.php',
                root_file_routing: false,
                auto_create_ssl: true
            },
            {
                name: 'Локальный сайт',
                host: '127.0.0.1',
                alias: ['localhost'],
                status: 'active',
                root_file: 'index.html',
                root_file_routing: true,
                auto_create_ssl: false
            }
        ];
        this.mockCerts = {
            'voxsel.ru': { has_cert: true, is_expired: false, days_left: 79, dns_names: ['*.voxsel.com', '*.voxsel.ru', 'voxsel.com', 'voxsel.ru'] },
            'finance.voxsel.ru': { has_cert: true, is_expired: false, days_left: 89, dns_names: ['finance.voxsel.ru'] }
        };
    }

    // Загрузить список сайтов
    async load() {
        if (isWailsAvailable()) {
            this.sitesData = await api.getSitesList();
            await this.loadCertsInfo();
        } else {
            this.sitesData = this.mockData;
            this.certsCache = this.mockCerts;
        }
        this.render();
    }

    // Загрузить информацию о сертификатах
    async loadCertsInfo() {
        const allCerts = await api.getAllCertsInfo();
        this.certsCache = {};
        for (const cert of allCerts) {
            this.certsCache[cert.domain] = cert;
        }
    }

    // Проверить соответствие домена wildcard паттерну
    matchesWildcard(domain, pattern) {
        if (pattern.startsWith('*.')) {
            const wildcardBase = pattern.slice(2);
            const domainParts = domain.split('.');
            if (domainParts.length >= 2) {
                const domainBase = domainParts.slice(1).join('.');
                return domainBase === wildcardBase;
            }
        }
        return domain === pattern;
    }

    // Найти сертификат для домена (включая wildcard)
    findCertForDomain(domain) {
        if (this.certsCache[domain]?.has_cert) {
            return this.certsCache[domain];
        }
        
        const domainParts = domain.split('.');
        if (domainParts.length >= 2) {
            const wildcardDomain = '*.' + domainParts.slice(1).join('.');
            if (this.certsCache[wildcardDomain]?.has_cert) {
                return this.certsCache[wildcardDomain];
            }
        }
        
        for (const [certDomain, cert] of Object.entries(this.certsCache)) {
            if (cert.has_cert && cert.dns_names) {
                for (const dnsName of cert.dns_names) {
                    if (this.matchesWildcard(domain, dnsName)) {
                        return cert;
                    }
                }
            }
        }
        
        return null;
    }

    // Получить иконку сертификата для домена
    getCertIcon(host, aliases = []) {
        const allDomains = [host, ...aliases.filter(a => !a.includes('*'))];
        
        for (const domain of allDomains) {
            const cert = this.findCertForDomain(domain);
            if (cert) {
                if (cert.is_expired) {
                    return `<span class="cert-icon cert-expired" title="SSL сертификат истёк"><i class="fas fa-shield-alt"></i></span>`;
                } else {
                    return `<span class="cert-icon cert-valid" title="SSL активен (${cert.days_left} дн.)"><i class="fas fa-shield-alt"></i></span>`;
                }
            }
        }
        return '';
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
            const certIcon = this.getCertIcon(site.host, site.alias);

            row.innerHTML = `
                <td>${certIcon}${site.name}</td>
                <td><code class="clickable-link" data-url="http://${site.host}">${site.host} <i class="fas fa-external-link-alt"></i></code></td>
                <td><code>${aliases}</code></td>
                <td><span class="badge ${statusBadge}">${site.status}</span></td>
                <td><code>${site.root_file}</code></td>
                <td>
                    <button class="icon-btn" data-action="open-folder" data-host="${site.host}" title="Открыть папку"><i class="fas fa-folder-open"></i></button>
                    <button class="icon-btn" data-action="edit-vaccess" data-host="${site.host}" data-is-proxy="false" title="vAccess"><i class="fas fa-user-lock"></i></button>
                    <button class="icon-btn" data-action="open-certs" data-host="${site.host}" data-aliases="${site.alias.join(',')}" data-is-proxy="false" title="SSL сертификаты"><i class="fas fa-shield-alt"></i></button>
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
            case 'open-certs':
                if (window.openCertManager) {
                    const aliasesStr = btn.getAttribute('data-aliases') || '';
                    const aliases = aliasesStr ? aliasesStr.split(',').filter(a => a) : [];
                    window.openCertManager(host, isProxy, aliases);
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

