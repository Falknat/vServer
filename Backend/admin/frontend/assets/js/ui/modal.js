/* ============================================
   Modal Manager
   Управление модальными окнами
   ============================================ */

import { $, addClass, removeClass } from '../utils/dom.js';

/**
 * Класс для управления модальными окнами
 */
export class Modal {
    constructor() {
        this.overlay = $('modalOverlay');
        this.title = $('modalTitle');
        this.content = $('modalContent');
        this.closeBtn = $('modalCloseBtn');
        this.cancelBtn = $('modalCancelBtn');
        this.saveBtn = $('modalSaveBtn');
        this.fieldEditorOverlay = $('fieldEditorOverlay');
        this.init();
    }

    init() {
        if (this.closeBtn) {
            this.closeBtn.addEventListener('click', () => this.close());
        }

        if (this.cancelBtn) {
            this.cancelBtn.addEventListener('click', () => this.close());
        }

        if (this.saveBtn) {
            this.saveBtn.addEventListener('click', () => {
                if (window.saveModalData) {
                    window.saveModalData();
                }
            });
        }

        if (this.overlay) {
            this.overlay.addEventListener('click', (e) => {
                if (e.target === this.overlay) {
                    this.close();
                }
            });
        }
    }

    /**
     * Открыть модальное окно
     * @param {string} title - Заголовок
     * @param {string} htmlContent - HTML контент
     */
    open(title, htmlContent) {
        if (this.title) this.title.textContent = title;
        if (this.content) this.content.innerHTML = htmlContent;
        if (this.overlay) addClass(this.overlay, 'show');
    }

    /**
     * Закрыть модальное окно
     */
    close() {
        if (this.overlay) removeClass(this.overlay, 'show');
    }

    /**
     * Установить обработчик сохранения
     * @param {Function} callback - Функция обратного вызова
     */
    onSave(callback) {
        if (this.saveBtn) {
            this.saveBtn.onclick = callback;
        }
    }

    /**
     * Открыть редактор поля
     * @param {string} title - Заголовок
     * @param {string} htmlContent - HTML контент
     */
    openFieldEditor(title, htmlContent) {
        const fieldTitle = $('fieldEditorTitle');
        const fieldContent = $('fieldEditorContent');
        
        if (fieldTitle) fieldTitle.textContent = title;
        if (fieldContent) fieldContent.innerHTML = htmlContent;
        if (this.fieldEditorOverlay) addClass(this.fieldEditorOverlay, 'show');
    }

    /**
     * Закрыть редактор поля
     */
    closeFieldEditor() {
        if (this.fieldEditorOverlay) removeClass(this.fieldEditorOverlay, 'show');
    }
}

// Глобальный экземпляр модального окна
export const modal = new Modal();

