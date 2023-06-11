# Automation Frameworks

Most engineers start their automation journey by writing small ad hoc scripts. Over time, as these scripts grow in size and number, we need to think about the operating model for the solutions we create and how strong the foundations we are building upon are. Ultimately, we have to coordinate automation practices across different teams to generate business outcomes at scale.

To reduce the time and effort spent automating their use cases, some organizations try to standardize their tools and reuse generic components in their solutions, which often leads them to automation frameworks.

Automation frameworks allow different teams to come together under the same umbrella, break silos that may lead to inefficiencies, embrace common practices and code reusability, and enforce policies across domains to make the developed solutions more secure.

When choosing what best fits your environment and use cases, make sure you evaluate different automation frameworks. In this chapter, we will review some of them and focus specifically on how they can integrate with Go. In particular, we will look at the following:

-   How Go programs can become Ansible modules
-   The development of a custom Terraform provider
-   An overview of the rest of the well-known Go-based frameworks

We close this chapter by looking at the current trends in the industry and how the new generation of automation frameworks may develop in the future.

Just Imagine

# Technical requirements

You can find the code examples for this chapter in the book’s GitHub repository (see the _Further reading_ section), in the `ch07` folder.

Important Note

We recommend you execute the Go programs in this chapter in a virtual lab environment. Refer to the appendix for the prerequisites and instructions on how to build it.

Just Imagine

# Ansible

Ansible is an open source project, framework, and automation platform. Its descriptive automation language has captured the attention of many network engineers who see it as an introduction with minimal friction into the world of network automation and something that can help them become productive relatively quickly.

Ansible has an agentless push-based architecture. It connects to the hosts it manages via SSH and runs a series of tasks. These tasks are small programs that we call Ansible modules, which are the units of code that Ansible abstracts away from the user. A user only has to give the input arguments and can rely on Ansible modules to do all the heavy work for them. Although the level of abstraction may vary, Ansible modules allow users to focus more on the desired state of their infrastructure and less on the individual commands required to achieve that state.

## Overview of Ansible components

Playbooks are at the core of Ansible. These text-based declarative YAML files define a set of automation tasks that you can group in different plays. Each task runs a module that comes from either the Ansible code base or a third-party content collection:

