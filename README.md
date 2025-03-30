# Virtual World Wide Web (VWWW)

This programme provide a a mock environnement to test crawlers. 

VWWW operates by running an HTTP server with a predefined seed. From this seed, the server generates a web graph—a collection of pages with links between them.

When the server receives a request, it deterministically generates the corresponding page throught the use of seed and returns it. If the page does not exist in the graph, the server responds with a 404 error.

Since the server only requires the seed to generate pages, it does not need to store the entire web graph in memory. This ensures a minimal and fixed memory footprint, allowing the server to simulate vast graphs while using very few resources.

Install
```bash
go install github.com/thebigroomxxl/vwww
```

Usage
```txt
Usage of vwww:
  -delay int
        delay in ms (default 1000)
  -pages int
        number of pages (default 10000)
  -seed int
        an int (default 42)
```


Todo
```
☐ response size parameter
☐ use of random distribution rather than fix parameters for delay and size
☐ mock library API
```
