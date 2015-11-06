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

## Personal motivation
Back in university, we learned M68k Assembler (which I'm very happy about).
There was a Windows program we used to execute the asm code with, but I wasn't happy
about it. I wasn't happy about the UX and that we had to use Windows (I was a software-hippie
back then ;-).

So then I tried to create my own interpreter written in Java using GTK bindings.
At the end, it helped me get to know Java, its standard library and Assembler.
And since it was in a very early semester, it helped me a lot by solving coding tasks
faster than the majority.

But on the other hand, I didn't have something to replace that Windows program with.

This project should also be something like the described one:

- learn XML and XSD / WSDL (and SOAP)
- make it easier to handle WSDLs in Go (for me)

So even if this project does not work out at the end, it may help the one or the other
to handle SOAP in Go a little bit more efficient.

# Solution
This Library is forked from [wsdl-go](https://code.google.com/p/wsdl-go/). In the
future, it should take out a lot of stuff from the auto-generating stuff to only
have the structs. The wsdl and the xsd packages will take care about generating
the requests.

# Status

- [x] support for generating basic requests
- [x] some Adwords API Endpoints still work (for get Requests)
- [ ] attributes
- [ ] validation ("minOccurs" and "maxOccurs")
- [ ] boil down code generation stuff
- [ ] retrieving of xsd schemes not already in the WSDL
- [ ] make the already working parts *nice* and *tested*
- [ ] use structs with proper xml tags for parameters, not map[string]interface{} (for simpler use of attributes)

# Example

```go
    c := makeNewOAuthHTTPClient()

    ws = goat.NewWebservice(c, map[string]interface{}{
        "RequestHeader/clientCustomerId": "CLIENT_CUSTOMER_ID",
        "RequestHeader/developerToken":   "DEVELOPER_TOKEN",
        "RequestHeader/userAgent":        "a random header",
        "RequestHeader/validateOnly":     true,
        "RequestHeader/partialFailure":   false,
    })

    err := ws.AddServices("https://adwords.google.com/api/adwords/mcm/v201509/ManagedCustomerService?wsdl")
    if err != nil {
        panic(err)
    }

    resp := struct {
        XMLName     xml.Name `xml:"getResponse"`
        GetResponse struct {
            XMLName         xml.Name `xml:"rval"`
            TotalNumEntries int      `xml:"totalNumEntries"`
            PageType        string   `xml:"Page.Type"`
            Entries         []struct {
                XMLName          xml.Name `xml:"entries"`
                AccountLabels    string   `xml:"accountLabels"`
                CanManageClients bool     `xml:"canManageClients"`
                CompanyName      string   `xml:"companyName"`
                CurrencyCode     string   `xml:"currencyCode"`
                CustomerId       string   `xml:"customerId"`
                DateTimeZone     string   `xml:"dateTimeZone"`
                Name             string   `xml:"name"`
                TestAccount      bool     `xml:"testAccount"`
            } `xml:"entries"`
        }
    }{}
    err = ws.Do("ManagedCustomerService", "get", &resp, map[string]interface{}{
        "get/serviceSelector/fields": []string{
            "AccountLabels",
            "CanManageClients",
            "CompanyName",
            "CurrencyCode",
            "CustomerId",
            "DateTimeZone",
            "Name",
            "TestAccount",
        },
    })
    if err != nil {
        panic(err)
    }
    // work with resp
```
