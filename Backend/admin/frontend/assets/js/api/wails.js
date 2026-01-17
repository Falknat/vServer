/* ============================================
   Wails API Wrapper
   Обёртка над Wails API
   ============================================ */

import { isWailsAvailable } from '../utils/helpers.js';

// Базовый класс для работы с Wails API
class WailsAPI {
    constructor() {
        this.available = isWailsAvailable();
    }

    // Проверка доступности API
    checkAvailability() {
        if (!this.available) {
            return false;
        }
        return true;
    }

    // Получить статус всех сервисов
    async getAllServicesStatus() {
        if (!this.checkAvailability()) return null;
        try {
            return await window.go.admin.App.GetAllServicesStatus();
        } catch (error) {
            return null;
        }
    }

    // Получить список сайтов
    async getSitesList() {
        if (!this.checkAvailability()) return [];
        try {
            return await window.go.admin.App.GetSitesList();
        } catch (error) {
            return [];
        }
    }

    // Получить список прокси
    async getProxyList() {
        if (!this.checkAvailability()) return [];
        try {
            return await window.go.admin.App.GetProxyList();
        } catch (error) {
            return [];
        }
    }

    // Получить правила vAccess
    async getVAccessRules(host, isProxy) {
        if (!this.checkAvailability()) return { rules: [] };
        try {
            return await window.go.admin.App.GetVAccessRules(host, isProxy);
        } catch (error) {
            return { rules: [] };
        }
    }

    // Сохранить правила vAccess
    async saveVAccessRules(host, isProxy, configJSON) {
        if (!this.checkAvailability()) return 'Error: API недоступен';
        try {
            return await window.go.admin.App.SaveVAccessRules(host, isProxy, configJSON);
        } catch (error) {
            return `Error: ${error.message}`;
        }
    }

    // Запустить сервер
    async startServer() {
        if (!this.checkAvailability()) return;
        try {
            await window.go.admin.App.StartServer();
        } catch (error) {
        }
    }

    // Остановить сервер
    async stopServer() {
        if (!this.checkAvailability()) return;
        try {
            await window.go.admin.App.StopServer();
        } catch (error) {
        }
    }

    // Проверить готовность сервисов
    async checkServicesReady() {
        if (!this.checkAvailability()) return false;
        try {
            return await window.go.admin.App.CheckServicesReady();
        } catch (error) {
            return false;
        }
    }

    // Открыть папку сайта
    async openSiteFolder(host) {
        if (!this.checkAvailability()) return;
        try {
            await window.go.admin.App.OpenSiteFolder(host);
        } catch (error) {
        }
    }

    // Создать новый сайт
    async createNewSite(siteJSON) {
        if (!this.checkAvailability()) return 'Error: API недоступен';
        try {
            return await window.go.admin.App.CreateNewSite(siteJSON);
        } catch (error) {
            return `Error: ${error.message}`;
        }
    }

    // Загрузить сертификат для сайта
    async uploadCertificate(host, certType, certDataBase64) {
        if (!this.checkAvailability()) return 'Error: API недоступен';
        try {
            return await window.go.admin.App.UploadCertificate(host, certType, certDataBase64);
        } catch (error) {
            return `Error: ${error.message}`;
        }
    }

    // Перезагрузить SSL сертификаты
    async reloadSSLCertificates() {
        if (!this.checkAvailability()) return 'Error: API недоступен';
        try {
            return await window.go.admin.App.ReloadSSLCertificates();
        } catch (error) {
            return `Error: ${error.message}`;
        }
    }

    // Удалить сайт
    async deleteSite(host) {
        if (!this.checkAvailability()) return 'Error: API недоступен';
        try {
            return await window.go.admin.App.DeleteSite(host);
        } catch (error) {
            return `Error: ${error.message}`;
        }
    }

    // Получить информацию о сертификате для домена
    async getCertInfo(domain) {
        if (!this.checkAvailability()) return { has_cert: false };
        try {
            return await window.go.admin.App.GetCertInfo(domain);
        } catch (error) {
            return { has_cert: false };
        }
    }

    // Получить информацию о всех сертификатах
    async getAllCertsInfo() {
        if (!this.checkAvailability()) return [];
        try {
            return await window.go.admin.App.GetAllCertsInfo();
        } catch (error) {
            return [];
        }
    }

    // Удалить сертификат
    async deleteCertificate(domain) {
        if (!this.checkAvailability()) return 'Error: API недоступен';
        try {
            return await window.go.admin.App.DeleteCertificate(domain);
        } catch (error) {
            return `Error: ${error.message}`;
        }
    }

    // Получить SSL сертификат через Let's Encrypt
    async obtainSSLCertificate(domain) {
        if (!this.checkAvailability()) return 'Error: API недоступен';
        try {
            return await window.go.admin.App.ObtainSSLCertificate(domain);
        } catch (error) {
            return `Error: ${error.message}`;
        }
    }
}

// Экспортируем единственный экземпляр
export const api = new WailsAPI();

