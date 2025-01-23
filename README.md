# Go Rules Engine

A lightweight and extensible rules engine for Go, designed to support conditional execution of actions based on user-defined rules.

## Features

- **Custom Conditions:** Define your own conditions for evaluating rules.
- **Flexible Actions:** Execute actions when conditions are met.
- **Execution Modes:** Supports `AllMatch`, `AnyMatch`, and `NoneMatch` execution strategies.
- **Rule Prioritization:** Rules are executed based on their priority.
- **Extensible Registry:** Centralized registration of conditions and actions.

## Installation

To use the rules engine in your project, add it as a dependency:

```bash
go get github.com/lvm/go-rules
```

## Usage

### Basic Setup

1. **Initialize the Rule Engine**  
   The rule engine holds the conditions and actions that the rules engine will use.

   ```go
   ruleEngine := NewRuleEngine(context.TODO(), AllMatch, func(msg string) {
       fmt.Println("Log:", msg)
   })
   ```

2. **Create Conditions and Actions**  
   Define conditions and actions to be used in rules.

   ```go
   isEven := func(args Arguments) bool {
       n, _ := args["number"].(int)
       return n%2 == 0
   }

   printSuccess := func(args Arguments) error {
       fmt.Println("number is even!")
       return nil
   }

   rule := NewRule(Condition{Name: "isEven", Fn: isEven}, Action{Name: "printSuccess", Fn: printSuccess}, 1)
   ```

3. **Add Rules and Execute**  
   Add rules to the engine and execute them.

   ```go
   ruleEngine.AddRules(rule)

   if err := ruleEngine.Execute(Arguments{"number": 4}); err != nil {
       fmt.Println("Execution failed:", err)
   }
   ```

4. **Managing Rule Engines**

    Use the `Registry` to manage multiple rule engines.

    ```go
    registry := NewRegistry()
    registry.AddEngine("default", *ruleEngine)

    defaultEngine := registry.GetEngine("default")
    ```


### Execution Modes

- **AllMatch:** Executes all matching rules.
- **AnyMatch:** Executes if any rule matches.
- **NoneMatch:** Executes if no rules match.

Example with different execution modes:

```go
NewRuleEngine(context.TODO(), AllMatch, func(msg string) {})

NewRuleEngine(context.TODO(), AnyMatch, func(msg string) {})

NewRuleEngine(context.TODO(), NoneMatch, func(msg string) {})
```

### Combining Conditions

Combine multiple conditions using `All`, `Any`, or `None`.

```go
isEven := func(args Arguments) bool { return args["number"].(int)%2 == 0 }
isPositive := func(args Arguments) bool { return args["number"].(int) > 0 }

allConditions := All(isEven, isPositive)
anyConditions := Any(isEven, isPositive)
nonConditions := None(isEven, isPositive)
```


### Context Management in Rules

The rules engine allows you to store and retrieve values from the context, enabling dynamic behavior based on the context during rule evaluation.

1. **Setting and Getting Context Values**  
   You can store context values using `SetContext` and retrieve them with `GetContext`. This is useful for passing additional data that might influence rule execution.

2. **Using Context in Conditions and Actions**  
   Rules can use context values to modify their behavior. For instance, a condition might check if a specific context value exists and decide whether to execute an action or skip the rule entirely.


```go
ruleEngine := NewRuleEngine(context.TODO(), AllMatch, func(msg string) {})

isEvenCondition := Condition{
    Name: "isEven",
    Fn: func(c context.Context, args Arguments) bool {
        if c.Value("ForcePass") != nil {
            return true
        }

        n, _ := args["number"].(int)
        return n%2 == 0
    },
}

printAction := Action{
    Name: "printSuccess",
    Fn: func(c context.Context, args Arguments) error {
        n, _ := args["number"].(int)
        fmt.Printf("Success: %d is even!\n", n)
        return nil
    },
}

rule := NewRule(isEvenCondition, printAction, 1)

ruleEngine.AddRules(rule)

ruleEngine.Execute(Arguments{"number": 4}) // this should pass: 4 is even.

ruleEngine.SetContext("ForcePass", 1)
ruleEngine.Execute(Arguments{"number": 5}) // this too should pass: even though 5 is odd, ForcePass is not nil
```


### Example: Middleware with Gin

Integrate the rules engine into a Gin application as middleware:

```go
func RulesMiddleware(ruleEngine *RuleEngine) gin.HandlerFunc {
    return func(c *gin.Context) {
        args := Arguments{
            "path": c.Request.URL.Path,
            "method": c.Request.Method,
        }

        if err := ruleEngine.Execute(args); err != nil {
            c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
                "error": err.Error(),
            })
            return
        }

        c.Next()
    }
}
```

### Tests

The project includes a suite of tests for validating functionality. To run the tests:

```bash
go test -v ./...
```

## Contributing

Contributions are welcome! Please open an issue or submit a pull request for any bug fixes or feature requests.

## License

See [LICENSE](LICENSE).