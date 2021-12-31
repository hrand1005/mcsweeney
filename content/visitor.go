package content

type Visitor interface {
	VisitClip(*Clip)
	VisitIntro(*Intro)
	VisitOutro(*Outro)
}
