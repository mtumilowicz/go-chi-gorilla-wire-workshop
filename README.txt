# go-chi-gorilla-wire-workshop

* references
    * https://www.oreilly.com/library/view/learning-go/9781492077206/

## preface
* goals of this workshop
    * understanding basics of golang
        * data types
        * syntax
        * error handling
        * code organisation
        * testing
    * introduction to go ecosystem
        * chi
        * gorilla
        * wire
* workshop task: implement endpoint for deleting customer

## golang
* 25 keywords
    * example: `break`, `continue`, `if`, `for`, etc
    * not keywords: predeclared identifiers in universe block
        * example
            * built-in types (like `int` and `string`)
            * constants (like `true` and `false`)
            * functions (like `make` or `close`)
            * `nil`
        * can be shadowed in other scopes
            * example: `true := 10`
* Go runtime is compiled into every Go binary
    * different from languages using a virtual machine
        * VM must be installed separately to allow programs to run
    * avoids worries about compatibility issues between the runtime and the program
    * drawback: even the simplest Go program produces a binary that’s about 2 MB
* every type in Go is a value type
    * sometimes the value is a pointer
    * variables are passed by value
* uses GC
* standard library
    * has a compatibility promise
        * programs written for Go 1.x will continue to compile and run correctly with any future 1.x version of Go
    * File I/O: `io.Reader` and `io.Writer`
    * time
        * period = `time.Duration`
        * moment of time = `time.Time`
            * with a time zone
        * `time.After` - channel that outputs once specified duration elapses
        * `time.Tick` - returns a new value every time the specified duration elapses
    * `json.Unmarshal`, `json.Marshal`
        * specification: struct tags
            * strings that are written after the fields
    * net/http
        * production-quality HTTP/2 client and server
    * log/slog
        * since Go 1.21
        * zap, logrus, go-kit log, and many others.
        * slog.Debug("debug log message")
          slog.Info("info log message")
          slog.Warn("warning log message")
          slog.Error("error log message")
        * userID := "fred"
          loginCount := 20
          slog.Info("user login",
          "id", userID,
          "login_count", loginCount)

          2023/04/20 23:36:38 INFO user login id=fred login_count=20
        * json
            * options := &slog.HandlerOptions{Level: slog.LevelDebug}
              handler := slog.NewJSONHandler(os.Stderr, options)
              mySlog := slog.New(handler)
              lastLogin := time.Date(2023, 01, 01, 11, 50, 00, 00, time.UTC)
              mySlog.Debug("debug message",
              "id", userID,
              "last_login", lastLogin)
            * {"time":"2023-04-22T23:30:01.170243-04:00","level":"DEBUG",
              "msg":"debug message","id":"fred","last_login":"2023-01-01T11:50:00Z"}
        * For
          improved performance with fewer allocations, use the LogAttrs method instead:
          mySlog.LogAttrs(ctx, slog.LevelInfo, "faster logging",
          slog.String("id", userID),
          slog.Time("last_login", lastLogin))

## collections
* `make` method
* array
    * size part of the type
