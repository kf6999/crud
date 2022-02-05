package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"kennethfan.net/snippetbox/pkg/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
    // Because Pat matches the "/" path exactly, we can now remove the manual check
    // of r.URL.Path != "/" from this handler.

    s, err := app.snippets.Latest()
    if err != nil {
        app.serverError(w, err)
        return
    }

    app.render(w, r, "home.page.tmpl", &templateData{
        Snippets: s,
    })
}

func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
    // Pat doesn't strip the colon from the named capture key, so we need to
    // get the value of ":id" from the query string instead of "id".
    id, err := strconv.Atoi(r.URL.Query().Get(":id"))
    if err != nil || id < 1 {
        app.notFound(w)
        return
    }

    s, err := app.snippets.Get(id)
    if err != nil {
        if errors.Is(err, models.ErrNoRecord) {
            app.notFound(w)
        } else {
            app.serverError(w, err)
        }
        return
    }

    app.render(w, r, "show.page.tmpl", &templateData{
        Snippet: s,
    })
}

// Add a new createSnippetForm handler, which for now returns a placeholder response.
func (app *application) createSnippetForm(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Create a new snippet..."))
}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
    // Checking if the request method is a POST is now superfluous and can be
    // removed.

    title := "O snail"
    content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
    expires := "7"

    id, err := app.snippets.Insert(title, content, expires)
    if err != nil {
        app.serverError(w, err)
        return
    }
    // Change the redirect to use the new semantic URL style of /snippet/:id
    http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
}
