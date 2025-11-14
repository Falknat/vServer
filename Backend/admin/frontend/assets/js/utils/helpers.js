/* ============================================
   Helper Utilities
   Вспомогательные функции
   ============================================ */

/**
 * Ждёт указанное время
 * @param {number} ms - Миллисекунды
 * @returns {Promise}
 */
export function sleep(ms) {
    return new Promise(resolve => setTimeout(resolve, ms));
}

/**
 * Debounce функция
 * @param {Function} func - Функция для debounce
 * @param {number} wait - Время задержки
 * @returns {Function}
 */
export function debounce(func, wait) {
    let timeout;
    return function executedFunction(...args) {
        const later = () => {
            clearTimeout(timeout);
            func(...args);
        };
        clearTimeout(timeout);
        timeout = setTimeout(later, wait);
    };
}

/**
 * Проверяет доступность Wails API
 * @returns {boolean}
 */
export function isWailsAvailable() {
    return typeof window.go !== 'undefined' && 
           window.go?.admin?.App !== undefined;
}

/**
 * Логирование с префиксом
 * @param {string} message - Сообщение
 * @param {string} type - Тип (log, error, warn, info)
 */
export function log(message, type = 'log') {
    const prefix = '🚀 vServer:';
    const styles = {
        log: '✅',
        error: '❌',
        warn: '⚠️',
        info: 'ℹ️'
    };
    console[type](`${prefix} ${styles[type]} ${message}`);
}