* slice
    * example
        ```
        var x = []int{1, 2} // slice
        var y = [2]int{1, 2} // array
        var z = y[:] // conversion: array -> slice
        var zz = [2]int(x) // conversion: slice -> array
        ```
    * used most of the time instead of array
    * has a capacity (number of consecutive memory locations reserved)
        * if length = capacity, `append` function uses the Go runtime to allocate a new backing array
          with a larger capacity
    * add - append
        ```
        var x = []int{1, 2}
        var y = append(x, 3) // usually shadowed: x = append(x, 4)
        fmt.Println(x) // 1, 2
        fmt.Println(y) // 1, 2, 3
        ```
    * emptying - clear
    * delete - involves more than just removing it, as slices are views into arrays
        * `s = append(s[:i], s[i+1:]...)`
    * isn’t comparable
        * `slices.Equal` returns true if slices are the same length and all of the elements are equal
        * `slices.EqualFunc` lets you pass in a function to compare elements
    * zero-length slice: `var x = []int{}`
        * is useful only when converting a slice to JSON
        * favor nil slices: `var x []int`
    * subslicing: `e := x[a:b]`
        * no starting offset => 0 is assumed
        * no ending offset => end of the slice is substituted
        * not making a copy of the data
            * changes to an element affect all slices that share that element
            * if needed => built-in copy function
        * extra confusing when combined with append
            * example
                ```
                original := []int{1, 2, 3, 4, 5}
                subslice := original[1:3]
                subslice = append(subslice, 6)
                fmt.Println(original) // [1 2 3 6 5]
                fmt.Println(subslice) // [2 3 6]
                ```
            * never use append with a subslice
                * or use full slice expression - if `append` exceeds the capacity => new array for slice is allocated (preserving the original slice)
                    ```
                    subslice := original[1:3:3]
                    ```
    * slice is implemented as a struct with three fields
        * an int field for length, an int field for capacity, and a pointer to a block of memory
        * when a slice is copied, copy is made of the length, capacity, and the pointer
        * ideal for reusable buffers
            * slice that’s passed to a function can have its contents modified
            * slice can’t be resized
        * by default, you should assume that a slice is not modified by a function
* string
    * Go uses a sequence of bytes to represent a string
* map
    * if key is not in the map => map returns the zero value
        * ok idiom to differentiate between a key in the map vs not in the map
    * used to simulate set
        * ok idiom + `map[string]bool` or `map[keyType]struct{}`
            ```
            if _, exists := map[key]; exists {
                ...
            }
            for key := range map {
                ...
            }
            ```
        * `m := make(map[string]bool)`
            * boolean uses one byte
        * `m := make(map[keyType]struct{})`
            * empty struct uses zero bytes
    * operations
        * set: `map[key] = value`
        * get: `map[key]`
        * delete: `delete(map, key)`
        * clear(map)
    * isn't comparable
        * `maps.Equal`, `maps.EqualFunc`
    * to mitigate Hash DoS attacks, Go's map implementation includes a random component in its hash function
        * no way for attacker to send many requests with keys designed to hash to the same bucket
        * new map is created => Go generates a random number and incorporates it into the hash function
            * the same keys will hash to different buckets in different instances of maps
    * map is implemented as a pointer to a struct
        * avoid using maps for input parameters or return values, especially on public APIs

## struct
* no inheritance => no classes
* example
    ```
    type person struct {
        name string
        age int
    }
    ```
* anonymous structs
    * common in two situations: unmarshaling and marshaling
        * example
            ```
            var pet struct {
                Name string `json:"name"`
                Kind string `json:"kind"`
            }

            err := json.Unmarshal(jsonData, &pet)
            ```

## syntax
* function - `func`
    * is a type
        * example
            ```
            var f func(string) int
            f := func(s string) { // anonymous function
                ...
            }
            ```
    * emulate named and optional parameters => define a struct
    * supports variadic parameters: `func max(first, rest ...int) int`
        * converted to a slice - slice can be supplied as the input
    * allows for multiple return values
        * example
            ```
            func sum(vals []int) (int, error) {
                if len(vals) == 0 {
                    return 0, errors.New("empty slice") // Return an error for empty slice
                }

                total := 0
                for _, v := range vals {
                    total += v
                }
                return total, nil // Return sum and nil error indicating success
            }
            ```
        * convention: last return value from a function is an error
            * no error => nil is returned for the error parameter
    * closures = functions declared inside functions
* Go uses capitalization to determine whether a package-level identifier is visible outside the package
    * identifier whose name starts with an uppercase letter is exported
* defer
    * example
        ```
        file, err := os.Open("example.txt")
        if err != nil {
        	fmt.Println("Error:", err)
        	return
        }
        defer file.Close()
        ```
    * code within defer functions runs after the return statement
        * LIFO order - the last defer registered runs first
    * common pattern: function that allocates a resource returns also closure that cleans up the resource
        * Go doesn’t allow unused variables => program will not compile if the function is not called
