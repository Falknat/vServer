/* ============================================
   Proxy Creator Component
   Управление созданием новых прокси сервисов
   ============================================ */

import { api } from '../api/wails.js';
import { configAPI } from '../api/config.js';
import { $, hide, show } from '../utils/dom.js';
import { notification } from '../ui/notification.js';
import { isWailsAvailable } from '../utils/helpers.js';
import { initCustomSelects } from '../ui/custom-select.js';

export class ProxyCreator {
    constructor() {
        this.certificates = {
            certificate: null,
            privatekey: null,
            cabundle: null
        };
    }

    open() {
        this.hideAllSections();
        show($('sectionAddProxy'));
        this.resetForm();
        this.attachEventListeners();
        setTimeout(() => initCustomSelects(), 100);
    }

    hideAllSections() {
        hide($('sectionServices'));
        hide($('sectionSites'));
        hide($('sectionProxy'));
        hide($('sectionSettings'));
        hide($('sectionVAccessEditor'));
        hide($('sectionAddSite'));
        hide($('sectionAddProxy'));
    }

    backToMain() {
        this.hideAllSections();
        show($('sectionServices'));
        show($('sectionSites'));
        show($('sectionProxy'));
    }

    resetForm() {
        $('newProxyDomain').value = '';
        $('newProxyLocalAddr').value = '127.0.0.1';
        $('newProxyLocalPort').value = '';
        $('newProxyStatus').value = 'enable';
        $('newProxyServiceHTTPS').checked = false;
        $('newProxyAutoHTTPS').checked = true;
        $('proxyCertMode').value = 'none';
        
        this.certificates = {
            certificate: null,
            privatekey: null,
            cabundle: null
        };
        
        hide($('proxyCertUploadBlock'));
        
        $('proxyCertFileStatus').innerHTML = '';
        $('proxyKeyFileStatus').innerHTML = '';
        $('proxyCaFileStatus').innerHTML = '';
        
        if ($('proxyCertFileName')) $('proxyCertFileName').textContent = 'Выберите файл...';
        if ($('proxyKeyFileName')) $('proxyKeyFileName').textContent = 'Выберите файл...';
        if ($('proxyCaFileName')) $('proxyCaFileName').textContent = 'Выберите файл...';
        
        if ($('proxyCertFile')) $('proxyCertFile').value = '';
        if ($('proxyKeyFile')) $('proxyKeyFile').value = '';
        if ($('proxyCaFile')) $('proxyCaFile').value = '';
        
        const labels = document.querySelectorAll('#sectionAddProxy .file-upload-btn');
        labels.forEach(label => label.classList.remove('file-uploaded'));
    }

    attachEventListeners() {
        const createBtn = $('createProxyBtn');
        if (createBtn) {
            createBtn.onclick = async () => await this.createProxy();
        }
        
        this.setupDragAndDrop();
    }

    setupDragAndDrop() {
        const fileWrappers = [
            { wrapper: document.querySelector('label[for="proxyCertFile"]')?.parentElement, input: $('proxyCertFile'), type: 'certificate' },
            { wrapper: document.querySelector('label[for="proxyKeyFile"]')?.parentElement, input: $('proxyKeyFile'), type: 'privatekey' },
            { wrapper: document.querySelector('label[for="proxyCaFile"]')?.parentElement, input: $('proxyCaFile'), type: 'cabundle' }
        ];

        fileWrappers.forEach(({ wrapper, input, type }) => {
            if (!wrapper || !input) return;

            ['dragenter', 'dragover', 'dragleave', 'drop'].forEach(eventName => {
                wrapper.addEventListener(eventName, (e) => {
                    e.preventDefault();
                    e.stopPropagation();
                });
            });

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

            wrapper.addEventListener('drop', (e) => {
                const files = e.dataTransfer.files;
                if (files.length > 0) {
                    const dataTransfer = new DataTransfer();
                    dataTransfer.items.add(files[0]);
                    input.files = dataTransfer.files;
                    
                    const event = new Event('change', { bubbles: true });
                    input.dispatchEvent(event);
                }
            });
        });
    }

    toggleCertUpload() {
        const mode = $('proxyCertMode')?.value;
        const block = $('proxyCertUploadBlock');
        
        if (mode === 'upload') {
            show(block);
        } else {
            hide(block);
        }
    }

