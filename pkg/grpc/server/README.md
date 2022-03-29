# b32svc

This is an example of a small concurrent handler for the back end of a
microservice or similar.

Because the API of this example only has two functions, they are manually
specified and the handlers are defined inline, but for a larger API, it would
usually be done with the handler implementations separately defined, and then a
generator to string the API specifications together and create an API handler
that contains all the elements in this implementation.

The service defined here is written so as to be a library that can be pulled 
in by a separate concrete launch system, such as, for this case, usually 
would be a tty based service that runs under the control of systemd or similar.
