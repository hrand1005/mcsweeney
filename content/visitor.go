package content

type Visitor interface {
	VisitClip(*Clip) error
	VisitIntro(*Intro) error
	VisitOutro(*Outro) error
}
