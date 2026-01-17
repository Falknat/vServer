/* ============================================
   Navigation
   Управление навигацией
   ============================================ */

import { $, $$, hide, show, removeClass, addClass } from '../utils/dom.js';

// Класс для управления навигацией
export class Navigation {
    constructor() {
        this.navItems = $$('.nav-item[data-page]');
        this.sections = {
            services: $('sectionServices'),
            sites: $('sectionSites'),
            proxy: $('sectionProxy'),
            settings: $('sectionSettings'),
            vaccess: $('sectionVAccessEditor'),
            addSite: $('sectionAddSite')
        };
        this.init();
    }

    init() {
        this.navItems.forEach(item => {
            item.addEventListener('click', () => {
                const page = item.dataset.page;
                this.navigate(page, item);
            });
        });
    }

    navigate(page, clickedItem) {
        // Убираем active со всех навигационных элементов
        this.navItems.forEach(nav => removeClass(nav, 'active'));
        addClass(clickedItem, 'active');

        // Скрываем все секции
        this.hideAllSections();

        // Показываем нужные секции по имени страницы
        switch (page) {
            case 'dashboard':
                show(this.sections.services);
                show(this.sections.sites);
                show(this.sections.proxy);
                break;
            case 'settings':
                show(this.sections.settings);
                if (window.loadConfig) {
                    window.loadConfig();
                }
                break;
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

