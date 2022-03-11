package b32svc

type Service struct {
	ServiceDefinition
	incoming chan interface{}
	quit     chan struct{}
}

type Handler struct {
	Chan interface{}
}

type Handlers map[string]Handler

type ServiceDefinition struct {
	Name string
	Handlers
}

func New(svc ServiceDefinition) *Service {

	return &Service{ServiceDefinition: svc}
}

func (s *Service) Run() func() {
	s.quit = make(chan struct{})
	go func() {
	out:
		for {
			select {
			case msg := <-s.incoming:
				switch msg := msg.(type) {
				case int:
					_ = msg
				default:
				}
			case <-s.quit:
				break out
			}
		}
		log.Println(s.ServiceDefinition.Name, "service runner is now shut down")
	}()
	return func() { close(s.quit) }
}
