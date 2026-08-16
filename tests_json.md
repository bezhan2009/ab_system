### **B2. Feature Flags и выдача вариантов**

#### **B2-1: default при отсутствии активного эксперимента**
```json
// POST /decisions/decide
{
  "user_id": "user123",
  "flags": ["button_color"]
}
```

#### **B2-2: default, если пользователь не подходит под таргетинг**
Предположим, эксперимент настроен на страну `RU`. Пользователь из Казахстана:
```json
POST /decisions/decide
{
  "user_id": "user_kz",
  "attributes": { "country": "KZ" },
  "flags": ["button_color"]
}
// Ответ: значение по умолчанию "green" (как в примере выше)
```

#### **B2-3: variant, если пользователь подходит**
Тот же эксперимент, пользователь из России:
```json
POST /decisions/decide
{
  "user_id": "user_ru",
  "attributes": { "country": "RU" },
  "flags": ["button_color"]
}
```

#### **B2-4: детерминизм**
Повторяем тот же запрос с `user_ru` – получаем тот же вариант "red". (Можно показать, что decision_id меняется, но значение флага то же.)

#### **B2-5: веса вариантов**
Для проверки весов можно выполнить 100 запросов с разными user_id и потом посмотреть распределение. Но для демо достаточно показать, что у эксперимента два варианта с весами 50/50 и что в ответе приходит то control, то test.

---

### **B3. Эксперименты – жизненный цикл и ревью**

#### **B3-1: черновик → на ревью**
```json
// PATCH /experiments/exp_123/review
{}
// Ответ:
{
  "message": "successfully sent to review"
}
// После этого статус эксперимента становится "in_review".
```

#### **B3-2: in_review → approved (после одобрений)**
```json
// POST /approvals
{
  "experiment_id": "exp_123",
  "approved": true,
  "comment": "OK, можно запускать"
}
// Когда наберётся нужное число одобрений, статус эксперимента станет "approved".
```

#### **B3-3: блокировка запуска без достаточных одобрений**
```json
// PATCH /experiments/exp_123/run
{}
// Ожидаемый ответ 403 или 400 с текстом "experiment not approved"
```

#### **B3-4: запрет недопустимых переходов**
Попытка перевести запущенный эксперимент обратно в черновик:
```json
PUT /experiments/exp_123
{
  "status": "draft"
}
// Ответ 400 с сообщением о недопустимом переходе.
```

#### **B3-5: только назначенные approver'ы могут одобрять**
Попытка одобрить эксперимент от пользователя, не входящего в группу аппруверов для данного экспериментера:
```json
POST /approvals
{
  "experiment_id": "exp_123",
  "approved": true
}
// Ответ 403 Forbidden
```

---

### **B4. События и атрибуция**

#### **B4-1, B4-2: валидация типов и обязательных полей**
```json
// POST /events/batch
[
  {
    "event_id": "evt_001",
    "event_type": "purchase",
    "decision_id": "dec_123",
    "user_id": "user_ru",
    "payload": {
      "amount": "сто",    // должно быть числом, а это строка
      "currency": "RUB"
    }
  }
]
```

#### **B4-3: дедупликация**
```json
// POST /events/batch
[
  {
    "event_id": "dup_001",
    "event_type": "impression",
    "decision_id": "dec_123",
    "user_id": "user_ru",
    "payload": {}
  },
  {
    "event_id": "dup_001",   // тот же event_id
    "event_type": "impression",
    "decision_id": "dec_123",
    "user_id": "user_ru",
    "payload": {}
  }
]
// Ответ:
{
  "accepted": 1,
  "duplicate": 1,
  "rejected": 0,
  "errors": []
}
```

#### **B4-4, B4-5: атрибуция (связь с decision_id)**
Это будет видно в отчёте. Например, в отчёте по эксперименту для каждого варианта считаются конверсии только при наличии события `impression` с тем же decision_id.

---

### **B5. Guardrails**

#### **B5-1 – B5-4: настройка, срабатывание и действие**
Сначала привязываем к эксперименту guardrail-метрику:
```json
// POST /experiments/exp_123/metrics
{
  "metric_id": "error_rate",
  "is_guardrail": true,
  "threshold": 5.0,
  "operator": ">",
  "window_min": 5,
  "action": "pause"
}
// Ответ 201
```

