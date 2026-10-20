package main
import (
	"fmt"
	"net/http"
	"strconv"
	//"html/template"
	"errors"
	"github.com/montenegrojds/snippetbox/internal/models"
)
/*
func (app *application) home(w http.ResponseWriter, r *http.Request){
	if r.URL.Path != "/"{
		app.notFound(w)
		return
	}
	
	files := []string{
		"./ui/html/base.html",
		"./ui/html/pages/home.html",
		"./ui/html/partials/nav.html",
	}

	ts,err:= template.ParseFiles(files...)
	if err != nil{
		app.serverError(w,err)
		return
	}

	err= ts.ExecuteTemplate(w,"base",nil)
	if err != nil{
		app.serverError(w,err)

	}
}*/

func (app *application) home(w http.ResponseWriter, r *http.Request){
	if r.URL.Path != "/"{
		app.notFound(w)
		return
	}
	snippets, err := app.snippets.Latest()

	if err != nil {
		app.serverError(w,err)
		return
	}

	for _, snippet := range snippets{
		fmt.Fprint(w, "%+v\n",snippet)
	}
}

func (app *application)snippetView(w http.ResponseWriter, r *http.Request){
	id,err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id<1{
		app.notFound(w)
		return
	}
	snippet, err := app.snippets.Get(id)
	if err != nil{
		if errors.Is(err,models.ErrNoRecord) {
			app.notFound(w)
		}else{
			app.serverError(w,err)
		}
		return
	}
	fmt.Fprintf(w, "%+v", snippet)
}



func (app *application)snippetCreate(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodPost{
		w.Header().Set("Allow",http.MethodPost)
		app.clientError(w,http.StatusMethodNotAllowed)
		return
	}
	title := "0 snail"
	content := "0 snail \nClimb Mount Fuji, \nBut sLowly slowly!\n\n- Kobayashi Issa"
	expires := 7

	id,err := app.snippets.Insert(title, content, expires)

	if err != nil{
		app.serverError(w,err)
		return
	}

	http.Redirect(w,r,fmt.Sprintf("/snippet/view?id=%d",id),http.StatusSeeOther)
	w.Write([]byte("Create a new snippet"))
}