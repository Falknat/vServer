/* ============================================
   DOM Utilities
   Утилиты для работы с DOM
   ============================================ */

/**
 * Получить элемент по ID
 * @param {string} id - ID элемента
 * @returns {HTMLElement|null}
 */
export function $(id) {
    return document.getElementById(id);
}

/**
 * Получить все элементы по селектору
 * @param {string} selector - CSS селектор
 * @param {HTMLElement} parent - Родительский элемент
 * @returns {NodeList}
 */
export function $$(selector, parent = document) {
    return parent.querySelectorAll(selector);
}

/**
 * Показать элемент
 * @param {HTMLElement|string} element - Элемент или ID
 */
export function show(element) {
    const el = typeof element === 'string' ? $(element) : element;
    if (el) el.style.display = 'block';
}

/**
 * Скрыть элемент
 * @param {HTMLElement|string} element - Элемент или ID
 */
export function hide(element) {
    const el = typeof element === 'string' ? $(element) : element;
    if (el) el.style.display = 'none';
}

/**
 * Переключить видимость элемента
 * @param {HTMLElement|string} element - Элемент или ID
 */
export function toggle(element) {
    const el = typeof element === 'string' ? $(element) : element;
    if (el) {
        el.style.display = el.style.display === 'none' ? 'block' : 'none';
    }
}

/**
 * Добавить класс
 * @param {HTMLElement|string} element - Элемент или ID
 * @param {string} className - Имя класса
 */
export function addClass(element, className) {
    const el = typeof element === 'string' ? $(element) : element;
    if (el) el.classList.add(className);
}

/**
 * Удалить класс
 * @param {HTMLElement|string} element - Элемент или ID
 * @param {string} className - Имя класса
 */
export function removeClass(element, className) {
    const el = typeof element === 'string' ? $(element) : element;
    if (el) el.classList.remove(className);
}

/**
 * Переключить класс
 * @param {HTMLElement|string} element - Элемент или ID
 * @param {string} className - Имя класса
 */
export function toggleClass(element, className) {
    const el = typeof element === 'string' ? $(element) : element;
    if (el) el.classList.toggle(className);
}

