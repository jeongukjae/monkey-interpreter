# monkey interpreter

A repository to implement [an interpreter for the Monkey programming language](https://interpreterbook.com).

## Run REPL in Web

I built a web page with WebAssembly to use REPL in web. Check out [this page](https://jeongukjae.github.io/monkey-interpreter/).

## Run REPL

```sh
$ go run main.go
Hello jeongukjae! This is the Monkey programming language!
Feel free to type in commands!
>> 1 + 1
2
>> let two = "two";
>> let hashmap = {"one": 10 - 9, two: 1 +1, "thr" + "ee": 6 / 2, 4 : 4, true: 5, false: 6}
>> hashmap[two]
2
>> puts(hashmap)
{three: 3, 4: 4, true: 5, false: 6, one: 1, two: 2}
null
>> let arr = [1,2,3,4,5]
>> map(arr, fn(x) {return x * 2})
[2, 4, 6, 8, 10]
>> reduce(map(arr, fn(x) { return x * 2;}), 0, fn(x, y) { return x + y; })
30
...
```

## Run test cases

```sh
make test
```