* no enumeration type
    * solution: `iota`
        * assign an increasing value to a set of constants
        * "internal" purposes only
            * new identifier => all subsequent ones will be renumbered
    * example
        ```
        type Status int

        const (
            Pending Status = iota // 0
            Running               // 1 (implicitly repeated from iota)
            Paused                // 2
            Finished              // 3
            Failed                // 4
        )
        ```
* methods
    * are functions associated with a particular type
        * can be defined only at the package block level
        * functions can be defined inside any block
    * example
        ```
        type Rectangle struct {
            width  float64
            height float64
        }

        func (r Rectangle) Area() float64 {
            return r.width * r.height
        }
        ```
    * must be declared in the same package as their associated type
        * Go doesn’t allow you to add methods to types you don’t control
    * convention: use pointer receiver to indicate that a parameter might be modified by the function
        * type has a pointer receiver => common practice is to be consistent and use pointer receivers for all methods
    * Go automatically takes the address of the local variable when calling the method
        * example: `c.Method()` is converted to `(&c).Method()`
    * Go automatically dereferences the pointer when calling the method
        * example: `c.Method()` is converted to `(*c).Method()`
    * Go allows to call a method on a nil receiver
        * most of the time => not very useful
        * use case: initializing data structure when nil
            * example: insert value into a tree (or create one when nil)
* interfaces
    * type-safe duck typing
        * type does not declare that it implements an interface
        * type implements the interface = method set contains all interface's methods
    * Go added any as a type alias for interface{}
        * use case: data placeholder
    * rule: accept interfaces, return structs
    * interfaces are implemented as a struct with two pointer fields: value and type of the value
        * type field is non-nil => interface is non-nil
        * value pointer is non-nil => type pointer is non-nil
            * you cannot have a variable without a type
        * interface is `nil` <=> both the type and the value must be nil
            * interface variable is nil => invoking any methods triggers a panic
            * interface variable is not nil but value is nil => invoking any methods triggers a panic
                * assuming that methods of the assigned type don’t properly handle nil
        * two instances of an interface type are equal <=> their types are equal and their values are equal
            * type isn’t comparable => panic
                * example: interface as a map key
    * accept interfaces, return structs
        * exception: returning error interface
        * concrete type is returned => new methods and fields can be added without breaking existing code
            * new fields and methods can be ignored if call-sites using interfaces
            * example: database/sql/driver in std lib
                * defines a set of interfaces that define what a database driver must provide
                    * responsibility of the database driver author to provide concrete implementations
                    * almost all methods on all interfaces return interfaces
                * problem: starting in Go 1.8, database drivers are expected to support additional features
                    * existing interfaces can’t be updated with new method
                    * existing methods on these interfaces can’t be updated to return different types
                    * solution: define new interfaces and tell database driver authors that they should implement both
        * invoking a function with parameters of interface types => heap allocation occurs for each interface parameter
* embedding
    * form of composition
        * is not inheritance - cannot assign a variable of one type to another
    * allows to invoke directly fields and methods of one struct from another
    * example
        ```
        type Manager struct {
            Employee
            Reports []Employee
        }
        ```
    * type assertion
        * check whether the concrete type behind an interface value also implements another interface
        * useful for specifying optional interfaces
            * type might implement additional methods beyond those required by the primary interface
        * example: `if p, ok := a.(SomeInterface); ok { ... }`
            * without "comma-ok" idiom: panic
    * type switch
        * used when interface could be one of multiple possible types
        * one of the few places where shadowing is a good idea
        * example
            ```
            switch a := a.(type) { // shadowing
                case Dog:
                    fmt.Printf("This is a dog named %s and it says %s\n", a.Name, a.Sound())
                case Cat:
                    fmt.Printf("This is a cat named %s and it says %s\n", a.Name, a.Sound())
                default:
                    fmt.Printf("Unknown animal\n")
                }
            ```
* ok idiom
    * check if operation was successful without causing a panic or error
        ```
        if value, ok := m["foo"]; ok {
        ```

## pointers
* variable that holds the location in memory where a value is stored
    * zero value for a pointer is nil
    * example
        ```
        x := 1
        pointerToX := &x
        fmt.Println(pointerToX) // memory address
        fmt.Println(*pointerToX) // value
        ```