    handleCertFile(input, certType) {
        const file = input.files[0];
        const statusId = certType === 'certificate' ? 'proxyCertFileStatus' :
                        certType === 'privatekey' ? 'proxyKeyFileStatus' : 'proxyCaFileStatus';
        const labelId = certType === 'certificate' ? 'proxyCertFileName' :
                       certType === 'privatekey' ? 'proxyKeyFileName' : 'proxyCaFileName';
        
        const statusDiv = $(statusId);
        const labelSpan = $(labelId);
        const labelBtn = input.nextElementSibling;
        
        if (!file) {
            this.certificates[certType] = null;
            statusDiv.innerHTML = '';
            if (labelSpan) labelSpan.textContent = 'Выберите файл...';
            if (labelBtn) labelBtn.classList.remove('file-uploaded');
            return;
        }

        if (file.size > 1024 * 1024) {
            statusDiv.innerHTML = '<span style="color: #e74c3c;"><i class="fas fa-times-circle"></i> Файл слишком большой (макс 1MB)</span>';
            this.certificates[certType] = null;
            input.value = '';
            if (labelSpan) labelSpan.textContent = 'Выберите файл...';
            if (labelBtn) labelBtn.classList.remove('file-uploaded');
            return;
        }

        if (labelSpan) labelSpan.textContent = file.name;
        if (labelBtn) labelBtn.classList.add('file-uploaded');

        const reader = new FileReader();
        reader.onload = (e) => {
            const content = e.target.result;
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

    validateForm() {
        const domain = $('newProxyDomain')?.value.trim();
        const localAddr = $('newProxyLocalAddr')?.value.trim();
        const localPort = $('newProxyLocalPort')?.value.trim();
        const certMode = $('proxyCertMode')?.value;

        if (!domain) {
            notification.error('❌ Укажите внешний домен');
            return false;
        }

        if (!localAddr) {
            notification.error('❌ Укажите локальный адрес');
            return false;
        }

        if (!localPort) {
            notification.error('❌ Укажите локальный порт');
            return false;
        }

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

    async createProxy() {
        if (!this.validateForm()) {
            return;
        }

        if (!isWailsAvailable()) {
            notification.error('Wails API недоступен');
            return;
        }

        const createBtn = $('createProxyBtn');
        const originalText = createBtn.querySelector('span').textContent;

        try {
            createBtn.disabled = true;
            createBtn.querySelector('span').textContent = 'Создание...';

            const certMode = $('proxyCertMode').value;

            const proxyData = {
                Enable: $('newProxyStatus').value === 'enable',
                ExternalDomain: $('newProxyDomain').value.trim(),
                LocalAddress: $('newProxyLocalAddr').value.trim(),
                LocalPort: $('newProxyLocalPort').value.trim(),
                ServiceHTTPSuse: $('newProxyServiceHTTPS').checked,
                AutoHTTPS: $('newProxyAutoHTTPS').checked,
                AutoCreateSSL: certMode === 'auto'
            };

            const config = await configAPI.getConfig();
            
            if (!config.Proxy_Service) {
                config.Proxy_Service = [];
            }
            
            config.Proxy_Service.push(proxyData);

            const result = await configAPI.saveConfig(JSON.stringify(config, null, 4));

            if (result.startsWith('Error')) {
                notification.error(result, 3000);
                return;
            }

            notification.success('✅ Прокси сервис создан!', 1500);

            if (certMode === 'upload') {
                createBtn.querySelector('span').textContent = 'Загрузка сертификатов...';
                
                if (this.certificates.certificate) {
                    await api.uploadCertificate(proxyData.ExternalDomain, 'certificate', this.certificates.certificate);
                }
                
                if (this.certificates.privatekey) {
                    await api.uploadCertificate(proxyData.ExternalDomain, 'privatekey', this.certificates.privatekey);
                }
                
                if (this.certificates.cabundle) {
                    await api.uploadCertificate(proxyData.ExternalDomain, 'cabundle', this.certificates.cabundle);
                }
                
                notification.success('🔒 Сертификаты загружены!', 1500);
            }

            createBtn.querySelector('span').textContent = 'Перезапуск серверов...';
            await configAPI.stopHTTPService();
            await configAPI.stopHTTPSService();
            await new Promise(resolve => setTimeout(resolve, 500));
            await configAPI.startHTTPService();
            await configAPI.startHTTPSService();

            notification.success('🚀 Серверы перезапущены! Прокси готов к работе!', 2000);

            setTimeout(() => {
                this.backToMain();
                if (window.proxyManager) {
                    window.proxyManager.load();
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
