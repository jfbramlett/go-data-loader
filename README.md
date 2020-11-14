[![Actions Status](https://github.com/jfbramlett/go-template/workflows/Go/badge.svg)](https://github.com/jfbramlett/go-template/actions)

# Go-Template
This project is a template project for new Go projects. The project generates a web server for serving JSON data.

Once you generate a new project from this template you will need to update the following files:

* README.md -> update the Action status reference to your action url
* go.mod -> update module name to reference your new project
* cmd/server/main.go -> update the Cobra.Command values
* cmd/server/server.go -> update the flag reference service_name

A script is included in the project that will make the necessary updates based on the new projects current working directory. You can run this with:

```
python3 update_template.py
```
