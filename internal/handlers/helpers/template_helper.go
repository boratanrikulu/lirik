package helpers

func GetTemplateFiles(controller_file string) []string {
	if controller_file == "" {
		panic("Controller file string can't be empty")
	}

	files := []string{
		controller_file,
		"./views/layouts/main_layout.html",
		"./views/partials/head_partial.html",
		"./views/partials/header_partial.html",
		"./views/partials/footer_partial.html",
		"./views/layouts/user_layout.html",
		"./views/partials/user/header_partial.html",
		"./views/partials/loading.html",
	}

	return files
}
