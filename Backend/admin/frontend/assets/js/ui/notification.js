/* ============================================
   Notification System
   Система уведомлений
   ============================================ */

import { $, addClass, removeClass } from '../utils/dom.js';

/**
 * Класс для управления уведомлениями
 */
export class NotificationManager {
    constructor() {
        this.container = $('notification');
        this.loader = $('appLoader');
    }

    /**
     * Показать уведомление
     * @param {string} message - Текст сообщения
     * @param {string} type - Тип (success, error)
     * @param {number} duration - Длительность показа (мс)
     */
    show(message, type = 'success', duration = 1000) {
        if (!this.container) return;

        const icon = type === 'success' 
            ? '<i class="fas fa-check-circle"></i>' 
            : '<i class="fas fa-exclamation-circle"></i>';
        
        this.container.innerHTML = `
            <div class="notification-content">
                <div class="notification-icon">${icon}</div>
                <div class="notification-text">${message}</div>
            </div>
        `;
        
        this.container.className = `notification show ${type}`;
        
        setTimeout(() => {
            removeClass(this.container, 'show');
        }, duration);
    }

    /**
     * Показать успешное уведомление
     * @param {string} message - Текст сообщения
     * @param {number} duration - Длительность
     */
    success(message, duration = 1000) {
        this.show(message, 'success', duration);
    }

    /**
     * Показать уведомление об ошибке
     * @param {string} message - Текст сообщения
     * @param {number} duration - Длительность
     */
    error(message, duration = 2000) {
        this.show(message, 'error', duration);
    }

    /**
     * Скрыть загрузчик приложения
     */
    hideLoader() {
        if (!this.loader) return;
        
        setTimeout(() => {
            addClass(this.loader, 'hide');
            setTimeout(() => {
                if (this.loader.parentNode) {
                    this.loader.remove();
                }
            }, 500);
        }, 1500);
    }
}

// Глобальный экземпляр менеджера уведомлений
export const notification = new NotificationManager();

