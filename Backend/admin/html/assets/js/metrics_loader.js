// Функция для создания HTML карточки метрики
function createMetricCard(type, icon, name, data) {
    // Используем уже вычисленный процент для всех типов
    let value = Math.round(data.usage || data.usage_percent || 0);
    const progressClass = type === 'cpu' ? 'cpu-progress' : 
                         type === 'ram' ? 'ram-progress' : 'disk-progress';
    
    // Определяем детали для каждого типа
    let details = '';
    if (type === 'cpu') {
        details = `
            <span class="metric-info">${data.model_name || 'CPU'}</span>
            <span class="metric-frequency">${data.frequency || ''} MHz</span>
        `;
         } else if (type === 'ram') {
         const usedGb = parseFloat(data.used_gb) || 0;
         details = `
             <span class="metric-info">${usedGb.toFixed(1)} GB из ${data.total_gb || 0} GB</span>
             <span class="metric-type">${data.type || 'RAM'}</span>
         `;
         } else if (type === 'disk') {
         const usedGb = parseFloat(data.used_gb) || 0;
         const freeGb = parseFloat(data.free_gb) || 0;
         details = `
             <span class="metric-info">Занято: ${usedGb.toFixed(0)} GB : Свободно: ${freeGb.toFixed(0)} GB</span>
             <span class="metric-speed">Размер: ${data.total_gb || 0}</span>
         `;
    }

    return `
        <div class="metric-card ${type}">
            <div class="metric-icon-wrapper">
                <span class="metric-icon">${icon}</span>
            </div>
            <div class="metric-content">
                <div class="metric-header">
                    <span class="metric-name">${name}</span>
                    <span class="metric-value">${value}%</span>
                </div>
                <div class="metric-progress-bar">
                    <div class="metric-progress ${progressClass}" style="width: ${value}%"></div>
                </div>
                <div class="metric-details">
                    ${details}
                </div>
            </div>
        </div>
    `;
}

// Функция для создания всех метрик
function renderMetrics(data) {
    const container = document.querySelector('.metrics-grid');
    if (!container) return;

    const html = [
        createMetricCard('cpu', '🖥️', 'Процессор', data.cpu || {}),
        createMetricCard('ram', '💾', 'Оперативная память', data.memory || {}),
        createMetricCard('disk', '💿', data.disk.type, data.disk || {})
    ].join('');

    container.innerHTML = html;
}

// Данные по умолчанию (будут заменены данными из API)
const staticData = {
    cpu: {
        usage: 0,
        model_name: 'Загрузка...',
        frequency: '0',
        cores: '0'
    },
    memory: {
        usage_percent: 0,
        used_gb: 0,
        total_gb: 0,
        type: 'Загрузка...'
    },
    disk: {
        usage_percent: 0,
        used_gb: 0,
        free_gb: 0,
        total_gb: '0',
        type: 'Загрузка...',
        read_speed: '520 MB/s'
    }
};




// Функция обновления метрик
async function updateMetrics() {
    try {
        const response = await fetch('/api/metrics');
        const data = await response.json();
        
        // Обновляем все данные из API
        if (data.cpu_name) staticData.cpu.model_name = data.cpu_name;
        if (data.cpu_ghz) staticData.cpu.frequency = data.cpu_ghz;
        if (data.cpu_cores) staticData.cpu.cores = data.cpu_cores;
        if (data.cpu_usage) staticData.cpu.usage = parseInt(data.cpu_usage);
        
        if (data.disk_name) staticData.disk.type = data.disk_name;
        if (data.disk_size) staticData.disk.total_gb = data.disk_size;
        if (data.disk_used) staticData.disk.used_gb = parseFloat(data.disk_used);
        if (data.disk_free) staticData.disk.free_gb = parseFloat(data.disk_free);

        if (data.ram_using) staticData.memory.used_gb = parseFloat(data.ram_using);
        if (data.ram_total) staticData.memory.total_gb = parseFloat(data.ram_total);
        
        // Вычисляем процент использования памяти
        if (staticData.memory.used_gb && staticData.memory.total_gb) {
            staticData.memory.usage_percent = Math.round((staticData.memory.used_gb / staticData.memory.total_gb) * 100);
        }
        
        // Вычисляем процент использования диска
        if (staticData.disk.used_gb && staticData.disk.total_gb) {
            const used = parseFloat(staticData.disk.used_gb.toString().replace(' GB', '')) || 0;
            const total = parseFloat(staticData.disk.total_gb.toString().replace(' GB', '')) || 1;
            staticData.disk.usage_percent = Math.round((used / total) * 100);
        }
        
        // Обновляем uptime
        if (data.server_uptime) {
            const uptimeElement = document.querySelector('.uptime-value');
            if (uptimeElement) {
                uptimeElement.textContent = data.server_uptime;
            }
        }
        
        // Перерисовываем только метрики
        renderMetrics(staticData);
    } catch (error) {
        console.error('Ошибка обновления метрик:', error);
    }

    
}

// Показываем статические данные когда DOM загружен
document.addEventListener('DOMContentLoaded', function() {
    renderMetrics(staticData);
    
    // Сразу загружаем актуальные данные
    updateMetrics();
    
    // Обновляем метрики каждые 5 секунд
    setInterval(updateMetrics, 5000);
});



