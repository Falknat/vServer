/* ============================================
   DOM Utilities
   Утилиты для работы с DOM
   ============================================ */

// Получить элемент по ID
export function $(id) {
    return document.getElementById(id);
}

// Получить все элементы по селектору
export function $$(selector, parent = document) {
    return parent.querySelectorAll(selector);
}

// Показать элемент
export function show(element) {
    const el = typeof element === 'string' ? $(element) : element;
    if (el) el.style.display = 'block';
}

// Скрыть элемент
export function hide(element) {
    const el = typeof element === 'string' ? $(element) : element;
    if (el) el.style.display = 'none';
}

// Переключить видимость элемента
export function toggle(element) {
    const el = typeof element === 'string' ? $(element) : element;
    if (el) {
        el.style.display = el.style.display === 'none' ? 'block' : 'none';
    }
}

// Добавить класс
export function addClass(element, className) {
    const el = typeof element === 'string' ? $(element) : element;
    if (el) el.classList.add(className);
}

// Удалить класс
export function removeClass(element, className) {
    const el = typeof element === 'string' ? $(element) : element;
    if (el) el.classList.remove(className);
}

// Переключить класс
export function toggleClass(element, className) {
    const el = typeof element === 'string' ? $(element) : element;
    if (el) el.classList.toggle(className);
}

