package app

type App struct { // TODO
}

type Logger interface { // TODO
}

type Storage interface { // TODO
}

func New() *App {
	return &App{}
}

func (a *App) CreateEvent() error {
	// TODO
	return nil
	// return a.storage.CreateEvent(storage.Event{ID: id, Title: title})
}

// TODO
