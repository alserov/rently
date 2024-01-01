package app

type App struct {
	port int
}

func NewApp() *App {
	return &App{}
}

func (a *App) MustStart() {

}
