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
- Пользователь - авторизованный пользователь с доступом к отправке отчетов на проверку.
- Нормоконтроллер - авторизованный пользователь с доступом к инструментам разметки и инструментам администрирования.
- Администратор - пользователь способный выдавать роли сущностям в системе, имеющий привелегии в добавлении и удалении проверяемых ошибок из датасета.

## Use-Case - диаграмма