* before dereferencing => make sure that the pointer is non-nil
    * panic if you attempt to dereference a nil pointer
* pointer type: written with a * before a type
* primitive literal (numbers, booleans, and strings) or a constant don’t have memory addresses
    * when you need a pointer to a primitive type, declare a variable and point to it
* lack of immutable declarations in Go might seem problematic
    * ability to choose between value and pointer parameter types addresses the issue
* are a last resort
* use cases
    * if a struct is large enough
        * time to pass a pointer into a function is constant for all data sizes
    * returning a pointer versus returning a value
        * memory for the object must be allocated on the heap
        * data structures that are smaller than 10 megabytes => slower to return a pointer
    * code predating generics
        * no way to know what type of value to create and return
        * example: `json.Unmarshal([]byte(`{"name": "Bob", "age": 30}`), &f)`
        * disclaimer: when returning values from a function, favor value types
    * control over memory allocation
        * if function returned value, calling it in a loop => one value would be created on each loop iteration
    * indicate the difference between variable/field that hasn’t been assigned a value at all
        * vs assigning the zero value
        * disclaimer: be careful when using this pattern
            * pointers indicate mutability
* escape analysis
    * when the compiler determines that the data can’t be stored on the stack compiler stores the data on the heap
        * called: data escapes the stack
    * isn’t perfect
        * sometimes data that could be stored on the stack escapes to the heap
* generics
    * since go 1.18
    * implemented using type parameters
    * type constraints
        * example
            ```
            type Number interface {
            	int | int8 | int16 | int32 | int64 | float32 | float64
            }

            func Add[T Number](a, b T) T {
            	return a + b
            }
            ```
        * allowed operators are the ones that are valid for all of the listed types
        * `~` denotes type constraint to work with any type that has a specific underlying type
            * example
                ```
                type Integer interface {
                	~int
                }

                type MyInt int

                func Multiply[T Integer](a, b T) T {
                	return a * b
                }

                func main() {
                	var x MyInt = 10
                	var y MyInt = 20
                	fmt.Println(Multiply(x, y))  // without ~ it is not working
                }
                ```
    * how Go compiler handles the generation of functions for different types
        * separate versions for each unique underlying type
            * example: generic function works on both `int` and `float64` => compiler generates
              two separate functions: one for `int` and another for `float64`
        * shared functions for pointer types
            * operates on `unsafe.Pointer`
                * has to perform extra checks to handle different pointer types correctly
                    * example: runtime check to determine the actual type of the pointer
            * example: generic function that takes a pointer => same generated function for `*int`, `*float64`, `*string`

## error
* is a built-in interface
    ```
    type error interface {
        Error() string
    }
    ```
* two ways to create an error from a string
    * `errors.New("...")`
    * `fmt.Errorf("%d ...", i)`
* sentinel errors = predefined errors signaling that processing cannot continue
    * convention
        * declared at the package level
        * naming that starts with "Err"
    * used to check with `errors.Is(err, ...)` for specific error conditions and handle them accordingly
* wrapping = add more context while preserving the original error
    * new message without wrapping: `fmt.Errorf("...: %v", err)`
        * can't retrieve the original error from the new one
    * new message with wrapping: `fmt.Errorf("...: %w", err)`
        * includes both the new context and the original error
        * unwrapping: `if originalErr := errors.Unwrap(wrappedErr); originalErr != nil { ... } `
    * custom error types: needs to implement the method `Unwrap() error`
    * multiple wrapped errors: `errors.Join(errs...)`
        * example: validating fields in a struct
        * custom error type: needs to implement the method `Unwrap() []error`
    * wrapped sentinel error
        * problem: cannot use `==`, type assertion or type switch to check for it
        * solution
            * `errors.Is(err, ...)`
                * checks if a given error is or wraps a specific target error
            * `errors.As(err, &...)`
                * checks if a given error is or matches a specific target error
                * second parameter anything other than a pointer to an error or a pointer to an interface => the method panics
