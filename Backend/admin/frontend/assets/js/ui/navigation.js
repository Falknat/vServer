/* ============================================
   Navigation
   Управление навигацией
   ============================================ */

import { $, $$, hide, show, removeClass, addClass } from '../utils/dom.js';

/**
 * Класс для управления навигацией
 */
export class Navigation {
    constructor() {
        this.navItems = $$('.nav-item');
        this.sections = {
            services: $('sectionServices'),
            sites: $('sectionSites'),
            proxy: $('sectionProxy'),
            settings: $('sectionSettings'),
            vaccess: $('sectionVAccessEditor')
        };
        this.init();
    }

    init() {
        this.navItems.forEach((item, index) => {
            item.addEventListener('click', () => this.navigate(index));
        });
    }

    navigate(index) {
        // Убираем active со всех навигационных элементов
        this.navItems.forEach(nav => removeClass(nav, 'active'));
        addClass(this.navItems[index], 'active');

        // Скрываем все секции
        this.hideAllSections();

        // Показываем нужные секции
        if (index === 0) {
            // Главная - всё кроме настроек
            show(this.sections.services);
            show(this.sections.sites);
            show(this.sections.proxy);
        } else if (index === 3) {
            // Настройки
            show(this.sections.settings);
            // Загружаем конфигурацию при открытии
            if (window.loadConfig) {
                window.loadConfig();
            }
        }
    }

    hideAllSections() {
        Object.values(this.sections).forEach(section => {
            if (section) hide(section);
        });
    }

    showDashboard() {
        this.hideAllSections();
        show(this.sections.services);
        show(this.sections.sites);
        show(this.sections.proxy);
    }
}

