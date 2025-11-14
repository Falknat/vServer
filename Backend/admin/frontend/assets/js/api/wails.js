/* ============================================
   Wails API Wrapper
   Обёртка над Wails API
   ============================================ */

import { isWailsAvailable, log } from '../utils/helpers.js';

/**
 * Базовый класс для работы с Wails API
 */
class WailsAPI {
    constructor() {
        this.available = isWailsAvailable();
    }

    /**
     * Проверка доступности API
     */
    checkAvailability() {
        if (!this.available) {
            log('Wails API недоступен', 'warn');
            return false;
        }
        return true;
    }

    /**
     * Получить статус всех сервисов
     */
    async getAllServicesStatus() {
        if (!this.checkAvailability()) return null;
        try {
            return await window.go.admin.App.GetAllServicesStatus();
        } catch (error) {
            log(`Ошибка получения статуса сервисов: ${error.message}`, 'error');
            return null;
        }
    }

    /**
     * Получить список сайтов
     */
    async getSitesList() {
        if (!this.checkAvailability()) return [];
        try {
            return await window.go.admin.App.GetSitesList();
        } catch (error) {
            log(`Ошибка получения списка сайтов: ${error.message}`, 'error');
            return [];
        }
    }

    /**
     * Получить список прокси
     */
    async getProxyList() {
        if (!this.checkAvailability()) return [];
        try {
            return await window.go.admin.App.GetProxyList();
        } catch (error) {
            log(`Ошибка получения списка прокси: ${error.message}`, 'error');
            return [];
        }
    }

    /**
     * Получить правила vAccess
     */
    async getVAccessRules(host, isProxy) {
        if (!this.checkAvailability()) return { rules: [] };
        try {
            return await window.go.admin.App.GetVAccessRules(host, isProxy);
        } catch (error) {
            log(`Ошибка получения правил vAccess: ${error.message}`, 'error');
            return { rules: [] };
        }
    }

    /**
     * Сохранить правила vAccess
     */
    async saveVAccessRules(host, isProxy, configJSON) {
        if (!this.checkAvailability()) return 'Error: API недоступен';
        try {
            return await window.go.admin.App.SaveVAccessRules(host, isProxy, configJSON);
        } catch (error) {
            log(`Ошибка сохранения правил vAccess: ${error.message}`, 'error');
            return `Error: ${error.message}`;
        }
    }

    /**
     * Запустить сервер
     */
    async startServer() {
        if (!this.checkAvailability()) return;
        try {
            await window.go.admin.App.StartServer();
        } catch (error) {
            log(`Ошибка запуска сервера: ${error.message}`, 'error');
        }
    }

    /**
     * Остановить сервер
     */
    async stopServer() {
        if (!this.checkAvailability()) return;
        try {
            await window.go.admin.App.StopServer();
        } catch (error) {
            log(`Ошибка остановки сервера: ${error.message}`, 'error');
        }
    }

    /**
     * Проверить готовность сервисов
     */
    async checkServicesReady() {
        if (!this.checkAvailability()) return false;
        try {
            return await window.go.admin.App.CheckServicesReady();
        } catch (error) {
            return false;
        }
    }

    /**
     * Открыть папку сайта
     */
    async openSiteFolder(host) {
        if (!this.checkAvailability()) return;
        try {
            await window.go.admin.App.OpenSiteFolder(host);
        } catch (error) {
            log(`Ошибка открытия папки: ${error.message}`, 'error');
        }
    }
}

// Экспортируем единственный экземпляр
export const api = new WailsAPI();

