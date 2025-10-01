// Универсальный класс для загрузки JSON данных
class JSONLoader {
    constructor(config) {
        this.url = config.url;
        this.container = config.container;
        this.requiredFields = config.requiredFields || [];
        this.displayTemplate = config.displayTemplate;
        this.errorMessage = config.errorMessage || 'Ошибка загрузки данных';
    }

    // Загрузка данных
    async load() {
        try {
            const response = await fetch(this.url);
            if (!response.ok) {
                throw new Error(`HTTP ${response.status}`);
            }
            
            const data = await response.json();
            
            // Проверяем структуру данных
            if (!this.validateData(data)) {
                throw new Error('Неверная структура данных');
            }
            
            this.display(data);
        } catch (error) {
            console.error('Ошибка загрузки:', error);
            this.displayError();
        }
    }

    // Простая проверка обязательных полей
    validateData(data) {
        if (!Array.isArray(data)) {
            return false;
        }
        
        for (let item of data) {
            for (let field of this.requiredFields) {
                if (!item.hasOwnProperty(field)) {
                    return false;
                }
            }
        }
        
        return true;
    }

    // Отображение данных по шаблону
    display(data) {
        const container = document.querySelector(this.container);
        container.innerHTML = '';
        
        data.forEach(item => {
            let html;
            
            // Если шаблон - функция, вызываем её
            if (typeof this.displayTemplate === 'function') {
                html = this.displayTemplate(item);
            } else {
                // Иначе используем строковый шаблон с заменой
                html = this.displayTemplate;
                for (let key in item) {
                    const value = item[key];
                    html = html.replace(new RegExp(`{{${key}}}`, 'g'), value);
                }
            }
            
            container.innerHTML += html;
        });
    }

    // Отображение ошибки
    displayError() {
        const container = document.querySelector(this.container);
        if (container) {
            container.innerHTML = `
                <div style="text-align: center; padding: 2rem; color: #ff6b6b;">
                    ⚠️ ${this.errorMessage}
                </div>
            `;
        }
    }

    // Перезагрузка данных
    reload() {
        this.load();
    }
} 