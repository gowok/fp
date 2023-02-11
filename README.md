# fp
Functional programming library for Go

# Usage

### Quality
It is well tested with high coverage:

<img width="700" alt="image" src="https://user-images.githubusercontent.com/16364286/218243184-2df92fdc-e433-433d-a285-5efcd3db0efc.png">

### Examples

```go
fmt.Println(strings.Repeat("*", 5))
// output: *****
```

```go
list := []int{1, 2, 3}
slices.ForEach(sample, func(s, i int) {
  fmt.Println(s, i)
})
// output:
// 1 0
// 2 1
// 3 2
```