![Figure 7.1 – Ansible high-level diagram](https://static.packt-cdn.com/products/9781800560925/graphics/image/B16971_07_01.jpg)

Figure 7.1 – Ansible high-level diagram

We use an Ansible inventory to describe the hosts or network devices we want to manage with Ansible. _Figure 7__.1_ provides a high-level overview of these elements.

### Inventory

An inventory is a list of managed hosts you can define statically in a text file or pull dynamically from an external system. You can manage hosts individually or collectively using groups. The following code snippet shows an Ansible inventory file:

```markup
[eos]
clab-netgo-ceos
[eos:vars]
ansible_user=admin
ansible_password=admin
ansible_connection=ansible.netcommon.network_cli
```

You can also use inventory to define group- and host-level variables that become available to Ansible playbooks.

### Playbooks, plays, and tasks

Ansible playbooks are files that you write using a YAML-based **Domain-Specific Language** (**DSL**). A playbook can have one or more plays on it. Each Ansible play targets a host or a group of hosts from an inventory to perform a series of tasks in a specific order. The following code output shows an example of a playbook with a single play and two tasks:

```markup
- name: First Play - Configure Routers
  hosts: routers
  gather_facts: true
  tasks:
    - name: Run Nokia Go module on local system with Go
      go_srl:
        host: "{{ inventory_hostname }}"
        user: "{{ ansible_user }}"
        password: "{{ ansible_password }}"
        input: "{{ hostvars[inventory_hostname] | string | b64encode }}"
      delegate_to: localhost
      when: ('srl' in group_names)
    - name: Run NVIDIA compiled Go module on remote system without Go
      go_cvx:
        host: localhost
        user: "{{ ansible_user }}"
        password: "{{ ansible_password }}"
        input: "{{ hostvars[inventory_hostname] | string | b64encode }}"
      when: ('cvx' in group_names)
```

The last example is a snippet from a larger playbook (see _Further reading_) included in the `ch07/ansible` folder of this book’s GitHub repository. That playbook has four tasks spread across two different plays. We use that playbook to review different concepts throughout this section.

### Modules

Each task executes an Ansible module. Although implementations may vary, the goal of an Ansible module is to be idempotent, so no matter how many times you run it against the same set of hosts, you always get the same outcome.

Ansible ships with several modules written mostly in Python, but it doesn’t stop you from using another programming language, which is what we explore in this section.

## Working with Ansible modules

The code of an Ansible module can execute either on a remote node, for hosts such as Linux servers, or locally, on the node running the playbook. The latter is what we typically do when the managed node is an API service or a network device because they both lack an execution environment with dependencies such as Linux shell and Python. Luckily, modern network operating systems meet those requirements, which give us both options of running the code locally or remotely.

If you look at the preceding playbook snippet, you can see how we implemented these two options. The first task invokes the `go_srl` module that gets delegated to the localhost. This means it runs from the machine running Ansible and targets a remote host provided in the host argument. The second task executes the `go_cvx` module, which is not delegated and thus runs on a remote node, targeting its API calls at the localhost.

The rest of the playbook uses a combination of local and remote execution environments, as denoted by the gear symbols in the following diagram:

![Figure 7.2 – Playbook example](https://static.packt-cdn.com/products/9781800560925/graphics/image/B16971_07_02.jpg)

Figure 7.2 – Playbook example

The Ansible playbook first runs an Ansible play to configure each node of the topology with these high-level objectives:

-   Configure the SR Linux node (`srl`) using a compiled Go code we execute locally on the machine running Ansible
-   Configure the NVIDIA Cumulus node (`cvx`) using a compiled Go code we execute on the remote node
-   Configure the Arista EOS node (`ceos`) using a compiled Go code we execute locally on the machine running Ansible

The choice of local or remote execution environments in the preceding playbook is random and only serves to show the two different approaches. Since all our lab devices are Linux-based, we can change this behavior without reworking the Ansible modules we use.

The second play has a single task that verifies the configured state on all three devices using a non-compiled code we execute using the `go run` command. We use this last task to show an alternative approach to concurrency that uses Go native primitives instead of Ansible forks to execute tasks on several nodes at the same time. We discuss this later in this section.

## Developing an Ansible module

While Ansible developers write most Ansible modules in Python, there are different reasons to write a module in another programming language:

-   Your company might use another programming language already.
-   Maybe you know or feel more comfortable writing in a different language.
-   The code is already available and there is no business justification to rewrite it in another programming language.
-   You want to take advantage of a feature that is not available in Python.

Ansible’s role is not to rip and replace everything that you have, especially if it’s working for you already. To illustrate this, we will take a set of Go programs from other chapters and turn them into Ansible modules we can execute in a playbook to configure our lab topology.

### Ansible module interface

You can extend Ansible by adding custom modules. Their implementation code should go into the `library` folder. When Ansible runs into a task with a module that is not installed in the system, it looks for a file with the module’s name in the `library` folder and tries to run it as a module, going through the following sequence of steps:

1.  It saves all module arguments in a temporary file, for example, `/tmp/foo`.
2.  It executes that module as a child process, passing it the filename as the first and only argument, for example, `./``library/my_module /tmp/foo`.
3.  It waits for the process to complete and expects to receive a structured response in its standard output.

While Ansible always expects a response in a JSON format, the input file format Ansible passes to the module depends on whether the module is a script or a binary. All binary modules get their input arguments as a JSON file, while script modules receive their input arguments as Bash files or just a list of key-value pairs.

From Go’s code perspective, to make this input behavior uniform, we normalize the input format to JSON before running any non-compiled Go programs. We do this using a wrapper Bash script that transforms the Bash input into JSON before calling the `go run` command, as you can see in the `ch07/ansible/library/go_state` file of this book’s GitHub repository (see _Further reading_).

### Adapting your Go code to interact with Ansible

Ultimately, a custom Ansible module can do anything as long as it understands how to parse the input arguments and knows how to return the expected output. We would need to change the Go programs from other chapters to make them an Ansible module. But the amount of changes necessary is minimal. Let’s examine this.

First, for this example, we need to create a struct to parse the module arguments we receive in the input JSON file. These arguments include login credentials and the input data model:

```markup
// ModuleArgs are the module inputs
type ModuleArgs struct {
  Host     string
  User     string
  Password string
  Input    string
}
func main() {
  if len(os.Args) != 2 {
    // generate error
  }
  argsFile := os.Args[1]
  text, err := os.ReadFile(argsFile)
  // check error
  var moduleArgs ModuleArgs
  err = json.Unmarshal(text, &moduleArgs)
  // check error
  /* ... <continues next > ... */
```

The input data model we use for Ansible remains the same as the one that we used in other chapters. This data is in the `ch07/ansible/host_vars` directory for this example. With Ansible, this data model becomes just a subset of all variables defined for each host. We pass it, along with the rest of the host variables, as a base64-encoded string. Inside our module, we decode the input string and decode it into the same `Model` struct we used before:

```markup
import (
  "encoding/base64"
  "gopkg.in/yaml.v2"
)
type Model struct {
  Uplinks  []Link `yaml:"uplinks"`
  Peers    []Peer `yaml:"peers"`
  ASN      int    `yaml:"asn"`
  Loopback Addr   `yaml:"loopback"`
}
func main() {
  /* ... <continues from before > ... */
  src, err :=
      base64.StdEncoding.DecodeString(moduleArgs.Input)
  // check error
  reader := bytes.NewReader(src)
  d := yaml.NewDecoder(reader)
  var input Model
  d.Decode(&input)
  /* ... <continues next > ... */
```

At this point, we’ve parsed enough information for our Go program to configure a network device. This part of the Go code does not require any modifications. The only thing you need to be mindful of is that instead of logging to the console, you now need to send any log messages as a response to Ansible.

When all the work is complete, we need to prepare and print the response object for Ansible. The following code snippet shows the _happy path_ when all changes have gone through:

```markup
// Response is the values returned from the module
type Response struct {
  Msg     string `json:"msg"`
  Busy    bool   `json:"busy"`
  Changed bool   `json:"changed"`
  Failed  bool   `json:"failed"`
}
func main() {
  /* ... <continues from before > ... */
  var r Response
  r.Msg = "Device Configured Successfully"
  r.Changed = true
  r.Failed = false
  response, err = json.Marshal(r)
  // check error
  fmt.Println(string(response))
  os.Exit(0)
}
```

Using a similar pattern to what we just described, we have created a custom module for each one of the three lab devices and one module to verify the state of the lab topology as we did in [_Chapter 6_](https://subscription.imaginedevops.io/book/cloud-and-networking/9781800560925/2B16971_06.xhtml#_idTextAnchor144), _Configuration Management_. You can find these modules in the `ch07/ansible/`{`srl`|`cvx`|`ceos`|`state`} directories of this book’s GitHub repository (see _Further reading_).

Before we move on to the execution, we want to show one way we can make use of Go’s built-in features to speed up and optimize concurrent task execution in Ansible.

### Taking advantage of Go’s concurrency

Ansible’s default behavior is to run each task on all hosts before moving on to the next one (linear strategy). Of course, it doesn’t just run one task on one host at a time; instead, it uses several independent processes attempting to run simultaneously on as many hosts as the number of forks you define in the Ansible configuration. Whether these processes run in parallel depends on the hardware resources available to them.

A less expensive approach from a resource utilization perspective is to leverage Go concurrency. This is what we do in the `go_state` Ansible module, where we target a single node from the inventory, the implicit localhost, and leave the concurrent communication with the remote nodes to Go.

For the following module, we reuse the code example from the _State validation_ section of [_Chapter 6_](https://subscription.imaginedevops.io/book/cloud-and-networking/9781800560925/2B16971_06.xhtml#_idTextAnchor144), _Configuration Management_ that has the access details embedded in the code already, but you could also pass these access details as arguments to the module to achieve the same result:

```markup
  - name: Run Validate module on Systems with Go installed
    go_state:
      host: "{{ inventory_hostname }}"
```

The trade - off of this approach is that we gain speed and get more efficient use of resources, but we lose the inventory management side of Ansible. Be mindful of this when trying to decide whether this is the right fit for your use case.

## Running the playbook

You can find the complete example involving four Go Ansible modules in the `ch07/ansible` directory. To run it, first make sure the lab topology is running from the root folder of the repository with `make lab-up`, then run the playbook with the `ansible-playbook` command:

```markup
ch07/ansible$ ansible-playbook playbook.yml 
# output omitted for brevity.
PLAY RECAP *********************************************************************************************************************************************************
clab-netgo-ceos            : ok=5    changed=0    unreachable=0    failed=0    skipped=4    rescued=0    ignored=0   
clab-netgo-cvx             : ok=2    changed=1    unreachable=0    failed=0    skipped=7    rescued=0    ignored=0   
clab-netgo-srl             : ok=2    changed=1    unreachable=0    failed=0    skipped=7    rescued=0    ignored=0   
localhost                  : ok=1    changed=0    unreachable=0    failed=0    skipped=0    rescued=0    ignored=0
```

Now that we’ve covered how Go programs can integrate with Ansible, we will move on to another popular automation framework: Terraform.

Just Imagine

# Terraform

Terraform is an open source software solution for declarative infrastructure management. It allows you to express and manage the desired state of your infrastructure with code. It has gained initial popularity as a framework to automate public cloud infrastructure but now supports a variety of on-premises and public cloud resources, platforms, services—almost anything that has an API.

One of the key distinctions of Terraform is the way it manages state. Once it creates a remote resource initially, it saves the resulting state in a file and relies on that state to be there for its next runs. As you update and develop your infrastructure code, the state file enables Terraform to manage the entire life cycle of a remote resource, calculating the precise sequence of API calls to transition between states. This ability to manage state and the declarative configuration language and the agentless, API-first architecture allowed Terraform to become deeply entrenched in the cloud infrastructure space and become a critical part of DevOps and Infrastructure-as-Code toolchains.

If we look at the Terraform registry (see _Further reading_), we can see over a hundred providers in the networking category ranging from SDN appliances and firewalls to various cloud services. This number is on a rising trend, as more people adopt a declarative approach to manage their infrastructure as code. This is why we believe it’s important for network automation engineers to know Terraform and be able to extend its capabilities using Go.

## Overview of Terraform components

The entire Terraform ecosystem is a collection of Go packages. They distribute the main CLI tool, often referred to as _Terraform Core_, as a statically compiled binary. This binary implements the command-line interface and can parse and evaluate instructions written in **Hashicorp Configuration Language** (**HCL**). On every invocation, it builds a resource graph and generates an execution plan to reach the desired state described in the configuration file. The main binary only includes a few plugins but can discover and download the required dependencies.

Terraform plugins are also distributed as standalone binaries. Terraform Core starts and terminates the required plugins as child processes and interacts with them using an internal gRPC-based protocol. Terraform defines two types of plugins:

-   **Providers**: Interact with a remote infrastructure provider and implement the required changes
-   **Provisioners**: Implement a set of imperative actions, declared as a set of terminal commands, to bootstrap a resource that a provider created before

The following diagram demonstrates what we have described and shows how different Terraform components communicate internally and externally:

![Figure 7.3 – Terraform high-level diagram](https://static.packt-cdn.com/products/9781800560925/graphics/image/B16971_07_03.jpg)

Figure 7.3 – Terraform high-level diagram

The vast majority of Terraform plugins are providers as they implement the declarative resource actuation and communicate with an upstream API. A provider defines two types of objects that you can use to interact with a remote API:

-   **Resources**: Represent the actual managed infrastructure objects, such as virtual machines, firewall policies, and DNS records
-   **Data Sources**: Offer a way to query information that is not managed by Terraform, such as a list of supported cloud regions, VM images, or **Identity and Access Management** (**IAM**) roles

It’s up to the Terraform provider maintainers to decide what resources and data sources to implement, so the coverage may vary, especially between official and community-supported providers.

## Working with Terraform

A typical Terraform workflow involves several stages that need to happen in sequence. We first need to define a provider that determines what infrastructure we would manage, and then describe the state of our infrastructure using a combination of resources and data sources. We will walk through these stages by following a configuration file, `ch07/terraform/main.tf`, we’ve created in this book’s GitHub repository (see _Further reading_).

### Defining a provider

Providers define connection details for the upstream API. They can point at the public AWS API URL or an address of a private vCenter instance. In the next example, we show how to manage the demo instance of Nautobot running at [https://demo.nautobot.com/](https://demo.nautobot.com/).

Terraform expects to find a list of required providers, along with their definition, in one file in the current working directory. For the sake of simplicity, we include those details at the top of the `main.tf` file and define credentials in the same file. In production environments, these details may live in a separate file, and you should source credentials externally, for example, from environment variables:

```markup
terraform {
  required_providers {
    nautobot = {
      version = "0.2.4"
      source  = "nleiva/nautobot"
    }
  }
}
provider "nautobot" {
  url = "https://demo.nautobot.com/api/"
  token = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
}
```

With this information defined, we can initialize Terraform. The following command instructs Terraform to perform plugin discovery and download any dependencies into a local `./``terraform` directory:

```markup
ch07/terraform$ terraform init -upgrade
Initializing the backend...
Initializing provider plugins...
- Finding nleiva/nautobot versions matching "0.2.4"...
- Installing nleiva/nautobot v0.2.4...
- Installed nleiva/nautobot v0.2.4 (self-signed, key ID A33D26E300F155FF)
```

At the end of this step, Terraform creates a lock file, `.terraform.lock.hcl`, to record the provider selections it just made. Include this file in your version control repository so that Terraform can guarantee to make the same selections by default when you run `terraform init` on a different machine.

### Creating a resource

To create a resource, we define it in a configuration block with zero or more arguments that assign values to resource fields. The following resource creates a new `Manufacturer` object in Nautobot with the specified name and description:

```markup
resource "nautobot_manufacturer" "new" {
  description = "Created with Terraform"
  name        = "New Vendor"
}
```

Now we can run `terraform plan` to check whether the current configuration matches the existing state. If they don’t match, Terraform creates an execution plan with the proposed changes to make the remote objects match the current configuration. We could skip the `terraform plan` command and move straight to `terraform apply`, which generates the plan and also executes it in a single step:

```markup
ch07/terraform$ terraform apply --auto-approve
Terraform used the selected providers to generate the following execution plan. Resource actions
are indicated with the following symbols:
  + create
Terraform will perform the following actions:
  # nautobot_manufacturer.new will be created
  + resource "nautobot_manufacturer" "new" {
      + created             = (known after apply)
      + description         = "Created with Terraform"
      + devicetype_count    = (known after apply)
      + display             = (known after apply)
      + id                  = (known after apply)
      + inventoryitem_count = (known after apply)
      + last_updated        = (known after apply)
      + name                = "New Vendor"
      + platform_count      = (known after apply)
      + slug                = (known after apply)
      + url                 = (known after apply)
    }
Plan: 1 to add, 0 to change, 0 to destroy.
```

You can see the result of running this plan in Nautobot’s web UI at [https://demo.nautobot.com/dcim/manufacturers/new-vendor/](https://demo.nautobot.com/dcim/manufacturers/new-vendor/), or you can check the resulting state using the following command:

```markup
ch07/terraform$ terraform state show 'nautobot_manufacturer.new'
# nautobot_manufacturer.new:
resource "nautobot_manufacturer" "new" {
    created             = "2022-05-04"
    description         = "Created with Terraform"
    devicetype_count    = 0
    display             = "New Vendor"
    id                  = "09219670-3e28-..."
    inventoryitem_count = 0
    last_updated        = "2022-05-04T18:29:06.241771Z"
    name                = "New Vendor"
    platform_count      = 0
    slug                = "new-vendor"
    url                 = "https://demo.nautobot.com/api/dcim/manufacturers/09219670-3e28-.../"
}
```

At the time of writing, there was no Terraform provider available for Nautobot, so the last example used a custom provider we created specifically for this book. Creating a new provider can enable many new use cases and it involves writing Go code, so this is what we cover next.

## Developing a Terraform provider

Eventually, you may come across a provider with limited or missing capabilities, or a provider may not even exist for a platform that is part of your infrastructure. This is when knowing how to build a provider can make a difference, to either extend or fix a provider or build a brand new one. The only prerequisite to get started is the availability of a Go SDK for the target platform. For example, Nautobot has a Go client package that gets automatically generated from its OpenAPI model, which we used already in the _Getting config inputs from other systems via HTTP_ section of [_Chapter 6_](https://subscription.imaginedevops.io/book/cloud-and-networking/9781800560925/2B16971_06.xhtml#_idTextAnchor144), _Configuration Management_, so we have all we need to develop its Terraform provider.

The recommended way to create a new Terraform provider is to start with the terraform-provider-scaffolding project (see _Further reading_). This repository provides enough boilerplate to allow you to focus on the internal logic while it provides function stubs and implements **Remote Procedure Call** (**RPC**) integration. We used this template to create the Nautobot provider, so you can compare our final result with the template to see what changes we made.

As a by-product of developing a Terraform provider using the scaffolding project, you can register your Git repository in the Terraform registry and get the benefit of automatically rendered provider documentation (see _Further reading_).

### Defining a provider

The provider’s internal code (`internal/provider/provider.go` (see _Further reading_)) starts with a schema definition for the provider itself as well as its managed resources and data sources. Inside the provider’s schema, we define two input arguments—`url` and `token`. You can extend each schema struct with more constraints, default values, and validation functions:

```markup
func New(version string) func() *schema.Provider {
  return func() *schema.Provider {
    p := &schema.Provider{
      Schema: map[string]*schema.Schema{
        "url": {
          Type:         schema.TypeString,
          Required:     true,
          DefaultFunc:
          schema.EnvDefaultFunc("NAUTOBOT_URL", nil),
          ValidateFunc: validation.IsURLWithHTTPorHTTPS,
          Description:  "Nautobot API URL",
        },
        "token": {
          Type:        schema.TypeString,
          Required:    true,
          Sensitive:   true,
          DefaultFunc:
            schema.EnvDefaultFunc("NAUTOBOT_TOKEN", nil),
          Description: "Admin API token",
        },
      },
      DataSourcesMap: map[string]*schema.Resource{
        "nautobot_manufacturers":
            dataSourceManufacturers(),
      },
      ResourcesMap: map[string]*schema.Resource{
        "nautobot_manufacturer": resourceManufacturer(),
      },
    }
    p.ConfigureContextFunc = configure(version, p)
    return p
  }
}
```

With login information defined, the provider can initialize an API client for the target platform. This happens inside a local function where `url` and `token` get passed to the Nautobot’s Go SDK, which creates a fully authenticated HTTP client. We save this client in a special `apiClient` struct, which gets passed as an argument to all provider resources, as we show later on:

```markup
import nb "github.com/nautobot/go-nautobot"
type apiClient struct {
  Client *nb.ClientWithResponses
  Server string
}
func configure(
  version string,
  p *schema.Provider,
) func(context.Context, *schema.ResourceData) (interface{}, diag.Diagnostics) {
  return func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
    serverURL := d.Get("url").(string)
    _, hasToken := d.GetOk("token")
    /* ... <omitted for brevity > ... */
    token, _ :=
        NewSecurityProviderNautobotToken(
          d.Get("token").(string))
    c, err := nb.NewClientWithResponses(
              serverURL,
              nb.WithRequestEditorFn(token.Intercept),
            )
    // process error
    return &apiClient{
      Client: c,
      Server: serverURL,
    }, diags
  }
}
```

Now that we have prepared a remote API client, we can start writing code for our managed resources.

### Defining resources

Just like how we defined a schema for our provider, we now need to define a schema for each managed resource and data source. For educational purposes, we only implement a single resource type, `Manufacturer`, and a corresponding data source you can use to retrieve the list of all existing manufacturers in Nautobot.

When we define a schema, our goal is to match the upstream API as closely as possible. This should reduce the number of required data transformations and make the implementation work much easier. Let’s look at Nautobot’s Go SDK code:

```markup
type Manufacturer struct {
  Created       *openapi_types.Date
    `json:"created,omitempty"`
  CustomFields  *Manufacturer_CustomFields
    `json:"custom_fields,omitempty"`
  Description   *string `json:"description,omitempty"`
  /* ... <omitted for brevity > ... */
  Url           *string `json:"url,omitempty"`
}
type Manufacturer_CustomFields struct {
  AdditionalProperties map[string]interface{} `json:"-"`
}
```

The schema that we define for the `Manufacturer` resource in `resource_manufacturer.go` closely follows the fields and types defined in the preceding output:

```markup
func resourceManufacturer() *schema.Resource {
  return &schema.Resource{
    Description: "This object manages a manufacturer",
    CreateContext: resourceManufacturerCreate,
    ReadContext:   resourceManufacturerRead,
    UpdateContext: resourceManufacturerUpdate,
    DeleteContext: resourceManufacturerDelete,
    Schema: map[string]*schema.Schema{
      "created": {
        Description: "Manufacturer's creation date.",
        Type:        schema.TypeString,
        Computed:    true,
      },
      "description": {
        Description: "Manufacturer's description.",
        Type:        schema.TypeString,
        Optional:    true,
      },
      "custom_fields": {
        Description: "Manufacturer custom fields.",
        Type:        schema.TypeMap,
        Optional:    true,
      },
      /* ... <omitted for brevity > ... */
      "url": {
        Description: "Manufacturer's URL.",
        Type:        schema.TypeString,
        Optional:    true,
        Computed:    true,
      },
    },
  }
}
```

Once we have defined all schemas with their constraints, types, and descriptions, we can start implementing resource operations. The scaffolding project provides stubs for each one of the CRUD functions, so we only need to fill them out with code.

### The create operation

We first look at the `resourceManufacturerCreate` function, which gets invoked when Terraform determines that it must create a new object. This function has two very important arguments:

-   `meta`: Stores the API client we created earlier
-   `d`: Stores all resource arguments defined in the HCL configuration file

We extract the user-defined configuration from `d` and use it to build a new `nb.Manufacturer` object from the Nautobot’s SDK. We can then use the API client to send that object to Nautobot and save the returned object ID:

```markup
func resourceManufacturerCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
    c := meta.(*apiClient).Client
    var m nb.Manufacturer
    name, ok := d.GetOk("name")
    n := name.(string)
    if ok {
        m.Name = n
    }
    /* ... <omitted for brevity > ... */
    rsp, err := c.DcimManufacturersCreateWithResponse(
        ctx,
        nb.DcimManufacturersCreateJSONRequestBody(m))
    // process error
    // process returned HTTP response
    d.SetId(id.String())
    return resourceManufacturerRead(ctx, d, meta)
}
```

Typically, we don’t define all optional fields when we create a new object. A remote provider assigns the unique ID and initializes default values as it creates a new object. Some platforms return the newly created object back, but there is no guarantee of that. Hence, it’s a common pattern in Terraform provider implementations to call a read function at the end of the create function to synchronize and update a local state.

### The read operation

The read function updates the local state to reflect the latest state of an upstream resource. We’ve seen in the preceding example how the create function calls the read at the end of its execution to update the state of a newly created object.

But the most important use of read is to detect configuration drift. When you do `terraform plan` or `terraform apply`, read is the first thing that Terraform executes and its goal is to retrieve the current upstream state and compare it with the state file. This allows Terraform to understand whether users have manually changed a remote object, so it needs to reconcile its state, or whether it’s up to date and no updates are necessary.

Read has the same signature as the rest of the CRUD functions, which means it gets the latest version of a managed resource as `*schema.ResourceData` and an API client stored in `meta`. The first thing we need to do in this function is fetch the upstream object:

```markup
import "github.com/deepmap/oapi-codegen/pkg/types"
func resourceManufacturerRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
    c := meta.(*apiClient).Client
    id := d.Get("id").(string)
    rsp, err := c.DcimManufacturersListWithResponse(
        ctx,
        &nb.DcimManufacturersListParams{
            IdIe: &[]types.UUID{types.UUID(id)},
        })
  /* ... <continues next > ... */
}
```

We use the data we get back to update the local Terraform state:

```markup
func resourceManufacturerRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
    /* ... <continues from before > ... */
    d.Set("name", item["name"].(string))
    d.Set("created", item["created"].(string))
    d.Set("description", item["description"].(string))
    d.Set("display", item["display"].(string))
    /* ... <omitted for brevity > ... */
    return diags
}
```

At this stage, our local state should be in sync with the upstream and Terraform can decide whether any changes are necessary as a result.

### Remaining implementations

In this chapter, we only cover a subset of the Nautobot provider code. The remaining sections we need to implement include the following:

-   The resource **update** and **delete** functions
-   **Data** **source** implementation

For the sake of brevity, we don’t include this code in the book, but the full implementation for the `Manufacturer` resource and data source is available in our demo Nautobot provider repository (see _Further reading_).

## Networking providers

Writing a provider and keeping it up to date is a major undertaking. At the beginning of this section, we mentioned that Terraform has several providers in the networking category of the Terraform registry (see _Further reading_). We invite you to explore them and always check whether there’s an existing provider before implementing your own.

Terraform’s guarantees of declarative configuration and state management are very appealing to network engineers trying to adopt DevOps and GitOps practices. As the interest grows, so does the number of new network-related providers, with the following notable recent additions:

-   **JUNOS Terraform Automation Framework** (see _Further reading_): Allows you to create a custom JunOS Terraform provider from YANG files
-   **Terraform Provider for Cisco IOS XE** (see _Further reading_): Manages the configuration of Cisco Catalyst IOS XE devices including switches, routers, and wireless LAN controllers
-   **terraform-provider-junos** (see _Further reading_): An unofficial Terraform provider for Junos OS devices with the NETCONF protocol
-   **terraform-provider-ciscoasa** (see _Further reading_): DevNet provider to configure Cisco ASA firewall rules

This completes the overview of Terraform and its network-related use cases. We hope that its adoption continues to increase and the number of networking providers grows. In the next section, we wrap up with a brief overview of a few other automation frameworks.

Just Imagine

# Other automation frameworks

Our industry has many more automation frameworks and solutions that we would have liked to cover in this chapter. The best we can do is just scratch the surface, leaving much of the exploration up to you. At the same time, we don’t want to leave you thinking there’s nothing out there besides Ansible and Terraform. This section gives you an overview of other automation frameworks and solutions that you can use or adapt to use within a networking context.

## Gornir

Nornir (see _Further reading_) is a popular network automation framework for Python that offers a pure programming experience by ditching DSL in favor of the Python API. It has a pluggable architecture where you can replace or extend almost any element of the framework, from inventory to device connections. It also has a flexible way to parallelize groups of tasks without having to deal with Python’s concurrency primitives directly.

Gornir (see _Further reading_) is a Nornir implementation in Go. Keeping with the same principles, it offers things such as inventory management, concurrent execution of tasks, and pluggable connection drivers. Gornir ships with a minimal set of drivers, but its core provides Go interfaces to improve upon and extend this feature. If you’re coming to Go from Python and are familiar with Nornir, Gornir may offer a very smooth transition through a familiar API and workflows.

## Consul-Terraform-Sync

In the preceding section, we examined how you can use Terraform to manage resources declaratively on a remote target, using Nautobot as an example. Hashicorp, the same company behind Terraform, has developed another automation solution that builds on top of it. It’s called Consul-Terraform-Sync (see _Further reading_) and it enables automatic infrastructure management by combining Terraform with Consul and linking them together with a synchronization agent.

Consul is a distributed key/value store used for service discovery, load balancing, and access control. It works by setting up a cluster of nodes that use the Raft consensus protocol to have a consistent view of their internal state. Server nodes communicate with their clients and broadcast relevant updates to make sure clients have an up-to-date version of the relevant part of the internal state. All this happens behind the scenes, with minimal configuration, which makes Consul a very popular choice for service discovery and data storage.

The main idea of the Consul-Terraform-Sync solution is to use Consul as a backend for Terraform configuration and state. The synchronization agent connects to Consul, waits for updates, and automatically triggers Terraform reconciliation as it detects any changes.

Consul-Terraform-Sync allows you to automate Terraform deployments for any of these providers and ensures that your state always matches your intent thanks to the automated reconciliation process.

## mgmt

`mgmt` (see _Further reading_) is another infrastructure automation and management framework written completely in Go. It has its own DSL and synchronizes its state using a baked-in etcd cluster. It uses a few interesting ideas, such as a declarative and functional DSL, resource graphs, and dynamic state transitions triggered by closed-loop feedback. Just like Gornir, `mgmt` ships with a set of plugins that users can extend, but none of these plugins is specifically for network devices since the main use case for mgmt is Linux server management.

## Looking into the future

In this chapter, we have covered popular network automation frameworks in use today. All these frameworks are at a different stage of development—some have already reached their peak while others are still crossing the chasm (see _Further reading_). But it’s important to remember that automation frameworks are not a solved problem with well-established projects and well-understood workflows. This field is constantly developing, and new automation approaches are emerging on the horizon.

These alternative approaches do not resemble what we had seen before. One big trend that we’re seeing lately is the departure from an imperative automation paradigm, where a human operator manually triggers actions and tasks. We briefly discussed this trend in [_Chapter 5_](https://subscription.imaginedevops.io/book/cloud-and-networking/9781800560925/2B16971_05.xhtml#_idTextAnchor128), _Network Automation_, and we want to revisit it here to show how the _closed-loop_ automation approach changes the landscape of infrastructure management systems. Most modern automation frameworks develop into systems that exhibit some or all the following characteristics:

-   Focus on the complete life cycle management of a system as opposed to individual stages, such as bootstrapping, provisioning, or decommissioning.
-   Exclusive use of declarative state definition and automatic reconciliation, or self-healing implemented internally.
-   Separation of state definitions from the platform managing this state through practices such as GitOps.
-   Offer a cloud-native self-service experience via APIs, reducing the friction in consuming of these services both manually and programmatically.

We’re currently at a point when these systems and their building blocks are becoming a reality, with some notable examples including Crossplane, Nokia Edge Network Controller, and Anthos Config Sync. They build these systems as Kubernetes controllers, leveraging the Operator model, allowing them to expose their APIs in a standard way, so other systems can talk to them with the same set of tools. We still don’t know whether these systems could become mainstream and displace the incumbent frameworks, since they increase the level of complexity and they introduce a steep learning curve. Regardless of that, it’s an area to explore, like other potential new trends that might develop, since infrastructure management is far from being a solved problem.

Just Imagine

# Summary

Whether to choose Ansible, Terraform, or a programming language to solve a particular use case depends on many variables. But don’t fall into the trap of looking at this as a binary decision. Most times, different technologies complement each other to offer solutions, as we showed in this chapter. In the next chapter, we will explore newer and more advanced techniques to interact with networking devices and Go.

Just Imagine

# Further reading

-   This book’s GitHub repository: https://github.com/ImagineDevOps DevOps/Network-Automation-with-Go
-   Playbook: https://github.com/ImagineDevOps DevOps/Network-Automation-with-Go/blob/main/ch07/ansible/playbook.yml
-   Terraform registry: [https://registry.terraform.io/browse/providers?category=networking](https://registry.terraform.io/browse/providers?category=networking%20)
-   terraform-provider-scaffolding project: [https://github.com/hashicorp/terraform-provider-scaffolding](https://github.com/hashicorp/terraform-provider-scaffolding%20)
-   Provider documentation: [https://registry.terraform.io/providers/nleiva/nautobot/latest/docs?pollNotifications=true](https://registry.terraform.io/providers/nleiva/nautobot/latest/docs?pollNotifications=true%20)
-   Provider’s internal code: [https://github.com/nleiva/terraform-provider-nautobot/blob/main/internal/provider/provider.go](https://github.com/nleiva/terraform-provider-nautobot/blob/main/internal/provider/provider.go%20)
-   `resource_manufacturer.go`: [https://github.com/nleiva/terraform-provider-nautobot/blob/main/internal/provider/resource\_manufacturer.go](https://github.com/nleiva/terraform-provider-nautobot/blob/main/internal/provider/resource_manufacturer.go%20)
-   Nautobot provider repository: [https://github.com/nleiva/terraform-provider-nautobot](https://github.com/nleiva/terraform-provider-nautobot%20)
-   JUNOS Terraform Automation Framework: [https://github.com/Juniper/junos-terraform](https://github.com/Juniper/junos-terraform%20)
-   Terraform Provider for Cisco IOS XE: [https://github.com/CiscoDevNet/terraform-provider-iosxe](https://github.com/CiscoDevNet/terraform-provider-iosxe%20)
-   terraform-provider-junos: [https://github.com/jeremmfr/terraform-provider-junos](https://github.com/jeremmfr/terraform-provider-junos%20)
-   terraform-provider-ciscoasa: https://github.com/CiscoDevNet/terraform-provider-ciscoasa
-   Nornir: [https://github.com/nornir-automation/nornir/](https://github.com/nornir-automation/nornir/%20)
-   Gornir: [https://github.com/nornir-automation/gornir](https://github.com/nornir-automation/gornir%20)
-   Consul-Terraform-Sync: [https://learn.hashicorp.com/tutorials/consul/consul-terraform-sync-intro?in=consul/network-infrastructure-automation](https://learn.hashicorp.com/tutorials/consul/consul-terraform-sync-intro?in=consul/network-infrastructure-automation%20)
-   `mgmt`: [https://github.com/purpleidea/mgmt](https://github.com/purpleidea/mgmt)
-   [https://en.wikipedia.org/wiki/Diffusion\_of\_innovations](https://en.wikipedia.org/wiki/Diffusion_of_innovations)