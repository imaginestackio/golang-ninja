# Network Monitoring

Despite the popularity of configuration management, we actually spend more time monitoring networks than configuring them. As networks become more and more complex, with new layers of encapsulation and IP address translations, our ability to understand whether a network functions correctly to let us meet customer **service-level agreements** (**SLAs**) is becoming increasingly difficult.

Engineers working in the cloud infrastructure space have come up with the term _observability_, referring to the ability to reason about the internal state of a system by observing its external outputs. Translated into networking terms, this may include passive monitoring through logs and state telemetry collection or active monitoring using distributed probing, data processing, and visualization.

The ultimate goal of all this is to reduce the **mean time to repair** (**MTTR**), adhere to customer SLAs, and shift to proactive problem resolution. Go is a very popular language of choice for these kinds of tasks, and in this chapter we will examine a few of the tools, packages, and platforms that can help you with network monitoring. Here are the highlights of this chapter:

-   We will explore traffic monitoring by looking at how to capture and parse network packets with Go.
-   Next, we will look at how to process and aggregate data plane telemetry to get meaningful insights into the current network behavior.
-   We show how you can use active probing to measure network performance, and how to produce, collect, and visualize performance metrics.

We will deliberately avoid talking about YANG-based telemetry, as we covered this already in [_Chapter 8_](https://subscription.imaginedevops.io/book/cloud-and-networking/9781800560925/2B16971_08.xhtml#_idTextAnchor182), _Network APIs_, and [_Chapter 9_](https://subscription.imaginedevops.io/book/cloud-and-networking/9781800560925/2B16971_09.xhtml#_idTextAnchor209), _OpenConfig_.

Another area that we haven’t touched on so far and that we want to discuss briefly in this chapter is the developer experience. As we write more code, maintaining existing software becomes an important part of our day-to-day operations. We introduce one tool per section of this chapter, acknowledging that we are just scratching the surface and that this topic could be the subject of an entire book. In the end, we don’t strive to give a comprehensive overview of all tools there are out there but just want to give you an idea of what developing Go code in production may feel like.

Just Imagine

# Technical requirements

You can find the code examples for this chapter in the book’s GitHub repository (see the _Further reading_ section), under the `ch10` folder.

Important Note

We recommend you execute the Go programs in this chapter in a virtual lab environment. Refer to the _Appendix_ for prerequisites and instructions on how to build the fully configured network topology.

The first example we will discuss in the following section explores packet capturing and parsing capabilities in Go.

Just Imagine

# Data plane telemetry processing

Network activities such as capacity planning, billing, or **distributed denial-of-service** (**DDoS**) attack monitoring require insights into the traffic flowing through a network. One way we can offer such visibility is by deploying a packet sampling technology. The premise is that at a high-enough rate, it’s possible to capture only a randomly sampled subset of packets to build a good understanding of the overall network traffic patterns.

While it’s the hardware that samples the packets, it’s the software that aggregates them into flows and exports them. NetFlow, sFlow, and **IP Flow Information Export** (**IPFIX**) are the three main protocols we use for this, and they define the structure of the payload and what metadata to include with each sampled packet.

One of the first steps in any telemetry processing pipeline is information ingestion. In our context, this means receiving and parsing data plane telemetry packets to extract and process flow records. In this section, we will look at how you can capture and process packets with the help of the `google/gopacket` package (see _Further reading_).

## Packet capturing

In [_Chapter 4_](https://subscription.imaginedevops.io/book/cloud-and-networking/9781800560925/2B16971_04.xhtml#_idTextAnchor109), _Networking (TCP/IP) with Go_, we discussed how to build a UDP ping application using the `net` package from Go’s standard library. And while we should probably take a similar approach when building an sFlow collector, we will do something different for the next example.

Instead of building a data plane telemetry collector, we designed our application to tap into an existing flow of telemetry packets, assuming the network devices in the topology are sending them to an existing collector somewhere in the network. This allows you to avoid changing the existing telemetry service configuration while still being able to capture and process telemetry traffic. You can use a program like this when you want a transparent tool that can run directly on a network device, on demand, and for a short period of time.

In the test lab topology, the `cvx` node runs an agent that exports sampled metrics using the sFlow protocol. The sFlow traffic flows toward `host-2`, where it gets intercepted by the example application using a tap:

![Figure 10.1 – sFlow application](https://static.packt-cdn.com/products/9781800560925/graphics/image/B16971_10_01.jpg)

Figure 10.1 – sFlow application

To show you the packet-capturing capabilities of the `google/gopacket` package, we intercept all sFlow packets using `pcapgo` – a native Go implementation of the traffic-capturing API in Linux. Although it’s less feature-rich than its counterpart `pcap` and `pfring` packages, the benefit of `pcapgo` is that it doesn’t rely on any external C libraries and can work natively on any Linux distribution.

In the first part of the `packet-capture` program, which you can find in the `ch10/packet-capture` folder of this book’s GitHub repository, we set up a new `af_packet` socket handler with the `pcapgo.NewEthernetHandle` function, passing it the name of the interface to monitor:

```markup
import (
     "github.com/google/gopacket/pcapgo"
)
var (
     intf = flag.String("intf", "eth0", "interface")
)
func main() {
     handle, err := pcapgo.NewEthernetHandle(*intf)
     /* ... <continues next > ... */
}
```

At this point, `handle` gives us access to all packets on the `eth0` interface.

## Packet filtering

While we could just capture all packets through the interface, for the sake of experimenting, we will include an example of how to filter the traffic we capture with a **Berkeley Packet Filter** (**BPF**) program in Go.

First, we generate a compiled packet-matching code in a human-readable format, using the `-d` option of the `tcpdump` command to filter IP and UDP packets:

```markup
$ sudo tcpdump -p -ni eth0 -d "ip and udp"
(000) ldh      [12]
(001) jeq      #0x800           jt 2    jf 5
(002) ldb      [23]
(003) jeq      #0x11            jt 4    jf 5
(004) ret      #262144
(005) ret      #0
```

Then, we convert each of the preceding instructions into a corresponding `bpf.Instruction` from the `golang.org/x/net/bpf` package. We assemble these instructions into a set of `[]bpf.RawInstruction` that are ready to load into a BPF virtual machine:

```markup
import (
  "golang.org/x/net/bpf"
)
 
func main() {
/* ... <continues from before > ... */
 
  rawInstructions, err := bpf.Assemble([]bpf.Instruction{
    // Load "EtherType" field from the ethernet header.
    bpf.LoadAbsolute{Off: 12, Size: 2},
    // Skip to last instruction if EtherType isn't IPv4.
    bpf.JumpIf{Cond: bpf.JumpNotEqual, Val: 0x800,
                    SkipTrue: 3},
    // Load "Protocol" field from the IPv4 header.
    bpf.LoadAbsolute{Off: 23, Size: 1},
    // Skip to the last instruction if Protocol is not UDP.
    bpf.JumpIf{Cond: bpf.JumpNotEqual, Val: 0x11,
                    SkipTrue: 1},
    // "send up to 4k of the packet to userspace."
    bpf.RetConstant{Val: 4096},
    // Verdict is "ignore packet and return to the stack."
    bpf.RetConstant{Val: 0},
  })
 
 
  handle.SetBPF(rawInstructions)
  /* ... <continues next > ... */
}
```

We can attach the result to the `EthernetHandle` function we created earlier, to act as a packet filter and reduce the number of packets received by the application.

In summary, we copy all packets that match the `0x800` EtherType and the `0x11` IP protocol to the user space process, where our Go program runs, while all the other packets, including the ones we match, continue through the network stack. This makes this program completely transparent to any existing traffic flows, and you can use it without having to change the configuration of the sFlow agent.

## Packet processing

All packets that the kernel sends to the user space become available in the Go application through the `PacketSource` type, which we build by combining the `EthernetHandle` function we created with an Ethernet packet decoder:

```markup
func main() {
  /* ... <continues from before > ... */
     packetSource := gopacket.NewPacketSource(
           handle,
           layers.LayerTypeEthernet,
     )
     /* ... <continues next > ... */
}
```

This `PacketSource` structure sends each received and decoded packet over a Go channel, which means we can use a `for` loop to iterate over them one by one. Inside this loop, we use `gopacket` to match packet layers and extract information about L2, L3, and L4 networking headers, including protocol-specific details such as the sFlow payload:

```markup
func main() {
  /* ... <continues from before > ... */
  for packet := range packetSource.Packets() {
    sflowLayer := packet.Layer(layers.LayerTypeSFlow)
    if sflowLayer != nil {
      sflow, ok := sflowLayer.(*layers.SFlowDatagram)
      if !ok {
        continue
      }
 
      for _, sample := range sflow.FlowSamples {
        for _, record := range sample.GetRecords() {
          p, ok := record.(layers.SFlowRawPacketFlowRecord)
          if !ok {
            log.Println("failed to decode sflow record")
            continue
          }
 
          srcIP, dstIP := p.Header.
            NetworkLayer().
            NetworkFlow().
            Endpoints()
          sPort, dPort := p.Header.
            TransportLayer().
            TransportFlow().
            Endpoints()
          log.Printf("flow record: %s:%s <-> %s:%s\n",
            srcIP,
            sPort,
            dstIP,
            dPort,
          )
        }
      }
     }
  }
}
```

The benefit of using `gopacket` specifically for sFlow decoding is that it can parse and create another `gopacket.Packet` based on the sampled packet’s headers.

## Generating traffic

To test this Go application, we need to generate some traffic in the lab topology, so the `cvx` device can generate sFlow records about it. Here, we use `microsoft/ethr` – a Go-based traffic generator that offers a user experience and features comparable to `iperf`. It can generate and receive a fixed volume of network traffic and measure bandwidth, latency, loss, and jitter. In our case, we only need it to generate a few low-volume traffic flows over the lab network to trigger the data plane flow sampling.

The `packet-capture` application taps into the existing sFlow traffic, parses and extracts flow records, and prints that information on the screen. To test the program, run `make capture-start` from the root of this book’s GitHub repository (see _Further reading_):

```markup
$ make capture-start
docker exec -d clab-netgo-cvx systemctl restart hsflowd
docker exec -d clab-netgo-host-3 ./ethr -s
docker exec -d clab-netgo-host-1 ./ethr -c 203.0.113.253 -b 900K -d 60s -p udp -l 1KB
docker exec -d clab-netgo-host-1 ./ethr -c 203.0.113.252 -b 600K -d 60s -p udp -l 1KB
docker exec -d clab-netgo-host-1 ./ethr -c 203.0.113.251 -b 400K -d 60s -p udp -l 1KB
cd ch10/packet-capture; go build -o packet-capture main.go
docker exec -it clab-netgo-host-2 /workdir/packet-capture/packet-capture
2022/02/28 21:50:25  flow record: 203.0.113.0:60087 <-> 203.0.113.252:8888
2022/02/28 21:50:25  flow record: 203.0.113.0:60087 <-> 
203.0.113.252:8888
2022/02/28 21:50:27  flow record: 203.0.113.0:40986 <-> 203.0.113.252:8888
2022/02/28 21:50:29  flow record: 203.0.113.0:60087 <-> 203.0.113.252:8888
2022/02/28 21:50:29  flow record: 203.0.113.0:49138 <-> 203.0.113.251:8888
2022/02/28 21:50:30  flow record: 203.0.113.0:60087 <-> 203.0.113.252:8888
2022/02/28 21:50:30  flow record: 203.0.113.0:49138 <-> 203.0.113.251:8888
```

As promised, before we move on to the next section, let’s review the first _developer experience_ tool of the chapter.

Just Imagine

# Debugging Go programs

Reading and reasoning about an existing code base is a laborious task, and it gets even harder as programs mature and evolve. This is why, when learning a new language, it’s very important to have at least a basic understanding of the debugging process. Debugging allows us to halt the execution of a program at a pre-defined place and step through the code line by line while examining in-memory variables and data structures.

In the following example, we use Delve to debug the `packet-capture` program we just ran. Before you can start, you need to generate some traffic through the lab topology with `make traffic-start`:

```markup
$ make traffic-start
docker exec -d clab-netgo-cvx systemctl restart hsflowd
docker exec -d clab-netgo-host-3 ./ethr -s
docker exec -d clab-netgo-host-1 ./ethr -c 203.0.113.253 -b 900K -d 60s -p udp -l 1KB
docker exec -d clab-netgo-host-1 ./ethr -c 203.0.113.252 -b 600K -d 60s -p udp -l 1KB
docker exec -d clab-netgo-host-1 ./ethr -c 203.0.113.251 -b 400K -d 60s -p udp -l 1KB
```

The Delve binary file is already pre-installed in the `host` lab containers, so you can connect to the `host-2` container with the `docker exec -it` command and start the Delve shell with the `dlv` `debug` command:

```markup
$ docker exec -it clab-netgo-host-2 bash
root@host-2:/# cd workdir/ch10/packet-capture/
root@host-2:/workdir/packet-capture# dlv debug main.go
```

Once in the `dlv` interactive shell, you can use different built-in commands to control the execution of the program (you can use `help` to view the full list of commands). Set a breakpoint at line 49 of `main.go` and run the program until the point where we receive the first packet:

```markup
(dlv) break main.go:49
Breakpoint 1 set at 0x5942ce for main.main() ./main.go:49
(dlv) continue
> main.main() ./main.go:49 (hits goroutine(1):1 total:1) (PC: 0x5942ce)
    44:    packetSource := gopacket.NewPacketSource(
    45:      handle,
    46:      layers.LayerTypeEthernet,
    47:    )
    48:    for packet := range packetSource.Packets() {
=>  49:      if l4 := packet.TransportLayer(); l4 == nil {
    50:        continue
    51:      }
    52:  
    53:      sflowLayer := packet.Layer(layers.LayerTypeSFlow)
    54:      if sflowLayer != nil {
```

When execution stops at a breakpoint, you can examine the local variables using the `locals` command:

```markup
(dlv) locals
err = error nil
handle = ("*github.com/google/gopacket/pcapgo.EthernetHandle")(0xc000162200)
rawInstructions = []golang.org/x/net/bpf.RawInstruction len: 6, cap: 6, [...]
packetSource = ("*github.com/google/gopacket.PacketSource")(0xc00009aab0)
packet = github.com/google/gopacket.Packet(*github.com/google/gopacket.eagerPacket) 0xc0000c3c08
```

You can print the variable contents on a screen, as in the following example for the `packet` variable:

```markup
(dlv) print packet
github.com/google/gopacket.Packet(*github.com/google/gopacket.eagerPacket) *{
  packet: github.com/google/gopacket.packet {
    data: []uint8 len: 758, cap: 758, [170,193,171,140,219,204,170,193,171,198,150,242,8,0,69,0,2,232,40,71,64,0,63,17,18,182,192,0,2,5,203,0,113,2,132,19,24,199,2,212,147,6,0,0,0,5,0,0,0,1,203,0,113,129,0,1,134,160,0,0,0,39,0,2,...+694 more],
    /* ... < omitted > ... */
    last: github.com/google/gopacket.Layer(*github.com/google/gopacket.DecodeFailure) ...,
    metadata: (*"github.com/google/gopacket.PacketMetadata")(0xc0000c6200),
    decodeOptions: (*"github.com/google/gopacket.DecodeOptions")(0xc0000c6250),
    link: github.com/google/gopacket.LinkLayer(*github.com/google/gopacket/layers.Ethernet) ...,
    network: github.com/google/gopacket.NetworkLayer(*github.com/google/gopacket/layers.IPv4) ...,
    transport: github.com/google/gopacket.TransportLayer(*github.com/google/gopacket/layers.UDP) ...,
    application: github.com/google/gopacket.ApplicationLayer nil,
    failure: github.com/google/gopacket.ErrorLayer(*github.com/google/gopacket.DecodeFailure) ...,},}
```

The text-based navigation and verbosity of the output may be intimidating for beginners, but luckily, we have alternative visualization options.

## Debugging from an IDE

If debugging in a console is not your preferred option, most of the popular **Integrated Development Environments** (**IDEs**) come with some form of support for Go debugging. For example, Delve integrates with **Visual Studio Code** (**VSCode**) and you can also configure it for remote debugging.

Although you can set up VSCode for remote debugging in different ways, in this example, we run Delve manually inside a container in the `headless` mode while specifying the port at which to listen for incoming connections:

```markup
$ docker exec -it clab-netgo-host-2 bash 
root@host-2:/# cd workdir/ch10/packet-capture/
root@host-2:/workdir/ch10/packet-capture#  dlv debug main.go --listen=:2345 --headless --api-version=2
API server listening at: [::]:2345
```

Now, we need to tell VSCode how to connect to the remote Delve process. You can do this by including a JSON config file in the `.vscode` folder next to the `main.go` file. Here’s an example file you can find in `ch10/packet-capture/.vscode/launch.json` in this book’s GitHub repository:

```markup
{
"version": "0.2.0",
"configurations": [
        {
            "name": "Connect to server",
            "type": "go",
            "request": "attach",
            "mode": "remote",
            "remotePath": "/workdir/ch10/packet-capture",
            "port": 2345,
            "host": "ec2-3-224-127-79.compute-1.amazonaws.com",  
        },
    ]
}
```

You need to replace the `host` value with the one where the lab is running and then start an instance of VSCode from the root of the Go program (`code ch10/packet-capture`):

![Figure 10.2 – VSCode development environment ](https://static.packt-cdn.com/products/9781800560925/graphics/image/B16971_10_02.jpg)

Figure 10.2 – VSCode development environment

In VSCode, now you can go to the debug icon in the left menu to get to **RUN AND DEBUG**, where you should see the **Connect to server** option that reads the preceding JSON config file. Click on the green arrow to connect to the remote debugging process.

At this point, you can navigate through the code and examine local variables inside the VSCode **user interface** (**UI**), while the debugging process is running inside a container:

![Figure 10.3 – VSCode debugging](https://static.packt-cdn.com/products/9781800560925/graphics/image/B16971_10_03.jpg)

Figure 10.3 – VSCode debugging

In the next section, we will look at how to add value to the data plane telemetry we collect and process by aggregating it to generate a report of the highest bandwidth consumers.

Just Imagine

# Data plane telemetry aggregation

After collecting and parsing data plane telemetry, we need to think about what to do with it next. Looking at raw data is not always helpful because of the sheer number of flows and lack of any meaningful context. Hence, the next logical step in a telemetry processing pipeline is data enrichment and aggregation.

Telemetry enrichment refers to the process of adding extra metadata to each flow based on some external source of information. For example, these external sources can provide a correlation between a public IP and its country of origin or BGP ASN, or between a private IP and its aggregate subnets or device identity.

Another technique that can help us interpret and reason about the telemetry we collect is aggregation. We can combine different flow records either based on the IP prefix boundary or flow metadata, such as a BGP ASN, to help network operators draw meaningful insights and create high-level views of the data.

You could build the entire telemetry processing pipeline out of open source components with ready-to-use examples (see _Further reading_) available on the internet, but sooner or later, you might need to write some code to meet your specific business requirements. In the following section, we will work on a scenario where we need to aggregate data plane telemetry to better understand the traffic patterns in our network.

## Top talkers

In the absence of long-term telemetry storage, getting a just-in-time snapshot of the highest bandwidth consumers can be quite helpful. We refer to this application as _top talkers_, and it works by displaying a list of network flows that are sorted based on their relative interface bandwidth utilization.

Let’s walk through an example Go application that implements this feature.

### Exploring telemetry data

In our `top-talkers` application, we collect sFlow records with `netsampler/goflow2`, a package designed specifically to collect, enrich, and save sFlow, IPFIX, or NetFlow telemetry. This package ingests raw protocol data and produces normalized (protocol-independent) flow records. By default, you can save these normalized records in a file or send them to a Kafka queue. In our case, we store them in memory for further processing.

To store the flow records in memory, we save the most relevant fields of each flow record we receive in a user-defined data structure we call `MyFlow`:

```markup
type MyFlow struct {
     Key         string
     SrcAddr     string `json:"SrcAddr,omitempty"`
     DstAddr     string `json:"DstAddr,omitempty"`
     SrcPort     int    `json:"SrcPort,omitempty"`
     DstPort     int    `json:"DstPort,omitempty"`
     Count       int    // times we've seen this flow sample
}
```

Additionally, we create a flow key as a concatenation of the ports and IP addresses of the source and destination to uniquely identify each flow:

![Figure 10.4 – A flow key](https://static.packt-cdn.com/products/9781800560925/graphics/image/B16971_10_04.jpg)

Figure 10.4 – A flow key

To help us calculate the final result, we create another data structure we call `topTalker`, which has two fields:

-   `flowMap`: A map to store a collection of `MyFlow`\-type flows. We use the key we created to index them.
-   `Heap`: A helper data structure that keeps track of the most frequently seen flows:
    
    ```markup
    type Heap []*MyFlow
    ```
    
    ```markup
    type topTalker struct {
    ```
    
    ```markup
         flowMap map[string]*MyFlow
    ```
    
    ```markup
         heap    Heap
    ```
    
    ```markup
    }
    ```
    

Since we use a high-level sFlow package (`goflow2`), we don’t need to worry about setting up a UDP listener or receiving and decoding packets, but we need to tell `goflow2` the format to report flow records (`json`) and point to a custom transport driver (`tt`) that determines what to do with the data after the sFlow package normalizes the received flow records:

```markup
import (
  "github.com/netsampler/goflow2/format"
  "github.com/netsampler/goflow2/utils"
)
func main() {
     tt := topTalker{
           flowMap: make(map[string]*MyPacket),
           heap:    make(Heap, 0),
     }
     formatter, err := format.FindFormat(ctx, "json")
     // process error
     sSFlow := &utils.StateSFlow{
           Format:    formatter,
           Logger:    log.StandardLogger(),
           Transport: &tt,
     }
     go sSFlow.FlowRoutine(1, hostname, 6343, false)
}
```

The `Transport` field in the `utils.StateSFlow` type of the preceding code snippet accepts any type that implements `TransportInterface`. This interface expects a single method (`Send()`) where all the enrichment and aggregation may take place:

```markup
type StateSFlow struct {
     Format    format.FormatInterface
     Transport transport.TransportInterface
     Logger    Logger
     /* ... < other fields > ... */
}
type TransportInterface interface {
     Send(key, data []byte) error
}
```

The `Send` method accepts two arguments, one representing the source IP of an sFlow datagram and the second one containing the actual flow record.

### Telemetry processing

In our implementation of the `Send` method (to satisfy the `TransportInterface` interface), we first parse the input binary data and deserialize it into a `MyFlow` data structure:

```markup
func (c *topTalker) Send(key, data []byte) error {
     var myFlow MyFlow
     json.Unmarshal(data, &myFlow)
     /* ... <continues next > ... */
}
```

Bearing in mind that sFlow can capture packets going in either direction, we need to ensure that both flows count toward the same in-memory flow record. This means creating a special flow key that satisfies the following two conditions:

-   It must be the same for both ingress and egress packets of the same flow.
-   It must be unique for all bidirectional flows.

We do this by sorting the source and destination IPs when constructing the bidirectional flow key, as the next code snippet shows:

```markup
var flowMapKey = `%s:%d<->%s:%d`
func (c *topTalker) Send(key, data []byte) error {
  /* ... <continues from before > ... */
  ips := []string{myFlow.SrcAddr, myFlow.DstAddr}
  sort.Strings(ips)
  var mapKey string
  if ips[0] != myFlow.SrcAddr {
    mapKey = fmt.Sprintf(
      flowMapKey,
      myFlow.SrcAddr,
      myFlow.SrcPort,
      myFlow.DstAddr,
      myFlow.DstPort,
    )
  } else {
    mapKey = fmt.Sprintf(
      flowMapKey,
      myFlow.DstAddr,
      myFlow.DstPort,
      myFlow.SrcAddr,
      myFlow.SrcPort,
    )
  }
  /* ... <continues next > ... */
}
```

With a unique key that represents both directions of a flow, we can save it in the map (`flowMap`) to store in memory. For each received flow record, the `Send` method performs the following checks:

-   If this is the first time we’ve seen this flow, then we save it on the map and set the count number to `1`.
-   Otherwise, we update the flow by incrementing its count by one:

```markup
func (c *topTalker) Send(key, data []byte) error {
  /* ... <continues from before > ... */
    myFlow.Key = mapKey
    foundFlow, ok := c.flowMap[mapKey]
    if !ok {
          myFlow.Count = 1
          c.flowMap[mapKey] = &myFlow
          heap.Push(&c.heap, &myFlow)
          return nil
    }
    c.heap.update(foundFlow)
    return nil
} 
```

Now, to display the top talkers in order, we need to sort the flow records we have saved. Here, we use the `container/heap` package from the Go standard library. It implements a sorting algorithm, offering O(log n) (logarithmic) upper-bound guarantees, which means it can do additions and deletions of data very efficiently.

To use this package, you only need to teach it how to compare your items. As you add, remove, or update elements, it will sort them automatically. In our example, we want to sort flow records saved as the `MyFlow` data type. We define `Heap` as a list of pointers to `MyFlow` records. The `Less()` method instructs the `container/heap` package to compare two `MyFlow` elements, based on the `Count` field that stores the number of times we have _seen_ a flow record:

```markup
type Heap []*MyFlow
func (h Heap) Less(i, j int) bool {
     return h[i].Count > h[j].Count
}
```

With this, we now have an in-memory flow record store with elements sorted according to their `Count`. We can now iterate over the `Heap` slice and print its elements on the screen. As in the earlier example with `gopacket`, we use `ethr` to generate three UDP flows with different throughputs to get a consistently sorted output. You can trigger the flows in the topology with `make top-talkers-start`:

```markup
Network-Automation-with-Go $ make top-talkers-start
docker exec -d clab-netgo-cvx systemctl restart hsflowd
docker exec -d clab-netgo-host-3 ./ethr -s
docker exec -d clab-netgo-host-1 ./ethr -c 203.0.113.253 -b 900K -d 60s -p udp -l 1KB
docker exec -d clab-netgo-host-1 ./ethr -c 203.0.113.252 -b 600K -d 60s -p udp -l 1KB
docker exec -d clab-netgo-host-1 ./ethr -c 203.0.113.251 -b 400K -d 60s -p udp -l 1KB
```

Then, run the Top-talkers Go application with `go run main.go` from within the `host-2` container (`clab-netgo-host-2`) to get a real-time Top-talkers table:

```markup
$ cd ch10/top-talkers; sudo ip netns exec clab-netgo-host-2 /usr/local/go/bin/go run main.go; cd ../../
Top Talkers
+---+-------------------+--------------------+------
| # | FROM              | TO                 | PROTO 
+---+-------------------+--------------------+------
| 1 | 203.0.113.253:8888 | 203.0.113.0:48494 | UDP | 
| 2 | 203.0.113.252:8888 | 203.0.113.0:42912 | UDP | 
| 3 | 203.0.113.251:8888 | 203.0.113.0:42882 | UDP | 
+---+-------------------+--------------------+------
```

Note that due to low traffic volume, random packet sampling, and limited test duration, your results may be slightly different but should converge to a similar distribution after several test iterations.

## Testing Go programs

Code testing is an integral part of any production software development process. Good test coverage improves application reliability and increases tolerance to bugs introduced at later stages of software development. Go has native support for testing with its `testing` package from the standard library and built-in command-line tool, `go test`. With test coverage built into the Go tool, it’s uncommon to see third-party packages used for testing Go code.

Table-driven testing is one of the most popular testing methodologies in Go. The idea is to describe test cases as a slice of custom data structures, with each one providing both inputs and expected results for each test case. Writing test cases as a table makes it easier to create new scenarios, consider corner cases, and interpret existing code behaviors.

We can test part of the code of the `top-talkers` example we just reviewed by building a set of table tests for the heap implementation we used to sort the flow records.

Let’s create a test file, `main_test.go`, with a single test function in it:

```markup
package main
import (
     "container/heap"
     "testing"
)
func TestHeap(t *testing.T) {
  // code tests
}
```

Both the `_test.go` filename suffix and the `Test<Name>` function prefix are naming conventions that allow Go to detect testing code and exclude it during binary compilation.

We design each test case to have all the relevant information, including the following:

-   A name to use in error messages
-   A set of unique flows described by their starting counters and resulting positions:
    
    ```markup
    type testFlow struct {
    ```
    
    ```markup
         startCount   int
    ```
    
    ```markup
         timesSeen    int
    ```
    
    ```markup
         wantPosition int
    ```
    
    ```markup
         wantCount    int
    ```
    
    ```markup
    }
    ```
    
    ```markup
    type testCase struct {
    ```
    
    ```markup
         name  string
    ```
    
    ```markup
         flows map[string]testFlow
    ```
    
    ```markup
    }
    ```
    

Given the preceding definitions, we create a test suite for a different combination of input and output values to cover as many non-repeating scenarios as possible:

```markup
 var testCases = []testCase{
  {
    name: "single packet",
    flows: map[string]testFlow{
      "1-1": {
        startCount:   1,
        timesSeen:    0,
        wantPosition: 0,
        wantCount:    1,
      },
    },
  },{
    name: "last packet wins",
    flows: map[string]testFlow{
      "2-1": {
        startCount:   1,
        timesSeen:    1,
        wantPosition: 1,
        wantCount:    2,
      },
      "2-2": {
        startCount:   2,
        timesSeen:    1,
        wantPosition: 0,
        wantCount:    3,
      },
    },
  },
```

We tie all this together in the body of the `TestHeap` function, where we iterate over all test cases. For each test case, we set up its preconditions, push all flows on the heap, and update their count `timeSeen` number of times:

```markup
func TestHeap(t *testing.T) {
     for _, test := range testCases {
           h := make(Heap, 0)
           // pushing flow on the heap
           for key, f := range test.flows {
                      flow := &MyFlow{
                           Count: f.startCount,
                           Key:   key,
                      }
                      heap.Push(&h, flow)
                      // updating packet counts
                      for j := 0; j < f.timesSeen; j++ {
                           h.update(flow)
                      }
           }
     /* ... <continues next > ... */
}
```

Once we have updated all flows, we remove them off the heap, one by one, based on the highest count, and check whether the resulting position and count match what we had described in the test case. In case of a mismatch, we generate an error message using the `*testing.T` type injected by the testing package:

```markup
func TestHeap(t *testing.T) {
  /* ... < continues from before > ... */
  for i := 0; h.Len() > 0; i++ {
                f := heap.Pop(&h).(*MyFlow)
                tf := test.flows[f.Key]
                if tf.wantPosition != i {
                           t.Errorf(
                             "%s: unexpected position for packet key %s: got %d, want %d", test.name, f.Key, i, tf.wantPosition)
                }
                if tf.wantCount != f.Count {
                           t.Errorf(
                                 "%s: unexpected count for packet key %s: got %d, want %d", test.name, f.Key, f.Count, tf.wantCount)
                }
           }
}
```

Thus far, we’ve only discussed data plane telemetry, which is crucial, but not the only element of network monitoring. In the following section, we will explore network control plane telemetry by building a complete end-to-end telemetry processing pipeline.

Just Imagine

# Measuring control plane performance

Most network engineers are familiar with tools such as `ping`, `traceroute`, and `iperf` to verify network data plane connectivity, reachability, and throughput. At the same time, control plane performance often remains a black box, and we can only assume how long it takes for our network to re-converge. In this section, we aim to address this problem by building a control plane telemetry solution.

Modern control plane protocols, such as BGP, distribute large volumes of information from IP routes to MAC addresses and flow definitions. As the size of our networks grows, so does the churn rate of the control plane state, with users, VMs, and applications constantly moving between different locations and network segments. Hence, it’s critical to have visibility of how well our control plane performs to troubleshoot network issues and take any preemptive actions.

The next code example covers the telemetry processing pipeline we built to monitor the control plane performance of the lab network. At the heart of it, there is a special `bgp-ping` application that allows us to measure the round-trip time of a BGP update. In this solution, we take advantage of the features of the following Go packages and applications:

-   `jwhited/corebgp`: A pluggable implementation of a BGP finite state machine that allows you to run arbitrary actions for different BGP states.
-   `osrg/gobgp`: One of the most popular BGP implementations in Go; we use it to encode and decode BGP messages.
-   `cloudprober/cloudprober`: A flexible distributed probing and monitoring framework.
-   `Prometheus` and `Grafana`: A popular monitoring and visualization software stack.

![Figure 10.5 – Telemetry pipeline architecture](https://static.packt-cdn.com/products/9781800560925/graphics/image/B16971_10_05.jpg)

Figure 10.5 – Telemetry pipeline architecture

To bring up this entire setup, you can run `make bgp-ping-start` from the root of this book’s GitHub repository (see _Further reading_):

```markup
Network-Automation-with-Go $ make bgp-ping-start
cd ch10/bgp-ping; go build -o bgp-ping main.go
docker exec -d clab-netgo-host-3 /workdir/bgp-ping/bgp-ping -id host-3 -nlri 100.64.0.2 -laddr 203.0.113.254 -raddr 203.0.113.129 -las 65005 -ras 65002 -p
docker exec -d clab-netgo-host-1 /workdir/bgp-ping/bgp-ping -id host-1 -nlri 100.64.0.0 -laddr 203.0.113.0 -raddr 203.0.113.1 -las 65003 -ras 65000 -p
docker exec -d clab-netgo-host-2 /cloudprober -config_file /workdir/workdir/cloudprober.cfg
cd ch10/bgp-ping; docker-compose up -d; cd ../../
Creating prometheus ... done
Creating grafana    ... done
http://localhost:3000
```

The final line of the preceding output shows the URL that you can use to access the deployed instance of Grafana, using `admin` as both `username` and `password`:

![Figure 10.6 – BGP ping dashboard](https://static.packt-cdn.com/products/9781800560925/graphics/image/B16971_10_06.jpg)

Figure 10.6 – BGP ping dashboard

This instance has a pre-created dashboard called `BGP-Ping` that plots the graph of BGP round-trip times in milliseconds.

It’s important to note that there’s a lot more to routing protocol convergence and performance than the update propagation time. Other important factors may include update churn due to transient events or **Forwarding Information Base** (**FIB)** programming time. We focus on a single-dimension metric in this example, but in reality, you may want to consider other performance metrics as well.

Just Imagine

# Measuring BGP Update propagation time

As the standard `ping`, the `bgp-ping` application works by sending and receiving probe messages. A sender embeds a probe in a BGP Update message and sends it to its BGP neighbor. We encode the probe as a custom BGP optional transitive attribute, which allows it to propagate transparently throughout the network until it reaches one of the `bgp-ping` responders.

A `bgp-ping` responder recognizes this custom transitive attribute and reflects it back to the sender. This gives the sender a measure of BGP Update propagation delay within the network, which is then reported to an external metric consumer or printed on a screen.

Since the `bgp-ping` application needs to inter-operate with real BGP stacks, at the very least it has to implement the initial exchange of `Open` messages to negotiate the BGP session capabilities, followed by the periodic exchange of `Keepalive` messages. We also need to do the following:

1.  Send BGP Update messages triggered by different events.
2.  Encode and decode custom BGP attributes.

Let’s see how we can implement these requirements using open source Go packages and applications.

## Event-driven BGP state machine

We use CoreBGP (`jwhited/corebgp`) to establish a BGP session with a peer and keep it alive until it’s shut down. This gets us the `Open` and `Keepalive` messages we just discussed.

Inspired by the popular DNS server CoreDNS, CoreBGP is a minimalistic BGP server that you can extend through event-driven plugins.

In practice, you extend the initial capabilities by building a custom implementation of the `Plugin` interface. This interface defines different methods that can implement user-defined behavior at certain points of the BGP **finite state** **machine** (**FSM**):

```markup
type Plugin interface {
     GetCapabilities(...) []Capability
     OnOpenMessage(...) *Notification
     OnEstablished(...) handleUpdate
     OnClose(...)
}
```

For the `bpg-ping` application, we only need to send and receive BGP Update messages, so we focus on implementing the following two methods:

-   `OnEstablished`: To send BGP Update messages.
-   `handleUpdate`: We use this to process received updates, identify ping requests, and send a response message.

The following diagram shows the main functional blocks of this application:

![Figure 10.7 – BGP Ping Design](https://static.packt-cdn.com/products/9781800560925/graphics/image/B16971_10_07.jpg)

Figure 10.7 – BGP Ping Design

Let’s start the code overview by examining the BGP Update handling logic (`handleUpdate`). Since our goal is to parse and process BGP ping probes, we can make sure we discard any other BGP updates early in the code. For every BGP Update message we receive, we check whether any of the BGP attributes have the custom `bgpPingType` transitive attribute we created to signal the probe or ping. We silently ignore BGP updates that don’t have this attribute with a `continue` statement:

```markup
import bgp "github.com/osrg/gobgp/v3/pkg/packet/bgp"
const (
     bgpPingType = 42
)
func (p *plugin) handleUpdate(
     peer corebgp.PeerConfig,
     update []byte,
) *corebgp.Notification {
 
     msg, err := bgp.ParseBGPBody(
           &bgp.BGPHeader{Type: bgp.BGP_MSG_UPDATE},
           update,
     )
     // process error
     for _, attr := range msg.Body.
                    (*bgp.BGPUpdate).PathAttributes {
           if attr.GetType() != bgpPingType {
                      continue
           }
     /* ... <continues next > ... */
}
```

Once we have determined that it’s a BGP ping message, we deal with two possible scenarios:

-   If it’s a **ping response**, we calculate the round-trip time using the timestamp extracted from the `bgpPingType` path attribute.
-   If it’s a **ping request**, we trigger a ping response by sending parsed data over a channel to the `OnEstablished` function:

```markup
func (p *plugin) handleUpdate(
  peer corebgp.PeerConfig,
  update []byte,
) *corebgp.Notification {
    /* ... < continues from before > ... */
    source, dest, ts, err := parseType42(attr)
    // process error
    sourceHost := string(bytes.Trim(source, "\x00"))
    destHost := string(bytes.Trim(dest, "\x00"))
    /* ... <omitted for brevity > ... */
 
    // if src is us, may be a response. id = router-id
    if sourceHost == *id {
      rtt := time.Since(ts).Nanoseconds()
      metric := fmt.Sprintf(
        "bgp_ping_rtt_ms{device=%s} %f\n",
        destHost,
        float64(rtt)/1e6,
      )
 
    p.store = append(p.store, metric)
return nil
    }
 
    p.pingCh <- ping{source: source, ts: ts.Unix()}
    return nil
}
```

The event-driven logic to send BGP updates lives in the `OnEstablished()` method that has a three-way select statement to listen for triggers over Go channels, representing three different states of the `bgp-ping` application:

-   Responding to a received ping request, triggered by a request coming from the `handleUpdate` function
-   Firing a new ping request, triggered by an external signal
-   Sending a scheduled withdraw message at the end of the probing cycle:

```markup
func (p *plugin) OnEstablished(
  peer corebgp.PeerConfig,
  writer corebgp.UpdateMessageWriter,
) corebgp.UpdateMessageHandler {
  log.Println("peer established, starting main loop")
  go func() {
    for {
      select {
      case pingReq := <-p.pingCh:
        // Build the ping response payload
        bytes, err := p.buildUpdate(
                      type42PathAttr,
                      peer.LocalAddress,
                      peer.LocalAS,
        )
        // process error
        writer.WriteUpdate(bytes)
        /* ... < schedule a withdraw > ... */
 
      case <-p.probeCh:
        // Build the ping request payload
        bytes, err := p.buildUpdate(
                      type42PathAttr,
                      peer.LocalAddress,
                      peer.LocalAS,
        )
        // process error
        writer.WriteUpdate(bytes)
        /* ... < schedule a withdraw > ... */
 
      case <-withdraw.C:
        bytes, err := p.buildWithdraw()
        // process error
        writer.WriteUpdate(bytes)
      }
    }
  }()
  return p.handleUpdate
}
```

One caveat of CoreBGP is that it doesn’t include its own BGP message parser or builder. It sends any raw bytes that may confuse or even crash a standard BGP stack, so always use it with caution.

Now, we need a way to parse and craft a BGP message, and here is where we can use another Go library called `GoBGP`.

## Encoding and decoding BGP messages

GoBGP is a full-blown BGP stack and supports most of the BGP address families and features. However, since we already use CoreBGP for BGP state management, we limit the use of GoBGP to message encoding and decoding.

For example, whenever we need to build a BGP withdraw update message, we call a helper function (`buildWithdraw`) that uses GoBGP to build the message. GoBGP allows us to include only the relevant information, such as a list of **Network Layer Reachability Information** (**NLRI**), while it takes care of populating the rest of the fields, such as type, length, and building a syntactically correct BGP message:

```markup
func (p *plugin) buildWithdraw() ([]byte, error) {
     myNLRI := bgp.NewIPAddrPrefix(32, p.probe.String())
     withdrawnRoutes := []*bgp.IPAddrPrefix{myNLRI}
     msg := bgp.NewBGPUpdateMessage(
           withdrawnRoutes,
           []bgp.PathAttributeInterface{},
           nil,
     )
     return msg.Body.Serialize()
}
```

Here’s another example of how to use GoBGP to parse a message received by CoreBGP. We take a slice of bytes and use the `ParseBGPBody` function to deserialize it into GoBGP’s `BGPMessage` type:

```markup
func (p *plugin) handleUpdate(
     peer corebgp.PeerConfig,
     update []byte,
) *corebgp.Notification {
     msg, err := bgp.ParseBGPBody(
           &bgp.BGPHeader{Type: bgp.BGP_MSG_UPDATE},
           update,
     )
     // process error
     if err := bgp.ValidateBGPMessage(msg); err != nil {
           log.Fatal("validate BGP message ", err)
     }
```

You can now further parse this BGP message to extract various path attributes and NLRIs, as we’ve seen in the earlier overview of the `handleUpdate` function.

## Collecting and exposing metrics

The `bgp-ping` application can run as a standalone process and print the results on a screen. We also want to be able to integrate our application into more general-purpose system monitoring solutions. To do that, it needs to expose its measurement results in a standard format that an external monitoring system can understand.

You can implement this capability natively by adding a web server and publishing your metrics for external consumers, or you can use an existing tool that collects and exposes metrics on behalf of your application. One tool that does this is Cloudprober, which enables automated and distributed probing and monitoring, and offers native Go integration with several external probes.

We integrate the `bgp-ping` application with the Cloudprober via its `serverutils` package, which allows you to exchange probe requests and replies over the **standard input (stdin)** and **standard output (stdout)**. When we start `bgp-ping` with a `-c` flag, it expects all probe triggers to come from Cloudprober and sends its results back in a `ProbeReply` message:

```markup
func main() {
  /* ... < continues from before > ... */
  probeCh := make(chan struct{})
  resultsCh := make(chan string)
  
  peerPlugin := &plugin{
              probeCh: probeCh,
            resultsCh: resultsCh,
  }
 
  if *cloudprober {
    go func() {
      serverutils.Serve(func(
        request *epb.ProbeRequest,
        reply *epb.ProbeReply,
      ) {
        probeCh <- struct{}{}
        reply.Payload = proto.String(<-resultsCh)
        if err != nil {
          reply.ErrorMessage = proto.String(err.Error())
        }
      })
    }()
  }
}
```

The Cloudprober application itself runs as a pre-compiled binary and requires minimal configuration to tell it about the `bgp-ping` application and its runtime options:

```markup
probe {
  name: "bgp_ping"
  type: EXTERNAL
  targets { dummy_targets {} }
  timeout_msec: 11000
  interval_msec: 10000
  external_probe {
    mode: SERVER
    command: "/workdir/bgp-ping/bgp-ping -id host-2 -nlri 100.64.0.1 -laddr 203.0.113.2 -raddr 203.0.113.3 -las 65004 -ras 65001 -c true"
  }
}
```

All measurement results are automatically published by Cloudprober in a format that most popular cloud monitoring systems can understand.

## Storing and visualizing metrics

The final stage in this control plane telemetry processing pipeline is metrics storage and visualization. Go is a very popular choice for these systems, with examples including Telegraf, InfluxDB, Prometheus, and Grafana.

The current telemetry processing example includes Prometheus and Grafana with their respective configuration files and pre-built dashboards. The following configuration snippet points Prometheus at the local Cloudprober instance and tells it to scrape all available metrics every 10 seconds:

```markup
scrape_configs:
  - job_name: 'bgp-ping'
    scrape_interval: 10s
    static_configs:
      - targets: ['clab-netgo-host-2:9313']
```

Although we discuss little of it here, building meaningful dashboards and alerts is as important as doing the measurements. Distributed systems observability is a big topic that is extensively covered in existing books and online resources. For now, we will stop at the point where we see a visual representation of the data in a Grafana dashboard but don’t want to imply that a continuous linear graph of absolute values is enough. Most likely, to make any reasonable assumptions, you’d want to present your data as an aggregated distribution and monitor its outlying values over time, as this would give a better sign of increasing system stress and may serve as a trigger for any further actions.

Just Imagine

# Developing distributed applications

Building a distributed application, such as `bgp-ping`, can be a major undertaking. Unit testing and debugging can help spot and fix a lot of bugs, but these processes can be time-consuming. In certain cases, when an application has different components, developing your code iteratively may require some manual orchestration. Steps such as building binary files and container images, starting the software process, enabling logging, and triggering events are now something you need to synchronize and repeat for all the components that include your application.

The final developer experience tool that we will cover in this chapter was specifically designed to address the preceding issues. Tilt helps developers automate manual steps, and it has native integration with container and orchestration platforms, such as Kubernetes or Docker Compose. You let it know which files to monitor, and it will automatically rebuild your binaries, swap out container images, and restart existing processes, all while showing you the output logs of all applications on a single screen.

It works by reading a special `Tiltfile` containing a set of instructions on what to build and how to do it. Here’s a snippet from a Tiltfile that automatically launches a `bgp-ping` process inside one of the host containers and restarts it every time it detects a change to `main.go`:

```markup
local_resource('host-1',
  serve_cmd='ip netns exec clab-netgo-host-1 go run main.go -id host-1 -nlri 100.64.0.0 -laddr 203.0.113.0 -raddr 203.0.113.1 -las 65003 -ras 65000 -p',
  deps=['./main.go'])
```

The full `Tiltfile` has two more resources for the other two hosts in our lab network. You can bring up all three parts of the application with `sudo` `tilt up`:

```markup
Network-Automation-with-Go $ cd ch10/bgp-ping
Network-Automation-with-Go/ch10/bgp-ping $ sudo tilt up
Tilt started on http://localhost:10350/
```

Tilt has both a console (text) and a web UI that you can use to view the logs of all resources:

![Figure 10.8 – Tilt](https://static.packt-cdn.com/products/9781800560925/graphics/image/B16971_10_08.jpg)

Figure 10.8 – Tilt

Any change to the source code of the `bgp-ping` application would trigger a restart of all affected resources. By automating a lot of manual steps and aggregating the logs, this tool can simplify the development process of any distributed application.

Just Imagine

# Summary

This concludes the chapter about network monitoring. We have only touched upon a few selected subjects and admit that the topic of this chapter is too vast to cover in this book. However, we hope we have provided enough resources, pointers, and ideas for you to continue the exploration of network monitoring, as it’s one of the most vibrant and actively growing areas of the network engineering discipline.

Just Imagine

# Further reading

-   Course’s GitHub repository: [https://github.com/ImagineDevOps DevOps/Network-Automation-with-Go](https://github.com/ImagineDevOps DevOps/Network-Automation-with-Go)
-   `google/gopacket` package: [https://github.com/google/gopacket](https://github.com/google/gopacket)
-   `gdb` documentation: [https://go.dev/doc/gdb](https://go.dev/doc/gdb)
-   `vscode-go`: [https://code.visualstudio.com/docs/languages/go](https://code.visualstudio.com/docs/languages/go)
-   `ch10/packet-capture/.vscode/launch.json`: [https://github.com/ImagineDevOps DevOps/Network-Automation-with-Go/blob/main/ch10/packet-capture/.vscode/launch.json](https://github.com/ImagineDevOps DevOps/Network-Automation-with-Go/blob/main/ch10/packet-capture/.vscode/launch.json)
-   Open source components with ready-to-use examples: [https://github.com/netsampler/goflow2/tree/main/compose/kcg](https://github.com/netsampler/goflow2/tree/main/compose/kcg)
-   CoreBGP documentation: [https://pkg.go.dev/github.com/jwhited/corebgp#section-readme](https://pkg.go.dev/github.com/jwhited/corebgp#section-readme)