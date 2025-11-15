/* ============================================
   Services Component
   Управление сервисами
   ============================================ */

import { api } from '../api/wails.js';
import { $, $$, addClass, removeClass } from '../utils/dom.js';
import { notification } from '../ui/notification.js';
import { sleep, isWailsAvailable } from '../utils/helpers.js';

// Класс для управления сервисами
export class ServicesManager {
    constructor() {
        this.serverRunning = true;
        this.isOperating = false;
        this.controlBtn = $('serverControlBtn');
        this.statusIndicator = document.querySelector('.status-indicator');
        this.statusText = document.querySelector('.status-text');
        this.btnText = document.querySelector('.btn-text');
        this.init();
    }

    init() {
        if (this.controlBtn) {
            this.controlBtn.addEventListener('click', () => this.toggleServer());
        }

        // Подписка на события
        if (window.runtime?.EventsOn) {
            window.runtime.EventsOn('service:changed', (status) => {
                this.renderServices(status);
            });

            window.runtime.EventsOn('server:already_running', () => {
                notification.error('vServer уже запущен!<br><br>Закройте другой экземпляр перед запуском нового.', 5000);
                this.setServerStatus(false, 'Уже запущен в другом процессе');
            });
        }
    }

    // Переключить состояние сервера
    async toggleServer() {
        if (this.serverRunning) {
            await this.stopServer();
        } else {
            await this.startServer();
        }
    }

    // Запустить сервер
    async startServer() {
        this.isOperating = true;
        this.controlBtn.disabled = true;
        this.statusText.textContent = 'Запускается...';
        this.btnText.textContent = 'Ожидайте...';
        this.setAllServicesPending('Запуск');

        await api.startServer();

        // Ждём пока все запустятся
        let attempts = 0;
        while (attempts < 20) {
            await sleep(500);
            if (await api.checkServicesReady()) {
                break;
            }
            attempts++;
        }

        this.isOperating = false;
        this.setServerStatus(true, 'Сервер запущен');
        removeClass(this.controlBtn, 'start-mode');
        this.btnText.textContent = 'Остановить';
        this.controlBtn.disabled = false;
    }

    // Остановить сервер
    async stopServer() {
        this.isOperating = true;
        this.controlBtn.disabled = true;
        this.statusText.textContent = 'Выключается...';
        this.btnText.textContent = 'Ожидайте...';
        this.setAllServicesPending('Остановка');

        await api.stopServer();
        await sleep(1500);

        this.isOperating = false;
        this.setServerStatus(false, 'Сервер остановлен');
        addClass(this.controlBtn, 'start-mode');
        this.btnText.textContent = 'Запустить';
        this.controlBtn.disabled = false;
    }

    // Установить статус сервера
    setServerStatus(isOnline, text) {
        this.serverRunning = isOnline;
        
        if (isOnline) {
            removeClass(this.statusIndicator, 'status-offline');
            addClass(this.statusIndicator, 'status-online');
        } else {
            removeClass(this.statusIndicator, 'status-online');
            addClass(this.statusIndicator, 'status-offline');
        }
        
        this.statusText.textContent = text;
    }

    // Установить всем сервисам статус pending
    setAllServicesPending(text) {
        const badges = $$('.service-card .badge');
        badges.forEach(badge => {
            badge.className = 'badge badge-pending';
            badge.textContent = text;
        });
    }

    // Отрисовать статусы сервисов
    renderServices(data) {
        const services = [data.http, data.https, data.mysql, data.php, data.proxy];
        const cards = $$('.service-card');

        services.forEach((service, index) => {
            const card = cards[index];
            if (!card) return;

            const badge = card.querySelector('.badge');
            const infoValues = card.querySelectorAll('.info-value');

            // Обновляем badge только если НЕ в процессе операции
            if (badge && !this.isOperating) {
                if (service.status) {
                    badge.className = 'badge badge-online';
                    badge.textContent = 'Активен';
                } else {
                    badge.className = 'badge badge-offline';
                    badge.textContent = 'Остановлен';
                }
            }

            // Обновляем значения
            if (service.name === 'Proxy') {
                if (infoValues[0] && service.info) {
                    infoValues[0].textContent = service.info;
                }
            } else {
                if (infoValues[0]) {
                    infoValues[0].textContent = service.port;
                }
            }
        });
    }

    // Загрузить статусы сервисов
    async loadStatus() {
        if (isWailsAvailable()) {
            const data = await api.getAllServicesStatus();
            if (data) {
                this.renderServices(data);
            }
        } else {
            // Используем тестовые данные если Wails недоступен
            const mockData = {
                http: { name: 'HTTP', status: true, port: '80' },
                https: { name: 'HTTPS', status: true, port: '443' },
                mysql: { name: 'MySQL', status: true, port: '3306' },
                php: { name: 'PHP', status: true, port: '8000-8003' },
                proxy: { name: 'Proxy', status: true, port: '', info: '1 из 3' }
            };
            this.renderServices(mockData);
        }
    }
}

