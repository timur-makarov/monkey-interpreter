### Interpreter for Monkey Programming Language

- Variables
- Conditions
- While Loops
- Functions
- Recursion
- Closures
- Arrays
- Hash tables
- Builtin functions

```monkey
let factorial = fn(n) {
    if (n < 1) {
        return 1
    } else {
        return n * factorial(n - 1)
    }
}

log("Factorial of 5:", factorial(5))

let makeCounter = fn() {
    let count = 0
    return fn() {
        count = count + 1
        return count
    }
}

let counter = makeCounter()

let arr = [1, 2, 3]
let hashTable = {"name": "Monkey", "type": "Language"}

log("Length of array:",  len(arr))
arr = append(arr, 4)

if (arr[2] > arr[0]) {
    log("Third element is greater than first.")
} else {
    log("Third element is not greater.")
}

while (counter() < len(arr) + 4) {
    log("Counter:", counter())
}

hashTable["version"] = "1.0"
log(hashTable["name"] + " version: " + hashTable["version"])
```