Затем отправляем много событий с ошибками (через `/events/batch`). Через некоторое время проверяем статус эксперимента:
```json
GET /experiments/exp_123
// Ответ:
{
  "id": "exp_123",
  "status": "paused",
  "guardrail_triggered": true,
  ...
}
```

#### **B5-5: аудит срабатывания**
```json
GET /experiments/exp_123/guardrails
// Ответ:
[
  {
    "id": "trigger_001",
    "metric_id": "error_rate",
    "threshold": 5.0,
    "operator": ">",
    "window_min": 5,
    "actual_value": 6.2,
    "action": "pause",
    "triggered_at": "2026-03-07T12:05:00Z"
  }
]
```

#### **B5-6: ограничение участия пользователя** (если реализовано)
Можно показать, что один пользователь не может одновременно участвовать в двух конфликтующих экспериментах (если есть конфликт доменов).

---

### **B6. Отчётность**

#### **B6-1, B6-2, B6-3: отчёт за период с разбивкой по вариантам и метрикам**
```json
GET /reports/exp_123?from=2026-03-01T00:00:00Z&to=2026-03-07T23:59:59Z
// Ответ:
{
  "experiment_id": "exp_123",
  "from": "2026-03-01T00:00:00Z",
  "to": "2026-03-07T23:59:59Z",
  "variants": [
    {
      "variant_id": "var_789",
      "variant_name": "test",
      "metric_values": [
        {
          "metric_id": "error_rate",
          "metric_title": "Доля ошибок",
          "value": 2.34,
          "unit": "%"
        },
        {
          "metric_id": "conversion_rate",
          "metric_title": "Конверсия",
          "value": 5.67,
          "unit": "%"
        }
      ]
    },
    {
      "variant_id": "var_456",
      "variant_name": "control",
      "metric_values": [
        {
          "metric_id": "error_rate",
          "value": 2.45,
          "unit": "%"
        },
        {
          "metric_id": "conversion_rate",
          "value": 4.89,
          "unit": "%"
        }
      ]
    }
  ]
}
```

#### **B6-4, B6-5: фиксация исхода с комментарием**
```json
// POST /experiments/exp_123/complete
{
  "conclusion": "winner",
  "comment": "Красная кнопка дала +15% к конверсии",
  "winner_variant_id": "var_789"
}
// Ответ:
{
  "id": "exp_123",
  "status": "completed",
  "conclusion": "winner",
  "comment": "Красная кнопка дала +15% к конверсии",
  "winner_variant_id": "var_789",
  ...
}
```

---

### **Edge‑case сценарии**

#### **1. События приходят с задержкой и не по порядку**
Можно показать, что при построении отчёта используется `client_time`, а не `received_at`. Например, отправить событие `click` с `client_time` более ранним, чем `impression`, и убедиться, что в отчёте порядок правильный. Для демонстрации достаточно привести пример запроса событий и затем отчёт.

#### **2. Эксперимент улучшает целевую метрику, но портит стабильность – guardrail срабатывает**
Это уже показано в B5. Можно дополнительно подчеркнуть, что даже при росте конверсии эксперимент останавливается из‑за ошибок.

#### **3. Конфликты экспериментов (если реализовано)**
Создать два эксперимента на одном конфликтном домене (например, оба на кнопку). Для одного пользователя проверить, что он получает вариант только одного из них, а второй игнорируется. Запрос:
```json
POST /decisions/decide
{
  "user_id": "user_conflict",
  "flags": ["button_color", "button_text"]
}
// Ответ должен содержать значение только одного эксперимента (например, цвет кнопки, а текст – default).
```

#### **4. Один и тот же пользователь слишком часто попадает в эксперименты**
Если реализовано ограничение, можно показать, что после нескольких экспериментов пользователь перестаёт попадать в новые тесты и получает default.

---

### **B9 – наблюдаемость (не JSON, а примеры логов и метрик)**

```bash
# Пример структурированного лога
{"time":"2026-03-07T12:10:00Z","level":"info","trace_id":"abc123","op":"service.Decide","msg":"decision made","user_id":"user_ru","flag":"button_color","value":"red"}

# Пример метрики Prometheus (эндпоинт /metrics)
http_requests_total{method="POST",path="/api/v1/decisions/decide",status="200"} 42
```

---

Этих примеров достаточно, чтобы продемонстрировать все ключевые сценарии и ответы системы.