/* ============================================
   Site Creator Component
   Управление созданием новых сайтов
   ============================================ */

import { api } from '../api/wails.js';
import { configAPI } from '../api/config.js';
import { $, hide, show } from '../utils/dom.js';
import { notification } from '../ui/notification.js';
import { isWailsAvailable } from '../utils/helpers.js';
import { initCustomSelects } from '../ui/custom-select.js';

/**
 * Класс для создания новых сайтов
 */
export class SiteCreator {
    constructor() {
        this.aliases = [];
        this.certificates = {
            certificate: null,
            privatekey: null,
            cabundle: null
        };
    }

    /**
     * Открыть страницу создания сайта
     */
    open() {
        // Скрываем все секции
        this.hideAllSections();
        
        // Показываем страницу создания
        show($('sectionAddSite'));
        
        // Очищаем форму
        this.resetForm();
        
        // Привязываем обработчики
        this.attachEventListeners();
        
        // Инициализируем кастомные select'ы
        setTimeout(() => initCustomSelects(), 100);
    }

    /**
     * Скрыть все секции
     */
    hideAllSections() {
        hide($('sectionServices'));
        hide($('sectionSites'));
        hide($('sectionProxy'));
        hide($('sectionSettings'));
        hide($('sectionVAccessEditor'));
        hide($('sectionAddSite'));
    }

    /**
     * Вернуться на главную
     */
    backToMain() {
        this.hideAllSections();
        show($('sectionServices'));
        show($('sectionSites'));
        show($('sectionProxy'));
    }

    /**
     * Очистить форму
     */
    resetForm() {
        $('newSiteName').value = '';
        $('newSiteHost').value = '';
        $('newSiteAliasInput').value = '';
        $('newSiteRootFile').value = 'index.html';
        $('newSiteStatus').value = 'active';
        $('newSiteRouting').checked = true;
        $('certMode').value = 'none';
        
        this.aliases = [];
        this.certificates = {
            certificate: null,
            privatekey: null,
            cabundle: null
        };
        
        // Скрываем блок загрузки сертификатов
        hide($('certUploadBlock'));
        
        // Очищаем статусы файлов
        $('certFileStatus').innerHTML = '';
        $('keyFileStatus').innerHTML = '';
        $('caFileStatus').innerHTML = '';
        
        // Очищаем labels файлов
        if ($('certFileName')) $('certFileName').textContent = 'Выберите файл...';
        if ($('keyFileName')) $('keyFileName').textContent = 'Выберите файл...';
        if ($('caFileName')) $('caFileName').textContent = 'Выберите файл...';
        
        // Очищаем input файлов
        if ($('certFile')) $('certFile').value = '';
        if ($('keyFile')) $('keyFile').value = '';
        if ($('caFile')) $('caFile').value = '';
        
        // Убираем класс uploaded
        const labels = document.querySelectorAll('.file-upload-btn');
        labels.forEach(label => label.classList.remove('file-uploaded'));
    }

    /**
     * Привязать обработчики событий
     */
    attachEventListeners() {
        const createBtn = $('createSiteBtn');
        if (createBtn) {
            createBtn.onclick = async () => await this.createSite();
        }
        
        // Drag & Drop для файлов сертификатов
        this.setupDragAndDrop();
    }

    /**
     * Настроить Drag & Drop для файлов
     */
    setupDragAndDrop() {
        const fileWrappers = [
            { wrapper: document.querySelector('label[for="certFile"]')?.parentElement, input: $('certFile'), type: 'certificate' },
            { wrapper: document.querySelector('label[for="keyFile"]')?.parentElement, input: $('keyFile'), type: 'privatekey' },
            { wrapper: document.querySelector('label[for="caFile"]')?.parentElement, input: $('caFile'), type: 'cabundle' }
        ];

        fileWrappers.forEach(({ wrapper, input, type }) => {
            if (!wrapper || !input) return;

            // Предотвращаем стандартное поведение
            ['dragenter', 'dragover', 'dragleave', 'drop'].forEach(eventName => {
                wrapper.addEventListener(eventName, (e) => {
                    e.preventDefault();
                    e.stopPropagation();
                });
            });

            // Подсветка при наведении файла
            ['dragenter', 'dragover'].forEach(eventName => {
                wrapper.addEventListener(eventName, () => {
                    wrapper.classList.add('drag-over');
                });
            });

            ['dragleave', 'drop'].forEach(eventName => {
                wrapper.addEventListener(eventName, () => {
                    wrapper.classList.remove('drag-over');
                });
            });

            // Обработка dropped файла
            wrapper.addEventListener('drop', (e) => {
                const files = e.dataTransfer.files;
                if (files.length > 0) {
                    // Создаём объект DataTransfer и присваиваем файлы input'у
                    const dataTransfer = new DataTransfer();
                    dataTransfer.items.add(files[0]);
                    input.files = dataTransfer.files;
                    
                    // Триггерим событие change
                    const event = new Event('change', { bubbles: true });
                    input.dispatchEvent(event);
                }
            });
        });
    }

    /**
     * Парсить aliases из строки (через запятую)
     */
    parseAliases() {
        const input = $('newSiteAliasInput');
        const value = input?.value.trim();
        
        if (!value) {
            this.aliases = [];
            return;
        }
        
        // Разделяем по запятой и очищаем
        this.aliases = value
            .split(',')
            .map(alias => alias.trim())
            .filter(alias => alias.length > 0);
    }

    /**
     * Переключить видимость блока загрузки сертификатов
     */
    toggleCertUpload() {
        const mode = $('certMode')?.value;
        const block = $('certUploadBlock');
        
        if (mode === 'upload') {
            show(block);
        } else {
            hide(block);
        }
    }

