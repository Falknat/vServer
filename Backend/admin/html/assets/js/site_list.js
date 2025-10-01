document.addEventListener('DOMContentLoaded', function() {
    const sitesList = document.querySelector('.sites-list');
    if (sitesList) {
        fetch('/service/Site_List')
        .then(r => r.json())
        .then(data => {
            const sites = data.sites || [];
            
            // Генерируем статистику
            updateSiteStats(sites);
            
            // Отображаем список сайтов
            sitesList.innerHTML = sites.map(site => `
                <div class="site-item">
                    <div class="site-status ${site.status === 'active' ? 'active' : 'inactive'}"></div>
                    <div class="site-info">
                        <span class="site-name">${site.host}</span>
                        <span class="site-details">${site.type.toUpperCase()} • Протокол</span>
                    </div>
                    <div class="site-actions">
                        <button class="btn-icon" title="Настройки">⚙️</button>
                    </div>
                </div>
            `).join('');
        });
    }
});

function updateSiteStats(sites) {
    const totalSites = sites.length;
    const activeSites = sites.filter(site => site.status === 'active').length;
    const inactiveSites = totalSites - activeSites;
    
    // Находим контейнер статистики
    const statsRow = document.querySelector('.stats-row');
    
    // Создаём всю статистику через JavaScript
    statsRow.innerHTML = `
        <div class="stat-item">
            <div class="stat-number">${totalSites}</div>
            <div class="stat-label">Всего сайтов</div>
        </div>
        <div class="stat-item">
            <div class="stat-number active-stat">${activeSites}</div>
            <div class="stat-label">Активных</div>
        </div>
        <div class="stat-item">
            <div class="stat-number inactive-stat">${inactiveSites}</div>
            <div class="stat-label">Неактивных</div>
        </div>
    `;
}