* panic
    * as soon as a panic happens, the current function exits immediately
        * defers attached to the current function start running
    * starting with Go 1.21, a `panic(nil)` is identical to `panic(new(runtime.PanicNilError))`
    * unrecoverable error
        * built-in `recover` function
            * recommended in one situation
                * do not let panics escape the boundaries of your public API
            * way to capture a panic to provide a more graceful shutdown (or prevent shutdown at all)
            * called from within a defer to check whether a panic happened
                * once a panic happens, only deferred functions are run
                * example
                    ```
                    defer func() {
                        if v := recover(); v != nil {
                            ....
                        }
                    }()
                    ```

## code organisation
* module = bundle of Go source code distributed and versioned as a single unit
    * consists of one or more packages (directories of source code)
    * `require` section lists the dependencies
    * conventional name: `go.mod`
    * example
        ```
        module example.com/myproject

        go 1.17  // Minimum Go version required by the module

        // lists direct dependencies
        require (
            github.com/some/dependency v1.2.3
            github.com/another/dependency v2.0.0
            ...
        )

        // lists indirect dependencies
        require (
            github.com/some/dependency v1.2.3 // indirect
            ...
        )
        ```
* modules can be grouped into two categories
    * intended as a single application
        * root of the project = main package
        * code in the main package should be minimal
            * all logic in an internal directory
                * code in the main function invokes code within internal
            * forbids to create a module that have dependency on application
    * intended as libraries
        * root of the project should have a package name that matches the repository name
* Go always builds applications from source code into a single binary file
    * includes the source code of module and all dependencies
* module system uses the principle of minimal version selection
    * assumption: all minor and patch versions of a module must be backward compatible
        * otherwise: bug
    * selects the lowest version of each dependency that satisfies the version requirements declared
      across all `go.mod` files involved
* vendoring = keeping copies of dependencies inside their module
    * ensures that a module always builds with identical dependencies
    * can make building your code faster on CI/CD
* workspace
    * introduced in Go 1.18
    * useful for development and testing
        * allows to make changes to one module and see the effects in another
            * without having to publish or tag new versions
        * allows you to use local import paths for your modules
    * workspace file (`go.work`) is used to define a set of modules that you are working on together
        * you can replace the dependencies of a module with local versions of those dependencies

## tests
* same directory and the same package as the production code
    * able to access and test unexported functions and variables
* written in a file whose name ends with `_test.go`
* test case start with the word `Test` and take single parameter of type `*testing.T`
    * use `t.Error` or `t.Errorf` to report a test failure and continue with the test execution
        * example
            ```
            if got != want {
                t.Error("Test failed: got", got, "want", want)
            }
            ```
    * use `t.Fatalf` to report a test failure and stop further execution of the test
* executing code before/after all tests: `TestMain` function
    1. `go test` calls it instead of the test functions
        * `go test ./...` run all tests
    1. `TestMain` function calls the `Run` method on `*testing.M`
        * runs the test functions in the package
        * `Run` method returns the exit code
            * 0 indicates that all tests passed
    1. `TestMain` function must call `os.Exit` with the exit code returned from `Run`
* temporary files: `TempDir` method on `*testing.T`
    * creates a new temporary directory every time it is invoked
    * returns the full path of the directory
    * registers a handler with Cleanup to delete the directory and its contents when the test exits
* setting environment variable for particular test: `t.Setenv()`
    * calls Cleanup to revert the environment variable to its previous state when the test exits
* sample data: subdirectory named `testdata`
    * used to keep fixtures, input files, and other data necessary for running tests
    * each package accesses its own `testdata` via a relative file path
        ```
        func TestReadFile(t *testing.T) {
            path := filepath.Join("testdata", "input.txt")
            ...
        }
        ```
* testing public API of the package: `packagename_test`
    * same directory as the production source code
* by default, unit tests are run sequentially
    * `t.Parallel()` makes tests run concurrently with other tests marked as parallel