    /**
     * Обработать выбор файла сертификата
     */
    handleCertFile(input, certType) {
        const file = input.files[0];
        const statusId = certType === 'certificate' ? 'certFileStatus' :
                        certType === 'privatekey' ? 'keyFileStatus' : 'caFileStatus';
        const labelId = certType === 'certificate' ? 'certFileName' :
                       certType === 'privatekey' ? 'keyFileName' : 'caFileName';
        
        const statusDiv = $(statusId);
        const labelSpan = $(labelId);
        const labelBtn = input.nextElementSibling; // label элемент
        
        if (!file) {
            this.certificates[certType] = null;
            statusDiv.innerHTML = '';
            if (labelSpan) labelSpan.textContent = 'Выберите файл...';
            if (labelBtn) labelBtn.classList.remove('file-uploaded');
            return;
        }

        // Проверяем размер файла (макс 1MB)
        if (file.size > 1024 * 1024) {
            statusDiv.innerHTML = '<span style="color: #e74c3c;"><i class="fas fa-times-circle"></i> Файл слишком большой (макс 1MB)</span>';
            this.certificates[certType] = null;
            input.value = '';
            if (labelSpan) labelSpan.textContent = 'Выберите файл...';
            if (labelBtn) labelBtn.classList.remove('file-uploaded');
            return;
        }

        // Обновляем UI
        if (labelSpan) labelSpan.textContent = file.name;
        if (labelBtn) labelBtn.classList.add('file-uploaded');

        // Читаем файл
        const reader = new FileReader();
        reader.onload = (e) => {
            const content = e.target.result;
            // Сохраняем как base64
            this.certificates[certType] = btoa(content);
            statusDiv.innerHTML = `<span style="color: #2ecc71;"><i class="fas fa-check-circle"></i> Загружен успешно</span>`;
        };
        reader.onerror = () => {
            statusDiv.innerHTML = '<span style="color: #e74c3c;"><i class="fas fa-times-circle"></i> Ошибка чтения файла</span>';
            this.certificates[certType] = null;
            if (labelSpan) labelSpan.textContent = 'Выберите файл...';
            if (labelBtn) labelBtn.classList.remove('file-uploaded');
        };
        reader.readAsText(file);
    }

    /**
     * Валидация формы
     */
    validateForm() {
        const name = $('newSiteName')?.value.trim();
        const host = $('newSiteHost')?.value.trim();
        const rootFile = $('newSiteRootFile')?.value;
        const certMode = $('certMode')?.value;

        if (!name) {
            notification.error('❌ Укажите название сайта');
            return false;
        }

        if (!host) {
            notification.error('❌ Укажите host (домен)');
            return false;
        }

        if (!rootFile) {
            notification.error('❌ Укажите root файл');
            return false;
        }

        // Проверка сертификатов если режим загрузки
        if (certMode === 'upload') {
            if (!this.certificates.certificate) {
                notification.error('❌ Загрузите файл certificate.crt');
                return false;
            }
            if (!this.certificates.privatekey) {
                notification.error('❌ Загрузите файл private.key');
                return false;
            }
        }

        return true;
    }

    /**
     * Создать сайт
     */
    async createSite() {
        if (!this.validateForm()) {
            return;
        }

        if (!isWailsAvailable()) {
            notification.error('Wails API недоступен');
            return;
        }

        const createBtn = $('createSiteBtn');
        const originalText = createBtn.querySelector('span').textContent;

        try {
            createBtn.disabled = true;
            createBtn.querySelector('span').textContent = 'Создание...';

            // Парсим aliases из поля ввода
            this.parseAliases();

            // Собираем данные сайта
            const siteData = {
                name: $('newSiteName').value.trim(),
                host: $('newSiteHost').value.trim(),
                alias: this.aliases,
                status: $('newSiteStatus').value,
                root_file: $('newSiteRootFile').value,
                root_file_routing: $('newSiteRouting').checked
            };

            // Создаём сайт
            const siteJSON = JSON.stringify(siteData);
            const result = await api.createNewSite(siteJSON);

            if (result.startsWith('Error')) {
                notification.error(result, 3000);
                return;
            }

            notification.success('✅ Сайт успешно создан!', 1500);

            // Загружаем сертификаты если нужно
            const certMode = $('certMode').value;
            if (certMode === 'upload') {
                createBtn.querySelector('span').textContent = 'Загрузка сертификатов...';
                
                // Загружаем certificate
                if (this.certificates.certificate) {
                    await api.uploadCertificate(siteData.host, 'certificate', this.certificates.certificate);
                }
                
                // Загружаем private key
                if (this.certificates.privatekey) {
                    await api.uploadCertificate(siteData.host, 'privatekey', this.certificates.privatekey);
                }
                
                // Загружаем ca bundle если есть
                if (this.certificates.cabundle) {
                    await api.uploadCertificate(siteData.host, 'cabundle', this.certificates.cabundle);
                }
                
                notification.success('🔒 Сертификаты загружены!', 1500);
            }

            // Перезапускаем HTTP/HTTPS
            createBtn.querySelector('span').textContent = 'Перезапуск серверов...';
            await configAPI.stopHTTPService();
            await configAPI.stopHTTPSService();
            await new Promise(resolve => setTimeout(resolve, 500));
            await configAPI.startHTTPService();
            await configAPI.startHTTPSService();

            notification.success('🚀 Серверы перезапущены! Сайт готов к работе!', 2000);

            // Возвращаемся на главную
            setTimeout(() => {
                this.backToMain();
                // Перезагружаем список сайтов
                if (window.sitesManager) {
                    window.sitesManager.load();
                }
            }, 1000);

        } catch (error) {
            notification.error('Ошибка: ' + error.message, 3000);
        } finally {
            createBtn.disabled = false;
            createBtn.querySelector('span').textContent = originalText;
        }
    }
}

