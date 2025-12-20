data "external_schema" "gorm" {
  program = [
    "go",
    "run",
    "-mod=mod",
    "./loader",
  ]
}

env "dev" {
  src = data.external_schema.gorm.url
  dev = "mysql://root:@127.0.0.1:3306/test"
  url = "mysql://root:@127.0.0.1:3306/{{.ProjectName}}"
  migration {
    dir = "file://migrations"
  }
  format {
    migrate {
      diff = "{{ "{{" }} sql . \"  \" {{ "}}" }}"
    }
  }
}

env "production" {
  src = data.external_schema.gorm.url
  url = env("DATABASE_URL")
}
