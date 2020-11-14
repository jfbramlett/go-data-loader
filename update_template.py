import os


def main():
    project_name = os.path.basename(os.getcwd())

    file_list = ['go.mod', 'cmd/server/main.go', 'cmd/server/server.go', 'README.md']

    for fpath in file_list:
        replace_in_file(fpath, project_name)


def replace_in_file(fpath, project_name):
    with open(fpath) as f:
        s = f.read()
    s = s.replace("go-template", project_name)
    with open(fpath, "w") as f:
        f.write(s)


if __name__ == "__main__":
    main()
