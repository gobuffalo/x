package mailer

//Deliverer interface for any upcomming mailers.
type Deliverer interface {
	Deliver(Message) error
}
