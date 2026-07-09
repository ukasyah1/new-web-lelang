package faq

type FAQ struct {
	ID       string
	Question string
	Answer   string
}

type Category struct {
	ID   string
	Name string
	FAQs []FAQ
}
