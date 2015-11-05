# goat
Generate SOAP requests for Go at runtime

# Approach
In my opinion, there are two main approaches for writing a soap client for Go
(any language):

1. generate native code from the WSDL withouth reflection
2. parse the WSDL at runtime and offer a general interface which then generates the requests and parses the responses

This library tries to combine the best from both approaches:

1. load WSDL at runtime, generate the requests from a parameter set and
2. generate structs to parse the responses into and work with them in a native way

# Motivation
REST is on the way, and all the services at JustWatch are REST APIs. But we need
to talk to the Adwords API, which is a SOAP api. Since there're no plans to change
this, I started changing the way how SOAP is handled in Go. The status quo is,
that XML in general is not as nice as handling JSON, but it is still possible.

Then there are some libraries around, which auto-generate code from a wsdl.
I prefer having native go-structs over using map[string]interface{}. Even if it
would compile all the time, putting the structs together still is a hazzle
because of the different name spaces and so on. I don't want to care about all
possible fields and their name spaces and if I have as many elements in a
sequence as there should be ("minOccurs" and "maxOccurs"). I just want to pass
the parameter fields and a library cares for putting them in the right order and
validating them.

Back to the Adwords API: there are some nice-ish libraries around, I only had a
look at the Python library: There it is like it is meant to be: name the service,
it gets downloaded and parsed by [suds](https://pypi.python.org/pypi/suds) and
the you just pass a dictionary with fields. Not so much code, and the library
(suds) knows what to do.

As long there is reflection in Go, this should be possible.

# Solution
This Library is forked from [wsdl-go](https://code.google.com/p/wsdl-go/). In the
future, it should take out a lot of stuff from the auto-generating stuff to only
have the structs. The wsdl and the xsd packages will take care about generating
the requests.

# Status
[x] support for generating basic requests
[x] some Adwords API Endpoints still work (for get Requests)
[ ] attributes
[ ] validation ("minOccurs" and "maxOccurs")
[ ] boil down code generation stuff
[ ] retrieving of xsd schemes not already in the WSDL
[ ] make the already working parts *nice* and *tested*
