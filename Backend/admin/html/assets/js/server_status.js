// Функция для получения HTML шаблона сервера
var patch_json = '/json/server_status.json';

function getServerTemplate(server) {
    // Определяем класс и текст статуса
    let statusClass, statusText;
    switch(server.Status.toLowerCase()) {
        case 'running':
            statusClass = 'running';
            statusText = 'Работает';
            break;
        case 'stopped':
            statusClass = 'stopped';
            statusText = 'Остановлен';
            break;
        case 'starting':
            statusClass = 'starting';
            statusText = 'Запускается';
            break;
        case 'stopping':
            statusClass = 'stopping';
            statusText = 'Завершается';
            break;
        default:
            statusClass = 'stopped';
            statusText = `Неизвестно (${server.Status})`;
    }
    
         // Определяем иконку и состояние кнопки
     let buttonIcon, buttonDisabled, buttonClass;
     
     if (statusClass === 'starting' || statusClass === 'stopping') {
         buttonIcon = '⏳';
         buttonDisabled = 'disabled';
         buttonClass = 'btn-icon disabled';
     } else if (statusClass === 'running') {
         buttonIcon = '⏹️';
         buttonDisabled = '';
         buttonClass = 'btn-icon';
     } else {
         buttonIcon = '▶️';
         buttonDisabled = '';
         buttonClass = 'btn-icon';
     }
     
     return `
         <div class="server-item">
             <div class="server-status ${statusClass}"></div>
             <div class="server-info">
                 <div class="server-name">${server.NameService}</div>
                 <div class="server-details">Port ${server.Port} - ${statusText}</div>
             </div>
             <div class="server-actions">
                 <button class="${buttonClass}" onclick="toggleServer('${server.NameService}')" ${buttonDisabled}>
                     ${buttonIcon}
                 </button>
             </div>
         </div>
     `;
}

// Создаем загрузчик серверов
const serversLoader = new JSONLoader({
    url: patch_json,
    container: '.servers-grid',
    requiredFields: ['NameService', 'Port', 'Status'],
    displayTemplate: getServerTemplate,
    errorMessage: 'Ошибка загрузки статуса серверов'
});

// Функция для показа временного статуса
function showTempStatus(serviceName, tempStatus) {
    // Ищем элемент этого сервера на странице
    const serverElements = document.querySelectorAll('.server-item');
    serverElements.forEach(element => {
        const nameElement = element.querySelector('.server-name');
        if (nameElement && nameElement.textContent === serviceName) {
            // Создаем временный объект сервера с новым статусом
            fetch(patch_json)
                .then(response => response.json())
                .then(servers => {
                    const server = servers.find(s => s.NameService === serviceName);
                    if (server) {
                        const tempServer = {...server, Status: tempStatus};
                        // Заменяем только этот элемент
                        element.outerHTML = getServerTemplate(tempServer);
                    }
                });
        }
    });
}

// Функция обновления одного сервера
function updateSingleServer(serviceName) {
    fetch(patch_json)
        .then(response => response.json())
        .then(servers => {
            const server = servers.find(s => s.NameService === serviceName);
            if (server) {
                // Ищем элемент этого сервера на странице
                const serverElements = document.querySelectorAll('.server-item');
                serverElements.forEach(element => {
                    const nameElement = element.querySelector('.server-name');
                    if (nameElement && nameElement.textContent === serviceName) {
                        // Заменяем только этот элемент
                        element.outerHTML = getServerTemplate(server);
                    }
                });
            }
        });
}

// Универсальная функция управления сервером
function serverAction(serviceName, startEndpoint, stopEndpoint, updateDelayMs) {
    // Получаем текущий статус сервера из JSON
    fetch(patch_json)
        .then(response => response.json())
        .then(servers => {
            const server = servers.find(s => s.NameService === serviceName);
            
            // Блокируем действие если сервер в процессе изменения
            if (server.Status === 'starting' || server.Status === 'stopping') {
                return;
            }
            
            if (server.Status === 'running') {
                // Сервер запущен - останавливаем
                showTempStatus(serviceName, 'stopping');
                fetch(stopEndpoint, {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    }
                }).then(() => {
                    setTimeout(() => {
                        updateSingleServer(serviceName); // Обновляем только этот сервер
                    }, updateDelayMs);
                });
            } else {
                // Сервер остановлен - запускаем
                showTempStatus(serviceName, 'starting');
                fetch(startEndpoint, {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    }
                }).then(() => {
                    setTimeout(() => {
                        updateSingleServer(serviceName); // Обновляем только этот сервер
                    }, updateDelayMs);
                });
            }
        });
}

// Функция для переключения сервера
function toggleServer(serviceName) {

    if (serviceName === 'MySQL Server') {
        serverAction('MySQL Server', '/service/MySql_Start', '/service/MySql_Stop', 2000);
    }

    if (serviceName === 'HTTP Server') {
        serverAction('HTTP Server', '/service/Http_Start', '/service/Http_Stop', 2000);
    }

    if (serviceName === 'HTTPS Server') {
        serverAction('HTTPS Server', '/service/Https_Start', '/service/Https_Stop', 2000);
    }

    if (serviceName === 'PHP Server') {
        serverAction('PHP Server', '/service/Php_Start', '/service/Php_Stop', 2000);
    }
    
}

// Запуск при загрузке страницы
document.addEventListener('DOMContentLoaded', function() {
    serversLoader.load();
}); 
