/* ============================================
   Window Controls
   Управление окном приложения
   ============================================ */

import { $, addClass } from '../utils/dom.js';

// Класс для управления окном
export class WindowControls {
    constructor() {
        this.minimizeBtn = $('minimizeBtn');
        this.maximizeBtn = $('maximizeBtn');
        this.closeBtn = $('closeBtn');
        this.init();
    }

    init() {
        if (this.minimizeBtn) {
            this.minimizeBtn.addEventListener('click', () => this.minimize());
        }

        if (this.maximizeBtn) {
            this.maximizeBtn.addEventListener('click', () => this.maximize());
        }

        if (this.closeBtn) {
            this.closeBtn.addEventListener('click', () => this.close());
        }
    }

    minimize() {
        if (window.runtime?.WindowMinimise) {
            window.runtime.WindowMinimise();
        }
    }

    maximize() {
        if (window.runtime?.WindowToggleMaximise) {
            window.runtime.WindowToggleMaximise();
        }
    }

    close() {
        if (window.runtime?.Quit) {
            window.runtime.Quit();
        }
    }
}

