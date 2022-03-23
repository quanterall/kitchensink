# cmd

In the cmd folder of a Go project, you put executables that can be invoked 
from the terminal, or potentially, from a GUI, if the app has a GUI front 
end available.

For standard command line interface apps, it is a good policy to use a 
convention that is widely used, if the executable is a service, add a `d` to 
the end of the name, and for clients, `cli` and if the application controls 
a service, use the suffix `ctl`.
