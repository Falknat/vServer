/* ============================================
   Config API
   Работа с конфигурацией
   ============================================ */

import { isWailsAvailable } from '../utils/helpers.js';

// Класс для работы с конфигурацией
class ConfigAPI {
    constructor() {
        this.available = isWailsAvailable();
    }

    // Получить конфигурацию
    async getConfig() {
        if (!this.available) return null;
        try {
            return await window.go.admin.App.GetConfig();
        } catch (error) {
            return null;
        }
    }

    // Сохранить конфигурацию
    async saveConfig(configJSON) {
        if (!this.available) return 'Error: API недоступен';
        try {
            return await window.go.admin.App.SaveConfig(configJSON);
        } catch (error) {
            return `Error: ${error.message}`;
        }
    }

    // Включить Proxy Service
    async enableProxyService() {
        if (!this.available) return;
        try {
            await window.go.admin.App.EnableProxyService();
        } catch (error) {
        }
    }

    // Отключить Proxy Service
    async disableProxyService() {
        if (!this.available) return;
        try {
            await window.go.admin.App.DisableProxyService();
        } catch (error) {
        }
    }

    // Включить ACME Service
    async enableACMEService() {
        if (!this.available) return;
        try {
            await window.go.admin.App.EnableACMEService();
        } catch (error) {
        }
    }

    // Отключить ACME Service
    async disableACMEService() {
        if (!this.available) return;
        try {
            await window.go.admin.App.DisableACMEService();
        } catch (error) {
        }
    }

    // Перезапустить все сервисы
    async restartAllServices() {
        if (!this.available) return;
        try {
            await window.go.admin.App.RestartAllServices();
        } catch (error) {
        }
    }

    // Запустить HTTP Service
    async startHTTPService() {
        if (!this.available) return;
        try {
            await window.go.admin.App.StartHTTPService();
        } catch (error) {
        }
    }

    // Остановить HTTP Service
    async stopHTTPService() {
        if (!this.available) return;
        try {
            await window.go.admin.App.StopHTTPService();
        } catch (error) {
        }
    }

    // Запустить HTTPS Service
    async startHTTPSService() {
        if (!this.available) return;
        try {
            await window.go.admin.App.StartHTTPSService();
        } catch (error) {
        }
    }

    // Остановить HTTPS Service
    async stopHTTPSService() {
        if (!this.available) return;
        try {
            await window.go.admin.App.StopHTTPSService();
        } catch (error) {
        }
    }
}

// Экспортируем единственный экземпляр
export const configAPI = new ConfigAPI();

