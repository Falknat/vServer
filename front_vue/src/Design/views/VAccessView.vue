<script setup>

const { t } = useI18n()
const route = useRoute()
const { success } = useNotification()

const props = defineProps({
  host: { type: String, required: true },
})

const isProxy = computed(() => route.query.proxy === 'true')

const activeTab = ref('rules')
const rules = ref([])
const loading = ref(true)

const { dragIndex, dragOverIndex, onDragStart, onDragOver, onDragEnter, onDragLeave, onDrop, onDragEnd } = useDraggable(rules)

onMounted(async () => {
  const data = await api.getVAccessRules(props.host, isProxy.value)
  rules.value = data?.rules || []
  loading.value = false
})

const saveRules = async () => {
  await api.saveVAccessRules(props.host, isProxy.value, JSON.stringify({ rules: rules.value }))
  success(t('notify.changesSaved'))
}

const addRule = () => {
  rules.value.push({
    type: 'Disable',
    type_file: [],
    path_access: [],
    ip_list: [],
    exceptions_dir: [],
    url_error: '404',
  })
}

const removeRule = (index) => {
  rules.value.splice(index, 1)
}

const formatList = (arr) => {
  if (!arr || arr.length === 0) return '—'
  return arr.join(', ')
}
</script>

<template>
  <div class="vaccess-page">
    <Breadcrumbs :items="[host]">
      <template #tabs>
        <button class="vaccess-tab" :class="{ active: activeTab === 'rules' }" @click="activeTab = 'rules'">
          <i class="fas fa-list"></i>
          <span>{{ t('vaccess.rulesTab') }}</span>
        </button>
        <button class="vaccess-tab" :class="{ active: activeTab === 'help' }" @click="activeTab = 'help'">
          <i class="fas fa-question-circle"></i>
          <span>{{ t('vaccess.helpTab') }}</span>
        </button>
      </template>
    </Breadcrumbs>

    <PageHeader icon="fas fa-shield-alt" :title="t('vaccess.title')" :subtitle="`${t('vaccess.subtitle')} — ${host}`">
      <template #actions>
        <VButton variant="success" icon="fas fa-save" @click="saveRules">{{ t('vaccess.save') }}</VButton>
        <VButton icon="fas fa-plus" @click="addRule">{{ t('vaccess.addRule') }}</VButton>
      </template>
    </PageHeader>

    <!-- Вкладка: Правила -->
    <div v-if="activeTab === 'rules'" class="vaccess-tab-content">
      <div v-if="rules.length === 0 && !loading" class="vaccess-empty">
        <div class="empty-icon"><i class="fas fa-shield-alt"></i></div>
        <h3>{{ t('vaccess.empty') }}</h3>
        <p>{{ t('vaccess.emptyDesc') }}</p>
        <VButton icon="fas fa-plus" @click="addRule">{{ t('vaccess.createRule') }}</VButton>
      </div>

      <div v-else class="vaccess-rules-container">
        <table class="vaccess-table">
          <thead>
            <tr>
              <th class="col-drag"></th>
              <th class="col-type"><span class="th-with-info">{{ t('vaccess.type') }} <VTooltip text="Allow — разрешить доступ, Disable — запретить. Клик для переключения" /></span></th>
              <th class="col-files"><span class="th-with-info">{{ t('vaccess.files') }} <VTooltip text="Расширения файлов через запятую" :items="['*.php', '*.exe', '*.sh', 'no_extension']" /></span></th>
              <th class="col-paths"><span class="th-with-info">{{ t('vaccess.paths') }} <VTooltip text="Пути доступа через запятую" :items="['/admin/*', '/api/*', '/uploads/*']" /></span></th>
              <th class="col-ips"><span class="th-with-info">{{ t('vaccess.ips') }} <VTooltip text="IP адреса через запятую" :items="['192.168.1.1', '10.0.0.0/24', '127.0.0.1']" /></span></th>
              <th class="col-exceptions"><span class="th-with-info">{{ t('vaccess.exceptions') }} <VTooltip text="Пути-исключения: правило НЕ применяется к ним" :items="['/public/*', '/bot/*', '/api/open/*']" /></span></th>
              <th class="col-error"><span class="th-with-info">{{ t('vaccess.error') }} <VTooltip text="Куда перенаправить при блокировке" :items="['404', 'https://site.com', '/error.html']" /></span></th>
              <th class="col-actions"></th>
            </tr>
          </thead>
          <tbody>
            <tr
              v-for="(rule, index) in rules"
              :key="index"
              draggable="true"
              :class="{ 'drag-over': dragOverIndex === index, 'dragging': dragIndex === index }"
              @dragstart="onDragStart(index, $event)"
              @dragover="onDragOver(index, $event)"
              @dragenter="onDragEnter(index, $event)"
              @dragleave="onDragLeave"
              @drop="onDrop(index)"
              @dragend="onDragEnd($event)"
            >
              <td class="drag-handle"><i class="fas fa-grip-vertical"></i></td>
              <td>
                <VBadge class="type-toggle" :variant="rule.type === 'Allow' ? 'yes' : 'no'" @click="rule.type = rule.type === 'Allow' ? 'Disable' : 'Allow'">
                  {{ rule.type }}
                </VBadge>
              </td>
              <td>
                <FieldEditor v-model="rule.type_file" placeholder="*.php" />
              </td>
              <td>
                <FieldEditor v-model="rule.path_access" placeholder="/admin/*" />
              </td>
              <td>
                <FieldEditor v-model="rule.ip_list" placeholder="192.168.1.1" />
              </td>
              <td>
                <FieldEditor v-model="rule.exceptions_dir" placeholder="/public/*" />
              </td>
              <td>
                <input v-model="rule.url_error" class="inline-input" placeholder="404" />
              </td>
              <td class="col-actions-cell">
                <button class="icon-btn-small" @click="removeRule(index)"><i class="fas fa-trash"></i></button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Вкладка: Инструкция -->
    <div v-if="activeTab === 'help'" class="vaccess-tab-content">
      <div class="vaccess-help">
        <div class="help-card">
          <h3><i class="fas fa-info-circle"></i> {{ t('vaccess.helpPrinciple') }}</h3>
          <ul>
            <li>Правила проверяются <strong>сверху вниз</strong> по порядку</li>
            <li>Первое подходящее правило срабатывает и завершает проверку</li>
            <li>Если ни одно правило не сработало — доступ <strong>разрешён</strong></li>
            <li>Перетаскивайте правила за <i class="fas fa-grip-vertical"></i> чтобы изменить порядок</li>
          </ul>
        </div>

        <div class="help-card">
          <h3><i class="fas fa-sliders-h"></i> {{ t('vaccess.helpParams') }}</h3>
          <div class="help-params">
            <div class="help-param">
              <strong>type:</strong>
              <p><VBadge variant="yes">Allow</VBadge> (разрешить) или <VBadge variant="no">Disable</VBadge> (запретить)</p>
            </div>
            <div class="help-param">
              <strong>Расширения файлов:</strong>
              <p>Список расширений через запятую (<code>*.php</code>, <code>*.exe</code>)</p>
            </div>
            <div class="help-param">
              <strong>Пути доступа:</strong>
              <p>Список путей через запятую (<code>/admin/*</code>, <code>/api/*</code>)</p>
            </div>
            <div class="help-param">
              <strong>IP адреса:</strong>
              <p>Список IP адресов через запятую (<code>192.168.1.1</code>, <code>10.0.0.5</code>)</p>
              <p class="help-warning"><i class="fas fa-exclamation-triangle"></i> Используется реальный IP соединения (не заголовки прокси!)</p>
            </div>
            <div class="help-param">
              <strong>Исключения:</strong>
              <p>Пути-исключения через запятую (<code>/bot/*</code>, <code>/public/*</code>). Правило НЕ применяется к этим путям</p>
            </div>
            <div class="help-param">
              <strong>Страница ошибки:</strong>
              <p>Куда перенаправить при блокировке:</p>
              <ul>
                <li><code>404</code> — стандартная страница ошибки</li>
                <li><code>https://site.com</code> — внешний редирект</li>
                <li><code>/error.html</code> — локальная страница</li>
              </ul>
            </div>
          </div>
        </div>

        <div class="help-card">
          <h3><i class="fas fa-search"></i> {{ t('vaccess.helpPatterns') }}</h3>
          <div class="help-patterns">
            <div class="pattern-item"><code>*.ext</code><span>Любой файл с расширением .ext</span></div>
            <div class="pattern-item"><code>no_extension</code><span>Файлы без расширения</span></div>
            <div class="pattern-item"><code>/path/*</code><span>Все файлы в папке /path/ и подпапках</span></div>
            <div class="pattern-item"><code>/file.php</code><span>Конкретный файл</span></div>
          </div>
        </div>

        <div class="help-card help-examples">
          <h3><i class="fas fa-lightbulb"></i> {{ t('vaccess.helpExamples') }}</h3>
          <div class="help-example">
            <h4>1. Запретить выполнение PHP в uploads</h4>
            <div class="example-rule">
              <div><strong>Тип:</strong> <VBadge variant="no">Disable</VBadge></div>
              <div><strong>Расширения:</strong> <code>*.php</code></div>
              <div><strong>Пути:</strong> <code>/uploads/*</code></div>
              <div><strong>Ошибка:</strong> <code>404</code></div>
            </div>
          </div>
          <div class="help-example">
            <h4>2. Разрешить админку только с определённых IP</h4>
            <div class="example-rule">
              <div><strong>Тип:</strong> <VBadge variant="yes">Allow</VBadge></div>
              <div><strong>Пути:</strong> <code>/admin/*</code></div>
              <div><strong>IP:</strong> <code>192.168.1.100, 127.0.0.1</code></div>
              <div><strong>Ошибка:</strong> <code>404</code></div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style>
.vaccess-tab {
  flex: 0 0 auto;
  padding: 10px 18px;
  background: transparent;
  border: none;
  border-radius: var(--radius);
  color: var(--text-muted);
  font-size: var(--text-base);
  font-weight: var(--font-medium);
  cursor: pointer;
  transition: all var(--transition-base);
  display: flex;
  align-items: center;
  gap: var(--space-sm);
}

.vaccess-tab:hover {
  background: rgba(var(--accent-rgb), 0.1);
  color: var(--text-primary);
}

.vaccess-tab.active {
  background: var(--accent-purple);
  color: white;
}

.vaccess-tab-content {
  animation: fadeIn var(--transition-slow);
}

/* Таблица правил */
.vaccess-rules-container {
  width: 100%;
  overflow-x: auto;
}

.vaccess-table {
  width: 100%;
  border-collapse: collapse;
}

.vaccess-table thead tr {
  background: rgba(var(--accent-rgb), 0.05);
}

.vaccess-table th {
  padding: 16px;
  text-align: left;
  font-size: var(--text-base);
  font-weight: var(--font-semibold);
  color: var(--text-muted);
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.vaccess-table tbody tr {
  background: var(--subtle-overlay);
  transition: all var(--transition-slow);
}

.vaccess-table tbody tr:hover {
  background: rgba(var(--accent-rgb), 0.08);
}

.vaccess-table td {
  padding: 20px 16px;
  font-size: var(--text-md);
  color: var(--text-primary);
  border-top: 1px solid rgba(var(--accent-rgb), 0.05);
  border-bottom: 1px solid rgba(var(--accent-rgb), 0.05);
}

.th-with-info {
  display: inline-flex;
  align-items: center;
  gap: 6px;
}

.col-drag { width: 3%; min-width: 40px; text-align: center; }
.col-type { width: 8%; min-width: 80px; }
.col-files { width: 15%; min-width: 120px; }
.col-paths { width: 18%; min-width: 150px; }
.col-ips { width: 15%; min-width: 120px; }
.col-exceptions { width: 15%; min-width: 120px; }
.col-error { width: 10%; min-width: 100px; }
.col-actions { width: 5%; min-width: 60px; text-align: center; }

.drag-handle {
  color: var(--text-muted);
  opacity: 0.3;
  cursor: grab;
  text-align: center;
  transition: all var(--transition-base);
}

.drag-handle:hover {
  opacity: 1;
  color: var(--accent-purple-light);
}

.vaccess-table tbody tr.dragging {
  opacity: 0.4;
}

.vaccess-table tbody tr.drag-over {
  border-top: 2px solid var(--accent-purple);
}

.vaccess-table tbody tr.drag-over td {
  border-top: 2px solid var(--accent-purple);
}

.mini-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
}

.mini-tags code {
  padding: 2px 6px;
  background: rgba(var(--accent-rgb), 0.15);
  border-radius: var(--radius);
  font-size: var(--text-sm);
  color: var(--accent-purple-light);
}

.empty-field {
  color: var(--text-muted);
  opacity: 0.4;
  font-style: italic;
}

.type-toggle {
  cursor: pointer;
  transition: all var(--transition-fast);
}

.type-toggle:hover {
  opacity: 0.7;
}

.inline-input {
  padding: 4px 8px;
  background: var(--glass-bg-dark);
  border: 1px solid var(--glass-border);
  border-radius: var(--radius);
  color: var(--text-primary);
  font-size: 12px;
  font-family: var(--font-mono);
  outline: none;
  width: 80px;
}

.inline-input:focus {
  border-color: rgba(var(--accent-rgb), 0.5);
}

.inline-input::placeholder {
  color: var(--text-muted);
  opacity: 0.4;
}

.col-actions-cell {
  text-align: center;
}

.icon-btn-small {
  width: 28px;
  height: 28px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  background: rgba(var(--danger-rgb), 0.1);
  border: 1px solid rgba(var(--danger-rgb), 0.3);
  border-radius: var(--radius);
  color: var(--accent-red);
  font-size: 12px;
  cursor: pointer;
  transition: all var(--transition-base);
}

.icon-btn-small:hover {
  background: rgba(var(--danger-rgb), 0.2);
}

/* Пустое состояние */
.vaccess-empty {
  text-align: center;
  padding: 80px 40px;
  color: var(--text-muted);
}

.empty-icon {
  font-size: 64px;
  margin-bottom: var(--space-lg);
  opacity: 0.3;
}

.vaccess-empty h3 {
  font-size: var(--text-2xl);
  font-weight: var(--font-semibold);
  color: var(--text-primary);
  margin-bottom: 12px;
}

.vaccess-empty p {
  font-size: var(--text-md);
  margin-bottom: var(--space-lg);
}

/* Help секция */
.vaccess-help {
  display: flex;
  flex-direction: column;
  gap: var(--space-lg);
}

.help-card {
  background: var(--subtle-overlay);
  border-radius: var(--radius);
  padding: var(--space-xl);
  border: 1px solid var(--glass-border);
  transition: all var(--transition-slow);
}

.help-card:hover {
  border-color: var(--glass-border-hover);
}

.help-card h3 {
  font-size: var(--text-2xl);
  font-weight: var(--font-semibold);
  color: var(--accent-purple-light);
  margin: 0 0 20px 0;
  display: flex;
  align-items: center;
  gap: var(--space-lg);
}

.help-card ul {
  list-style: none;
  padding: 0;
  margin: 0;
}

.help-card li {
  padding: 12px 0;
  color: var(--text-secondary);
  line-height: 1.6;
  border-bottom: 1px solid rgba(var(--accent-rgb), 0.05);
}

.help-card li:last-child {
  border-bottom: none;
}

.help-params {
  display: flex;
  flex-direction: column;
  gap: var(--space-lg);
}

.help-param {
  padding: 20px;
  background: rgba(var(--accent-rgb), 0.03);
  border-radius: var(--radius);
  border-left: 3px solid var(--accent-purple);
}

.help-param strong {
  display: block;
  font-size: 15px;
  color: var(--accent-purple-light);
  margin-bottom: var(--space-sm);
}

.help-param p {
  margin: var(--space-sm) 0 0 0;
  color: var(--text-secondary);
  line-height: 1.6;
}

.help-param code {
  padding: 3px 8px;
  background: rgba(var(--accent-rgb), 0.15);
  border-radius: var(--radius);
  font-size: var(--text-base);
  color: var(--accent-purple-light);
}

.help-param ul {
  margin: 12px 0 0 20px;
  list-style: disc;
}

.help-warning {
  color: var(--warning-icon-color) !important;
  margin-top: var(--space-sm);
  font-size: var(--text-base);
  display: flex;
  align-items: center;
  gap: var(--space-sm);
}

.help-patterns {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: var(--space-md);
}

.pattern-item {
  display: flex;
  flex-direction: column;
  gap: var(--space-sm);
  padding: var(--space-md);
  background: rgba(var(--accent-rgb), 0.05);
  border-radius: 10px;
  transition: all var(--transition-base);
}

.pattern-item:hover {
  background: rgba(var(--accent-rgb), 0.1);
}

.pattern-item code {
  font-size: var(--text-md);
  font-weight: var(--font-semibold);
  color: var(--accent-purple-light);
}

.pattern-item span {
  font-size: var(--text-base);
  color: var(--text-muted);
}

.help-examples {
  background: linear-gradient(135deg, rgba(var(--accent-rgb), 0.05) 0%, rgba(var(--accent-rgb), 0.03) 100%);
}

.help-example {
  margin-bottom: var(--space-xl);
}

.help-example:last-child {
  margin-bottom: 0;
}

.help-example h4 {
  font-size: var(--text-lg);
  font-weight: var(--font-semibold);
  color: var(--text-primary);
  margin-bottom: var(--space-md);
}

.example-rule {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 12px;
  padding: 20px;
  background: var(--subtle-overlay);
  border-radius: 10px;
  border: 1px solid var(--glass-border);
}

.example-rule div {
  font-size: var(--text-md);
  color: var(--text-secondary);
}

.example-rule strong {
  color: var(--text-muted);
  margin-right: var(--space-sm);
}

.example-rule code {
  color: var(--accent-purple-light);
}
</style>
