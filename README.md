# wordcount
Go implementation of word count using a simplified MapReduce approach.
## Installation
Install using the "go get" command:
```
go get -u github.com/redefik/wordcount/...
```
## Example Usage
* Move to the ```bin``` subdirectory of your go workspace.
* Launch the master node specifying its address:
  ```
  ./master localhost:1234 &
  ```
* Launch the worker nodes specifying address for each one:
  ```
  ./worker localhost:1235 &
  ./worker localhost:1236 &
  ./worker localhost:1237 &
  ./worker localhost:1238 &
  ./worker localhost:1239 &
  ```
* Prepare the JSON configuration file:
  ```
  {
  "master": [
    "localhost:1234"
  ],
  "mapper":[
    "localhost:1235", "localhost:1236", "localhost:1237"
  ],
  "reducer":[
    "localhost:1238", "localhost:1239"
  ],
  "outdir": "~/out"
  }
  ```
  In the example above, the first three workers act as mapper and the last two act as reducer.
  In the field ```outdir``` you can specify the directory where the output files will be stored.
  
* Launch the client specifying configuration and data files location:
  ```
  ./client -config=~/conf.json file1.txt file2.txt file3.txt file4.txt
  ```
* At completion, the program shows created files:
  ```
  Work done!
  You can find results in:
  out/out0.txt
  out/out1.txt
  ```
## Implementation Overview

* The client passes N input files to the master node using a synchronous RPC invocation.
* The master distribute the N files among M mappers.
  Each mapper is assigned R temporary files. The i-th file will be then processed by the i-th reducer.
* Each mapper counts the occurrences for the given files and stores the results into the R temporary files.
  The distribution among the temporary files is made by hashing.
* When all the mappers have completed, the master activates the R reducers and creates R output files.
  Each reducer receives M temporary files from the master and merges the intermediate results previously computed by the mappers.
  Final results are stored in the output file passed by the master.
* When all the reducers have completed, the master removes the temporary files and returns to the client the output files.
* Mapper and reducer are invoked asynchronously using Go method ```Go```.

The following assumptions have been made:
- The set of workers is known and does not change during computation;
- No worker fails;
- The master does not fail.
