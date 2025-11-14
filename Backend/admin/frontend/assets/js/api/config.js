/* ============================================
   Config API
   Работа с конфигурацией
   ============================================ */

import { isWailsAvailable, log } from '../utils/helpers.js';

/**
 * Класс для работы с конфигурацией
 */
class ConfigAPI {
    constructor() {
        this.available = isWailsAvailable();
    }

    /**
     * Получить конфигурацию
     */
    async getConfig() {
        if (!this.available) return null;
        try {
            return await window.go.admin.App.GetConfig();
        } catch (error) {
            log(`Ошибка получения конфигурации: ${error.message}`, 'error');
            return null;
        }
    }

    /**
     * Сохранить конфигурацию
     */
    async saveConfig(configJSON) {
        if (!this.available) return 'Error: API недоступен';
        try {
            return await window.go.admin.App.SaveConfig(configJSON);
        } catch (error) {
            log(`Ошибка сохранения конфигурации: ${error.message}`, 'error');
            return `Error: ${error.message}`;
        }
    }

    /**
     * Включить Proxy Service
     */
    async enableProxyService() {
        if (!this.available) return;
        try {
            await window.go.admin.App.EnableProxyService();
        } catch (error) {
            log(`Ошибка включения Proxy: ${error.message}`, 'error');
        }
    }

    /**
     * Отключить Proxy Service
     */
    async disableProxyService() {
        if (!this.available) return;
        try {
            await window.go.admin.App.DisableProxyService();
        } catch (error) {
            log(`Ошибка отключения Proxy: ${error.message}`, 'error');
        }
    }

    /**
     * Перезапустить все сервисы
     */
    async restartAllServices() {
        if (!this.available) return;
        try {
            await window.go.admin.App.RestartAllServices();
        } catch (error) {
            log(`Ошибка перезапуска сервисов: ${error.message}`, 'error');
        }
    }

    /**
     * Запустить HTTP Service
     */
    async startHTTPService() {
        if (!this.available) return;
        try {
            await window.go.admin.App.StartHTTPService();
        } catch (error) {
            log(`Ошибка запуска HTTP: ${error.message}`, 'error');
        }
    }

    /**
     * Остановить HTTP Service
     */
    async stopHTTPService() {
        if (!this.available) return;
        try {
            await window.go.admin.App.StopHTTPService();
        } catch (error) {
            log(`Ошибка остановки HTTP: ${error.message}`, 'error');
        }
    }

    /**
     * Запустить HTTPS Service
     */
    async startHTTPSService() {
        if (!this.available) return;
        try {
            await window.go.admin.App.StartHTTPSService();
        } catch (error) {
            log(`Ошибка запуска HTTPS: ${error.message}`, 'error');
        }
    }

    /**
     * Остановить HTTPS Service
     */
    async stopHTTPSService() {
        if (!this.available) return;
        try {
            await window.go.admin.App.StopHTTPSService();
        } catch (error) {
            log(`Ошибка остановки HTTPS: ${error.message}`, 'error');
        }
    }
}

// Экспортируем единственный экземпляр
export const configAPI = new ConfigAPI();

