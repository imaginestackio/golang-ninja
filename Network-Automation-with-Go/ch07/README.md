# Network APIs

As the ways we build, deploy, and operate networks evolve, new protocols and interfaces are emerging to ease machine-to-machine communication—a primary enabler of network automation. In this and the following chapters, we’ll navigate through some of these new capabilities and explore how to take advantage of them in the context of the Go programming language.

The network **Command-Line Interface** (**CLI**) is what we, network engineers, have used for decades to operate and manage network devices. As we move toward a more programmatic approach to managing networks, simply relying on faster CLI command execution might not be enough to deploy network automation solutions at scale.

Solutions that don’t have a strong foundation are brittle and unstable. Hence, when possible, we prefer to build network automation projects based on structured data and machine-friendly **Application Programming Interfaces** (**APIs**). The target use case for these interfaces isn’t direct human interaction, so you can rely on Go to translate between remote API calls and a local, user-facing interface.

When we talk about APIs, we generally refer to different things that make up the API developer experience, which you need to consider when evaluating an API:

-   A set of **Remote Procedure Calls** (**RPCs**) defining the rules of interaction between a client and a server—at the very least, this would include a standard set of create, get, update, and delete operations.
-   The structure and data type exchanged—product vendors can define this using data model specification languages such as YANG or OpenAPI.
-   The underlying protocol that wraps the modeled data, which you can serialize into one of the standard formats, such as XML or JSON, and transports it between a client and a server—this could be SSH or, more often these days, HTTP.

In the networking world, we have another dimension in the API landscape that determines the origin of a model specification document. While every networking vendor is free to write their own data models, there are two sources of vendor-agnostic models—IETF and OpenConfig—that strive to offer a vendor-neutral way of configuring and monitoring network devices. Because of this variability in the API ecosystem, it’s impossible to cover all protocols and standards, so in this chapter, we’ll only cover a subset of network APIs, selected based on availability, practicality, and usefulness:

-   We’ll start by looking at OpenAPI as one of the most prevalent API specification standards in a wider infrastructure landscape.
-   We’ll then move on to JSON-RPC, which uses vendor-specific YANG models.
-   After that, we’ll show an example of an RFC-standard HTTP-based protocol called RESTCONF.
-   Finally, we’ll look at how you can leverage **Protocol Buffers** (**protobuf**) and gRPC to interact with network devices and stream telemetry.

In this chapter, we’ll focus only on these network APIs, as the others are outside of the scope. The most notable absentee is the **Network Configuration Protocol** (**NETCONF**)—one of the oldest network APIs, defined originally by IETF in 2006. We’re skipping NETCONF mainly because of the lack of support for XML in some Go packages we use throughout this chapter. Although NETCONF is in use today and offers relevant capabilities, such as different configuration datastores, configuration validation, and network-wide configuration transactions, in the future, it may get displaced by technologies running over HTTP and TLS, such as RESTCONF, gNMI, and various proprietary network APIs.

Just Imagine

# Technical requirements

You can find the code examples for this chapter in the book’s GitHub repository (refer to the _Further reading_ section), under the `ch08` folder.

Important Note

We recommend you execute the Go programs in this chapter in a virtual lab environment. Refer to the appendix for prerequisites and instructions on how to build it.

Just Imagine

# API data modeling

Before we look at any code, let’s review what data modeling is, what its key components are, and their relationships. While we focus on the configuration management side of model-driven APIs for this explanation, similar rules and assumptions apply to workflows involving state data retrieval and verification.

The main goal of a configuration management workflow is to transform some input into a serialized data payload whose structure adheres to a data model. This input is usually some user-facing data, which has its own structure and may contain only a small subset of the total number of configuration values. But this input has a one-to-one relationship with the resulting configuration, meaning that rerunning the same workflow should result in the same set of RPCs with the same payloads and the same configuration state on a network device.

At the center of it all is a data model—a text document that describes the hierarchical structure and types of values of a (configuration) data payload. This document becomes a contract with all potential clients—as long as they send their data in the right format, a server should be able to understand it and parse it. This contract works both ways so that when a client requests some information from a server, it can expect to receive it in a predetermined format.

The following diagram shows the main components of a model-driven configuration management workflow and their relationships:

