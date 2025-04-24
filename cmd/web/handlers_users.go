package main

import (
	"fmt"
	"github/jordani-alpuche/test2/internal/data"
	"github/jordani-alpuche/test2/internal/validator"
	"net/http"
	"strconv"

	"github.com/gorilla/csrf"
)

/*******************************************************************************************************************************************************
* User Handlers                                                                                                                                                                                                                                                                                                                                                                                    *
********************************************************************************************************************************************************/

func (app *application) GETUsers(w http.ResponseWriter, r *http.Request) {
        data := NewTemplateData()
        data.CSRFField = csrf.TemplateField(r)
        data.CurrentPage = "/users" // Set CurrentPage here

        users, err := app.userInfo.GET(0)
        if err != nil {
                app.logger.Error("failed to fetch users", "error", err)
                http.Error(w, "Internal Server Error", http.StatusInternalServerError)
                return
        }

        data.User = users

        err = app.render(w,r, http.StatusOK, "userlist.tmpl", data)
        if err != nil {
                app.logger.Error("failed to render users page", "template", "userlist.tmpl", "error", err)
                http.Error(w, "Internal Server Error", http.StatusInternalServerError)
                return
        }
}

func (app *application) createUserForm(w http.ResponseWriter, r *http.Request) {
        data := NewTemplateData()
        data.CSRFField = csrf.TemplateField(r)
        data.CurrentPage = "/user" // Set CurrentPage here
        data.CurrentPageType="create"

        err := app.render(w,r, http.StatusOK, "addupdateuser.tmpl", data)
        if err != nil {
                app.logger.Error("failed to render add user page", "template", "addupdateuser.tmpl", "error", err)
                http.Error(w, "Internal Server Error", http.StatusInternalServerError)
                return
        }
}

func (app *application) createUser(w http.ResponseWriter, r *http.Request) {
        err := r.ParseForm()
        if err != nil {
                app.logger.Error("failed to parse form", "error", err)
                http.Error(w, "Bad Request", http.StatusBadRequest)
                return
        }

        firstName := r.PostForm.Get("FirstName")
        lastName := r.PostForm.Get("LastName")
        username := r.PostForm.Get("Username")
        password := r.PostForm.Get("Password")
        email := r.PostForm.Get("Email")
        phoneNumber := r.PostForm.Get("PhoneNumber")
        role := r.PostForm.Get("Role")
        status := r.PostForm.Get("Status")

        //hashpassword
        hashedPassword,err := validator.HashPassword(password)
        if err != nil {
                app.logger.Error("failed to hash password", "error", err)
                http.Error(w, "Internal Server Error", http.StatusInternalServerError)
                return
        }
        // check if username already exists
        existingUser, err := app.userInfo.FindByUsername(username,"users")
        if err != nil {
                app.logger.Error("failed to check if username exists", "error", err)
                http.Error(w, "Internal Server Error", http.StatusInternalServerError)
                return
        }
        if existingUser != nil {
                v := validator.NewValidator()
                v.AddError("Username", "Username already exists")
                data := NewTemplateData()
                data.CSRFField = csrf.TemplateField(r)
                data.CurrentPage = "/user" // Set CurrentPage here
                data.CurrentPageType="create"
                data.PasswordFieldName = "Password"
                data.FormErrors = v.Errors
                data.FormData = map[string]string{
                        "FirstName":   firstName,
                        "LastName":    lastName,
                        "Username":    username,
                        "Password":    password,
                        "Email":       email,
                        "PhoneNumber": phoneNumber,
                        "Role":        role,
                        "Status":      status,
                }
                err := app.render(w,r, http.StatusUnprocessableEntity, "addupdateuser.tmpl", data)
                if err != nil {
                        app.logger.Error("failed to render user Form", "template", "addupdateuser.tmpl", "error", err, "url", r.URL.Path, "method", r.Method)
                        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
                        return
                }
                return
        }

        user := &data.UsersData{
                FirstName:   firstName,
                LastName:    lastName,
                Username:    username,
                Password:    string(hashedPassword),
                Email:       email,
                PhoneNumber: phoneNumber,
                Role:        role,
                Status:      status,
        }

        // validate data
        v := validator.NewValidator()
        data.ValidateUsers(v, user)
        //
        // Check for validation errors
        if !v.ValidData() {
                data := NewTemplateData()
                data.CSRFField = csrf.TemplateField(r)
                data.CurrentPage = "/user" // Set CurrentPage here
                data.CurrentPageType="create"
                data.PasswordFieldName = "Password"
                data.FormErrors = v.Errors
                data.FormData = map[string]string{
                        "FirstName":   firstName,
                        "LastName":    lastName,
                        "Username":    username,
                        "Password":    password,
                        "Email":       email,
                        "PhoneNumber": phoneNumber,
                        "Role":        role,
                        "Status":      status,
                }

                err := app.render(w,r, http.StatusUnprocessableEntity, "addupdateuser.tmpl", data)
                if err != nil {
                        app.logger.Error("failed to render user Form", "template", "addupdateuser.tmpl", "error", err, "url", r.URL.Path, "method", r.Method)
                        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
                        return
                }
                return
        }

        err = app.userInfo.POST(user)

        if err != nil {
                app.logger.Error("failed to insert user", "error", err)
                http.Error(w, "Internal Server Error", http.StatusInternalServerError)
                return
        }
        http.Redirect(w, r, "/users", http.StatusSeeOther)
}

