# Описание GoWebCore v0.4.0
Этот репозиторий содержит описание библиотеки GoWebCore.

## Статус библиотеки
Библиотека находится в стадии разработки.

## Описание библиотеки
Библиотека с базовой функциональностью для разработки web сервисов, в которую входят:
- общие интерфейсы, такие как logger, router, validator и другие, которые реализовываются уже в конкретных проектах;
- часто используемые программные, системные и пользовательские ошибки, которые возникают в разных слоях программы;
- пакеты с часто используемыми функциями: генерация токенов, преобразование IP и т.д.;
- парсеры для некоторых типов данных, которые поступают из вне;
- сохранение и извлечения значений некоторых типов данных в контексте;
 
## Подключение библиотеки
go get github.com/mondegor/go-webcore