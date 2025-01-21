
# Go Rules Engine

A lightweight and extensible rules engine for Go, designed to support conditional execution of actions based on user-defined rules.

## Features

- **Custom Conditions:** Define your own conditions for evaluating rules.
- **Flexible Actions:** Execute actions when conditions are met.
- **Execution Modes:** Supports `FirstMatch`, `AllMatch`, `AnyMatch`, and `NoneMatch` execution strategies.
- **Rule Prioritization:** Rules are executed based on their priority.
- **Extensible Registry:** Centralized registration of conditions and actions.

## Installation

To use the rules engine in your project, add it as a dependency:

```bash
go get github.com/lvm/go-rules
```

## Usage

### Basic Setup

1. **Initialize the Registry**  
   The registry holds the conditions and actions that the rules engine will use.

   ```go
   registry := NewRegistry()
   ```

2. **Register Conditions and Actions**  
   Define conditions and actions to be used in rules.

   ```go
   registry.AddCondition("isEven", func(args Arguments) bool {
       n, _ := args["number"].(int)
       return n%2 == 0
   })

   registry.AddAction("printSuccess", func(args Arguments) error {
       fmt.Println("number is even!")
       return nil
   })
   ```

3. **Create Rules**  
   Define rules with conditions, arguments, and actions.

   ```go
   rule := When("isEven", Arguments{"number": 4}, 1)
   rule.Action = "printSuccess"
   ```

4. **Initialize the Rule Engine**  
   Set the execution mode and context for the engine.

   ```go
   ruleEngine := NewRuleEngine(AllMatch, func(msg string) {
       fmt.Println("Log:", msg)
   }, registry, context.TODO())
   ```

5. **Add Rules and Execute**  
   Add rules to the engine and execute them.

   ```go
   ruleEngine.AddRule(rule)

   if err := ruleEngine.Execute(Arguments{"number": 4}); err != nil {
       fmt.Println("Execution failed:", err)
   }
   ```

### Execution Modes

- **AllMatch:** Executes all matching rules.
- **AnyMatch:** Executes if any rule matches.
- **NoneMatch:** Executes if no rules match.

Example with different execution modes:

```go
ruleEngine := NewRuleEngine(FirstMatch, func(msg string) {
    fmt.Println("Log:", msg)
}, registry, context.TODO())
```

### Combining Conditions

Combine multiple conditions using `All`, `Any`, or `None`.

```go
isEven := func(args Arguments) bool { return args["number"].(int)%2 == 0 }
isPositive := func(args Arguments) bool { return args["number"].(int) > 0 }

registry.AddCondition("isEvenAndPositive", All(isEven, isPositive))
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
