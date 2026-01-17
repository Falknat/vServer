/* ============================================
   Proxy Component
   Управление прокси
   ============================================ */

import { api } from '../api/wails.js';
import { isWailsAvailable } from '../utils/helpers.js';
import { $ } from '../utils/dom.js';

// Класс для управления прокси
export class ProxyManager {
    constructor() {
        this.proxiesData = [];
        this.certsCache = {};
        this.mockData = [
            {
                enable: true,
                external_domain: 'git.example.ru',
                local_address: '127.0.0.1',
                local_port: '3333',
                service_https_use: false,
                auto_https: true,
                auto_create_ssl: true,
                status: 'active'
            },
            {
                enable: true,
                external_domain: 'api.example.com',
                local_address: '127.0.0.1',
                local_port: '8080',
                service_https_use: true,
                auto_https: false,
                auto_create_ssl: false,
                status: 'active'
            },
            {
                enable: false,
                external_domain: 'test.example.net',
                local_address: '127.0.0.1',
                local_port: '5000',
                service_https_use: false,
                auto_https: false,
                auto_create_ssl: false,
                status: 'disabled'
            }
        ];
        this.mockCerts = {
            'git.example.ru': { has_cert: true, is_expired: false, days_left: 60 }
        };
    }

    // Загрузить список прокси
    async load() {
        if (isWailsAvailable()) {
            this.proxiesData = await api.getProxyList();
            await this.loadCertsInfo();
        } else {
            this.proxiesData = this.mockData;
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
    getCertIcon(domain) {
        const cert = this.findCertForDomain(domain);
        if (cert) {
            if (cert.is_expired) {
                return `<span class="cert-icon cert-expired" title="SSL сертификат истёк"><i class="fas fa-shield-alt"></i></span>`;
            } else {
                return `<span class="cert-icon cert-valid" title="SSL активен (${cert.days_left} дн.)"><i class="fas fa-shield-alt"></i></span>`;
            }
        }
        return '';
    }

    // Отрисовать список прокси
    render() {
        const tbody = $('proxyTable')?.querySelector('tbody');
        if (!tbody) return;

        tbody.innerHTML = '';

        this.proxiesData.forEach((proxy, index) => {
            const row = document.createElement('tr');
            const statusBadge = proxy.status === 'active' ? 'badge-online' : 'badge-offline';
            const httpsBadge = proxy.service_https_use ? 'badge-yes">HTTPS' : 'badge-no">HTTP';
            const autoHttpsBadge = proxy.auto_https ? 'badge-yes">Да' : 'badge-no">Нет';
            const protocol = proxy.auto_https ? 'https' : 'http';
            const certIcon = this.getCertIcon(proxy.external_domain);

            row.innerHTML = `
                <td>${certIcon}<code class="clickable-link" data-url="${protocol}://${proxy.external_domain}">${proxy.external_domain} <i class="fas fa-external-link-alt"></i></code></td>
                <td><code>${proxy.local_address}:${proxy.local_port}</code></td>
                <td><span class="badge ${httpsBadge}</span></td>
                <td><span class="badge ${autoHttpsBadge}</span></td>
                <td><span class="badge ${statusBadge}">${proxy.status}</span></td>
                <td>
                    <button class="icon-btn" data-action="edit-vaccess" data-host="${proxy.external_domain}" data-is-proxy="true" title="vAccess"><i class="fas fa-user-lock"></i></button>
                    <button class="icon-btn" data-action="open-certs" data-host="${proxy.external_domain}" data-is-proxy="true" title="SSL сертификаты"><i class="fas fa-shield-alt"></i></button>
                    <button class="icon-btn" data-action="edit-proxy" data-index="${index}" title="Редактировать"><i class="fas fa-edit"></i></button>
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
    handleAction(action, btn) {
        const host = btn.getAttribute('data-host');
        const index = parseInt(btn.getAttribute('data-index'));
        const isProxy = btn.getAttribute('data-is-proxy') === 'true';

        switch (action) {
            case 'edit-vaccess':
                if (window.editVAccess) {
                    window.editVAccess(host, isProxy);
                }
                break;
            case 'open-certs':
                if (window.openCertManager) {
                    window.openCertManager(host, isProxy, []);
                }
                break;
            case 'edit-proxy':
                if (window.editProxy) {
                    window.editProxy(index);
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