func (app *application) UserItem(w http.ResponseWriter, r *http.Request) {
        idStr := r.PathValue("id") // 👈 get ID from URL path like /brand-item/25
        id, err := strconv.Atoi(idStr)
        if err != nil || id < 1 {
                http.NotFound(w, r)
                return
        }

        users, err := app.userInfo.GET(id)

        if err != nil || len(users) == 0 {
                app.logger.Error("user not found", "id", id, "error", err)

                data := NewTemplateData()
                data.CSRFField = csrf.TemplateField(r)
                data.FormData = map[string]string{
                        "Message": "The user you're looking for doesn't exist.",
                }

                // Render custom 404 page
                err = app.render(w,r, http.StatusNotFound, "error-404.tmpl", data)
                if err != nil {
                        app.logger.Error("failed to render 404 page", "error", err)
                        http.Error(w, "Page not found", http.StatusNotFound)
                }
                return
        }

        user := users[0]

        data := NewTemplateData()
        data.CSRFField = csrf.TemplateField(r)
		data.CurrentPage =  "/users"
                data.CurrentPageType="update"
        data.FormData = map[string]string{
                "FirstName":        user.FirstName,
                "LastName":         user.LastName,
                "Username":         user.Username,
                "Password":         user.Password,
                "Email":            user.Email,
                "PhoneNumber":      user.PhoneNumber,
                "Role":             user.Role,
                "Status":           user.Status,                
        }

        err = app.render(w,r, http.StatusOK, "user-details.tmpl", data)
        if err != nil {
                app.logger.Error("failed to render viewer", "error", err)
                http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        }
}

func (app *application) updateUserForm(w http.ResponseWriter, r *http.Request) {
        idStr := r.PathValue("id")
        id, err := strconv.Atoi(idStr)
        if err != nil || id < 1 {
                http.NotFound(w, r)
                return
        }

        users, err := app.userInfo.GET(id)

        if err != nil || len(users) == 0 {
                app.logger.Error("users not found", "id", id, "error", err)

                data := NewTemplateData()
                data.CSRFField = csrf.TemplateField(r)
                data.FormData = map[string]string{
                        "Message": "The User you're looking for doesn't exist.",
                }

                err = app.render(w,r, http.StatusNotFound, "error-404.tmpl", data)
                if err != nil {
                        app.logger.Error("failed to render 404 page", "error", err)
                        http.Error(w, "Page not found", http.StatusNotFound)
                }
                return
        }

        user := users[0]
        data := NewTemplateData()
        data.CSRFField = csrf.TemplateField(r)
        data.CurrentPage = "/users" //Set CurrentPage here.
        data.CurrentPageType="update"
        data.PasswordFieldName = "NewPassword"
        data.FormData = map[string]string{
                "ID":               strconv.FormatInt(user.ID, 10),
                "FirstName":        user.FirstName,
                "LastName":         user.LastName,
                "Username":         user.Username,
                // "Password":         user.Password,
                "Email":            user.Email,
                "PhoneNumber":      user.PhoneNumber,
                "Role":             user.Role,
                "Status":           user.Status,
        }

        err = app.render(w,r, http.StatusOK, "addupdateuser.tmpl", data)
        if err != nil {
                app.logger.Error("failed to render viewer", "error", err)
                http.Error(w, "Internal Server Error", http.StatusInternalServerError)
                return
        }
}

