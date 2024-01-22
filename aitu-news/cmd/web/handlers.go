package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"alexedwards.net/snippetbox/pkg/forms" // New import
	"alexedwards.net/snippetbox/pkg/models"
)

func home(w http.ResponseWriter, r *http.Request) {
	session, err := app.sessionStore.Get(r, "user-session")
	if checkError(err, &w, "Internal Server Error", http.StatusInternalServerError) {
		return
	}
	email := session.Values["userEmail"].(string)
	role := session.Values["userRole"].(string)
	user := &models.User{Email: email, Role: role}

	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	articles, err := app.articles.Latest()
	if err != nil {
		app.errorLog.Println(err)
	}
	data := &templateData{Articles: articles, User: user}

	files := []string{

		"./ui/html/home.page.tmpl",

		"./ui/html/base.layout.tmpl",

		"./ui/html/header.partial.tmpl",
		"./ui/html/footer.partial.tmpl",
	}

	ts, err := template.ParseFiles(files...)
	if checkError(err, &w, "Internal Server Error", 500) {
		return
	}

	err = ts.Execute(w, data)
	checkError(err, &w, "Internal Server Error", 500)
}

func article(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if checkError(err, &w, "Not found", 404) {
		return
	}

	artic, err := app.articles.Get(id)
	if err != nil {
		app.errorLog.Println(err)
	}
	data := &templateData{Articles: []*models.Article{artic}}

	files := []string{
		"./ui/html/article.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/header.partial.tmpl",
		"./ui/html/footer.partial.tmpl",
	}

	ts, err := template.ParseFiles(files...)
	if checkError(err, &w, "Internal Server Error", 500) {
		return
	}

	err = ts.Execute(w, data)
	checkError(err, &w, "Internal Server Error", 500)
}

func articleCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {

		err := r.ParseMultipartForm(10 << 20) // Максимальный размер файла 10 MB
		if err != nil {
			app.errorLog.Println(w, "File size greater than 10 mb. ", http.StatusInternalServerError)
			return
		}

		file, _, err := r.FormFile("image")
		if err != nil {
			app.errorLog.Println(w, "Cannot upload file", http.StatusInternalServerError)
			return
		}
		defer file.Close()

		title := r.FormValue("title")
		text := r.FormValue("text")
		categoryId := r.FormValue("category")

		insertedId, err := app.articles.Insert(title, text, categoryId)
		if checkError(err, &w, "Article not created.", 404) {
			return
		}
		filename := app.config.SERVER.StaticDir + "img/" + "articleImg" + strconv.Itoa(insertedId) + ".jpg"
		f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			app.errorLog.Println(w, "Error when saving image", http.StatusInternalServerError)
			return
		}
		defer f.Close()

		io.Copy(f, file)
		app.infoLog.Println(w, "Image %s successfully uploaded загружено!", filename)

	}
	http.Redirect(w, r, "/admin", 200)
}

func register(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			app.errorLog.Println(w, "ParseForm() err: %v", err)
			return
		}
		name := r.FormValue("name")
		lastname := r.FormValue("lastname")
		email := r.FormValue("email")
		password := r.FormValue("password")

		id, err := app.users.Insert(name, lastname, email, password)
		if checkError(err, &w, "User not created", 500) {
			return
		}
		app.infoLog.Printf("User with id: %d was created.\n", id)
	}
	if r.Method == "GET" {
		files := []string{
			"./ui/html/register.page.tmpl",
			"./ui/html/base.layout.tmpl",
			"./ui/html/header.partial.tmpl",
			"./ui/html/footer.partial.tmpl",
		}

		ts, err := template.ParseFiles(files...)
		if checkError(err, &w, "Internal Server Error", 500) {
			return
		}

		err = ts.Execute(w, nil)
		checkError(err, &w, "Internal Server Error", 500)
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		files := []string{
			"./ui/html/login.page.tmpl",
			"./ui/html/base.layout.tmpl",
			"./ui/html/header.partial.tmpl",
			"./ui/html/footer.partial.tmpl",
		}

		ts, err := template.ParseFiles(files...)
		if checkError(err, &w, "Internal Server Error", 500) {
			return
		}

		err = ts.Execute(w, nil)
		checkError(err, &w, "Internal Server Error", 500)
	}

	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			app.errorLog.Println(w, "ParseForm() err: %v", err)
			return
		}
		email := r.FormValue("email")
		password := r.FormValue("password")
		u, err := app.users.Get(email, password)
		if checkError(err, &w, "User not created", 500) {
			return
		}
		session, err := app.sessionStore.Get(r, "user-session")
		if checkError(err, &w, "Cannot create session", 500) {
			return
		}
		session.Values["userEmail"] = u.Email
		session.Values["userRole"] = u.Role
		err = session.Save(r, w)
		if err != nil {
			app.errorLog.Println(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func articles(w http.ResponseWriter, r *http.Request) {
	filter := 0
	if r.URL.Path == "/articles/students" {
		filter = 1
	} else if r.URL.Path == "/articles/staff" {
		filter = 2
	} else if r.URL.Path == "/articles/applicants" {
		filter = 3
	} else if r.URL.Path == "/articles/researchers" {
		filter = 4
	}
	articles, err := app.articles.GetAll(filter)
	if err != nil {
		app.errorLog.Println(err)
	}
	data := &templateData{Articles: articles}

	files := []string{

		"./ui/html/articles.page.tmpl",

		"./ui/html/base.layout.tmpl",

		"./ui/html/header.partial.tmpl",
		"./ui/html/footer.partial.tmpl",
	}

	ts, err := template.ParseFiles(files...)
	if checkError(err, &w, "Internal Server Error", 500) {
		return
	}

	err = ts.Execute(w, data)
	checkError(err, &w, "Internal Server Error", 500)
}

func contacts(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from /contacts"))
}

func adminPanel(w http.ResponseWriter, r *http.Request) {

	files := []string{

		"./ui/html/admin.page.tmpl",

		"./ui/html/base.layout.tmpl",

		"./ui/html/header.partial.tmpl",
		"./ui/html/footer.partial.tmpl",
	}

	ts, err := template.ParseFiles(files...)
	if checkError(err, &w, "Internal Server Error", 500) {
		return
	}

	err = ts.Execute(w, nil)
	checkError(err, &w, "Internal Server Error", 500)
}

func checkError(err error, w *http.ResponseWriter, msg string, status int) bool {
	if err != nil {
		app.errorLog.Println(err.Error())
		http.Error(*w, msg, status)
	}
	return err != nil
}

func (app *application) createSnippetForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "create.page.tmpl", &templateData{
		// Pass a new empty forms.Form object to the template.
		Form: forms.New(nil),
	})
}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Create a new forms.Form struct containing the POSTed data from the
	// form, then use the validation methods to check the content.
	form := forms.New(r.PostForm)
	form.Required("title", "content", "expires")
	form.MaxLength("title", 100)
	form.PermittedValues("expires", "365", "7", "1")

	// If the form isn't valid, redisplay the template passing in the
	// form.Form object as the data.
	if !form.Valid() {
		app.render(w, r, "create.page.tmpl", &templateData{Form: form})
		return
	}

	// Because the form data (with type url.Values) has been anonymously embedded
	// in the form.Form struct, we can use the Get() method to retrieve
	// the validated value for a particular form field.
	id, err := app.snippets.Insert(form.Get("title"), form.Get("content"), form.Get("expires"))
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
}