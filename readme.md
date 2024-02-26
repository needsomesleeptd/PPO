## Название проекта
Создание системы проверки соответствия документов выделенным требованиям (например ГОСТ).

## Краткое описание идеи проекта
Выделяются 2 роли: контроллер и пользователь. Контроллер размечает распространенные ошибки пользователей в соответствии с определенными стандартами (например ГОСТ 7.32), после чего на размеченных данных обучается нейросеть, решающая задачу детекции размеченных ошибок. Нейросетью рассматривается только визуальные части документов (таблицы, графики, схемы алгоритмов, формулы), датасет ошибок пользователей может дополняться во время работы системы.
Пользователь может предложить свой документ на проверку, и получить предсказания ошибок в документе от обученной нейросети.

**Проблема: большое количество времени на ручную проверку соответствия документов стандарту (например курсача).**

## Краткое описание предметной области
Документы, создаваемые и оформляемые в соответствии с определенными установленными стандартами форматирования, структурирования и содержания.
Такие стандарты могут включать в себя требования к шрифтам, отступам, заголовкам, нумерации страниц, использованию логотипов и другие аспекты визуальной и текстовой атрибутики.

## Краткий анализ аналогичных решений по минимум 3 критериям
Решений решающих конкретно выделенную проблему небыло найдено, однако были рассмотрены решения для проверки работ на соответствие ГОСТ 7.32 (ВКР ВУЗ, TestVkr) и приложение для совершения "визуального" тестирования ().
Решение | Проверка элементов отчета | Наличие общего хранилища работ| Возможность динамического добавления правил
----------------------------------------------------------- | -------------------- | ----------------- | ----------------------
[ВКР ВУЗ](http://www.vkr-vuz.ru/)         | -                    | -                 | -     |
[TestVkr](https://labelme.ru/)                              | -                    | -                 | -                      |
[Applitools](https://applitools.com/)                         | +                    | -                 | +                    |
Предлагаемое решение                                        |+                    | +                 | +  |

## Краткое обоснование целесообразности и актуальности проекта

* Эффективность и точность: ручная проверка с может быть трудоемкой и подверженной человеческой ошибке. 

* Экономия времени и ресурсов: Автоматизация процесса проверки соответствия документов позволит сократить временные и человеческие затраты, ускоряя процессы разработки, производства и контроля качества.

## Краткое описание акторов
- Гость - неавторизованный посетитель.
- Авторизованный пользователь - пользователь, прошудший аутенфикацию, имеющий доступ к отправке документа на проверку и возможностью комментирования результатов проверки.
- Нормоконтроллер - авторизованный пользователь с доступом к инструментам разметки, добавления ошибок.
- Администратор - пользователь способный управлять определением ролей в системе, имеющий привелегии в добавлении и удалении проверяемых ошибок из датасета.

## Use-Case - диаграмма
![Диаграмма использования приложения](imgs/PPO_use_case.svg)

## ER-диаграмма сущностей
![Диаграмма использования приложения](imgs/PPO_ER.svg)

## Пользовательские сценарии
1. Сценарий добавления разметки на не существующий до этого тип ошибки:
   1. пользователь авторизируется как нормконтроллер;
   2. пользователь загружает документы для разметки;
   3. пользователь создает новый тип ошибки;
   4. пользователь выбирает неразмеченную страницу;
   5. пользователь выбирает требуемый тип ошибки;
   6. пользователь размечает данные нужным типом ошибки;
   7. пользователь подтверждает создание разметки.
2. Сценарий отправления документа на проверку:
   1. пользователь заходит в систему;
   2. пользователь загружает документ в систему;
   3. пользователь получает список возможных ошибок.
3. Сценарий изменения роли пользователя user администратором:
   1. пользователь авторизируется как администратор;
   2. пользователь выбирает требуемого пользователя (user) для изменения роли;
   3. пользователь назначает user необходимую роль;
   4. пользователь выбирает сохранить/удалить ли данные внесенные user (либо разметка и типы ошибок, либо документы);
   5. пользователь подтверждает изменение роли user.
4. Сценарий удаления разметки нормоконтроллером
   1. пользователь региструрется как нормконтроллер;
   2. пользователь получает список документов с данной ошибок;
   3. пользователь просматривает размеченные документы
   4. пользователь выбирает разметки для удаления по страницу документа + номер разметки;
   5. пользователь подтверждает удаление.

## Формализация ключевых бизнес-процессов
![Диаграмма аутенфикации](imgs/PPO_reg.svg)
![Диаграмма разметки](imgs/PPO_mark.svg)