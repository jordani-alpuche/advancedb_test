package main

import (
	"fmt"
	"github/jordani-alpuche/test2/internal/data"
	"github/jordani-alpuche/test2/internal/validator"
	"html/template"
	"net/http"
	"strings"

	"github.com/gorilla/csrf"
	"golang.org/x/crypto/bcrypt"
)


func (app *application) LoginForm(w http.ResponseWriter, r *http.Request) {

    // Check if the user is already authenticated
    session, err := app.sessionStore.Get(r, "session-name")
    if err != nil {
        app.logger.Error("error getting session", "error", err)
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }
    isAuthenticated, ok := session.Values["authenticated"].(bool)
    if ok && isAuthenticated {
        // User is already authenticated, redirect to the homepage or another page
        http.Redirect(w, r, "/", http.StatusSeeOther)
        return
    }

    // If not authenticated, render the login form
	data := NewTemplateData()
	
	data.CSRFField = csrf.TemplateField(r)

	session, _ = app.sessionStore.Get(r, "signup-data")
	if session.Values["username"] != nil {
		// Concatenate the strings as usual
		alertMessage := "Sign up was successful (" +
			"Username: " + session.Values["username"].(string) + " , "+	
			"Password: " + session.Values["password"].(string) +
			")"
		
		// Assign the string to AlertMessage
		data.AlertMessage = alertMessage
		data.AlertType = "alert-success"
	
		session.Options.MaxAge = -1 // Clear session data after use
		session.Save(r, w)
	}
	
	fmt.Println("CSRF token:", csrf.Token(r))
	
	err = app.render(w, r, http.StatusOK, "signin.tmpl", data)
	if err != nil {
			app.logger.Error("failed to render login page", "template", "signin.tmpl", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
	}
}

func (app *application) login(w http.ResponseWriter, r *http.Request) {


	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	username := r.Form.Get("Username")
	password := r.Form.Get("Password")

	userData := &data.UsersData{
		Username: username,
		Password: password,
	}

	v := validator.NewValidator()
	data.ValidateLogin(v, userData)
	

	if !v.ValidData() {
		data := NewTemplateData()
		data.CSRFField = template.HTML(csrf.TemplateField(r)) 
		
		data.FormErrors = v.Errors
		data.FormData = map[string]string{
			"Username": username,
		}
		data.AlertMessage = "Please correct the errors below."
		data.AlertType = "alert-warning"

		err := app.render(w, r, http.StatusUnprocessableEntity, "signin.tmpl", data)
		if err != nil {
			app.logger.Error("failed to render user Form with validation errors", "template", "signin.tmpl", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		return // Stop processing
	}


	users, err := app.userInfo.FindByUsername(username, "login")

	if err != nil || len(users) == 0 {
		data := NewTemplateData()
		data.CSRFField = template.HTML(csrf.TemplateField(r)) 
		
		data.AlertMessage = "Invalid username or password."
		data.AlertType = "alert-danger"
		data.FormData = map[string]string{
			"Username": username,
		}
		app.render(w, r, http.StatusUnauthorized, "signin.tmpl", data)
		return
	}

	user := users[0]

	errf := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if errf != nil {
		if errf == bcrypt.ErrMismatchedHashAndPassword {
			data := NewTemplateData()
			data.CSRFField = template.HTML(csrf.TemplateField(r)) 
			
			data.AlertMessage = "Invalid username or password."
			data.AlertType = "alert-danger"
			data.FormData = map[string]string{
				"Username": username,
			}
			app.render(w, r, http.StatusUnauthorized, "signin.tmpl", data)
		} else {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return // Stop processing
	}

	// --- Login Successful ---
	// Get a session. We're using the default "session-name"
	session, err := app.sessionStore.Get(r, "session-name")
	if err != nil {
		app.logger.Error("failed to get session", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Set user information in the session
	session.Values["authenticated"] = true
	session.Values["userID"] = user.ID // Assuming your user struct has an ID field
    session.Values["userRole"] = user.Role // Assuming your user struct has a Role field
    // IMPORTANT: Remove or disable the Secure flag for non-HTTPS
    // session.Options.Secure = false // Set to false for local testing
	session.Options.MaxAge = 86400 // Session expires in 1 day
    // Set SameSite to Lax or Strict to prevent CSRF attacks
    // session.Options.SameSite = http.SameSiteLaxMode // Or http.SameSiteStrictMode used for https or set to none for https


	// Save the session
	session_err := session.Save(r, w)
	if session_err != nil {
		app.logger.Error("failed to save session", "error", session_err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}    

	// Redirect to the homepage
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) logout(w http.ResponseWriter, r *http.Request) {
	// Get the session
	session, err := app.sessionStore.Get(r, "session-name")
	if err != nil {
		// If there's no session, the user is already logged out,
		// or handle the error as appropriate for your application.
		http.Redirect(w, r, "/login", http.StatusSeeOther) // Redirect to login page
		return
	}

	// Clear session values
     // IMPORTANT: Remove or disable the Secure flag for non-HTTPS
    //  session.Options.Secure = false // Set to false for local testing
    //  // Set SameSite to Lax or Strict to prevent CSRF attacks
    //  session.Options.SameSite = http.SameSiteLaxMode // Or http.SameSiteStrictMode used for https or set to none for https
	// session.Values = make(map[interface{}]interface{})
	session.Options.MaxAge = -1 // Expire the cookie immediately


	// Save the empty/invalidated session
	err = session.Save(r, w)
	if err != nil {
		app.logger.Error("failed to save invalidated session", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Redirect the user to the login page
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}


func (app *application) signupForm(w http.ResponseWriter, r *http.Request) {
	data := NewTemplateData()
	data.CSRFField = csrf.TemplateField(r)
	data.CurrentPage = "/user" // Set CurrentPage here
	

	err := app.render(w,r, http.StatusOK, "signup.tmpl", data)
	if err != nil {
			app.logger.Error("failed to render signup page", "template", "signup.tmpl", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
	}
}

func (app *application) signupUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
			app.logger.Error("failed to parse form", "error", err)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
	}

	fmt.Printf("Form values: %v\n", r.PostForm)

	firstName := r.PostForm.Get("FirstName")
	lastName := r.PostForm.Get("LastName")
	password := r.PostForm.Get("Password")
	email := r.PostForm.Get("Email")
	phoneNumber := "0"
	role := "user"
	status := "active"	
	// generate username from firstname and last name

	user := &data.UsersData{
		FirstName:   firstName,
		LastName:    lastName,
		Username:    "username", // Placeholder, will be generated later
		Password:    password,
		Email:       email,
		PhoneNumber: phoneNumber,
		Role:        role,
		Status:      status,
}

	v := validator.NewValidator()
	data.ValidateSignup(v, user)

	if !v.ValidData() {
		data := NewTemplateData()
		data.CSRFField = csrf.TemplateField(r)
		data.AlertMessage = "Ensure all Fields are filled!"
			data.AlertType = "alert-warning"
		data.FormErrors = v.Errors
		data.FormData = map[string]string{
				"FirstName":   firstName,
				"LastName":    lastName,
				"Username":    "username", // Placeholder, will be generated later
				"Password":    password,
				"Email":       email,
				"PhoneNumber": phoneNumber,
				"Role":        role,
				"Status":      status,
		}

		err := app.render(w,r, http.StatusUnprocessableEntity, "signup.tmpl", data)
		if err != nil {
				app.logger.Error("failed to render user Form", "template", "signup.tmpl", "error", err, "url", r.URL.Path, "method", r.Method)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
		}
		return
}

		// Get first letter of first name in lowercase
		var firstInitial string
		if len(firstName) > 0 {
			firstInitial = strings.ToLower(string(firstName[0]))
		}
	
		// Get full last name in lowercase
		lastNameLower := strings.ToLower(lastName)
	
		// Concatenate result
		baseUsername := firstInitial + lastNameLower
		username := baseUsername

		fmt.Printf("username before pass: %s\n", username)

		// Retry logic to find a unique username
		counter := 1
		for {
			// Check if the username exists
			existingUser, err := app.userInfo.FindByUsername(username, "users")
			if err != nil {
				app.logger.Error("failed to check if username exists", "error", err, "username", username)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		
			// If no users are found, break out of the loop
			if len(existingUser) == 0 {
				break
			}
		
			// If the username exists, increment counter and try again
			username = fmt.Sprintf("%s%d", baseUsername, counter)
			counter++
		
			// Optional: Prevent an infinite loop after a max number of attempts
			if counter > 1000 { // Arbitrary max attempts to prevent infinite loop
				app.logger.Error("maximum attempts reached to generate a unique username")
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		}

	username = strings.ToLower(username)
	
	fmt.Printf("Final username: %s\n", username)

	//hashpassword
	hashedPassword,err := validator.HashPassword(password)
	if err != nil {
			app.logger.Error("failed to hash password", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
	}



	user = &data.UsersData{
			FirstName:   firstName,
			LastName:    lastName,
			Username:    username,
			Password:    string(hashedPassword),
			Email:       email,
			PhoneNumber: phoneNumber,
			Role:        role,
			Status:      status,
	}

	fmt.Printf("User to be created: %+v\n", user)

	err = app.userInfo.POST(user)

	if err != nil {
			app.logger.Error("failed to create user", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
	}

	// Store signup credentials in session
	session, _ := app.sessionStore.Get(r, "signup-data")
	session.Values["username"] = user.Username
	session.Values["password"] = password // plain text is okay temporarily
	session.Save(r, w)

	// Redirect to login page
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}