func (app *application) updateUser(w http.ResponseWriter, r *http.Request) {
        idStr := r.PathValue("id")
        id, err := strconv.Atoi(idStr)
        if err != nil || id < 1 {
            http.NotFound(w, r)
            return
        }
    
        err = r.ParseForm()
        if err != nil {
            app.logger.Error("failed to parse form", "error", err, "url", r.URL.Path, "method", r.Method)
            http.Error(w, "Bad Request", http.StatusBadRequest)
            return
        }

        existingUsers, err := app.userInfo.GET(id)
        if err != nil || len(existingUsers) == 0 {
            app.logger.Error("failed to fetch existing user", "error", err)
            http.Error(w, "Internal Server Error", http.StatusInternalServerError)
            return
        }
        existingUser := existingUsers[0]

        fmt.Printf("\n existing userPassword: %v", existingUser.Password)
    
        firstName := r.PostForm.Get("FirstName")
        lastName := r.PostForm.Get("LastName")
        username := r.PostForm.Get("Username")
        newPassword := r.PostForm.Get("NewPassword")
        email := r.PostForm.Get("Email")
        phoneNumber := r.PostForm.Get("PhoneNumber")
        role := r.PostForm.Get("Role")
        status := r.PostForm.Get("Status")

        if newPassword == "" {
                newPassword = existingUser.Password
        }else{
                if err := validator.ValidatePassword(newPassword); 
                err != nil {
                        v := validator.NewValidator()
                                v.AddError("NewPassword", "Password must be at least 8 characters long")
                        data := NewTemplateData()
                        data.CSRFField = csrf.TemplateField(r)
                        data.FormErrors = v.Errors
                        data.FormData = map[string]string{
                                "ID":          idStr,
                                "FirstName":   firstName,
                                "LastName":    lastName,
                                "Username":    username,
                                "Password":    "", // don't prefill password field
                                "Email":       email,
                                "PhoneNumber": phoneNumber,
                                "Role":        role,
                                "Status":      status,
                        }
                        data.CurrentPage = "/users" //Set CurrentPage here.
                        data.CurrentPageType="update"
                        data.PasswordFieldName = "NewPassword"

                        err := app.render(w,r, http.StatusUnprocessableEntity, "addupdateuser.tmpl", data)
                        if err != nil {
                                app.logger.Error("failed to render Users Form", "template", "addupdateuser.tmpl", "error", err, "url", r.URL.Path, "method", r.Method)
                                http.Error(w, "Internal Server Error", http.StatusInternalServerError)
                        }
                        return
                }
                // Hash the new password
                hashedPassword,err:= validator.HashPassword(newPassword)
                if err != nil {
                        app.logger.Error("failed to hash password", "error", err)
                        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
                        return
                }
                newPassword = string(hashedPassword)   
        }

        if existingUser.Username != username {
          // check if username already exists
          findExistingUser, err := app.userInfo.FindByUsername(username,"users")
          if err != nil {
                  app.logger.Error("failed to check if username exists", "error", err)
                  http.Error(w, "Internal Server Error", http.StatusInternalServerError)
                  return
          }
          if findExistingUser != nil {
                  v := validator.NewValidator()
                  v.AddError("Username", "Username already exists")
                  data := NewTemplateData()
                  data.CSRFField = csrf.TemplateField(r)
                  data.CurrentPage = "/users" // Set CurrentPage here
                  data.CurrentPageType="update"
                  data.PasswordFieldName = "NewPassword"
                  data.FormErrors = v.Errors
                  data.FormData = map[string]string{
                                "ID":          idStr,
                          "FirstName":   firstName,
                          "LastName":    lastName,
                          "Username":    username,
                          "Password":    "",
                          "Email":       email,
                          "PhoneNumber": phoneNumber,
                          "Role":        role,
                          "Status":      status,
                  }
                  err := app.render(w,r, http.StatusUnprocessableEntity, "addupdateuser.tmpl", data)
                  if err != nil {
                          app.logger.Error("failed to render user Form", "template", "addupdateuser.tmpl", "error", err, "url", r.URL.Path, "method", r.Method)
                          http.Error(w, "Internal Server Error", http.StatusInternalServerError)
                          return
                  }
                  return
          }
        }

        fmt.Printf("\n form password: %s", newPassword)
       
       
        user := &data.UsersData{
            FirstName:   firstName,
            LastName:    lastName,
            Username:    username,
            Password:    newPassword, 
            Email:       email,
            PhoneNumber: phoneNumber,
            Role:        role,
            Status:      status,
        }
    
        v := validator.NewValidator()
        data.ValidateUsers(v, user)
    
        if !v.ValidData() {
            data := NewTemplateData()
            data.CSRFField = csrf.TemplateField(r)
            data.FormErrors = v.Errors
            data.FormData = map[string]string{
                "ID":               idStr,
                "FirstName":        firstName,
                "LastName":         lastName,
                "Username":         username,
                "Password":         newPassword,
                "Email":            email,
                "PhoneNumber":      phoneNumber,
                "Role":             role,
                "Status":           status,
            }
            data.CurrentPage = "/users" //Set CurrentPage here.
            data.CurrentPageType="update"
            data.PasswordFieldName = "NewPassword"
    
            err := app.render(w,r, http.StatusUnprocessableEntity, "addupdateuser.tmpl", data)
            if err != nil {
                app.logger.Error("failed to render Users Form", "template", "addupdateuser.tmpl", "error", err, "url", r.URL.Path, "method", r.Method)
                http.Error(w, "Internal Server Error", http.StatusInternalServerError)
            }
            return
        } 

 
        fmt.Printf("\nPassword to be saved: %s", user.Password)
        fmt.Print("\n\n")
        err = app.userInfo.PUT(id, user)
        if err != nil {
            app.logger.Error("failed to update user", "error", err)
            http.Error(w, "Internal Server Error", http.StatusInternalServerError)
            return
        }
        http.Redirect(w, r, "/users", http.StatusSeeOther)
    }
    

func (app *application) deleteUser(w http.ResponseWriter, r *http.Request) {
        idStr := r.PathValue("id") // 👈 get ID from URL path like /brand-item/25
        id, err := strconv.Atoi(idStr)
        if err != nil || id < 1 {
                http.NotFound(w, r)
                return
        }

        err = app.userInfo.DELETE(id)
        if err != nil {
                app.logger.Error("failed to delete user", "error", err)
                http.Error(w, "Internal Server Error", http.StatusInternalServerError)
                return
        }
        http.Redirect(w, r, "/users", http.StatusSeeOther)
}