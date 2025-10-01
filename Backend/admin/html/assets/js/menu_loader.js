// Функция для создания HTML пункта меню
function getMenuTemplate(item) {
    const isActive = item.active ? 'active' : '';
    
    return `
        <li class="nav-item ${isActive}">
            <a href="${item.url}" class="nav-link">
                <span class="nav-icon">${item.icon}</span>
                <span class="nav-text">${item.name}</span>
            </a>
        </li>
    `;
}

// Создаем загрузчик меню
const menuLoader = new JSONLoader({
    url: '/json/menu.json',
    container: '.nav-menu',
    requiredFields: ['name', 'icon', 'url', 'active'],
    displayTemplate: getMenuTemplate,
    errorMessage: 'Ошибка загрузки меню'
});

// Запуск при загрузке страницы
document.addEventListener('DOMContentLoaded', function() {
    menuLoader.load();
}); 