![Figure 8.1 – Data modeling concepts](https://static.packt-cdn.com/products/9781800560925/graphics/image/B16971_08_01.jpg)

Figure 8.1 – Data modeling concepts

Thus far, we’ve discussed a model, its input, and the resulting configuration. The only thing we haven’t mentioned until now is the _bindings_. We use this term to refer to a broad set of tools and libraries that can help us generate the final configuration data payload programmatically, that is, without resorting to a set of text templates or building these data payloads manually, both of which we consider an anti-pattern in any network automation workflow. We produce these bindings based on the data model and they represent a programmatic view of the model. They may also include several helper functions to serialize and deserialize data structures into one of the expected output formats, for example, JSON or protobuf. We’ll spend most of this chapter discussing and interacting with bindings as they become the main interface for a data model inside of the programming language.

Now that we’ve covered some theory, it’s time to put it into practice. In the following section, we’ll examine OpenAPI models and one way you can instantiate and validate them.

Just Imagine

# OpenAPI

Within a greater infrastructure landscape, HTTP and JSON are two commonly used standards for machine-to-machine communication. Most web-based services, including public and private clouds, use a combination of these technologies to expose their externally facing APIs.

The OpenAPI Specification allows us to define and consume RESTful APIs. It lets us describe the enabled HTTP paths, responses, and JSON schemas for the corresponding payloads. It serves as a contract between an API provider and its clients to allow for a more stable and reliable API consumer experience and enables API evolution through versioning.

We don’t widely use OpenAPI in networking, arguably for historical reasons. YANG and its ecosystem of protocols predate OpenAPI and the rate of change in network operating systems is not as fast as you might expect. But we often find OpenAPI support in network appliances—SDN controllers, monitoring and provisioning systems or **Domain Name System** (**DNS**), **Dynamic Host Configuration Protocol** (**DHCP**), and **IP Address Management** (**IPAM**) products. This makes working with OpenAPI a valuable skill to have for any network automation engineer.

In _Chapters 6_ and _7_, we went through an example of how to interact with Nautobot’s external OpenAPI-based interface. We used a Go package produced by an open source code generation framework based on Nautobot’s OpenAPI specification. One thing to be mindful of with automatic code generation tools is that they rely on a certain version of the OpenAPI Specification. If the version of your API specification is different (there are nine different OpenAPI versions today; refer to the _Further reading_ section), the tool may not generate the Go code. Hence, we want to explore an alternative approach.

In this section, we’ll configure NVIDIA’s Cumulus Linux device (`cvx`), which has an OpenAPI-based HTTP API, using **Configure Unify Execute** (**CUE**; refer to the _Further reading_ section)—an open source **Domain-Specific Language** (**DSL**) designed to define, generate, and validate structured data.

CUE’s primary user-facing interface is CLI, but it also has first-class Go API support, so we’ll focus on how to interact with it entirely within Go code while providing the corresponding shell commands where appropriate.

The following figure shows a high-level overview of the Go program we’ll discuss next:

![Figure 8.2 – Working with OpenAPI data models](https://static.packt-cdn.com/products/9781800560925/graphics/image/B16971_08_02.jpg)

Figure 8.2 – Working with OpenAPI data models

## Data modeling

Starting from the top of the diagram, the first thing we need to do is produce the CUE code we can use to generate the data structures to configure a network device.

Although CUE can import existing structured data and generate CUE code, it may take a few iterations to get to a point where the code organization is optimal. It turned out to be faster to write this code from scratch for the example we present here. The result is in the `ch08/cue/template.cue` file (refer to the _Further_ _reading_ section).

Important Note

We won’t cover CUE syntax or any of its core concepts and principles in this book but will instead focus on its Go API. For more details about the language, please refer to CUE’s official documentation, linked in the _Further_ _reading_ section.

CUE resembles JSON with heavy influences from Go. It allows you to define data structures and map values between different data structures via references. Data generation in CUE thus becomes an exercise of data transformation with strict value typing and schema validation. Here’s a snippet from the `template.cue` file mentioned earlier, which defines three top-level objects for interfaces, routing, and VRF configuration:

```markup
package cvx
import "network.automation:input"
interface: _interfaces
router: bgp: {
    _global_bgp
}
vrf: _vrf
_global_bgp: {
    "autonomous-system": input.asn
    enable:              "on"
    "router-id":         input.loopback.ip
}
_interfaces: {
    lo: {
        ip: address: "\(input.LoopbackIP)": {}
        type: "loopback"
    }
    for intf in input.uplinks {
        "\(intf.name)": {
            type: "swp"
            ip: address: "\(intf.prefix)": {}
        }
    }
}
/* ... omitted for brevity ... */
```

Important Note

You can refer to CUE’s _References and Visibility_ tutorial (linked in the _Further reading_ section) for explanations about emitted values, references, and the use of underscores.

This file has references to an external CUE package called input, which provides the required input data for the data model in the preceding output. This separation of data templates and their inputs allows you to distribute these files separately and potentially have them come from different sources. CUE provides a guarantee that the result is always the same, no matter the order you follow to assemble those files.

## Data input

Now, let’s see how we define and provide inputs to the preceding data model. We use the same data structure we used in _Chapters 6_, _Configuration Management_, and [_Chapter 7_](https://subscription.imaginedevops.io/book/cloud-and-networking/9781800560925/2B16971_07.xhtml#_idTextAnchor161), _Automation Frameworks_, in a YAML file (`input.yaml`), which for the `cvx` lab device looks as follows:

```markup
# input.yaml
asn: 65002
loopback: 
  ip: "198.51.100.2"
uplinks:
  - name: "swp1"
    prefix: "192.0.2.3/31"
peers:
  - ip: "192.0.2.2"
    asn: 65001
```

Using CUE, we can validate that this input data is correct by building a corresponding object and introducing constraints, for example, a valid ASN range or IPv4 prefix format. CUE allows you to define extra values directly inside the schema definition, either by hardcoding defaults (`input.VRFs`) or referencing other values from the same context (`input.LoopbackIP`):

```markup
package input
import (
    "net"
)
asn: <=65535 & >=64512
loopback: ip: net.IPv4 & string
uplinks: [...{
    name:   string
    prefix: net.IPCIDR & string
}]
peers: [...{
    ip:  net.IPv4 & string
    asn: <=65535 & >=64512
}]
LoopbackIP: "\(loopback.ip)/32"
VRFs: [{name: "default"}]
```

In the main function of the example program, we use the `importInput` helper function to read the input YAML file and generate a corresponding CUE file:

```markup
import "cuelang.org/go/cue/load"
func main() {
    err := importInput()
    /* ... <continues next > ... */
}
```

The program saves the resulting file as `input.cue` in the local directory. The implementation details of this function are not too important as you can perform the same action from the command line with `cue import input.yaml -``p input`.

At this stage, we can validate that our input conforms to the schema and constraints shown earlier. For example, if we had set the `asn` value in `input.yaml` to something outside of the expected range, CUE would’ve caught and reported this error:

```markup
ch08/cue$ cue eval network.automation:input -c
asn: invalid value 10 (out of bound >=64512):
    ./schema.cue:7:16
    ./input.cue:3:6
```

## Device configuration

Now we have all the pieces in place to configure our network device. We produce the final configuration instance by compiling the template defined in the `cvx` package into a concrete CUE value. We do this in three steps.

First, we load all CUE files from the local directory, specifying the name of the package containing the template (`cvx`):

```markup
func main() {
    /* ... <continues from before > ... */
    bis := load.Instances([]string{"."}, &load.Config{
        Package: "cvx",
    })
    /* ... <continues next > ... */
}
```

Second, we compile all loaded files into a CUE value, which resolves all imports and combines the input with the template:

```markup
func main() {
    /* ... <continues from before > ... */
    ctx := cuecontext.New()
    i := ctx.BuildInstance(instances[0])
    if i.Err() != nil {
        msg := errors.Details(i.Err(), nil)
        fmt.Printf("Compile Error:\n%s\n", msg)
    }
    /* ... <continues next > ... */
}
```

Finally, we validate that we can resolve all references and that the input provides all the required fields:

```markup
func main() {
    /* ... <continues from before > ... */
    if err := i.Validate(
        cue.Final(),
        cue.Concrete(true),
    ); err != nil {
        msg := errors.Details(err, nil)
        fmt.Printf("Validate Error:\n%s\n", msg)
    }
    /* ... <continues next > ... */
}
```

Once we know the CUE value is concrete, we can safely marshal it into JSON and send it directly to the `cvx` device. The body of the `sendBytes` function implements the three-stage commit process we discussed in [_Chapter 6_](https://subscription.imaginedevops.io/book/cloud-and-networking/9781800560925/2B16971_06.xhtml#_idTextAnchor144), _Configuration Management_:

```markup
func main() {
    /* ... <continues from before > ... */
    data, err := e.MarshalJSON()
    // check error
    if err := sendBytes(data); err != nil {
        log.Fatal(err)
    }
    log.Printf("Successfully configured the device")
}
```

You can find the full program in the `ch08/cue` directory (refer to the _Further reading_ section) of this book’s GitHub repository (refer to the _Further reading_ section). The same directory includes the complete version of the CUE files with a data template and input schema and the input YAML file. Successful execution of this program should produce an output like this:

```markup
ch08/cue$ go run main.go
main.go:140: Created revisionID: changeset/cumulus/2022-05-25_20.56.51_KF9A
{
  "state": "apply",
  "transition": {
    "issue": {},
    "progress": ""
  }
}
main.go:69: Successfully configured the device
```

Keep in mind that although we focus on CUE’s Go API in this chapter, you can do the same set of actions using the CUE CLI (executable binary). This even includes the three-stage commit to submit and apply the `cvx` configuration. Using the built-in CUE scripting language, you can define any sequence of tasks, such as making HTTP calls or checking and parsing responses. You can save these actions or tasks in a special _tool_ file and they automatically become available in the `cue` binary. You can read more about this in the `ch08/cue` readme document and find example source code in the `ch08/cue/cue_tool.cue` file (refer to the _Further_ _reading_ section).

CUE has many use cases outside of what we’ve just described and different open source projects such as **Istio** and **dagger.io** (refer to the _Further reading_ section) have adopted it and use it in their products. We encourage you to explore other CUE use cases beyond what’s covered in this book, as well as similar configuration languages such as **Jsonnet** and **Dhall** (refer to the _Further_ _reading_ section).

We’ve covered a few different ways of interacting with an OpenAPI provider. For the rest of this chapter, we’ll focus on YANG-based APIs. The first one we’ll introduce is a JSON-RPC interface implementation from Nokia.

Just Imagine

# JSON-RPC

JSON-RPC is a lightweight protocol you can use to exchange structured data between a client and a server. It can work over different transport protocols, but we’ll focus only on HTTP. Although JSON-RPC is a standard, it only defines the top-level RPC layer, while payloads and operations remain specific to each implementation.

In this section, we’ll show how to use Nokia-specific YANG models to configure the srl device from our lab topology, as SR Linux supports sending and receiving YANG payloads over JSON-RPC (refer to the _Further_ _reading_ section).

We’ll try to avoid building YANG data payloads manually or relying on traditional text templating methods. The sheer size of some YANG models, as well as model deviations and augmentations, make it impossible to build the payloads manually. To do this at scale, we need to rely on a programmatic approach to build configuration instances and retrieve state data. This is where we use openconfig/ygot (YANG Go Tools) (refer to the _Further reading_ section)—a set of tools and APIs for automatic code generation from a collection of YANG models.

At a high level, the structure of the example program is analogous to the one in the _OpenAPI_ section. _Figure 8__.3_ shows the building blocks of the program we’ll review in this section:

![Figure 8.3 – Working with YANG data models](https://static.packt-cdn.com/products/9781800560925/graphics/image/B16971_08_03.jpg)

Figure 8.3 – Working with YANG data models

We’ll start by combining the auto-generated Go bindings with the input data and building a configuration instance to provision the `srl` device.

## Code generation

Starting from the top of the preceding diagram, the first step is to generate the corresponding Go code from a set of Nokia’s YANG models (refer to the _Further reading_ section). We’ll only use a subset of Nokia’s YANG models to generate the bindings to configure what we need, namely L3 interfaces, BGP, and route redistribution. This way, we keep the size of the generated Go package small and constrained to our specific use case.

Sadly, there is no universal rule for how to pinpoint the list of models you need apart from reading and understanding YANG models or reverse-engineering them from an existing configuration. Thankfully, Nokia has developed a YANG browser (refer to the _Further reading_ section) that includes a pattern-matching search that highlights the relevant XPaths and can help you find the right set of YANG models.

Once we’ve identified which models we need, we can use the ygot generator tool to build a Go package based on them. We won’t describe all the flags of this tool, as ygot’s official documentation (refer to the _Further reading_ section) covers them. Still, we want to highlight the most important options we’ll use:

-   `generate_fakeroot`: This encapsulates all generated Go data structures in a top-level _fake_ root data structure called `Device` to join all modules in a common hierarchy. Because there isn’t a YANG model that defines a universal root top-level container for all devices, network devices just add the YANG modules they support at the root (`/`). `ygot` represents the root via this _fake_ root container.
-   `path`: This flag helps `ygot` find and resolve any YANG data model imports.

The complete command to auto-generate the `srl` package and place it in the `./pkg/srl/` directory we used is this:

```markup
ch08/json-rpc$ go run \
  github.com/openconfig/ygot/generator \
    -path=yang \
    -generate_fakeroot -fakeroot_name=device \
    -output_file=pkg/srl/srl.go \
    -package_name=srl \
    yang/srl_nokia/models/network-instance/srl_nokia-bgp.yang \
    yang/srl_nokia/models/routing-policy/srl_nokia-routing-policy.yang \
    yang/srl_nokia/models/network-instance/srl_nokia-ip-route-tables.yang
```

Since the preceding command has several flags, it may be desirable to remember their exact set to make the build reproducible in the future. One alternative is to include it in a code build utility, such as make. Another, more Go-native option is to include it in the source code using the `//go:generate` directive, as you can see in the `ch08/json-rpc/main.go` file (refer to the _Further reading_ section). Thus, you can generate the same `srl` repeatedly using this command:

```markup
ch08/json-rpc$ go generate ./...
```

## Building configuration

Now that we’ve built a YANG-based Go package, we can create a programmatic instance of our desired configuration state and populate it. We do all this within Go, with the full flexibility of a general-purpose programming language at our disposal.

For example, we can design the configuration program as a set of methods, with the input model being the receiver argument. After we read and decode the input data, we create an empty _fake_ root device we extend iteratively until we build the complete YANG instance with all the relevant values we want to configure.

The benefit of using a root device is that we don’t need to worry about individual paths. We can send our payload to `/`, assuming that the resulting YANG tree hierarchy starts from the root:

```markup
import (
  api "json-rpc/pkg/srl"
)
// Input Data Model
type Model struct {
  Uplinks  []Link `yaml:"uplinks"`
  Peers    []Peer `yaml:"peers"`
  ASN      int    `yaml:"asn"`
  Loopback Addr   `yaml:"loopback"`
}
func main() {
  /* ... <omitted for brevity > ... */
  var input Model
  d.Decode(&input)
  device := &api.Device{}
  input.buildDefaultPolicy(device)
  input.buildL3Interfaces(device)
  input.buildNetworkInstance(device)
  /* ... <continues next (main) > ... */
}
```

The preceding code calls three methods on input. Let’s zoom in on `buildNetworkInstance`, responsible for L3 routing configuration. This method is where we define a _network instance_, which is a commonly used abstraction for **VPN Routing and Forwarding** (**VRF**) instances and **Virtual Switch Instances** (**VSIs**). We create a new network instance from the top-level root device to ensure we attach it to the top of the YANG tree:

```markup
func (m *Model) buildNetworkInstance(dev *api.Device) error {
  ni, err := dev.NewNetworkInstance(defaultNetInst)
  /* ... <continues next (buildNetworkInstance) > ... */
}
```

In the next code snippet, we move all uplinks and a loopback interface into the newly created network instance by defining each subinterface as a child of the default network instance:

```markup
func (m *Model) buildNetworkInstance(dev *api.Device) error {
  // ... <continues from before (buildNetworkInstance) > 
  links := m.Uplinks
  links = append(
    links,
    Link{
      Name:   srlLoopback,
      Prefix: fmt.Sprintf("%s/32", m.Loopback.IP),
    },
  )
  for _, link := range links {
    linkName := fmt.Sprintf("%s.%d", link.Name,
                            defaultSubIdx)
    ni.NewInterface(linkName)
  }
  /* ... <continues next (buildNetworkInstance) > ... */
}
```

Next, we define the global BGP settings by manually populating the BGP struct and attaching it to the `Protocols.Bgp` field of the `default` network instance:

```markup
func (m *Model) buildNetworkInstance(dev *api.Device) error {
  // ... <continues from before (buildNetworkInstance) > 
  ni.Protocols =
  &api.SrlNokiaNetworkInstance_NetworkInstance_Protocols{
    Bgp: 
    &api.
    SrlNokiaNetworkInstance_NetworkInstance_Protocols_Bgp{
      AutonomousSystem: ygot.Uint32(uint32(m.ASN)),
      RouterId:         ygot.String(m.Loopback.IP),
      Ipv4Unicast: 
      &api. 
SrlNokiaNetworkInstance_NetworkInstance_Protocols_Bgp_Ipv4Unicast{
        AdminState: api.SrlNokiaBgp_AdminState_enable,
      },
    },
  }
  /* ... <continues next (buildNetworkInstance) > ... */
}
```

The final part of the configuration is BGP neighbors. We iterate over a list of peers defined in the input data model and add a new entry under the BGP struct we set up earlier:

```markup
func (m *Model) buildNetworkInstance(dev *api.Device) error {
  // ... <continues from before (buildNetworkInstance) > 
  ni.Protocols.Bgp.NewGroup(defaultBGPGroup)
  for _, peer := range m.Peers {
    n, err := ni.Protocols.Bgp.NewNeighbor(peer.IP)
    // check error
    n.PeerAs = ygot.Uint32(uint32(peer.ASN))
    n.PeerGroup = ygot.String(defaultBGPGroup)
  }
  /* ... <continues next (buildNetworkInstance) > ... */
}
```

When we finish populating the Go structs, we make sure that all provided values are correct and match the YANG constraints. We can do this with a single call to the `Validate` method on the parent container:

```markup
func (m *Model) buildNetworkInstance(dev *api.Device) error {
    /* ... <continues from before (buildNetworkInstance) > ... */
    if err := ni.Validate(); err != nil {
        return err
    }
    return nil
}
```

## Device configuration

Once we have populated a YANG model instance with all the input values, the next step is to send it to the target device. We do this in a few steps:

1.  We use a `ygot` helper function to produce a map from the current YANG instance. This map is ready to be serialized into JSON according to the rules defined in RFC7951.
2.  We use the standard `encoding/json` library to build a single JSON-RPC request that updates the entire YANG tree with our configuration changes.
3.  Using the standard `net/http` package, we send this request to the `srl` device:
    
    ```markup
    func main() {
    ```
    
    ```markup
        /* ... <continues from before (main) > ... */
    ```
    
    ```markup
        v, err := ygot.ConstructIETFJSON(device, nil)
    ```
    
    ```markup
        // check error
    ```
    
    ```markup
        value, err := json.Marshal(RpcRequest{
    ```
    
    ```markup
            Version: "2.0",
    ```
    
    ```markup
            ID:      0,
    ```
    
    ```markup
            Method:  "set",
    ```
    
    ```markup
            Params: Params{
    ```
    
    ```markup
                Commands: []*Command{
    ```
    
    ```markup
                    {
    ```
    
    ```markup
                        Action: "update",
    ```
    
    ```markup
                        Path:   "/",
    ```
    
    ```markup
                        Value:  v,
    ```
    
    ```markup
                    },
    ```
    
    ```markup
                },
    ```
    
    ```markup
            },
    ```
    
    ```markup
        })
    ```
    
    ```markup
        // check error
    ```
    
    ```markup
        req, err := http.NewRequest(
    ```
    
    ```markup
            "POST",
    ```
    
    ```markup
            hostname,
    ```
    
    ```markup
            bytes.NewBuffer(value),
    ```
    
    ```markup
        )
    ```
    
    ```markup
        resp, err := client.Do(req)
    ```
    
    ```markup
         // check error
    ```
    
    ```markup
        defer resp.Body.Close()
    ```
    
    ```markup
        if resp.StatusCode != http.StatusOK {
    ```
    
    ```markup
            log.Printf("Status: %s", resp.Status)
    ```
    
    ```markup
        }
    ```
    

You can find the complete program that configures the srl device in the `ch08/json-rpc` directory (refer to the _Further reading_ section) of this book’s GitHub repository. To run it, `cd` into this folder and run the following command:

```markup
ch08/json-rpc$ go run main.go
2022/04/26 13:09:03 Successfully configured the device
```

This program only verifies that we executed the RPC successfully; it doesn’t yet check to confirm that it had the desired effect, which we will discuss later in this chapter. As with most HTTP-based protocols, a single RPC is a single transaction, so you can assume the target device applied the changes, as long as you receive a successful response. It’s worth mentioning that some JSON-RPC implementations have more session control functions that allow multistage commits, rollbacks, and other features.

In the following section, we’ll take a similar approach of configuring a network device based on its YANG models but introduce a couple of twists to show OpenConfig models and the RESTCONF API.

Just Imagine

# RESTCONF

The IETF designed RESTCONF as an HTTP-based alternative to NETCONF that offers **Create, Read, Update, and Delete** (**CRUD**) operations on a conceptual datastore containing YANG-modeled data. It may lack some NETCONF features, such as different datastores, exclusive configuration locking, and batch and rollback operations, but the exact set of supported and unsupported features depends on the implementation and network device capabilities. That said, because it uses HTTP methods and supports JSON encoding, RESTCONF reduces the barrier of entry for external systems to integrate and inter-operate with a network device.

RESTCONF supports a standard set of CRUD operations through HTTP methods: POST, PUT, PATCH, GET, and DELETE. RESTCONF builds HTTP messages with the YANG XPath translated into a REST-like URI and it transports the payload in the message body. Although RESTCONF supports both XML and JSON encoding, we will only focus on the latter, with the rules of the encoding defined in RFC7951. We’ll use Arista’s EOS as a test device, which has its RESTCONF API enabled when launching the lab topology.

The structure of the program we’ll create in this section is the same as for the JSON-RPC example illustrated in _Figure 8__.3_.

## Code generation

The code generation process is almost the same as the one we followed in the _JSON-RPC_ section. We use openconfig/ygot (refer to the _Further reading_ section) to generate a Go package from a set of YANG models that EOS supports. But there are a few notable differences that are worth mentioning before moving forward:

-   Instead of vendor-specific YANG models, we use vendor-neutral OpenConfig models, which Arista EOS supports.
-   When generating Go code with openconfig/ygot (refer to the _Further reading_ section), you might run into situations when more than one model is defined in the same namespace. In those cases, you can use the `-exclude_modules` flag to ignore a certain YANG model without having to remove its source file from the configured search path.
-   We enable OpenConfig path compression to optimize the generated Go code by removing the YANG containers containing `list` nodes. Refer to the `ygen` library design documentation for more details (_Further reading_).
-   We also show an alternative approach where we don’t generate a _fake_ root device. As a result, we can’t apply all the changes in a single RPC. Instead, we have to make more than one HTTP call, each with its own unique URI path.

Before we can generate the Go code, we need to identify the supported set of Arista YANG models (refer to the _Further reading_ section) and copy them into the `yang` directory. We use the following command to generate the `eos` Go package from that list of models:

```markup
ch08/restconf$ go run github.com/openconfig/ygot/generator \
  -path=yang \
  -output_file=pkg/eos/eos.go \
  -compress_paths=true \
  -exclude_modules=ietf-interfaces \
  -package_name=eos \
  yang/openconfig/public/release/models/bgp/openconfig-bgp.yang \
  yang/openconfig/public/release/models/interfaces/openconfig-if-ip.yang \
  yang/openconfig/public/release/models/network-instance/openconfig-network-instance.yang \
  yang/release/openconfig/models/interfaces/arista-intf-augments-min.yang
```

For the same reasons we described in the _JSON-RPC_ section, we can also embed this command into the Go source code to generate the same Go package using the following command instead:

```markup
ch08/restconf$ go generate ./...
```

## Building configuration

In this example, we won’t apply all changes in a single HTTP call so that we can show you how to update a specific part of a YANG tree without affecting other, unrelated parts. In the preceding section, we worked around that by using an `Update` operation, which merges the configuration we send with the existing configuration on the device.

But in certain cases, we want to avoid the _merge_ behavior and ensure that only the configuration we send is present on the device (declarative management). For that, we could’ve imported all existing configurations and identified the parts that we want to keep or replace before sending a new configuration version to the target device. Instead, we create a configuration for the specific parts of a YANG tree via a series of RPCs.

To simplify RESTCONF API calls, we create a special `restconfRequest` type that holds a URI path and a corresponding payload to send to the device. The `main` function starts with parsing the inputs for the data model and preparing a variable to store a set of RESTCONF RPCs:

```markup
type restconfRequest struct {
    path    string
    payload []byte
}
func main() {
    /* ... <omitted for brevity > ... */
    var input Model
    err = d.Decode(&input)
    // check error
    var cmds []*restconfRequest
    /* ... <continues next > ... */
}
```

As in the JSON-RPC example, we build the desired configuration instance in a series of method calls. This time, each method returns one `restConfRequest` that has enough details to build an HTTP request:

```markup
func main() {
    /* ... <continues from before > ... */ 
    l3Intfs, err := input.buildL3Interfaces()
    // check error
    cmds = append(cmds, l3Intfs...)
    bgp, err := input.buildBGPConfig()
    // check error
    cmds = append(cmds, bgp)
    redistr, err := input.enableRedistribution()
    // check error
    cmds = append(cmds, redistr)
    /* ... <continues next > ... */
}
```

Let’s examine one of these methods that creates a YANG configuration from our inputs. The `enableRedistribution` method generates a configuration to enable redistribution between a directly connected table and the BGP **Routing Information Base** (**RIB**). OpenConfig defines a special `TableConnection` struct that uses a pair of YANG enums to identify the redistribution source and destination:

```markup
const defaultNetInst = "default"
func (m *Model) enableRedistribution() (*restconfRequest, error) {
    netInst := &api.NetworkInstance{
        Name: ygot.String(defaultNetInst),
    }
    _, err := netInst.NewTableConnection(
        api.OpenconfigPolicyTypes_INSTALL_PROTOCOL_TYPE_DIRECTLY_CONNECTED,
        api.OpenconfigPolicyTypes_INSTALL_PROTOCOL_TYPE_BGP,
        api.OpenconfigTypes_ADDRESS_FAMILY_IPV4,
    )
    
    /* ... <omitted for brevity > ... */
    value, err := ygot.Marshal7951(netInst)
    // check error
    return &restconfRequest{
        path: fmt.Sprintf(
            "/network-instances/network-instance=%s",
            defaultNetInst,
        ),
        payload: value,
    }, nil
}
```

The rest of the code in _Figure 8__.3_ shows the building blocks of the program we review in this section.

## Device configuration

Once we’ve prepared all the required RESTCONF RPCs, we can send them to the device. We iterate over each `restconfRequest` and pass it to a helper function, catching any returned errors.

The `restconfPost` helper function has just enough code to build an HTTP request using the `net/http` package and send it to the `ceos` device:

```markup
const restconfPath = "/restconf/data"
func restconfPost(cmd *restconfRequest) error {
  baseURL, err := url.Parse(
    fmt.Sprintf(
      "https://%s:%d%s",
      ceosHostname,
      defaultRestconfPort,
      restconfPath,
    ),
  )
  // return error if not nil
  baseURL.Path = path.Join(restconfPath, cmd.path)
  req, err := http.NewRequest(
    "POST",
    baseURL.String(),
    bytes.NewBuffer(cmd.payload),
  )
  // return error if not nil
  req.Header.Add("Content-Type", "application/json")
  req.Header.Add(
    "Authorization",
    "Basic "+base64.StdEncoding.EncodeToString(
      []byte(
        fmt.Sprintf("%s:%s", ceosUsername, ceosPassword),
      ),
    ),
  )
  client := &http.Client{Transport: &http.Transport{
        TLSClientConfig: 
          &tls.Config{
            InsecureSkipVerify: true
          },
      }
  }
  resp, err := client.Do(req)
  /* ... <omitted for brevity > ... */
}
```

You can find the complete program in the `ch08/restconf` directory (refer to the _Further reading_ section) of this book’s GitHub repository. Running it from a host running the lab topology should produce a similar output to this:

```markup
ch08/restconf$ go run main.go
2022/04/28 20:49:16 Successfully configured the device
```

At this point, we should have all three nodes of our lab topology fully configured. Still, we haven’t confirmed that what we’ve done has had the desired effect. In the next section, we’ll go through a process of state validation and show how you can do it using network APIs.

Just Imagine

# State validation

In the last three sections of this chapter, we pushed device configs without verifying that the configuration changes had the desired effect. This is because we need all devices configured before we can validate the resulting converged operational state. Now, with all the code examples from the _OpenAPI_, _JSON-RPC_, and _RESTCONF_ sections executed against the lab topology, we can verify whether we achieved our configuration intent—establish end-to-end reachability between loopback IP addresses of all three devices.

In this section, we’ll use the same protocols and modeling language we used earlier in this chapter to validate that each lab device can see the loopback IP address of the other two lab devices in its **Forwarding Information Base** (**FIB**) table. You can find the complete code for this section in the `ch08/state` directory (refer to the _Further reading_ section) of this book’s GitHub repository. Next, we’ll examine a single example of how you can do this with Arista’s cEOS (`ceos`) lab device.

## Operational state modeling

One thing we need to be mindful of when talking about the operational state of a network element is the difference between the applied and the derived state, as described by the YANG operational state IETF draft (refer to the _Further reading_ section). The former refers to the currently active device configuration and should reflect what an operator has already applied. The latter is a set of read-only values that result from the device’s internal operations, such as CPU or memory utilization, and interaction with external elements, such as packet counters or BGP neighbor state. Although we aren’t explicitly mentioning it when we’re talking about an operational state, assume we’re referring to the derived state unless we state otherwise.

Historically, there’ve been different ways to model the device’s operational state in YANG:

-   You could either enclose everything in a top-level container or read from a separate `state` datastore, completely distinct from the `config` container/datastore we use for configuration management.
-   Another way is to create a separate `state` container for every YANG sub-tree alongside the `config` container. This is what the YANG operational state IETF draft (refer to the _Further reading_ section) describes.

Depending on which approach you use, you may need to adjust how you construct your RPC request. For example, the `srl` device needs an explicit reference to the `state` datastore. What we show in the next code example is the alternative approach, where you retrieve a part of the YANG sub-tree and extract the relevant state information from it.

It’s worth noting that OpenAPI is less strict about the structure and composition of its models and the state may come from a different part of a tree or require a specific query parameter to reference the operational datastore, depending on the implementation.

## Operational state processing

Configuration management workflows typically involve the processing of some input data to generate a device-specific configuration. This is a common workflow that we often use to show the capabilities of an API. But there is an equally important workflow that involves operators retrieving state data from a network device, which they process and verify. In that case, the information flows in the opposite direction—from a network device to a client application.

At the beginning of this chapter, we discussed the configuration management workflow, so now we want to give a high-level overview of the state retrieval workflow:

1.  We start by querying a remote API endpoint, represented by a set of URL and HTTP query parameters.
2.  We receive an HTTP response, which has a binary payload attached to it.
3.  We unmarshal this payload into a Go struct that follows the device’s data model.
4.  Inside this struct, we look at the relevant parts of the state we can extract and evaluate.

The following code snippet from the `ch08/state` program (refer to the _Further reading_ section) is a concrete example of this workflow. The program structure follows the same pattern we described in the _State validation_ section of [_Chapter 6_](https://subscription.imaginedevops.io/book/cloud-and-networking/9781800560925/2B16971_06.xhtml#_idTextAnchor144), _Configuration Management_. Hence, in this chapter, we’ll only zoom in on the most relevant part—the `GetRoutes` function, which connects to the `ceos` device and retrieves the content of its routing table.

It starts by building an HTTP request with the device-specific login information:

```markup
func (r CEOS) GetRoutes(wg *sync.WaitGroup) {
  client := resty.NewWithClient(&http.Client{
    Transport: &http.Transport{
      TLSClientConfig: &tls.Config{
        InsecureSkipVerify: true},
    },
  })
  client.SetBaseURL("https://" + r.Hostname + ":6020")
  client.SetBasicAuth(r.Username, r.Password)
  resp, err := client.R().
    SetHeader("Accept", "application/yang-data+json").
    Get(fmt.Sprintf("/restconf/data/network-instances/network-instance=%s/afts", "default"))
  /* ... <continues next > ... */
}
```

The **Abstract Forwarding Table** (**AFT**) in the code example is an OpenConfig representation of the FIB (routing) table and the GET API call retrieves a JSON representation of the default **Virtual Routing and Forwarding** (**VRF**) routing table.

Next, we create an instance of the Go struct corresponding to the part of the YANG tree we queried and pass it to the `Unmarshal` function for deserialization. The resulting Go struct now has one `Ipv4Entry` value for each entry in the default FIB and we store that list of prefixes in the `out` slice:

```markup
import eosAPI "restconf/pkg/eos"
func (r CEOS) GetRoutes(wg *sync.WaitGroup) {
  /* ... <continues from before > ... */
  response := &eosAPI.NetworkInstance_Afts{}
  err := eosAPI.Unmarshal(resp.Body(), response)
  // process error
  out := []string{}
  for key := range response.Ipv4Entry {
    out = append(out, key)
  }
  /* ... <omitted for brevity > ... */
  go checkRoutes(r.Hostname, out, expectedRoutes, wg)
}
```

In this example, we import the `eos` package (`restconf/pkg/eos`) we auto-generated in the _RESTCONF_ section of this chapter, which lives outside the root directory of this program. To do this, we add the `replace restconf => ../restconf/` instruction to this program’s `go.mod` file (`ch08/state/go.mod`; refer to the _Further_ _reading_ section).

For the remaining lab devices, we follow a similar state retrieval workflow. The only difference is in the YANG paths and the model-based Go structs we use for deserialization. You can find the full program code in the `ch08/state` directory (refer to the _Further reading_ section) of this book’s GitHub repository.

In this chapter, we have covered network APIs based on HTTP version 1.1 that use common encoding formats, such as JSON. Although HTTP is still very popular and this is unlikely to change soon, it has its own limitations that may manifest themselves in large-scale deployments. HTTP 1.1 is a text-based protocol, which means it’s not efficient on the wire and its client-server origins make it difficult to adapt it for bi-directional streaming. The next version of this protocol, HTTP/2, overcomes these shortcomings. HTTP/2 is the transport protocol of the gRPC framework, which is what we’ll examine in the next section.

Just Imagine

# gRPC

Network automation opens a door that until recently seemed closed or at least prevented network engineers from reusing technologies that have had success in other areas, such as microservices or cloud infrastructure.

One of the most recent advances in network device management is the introduction of gRPC. We can use this high-performance RPC framework for a wide range of network operations, from configuration management to state streaming and software management. But performance is not the only thing that is appealing about gRPC. Just like with YANG and OpenAPI apps, gRPC auto-generates client and server stubs in different programming languages, which enables us to create an ecosystem of tools around the API.

In this section, we’ll go over the following topics to help you understand the gRPC API better:

-   Protobuf
-   gRPC transport
-   Defining gRPC services
-   Configuring network devices with gRPC
-   Streaming telemetry from a network device with gRPC

## Protobuf

gRPC uses protobuf as its **Interface Definition Language** (**IDL**) to allow you to share structured data between remote software components that may be written in different programming languages.

When working with protobuf, one of the first steps is to model the information you’re serializing by creating a protobuf file. This file has a list of _messages_ defining the structure and type of data to exchange.

If we take the input data model we have been using throughout this book as an example and encode it in a `.proto` file, it would look something like this:

```markup
message Router {
  repeated Uplink uplinks = 1;
  repeated Peer peers = 2;
  int32 asn = 3;
  Addr loopback = 4; 
}
message Uplink {
    string name = 1;
    string prefix = 2;
}
message Peer {
    string ip = 1;
    int32 asn = 2;
}
message Addr {
  string ip = 1;
}
```

Each field has an explicit type and a unique sequence number that identifies it within the enclosing message.

The next step in the workflow, just like with OpenAPI or YANG, is to generate bindings for Go (or any other programming language). For this, we use the protobuf compiler, protoc, which generates the source code with data structures and methods to access and validate different fields:

```markup
ch08/protobuf$ protoc --go_out=. model.proto
```

The preceding command saves the bindings in a single file, `pb/model.pb.go`. You can view the contents of this file to see what structs and functions you can use. For example, we automatically get this `Router` struct, which is what we had to define manually before:

```markup
type Router struct {
  Uplinks  []*Uplink 
  Peers    []*Peer   
  Asn      int32     
  Loopback *Addr
}
```

Protobuf encodes a series of key-value pairs in a binary format similar to how routing protocols encode **Type-Length-Values** (**TLVs**). But instead of sending the key name and a declared type for each field, it just sends the field number as the key with its value appended to the end of the byte stream.

As with TLVs, Protobuf needs to know the length of each value to encode and decode a message successfully. For this, Protobuf encodes a wire type in the 8-bit key field along with the field number that comes from the `.proto` file. The following table shows the wire types available:

<table id="table001-2" class="No-Table-Style _idGenTablePara-1"><colgroup><col> <col> <col></colgroup><tbody><tr class="No-Table-Style"><td class="No-Table-Style"><p><span class="No-Break"><strong class="bold">Type</strong></span></p></td><td class="No-Table-Style"><p><span class="No-Break"><strong class="bold">Meaning</strong></span></p></td><td class="No-Table-Style"><p><span class="No-Break"><strong class="bold">Used For</strong></span></p></td></tr><tr class="No-Table-Style"><td class="No-Table-Style"><p>0</p></td><td class="No-Table-Style"><p><span class="No-Break">Varint</span></p></td><td class="No-Table-Style"><p>int32, int64, uint32, uint64, sint32, sint64, <span class="No-Break">bool, enum</span></p></td></tr><tr class="No-Table-Style"><td class="No-Table-Style"><p>1</p></td><td class="No-Table-Style"><p><span class="No-Break">64-bit</span></p></td><td class="No-Table-Style"><p>fixed64, <span class="No-Break">sfixed64, double</span></p></td></tr><tr class="No-Table-Style"><td class="No-Table-Style"><p>2</p></td><td class="No-Table-Style"><p><span class="No-Break">Length-delimited</span></p></td><td class="No-Table-Style"><p>string, bytes, embedded messages, packed <span class="No-Break">repeated fields</span></p></td></tr><tr class="No-Table-Style"><td class="No-Table-Style"><p>5</p></td><td class="No-Table-Style"><p><span class="No-Break">32-bit</span></p></td><td class="No-Table-Style"><p>fixed32, <span class="No-Break">sfixed32, float</span></p></td></tr></tbody></table>

Table 8.1 – Protobuf wire types

This generates a dense message (small output) that a CPU can process faster compared to a JSON- or XML-encoded message. The downside is the message you generate is not human-readable in its native format and it’s only meaningful if you have the message definition (proto file) to find out the name and type for each field.

### Protobuf on the wire

One of the easiest ways to see how protobuf looks in a binary format is to save it into a file. In our book’s GitHub repository, we have an example in the `ch08/protobuf/write` directory (refer to the _Further reading_ section) that reads a sample `input.yaml` file and populates the data structure generated from the `.proto` file we discussed earlier. We then serialize and save the result into a file we name `router.data`. You can use the following command to execute this example:

```markup
ch08/protobuf/write$ go run protobuf
```

You can see the content of the generated protobuf message by viewing the file with `hexdump -C router.data`. If we group some bytes for convenience and refer to the proto definition file, we can make sense of the data, as shown in the following figure:

![Figure 8.4 – Protobuf-encoded message](https://static.packt-cdn.com/products/9781800560925/graphics/image/B16971_08_04.jpg)

Figure 8.4 – Protobuf-encoded message

To give you an idea of how efficient the protobuf encoding is, we’ve included a couple of JSON files encoding the same data. The `router.json` file is a compact (space-free) JSON encoding. The second version, called `router_ident.json`, has the same JSON payload indented with extra spaces, which can happen if you generate JSON from a text template or use _pretty print_ functions before sending the data over the network:

```markup
ch08/protobuf$ ls -ls router* | awk '{print $6, $10}'
108 router.data
454 router_indent.json
220 router.json
```

The difference between JSON and protobuf is quite stark and can become very important when transferring and encoding/decoding large datasets.

Now that we know some basics about gRPC data encoding, we can move on to the protocol used to transfer these messages.

## gRPC transport

Besides efficient binary encoding and enabling simpler framing to serialize your data—compared to newline-delimited plain text—the gRPC framework also attempts to exchange those messages as efficiently as possible over the network.

While you can only process one request/response message at a time with HTTP/1.1, gRPC makes use of HTTP/2 to multiplex parallel requests over the same TCP connection. Another benefit of HTTP/2 is that it supports header compression. _Table 8.2_ shows the various transport methods used by different APIs:

<table id="table002" class="No-Table-Style _idGenTablePara-1"><colgroup><col> <col> <col></colgroup><tbody><tr class="No-Table-Style"><td class="No-Table-Style"><p><span class="No-Break"><strong class="bold">API</strong></span></p></td><td class="No-Table-Style"><p><span class="No-Break"><strong class="bold">Transport</strong></span></p></td><td class="No-Table-Style"><p><span class="No-Break"><strong class="bold">RPC/Methods</strong></span></p></td></tr><tr class="No-Table-Style"><td class="No-Table-Style"><p><span class="No-Break">NETCONF</span></p></td><td class="No-Table-Style"><p><span class="No-Break">SSH</span></p></td><td class="No-Table-Style"><p>get-config, edit-config, <span class="No-Break">commit, lock</span></p></td></tr><tr class="No-Table-Style"><td class="No-Table-Style"><p><span class="No-Break">RESTCONF</span></p></td><td class="No-Table-Style"><p><span class="No-Break">HTTP</span></p></td><td class="No-Table-Style"><p>GET, POST, <span class="No-Break">DELETE, PUT</span></p></td></tr><tr class="No-Table-Style"><td class="No-Table-Style"><p><span class="No-Break">gRPC</span></p></td><td class="No-Table-Style"><p><span class="No-Break">HTTP/2</span></p></td><td class="No-Table-Style"><p>Unary, server streaming, client streaming, <span class="No-Break">bidirectional streaming</span></p></td></tr></tbody></table>

Table 8.2 – API comparative table

Compared to the older network APIs, gRPC not only allows you to make unary or single requests, but it also supports full-duplex streaming. Both the client and server can stream data simultaneously, so you no longer need to work around the limitations of the traditional client-server mode of interaction.

## Defining gRPC services

gRPC uses Protobuf to define statically typed services and messages in a file that we can use to generate the code for client and server applications to consume. gRPC abstracts the underlying transport and serialization details, allowing developers to focus on the business logic of their applications instead.

A gRPC service is a collection of RPCs that accept and return protobuf messages. In the following output, you can see a snippet from Cisco IOS XR’s proto file called `ems_grpc.proto` (refer to the _Further reading_ section). This file defines a gRPC service called `gRPCConfigOper` with several RPCs to perform a standard set of configuration management operations:

```markup
syntax = "proto3";
service gRPCConfigOper {
  rpc GetConfig(ConfigGetArgs) returns(stream ConfigGetReply) {};
        
  rpc MergeConfig(ConfigArgs) returns(ConfigReply) {};
    
  rpc DeleteConfig(ConfigArgs) returns(ConfigReply) {};
    
  rpc ReplaceConfig(ConfigArgs) returns(ConfigReply) {};
  /* ... <omitted for brevity > ... */
  rpc CreateSubs(CreateSubsArgs) returns(stream CreateSubsReply) {};
}
```

As well as the configuration management operations, this Cisco IOS XR protobuf definition includes a streaming telemetry subscription (`CreateSubs`) RPC. The message format for the request and response is also part of the `ems_grpc.proto` file (refer to the _Further reading_ section). For example, to invoke the telemetry subscription RPC, the client has to send a `ConfigArgs` message and the server (router) should reply with a stream of `CreateSubsReply` messages.

Unlike with NETCONF, where **Request for Comments** (**RFC**) documents predefine all RPCs, networking vendors didn’t initially push for a standard set of gRPC services. This flexibility comes with a cost, as any other vendor could define a similar service, but with different names and message types. Here, you can see a snippet from Juniper’s telemetry protobuf file called `telemetry.proto` (refer to the _Further_ _reading_ section):

```markup
syntax = "proto3";
service OpenConfigTelemetry {
  rpc telemetrySubscribe(SubscriptionRequest) returns (stream OpenConfigData) {}
  /* ... <omitted for brevity > ... */
  rpc getTelemetryOperationalState(GetOperationalStateRequest) returns(GetOperationalStateReply) {}
  rpc getDataEncodings(DataEncodingRequest) returns (DataEncodingReply) {}
}
```

This is something that the OpenConfig community is addressing with the definition of vendor-agnostic services, such as gNMI (`gnmi.proto`; refer to the _Further reading_ section), which we will explore in the next chapter:

```markup
service gNMI {
  rpc Capabilities(CapabilityRequest) returns (CapabilityResponse);
  rpc Get(GetRequest) returns (GetResponse);
  rpc Set(SetRequest) returns (SetResponse);
  rpc Subscribe(stream SubscribeRequest) returns (stream SubscribeResponse);
}
```

Now, let’s see how you can use these RPCs with Go.

## Configuring network devices with gRPC

In our example program, we configure an IOS XR device with the `ReplaceConfig` RPC, defined in a service called `gRPCConfigOper`. You can find all the source code for this program in the `ch08/grpc` directory of this book’s GitHub repository (refer to the _Further reading_ section). You can use the following command to execute this program against a test device in Cisco’s DevNet sandbox:

```markup
ch08/grpc$ go run grpc
```

Following the same configuration management workflow we’ve used throughout this chapter, we’ll start by generating the code for the following gRPC service:

```markup
service gRPCConfigOper { 
  rpc ReplaceConfig(ConfigArgs) returns(ConfigReply) {};
}
message ConfigArgs {
  int64 ReqId = 1;
  string yangjson = 2;
  bool   Confirmed = 3;
  uint32  ConfirmTimeout = 4;
}
```

One thing to remember when working with gRPC-based network APIs is that they might not define the full data tree natively as protobuf schemas. In the preceding example, one field defines a string called `yangjson` that expects a YANG-based JSON payload, not exploring any further what might be inside that “string.” Carrying a YANG-based JSON payload is what we also did in the JSON-RPC and RESTCONF examples. In a sense, gRPC serves as a thin RPC wrapper in this example, not too different from JSON-RPC. We are still doing the configuration management work with YANG-based data structures.

Since we’re now using both gRPC and YANG schemas, we have to use `protoc` together with `ygot` to generate their respective bindings. We run the `protoc` command to generate the code from the proto definition in `ch08/grpc/proto` (refer to the _Further reading_ section) and `ygot` to generate code from a set of OpenConfig YANG models. You can find the exact set of commands in the `ch08/grpc/generate_code` file (refer to the _Further_ _reading_ section).

Before we can connect to the target device, we need to gather all the information we need to run the program, so we reuse the data structures from [_Chapter 6_](https://subscription.imaginedevops.io/book/cloud-and-networking/9781800560925/2B16971_06.xhtml#_idTextAnchor144), _Configuration Management_, to store this data:

```markup
type Authentication struct {
  Username string
  Password string
}
type IOSXR struct {
  Hostname string
  Authentication
}
type xrgrpc struct {
  IOSXR
  conn *grpc.ClientConn
  ctx  context.Context
}
```

We start the `main` function of the program by populating the access credentials and processing the device configuration inputs, just like in other examples in the book:

```markup
func main() {
  iosxr := xrgrpc{
    IOSXR: IOSXR{
      Hostname: "sandbox-iosxr-1.cisco.com",
      Authentication: Authentication{
        Username: "admin",
        Password: "C1sco12345",
      },
    },
  }
  src, err := os.Open("input.yml")
  // process error
  defer src.Close()
  d := yaml.NewDecoder(src)
  var input Model
  err = d.Decode(&input)
  /* ... <continues next > ... */
}
```

Next, we use the `ygot` Go bindings from the `grpc/pkg/oc` package to prepare the `yangjson` payload. We build the BGP configuration in the `buildNetworkInstance` method in the same way we showed in the _JSON-RPC_ section of this chapter. Once the `oc.Device` struct is fully populated, we serialize it into a JSON string:

```markup
func main() {
  /* ... <continues from before > ... */
  device := &oc.Device{}
  input.buildNetworkInstance(device)
  payload, err := ygot.EmitJSON(device,
  &ygot.EmitJSONConfig{
    Format: ygot.RFC7951,
    Indent: "  ",
    RFC7951Config: &ygot.RFC7951JSONConfig{
      AppendModuleName: true,
    },
  })
  /* ... <continues next > ... */
}
```

To simplify the interactions with the target device, we created a thin wrapper around the gRPC API. We define a handful of method receivers for the `xrgrpc` type that implement things such as initial connection establishment and deleting or replacing RPCs. This is how we connect and replace the target device’s configuration:

```markup
func main() {
  /* ... <continues from before > ... */
  iosxr.Connect()
  defer router.conn.Close()
  iosxr.ReplaceConfig(payload)
  /* ... <continues next > ... */ 
}
```

Looking closer at the `ReplaceConfig` method, we can see exactly how to invoke the required RPC. We dynamically generate a random ID and populate the `ConfigArg` message with the YANG-based JSON payload that we generated with `ygot` a couple of steps before. The inner `ReplaceConfig` method is the one that the `protoc` command automatically generated for us:

```markup
func (x *xrgrpc) ReplaceConfig(json string) error {
  // Random int64 for id
  id := rand.Int63()
  // 'g' is the gRPC stub.
  g := xr.NewGRPCConfigOperClient(x.conn)
  // We send 'a' to the router via the stub.
  a := xr.ConfigArgs{ReqId: id, Yangjson: json}
  // 'r' is the result that comes back from the target.
  r, err := g.ReplaceConfig(x.ctx, &a)
  // process error
  return nil
}
```

The configuration payload we send in this case is a string blob, but we can also encode the content fields with protobuf if the target devices support this. This is what we’ll examine next with a streaming telemetry example.

## Streaming telemetry from a network device with gRPC

gRPC streaming capabilities allow network devices to send data over a persistent TCP connection either continuously (stream) or on demand (poll). We’ll continue with the same program we started earlier and reuse the same connection we set up to configure a network device to subscribe to a telemetry stream.

Even though we initiated a connection to the Cisco IOS XR device, the data now flows in the opposite direction. This means we need to be able to decode the information we receive and there are two different ways of doing this.

Once we’ve configured the device, we request it to stream the operational state of all BGP neighbors. In the first scenario, we’ll cover the case where you have the BGP neighbor proto definition to decode the messages you get. Then, we’ll examine a less efficient option where a proto definition is unnecessary.

### Decoding YANG-defined data with Protobuf

We use the `CreateSubs` RPC to subscribe to a telemetry stream. We need to submit the subscription ID we want to stream and choose an encoding option between `gpb` for protobuf or `gpbkv` for an option we’ll explore at the end of this chapter. The following output shows the proto definition of this RPC and its message types:

```markup
service gRPCConfigOper { 
  rpc CreateSubs(CreateSubsArgs) returns(stream CreateSubsReply) {};
}
message CreateSubsArgs {
  int64 ReqId = 1;
  int64 encode = 2;
  string subidstr = 3;
  QOSMarking qos = 4;
  repeated string Subscriptions = 5;
}
message CreateSubsReply {
  int64 ResReqId = 1;
  bytes data = 2;
  string errors = 3;
}
```

Similar to the configuration part of the program, we create a helper function to submit the request to the router. The main difference is that now the reply is a data stream. We store the result of `CreateSubs` in a variable we call `st`.

For data streams, gRPC gives us the `Recv` method, which blocks until it receives a message. To continue processing in the main thread, we run an anonymous function in a separate goroutine that calls the auto-generated `GetData` method. This method returns the `data` field of each message we get and we send it over a channel (`b`) back to the main goroutine:

```markup
func (x *xrgrpc) GetSubscription(sub, enc string) (chan []byte, chan error, error) {
  /* ... <omitted for brevity > ... */
  
  // 'c' is the gRPC stub.
  c := xr.NewGRPCConfigOperClient(x.conn)
  // 'b' is the bytes channel where telemetry is sent.
  b := make(chan []byte)
  a := xr.CreateSubsArgs{
        ReqId: id, Encode: encoding, Subidstr: sub}
  // 'r' is the result that comes back from the target.
  st, err := c.CreateSubs(x.ctx, &a)
  // process error
  go func() {
    r, err := st.Recv()
    /* ... <omitted for brevity > ... */
    for {
      select {
      /* ... <omitted for brevity > ... */
      case b <- r.GetData():
      /* ... <omitted for brevity > ... */
      }
    }
  }()
  return b, e, err
}
```

The `data` field, and hence the data we receive in channel `b`, consist of arrays of bytes that we need to decode. We know this is a streaming telemetry message, so we use its proto-generated code to decode its fields. _Figure 8__.5_ shows an example of how we can get to BGP state information by following the proto file definitions:

![Figure 8.5 – Protobuf telemetry message (protobuf)](https://static.packt-cdn.com/products/9781800560925/graphics/image/B16971_08_06.jpg)

Figure 8.5 – Protobuf telemetry message (protobuf)

Back in the main goroutine, we listen out for what the `GetSubscription` channel returns and iterate over each message we get. We unmarshal the data received into a `Telemetry` message. At this point, we have access to the general telemetry data, so we can use the auto-generated functions to access some of its fields, such as the timestamp and the encoding path:

```markup
func main() {
  /* ... <omitted for brevity > ... */
  ch, errCh, err := router.GetSubscription("BGP", "gpb")
  // process error
    
  for msg := range ch {
    message := new(telemetry.Telemetry)
    proto.Unmarshal(msg, message)
        
    t := time.UnixMilli(int64(message.GetMsgTimestamp()))
    fmt.Printf(
      "Time: %v\nPath: %v\n\n",
      t.Format(time.ANSIC),
      message.GetEncodingPath(),
    )
    /* ... <continues next > ... */
  }
}
```

Following that, we extract the content of the `data_bgp` field to access the BGP data encoded with protobuf. Cisco IOS XR lists the items in rows, so for each one, we unmarshal the content into the auto-generated `BgpNbrBag` data structure, from where we can access all operational information of a BGP neighbor. This way, we get the connection state and the IPv4 address of the BGP peer, which we print to the screen:

```markup
func main() {
  for msg := range ch {
    /* ... <continues from before > ... */  
    for _, row := range message.GetDataGpb().GetRow() {
      content := row.GetContent()
      nbr := new(bgp.BgpNbrBag)
      err = proto.Unmarshal(content, nbr)
      if err != nil {
        fmt.Printf("could decode Content: %v\n", err)
        return
      }
      state := nbr.GetConnectionState()
      addr := nbr.GetConnectionRemoteAddress().Ipv4Address
      fmt.Println("  Neighbor: ", addr)
      fmt.Println("  Connection state: ", state)
    }
  }
}
```

If you don’t have access to the BGP message definition (proto file), gRPC can still represent the fields with protobuf, but it has to add the name and value type for each one, so the receiving end can parse them. This is what we’ll examine next.

### Protobuf self-describing messages

While self-describing messages in a way defeat the purpose of protobuf by sending unnecessary data, we’ve included an example here to contrast how you could parse a message in this scenario:

![Figure 8.6 – Protobuf self-describing telemetry message (JSON)](https://static.packt-cdn.com/products/9781800560925/graphics/image/B16971_08_05.jpg)

Figure 8.6 – Protobuf self-describing telemetry message (JSON)

The telemetry header is the same, but when you choose `gpbkv` as the encoding format, Cisco IOS XR sends the data in the `data_bgpkv` field instead:

```markup
func main() {
  for msg := range ch {
    message := new(telemetry.Telemetry)
    err := proto.Unmarshal(msg, message)
    /* ... <omitted for brevity > ... */
    b, err := json.Marshal(message.GetDataGpbkv())
    check(err)
    j := string(b)
    // https://go.dev/play/p/uyWenG-1Keu
    data := gjson.Get(
      j,
      "0.fields.0.fields.#(name==neighbor-address).ValueByType.StringValue",
    )
    fmt.Println("  Neighbor: ", data)
    data = gjson.Get(
      j,
      "0.fields.1.fields.#(name==connection-state).ValueByType.StringValue",
    )
    fmt.Println("  Connection state: ", data)
  }
}
```

At this point, what you have is a big JSON file you can navigate using a Go package of your preference. Here, we’ve used `gjson`. To test this program, you can rerun the same program we described earlier with an extra flag to enable the self-describing key-value messages:

```markup
ch08/grpc$ go run grpc -kvmode=true
```

While this method might seem less involved, not only do you compromise the performance benefits but also, by not knowing the Go data structures beforehand, it opens up room for bugs and typos, it prevents you from taking advantage of the auto-completion features of most IDEs, and it makes your code less explicit. All of that has a negative impact on code development and troubleshooting.

Just Imagine

# Summary

In this chapter, we explored different ways to use APIs and RPCs to interact with network devices. One common theme we saw throughout this chapter was having a model for any data we exchange. Although the network community has embraced YANG as the standard language to model network configuration and operational state data, the implementation differences across networking vendors still impede its wide adoption.

In the next chapter, we’ll look at how OpenConfig tries to increase the adoption of declarative configuration and model-driven management and operations by defining a set of vendor-neutral models and protocols.

Just Imagine

# Further reading

-   The book’s GitHub repository: [https://github.com/ImagineDevOps DevOps/Network-Automation-with-Go](https://github.com/ImagineDevOps DevOps/Network-Automation-with-Go)
-   OpenAPI versions: [https://swagger.io/specification/#appendix-a-revision-history](https://swagger.io/specification/#appendix-a-revision-history)
-   CUE: [https://cuelang.org/](https://cuelang.org/)
-   `ch08/cue/template.cue`: [https://github.com/ImagineDevOps DevOps/Network-Automation-with-Go/blob/main/ch08/cue/template.cue](https://github.com/ImagineDevOps DevOps/Network-Automation-with-Go/blob/main/ch08/cue/template.cue)
-   CUE’s _References and Visibility_ tutorial: [https://cuelang.org/docs/tutorials/tour/references/](https://cuelang.org/docs/tutorials/tour/references/)
-   The `ch08/cue` directory: [https://github.com/ImagineDevOps DevOps/Network-Automation-with-Go/blob/main/ch08/cue](https://github.com/ImagineDevOps DevOps/Network-Automation-with-Go/blob/main/ch08/cue)
-   `ch08/cue/cue_tool.cue`: [https://github.com/ImagineDevOps DevOps/Network-Automation-with-Go/blob/main/ch08/cue/cue\_tool.cue](https://github.com/ImagineDevOps DevOps/Network-Automation-with-Go/blob/main/ch08/cue/cue_tool.cue)
-   Istio: [https://istio.io/](https://istio.io/)
-   dagger.io: [https://dagger.io/](https://dagger.io/)
-   Jsonnet: [https://github.com/google/go-jsonnet](https://github.com/google/go-jsonnet)
-   Dhall: [https://github.com/philandstuff/dhall-golang](https://github.com/philandstuff/dhall-golang)
-   JSON-RPC: [https://documentation.nokia.com/srlinux/SR\_Linux\_HTML\_R21-11/SysMgmt\_Guide/json-interface.html](https://documentation.nokia.com/srlinux/SR_Linux_HTML_R21-11/SysMgmt_Guide/json-interface.html)
-   openconfig/ygot: https://github.com/openconfig/ygot
-   Nokia’s YANG models: [https://github.com/nokia/srlinux-yang-models](https://github.com/nokia/srlinux-yang-models)
-   The YANG browser: [https://yang.srlinux.dev/v21.6.4/](https://yang.srlinux.dev/v21.6.4/)
-   ygot’s official documentation: [https://github.com/openconfig/ygot#introduction](https://github.com/openconfig/ygot#introduction)
-   The `ch08/json-rpc/main.go` file: [https://github.com/ImagineDevOps DevOps/Network-Automation-with-Go/blob/main/ch08/json-rpc/main.go](https://github.com/ImagineDevOps DevOps/Network-Automation-with-Go/blob/main/ch08/json-rpc/main.go)
-   The `ch08/json-rpc` directory: [https://github.com/ImagineDevOps DevOps/Network-Automation-with-Go/blob/main/ch08/json-rpc](https://github.com/ImagineDevOps DevOps/Network-Automation-with-Go/blob/main/ch08/json-rpc)
-   The `yget` library design documentation: [https://github.com/openconfig/ygot/blob/master/docs/design.md#openconfig-path-compression](https://github.com/openconfig/ygot/blob/master/docs/design.md#openconfig-path-compression)
-   Arista YANG models: [https://github.com/aristanetworks/yang](https://github.com/aristanetworks/yang)
-   The `ch08/restconf` directory: [https://github.com/ImagineDevOps DevOps/Network-Automation-with-Go/tree/main/ch08/restconf](https://github.com/ImagineDevOps DevOps/Network-Automation-with-Go/tree/main/ch08/restconf)
-   The `ch08/state` directory: [https://github.com/ImagineDevOps DevOps/Network-Automation-with-Go/tree/main/ch08/state](https://github.com/ImagineDevOps DevOps/Network-Automation-with-Go/tree/main/ch08/state)
-   IETF draft: [https://datatracker.ietf.org/doc/html/draft-openconfig-netmod-opstate-01](https://datatracker.ietf.org/doc/html/draft-openconfig-netmod-opstate-01)
-   The `ch08/state` program: [https://github.com/ImagineDevOps DevOps/Network-Automation-with-Go/tree/main/ch08/state](https://github.com/ImagineDevOps DevOps/Network-Automation-with-Go/tree/main/ch08/state)
-   `ch08/state/go.mod`: [https://github.com/ImagineDevOps DevOps/Network-Automation-with-Go/blob/main/ch08/state/go.mod](https://github.com/ImagineDevOps DevOps/Network-Automation-with-Go/blob/main/ch08/state/go.mod)
-   The `ch08/protobuf/write` directory: [https://github.com/ImagineDevOps DevOps/Network-Automation-with-Go/tree/main/ch08/protobuf/write](https://github.com/ImagineDevOps DevOps/Network-Automation-with-Go/tree/main/ch08/protobuf/write)
-   `ems_grpc.proto`: [https://github.com/nleiva/xrgrpc/blob/master/proto/ems/ems\_grpc.proto](https://github.com/nleiva/xrgrpc/blob/master/proto/ems/ems_grpc.proto)
-   `telemetry.proto`: [https://github.com/Juniper/jtimon/blob/master/telemetry/telemetry.proto](https://github.com/Juniper/jtimon/blob/master/telemetry/telemetry.proto)
-   `gnmi.proto`: [https://github.com/openconfig/gnmi/blob/master/proto/gnmi/gnmi.proto](https://github.com/openconfig/gnmi/blob/master/proto/gnmi/gnmi.proto)
-   `ch08/grpc/proto`: [https://github.com/ImagineDevOps DevOps/Network-Automation-with-Go/tree/main/ch08/grpc/proto](https://github.com/ImagineDevOps DevOps/Network-Automation-with-Go/tree/main/ch08/grpc/proto)
-   `ch08/grpc/generate_code`: [https://github.com/ImagineDevOps DevOps/Network-Automation-with-Go/blob/main/ch08/grpc/generate\_code](https://github.com/ImagineDevOps DevOps/Network-Automation-with-Go/blob/main/ch08/grpc/generate_code)