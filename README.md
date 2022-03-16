# Validate

A simple library to register a validation for a type once, then validate that type anywhere.

*This library is just for playing around with generics, should not actually be used.*

## Example Usage

```go
type ISBN string

type Book struct {
    ISBN ISBN
    Author string
    // Other fields ommited for clarity
}


func main() {
    validate.MustRegisterValidatorFunc[ISBN](func (isbn ISBN)) bool {
        // For brevity, lets just implement one ISBN rule
        // It must be length of 10 or 13 
        return len(isbn) == 10 || len(isbn) == 13
    })

    err := validate.RegisterValidatorFunc[Book](func(b Book) bool {
        // We can compose validators easily
        return validate.Validate(b.ISBN) && b.Author != ""
    })

    if err != nil {
        panic(err)
    }

    // setup rest of application
}

// Post handler to add a new Book
func AddBook(w http.ResponseWriter, r *http.Request) {

    defer r.Body.Close()
    
    var book Book
    err := json.NewDecoder(r.Body).Decode(&book)

    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    ok := validate.Validate(book)

    if !ok {
        http.Error(w, "book is not valid", http.StatusBadRequest)
        return
    }
    
    // Continue on with handler
}
```

