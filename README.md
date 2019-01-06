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
* At completion, the program shows you created files:
  ```
  Work done!
  You can find results in:
  out/out0.txt
  out/out1.txt
  ```
