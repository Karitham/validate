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
    validate.MustRegisterValidatorFunc[ISBN](func (isbn ISBN) error {
        // For brevity, lets just implement one ISBN rule
        // It must be length of 10 or 13 
        if len(isbn) != 10 || len(isbn) != 13 {
            return fmt.Errorf("ISBN is of wrong length. Expected either 10 or 13 got %d", len(isbn))
        }
        return nil
    })

    err := validate.RegisterValidatorFunc[Book](func(b Book) error {
        // We can compose validators from other validators
        if err := validate.Validate(b.ISBN); err != nil {
            return err
        }
        
        if b.Author == "" {
            return errors.New("Author should not be empty")
        }

        return nil
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
    if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    if err := validate.Validate(book); err != nil {
        http.Error(w, fmt.Sprintf("book is not valid: %s", err.Error()), http.StatusBadRequest)
        return
    }
    
    // Continue on with handler
}
```
