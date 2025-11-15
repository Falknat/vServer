/* ============================================
   Helper Utilities
   Вспомогательные функции
   ============================================ */

// Ждёт указанное время
export function sleep(ms) {
    return new Promise(resolve => setTimeout(resolve, ms));
}

// Debounce функция
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

// Проверяет доступность Wails API
export function isWailsAvailable() {
    return typeof window.go !== 'undefined' && 
           window.go?.admin?.App !== undefined;
}