* table-driven tests
    * allows you to test multiple scenarios with a single test function
    * example: define a table (slice) of test cases and iterate over them within a single test function
        ```
        type Person struct {
            Name    string
            Age     int
            Address string
        }

        func (p *Person) IsAdult() bool {
            return p.Age >= 18
        }

        func TestIsAdult(t *testing.T) {
            tests := []struct {
                name     string
                person   Person
                expected bool
            }{
                {
                    name:     "Adult",
                    person:   Person{Name: "Alice", Age: 20, Address: "123 Main St"},
                    expected: true,
                },
                {
                    name:     "Minor",
                    person:   Person{Name: "Bob", Age: 17, Address: "456 Elm St"},
                    expected: false,
                },
                {
                    name:     "Edge case - Exactly 18",
                    person:   Person{Name: "Charlie", Age: 18, Address: "789 Oak St"},
                    expected: true,
                },
                {
                    name:     "Negative Age",
                    person:   Person{Name: "Dave", Age: -1, Address: "000 Zero St"},
                    expected: false,
                },
            }

            for _, tt := range tests {
                t.Run(tt.name, func(t *testing.T) {
                    result := tt.person.IsAdult()
                    if result != tt.expected {
                        t.Errorf("isAdult(%+v) = %v; want %v", tt.person, result, tt.expected)
                    }
                })
            }
        }
        ```
* generating random data: fuzzing
    * can handle various types of input, such as integers, floats, structs, and more complex data types
    * `f.Add`: seeds the fuzzer with initial inputs
        * inputs are used as starting points for the fuzzing process
            * guides its exploration of the input space
        * no initial seeds => fuzzer starts with entirely random inputs
            * less likely to quickly hit edge cases or meaningful test scenarios
    * `f.Fuzz`: repeatedly called with different inputs generated by the fuzzer
    * example
        ```
        type Person struct {
            Name    string
            Age     int
            Address string
        }

        func (p *Person) IsAdult() bool {
            return p.Age >= 18
        }

        func FuzzIsAdult(f *testing.F) {
            f.Add("Alice", 20)
            f.Add("Bob", 17)
            f.Add("", 0)
            f.Add("Charlie", -1)

            f.Fuzz(func(t *testing.T, name string, age int, address string) {
                if !utf8.ValidString(name) {
                    t.Skip("invalid UTF-8 string")
                }

                person := &Person{
                    Name:    name,
                    Age:     age,
                    Address: address,
                }

                legal := person.IsAdult()

                expected := age >= 18

                if legal != expected {
                    t.Errorf("isAdult(%+v) = %v; want %v", person, legal, expected)
                }
            })
        }
        ```

## commands
* go install
    * compiles and installs packages and dependencies
* go generate
    * runs commands specified in specially formatted comments in the source code
    * example: protobufs
    * good idea to automate calling go generate before go build
* go test
    * allows you to specify which packages to test
* go get
    * fetches, compiles, and installs packages
    * by default doesn’t fetch code directly from source code repositories
        * sends requests to a proxy server run by Google and checks its cache
            * Google also maintains a checksum database
* go build
    * generates binary executables


## libs
* https://github.com/shopspring/decimal
    * arbitrary-precision fixed-point decimal numbers in Go
* https://pkg.go.dev/github.com/qfornaguera/goimports
    * same as gofmt + fixes imports
* https://github.com/dominikh/go-tools
    * staticcheck - The advanced Go linter
    * another linter: https://github.com/mgechev/revive
* https://pkg.go.dev/golang.org/x/tools/cmd/stringer
    * add `String()` method to enumeration’s values
* github.com/samber/lo
    * map, filter, contains, find...
* https://github.com/go-chi/chi
    * router for building Go HTTP services
* https://github.com/gorilla/handlers
    * middleware for Go HTTP
* https://github.com/google/go-cmp
    * package for comparing Go values in tests
    * returns a detailed description of what does not match (diffs)
* https://github.com/IBM/sarama
    * kafka client
* https://github.com/google/wire
    * compile-time dependency injection
* https://github.com/google/uuid
    * UUIDs based on RFC 4122 and DCE 1.1
* https://github.com/stretchr/testify
    * common assertions and mocks
* https://github.com/go-playground/validator
    * Go Struct and Field validation