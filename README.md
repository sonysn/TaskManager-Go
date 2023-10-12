# Task Manager

Working on the Task Manager Project Written in Golang by Stephen.  

## Notes on Go

- Functions are exported by beginning with a Capital letter

```Go
//This is a Private Function
func printName(){
    fmt.Println("Stephen")
}

//This is a Public Exported Function
func PrintName(){
    fmt.Println("Stephen")
}
```

- Variables and Constants are exported by beginning with a Capital letter

```Go
//This is a Private Variable and Constant
var name = "Stephen"
const name = "Stephen"

//This is a Public Exported Variable and Constant
var Name = "Stephen"
const Name = "Stephen"
```

- I import functions and parameters in two ways

```Go
//With the import keyword
import "fmt"

//Multiple Imports
import (
    "fmt"
    "net/http"
    "TaskManager/functions" //This is an example project directory
)

//Let's call the function PrintName() from the functions package
//I can do it in two ways

//In another file

//This would print Stephen as defined above
func AnotherFunction(){
    functions.PrintName()
}

//And this
import (
    "fmt"
    "net/http"
    P "TaskManager/functions" //This is an example project directory
)

//This would print Stephen as defined above
func AnotherFunction(){
    P.PrintName()
}

//The both do the same thing
```

- I can also declare and infer types

```Go
//Like
var name string

//or like
name := "Stephen"

//or from a function that returns a string for example
name := functions.PrintName()
```

- Functions can return values

```Go
func returnString() string{}

func returnInt() int{}

name := returnString() //name will be a string
age := returnInt() //age will be an integer
```

With that out of the way let's get started!

## Dependencies

- [Redis](https://github.com/go-redis/redis)
- [MongoDB ODM](https://github.com/kamva/mgm)
- [MongoDB Driver](https://github.com/mongodb/mongo-go-driver)

## Usage

Started out with Golang completed with NodeJS.
The usage and endpoints are the same so I'll be linking the NodeJS project here [TaskManager NodeJS](https://github.com/sonysn/TaskManager-Node.git)

Thanks for reading!
