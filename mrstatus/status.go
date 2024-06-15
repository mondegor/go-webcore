package mrstatus

type (
	// Flow - интерфейс для управления статусами их переключениями.
	Flow interface {
		Exists(status Getter) bool
		Check(from, to Getter) bool
		PossibleToStatuses(from Getter) []Getter
		PossibleFromStatuses(to Getter) []Getter
	}

	// Getter - статус участвующий во Flow.
	Getter interface {
		String() string
	}
)
