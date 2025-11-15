/* ============================================
   Custom Select Component
   Кастомные выпадающие списки
   ============================================ */

import { $ } from '../utils/dom.js';

// Инициализация всех кастомных select'ов на странице
export function initCustomSelects() {
    const selects = document.querySelectorAll('select.form-input');
    selects.forEach(select => {
        if (!select.dataset.customized) {
            createCustomSelect(select);
        }
    });
}

// Создать кастомный select из нативного
function createCustomSelect(selectElement) {
    // Помечаем как обработанный
    selectElement.dataset.customized = 'true';
    
    // Создаём контейнер
    const wrapper = document.createElement('div');
    wrapper.className = 'custom-select';
    
    // Получаем выбранное значение
    const selectedOption = selectElement.options[selectElement.selectedIndex];
    const selectedText = selectedOption ? selectedOption.text : '';
    
    // Создаём кнопку (видимая часть)
    const button = document.createElement('div');
    button.className = 'custom-select-trigger';
    button.innerHTML = `
        <span class="custom-select-value">${selectedText}</span>
        <i class="fas fa-chevron-down custom-select-arrow"></i>
    `;
    
    // Создаём выпадающий список
    const dropdown = document.createElement('div');
    dropdown.className = 'custom-select-dropdown';
    
    // Заполняем опции
    Array.from(selectElement.options).forEach((option, index) => {
        const item = document.createElement('div');
        item.className = 'custom-select-option';
        item.textContent = option.text;
        item.dataset.value = option.value;
        item.dataset.index = index;
        
        if (option.selected) {
            item.classList.add('selected');
        }
        
        // Клик по опции
        item.addEventListener('click', () => {
            selectOption(selectElement, wrapper, item, index);
        });
        
        dropdown.appendChild(item);
    });
    
    // Клик по кнопке - открыть/закрыть
    button.addEventListener('click', (e) => {
        e.stopPropagation();
        toggleDropdown(wrapper);
    });
    
    // Собираем вместе
    wrapper.appendChild(button);
    wrapper.appendChild(dropdown);
    
    // Скрываем оригинальный select
    selectElement.style.display = 'none';
    
    // Вставляем кастомный select после оригинального
    selectElement.parentNode.insertBefore(wrapper, selectElement.nextSibling);
    
    // Закрываем при клике вне
    document.addEventListener('click', (e) => {
        if (!wrapper.contains(e.target)) {
            closeDropdown(wrapper);
        }
    });
}

// Открыть/закрыть dropdown
function toggleDropdown(wrapper) {
    const isOpen = wrapper.classList.contains('open');
    
    // Закрываем все открытые
    document.querySelectorAll('.custom-select.open').forEach(el => {
        el.classList.remove('open');
    });
    
    if (!isOpen) {
        wrapper.classList.add('open');
    }
}

// Закрыть dropdown
function closeDropdown(wrapper) {
    wrapper.classList.remove('open');
}

// Выбрать опцию
function selectOption(selectElement, wrapper, optionElement, index) {
    // Обновляем оригинальный select
    selectElement.selectedIndex = index;
    
    // Триггерим событие change
    const event = new Event('change', { bubbles: true });
    selectElement.dispatchEvent(event);
    
    // Обновляем UI
    const valueSpan = wrapper.querySelector('.custom-select-value');
    valueSpan.textContent = optionElement.textContent;
    
    // Убираем selected у всех опций
    wrapper.querySelectorAll('.custom-select-option').forEach(opt => {
        opt.classList.remove('selected');
    });
    
    // Добавляем selected к выбранной
    optionElement.classList.add('selected');
    
    // Закрываем dropdown
    closeDropdown(wrapper);
}

