package controllers

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"playzone/backend/middleware"
	"playzone/backend/repositories"
	"playzone/backend/services"
)

// ShowCategory liste les fils d'une categorie
func ShowCategory(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	cat, err := repositories.FindCategoryBySlug(slug)
	if err != nil {
		Render(w, r, "error", PageData{Title: "404", Error: "Categorie introuvable."})
		return
	}

	limit, offset, perPage := parsePagination(r)
	total := repositories.CountThreadsByCategory(slug)
	threads, _ := repositories.ListThreadsByCategory(slug, limit, offset)

	totalPages := 1
	if limit > 0 {
		totalPages = (total + limit - 1) / limit
		if totalPages == 0 {
			totalPages = 1
		}
	}

	currentPage, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if currentPage < 1 {
		currentPage = 1
	}

	Render(w, r, "category", PageData{
		Title: cat.Name,
		Data: map[string]interface{}{
			"Category":    cat,
			"Threads":     threads,
			"Total":       total,
			"PerPage":     perPage,
			"CurrentPage": currentPage,
			"TotalPages":  totalPages,
		},
	})
}

// ShowByTag liste les fils d'un tag
func ShowByTag(w http.ResponseWriter, r *http.Request) {
	tag := chi.URLParam(r, "name")
	threads, _ := repositories.ListThreadsByTag(tag)
	Render(w, r, "tag", PageData{
		Title: "Tag : " + tag,
		Data: map[string]interface{}{
			"Tag":     tag,
			"Threads": threads,
		},
	})
}

// NewThreadForm affiche le formulaire de creation
func NewThreadForm(w http.ResponseWriter, r *http.Request) {
	categories, _ := repositories.ListAllCategories()
	tags, _ := repositories.ListAllTags()
	Render(w, r, "thread_new", PageData{
		Title: "Nouveau fil",
		Data: map[string]interface{}{
			"Categories": categories,
			"Tags":       tags,
		},
	})
}

// CreateThread cree un fil
func CreateThread(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	title := r.FormValue("title")
	content := r.FormValue("content")
	catID, _ := strconv.Atoi(r.FormValue("category_id"))
	tagIDs := parseTagIDs(r.Form["tags"])

	id, err := services.CreateThread(title, content, middleware.GetUserID(r), catID, tagIDs)
	if err != nil {
		SetFlash(w, "error", err.Error())
		http.Redirect(w, r, "/threads/new", http.StatusFound)
		return
	}
	SetFlash(w, "success", "Fil cree !")
	http.Redirect(w, r, "/threads/"+strconv.FormatInt(id, 10), http.StatusFound)
}

// ShowThread affiche un fil + ses messages
func ShowThread(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.NotFound(w, r)
		return
	}

	thread, err := services.GetThread(id)
	if err != nil {
		Render(w, r, "error", PageData{Title: "Erreur", Error: err.Error()})
		return
	}

	// Tri et pagination
	sortBy := r.URL.Query().Get("sort")
	if sortBy == "" {
		sortBy = "date" // par defaut : du plus recent au plus ancien
	}
	limit, offset, perPage := parsePagination(r)
	total := repositories.CountMessagesByThread(id)
	messages, _ := repositories.ListMessagesByThread(id, sortBy, limit, offset, middleware.GetUserID(r))

	totalPages := 1
	if limit > 0 {
		totalPages = (total + limit - 1) / limit
		if totalPages == 0 {
			totalPages = 1
		}
	}
	currentPage, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if currentPage < 1 {
		currentPage = 1
	}

	Render(w, r, "thread_detail", PageData{
		Title: thread.Title,
		Data: map[string]interface{}{
			"Thread":      thread,
			"Messages":    messages,
			"Sort":        sortBy,
			"PerPage":     perPage,
			"Total":       total,
			"CurrentPage": currentPage,
			"TotalPages":  totalPages,
		},
	})
}

// EditThreadForm formulaire d'edition
func EditThreadForm(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	thread, err := repositories.FindThreadByID(id)
	if err != nil {
		Render(w, r, "error", PageData{Title: "404", Error: "Fil introuvable"})
		return
	}
	if thread.AuthorID != middleware.GetUserID(r) && !middleware.IsAdmin(r) {
		http.Error(w, "Acces refuse", http.StatusForbidden)
		return
	}
	tags, _ := repositories.ListAllTags()
	Render(w, r, "thread_edit", PageData{
		Title: "Modifier le fil",
		Data: map[string]interface{}{
			"Thread": thread,
			"Tags":   tags,
		},
	})
}

// UpdateThread modifie le fil
func UpdateThread(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	r.ParseForm()
	title := r.FormValue("title")
	content := r.FormValue("content")
	tagIDs := parseTagIDs(r.Form["tags"])

	err := services.UpdateThread(id, middleware.GetUserID(r), middleware.IsAdmin(r), title, content, tagIDs)
	if err != nil {
		SetFlash(w, "error", err.Error())
	} else {
		SetFlash(w, "success", "Fil modifie")
	}
	http.Redirect(w, r, "/threads/"+strconv.Itoa(id), http.StatusFound)
}

// DeleteThread supprime le fil
func DeleteThread(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	err := services.DeleteThread(id, middleware.GetUserID(r), middleware.IsAdmin(r))
	if err != nil {
		SetFlash(w, "error", err.Error())
		http.Redirect(w, r, "/threads/"+strconv.Itoa(id), http.StatusFound)
		return
	}
	SetFlash(w, "success", "Fil supprime")
	http.Redirect(w, r, "/", http.StatusFound)
}

// ChangeState change l'etat (ouvert/ferme/archive)
func ChangeState(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	r.ParseForm()
	state := r.FormValue("state")
	err := services.ChangeState(id, middleware.GetUserID(r), middleware.IsAdmin(r), state)
	if err != nil {
		SetFlash(w, "error", err.Error())
	} else {
		SetFlash(w, "success", "Etat change : "+state)
	}
	http.Redirect(w, r, "/threads/"+strconv.Itoa(id), http.StatusFound)
}

// ---- helpers pagination & tags ----

// parsePagination lit ?per_page=10|20|30|all&page=2 -> retourne limit, offset, perPage label
func parsePagination(r *http.Request) (int, int, string) {
	perPage := r.URL.Query().Get("per_page")
	if perPage == "" {
		perPage = "10"
	}

	limit := 10
	if perPage == "all" {
		limit = 0
	} else if v, err := strconv.Atoi(perPage); err == nil && (v == 10 || v == 20 || v == 30) {
		limit = v
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit
	return limit, offset, perPage
}

// parseTagIDs convertit les valeurs texte du formulaire (tags coches)
// en une liste d'identifiants entiers, en ignorant les valeurs invalides.
func parseTagIDs(raw []string) []int {
	var ids []int
	for _, s := range raw {
		if id, err := strconv.Atoi(s); err == nil {
			ids = append(ids, id)
		}
	}
	return ids
}
