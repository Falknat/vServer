/* Класс для показа уведомлений */
class MessageUp {
    // Время автозакрытия уведомлений (в миллисекундах)
    static autoCloseTime = 3000;

    // Типы уведомлений (легко редактировать тут)
    static TYPES = {
        info: { 
            borderColor: 'rgb(52, 152, 219)',
            background: 'linear-gradient(135deg, rgb(30, 50, 70), rgb(35, 55, 75))'
        },
        success: { 
            borderColor: 'rgb(46, 204, 113)',
            background: 'linear-gradient(135deg, rgb(25, 55, 35), rgb(30, 60, 40))'
        },
        warning: { 
            borderColor: 'rgb(243, 156, 18)',
            background: 'linear-gradient(135deg, rgb(60, 45, 20), rgb(65, 50, 25))'
        },
        error: { 
            borderColor: 'rgb(231, 76, 60)',
            background: 'linear-gradient(135deg, rgb(60, 25, 25), rgb(65, 30, 30))'
        }
    };

    constructor() {
        this.addStyles();
    }

    /* Показать уведомление */
    show(message, type = 'info') {
        const notification = document.createElement('div');
        notification.className = `message-up message-up-${type}`;
        
        notification.innerHTML = `
            <div class="message-content">
                <span class="message-text">${message}</span>
            </div>
        `;

        // Вычисляем позицию для нового уведомления
        const existingNotifications = document.querySelectorAll('.message-up');
        let topPosition = 2; // начальная позиция в rem
        
        existingNotifications.forEach(existing => {
            const rect = existing.getBoundingClientRect();
            const currentTop = parseFloat(existing.style.top) || 2;
            const height = rect.height / 16; // переводим px в rem (примерно)
            topPosition = Math.max(topPosition, currentTop + height + 1); // добавляем отступ
        });
        
        notification.style.top = `${topPosition}rem`;
        document.body.appendChild(notification);

        // Показываем с анимацией
        setTimeout(() => {
            notification.classList.add('message-up-show');
        }, 10);

        // Автоматическое закрытие и удаление
        if (MessageUp.autoCloseTime > 0) {
            setTimeout(() => {
                if (notification && notification.parentNode) {
                    notification.classList.remove('message-up-show');
                    notification.classList.add('message-up-hide');
                    
                    setTimeout(() => {
                        if (notification.parentNode) {
                            notification.remove();
                            // Пересчитываем позиции оставшихся уведомлений
                            this.repositionNotifications();
                        }
                    }, 300);
                }
            }, MessageUp.autoCloseTime);
        }
    }

    /* Пересчитать позиции всех уведомлений */
    repositionNotifications() {
        const notifications = document.querySelectorAll('.message-up');
        let currentTop = 2; // начальная позиция
        
        notifications.forEach(notification => {
            notification.style.transition = 'all 0.3s ease';
            notification.style.top = `${currentTop}rem`;
            
            const rect = notification.getBoundingClientRect();
            const height = rect.height / 16; // переводим px в rem
            currentTop += height + 1; // добавляем отступ
        });
    }

    /* Показать сообщение напрямую или привязать к элементам */
    send(messageOrSelector, typeOrMessage = 'info', type = 'info') {
        // Если первый параметр строка и нет других элементов на странице с таким селектором
        // то показываем сообщение напрямую
        if (typeof messageOrSelector === 'string' && 
            document.querySelectorAll(messageOrSelector).length === 0) {
            // Показываем сообщение напрямую
            this.show(messageOrSelector, typeOrMessage);
            return;
        }
        
        // Иначе привязываем к элементам (старый способ)
        this.bindToElements(messageOrSelector, typeOrMessage, type);
    }

    /* Привязать уведомления к элементам */
    bindToElements(selector, message = 'Страница в разработке', type = 'info') {
        document.querySelectorAll(selector).forEach(element => {
            element.addEventListener('click', function(e) {
                e.preventDefault();
                e.stopPropagation();
                
                // Вызываем нужный тип уведомления
                window.messageUp.show(message, type);
            });
        });
    }

    /* Добавить стили для уведомлений */
    addStyles() {
        if (document.querySelector('#message-up-styles')) return;

        const cfg = MessageUp.TYPES;
        
        // Генерируем стили для типов уведомлений
        const typeStyles = Object.entries(cfg).map(([type, style]) => `
            .message-up-${type} {
                border-color: ${style.borderColor};
                background: ${style.background};
            }
        `).join('');

        const style = document.createElement('style');
        style.id = 'message-up-styles';
        style.textContent = `
            .message-up {
                position: fixed;
                top: 2rem;
                right: 2rem;
                min-width: 300px;
                max-width: 400px;
                background: rgb(26, 37, 47);
                backdrop-filter: blur(15px);
                border-radius: 12px;
                padding: 1rem;
                color: #ecf0f1;
                font-size: 0.9rem;
                border: 2px solid;
                box-shadow: 0 8px 32px rgba(0, 0, 0, 0.3);
                z-index: 10000;
                opacity: 0;
                transform: translateX(100%);
                transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
                font-family: 'Segoe UI', system-ui, sans-serif;
            }

            .message-up-show {
                opacity: 1;
                transform: translateX(0);
            }

            .message-up-hide {
                opacity: 0;
                transform: translateX(100%);
            }

            ${typeStyles}

            .message-content {
                display: flex;
                align-items: center;
                gap: 1rem;
            }

            .message-text {
                flex: 1;
                line-height: 1.4;
                font-weight: 500;
            }

            @media (max-width: 768px) {
                .message-up {
                    top: 1rem;
                    right: 1rem;
                    left: 1rem;
                    min-width: auto;
                    max-width: none;
                }
            }
        `;
        
        document.head.appendChild(style);
    }

}

// Создаем глобальный экземпляр
window.messageUp = new MessageUp(); 